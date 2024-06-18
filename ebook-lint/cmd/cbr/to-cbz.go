package cbr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var dir string

const (
	DirArgEmpty = "directory must have a non-whitespace value"
)

// cbrToCbzCmd represents the toCbz command
var cbrToCbzCmd = &cobra.Command{
	Use:   "to-cbz",
	Short: "Converts all of the cbr files to cbz files in the specified directory.",
	Example: heredoc.Doc(`To convert all cbrs to cbzs in a folder:
	ebook-lint cbr to-cbz -d folder
	
	To convert all cbrs to cbzs in the current directory:
	ebook-lint cbr to-cbz 
	`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateToCbrFlags(dir)
		if err != nil {
			logger.WriteError(err.Error())
		}

		logger.WriteInfo("Starting converting cbr files to cbz files\n")

		cbrs, err := filehandler.MustGetAllFilesWithExtInASpecificFolder(dir, ".cbr")
		if err != nil {
			logger.WriteError(err.Error())
		}

		for _, cbr := range cbrs {
			logger.WriteInfo(fmt.Sprintf("starting to convert %s to a cbz file...", cbr))

			err = filehandler.ConvertRarToCbz(cbr)
			if err != nil {
				logger.WriteError(err.Error())
			}
		}

		logger.WriteInfo("\nFinished converting cbr files to cbz files")
	},
}

func init() {
	CbrCmd.AddCommand(cbrToCbzCmd)

	cbrToCbzCmd.Flags().StringVarP(&dir, "directory", "d", ".", "the folder where all cbr files should be converted to cbz files")
}

func ValidateToCbrFlags(dir string) error {
	if strings.TrimSpace(dir) == "" {
		return errors.New(DirArgEmpty)
	}

	return nil
}
