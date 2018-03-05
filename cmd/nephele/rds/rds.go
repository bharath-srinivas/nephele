// Package rds implements AWS RDS related operations.
package rds

import (
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/nephele/cmd/nephele/command"
)

// rds command.
var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "Perform AWS RDS specific operations",
	Long:  `List AWS RDS instances`,
	Args:  cobra.NoArgs,
	Example: `  nephele rds list
  nephele rds start test-db-instance
  nephele rds stop test-db-instance`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	command.AddCommand(rdsCmd)
}
