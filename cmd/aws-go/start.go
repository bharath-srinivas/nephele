package cmd

import (
	"github.com/spf13/cobra"
	"aws-go/function"
)

var dryRun bool

var startCmd = &cobra.Command{
	Use: "start [instance id]",
	Short: "Start the specified EC2 instance",
	Args: cobra.ExactArgs(1),
	Example: "aws-go start i-0a12b345c678de",
	Run: func(cmd *cobra.Command, args []string) {
		function.StartInstance(args[0], dryRun)
	},
}

func init() {
	Command.AddCommand(startCmd)
	startCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "perform the operation with dry run enabled")
}