package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Prints the version information and exit",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aws-go:")
		fmt.Println(" version 	: v0.1.0")
		fmt.Println(" build date	: 2017-12-31")
		fmt.Println(" go version	: go1.9.2")
		fmt.Println(" platform	: linux/amd64")
	},
}

func init() {
	Command.AddCommand(versionCmd)
}