// Package cmd implements all the commands used by aws-go.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

// Main command.
var Command = &cobra.Command{
	Use:  "aws-go",
	Long: description,
	RunE: run,
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
