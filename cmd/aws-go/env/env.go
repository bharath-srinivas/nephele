// Package env manages the environments for aws-go.
package env

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/bharath-srinivas/aws-go/store"
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
	Run: run,
}

func init() {
	command.AddCommand(envCmd)

	envCmd.AddCommand(createCmd)
	envCmd.AddCommand(useCmd)

	envCmd.Flags().BoolVarP(&listEnv, "list", "l", false, "list all the available profiles")
	envCmd.Flags().StringVarP(&delEnv, "delete", "d", "", "delete the specified profile")

	createCmd.Flags().StringVarP(&store.Profile, "profile", "p", "", "the name of the profile")
	createCmd.Flags().StringVarP(&store.Region, "region", "r", "us-east-1", "the region to use")

	useCmd.Flags().StringVarP(&store.Profile, "profile", "p", "", "the name of the profile")
	useCmd.Flags().StringVarP(&store.Region, "region", "r", "us-east-1", "the region to use")
}

// run command.
func run(cmd *cobra.Command, args []string) {
	if !listEnv && delEnv == "" {
		cmd.Usage()
		os.Exit(0)
	}

	if listEnv {
		store.ListProfiles()
	} else if delEnv != "" {
		store.DeleteProfile(delEnv)
	}
}
