package cmd

import (
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/spf13/cobra"
)

// rds command.
var rdsCmd = &cobra.Command{
	Use:     "rds",
	Short:   "Perform AWS RDS specific operations",
	Long:    `List, start or stop AWS RDS instances`,
	Args:    cobra.NoArgs,
	Example: `  aws-go rds list`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// rds list command.
var listRdsCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all the available AWS RDS instances",
	Args:    cobra.NoArgs,
	Example: "  aws-go rds list",
	Run: func(cmd *cobra.Command, args []string) {
		function.ListRDSInstances()
	},
}

func init() {
	Command.AddCommand(rdsCmd)
	rdsCmd.AddCommand(listRdsCmd)
}
