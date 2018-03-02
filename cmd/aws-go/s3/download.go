package s3

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"

	"github.com/bharath-srinivas/aws-go/cmd/aws-go/command"
	"github.com/bharath-srinivas/aws-go/function"
	"github.com/bharath-srinivas/aws-go/internal/spinner"
)

// objects file.
var objectsFile string

// s3 download command.
var downloadCmd = &cobra.Command{
	Use:   "download [src] [dst]",
	Short: "Download a S3 object from the specified bucket",
	Args:  cobra.MaximumNArgs(2),
	Example: ` To download a S3 object:
  aws-go s3 download [bucket-name:object-name] [dst-file-name]

To download an object from sub directory of a bucket:
  aws-go s3 download [bucket-name/sub-dir/:object-name] [dst-file-name]

To download multiple objects concurrently:
  aws-go s3 download -o objects-file.json

Note: Sub-directory name is case-sensitive and requires '/' at the end`,
	PreRun: command.PreRun,
	RunE:   download,
}

func init() {
	s3Cmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&objectsFile, "objects-file", "o", "", "JSON file containing the objects")
}

// run command.
func download(cmd *cobra.Command, args []string) error {
	if objectsFile != "" {
		return multiDownload()
	}
	return downloadObject(args)
}

// download s3 object.
func downloadObject(args []string) error {
	if len(args) < 2 {
		return errors.New("error: download command requires both source and destination")
	}

	if !strings.Contains(args[0], ":") {
		return errors.New("error: invalid src: '" + args[0] + "'")
	}

	sp := spinner.Default(spinner.Prefix[3])
	sp.Start()

	sess := s3manager.NewDownloader(command.Session)
	input := strings.Split(args[0], ":")

	downloader := function.S3{
		Name:       input[0],
		Key:        input[1],
		FileName:   args[1],
		Downloader: sess,
	}

	downloaderService := &function.S3Service{
		S3: downloader,
	}

	_, err := downloaderService.DownloadObject()
	if err != nil {
		sp.Stop()
		return err
	}

	sp.Stop()
	fmt.Println("The requested file has been downloaded successfully")
	return nil
}

// download multiple s3 objects.
func multiDownload() error {
	sp := spinner.Default(spinner.Prefix[3])
	sp.Start()

	sess := s3.New(command.Session)
	s3Service := &function.S3Service{
		Service: sess,
	}

	if err := s3Service.MultiObjectDownload(objectsFile); err != nil {
		sp.Stop()
		return err
	}

	sp.Stop()
	fmt.Println("The requested files have been downloaded successfully")
	return nil
}
