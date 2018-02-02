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
	"gopkg.in/cheggaaa/pb.v1"
)

// GetProgressBar returns an instance of ProgressBar with predefined config.
func GetProgressBar(totalSize int) *pb.ProgressBar {
	progressBar := pb.New(totalSize).SetUnits(pb.U_BYTES)
	progressBar.ShowPercent = true
	progressBar.ShowBar = true
	progressBar.ShowTimeLeft = true
	progressBar.ShowSpeed = true
	return progressBar
}

// WordWrap wraps the given string according to the provided parts with the separator sep and returns the wrapped
// string if and only if the given string has the character `.` or `-`. Currently WordWrap is very naive and it'll
// break the string if the separator position is greater than the half length of the provided string. It has been
// written solely for the purpose of wrapping text for rendering in table writer and not recommended for normal use.
func WordWrap(s string, sep byte, parts int) string {
	var wrapped []byte

	if parts <= 0 || !hasSeparator(s) {
		return s
	}

	halfLength := len(s) / parts
	if halfLength <= 10 {
		return s
	}

	broken := false
	for i, char := range s {
		wrapped = append(wrapped, byte(char))
		if char == rune(sep) && !broken {
			if i >= halfLength {
				wrapped = append(wrapped, byte('\n'))
				broken = true
			}
		}
	}

	return string(wrapped)
}

// Upgrade checks for latest version of aws-go and downloads the latest version for the current platform, if available.
func Upgrade(version string) error {
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
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	scriptDir := filepath.Dir(scriptPath)
	tmpScriptPath := filepath.Join(scriptDir, "aws_go_latest.sh")
	tmpScriptFile, err := os.OpenFile(tmpScriptPath, os.O_CREATE|os.O_RDWR, 0755)
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

// hasSeparator is a helper function for WordWrap which will return true if the given string has any one of the
// `.` or `-` separator.
func hasSeparator(s string) bool {
	for _, c := range s {
		if c == '.' || c == '-' {
			return true
		}
	}
	return false
}
