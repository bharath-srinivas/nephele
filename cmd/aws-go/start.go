package cmd

import (
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/spf13/cobra"
)

// dryRun enabled.
var dryRun bool

// start instance command.
var startCmd = &cobra.Command{
	Use:     "start [instance id]",
	Short:   "Start the specified EC2 instance",
	Args:    cobra.ExactArgs(1),
	Example: "aws-go start i-0a12b345c678de",
	Run: func(cmd *cobra.Command, args []string) {
		function.StartInstance(args[0], dryRun)
	},
}

func init() {
	Command.AddCommand(startCmd)
	startCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "perform the operation with dry run enabled")
}
