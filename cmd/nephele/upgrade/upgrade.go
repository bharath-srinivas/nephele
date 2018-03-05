package upgrade

import (
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/nephele/cmd/nephele/command"
	"github.com/bharath-srinivas/nephele/cmd/nephele/version"
	"github.com/bharath-srinivas/nephele/upgrade"
)

// upgrade command.
var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Short:   "Upgrade nephele to the latest version",
	Args:    cobra.NoArgs,
	Example: "  nephele upgrade",
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
