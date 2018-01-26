// Package utils contains all the utility functions used by aws-go
package utils

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/github"
	//"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/cheggaaa/pb.v1"
)

//func getTermSize() (int, int, error) {
//	fd := int(os.Stdin.Fd())
//	if !terminal.IsTerminal(fd) {
//		return 0, 0, errors.New("not a terminal")
//	}
//
//	width, height, err := terminal.GetSize(fd)
//	if err != nil {
//		return width, height, err
//	}
//
//	return width, height, nil
//}

// GetProgressBar returns an instance of ProgressBar with predefined config.
func GetProgressBar(totalSize int) *pb.ProgressBar {
	progressBar := pb.New(totalSize).SetUnits(pb.U_BYTES)
	progressBar.ShowPercent = true
	progressBar.ShowBar = true
	progressBar.ShowTimeLeft = true
	progressBar.ShowSpeed = true
	return progressBar
}

// WordWrap wraps the given string according to the provided limit with the separator sep and returns the wrapped string.
func WordWrap(s string, sep string, limit int) string {
	var wrapped string

	formattedStr := strings.Replace(s, sep, " ", -1)
	if strings.TrimSpace(formattedStr) == "" {
		return s
	}

	strSlice := strings.Fields(formattedStr)

	for len(strSlice) >= 1 {
		if wrapped == "" {
			wrapped = wrapped + strings.Join(strSlice[:limit], sep)
		} else {
			wrapped = wrapped + "\n" + sep + strings.Join(strSlice[:limit], sep)
		}

		strSlice = strSlice[limit:]

		if len(strSlice) < limit {
			limit = len(strSlice)
		}
	}

	return wrapped
}

// Upgrade checks for latest version of aws-go and downloads the latest version for the current platform, if available.
func Upgrade(version string) (error){
	gitClient := github.NewClient(nil)
	releases, _, err := gitClient.Repositories.ListReleases(context.Background(), "bharath-srinivas", "aws-go", nil)
	autoCompScript, _, _, err := gitClient.Repositories.GetContents(context.Background(), "bharath-srinivas", "aws-go", "aws_go.sh", nil)

	if err != nil {
		return err
	}

	latestRelease := releases[0]
	latestAutoCompScript := *autoCompScript.DownloadURL

	if *latestRelease.TagName == version {
		fmt.Println("aws-go is already up to date")
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("A newer version %s is available! View the release notes here: %s\n", *latestRelease.TagName,
			*latestRelease.HTMLURL)
		fmt.Print("Do you want to upgrade? [yes/no] ")
		userChoice, _ := reader.ReadString('\n')
		userChoice = strings.TrimRight(userChoice, "\n")

		switch userChoice = strings.ToLower(userChoice); userChoice {
		case "yes":
			if err := downloadRelease(latestRelease, latestAutoCompScript); err != nil {
				return err
			}
		case "no":
			break
		default:
			return errors.New("error: invalid input. please enter either \"yes\" or \"no\"")
		}

	}

	return nil
}

// downloadRelease downloads the latest version of aws-go for the current platform, if available.
// It will return error, if any.
func downloadRelease(release *github.RepositoryRelease, autoComp string) error {
	assetInfo := getAssetInfo(release)
	if assetInfo == nil {
		return errors.New("cannot find binary compatible for your system")
	}

	cmdPath, err := exec.LookPath("aws-go")
	if err != nil {
		return err
	}

	scriptPath := filepath.Join("/etc/bash_completion.d", "aws_go.sh")

	resp, err := http.Get(*assetInfo.BrowserDownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scriptResp, err := http.Get(autoComp)
	if err != nil {
		return err
	}

	cmdDir := filepath.Dir(cmdPath)
	tmpPath := filepath.Join(cmdDir, "aws-go-tmp")
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE | os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	scriptDir := filepath.Dir(scriptPath)
	tmpScriptPath := filepath.Join(scriptDir, "aws_go_latest.sh")
	tmpScriptFile, err := os.OpenFile(tmpScriptPath, os.O_CREATE | os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	fmt.Println("\nDownloading the latest version now")
	progressBar := GetProgressBar(int(resp.ContentLength))
	progressBar.Start()

	src := progressBar.NewProxyReader(resp.Body)
	if _, err := io.Copy(tmpFile, src); err != nil {
		return err
	}

	progressBar.Finish()

	scriptSrc := io.Reader(scriptResp.Body)
	if _, err := io.Copy(tmpScriptFile, scriptSrc); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, cmdPath); err != nil {
		return err
	}

	if err := os.Rename(tmpScriptPath, scriptPath); err != nil {
		return err
	}

	fmt.Println("\nVisit https://github.com/bharath-srinivas/aws-go/releases to read the changelog")

	return nil
}

// getAssetInfo returns the asset info related to the current platform.
func getAssetInfo(release *github.RepositoryRelease) *github.ReleaseAsset {
	currentBinary := "aws_go_" + runtime.GOOS + "_" + runtime.GOARCH

	for _, asset := range release.Assets {
		if *asset.Name == currentBinary {
			return &asset
		}
	}

	return nil
}