package s3

import (
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/internal/spinner"
)

// s3 list command.
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all the available S3 buckets",
	Args:    cobra.NoArgs,
	Example: `  aws-go s3 list`,
	PreRun:  command.PreRun,
	RunE:    list,
}

func init() {
	s3Cmd.AddCommand(listCmd)
}

// run command.
func list(cmd *cobra.Command, args []string) error {
	sp := spinner.Default(spinner.Prefix[1])
	sp.Start()

	sess := s3.New(command.Session)
	s3Service := &function.S3Service{
		Service: sess,
	}

	resp, err := s3Service.GetBuckets()
	if err != nil {
		sp.Stop()
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{
		"Bucket Name",
		"Creation Date",
	})

	for _, bucket := range resp.Buckets {
		tableData := []string{
			*bucket.Name,
			bucket.CreationDate.String(),
		}
		table.Append(tableData)
	}
	sp.Stop()
	table.Render()

	return nil
}
