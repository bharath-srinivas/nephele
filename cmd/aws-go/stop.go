package cmd

import (
	"github.com/spf13/cobra"
	"function"
)

var stopCmd = &cobra.Command{
	Use: "stop [instance id]",
	Short: "Stop the specified EC2 instance",
	Args:cobra.ExactArgs(1),
	Example: "aws-go stop i-0a12b345c678de",
	Run: func(cmd *cobra.Command, args []string) {
		function.StopInstance(args[0], dryRun)
	},
}

func init() {
	Command.AddCommand(stopCmd)
	stopCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "perform the operation with dry run enabled")
}