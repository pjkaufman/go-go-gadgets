package cmd

import (
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// ListCmd represents the add book info command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the names of each of the series that is currently being tracked",
	Example: heredoc.Doc(`To show a list of all series names that are being tracked:
	magnum series list

	To include information like publisher, status, series, etc.:
	magnum series list -v
	`),
	Run: func(cmd *cobra.Command, args []string) {
		seriesInfo := config.GetConfig()

		if len(seriesInfo.Series) == 0 {
			logger.WriteInfo("No series have been added to the list to keep track of.")

			return
		}

		var (
			filterOnPublisher    = strings.TrimSpace(seriesPublisher) != "" && config.IsPublisherType(seriesPublisher)
			publisherType        = config.PublisherType(seriesPublisher)
			filterOnSeriesType   = strings.TrimSpace(seriesType) != "" && config.IsSeriesType(seriesType)
			typeOfSeries         = config.SeriesType(seriesType)
			filterOnSeriesStatus = strings.TrimSpace(seriesStatus) != "" && config.IsSeriesStatus(seriesStatus)
			statusOfSeries       = config.SeriesStatus(seriesStatus)
		)
		for _, series := range seriesInfo.Series {
			if (filterOnPublisher && publisherType != series.Publisher) || (filterOnSeriesType && typeOfSeries != series.Type) || (filterOnSeriesStatus && statusOfSeries != series.Status) {
				continue
			}

			logger.WriteInfo(series.Name)
			if verbose {
				logger.WriteInfo("Status: " + config.SeriesStatusToDisplayText(series.Status))
				logger.WriteInfo("Publisher: " + string(series.Publisher))
				logger.WriteInfo("Type: " + config.SeriesTypeToDisplayText(series.Type))
				logger.WriteInfof("Total Volumes: %d\n", series.TotalVolumes)

				var slugOverride = "N/A"
				if series.SlugOverride != nil {
					slugOverride = *series.SlugOverride
				}
				logger.WriteInfo("Slug Override: " + slugOverride)

				logger.WriteInfo("")
			}
		}
	},
}

func init() {
	seriesCmd.AddCommand(ListCmd)

	ListCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show the publisher and other info about the series")
	ListCmd.Flags().StringVarP(&seriesPublisher, "publisher", "p", "", "show series with the specified publisher")
	ListCmd.Flags().StringVarP(&seriesType, "type", "t", "", "show series with the specified type")
	ListCmd.Flags().StringVarP(&seriesStatus, "status", "r", "", "show series with the specified status")
}
