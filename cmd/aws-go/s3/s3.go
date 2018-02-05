// Package s3 implements AWS S3 related operations.
package s3

import (
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
)

// s3 command.
var s3Cmd = &cobra.Command{
	Use:     "s3",
	Short:   "Perform AWS S3 specific operations",
	Long:    `List AWS S3 buckets`,
	Example: `  aws-go s3 list`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	command.AddCommand(s3Cmd)
}
