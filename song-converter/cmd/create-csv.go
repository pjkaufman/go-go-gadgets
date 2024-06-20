package cmd

import (
	"errors"
	"fmt"
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
	Use:   "create-csv",
	Short: `Creates a "|" delimited csv file that includes metadata about songs like whether they are in the church or copyrighted`,
	Example: heredoc.Doc(`To write the output of converting the files in the specified folder into a csv format to a file:
	song-converter create-csv -d working-dir -o churchSongs.csv

	To write the output of converting the files in the specified folder into a csv format to std out:
	song-converter create-csv -d working-dir
	`),
	Long: heredoc.Doc(`How it works:
	- Reads in all of the files in the specified folder.
	- Sorts the files alphabetically
	- Converts each file into a CSV row
	- Writes the content to the specified source
	`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateCreateCsvFlags(stagingDir)
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FolderMustExist(stagingDir, "working-dir")
		if err != nil {
			logger.WriteError(err.Error())
		}

		var isWritingToFile = strings.TrimSpace(coverOutputFile) == ""
		if isWritingToFile {
			logger.WriteInfo("Converting Markdown files to csv")
		}

		files, err := filehandler.MustGetAllFilesWithExtInASpecificFolder(stagingDir, ".md")
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
	SongConverterCmd.AddCommand(createCsvCmd)

	createCsvCmd.Flags().StringVarP(&stagingDir, "working-dir", "d", "", "the directory where the Markdown files are located")
	err := createCsvCmd.MarkFlagRequired("working-dir")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "working-dir" as required on create csv command: %v`, err))
	}

	err = createCsvCmd.MarkFlagDirname("working-dir")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "working-dir" as a directory on create csv command: %v`, err))
	}

	createCsvCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "the file to write the csv to")
	err = createCsvCmd.MarkFlagFilename("output-file", "csv")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "output-file" as looking for specific file types on create csv command: %v`, err))
	}
}

func ValidateCreateCsvFlags(stagingDir string) error {
	if strings.TrimSpace(stagingDir) == "" {
		return errors.New(StagingDirArgEmpty)
	}

	return nil
}
