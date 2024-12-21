package epub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	commandhandler "github.com/pjkaufman/go-go-gadgets/pkg/command-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [epub-file]",
	Short: "Validate an EPUB file using EPUBCheck",
	Long: `Validates an EPUB file using W3C EPUBCheck tool. 
If EPUBCheck is not installed, it will automatically download and install the latest version.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		epubFile := args[0]

		if _, err := os.Stat(epubFile); os.IsNotExist(err) {
			logger.WriteErrorf("EPUB file does not exist: %s", epubFile)
		}

		epubcheckDir, err := filehandler.GetDataDir("epubcheck")
		if err != nil {
			logger.WriteError(err.Error())
		}

		epubcheckExists, _ := filehandler.FileExists(epubcheckDir)
		if !epubcheckExists {
			logger.WriteInfo("EPUBCheck not found. Installing...")
			if err := downloadEPUBCheck(epubcheckDir); err != nil {
				logger.WriteError(err.Error())
			}
		}

		jarPath := filepath.Join(epubcheckDir, "epubcheck.jar")
		output := commandhandler.MustGetCommandOutputEvenIfExitError("java", "failed to run EPUBCheck", "-jar", jarPath, epubFile)

		logger.WriteInfo(output)
	},
}

func init() {
	EpubCmd.AddCommand(validateCmd)
}

func downloadEPUBCheck(epubcheckDir string) error {
	type githubRelease struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"` // Fixed the json tag name
		} `json:"assets"`
	}

	err := filehandler.CreateFolderIfNotExists(epubcheckDir)
	if err != nil {
		return err
	}

	// Get latest release info from GitHub API
	resp, err := http.Get("https://api.github.com/repos/w3c/epubcheck/releases/latest")
	if err != nil {
		return fmt.Errorf("failed to get latest release info: %w", err)
	}
	defer resp.Body.Close()

	var release githubRelease
	if err = json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to decode release info: %w", err)
	}

	// Find the .zip asset
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

	// Download the zip file
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

	// Save zip file
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save downloaded file: %w", err)
	}

	// Unzip the file using our Go implementation
	logger.WriteInfo("Extracting EPUBCheck...")
	if err := filehandler.UnzipFile(tmpFile, epubcheckDir); err != nil {
		return fmt.Errorf("failed to extract EPUBCheck: %w", err)
	}

	// Find and move the jar file to the right location
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

			// Copy library files as well
			srcLib := filehandler.JoinPath(epubcheckDir, folder, "lib")
			destLib := filehandler.JoinPath(epubcheckDir, "lib")
			err = filehandler.Rename(srcLib, destLib)
			if err != nil {
				return fmt.Errorf("failed to move epubcheck libraries: %w", err)
			}

			// Clean up the extracted directory
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
