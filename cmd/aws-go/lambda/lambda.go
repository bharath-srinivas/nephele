// Package lambda implements AWS Lambda related operations.
package lambda

import (
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
)

// lambda command.
var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Perform AWS Lambda specific operations",
	Long:  `List all the AWS Lambda functions or invoke a AWS Lambda function`,
	Args:  cobra.NoArgs,
	Example: `  aws-go lambda list
  aws-go lambda invoke testFunction`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	command.AddCommand(lambdaCmd)
}
