// Package cmd implements all the commands used by nephele.
package command

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/nephele/function"
)

// logo for nephele
var logo = `
                  _          _      
                 | |        | |     
 _ __   ___ _ __ | |__   ___| | ___ 
| '_ \ / _ \ '_ \| '_ \ / _ \ |/ _ \
| | | |  __/ |_) | | | |  __/ |  __/
|_| |_|\___| .__/|_| |_|\___|_|\___|
           | |                      
           |_|

`

// long description for nephele
var description = logo + `
Nephelê is a CLI tool for managing AWS services without the need
to login to the AWS console built to be fast and easy to use.`

// AWS Session instance.
var Session *session.Session

// Main command.
var Command = &cobra.Command{
	Use:           "nephele",
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
