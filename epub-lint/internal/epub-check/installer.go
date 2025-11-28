package epubcheck

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func EnsureEPUBCheckIsInstalled(epubcheckDir string) error {
	epubcheckExists, _ := filehandler.FileExists(epubcheckDir)
	if epubcheckExists {
		return nil
	}

	logger.WriteInfo("EPUBCheck not found. Installing...")

	type githubRelease struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}

	err := filehandler.CreateFolderIfNotExists(epubcheckDir)
	if err != nil {
		return err
	}

	resp, err := http.Get("https://api.github.com/repos/w3c/epubcheck/releases/latest")
	if err != nil {
		return fmt.Errorf("failed to get latest release info: %w", err)
	}
	defer resp.Body.Close()

	var release githubRelease
	if err = json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to decode release info: %w", err)
	}

	var downloadURL string
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".zip") {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("could not find EPUBCheck zip file in release %s", release.TagName)
	}

	logger.WriteInfof("Downloading EPUBCheck %s...\n", release.TagName)
	resp, err = http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download EPUBCheck: %w", err)
	}
	defer resp.Body.Close()

	tmpFile, err := filehandler.CreateTemp("", "epubcheck-*.zip")
	if err != nil {
		return err
	}

	defer filehandler.DeleteFile(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save downloaded file: %w", err)
	}

	logger.WriteInfo("Extracting EPUBCheck...")
	if err := filehandler.UnzipFile(tmpFile, epubcheckDir); err != nil {
		return fmt.Errorf("failed to extract EPUBCheck: %w", err)
	}

	folders, err := filehandler.GetFoldersInCurrentFolder(epubcheckDir)
	if err != nil {
		return err
	}

	for _, folder := range folders {
		if strings.HasPrefix(folder, "epubcheck-") {
			srcJar := filehandler.JoinPath(epubcheckDir, folder, "epubcheck.jar")
			destJar := filehandler.JoinPath(epubcheckDir, "epubcheck.jar")
			err = filehandler.Rename(srcJar, destJar)
			if err != nil {
				return fmt.Errorf("failed to move epubcheck.jar: %w", err)
			}

			srcLib := filehandler.JoinPath(epubcheckDir, folder, "lib")
			destLib := filehandler.JoinPath(epubcheckDir, "lib")
			err = filehandler.Rename(srcLib, destLib)
			if err != nil {
				return fmt.Errorf("failed to move epubcheck libraries: %w", err)
			}

			err = filehandler.DeleteFolder(filehandler.JoinPath(epubcheckDir, folder))
			if err != nil {
				return err
			}

			break
		}
	}

	logger.WriteInfo("EPUBCheck installed successfully!")

	return nil
}
