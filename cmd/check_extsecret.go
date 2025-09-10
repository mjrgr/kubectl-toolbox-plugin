package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubectl-toolbox-plugin/pkg/utils"
)

var extSecretNamespace = "external-secrets"

type CheckExtSecret struct {
}

var checkExtSecret = RegisterCheck(CheckExtSecret{})

func (c CheckExtSecret) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&extSecretNamespace, "es-namespace", "e", extSecretNamespace, "ExternalSecrets secret namespace on Kubernetes")
	viper.BindPFlag("TB_ES_NAMESPACE", pflag.Lookup("es-namespace"))
}

func (c CheckExtSecret) doCheck(cmd *cobra.Command, args []string) error {

	logrus.Infof("extSecretNamespace = [%s]", extSecretNamespace)
	esNamespace, err := client.GetKubeClient().CoreV1().Namespaces().Get(context.Background(), extSecretNamespace, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if esNamespace == nil {
		return errors.New(fmt.Sprintf("ExternalSecrets namespace [%s] could not be found", esNamespace.Name))
	} else {
		logrus.Infof("ExternalSecrets namespace [%s] found", esNamespace.Name)
		return c.checkAllPodsRunning(esNamespace)
	}

}

func (c CheckExtSecret) checkAllPodsRunning(esNamespace *v1.Namespace) error {
	podList, err := client.GetKubeClient().CoreV1().Pods(esNamespace.Name).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	} else {
		var pods []string
		for _, p := range podList.Items {
			if p.Status.Phase != "Running" {
				pods = append(pods, p.Name)
			}
		}
		if len(pods) > 0 {
			return errors.New(fmt.Sprintf("ExternalSecrets Pods %s are not running", pods))
		} else {
			logrus.Infof("ExternalSecrets Pods are running")
		}
	}
	return nil
}

// checkExtSecretCmd represents the check external-secret command
var checkExtSecretCmd = &cobra.Command{
	Use:   "external-secrets",
	Short: "Checks that kubernetes has ExternalSecrets Operator",
	Long:  ``,
	Run:   utils.CreateCommandFunc(checkExtSecret.doCheck),
}

func init() {
	checkCmd.AddCommand(checkExtSecretCmd)
}
