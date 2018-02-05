package ec2

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/internal/spinner"
)

// list filters.
var filters []string

// filters file.
var filtersFile string

// list all.
var listAll bool

// list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the available EC2 instances",
	Args:  cobra.NoArgs,
	Example: `  aws-go ec2 list
  aws-go ec2 list --filters name=web,az=us-east-1a
  aws-go ec2 list -F filters.json
  aws-go ec2 list --all`,
	PreRun: command.PreRun,
	RunE:   list,
}

func init() {
	ec2Cmd.AddCommand(listCmd)
	listCmd.Flags().StringSliceVarP(&filters, "filters", "f", nil, "filter list output")
	listCmd.Flags().StringVarP(&filtersFile, "filters-file", "F", "", "JSON file containing the filters")
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "list all the instances in JSON format")
}

// run command.
func list(cmd *cobra.Command, args []string) error {
	if listAll {
		if err := listJSON(); err != nil {
			return err
		}
	} else {
		if err := listInstances(); err != nil {
			return err
		}
	}

	return nil
}

// list all instances in JSON format.
func listJSON() error {
	sess := ec2.New(command.Session)
	ec2Service := &function.EC2Service{
		Service: sess,
	}

	resp, err := ec2Service.GetAllInstances()
	if err != nil {
		return err
	}

	result, _ := json.MarshalIndent(resp.Reservations, "", "    ")
	fmt.Println(string(result))

	return nil
}

// list instances in table format.
func listInstances() error {
	sp := spinner.Default(spinner.Prefix[1])
	sp.Start()

	sess := ec2.New(command.Session)
	ec2Service := &function.EC2Service{
		Service: sess,
	}

	if filters != nil {
		if err := ec2Service.SetFilters(filters); err != nil {
			sp.Stop()
			return err
		}
	} else if filtersFile != "" {
		if err := ec2Service.LoadFiltersFromFile(filtersFile); err != nil {
			sp.Stop()
			return err
		}
	}

	resp, err := ec2Service.GetInstances()
	if err != nil {
		sp.Stop()
		return err
	}

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

	return nil
}
