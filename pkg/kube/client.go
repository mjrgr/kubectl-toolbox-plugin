package kube

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
	"kubectl-toolbox-plugin/pkg/utils"
	"os"
	"os/signal"
	"syscall"
)

type Client struct {
	Factory    cmdutil.Factory
	Namespace  string
	kubeClient *kubernetes.Clientset
}

// New creates a new Client.
func New(getter genericclioptions.RESTClientGetter) *Client {
	if getter == nil {
		getter = genericclioptions.NewConfigFlags(true)
	}
	return &Client{
		Factory: cmdutil.NewFactory(getter),
	}
}

// GetKubeClient returns Kubernetes client
func (c *Client) GetKubeClient() *kubernetes.Clientset {
	if c.kubeClient == nil {
		cli, err := c.Factory.KubernetesClientSet()
		if err != nil {
			logrus.WithError(err).Panic("Unable to get Kubernetes client")
		} else {
			c.kubeClient = cli
		}
	}

	return c.kubeClient
}

// GetNamespace returns default namespace used
func (c *Client) GetNamespace() string {
	if c.Namespace != "" {
		return c.Namespace
	}
	if ns, _, err := c.Factory.ToRawKubeConfigLoader().Namespace(); err == nil {
		return ns
	}
	return v1.NamespaceDefault
}

// ServerVersion returns Kubernetes cluster version
func (c *Client) ServerVersion() (*version.Info, error) {
	return c.GetKubeClient().ServerVersion()
}

func generateName() string {
	return utils.CommandName + "-" + rand.String(7)
}

func (c *Client) RunRemotePod(image string, command string) error {
	// Create a Kubernetes core/v1 client.
	// Create a busybox Pod.  By running `cat`, the Pod will sit and do nothing.
	customName := generateName()
	pod, err := c.GetKubeClient().CoreV1().Pods(c.GetNamespace()).Create(context.Background(), &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: customName,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    customName,
					Image:   image,
					Command: []string{"cat"},
					Stdin:   true,
				},
			},
			TerminationGracePeriodSeconds: &utils.Zero,
		},
	}, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	deletePod := func() {
		c.GetKubeClient().CoreV1().Pods(c.GetNamespace()).Delete(context.Background(), pod.Name, metav1.DeleteOptions{
			GracePeriodSeconds: &utils.Zero,
		})
	}

	// Delete the Pod before we exit.
	defer deletePod()

	// Delete the Pod before on interrupt signal (Ctrl+C)
	cn := make(chan os.Signal)
	signal.Notify(cn, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cn
		deletePod()
	}()

	// Wait for the Pod to indicate Ready state
	watcher, err := c.GetKubeClient().CoreV1().Pods(c.GetNamespace()).Watch(context.Background(),
		metav1.SingleObject(pod.ObjectMeta),
	)

	if err != nil {
		return err
	}

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Modified:
			pod = event.Object.(*v1.Pod)

			// If the Pod contains a status condition Ready == True, stop
			// watching.
			for _, cond := range pod.Status.Conditions {
				if cond.Type == v1.PodReady &&
					cond.Status == v1.ConditionTrue {
					watcher.Stop()
				}
			}
		default:
			return errors.New(fmt.Sprintf("Unexpected event type %s", event.Type))
		}
	}

	// Prepare the API URL used to execute another process within the Pod.  In
	// this case, we'll run a remote shell.
	req := c.GetKubeClient().CoreV1().RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: pod.Spec.Containers[0].Name,
			Command:   []string{command},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	config, _ := c.Factory.ToRESTConfig()
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	// Connect this process' std{in,out,err} to the remote shell process.
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		return err
	}
	return nil
}
