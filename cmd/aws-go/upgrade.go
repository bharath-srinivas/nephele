package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/bharath-srinivas/aws-go/utils"
)

// upgrade command.
var upgradeCmd = &cobra.Command{
	Use: "upgrade",
	Short: "Upgrade aws-go to the latest version",
	Args: cobra.NoArgs,
	Example: "aws-go upgrade",
	Run: func(cmd *cobra.Command, args []string) {
		if err := utils.Upgrade(version); err != nil {
			fmt.Println(err)
		}
	},

}

func init() {
	Command.AddCommand(upgradeCmd)
}