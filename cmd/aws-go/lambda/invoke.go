package lambda

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/internal/spinner"
)

// lambda invoke command.
var invokeLambdaCmd = &cobra.Command{
	Use:   "invoke [function name]",
	Short: "Invoke the specified AWS Lambda function",
	Long: `Invokes the specified AWS Lambda function and returns the status code of the function call.
It's important to note that invoke command invokes the $LATEST version of the lambda function
available with RequestResponse invocation type`,
	Args:    cobra.ExactArgs(1),
	Example: "  aws-go lambda invoke testLambdaFunction",
	PreRun:  command.PreRun,
	RunE:    invokeFunction,
}

func init() {
	lambdaCmd.AddCommand(invokeLambdaCmd)
}

// run command for invoke function.
func invokeFunction(cmd *cobra.Command, args []string) error {
	sp := spinner.Default(spinner.Prefix[2])
	sp.Start()
	sess := lambda.New(command.Session)

	functionName := function.Function{
		Name: args[0],
	}

	lambdaService := &function.LambdaService{
		Function: functionName,
		Service:  sess,
	}

	resp, err := lambdaService.InvokeFunction()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			sp.Stop()
			return aerr
		} else {
			sp.Stop()
			return err
		}
	}

	sp.Stop()
	fmt.Printf("Status Code: %d\n", *resp.StatusCode)

	return nil
}
