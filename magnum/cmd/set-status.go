package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	includeCompleted bool
	bookName         string
	bookStatus       string
)

// SetStatus represents the set book status command
var SetStatus = &cobra.Command{
	Use:   "set-status",
	Short: "Sets the status of the provided/selected book name",
	Example: heredoc.Doc(`To set the status of a book you know the name of:
	magnum set-status -n "book_name"
	This will result in being prompted for a status for that book.

	To set the status of a book you know the name and status of:
	magnum set-status -n "book_name" -s C

	To set the status of a book by using the cli selection options:
	magnum set-status

	To set the status of a book and include the completed series:
	magnum set-status -c
	`),
	Run: func(cmd *cobra.Command, args []string) {
		seriesInfo := config.GetConfig()

		if len(seriesInfo.Series) == 0 {
			logger.WriteInfo("No series have been added to the list to keep track of.")

			return
		}

		var name = bookName
		if strings.TrimSpace(name) == "" {
			name = selectBookName(seriesInfo.Series, false)

			logger.WriteInfo(fmt.Sprintf("\"%s\" selected", name))
		}

		var status = config.SeriesStatus(bookStatus)
		if !config.IsSeriesStatus(bookStatus) {
			logger.WriteWarn(fmt.Sprintf(`Status "%s" is not a valid book status, so it is being ignored`, bookStatus))

			bookStatus = ""
		}

		logger.WriteInfo(bookStatus)
		if strings.TrimSpace(bookStatus) == "" {
			status = selectBookStatus()

			logger.WriteInfo(fmt.Sprintf("\"%s\" selected", status))
		}

		var foundSeriesToUpdate = false
		for i, series := range seriesInfo.Series {
			if name == series.Name {
				foundSeriesToUpdate = true
				seriesInfo.Series[i].Status = status
				break
			}
		}

		if !foundSeriesToUpdate {
			logger.WriteError(fmt.Sprintf("\n"+`Failed to find "%s" to set the status to %s.`, seriesName, status))
		}

		config.WriteConfig(seriesInfo)

		logger.WriteInfo(fmt.Sprintf("\n"+`Successfully set "%s" to have a status of %s.`, name, status))
	},
}

func init() {
	rootCmd.AddCommand(SetStatus)

	SetStatus.Flags().BoolVarP(&includeCompleted, "include-completed", "c", false, "include completed series in the books to search")
	SetStatus.Flags().StringVarP(&bookName, "name", "n", "", "name of the book to set the status for")
	SetStatus.Flags().StringVarP(&bookStatus, "status", "s", "", "status to set for the selected book (O/H/C)")
}
