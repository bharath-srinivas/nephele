// Package s3 implements AWS S3 related operations.
package s3

import (
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/nephele/cmd/nephele/command"
)

// s3 command.
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Perform AWS S3 specific operations",
	Long:  `List AWS S3 buckets`,
	Example: `  nephele s3 list
  nephele s3 list [bucket-name]
  nephele s3 download [bucket-name:object-name] [dst-file-name]
  nephele s3 download -o objects.json`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	command.AddCommand(s3Cmd)
}
