package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc"
	commandhandler "github.com/pjkaufman/go-go-gadgets/pkg/command-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	outputToFile string
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an EPUB file using EPUBCheck",
	Long: heredoc.Doc(`Validates an EPUB file using W3C EPUBCheck tool.
	If EPUBCheck is not installed, it will automatically download and install the latest version.`),
	Example: heredoc.Doc(`
	epub-lint validate -f test.epub
	will run EPUBCheck against the file specified.
`),
	Run: func(cmd *cobra.Command, args []string) {
		err := validateCommonEpubFlags(epubFile)
		if err != nil {
			logger.WriteError(err.Error())
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

		jarPath := filehandler.JoinPath(epubcheckDir, "epubcheck.jar")
		extraInputs := []string{"-jar", jarPath, epubFile}

		output := commandhandler.MustGetCommandOutputEvenIfExitError("java", "failed to run EPUBCheck", extraInputs...)

		if outputToFile != "" {
			err = filehandler.WriteFileContents(outputToFile, output)

			if err != nil {
				logger.WriteError(err.Error())
			}
		} else {
			logger.WriteInfo(output)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&epubFile, "file", "f", "", "the epub file to validate")
	err := validateCmd.MarkFlagRequired("file")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"file\" as required on validate command: %v\n", err)
	}

	err = validateCmd.MarkFlagFilename("file", "epub")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"file\" as looking for specific file types on validate command: %v\n", err)
	}

	validateCmd.Flags().StringVarP(&outputToFile, "output-file", "", "", "specifies that the validation output should be in the specified file")
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
