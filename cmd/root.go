package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kubectl-toolbox-plugin/pkg/cli"
	"kubectl-toolbox-plugin/pkg/kube"
	"kubectl-toolbox-plugin/pkg/log"
	"kubectl-toolbox-plugin/pkg/utils"
	"os"
)

var settings = cli.New()
var client = kube.New(settings.RESTClientGetter())

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   utils.CommandName,
	Short: "Kubectl toolbox plugin command",
	Long:  `Lightweight init container tool for Kubernetes checks`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Error("Failed to execute")
		os.Exit(1)
	}
}

func PrintCompletion() {
	var args []string
	args = append(args, cobra.ShellCompRequestCmd)
	args = append(args, os.Args[1:]...)

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.AutomaticEnv()
	flags := rootCmd.PersistentFlags()
	settings.AddFlags(flags)
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	lvl, err := logrus.ParseLevel(settings.LogLevel)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logrus.SetLevel(lvl)
	logrus.SetFormatter(&log.PwcCtlLogFormat{
		Color: true,
	})
}
