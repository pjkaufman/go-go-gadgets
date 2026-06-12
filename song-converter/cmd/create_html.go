package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/spf13/cobra"
)

const (
	CoverPathArgEmpty  = "cover-file must have a non-whitespace value"
	CoverPathNotMdFile = "cover-file must be an md file"
	StagingDirArgEmpty = "working-dir must have a non-whitespace value"
	fileFormat         = `<html>
  <body>
    <section id="cover">
      %s
    </section>
    <section id="contents">
    <h1 class="toc">Table of Contents</h1>
    <ul>
        %s
    </ul>
    </section>
    <section id="songs">
      %s
    </section>
  </body>
</html>`
)

var (
	stagingDir         string
	bodyHtmlOutputFile string
	coverOutputFile    string
	coverInputFilePath string
	versionDescriptor  string
	location           string
	commonBookFlags    = flags.Flags{
		Flags: []flags.Flag{
			flags.NewFileFlag(true, false, &coverInputFilePath, "cover-file", "c", "", "the markdown cover file to use", []string{"md"}, true),
			flags.NewFileFlag(false, false, &bodyHtmlOutputFile, "output", "o", "", "the html file to write the output to", []string{"html"}, false),
		},
	}
	createHtmlFlags = flags.Flags{
		Flags: []flags.Flag{
			flags.NewEnumFlag(true, false, &versionDescriptor, "format", "", "", "the version descriptor for the type of songs to generate (Abridged or Unabridged)", []string{"Abridged", "Unabridged"}),
		},
	}
)

// createHtmlCmd represents the CreateSongs command
var createHtmlCmd = &cobra.Command{
	Use:   "html",
	Short: "Converts the cover and all Markdown files in the specified folder into html in alphabetical order generating three sections: the cover, table of contents, and songs",
	Example: heredoc.Doc(`To write the output of converting the files in the specified folder to html to a file:
	song-converter create html -d working-dir -c cover.md -o songs.html

	To write the output of converting the files in the specified folder to html to std out:
	song-converter create html -d working-dir -c cover.md
	`),
	Long: heredoc.Doc(`How it works:
	- Reads in all of the files in the specified folder
	- Sorts the files alphabetically
	- Adds the cover to the start of the content after converting it to html
	- Converts each file into html
	- Writes the content to the specified source
	`),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := validateCreateHtmlFile()
		if err != nil {
			return err
		}

		return createHtmlFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		createHtmlFile(stagingDir, coverInputFilePath, coverOutputFile, bodyHtmlOutputFile, versionDescriptor, "", false)
	},
}

func createHtmlFile(stagingDir, coverInputFilePath, coverOutputFile, bodyHtmlOutputFile, bookType, extraCss string, isBook bool) {
	var isWritingToFile = strings.TrimSpace(coverOutputFile) == ""
	if isWritingToFile {
		logger.WriteInfo("Converting file to html cover")
	}
	coverMd, err := filehandler.ReadInFileContents(coverInputFilePath)
	if err != nil {
		logger.WriteError(err.Error())
	}

	coverHtml := converter.BuildHtmlCover(coverMd, bookType, extraCss, time.Now())

	if isWritingToFile {
		logger.WriteInfo("Finished creating html cover file")
	}

	if isWritingToFile {
		logger.WriteInfo("Converting Markdown files to html")
	}

	files, err := filehandler.GetAllFilesWithExtInASpecificFolder(stagingDir, ".md")
	if err != nil {
		logger.WriteError(err.Error())
	}

	sort.Strings(files)

	var mdInfo = make([]converter.MdFileInfo, len(files))

	for i, fileName := range files {
		var filePath = filehandler.JoinPath(stagingDir, fileName)
		fileContents, err := filehandler.ReadInFileContents(filePath)
		if err != nil {
			logger.WriteError(err.Error())
		}

		mdInfo[i] = converter.MdFileInfo{
			FilePath:     filePath,
			FileName:     fileName,
			FileContents: fileContents,
		}
	}

	var converterType = converter.Digital
	if isBook {
		mdInfo, err = converter.FilterAndSortSongs(mdInfo, location)
		if err != nil {
			logger.WriteError(err.Error())
		}

		converterType = converter.Book
	}

	songsHtml, headerIds, err := converter.BuildHtmlSongs(mdInfo, converterType)
	if err != nil {
		logger.WriteError(err.Error())
	}

	if isBook {
		writeToFileOrStdOut(fmt.Sprintf(bookFormat, coverHtml, buildBookListItems(mdInfo), songsHtml), bodyHtmlOutputFile)
	} else {
		writeToFileOrStdOut(fmt.Sprintf(fileFormat, coverHtml, buildListItems(headerIds), songsHtml), bodyHtmlOutputFile)
	}

	if isWritingToFile {
		logger.WriteInfo("Finished converting Markdown files to html")
	}
}

func validateCreateHtmlFile() error {
	err := createFlags.Validate()
	if err != nil {
		return err
	}

	err = commonBookFlags.Validate()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	createCmd.AddCommand(createHtmlCmd)

	err := commonBookFlags.AddToCmd(createHtmlCmd)
	if err != nil {
		logger.WriteError(err.Error())
	}

	err = createHtmlFlags.AddToCmd(createHtmlCmd)
	if err != nil {
		logger.WriteError(err.Error())
	}
}

func buildListItems(headerIds []string) string {
	if len(headerIds) == 0 {
		return ""
	}

	var listItems = strings.Builder{}
	for _, headerId := range headerIds {
		fmt.Fprintf(&listItems, `<li><a href="#%s"></a></li>`+"\n", headerId)
	}

	return listItems.String()
}
