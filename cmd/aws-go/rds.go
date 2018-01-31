package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/spinner"
)

// rds command.
var rdsCmd = &cobra.Command{
	Use:     "rds",
	Short:   "Perform AWS RDS specific operations",
	Long:    `List, start or stop AWS RDS instances`,
	Args:    cobra.NoArgs,
	Example: `  aws-go rds list`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// rds list command.
var listRdsCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all the available AWS RDS instances",
	Args:    cobra.NoArgs,
	Example: "  aws-go rds list",
	Run:     listRDSInstances,
}

func init() {
	Command.AddCommand(rdsCmd)
	rdsCmd.AddCommand(listRdsCmd)
}

// run command.
func listRDSInstances(cmd *cobra.Command, args []string) {
	sp := spinner.Default(spinnerPrefix[1])
	sp.Start()
	sess := rds.New(Session)

	rdsService := &function.RDSService{
		Service: sess,
	}

	resp, err := rdsService.GetRDSInstances()

	if err != nil {
		sp.Stop()
		fmt.Println(err.Error())
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetColWidth(20)
		table.SetRowLine(true)
		table.SetHeader([]string{
			"DB Instance ID",
			"DB Instance Status",
			"Endpoint",
			"DB Instance Class",
			"Engine",
			"Engine Version",
			"Multi-AZ",
		})
		table.AppendBulk(resp)
		sp.Stop()
		table.Render()
	}
}
