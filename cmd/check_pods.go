package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/multierr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubectl-toolbox-plugin/pkg/utils"
	"strings"
)

var labelSelector = ""

type CheckPods struct {
}

var checkPods = RegisterCheck(CheckPods{})

func (c CheckPods) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&labelSelector, "label-selector", "s", labelSelector, "Label selector for pods")
}

func (c CheckPods) doCheck(cmd *cobra.Command, args []string) error {

	return c.CheckPods()
}

func (c CheckPods) CheckPods() error {
	opts := metav1.ListOptions{
		LabelSelector: labelSelector,
	}
	podList, err := client.GetKubeClient().CoreV1().Pods(client.GetNamespace()).List(context.Background(), opts)
	if err != nil {
		return err
	} else {
		var podErrors []error
		if len(podList.Items) == 0 {
			podErrors = append(podErrors, errors.New("No pods found"))
		}
		for _, pod := range podList.Items {
			for _, container := range pod.Spec.Containers {
				err := c.CheckContainer(pod, container)
				if err != nil {
					podErrors = append(podErrors, err)
				}
			}
		}
		logrus.Debugf("podErrors: %s", podErrors)
		if len(podErrors) > 0 {
			return multierr.Combine(podErrors...)
		} else {
			logrus.Infof("All Pods are OK")
		}
	}
	return nil
}

func (c CheckPods) CheckContainer(pod v1.Pod, container v1.Container) error {

	var errs []string

	if container.LivenessProbe == nil {
		errs = append(errs, "no LivenessProbe defined")
	}
	if container.ReadinessProbe == nil {
		errs = append(errs, "no ReadinessProbe defined")
	}
	if container.Resources.Limits == nil {
		errs = append(errs, "no ResourcesLimits defined")
	}
	if container.Resources.Requests == nil {
		errs = append(errs, "no ResourcesRequests defined")
	}

	containerStatus, err := findContainerStatus(pod.Status.ContainerStatuses, container.Name)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		if !containerStatus.Ready {
			errs = append(errs, "container not ready")
		}
		if !*containerStatus.Started {
			errs = append(errs, "container not started")
		}
	}

	if len(errs) > 0 {
		return errors.New(fmt.Sprintf("Pod [%-48s] - Container [%-48s] error: [%s] defined", pod.Name, container.Name, strings.Join(errs, ", ")))
	}
	return nil
}

func findContainerStatus(containerStatuses []v1.ContainerStatus, containerName string) (v1.ContainerStatus, error) {
	for _, status := range containerStatuses {
		if status.Name == containerName {
			return status, nil
		}
	}
	return v1.ContainerStatus{}, errors.New("no ContainerStatus found")
}

// checkVersionCmd represents the check version command
var checkPodsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Checks that Pods are compliant",
	Long:  ``,
	Run:   utils.CreateCommandFunc(checkPods.doCheck),
}

func init() {
	checkPods.AddFlags(checkPodsCmd.Flags())
	checkCmd.AddCommand(checkPodsCmd)
}
