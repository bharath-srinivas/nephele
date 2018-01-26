package cmd

import (
	"github.com/spf13/cobra"
	"github.com/bharath-srinivas/aws-go/function"
)

// list command.
var listCmd = &cobra.Command{
	Use: "list",
	Short: "List all the available EC2 instances",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		function.ListInstances()
	},
}

func init() {
	Command.AddCommand(listCmd)
}