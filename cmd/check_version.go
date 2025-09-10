package cmd

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"kubectl-toolbox-plugin/pkg/utils"
)

var minimalVersion = "v1.28.0"

type CheckVersion struct {
}

var checkVersion = RegisterCheck(CheckVersion{})

func (c CheckVersion) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&minimalVersion, "minimal", "m", minimalVersion, "Minimal version of Kubernetes to have")
}

func (c CheckVersion) doCheck(cmd *cobra.Command, args []string) error {
	minVersion, err := version.NewSemver(minimalVersion)
	if err != nil {
		return err
	} else {
		v, err := client.ServerVersion()
		if err != nil {
			return err
		} else {
			currentVersion, err := version.NewSemver(v.GitVersion)
			if err != nil {
				return err
			} else {
				logrus.Debugf("Kubernetes minimal version [%s]", minVersion)
				logrus.Debugf("Kubernetes current version [%s]", currentVersion)
				if minVersion.GreaterThan(currentVersion) {
					return errors.New(fmt.Sprintf("Kubernetes version [%s] is lower than minimal required: %s", currentVersion, minimalVersion))
				} else {
					logrus.Infof("Kubernetes version [%s] is OK", currentVersion)
				}
				return nil
			}
		}
	}

}

// checkVersionCmd represents the check version command
var checkVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Checks that kubernetes version is ok",
	Long:  ``,
	Run:   utils.CreateCommandFunc(checkVersion.doCheck),
}

func init() {
	checkVersion.AddFlags(checkVersionCmd.Flags())
	checkCmd.AddCommand(checkVersionCmd)
}
