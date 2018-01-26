package cmd

import (
	"os"

	"github.com/bharath-srinivas/aws-go/store"
	"github.com/spf13/cobra"
)

// enable env listing.
var listEnv bool

// delete a env.
var delEnv string

// env command.
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage AWS profile configurations",
	Example: `  aws-go env --list
  aws-go env --delete staging`,
	Run: func(cmd *cobra.Command, args []string) {
		if !listEnv && delEnv == "" {
			cmd.Usage()
			os.Exit(0)
		}

		if listEnv {
			store.ListProfiles()
		} else if delEnv != "" {
			store.DeleteProfile(delEnv)
		}
	},
}

// env create command.
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new AWS profile with specified region (if provided)",
	Example: "  aws-go env create --profile staging --region us-west-1",
	Run: func(cmd *cobra.Command, args []string) {
		if store.Profile == "" {
			cmd.Usage()
			os.Exit(0)
		}
		store.SetCredentials()
	},
}

// env use command.
var useCmd = &cobra.Command{
	Use:     "use",
	Short:   "Use the specified AWS profile and region (if provided)",
	Example: "  aws-go env use --profile staging --region eu-west-1",
	Run: func(cmd *cobra.Command, args []string) {
		if store.Profile == "" {
			cmd.Usage()
			os.Exit(0)
		}
		store.UseProfile()
	},
}

func init() {
	Command.AddCommand(envCmd)

	envCmd.AddCommand(createCmd)
	envCmd.AddCommand(useCmd)

	envCmd.Flags().BoolVarP(&listEnv, "list", "l", false, "list all the available profiles")
	envCmd.Flags().StringVarP(&delEnv, "delete", "d", "", "delete the specified profile")

	createCmd.Flags().StringVarP(&store.Profile, "profile", "p", "", "the name of the profile")
	createCmd.Flags().StringVarP(&store.Region, "region", "r", "us-east-1", "the region to use")

	useCmd.Flags().StringVarP(&store.Profile, "profile", "p", "", "the name of the profile")
	useCmd.Flags().StringVarP(&store.Region, "region", "r", "us-east-1", "the region to use")
}
