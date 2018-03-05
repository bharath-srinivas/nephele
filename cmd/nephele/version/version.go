package version

import (
	"fmt"
	"runtime"

	"github.com/bharath-srinivas/nephele/cmd/nephele/command"
	"github.com/spf13/cobra"
)

// current version.
const Version = "v0.4.0"

// build date.
const buildDate = "2018-03-03"

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
	fmt.Println("nephele:")
	fmt.Println(" version 	:", Version)
	fmt.Println(" build date	:", buildDate)
	fmt.Println(" go version	:", runtime.Version())
	fmt.Println(" platform	:", runtime.GOOS+"/"+runtime.GOARCH)
}
