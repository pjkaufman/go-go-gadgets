package cmd

import (
	"errors"
	"sort"
	"strings"

	"github.com/MakeNowJust/heredoc"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/spf13/cobra"
)

var (
	outputFile string
)

// createCsvCmd represents the createCsv command
var createCsvCmd = &cobra.Command{
	Use:   "csv",
	Short: `Creates a "|" delimited csv file that includes metadata about songs like whether they are in the church or copyrighted`,
	Example: heredoc.Doc(`To write the output of converting the files in the specified folder into a csv format to a file:
	song-converter create csv -d working-dir -o churchSongs.csv

	To write the output of converting the files in the specified folder into a csv format to std out:
	song-converter create csv -d working-dir
	`),
	Long: heredoc.Doc(`How it works:
	- Reads in all of the files in the specified folder.
	- Sorts the files alphabetically
	- Converts each file into a CSV row
	- Writes the content to the specified source
	`),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := ValidateCreateCsvFlags(stagingDir)
		if err != nil {
			return err
		}

		err = filehandler.FolderArgExists(stagingDir, "working-dir")
		if err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		var isWritingToFile = strings.TrimSpace(coverOutputFile) == ""
		if isWritingToFile {
			logger.WriteInfo("Converting Markdown files to csv")
		}

		files, err := filehandler.GetAllFilesWithExtInASpecificFolder(stagingDir, ".md")
		if err != nil {
			logger.WriteError(err.Error())
		}

		sort.Strings(files)

		var mdInfo = make([]converter.MdFileInfo, len(files))

		for i, fileName := range files {
			var filePath = filehandler.JoinPath(stagingDir, fileName)
			contents, err := filehandler.ReadInFileContents(filePath)
			if err != nil {
				logger.WriteError(err.Error())
			}

			mdInfo[i] = converter.MdFileInfo{
				FilePath:     filePath,
				FileName:     fileName,
				FileContents: contents,
			}
		}

		csvFile, err := converter.BuildCsv(mdInfo)
		if err != nil {
			logger.WriteError(err.Error())
		}

		writeToFileOrStdOut(csvFile, outputFile)

		if isWritingToFile {
			logger.WriteInfo("Finished converting Markdown files to csv")
		}
	},
}

func init() {
	createCmd.AddCommand(createCsvCmd)

	createCsvCmd.Flags().StringVarP(&outputFile, "output", "o", "", "the file to write the csv to")
	err := createCsvCmd.MarkFlagFilename("output", "csv")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"output\" as looking for specific file types on create csv command: %v\n", err)
	}
}

func ValidateCreateCsvFlags(stagingDir string) error {
	if strings.TrimSpace(stagingDir) == "" {
		return errors.New(StagingDirArgEmpty)
	}

	return nil
}
