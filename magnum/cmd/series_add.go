package cmd

import (
	"errors"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	NameArgEmpty = "name must have a non-whitespace value"
)

var (
	seriesName                     string
	seriesType                     string
	seriesPublisher                string
	slugOverride                   string
	seriesStatus                   string
	wikipediaTablesToParseOverride int
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
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateAddSeriesFlags(seriesName)
		if err != nil {
			logger.WriteError(err.Error())
		}

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

	AddCmd.Flags().StringVarP(&seriesName, "name", "n", "", "the name of the series")
	err := AddCmd.MarkFlagRequired("name")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"name\" as required on add command: %v\n", err)
	}

	AddCmd.Flags().StringVarP(&seriesPublisher, "publisher", "p", "", "the publisher of the series")
	AddCmd.Flags().StringVarP(&seriesType, "type", "t", "", "the series type")
	AddCmd.Flags().StringVarP(&slugOverride, "slug", "r", "", "the slug for the series to use instead of the one based on the series name")
	AddCmd.Flags().StringVarP(&seriesStatus, "status", "s", string(config.Ongoing), "the status of the series (defaults to Ongoing)")
	AddCmd.Flags().IntVarP(&wikipediaTablesToParseOverride, "wikipedia-table-parse-override", "o", 0, "the amount of tables that should parsed in the light novels section of the wikipedia page if it should not be all of them")
}

func ValidateAddSeriesFlags(seriesName string) error {
	if strings.TrimSpace(seriesName) == "" {
		return errors.New(NameArgEmpty)
	}

	return nil
}
