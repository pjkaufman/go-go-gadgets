package cmd

import (
	"sort"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/spf13/cobra"
)

var (
	outputFile     string
	createCsvFlags = flags.Flags{
		Flags: []flags.Flag{
			flags.NewFileFlag(false, false, &outputFile, "output", "o", "", "the file to write the csv to", []string{"csv"}, false),
		},
	}
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
		err := createFlags.Validate()
		if err != nil {
			return err
		}

		return createCsvFlags.Validate()
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

	err := createCsvFlags.AddToCmd(createCsvCmd)
	if err != nil {
		logger.WriteError(err.Error())
	}
}
