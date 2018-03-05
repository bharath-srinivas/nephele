package rds

import (
	"os"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/nephele/cmd/nephele/command"
	"github.com/bharath-srinivas/nephele/function"
	"github.com/bharath-srinivas/nephele/internal/spinner"
)

// rds list command.
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all the available AWS RDS instances",
	Args:    cobra.NoArgs,
	Example: "  nephele rds list",
	PreRun:  command.PreRun,
	RunE:    listInstances,
}

func init() {
	rdsCmd.AddCommand(listCmd)
}

// run command.
func listInstances(cmd *cobra.Command, args []string) error {
	sp := spinner.Default(spinner.Prefix[1])
	sp.Start()
	sess := rds.New(command.Session)

	rdsService := &function.RDSService{
		Service: sess,
	}

	resp, err := rdsService.GetRDSInstances()

	if err != nil {
		sp.Stop()
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(20)
	table.SetRowLine(true)
	table.SetHeader([]string{
		"DB Instance ID",
		"Status",
		"Endpoint",
		"Instance Class",
		"Engine/Version",
		"Multi-AZ",
	})
	table.AppendBulk(resp)
	sp.Stop()
	table.Render()

	return nil
}
