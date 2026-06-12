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
	seriesName                     string
	seriesType                     string
	seriesPublisher                string
	slugOverride                   string
	seriesStatus                   string
	wikipediaTablesToParseOverride int
	seriesAddFlags                 = flags.Flags{
		Flags: []flags.Flag{
			flags.NewStringFlag(true, false, &seriesName, "name", "n", "", "the name of the series"),
			flags.NewEnumFlag(false, false, &seriesPublisher, "publisher", "p", "", "the publisher of the series", config.AllPublisherTypes()),
			flags.NewEnumFlag(false, false, &seriesType, "type", "t", "", "the series type", config.AllSeriesTypes()),
			flags.NewStringFlag(false, false, &slugOverride, "slug", "r", "", "the slug for the series to use instead of the one based on the series name"),
			flags.NewEnumFlag(false, false, &seriesStatus, "status", "s", string(config.Ongoing), "the status of the series (defaults to Ongoing)", config.AllStatuses()),
			flags.NewIntFlag(false, false, &wikipediaTablesToParseOverride, "wikipedia-table-parse-override", "", 0, "the amount of tables that should parsed in the light novels section of the wikipedia page if it should not be all of them"),
		},
	}
)

// AddCmd represents the add book info command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds the provided series info to the list of series to keep track of",
	Example: heredoc.Doc(`To add a series with just a name and other information to be filled out:
	magnum series add -n "Lady and the Tramp"
	Note: that the other fields will be filled in via prompts except the series status which is assumed to be ongoing

	To add a series with a special URL slug that does not follow the normal pattern for the publisher in question or is on its own page:
	magnum series add -n "Re:ZERO -Starting Life in Another World" -s "re-starting-life-in-another-world"

	To add a series that is not ongoing (for example Completed):
	magnum series add -n "Demon Slayer" -r "C"
	`),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return seriesAddFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		seriesInfo := config.GetConfig()
		if seriesInfo.HasSeries(seriesName) {
			logger.WriteInfo("The series already exists in the list.")

			return
		}

		var publisher = config.PublisherType(seriesPublisher)
		if strings.TrimSpace(seriesPublisher) == "" || !config.IsPublisherType(seriesPublisher) {
			publisher = selectPublisher(nil)
		}

		var typeOfSeries = config.SeriesType(seriesType)
		if strings.TrimSpace(seriesType) == "" || !config.IsSeriesType(seriesType) {
			typeOfSeries = selectSeriesType(nil)
		}

		var status = config.SeriesStatus(seriesStatus)
		if strings.TrimSpace(seriesStatus) == "" || !config.IsSeriesStatus(seriesStatus) {
			status = selectBookStatus(nil)
		}

		var override *string
		if strings.TrimSpace(slugOverride) != "" {
			override = &slugOverride
		}

		newSeries := config.SeriesInfo{
			Name:         seriesName,
			Publisher:    publisher,
			Type:         typeOfSeries,
			SlugOverride: override,
			Status:       status,
		}

		warning := seriesInfo.AddSeries(newSeries, wikipediaTablesToParseOverride)
		if warning != "" {
			logger.WriteWarn(warning)
		}

		config.WriteConfig(seriesInfo)
	},
}

func init() {
	seriesCmd.AddCommand(AddCmd)

	err := seriesAddFlags.AddToCmd(AddCmd)
	if err != nil {
		logger.WriteError(err.Error())
	}
}
