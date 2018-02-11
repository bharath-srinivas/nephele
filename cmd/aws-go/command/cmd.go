// Package cmd implements all the commands used by aws-go.
package command

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/function"
)

// logo for aws-go
var logo = `
__________       _________    _________
___    |_ |     / /_  ___/    __  ____/_____
__  /| |_ | /| / /_____ \     _  / __ _  __ \
_  ___ |_ |/ |/ / ____/ /     / /_/ / / /_/ /
/_/  |_|___/|__/  /____/      \____/  \____/

`

// long description for aws-go
var description = logo + `
AWS Go is a CLI tool for managing AWS services without the need
to login to the AWS console built to be fast and easy to use.`

// AWS Session instance.
var Session *session.Session

// Main command.
var Command = &cobra.Command{
	Use:           "aws-go",
	Long:          description,
	RunE:          run,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func AddCommand(cmd *cobra.Command) {
	Command.AddCommand(cmd)
}

// Execute executes the provided command.
func Execute() {
	Command.Execute()
}

// run is the actual function for the main command.
func run(cmd *cobra.Command, args []string) error {
	fmt.Println(logo)
	return cmd.Usage()
}

// preRun will initialize the session required for all the child commands.
func PreRun(cmd *cobra.Command, args []string) {
	Session = function.NewSession()
}
