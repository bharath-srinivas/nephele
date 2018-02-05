// Package rds implements AWS RDS related operations.
package rds

import (
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
)

// rds command.
var rdsCmd = &cobra.Command{
	Use:     "rds",
	Short:   "Perform AWS RDS specific operations",
	Long:    `List AWS RDS instances`,
	Args:    cobra.NoArgs,
	Example: `  aws-go rds list`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	command.AddCommand(rdsCmd)
}
