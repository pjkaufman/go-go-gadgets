package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// ShowInfoCmd represents the add book info command
var ShowInfoCmd = &cobra.Command{
	Use:   "show-info",
	Short: "Shows each series that has upcoming releases along with when the releases are in the order they are going to be released",
	Example: heredoc.Doc(`To show upcoming releases in order of when they are releasing:
	magnum show-info
	`),
	Run: func(cmd *cobra.Command, args []string) {
		seriesInfo := config.GetConfig()

		if len(seriesInfo.Series) == 0 {
			logger.WriteInfo("No series have been added to the list to keep track of.")

			return
		}

		var unreleasedVolumes []config.ReleaseInfo
		for _, series := range seriesInfo.Series {
			if len(series.UnreleasedVolumes) == 0 {
				continue
			}

			for i, unreleasedVolume := range series.UnreleasedVolumes {
				if unreleasedVolume.ReleaseDate == defaultReleaseDate {
					continue
				}

				if strings.HasPrefix(unreleasedVolume.Name, "Vol") {
					series.UnreleasedVolumes[i].Name = series.Name + ": " + unreleasedVolume.Name
				}

				unreleasedVolumes = append(unreleasedVolumes, series.UnreleasedVolumes[i])
			}
		}

		if len(unreleasedVolumes) == 0 {
			logger.WriteInfo("No release are upcoming")
			return
		}

		logger.WriteInfo("Upcoming releases:")
		logger.WriteInfo("")
		sort.Slice(unreleasedVolumes, func(i, j int) bool {
			if unreleasedVolumes[i].ReleaseDate == defaultReleaseDate {
				return false
			} else if unreleasedVolumes[j].ReleaseDate == defaultReleaseDate {
				return true
			}

			date1 := parseVolumeReleaseDate(unreleasedVolumes[i].Name, unreleasedVolumes[i].ReleaseDate)

			date2 := parseVolumeReleaseDate(unreleasedVolumes[j].Name, unreleasedVolumes[j].ReleaseDate)

			return date1.Before(date2)
		})

		var today = time.Now()
		var oneWeekAgo = today.AddDate(0, 0, -7)
		var nextMonth = today.AddDate(0, 1, 0)
		for _, unreleasedVolume := range unreleasedVolumes {
			var displayText = getUnreleasedVolumeDisplayText(unreleasedVolume.Name, unreleasedVolume.ReleaseDate)
			if unreleasedVolume.ReleaseDate == defaultReleaseDate {
				logger.WriteInfo(displayText)
				continue
			}

			date := parseVolumeReleaseDate(unreleasedVolume.Name, unreleasedVolume.ReleaseDate)
			if date.Before(oneWeekAgo) {
				logger.WriteInfoWithColor(displayText, color.FgRed)
			} else if date.Before(nextMonth) {
				logger.WriteInfoWithColor(displayText, color.FgYellow)
			} else {
				logger.WriteInfo(displayText)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(ShowInfoCmd)
}

func parseVolumeReleaseDate(name, releaseDate string) time.Time {
	date, err := time.Parse(releaseDateFormat, releaseDate)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to parse release date \"%s\" for \"%s\": %s", name, releaseDate, err))
	}

	return date
}
