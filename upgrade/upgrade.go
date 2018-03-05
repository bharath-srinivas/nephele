// Package upgrade downloads the latest binary of nephele and installs it to the system.
package upgrade

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

	"github.com/bharath-srinivas/nephele/utils"
)

// Upgrade checks for latest version of nephele and downloads the latest version for the current platform, if available.
func Upgrade(version string) error {
	gitClient := github.NewClient(nil)
	releases, _, err := gitClient.Repositories.ListReleases(context.Background(), "bharath-srinivas", "nephele", nil)
	autoCompScript, _, _, err := gitClient.Repositories.GetContents(context.Background(), "bharath-srinivas", "nephele", "nephele.sh", nil)

	if err != nil {
		return err
	}

	latestRelease := releases[0]
	latestAutoCompScript := *autoCompScript.DownloadURL

	if *latestRelease.TagName == version {
		fmt.Println("nephele is already up to date")
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

// downloadRelease downloads the latest version of nephele for the current platform, if available.
// It will return error, if any.
func downloadRelease(release *github.RepositoryRelease, autoComp string) error {
	assetInfo := getAssetInfo(release)
	if assetInfo == nil {
		return errors.New("cannot find binary compatible for your system")
	}

	cmdPath, err := exec.LookPath("nephele")
	if err != nil {
		return err
	}

	scriptPath := filepath.Join("/etc/bash_completion.d", "nephele.sh")
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
	tmpPath := filepath.Join(cmdDir, "nephele-tmp")
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	scriptDir := filepath.Dir(scriptPath)
	tmpScriptPath := filepath.Join(scriptDir, "nephele_latest.sh")
	tmpScriptFile, err := os.OpenFile(tmpScriptPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	fmt.Println("\nDownloading the latest version now")
	progressBar := utils.GetProgressBar(int(resp.ContentLength))
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

	fmt.Println("\nVisit https://github.com/bharath-srinivas/nephele/releases to read the changelog")
	return nil
}

// getAssetInfo returns the asset info related to the current platform.
func getAssetInfo(release *github.RepositoryRelease) *github.ReleaseAsset {
	currentBinary := "nephele_" + runtime.GOOS + "_" + runtime.GOARCH

	for _, asset := range release.Assets {
		if *asset.Name == currentBinary {
			return &asset
		}
	}
	return nil
}
