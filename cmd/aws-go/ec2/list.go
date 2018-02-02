package ec2

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/internal/spinner"
)

// list command.
var listCmd = &cobra.Command{
	Use:    "list",
	Short:  "List all the available EC2 instances",
	Args:   cobra.NoArgs,
	PreRun: command.PreRun,
	Run:    listInstances,
}

func init() {
	ec2Cmd.AddCommand(listCmd)
}

// run command.
func listInstances(cmd *cobra.Command, args []string) {
	sp := spinner.Default(spinner.Prefix[1])
	sp.Start()
	sess := ec2.New(command.Session)

	ec2Service := &function.EC2Service{
		Service: sess,
	}

	resp, err := ec2Service.GetInstances()

	if err != nil {
		sp.Stop()
		fmt.Println(err.Error())
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetRowLine(true)
		table.SetHeader([]string{
			"Instance Name",
			"Instance ID",
			"Instance State",
			"Private IPv4 Address",
			"Public IPv4 Address",
			"Instance Type",
		})
		table.AppendBulk(resp)
		sp.Stop()
		table.Render()
	}
}
