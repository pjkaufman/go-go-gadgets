package cmd

import (
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	interactiveMode bool
	bookStatus      string
	seriesEditFlags = flags.Flags{
		Flags: []flags.Flag{
			flags.NewStringFlag(false, false, &seriesName, "name", "n", "", "the name of the series to edit"),
			flags.NewEnumFlag(false, false, &seriesPublisher, "publisher", "p", "", "the publisher of the series", config.AllPublisherTypes()),
			flags.NewEnumFlag(false, false, &seriesType, "type", "t", "", "the series type", config.AllSeriesTypes()),
			flags.NewEnumFlag(false, false, &seriesStatus, "status", "s", "", "status to set for the selected book (O/H/C)", config.AllStatuses()),
			flags.NewStringFlag(false, false, &slugOverride, "slug", "r", "", "the slug for the series to use instead of the one based on the series name"),
			flags.NewIntFlag(false, false, &wikipediaTablesToParseOverride, "wikipedia-table-parse-override", "o", 0, "the amount of tables that should parsed in the light novels section of the wikipedia page if it should not be all of them"),
			flags.NewBoolFlag(false, false, &interactiveMode, "interactive", "i", false, "gets the name, publisher, series type, and series status interactively when not provided"),
			flags.NewBoolFlag(false, false, &includeCompleted, "include-completed", "c", false, "include completed series in the books to search"),
		},
	}
)

// EditCmd represents the edit series command
var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edits the provided/selected book name",
	Example: heredoc.Doc(`To set the status, publisher, and/or series type of a series:
	magnum series edit -n "series_name" -i
	This will result in being prompted for a status, publisher, and series type for the series.

	To set the status of a series you know the name and status of without wanting to use any interactive prompt:
	magnum series edit -n "series_name" -s C
	`),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return seriesEditFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		seriesInfo := config.GetConfig()

		if len(seriesInfo.Series) == 0 {
			logger.WriteInfo("No series have been added to the list to keep track of.")

			return
		}

		var name = seriesName
		if strings.TrimSpace(name) == "" {
			if !interactiveMode {
				logger.WriteFatal("No series name was provided and interactive was not specified, so no series change could be made.")
			}

			name = selectBookName(seriesInfo.Series, includeCompleted)

			logger.WriteInfof("%q selected\n", name)
		}

		if !seriesInfo.HasSeries(name) {
			logger.WriteWarnf("No series with the name %q is in the series list.\n", name)
			return
		}

		var (
			foundSeriesToUpdate = false
			updatedSeries       config.SeriesInfo
			indexToUpdate       int
		)
		for i, series := range seriesInfo.Series {
			if name == series.Name {
				foundSeriesToUpdate = true
				updatedSeries = series
				indexToUpdate = i
				break
			}
		}

		if !foundSeriesToUpdate {
			logger.WriteFatalf("No series with the name %q is in the series list.\n", name)
		}

		var changeMade bool
		var (
			status       = config.SeriesStatus(bookStatus)
			updateStatus = interactiveMode || cmd.Flags().Changed("status")
		)
		if strings.TrimSpace(bookStatus) == "" && interactiveMode {
			status = selectBookStatus(&updatedSeries.Status)
			bookStatus = string(status)

			logger.WriteInfof("%q status selected\n", config.SeriesStatusToDisplayText(status))
		}

		if updateStatus && status != updatedSeries.Status {
			if !config.IsSeriesStatus(bookStatus) {
				logger.WriteWarnf("Status %q is not a valid book status, so it is being ignored\n", bookStatus)

				bookStatus = ""
			} else {
				updatedSeries.Status = status
				changeMade = true
			}
		}

		var (
			publisher       = config.PublisherType(seriesPublisher)
			updatePublisher = interactiveMode || cmd.Flags().Changed("publisher")
		)
		if strings.TrimSpace(seriesPublisher) == "" && interactiveMode {
			publisher = selectPublisher(&updatedSeries.Publisher)
			seriesPublisher = string(publisher)

			logger.WriteInfof("%q publisher type selected\n", config.PublisherToDisplayString(publisher))
		}

		if updatePublisher && publisher != updatedSeries.Publisher {
			if !config.IsPublisherType(seriesPublisher) {
				logger.WriteWarnf("Publisher %q is not a valid book publisher, so it is being ignored\n", seriesPublisher)

				seriesPublisher = ""
			} else {
				updatedSeries.Publisher = publisher
				changeMade = true
			}
		}

		var (
			typeOfSeries = config.SeriesType(seriesType)
			updateType   = interactiveMode || cmd.Flags().Changed("type")
		)
		if strings.TrimSpace(seriesType) == "" && interactiveMode {
			typeOfSeries = selectSeriesType(&updatedSeries.Type)
			seriesType = string(typeOfSeries)

			logger.WriteInfof("%q series type selected\n", config.SeriesTypeToDisplayText(typeOfSeries))
		}

		if updateType && typeOfSeries != updatedSeries.Type {
			if !config.IsSeriesType(seriesType) {
				logger.WriteWarnf("Series type %q is not a valid book type, so it is being ignored\n", seriesType)

				seriesType = ""
			} else {
				updatedSeries.Type = typeOfSeries
				changeMade = true
			}
		}

		if cmd.Flags().Changed("slug") {
			if strings.TrimSpace(slugOverride) != "" {
				updatedSeries.SlugOverride = &slugOverride
			} else {
				updatedSeries.SlugOverride = nil
			}

			changeMade = true
		}

		if cmd.Flags().Changed("wikipedia-table-parse-override") {
			if wikipediaTablesToParseOverride > 0 {
				updatedSeries.WikipediaTablesToParseOverride = &wikipediaTablesToParseOverride
			} else {
				updatedSeries.WikipediaTablesToParseOverride = nil
			}

			changeMade = true
		}

		if changeMade {
			seriesInfo.Series[indexToUpdate] = updatedSeries

			config.WriteConfig(seriesInfo)
			logger.WriteInfof("Successfully updated %q.\n", name)
		} else {
			logger.WriteInfof("No changes made for %q.\n", name)
		}
	},
}

func init() {
	seriesCmd.AddCommand(EditCmd)

	err := seriesEditFlags.AddToCmd(EditCmd)
	if err != nil {
		logger.WriteFatal(err.Error())
	}
}
