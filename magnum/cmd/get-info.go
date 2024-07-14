package cmd

import (
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	jnovelclub "github.com/pjkaufman/go-go-gadgets/magnum/internal/jnovel-club"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/sevenseasentertainment"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/vizmedia"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/wikipedia"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/yenpress"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// GetInfoCmd represents the get book info command
var GetInfoCmd = &cobra.Command{
	Use:   "get-info",
	Short: "Gets the book release info for books that have been added to the list of series to track",
	Example: heredoc.Doc(`To get all of the release data for non-completed series:
	magnum get-info`),
	Run: func(cmd *cobra.Command, args []string) {
		seriesInfo := config.GetConfig()

		for i, series := range seriesInfo.Series {
			if series.Status != config.Completed || includeCompleted {
				seriesInfo.Series[i] = getSeriesVolumeInfo(series)
			}
		}

		config.WriteConfig(seriesInfo)
	},
}

func init() {
	rootCmd.AddCommand(GetInfoCmd)

	GetInfoCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show more info about what is going on")
	GetInfoCmd.Flags().BoolVarP(&includeCompleted, "include-completed", "c", false, "get info for completed series")
}

func getSeriesVolumeInfo(seriesInfo config.SeriesInfo) config.SeriesInfo {
	logger.WriteInfof("Checking for volume info for %q\n", seriesInfo.Name)

	switch seriesInfo.Publisher {
	case config.YenPress:
		return yenPressGetSeriesVolumeInfo(seriesInfo)
	case config.JNovelClub:
		return jNovelClubGetSeriesVolumeInfo(seriesInfo)
	case config.SevenSeasEntertainment:
		return sevenSeasEntertainmentGetSeriesVolumeInfo(seriesInfo)
	case config.OnePeaceBooks, config.HanashiMedia:
		return wikipediaGetSeriesVolumeInfo(seriesInfo)
	case config.VizMedia:
		return vizMediaGetSeriesVolumeInfo(seriesInfo)
	default:
		return seriesInfo
	}
}

func yenPressGetSeriesVolumeInfo(seriesInfo config.SeriesInfo) config.SeriesInfo {
	volumes, numVolumes := yenpress.GetVolumes(seriesInfo.Name, seriesInfo.SlugOverride, verbose)

	if len(volumes) == 0 {
		logger.WriteInfo("The yen press light novels do not exist for this series.")

		return seriesInfo
	}

	if numVolumes == seriesInfo.TotalVolumes {
		return handleNoChangeDisplayAndSeriesInfoUpdates(seriesInfo)
	}

	var today = time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	var unreleasedVolumes = []string{}
	var releaseDateInfo = []string{}
	for _, info := range volumes {
		releaseDate := yenpress.GetReleaseDateInfo(info, verbose)

		if releaseDate != nil {
			if releaseDate.Before(today) {
				break
			} else {
				releaseDateInfo = append(releaseDateInfo, releaseDate.Format("January 2, 2006"))
				unreleasedVolumes = append(unreleasedVolumes, info.Name)
			}
		}
	}

	return printReleaseInfoAndUpdateSeriesInfo(seriesInfo, unreleasedVolumes, releaseDateInfo, numVolumes, volumes[0].Name)
}

func jNovelClubGetSeriesVolumeInfo(seriesInfo config.SeriesInfo) config.SeriesInfo {
	volumeInfo := jnovelclub.GetVolumeInfo(seriesInfo.Name, seriesInfo.SlugOverride, verbose)

	if len(volumeInfo) == 0 {
		logger.WriteInfo("The jnovel club light novels do not exist for this series.")

		return seriesInfo
	}

	if len(volumeInfo) == seriesInfo.TotalVolumes {
		return handleNoChangeDisplayAndSeriesInfoUpdates(seriesInfo)
	}

	var today = time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	var unreleasedVolumes = []string{}
	var releaseDateInfo = []string{}
	for _, info := range volumeInfo {
		if info.ReleaseDate.Before(today) {
			break
		} else {
			releaseDateInfo = append(releaseDateInfo, info.ReleaseDate.Format(releaseDateFormat))
			unreleasedVolumes = append(unreleasedVolumes, info.Name)
		}

	}

	return printReleaseInfoAndUpdateSeriesInfo(seriesInfo, unreleasedVolumes, releaseDateInfo, len(volumeInfo), volumeInfo[0].Name)
}

func wikipediaGetSeriesVolumeInfo(seriesInfo config.SeriesInfo) config.SeriesInfo {
	volumeInfo := wikipedia.GetVolumeInfo(userAgent, seriesInfo.Name, seriesInfo.SlugOverride, seriesInfo.WikipediaTablesToParseOverride, verbose)

	if len(volumeInfo) == 0 {
		logger.WriteInfo("The wikipedia light novels do not exist for this series.")

		return seriesInfo
	}

	if len(volumeInfo) == seriesInfo.TotalVolumes && (len(seriesInfo.UnreleasedVolumes) == 0 || seriesInfo.UnreleasedVolumes[0].ReleaseDate != defaultReleaseDate) {
		return handleNoChangeDisplayAndSeriesInfoUpdates(seriesInfo)
	}

	var today = time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	var unreleasedVolumes = []string{}
	var releaseDateInfo = []string{}
	for _, info := range volumeInfo {
		if info.ReleaseDate != nil && info.ReleaseDate.Before(today) {
			break
		} else {
			var releaseDate = defaultReleaseDate
			if info.ReleaseDate != nil {
				releaseDate = info.ReleaseDate.Format("January 2, 2006")
			}

			releaseDateInfo = append(releaseDateInfo, releaseDate)
			unreleasedVolumes = append(unreleasedVolumes, info.Name)
		}

	}

	return printReleaseInfoAndUpdateSeriesInfo(seriesInfo, unreleasedVolumes, releaseDateInfo, len(volumeInfo), volumeInfo[0].Name)
}

func sevenSeasEntertainmentGetSeriesVolumeInfo(seriesInfo config.SeriesInfo) config.SeriesInfo {
	volumeInfo := sevenseasentertainment.GetVolumeInfo(seriesInfo.Name, seriesInfo.SlugOverride, verbose)

	if len(volumeInfo) == 0 {
		logger.WriteInfo("The seven seas entertainment light novels do not exist for this series.")

		return seriesInfo
	}

	if len(volumeInfo) == seriesInfo.TotalVolumes && (len(seriesInfo.UnreleasedVolumes) == 0 || seriesInfo.UnreleasedVolumes[0].ReleaseDate != defaultReleaseDate) {
		return handleNoChangeDisplayAndSeriesInfoUpdates(seriesInfo)
	}

	var today = time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	var unreleasedVolumes = []string{}
	var releaseDateInfo = []string{}
	for _, info := range volumeInfo {
		if info.ReleaseDate != nil && info.ReleaseDate.Before(today) {
			break
		} else {
			var releaseDate = defaultReleaseDate
			if info.ReleaseDate != nil {
				releaseDate = info.ReleaseDate.Format("January 2, 2006")
			}

			releaseDateInfo = append(releaseDateInfo, releaseDate)
			unreleasedVolumes = append(unreleasedVolumes, info.Name)
		}

	}

	return printReleaseInfoAndUpdateSeriesInfo(seriesInfo, unreleasedVolumes, releaseDateInfo, len(volumeInfo), volumeInfo[0].Name)
}

func vizMediaGetSeriesVolumeInfo(seriesInfo config.SeriesInfo) config.SeriesInfo {
	volumeInfo := vizmedia.GetVolumeInfo(seriesInfo.Name, seriesInfo.SlugOverride, verbose)

	if len(volumeInfo) == 0 {
		logger.WriteInfo("The viz media series does not exist.")

		return seriesInfo
	}

	if len(volumeInfo) == seriesInfo.TotalVolumes && (len(seriesInfo.UnreleasedVolumes) == 0 || seriesInfo.UnreleasedVolumes[0].ReleaseDate != defaultReleaseDate) {
		return handleNoChangeDisplayAndSeriesInfoUpdates(seriesInfo)
	}

	var today = time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	var unreleasedVolumes = []string{}
	var releaseDateInfo = []string{}
	for _, info := range volumeInfo {
		if info.ReleaseDate.Before(today) {
			break
		} else {
			var releaseDate = info.ReleaseDate.Format("January 2, 2006")

			releaseDateInfo = append(releaseDateInfo, releaseDate)
			unreleasedVolumes = append(unreleasedVolumes, info.Name)
		}
	}

	return printReleaseInfoAndUpdateSeriesInfo(seriesInfo, unreleasedVolumes, releaseDateInfo, len(volumeInfo), volumeInfo[0].Name)
}

func printReleaseInfoAndUpdateSeriesInfo(seriesInfo config.SeriesInfo, unreleasedVolumes, releaseDateInfo []string, totalVolumes int, latestVolumeName string) config.SeriesInfo {
	var releaseInfo = []config.ReleaseInfo{}
	for i, unreleasedVol := range unreleasedVolumes {
		releaseInfo = append(releaseInfo, config.ReleaseInfo{
			Name:        unreleasedVol,
			ReleaseDate: releaseDateInfo[i],
		})

		logger.WriteInfo(getUnreleasedVolumeDisplayText(unreleasedVol, releaseDateInfo[i]))
	}

	seriesInfo.TotalVolumes = totalVolumes
	seriesInfo.LatestVolume = latestVolumeName
	seriesInfo.UnreleasedVolumes = releaseInfo

	return seriesInfo
}

func handleNoChangeDisplayAndSeriesInfoUpdates(seriesInfo config.SeriesInfo) config.SeriesInfo {
	logger.WriteWarn("No change in list of volumes from last check.")

	var updatedUnreleasedVolumes = []config.ReleaseInfo{}
	var today = time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	for _, unreleasedVol := range seriesInfo.UnreleasedVolumes {
		if !unreleasedDateIsBeforeDate(unreleasedVol.ReleaseDate, today) {
			logger.WriteInfo(getUnreleasedVolumeDisplayText(unreleasedVol.Name, unreleasedVol.ReleaseDate))
			updatedUnreleasedVolumes = append(updatedUnreleasedVolumes, unreleasedVol)
		}
	}

	seriesInfo.UnreleasedVolumes = updatedUnreleasedVolumes

	return seriesInfo
}
