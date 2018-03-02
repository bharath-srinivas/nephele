package env

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/store"
)

// env create command.
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new AWS profile with specified region (if provided)",
	Args:    cobra.NoArgs,
	Example: "  aws-go env create --profile staging --region us-west-1",
	RunE:    createEnv,
}

func init() {
	envCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&store.Profile, "profile", "p", "", "the name of the profile")
	createCmd.Flags().StringVarP(&store.Region, "region", "r", "us-east-1", "the region to use")
}

// env create run command.
func createEnv(cmd *cobra.Command, args []string) error {
	if store.Profile == "" {
		return cmd.Usage()
	}

	db := store.NewSession()
	defer db.Close()

	if db.EntryExists(store.Profile) {
		return errors.New("error: profile already exists")
	}

	var accessId string
	var secretKey string
	fmt.Print("AWS Access Key ID: ")
	_, err := fmt.Scanln(&accessId)
	if err != nil {
		return err
	}

	fmt.Print("AWS Secret Access Key: ")
	_, err1 := fmt.Scanln(&secretKey)
	if err1 != nil {
		return err1
	}

	db.SetCredentials(accessId, secretKey)
	fmt.Printf("\nSuccessfully added new profile: %s\n", store.Profile)
	fmt.Printf("Switched to '%s' with region '%s'\n", store.Profile, store.Region)
	return nil
}
