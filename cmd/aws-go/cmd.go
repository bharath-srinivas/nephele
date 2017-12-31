package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logo = `
__________       _________    _________
___    |_ |     / /_  ___/    __  ____/_____
__  /| |_ | /| / /_____ \     _  / __ _  __ \
_  ___ |_ |/ |/ / ____/ /     / /_/ / / /_/ /
/_/  |_|___/|__/  /____/      \____/  \____/

`
var description = logo + `
AWS Go is a CLI tool for managing AWS services without the need
to login to the AWS console built to be fast and easy to use.`

var Command = &cobra.Command{
	Use:   "aws-go",
	Long: description,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(logo)
		cmd.Usage()
	},
}

func Execute() {
	Command.Execute()
}