package cmd

import (
	"github.com/spf13/cobra"
	"kubectl-toolbox-plugin/pkg/utils"
)

var image = "docker.io/busybox"
var command = "bash"

// busyboxCmd represents the busybox command
var busyboxCmd = &cobra.Command{
	Use:   "busybox",
	Short: "Starts a busybox pod in Kubernetes",
	Long:  ``,
	Run:   utils.CreateCommandFunc(startBusyboxPod),
}

func init() {
	busyboxCmd.Flags().StringVarP(&image, "image", "i", image, "Image to use for the busybox Pod")
	busyboxCmd.Flags().StringVarP(&command, "command", "c", command, "Command to execute in the busybox Pod")
	rootCmd.AddCommand(busyboxCmd)
}

func startBusyboxPod(cmd *cobra.Command, args []string) error {
	return client.RunRemotePod(image, command)
}
