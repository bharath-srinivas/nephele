package lambda

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/internal/spinner"
)

// lambda list command.
var listLambdaCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all available AWS Lambda functions and their configurations",
	Args:    cobra.NoArgs,
	Example: "  aws-go lambda list",
	PreRun:  command.PreRun,
	RunE:    listFunctions,
}

func init() {
	lambdaCmd.AddCommand(listLambdaCmd)
}

// run command for lambda list.
func listFunctions(cmd *cobra.Command, args []string) error {
	sp := spinner.Default(spinner.Prefix[1])
	sp.Start()
	sess := lambda.New(command.Session)

	lambdaService := &function.LambdaService{
		Service: sess,
	}

	resp, err := lambdaService.GetFunctions()

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

	return nil
}
