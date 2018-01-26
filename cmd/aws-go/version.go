package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// current version.
const version = "v0.2.1"

// version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version information and exit",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("aws-go:")
		fmt.Println(" version 	:", version)
		fmt.Println(" build date	: 2018-01-26")
		fmt.Println(" go version	:", runtime.Version())
		fmt.Println(" platform	:", runtime.GOOS+"/"+runtime.GOARCH)
	},
}

func init() {
	Command.AddCommand(versionCmd)
}
