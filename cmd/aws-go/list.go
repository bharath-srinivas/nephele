package cmd

import (
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/spf13/cobra"
)

// list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the available EC2 instances",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		function.ListInstances()
	},
}

func init() {
	Command.AddCommand(listCmd)
}
