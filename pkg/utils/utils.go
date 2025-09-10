package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"os"
)

var Zero int64 = 0

const CommandName string = "kubectl-toolbox-plugin"

func CreateCommandFunc(cmdFunc func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		errors := multierr.Errors(cmdFunc(cmd, args))
		if errors != nil {
			for _, err := range errors {
				logrus.Error(err)
			}
			os.Exit(1)
		}
	}
}
