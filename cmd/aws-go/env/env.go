// Package env manages the environments for aws-go.
package env

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/bharath-srinivas/aws-go/internal/colors"
	"github.com/bharath-srinivas/aws-go/store"
)

// constant for string color wrapping
const escape = "\x1b"

// enable env listing.
var listEnv bool

// delete a env.
var delEnv string

// env command.
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage AWS profile configurations",
	Args:  cobra.NoArgs,
	Example: `  aws-go env --list
  aws-go env --delete staging`,
	RunE: run,
}

func init() {
	command.AddCommand(envCmd)
	envCmd.Flags().BoolVarP(&listEnv, "list", "l", false, "list all the available profiles")
	envCmd.Flags().StringVarP(&delEnv, "delete", "d", "", "delete the specified profile")
}

// run command.
func run(cmd *cobra.Command, args []string) error {
	if !listEnv && delEnv == "" {
		return cmd.Usage()
	}

	if listEnv {
		return listProfiles()
	}
	return deleteProfile()
}

// list profiles.
func listProfiles() error {
	db := store.NewSession()
	defer db.Close()

	rows := db.ListProfiles()
	for rows.Next() {
		var profile string
		rows.Scan(&profile)

		if db.CurrentProfile(profile) {
			fmt.Printf("%s[1;%dm%s* (current)%s[m\n", escape, colors.Green, profile, escape)
		} else {
			fmt.Println(profile)
		}
	}
	return nil
}

// delete profile.
func deleteProfile() error {
	db := store.NewSession()
	defer db.Close()

	if !db.EntryExists(delEnv) {
		return errors.New("error: no such profile found")
	} else if db.CurrentProfile(delEnv) {
		return errors.New("error: cannot delete currently active profile")
	}

	db.DeleteProfile(delEnv)
	fmt.Printf("Successfully deleted '%s' from config\n", delEnv)
	return nil
}
