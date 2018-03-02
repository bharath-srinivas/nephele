package env

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/store"
)

// env use command.
var useCmd = &cobra.Command{
	Use:     "use",
	Short:   "Use the specified AWS profile and region (if provided)",
	Args:    cobra.NoArgs,
	Example: "  aws-go env use --profile staging --region eu-west-1",
	RunE:    useEnv,
}

func init() {
	envCmd.AddCommand(useCmd)
	useCmd.Flags().StringVarP(&store.Profile, "profile", "p", "", "the name of the profile")
	useCmd.Flags().StringVarP(&store.Region, "region", "r", "us-east-1", "the region to use")
}

// env use run command.
func useEnv(cmd *cobra.Command, args []string) error {
	if store.Profile == "" {
		return cmd.Usage()
	}

	db := store.NewSession()
	defer db.Close()

	if !db.EntryExists(store.Profile) {
		return errors.New("error: no such profile found")
	}

	db.UseProfile()
	fmt.Printf("Switched to '%s' with region '%s'\n", store.Profile, store.Region)
	return nil
}
