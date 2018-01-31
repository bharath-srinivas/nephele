package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/spinner"
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
	Run:     listFunctions,
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
	Run:     invokeFunction,
}

func init() {
	Command.AddCommand(lambdaCmd)
	lambdaCmd.AddCommand(listLambdaCmd)
	lambdaCmd.AddCommand(invokeLambdaCmd)
}

// run command for lambda list.
func listFunctions(cmd *cobra.Command, args []string) {
	sp := spinner.Default(spinnerPrefix[1])
	sp.Start()
	sess := lambda.New(Session)

	lambdaService := function.LambdaService{
		Service: sess,
	}

	resp, err := lambdaService.GetFunctions()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			sp.Stop()
			fmt.Println(aerr.Error())
		} else {
			sp.Stop()
			fmt.Println(err.Error())
		}
	} else {
		sp.Stop()
		for index, fn := range resp.Functions {
			var functionDescription string
			if fn.Description != nil {
				functionDescription = *fn.Description
			}

			fmt.Fprintln(os.Stdout, *fn.FunctionName, "\n", " -description:", functionDescription, "\n",
				" -runtime:", *fn.Runtime, "\n", " -memory:", *fn.MemorySize, "\n",
				" -timeout:", *fn.Timeout, "\n", " -handler:", *fn.Handler, "\n",
				" -role:", *fn.Role, "\n", " -version:", *fn.Version)
			if index < len(resp.Functions)-1 {
				fmt.Printf("\n")
			}
		}
	}
}

// run command for invoke function.
func invokeFunction(cmd *cobra.Command, args []string) {
	sp := spinner.Default(spinnerPrefix[2])
	sp.Start()
	sess := lambda.New(Session)

	lambdaService := function.LambdaService{
		Service: sess,
	}

	resp, err := lambdaService.InvokeFunction(args[0])

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			sp.Stop()
			fmt.Println(aerr.Error())
		} else {
			sp.Stop()
			fmt.Println(err.Error())
		}
	} else {
		sp.Stop()
		fmt.Printf("Status Code: %d\n", *resp.StatusCode)
	}
}
