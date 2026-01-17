package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/spf13/cobra"
)

const (
	bookFormat = `<html>
  <body>
    <section id="cover">
      %s
    </section>
    <section id="contents">
    <h1 class="toc">Index</h1>
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

// CreateBookCmd represents the CreateSongs command
var CreateBookCmd = &cobra.Command{
	Use: "book",
	// Short: "Converts the cover and all Markdown files in the specified folder into html in alphabetical order generating three sections: the cover, table of contents, and songs",
	// Example: heredoc.Doc(`To write the output of converting the files in the specified folder to html to a file:
	// song-converter create-html -d working-dir -c cover.md -o songs.html

	// To write the output of converting the files in the specified folder to html to std out:
	// song-converter create-html -d working-dir -s cover.md
	// `),
	// Long: heredoc.Doc(`How it works:
	// - Reads in all of the files in the specified folder
	// - Sorts the files alphabetically
	// - Adds the cover to the start of the content after converting it to html
	// - Converts each file into html
	// - Writes the content to the specified source
	// `),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: refactor most of this logic into a general function that can be used for create book and create html
		err := ValidateCreateBookFlags(stagingDir, coverInputFilePath)
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FolderArgExists(stagingDir, "working-dir")
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FileArgExists(coverInputFilePath, "cover-file")
		if err != nil {
			logger.WriteError(err.Error())
		}

		var isWritingToFile = strings.TrimSpace(coverOutputFile) == ""
		if isWritingToFile {
			logger.WriteInfo("Converting file to html cover")
		}
		coverMd, err := filehandler.ReadInFileContents(coverInputFilePath)
		if err != nil {
			logger.WriteError(err.Error())
		}

		coverHtml := converter.BuildHtmlCover(coverMd, "", "font-size: 52pt;", time.Now())

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

		mdInfo, err = converter.FilterAndSortSongs(mdInfo, location)
		if err != nil {
			logger.WriteError(err.Error())
		}

		songsHtml, _, err := converter.BuildHtmlSongs(mdInfo, converter.Book)
		if err != nil {
			logger.WriteError(err.Error())
		}

		writeToFileOrStdOut(fmt.Sprintf(bookFormat, coverHtml, buildBookListItems(mdInfo), songsHtml), bodyHtmlOutputFile)

		if isWritingToFile {
			logger.WriteInfo("Finished converting Markdown files to html")
		}
	},
}

func init() {
	createCmd.AddCommand(CreateBookCmd)

	CreateBookCmd.Flags().StringVarP(&stagingDir, "working-dir", "d", "", "the directory where the Markdown files are located")
	err := CreateBookCmd.MarkFlagRequired("working-dir")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"working-dir\" as required on create book command: %v\n", err)
	}

	err = CreateBookCmd.MarkFlagDirname("working-dir")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"working-dir\" as a directory on create book command: %v\n", err)
	}

	CreateBookCmd.Flags().StringVarP(&coverInputFilePath, "cover-file", "c", "", "the markdown cover file to use")
	err = CreateBookCmd.MarkFlagRequired("cover-file")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"cover-file\" as required on create book command: %v\n", err)
	}

	err = CreateBookCmd.MarkFlagFilename("cover-file", "md")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"cover-file\" as a looking for specific file types on create book command: %v\n", err)
	}

	CreateBookCmd.Flags().StringVarP(&bodyHtmlOutputFile, "output", "o", "", "the html file to write the output to")
	err = CreateBookCmd.MarkFlagFilename("output", "html")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"output\" as a looking for specific file types on create book command: %v\n", err)
	}

	CreateBookCmd.Flags().StringVarP(&location, "location", "l", "", "the specific book to recreate by filtering songs down to just that book location")
	err = CreateBookCmd.MarkFlagRequired("location")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"location\" as required on create book command: %v\n", err)
	}
	// TODO: add a sort type for alphabetical and a sort type for in order for the TOC order...
}

func ValidateCreateBookFlags(stagingDir, coverInputFilePath string) error {
	if strings.TrimSpace(stagingDir) == "" {
		return errors.New(StagingDirArgEmpty)
	}

	if strings.TrimSpace(coverInputFilePath) == "" {
		return errors.New(CoverPathArgEmpty)
	}

	if !strings.HasSuffix(coverInputFilePath, ".md") {
		return errors.New(CoverPathNotMdFile)
	}

	return nil
}

// TODO: this needs to take in option for putting the songs in page order or putting them in alphabetical order...
// Note: that alphabetical order will not be perfect given the discrepancy between some of the names in the digital vs. book versions
func buildBookListItems(headerInfo []converter.MdFileInfo) string {
	if len(headerInfo) == 0 {
		return ""
	}

	var (
		pageNumberIndex = make(map[string]int)
		listItems       = strings.Builder{}
		pageNumber      int
	)
	for _, headerData := range headerInfo {
		if val, ok := pageNumberIndex[headerData.FileName]; ok {
			pageNumber = headerData.PageNumbers[val]
		} else {
			pageNumber = headerData.PageNumbers[0]
			pageNumberIndex[headerData.FileName] = 1
		}

		listItems.WriteString(fmt.Sprintf(`<li><span class="name">%s</span><span class="page">%d</span></li>`+"\n", headerData.Header, pageNumber))
	}

	return listItems.String()
}
