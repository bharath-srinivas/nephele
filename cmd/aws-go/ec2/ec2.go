// Package ec2 implements AWS EC2 related operations.
package ec2

import (
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
)

// ec2 command.
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Perform AWS EC2 specific operations",
	Long:  `List, start or stop AWS EC2 instances`,
	Example: `  aws-go ec2 list
  aws-go ec2 start i-0a12b345c678de
  aws-go ec2 stop i-0a12b345c678de`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	command.AddCommand(ec2Cmd)
}
