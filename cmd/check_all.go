package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"kubectl-toolbox-plugin/pkg/utils"
)

// checkAllCmd represents the check all command
var checkAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Run all checks",
	Long:  ``,
	Run:   utils.CreateCommandFunc(doCheckAll),
}

func init() {
	for _, c := range checks {
		c.AddFlags(checkAllCmd.Flags())
	}
	checkCmd.AddCommand(checkAllCmd)
}

func doCheckAll(cmd *cobra.Command, args []string) error {
	var allErrors []error
	for _, c := range checks {
		allErrors = append(allErrors, c.doCheck(cmd, args))
	}
	return multierr.Combine(allErrors...)
}
