package version

import (
	"fmt"
	"runtime"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/spf13/cobra"
)

// current version.
const Version = "v0.3.0"

// version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version information and exit",
	Args:  cobra.NoArgs,
	Run:   printVersion,
}

func init() {
	command.AddCommand(versionCmd)
}

// run command.
func printVersion(cmd *cobra.Command, args []string) {
	fmt.Println("aws-go:")
	fmt.Println(" version 	:", Version)
	fmt.Println(" build date	: 2018-01-26")
	fmt.Println(" go version	:", runtime.Version())
	fmt.Println(" platform	:", runtime.GOOS+"/"+runtime.GOARCH)
}
