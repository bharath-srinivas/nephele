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

func init() {
	envCmd.AddCommand(useCmd)
	useCmd.Flags().StringVarP(&store.Profile, "profile", "p", "", "the name of the profile")
	useCmd.Flags().StringVarP(&store.Region, "region", "r", "us-east-1", "the region to use")
}

// env use run command.
func useEnv(cmd *cobra.Command, args []string) {
	if store.Profile == "" {
		cmd.Usage()
		os.Exit(0)
	}
	store.UseProfile()
}
