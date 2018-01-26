package cmd

import (
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/spf13/cobra"
)

// lambda command.
var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Perform AWS Lambda specific operations",
	Long:  `List all the AWS Lambda functions or invoke a AWS Lambda function`,
	Args:  cobra.NoArgs,
	Example: `  aws-go lambda list
  aws-go lambda invoke testLambdaFunction`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// lambda list command.
var listLambdaCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all available AWS Lambda functions and their configurations",
	Args:    cobra.NoArgs,
	Example: "  aws-go lambda list",
	Run: func(cmd *cobra.Command, args []string) {
		function.ListLambdaFunctions()
	},
}

// lambda invoke command.
var invokeLambdaCmd = &cobra.Command{
	Use:   "invoke [function name]",
	Short: "Invoke the specified AWS Lambda function",
	Long: `Invokes the specified AWS Lambda function and returns the status code of the function call.
It's important to note that invoke command invokes the $LATEST version of the lambda function
available with RequestResponse invocation type`,
	Args:    cobra.ExactArgs(1),
	Example: "  aws-go lambda invoke testLambdaFunction",
	Run: func(cmd *cobra.Command, args []string) {
		function.InvokeLambdaFunction(args[0])
	},
}

func init() {
	Command.AddCommand(lambdaCmd)
	lambdaCmd.AddCommand(listLambdaCmd)
	lambdaCmd.AddCommand(invokeLambdaCmd)
}
