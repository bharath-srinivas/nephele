package env

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/store"
)

// env use command.
var useCmd = &cobra.Command{
	Use:     "use",
	Short:   "Use the specified AWS profile and region (if provided)",
	Example: "  aws-go env use --profile staging --region eu-west-1",
	Run:     useEnv,
}

// env use run command.
func useEnv(cmd *cobra.Command, args []string) {
	if store.Profile == "" {
		cmd.Usage()
		os.Exit(0)
	}
	store.UseProfile()
}
