package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Check interface {
	AddFlags(fs *pflag.FlagSet)
	doCheck(cmd *cobra.Command, args []string) error
}

var checks = []Check{}

func RegisterCheck(c Check) Check {
	checks = append(checks, c)
	return c
}

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Carries out different types of checks",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
