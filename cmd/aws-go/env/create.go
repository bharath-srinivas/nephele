package env

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/store"
)

// env create command.
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new AWS profile with specified region (if provided)",
	Example: "  aws-go env create --profile staging --region us-west-1",
	Run:     createEnv,
}

// env create run command.
func createEnv(cmd *cobra.Command, args []string) {
	if store.Profile == "" {
		cmd.Usage()
		os.Exit(0)
	}
	store.SetCredentials()
}
