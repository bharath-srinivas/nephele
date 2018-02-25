package version

import (
	"fmt"
	"runtime"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/spf13/cobra"
	"time"
)

// current version.
const Version = "v0.3.0"

// build date.
var buildDate = fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day())

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
	fmt.Println(" build date	:", buildDate)
	fmt.Println(" go version	:", runtime.Version())
	fmt.Println(" platform	:", runtime.GOOS+"/"+runtime.GOARCH)
}
