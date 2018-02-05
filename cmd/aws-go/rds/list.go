package rds

import (
	"os"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/internal/spinner"
)

// rds list command.
var listRdsCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all the available AWS RDS instances",
	Args:    cobra.NoArgs,
	Example: "  aws-go rds list",
	PreRun:  command.PreRun,
	RunE:    listRDSInstances,
}

func init() {
	rdsCmd.AddCommand(listRdsCmd)
}

// run command.
func listRDSInstances(cmd *cobra.Command, args []string) error {
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
