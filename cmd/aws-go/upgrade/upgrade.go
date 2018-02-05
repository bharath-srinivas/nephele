package upgrade

import (
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/bharath-srinivas/aws-go/cmd/aws-go/version"
	"github.com/bharath-srinivas/aws-go/upgrade"
)

// upgrade command.
var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Short:   "Upgrade aws-go to the latest version",
	Args:    cobra.NoArgs,
	Example: "aws-go upgrade",
	RunE:    doUpgrade,
}

func init() {
	command.AddCommand(upgradeCmd)
}

// run command.
func doUpgrade(cmd *cobra.Command, args []string) error {
	if err := upgrade.Upgrade(version.Version); err != nil {
		return err
	}
	return nil
}
