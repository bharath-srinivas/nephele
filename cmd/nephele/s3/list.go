package s3

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/nephele/cmd/nephele/command"
	"github.com/bharath-srinivas/nephele/function"
	"github.com/bharath-srinivas/nephele/internal/spinner"
)

// continuation token for the bucket.
var token string

// page size for s3 objects.
var maxCount int64

// prefix of the s3 objects.
var prefix string

// map for storing clear functions.
var clear map[string]func()

// s3 list command.
var listCmd = &cobra.Command{
	Use:   "list [Bucket name]",
	Short: "List all the available S3 buckets or objects in a bucket",
	Args:  cobra.MaximumNArgs(1),
	Example: `To list S3 buckets:  
  nephele s3 list

To list S3 objects in a bucket:
  nephele s3 list [bucket-name]

For fetching more objects than the default limit:
  nephele s3 list [bucket-name] -c <count>

For fetching the next set of objects in a bucket:
  nephele s3 list [bucket-name] -t <token>

Note: Maximum number of objects you can fetch per request is limited to 1000`,
	PreRun: command.PreRun,
	RunE:   list,
}

func init() {
	s3Cmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&token, "token", "t", "", "the token for fetching the next or previous set of s3 objects")
	listCmd.Flags().Int64VarP(&maxCount, "count", "c", 100, "number of objects to fetch per request")
	listCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "search for s3 objects containing specified prefix")

	// TODO: find a way to fix this properly so that this ugly hack can be removed.
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	clear["darwin"] = clear["linux"]
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// run command.
func list(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return listObjects(args)
	}
	return listBuckets()
}

// list s3 buckets.
func listBuckets() error {
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

// list s3 objects in a bucket.
func listObjects(args []string) error {
	sp := spinner.Default(spinner.Prefix[1])
	sp.Start()

	cmdPath, err := exec.LookPath("less")
	if err != nil {
		cmdPath, err = exec.LookPath("more")
		if err != nil {
			return err
		}
	}

	pager := exec.Command(cmdPath)
	stdin, err := pager.StdinPipe()
	if err != nil {
		return err
	}
	pager.Stdout = os.Stdout

	sess := s3.New(command.Session)
	bucketName := function.S3{
		Name: args[0],
	}

	var options function.S3Options
	options.MaxCount = maxCount
	if token != "" {
		options.ContinuationToken = token
	}

	if prefix != "" {
		options.Prefix = prefix
	}

	s3Service := &function.S3Service{
		S3:        bucketName,
		S3Options: options,
		Service:   sess,
	}

	resp, err := s3Service.GetObjects()
	if err != nil {
		sp.Stop()
		return err
	}

	clrScr()
	writer := tabwriter.NewWriter(stdin, 0, 0, 2, ' ', 0)
	go func() {
		defer stdin.Close()
		for _, object := range resp.Contents {
			io.WriteString(writer, object.LastModified.String()+"\t"+strconv.Itoa(int(*object.Size))+"\t"+
				*object.Key+"\n")
		}
		writer.Flush()

		io.WriteString(stdin, "\n\nTotal Objects: "+strconv.Itoa(int(*resp.KeyCount)))
		if resp.ContinuationToken != nil {
			io.WriteString(stdin, "\nToken for fetching the previous set of objects: "+*resp.ContinuationToken)
		}

		if resp.NextContinuationToken != nil {
			io.WriteString(stdin, "\nToken for fetching the next set of objects: "+*resp.NextContinuationToken)
		}
	}()
	sp.Stop()
	pager.Run()
	return nil
}

// clear screen.
func clrScr() {
	if value, ok := clear[runtime.GOOS]; ok {
		value()
	}
}
