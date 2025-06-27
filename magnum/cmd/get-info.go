package cmd

import (
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	jnovelclub "github.com/pjkaufman/go-go-gadgets/magnum/internal/jnovel-club"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/sevenseasentertainment"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/vizmedia"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/wikipedia"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/yenpress"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var promptForSeries bool

var (
	handlersInitialized           = false
	jNovelClubHandler             sitehandler.SiteHandler
	sevenSeasEntertainmentHandler sitehandler.SiteHandler
	wikipediaHandler              sitehandler.SiteHandler
	yenPressHandler               sitehandler.SiteHandler
	vizMediaHandler               sitehandler.SiteHandler
)

// GetInfoCmd represents the get book info command
var GetInfoCmd = &cobra.Command{
	Use:   "get-info",
	Short: "Gets the book release info for books that have been added to the list of series to track",
	Example: heredoc.Doc(`To get all of the release data for non-completed series:
	magnum get-info

	To get release data including completed series:
	magnum get-info -c

	To get release data for a specific series:
	magnum get-info -s "Series Name"

	To interactively select a series from a prompt:
	magnum get-info -p
	`),
	Run: func(cmd *cobra.Command, args []string) {
		seriesInfo := config.GetConfig()

		if promptForSeries {
			seriesName = selectBookName(seriesInfo.Series, includeCompleted)
		}

		if seriesName != "" {
			if !seriesInfo.HasSeries(seriesName) {
				logger.WriteWarnf("No series with the name %q is in the series list.", seriesName)
				return
			}

			setupHandlers()
			for i, series := range seriesInfo.Series {
				if strings.EqualFold(seriesName, series.Name) {
					seriesInfo.Series[i] = getSeriesVolumeInfo(series)
					return
				}
			}

			return // on the off chance that we somehow have it, but then don't find it
		}

		setupHandlers()
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
	GetInfoCmd.Flags().StringVarP(&seriesName, "series", "s", "", "get info for just the specified series")
	GetInfoCmd.Flags().BoolVarP(&promptForSeries, "prompt-name", "p", false, "get info for a series that you will select from a prompt")
}

func getSeriesVolumeInfo(seriesInfo config.SeriesInfo) config.SeriesInfo {
	logger.WriteInfof("Checking for volume info for %q\n", seriesInfo.Name)

	var handler sitehandler.SiteHandler
	switch seriesInfo.Publisher {
	case config.YenPress:
		handler = yenPressHandler

	case config.JNovelClub:
		handler = jNovelClubHandler
	case config.SevenSeasEntertainment:
		handler = sevenSeasEntertainmentHandler
	case config.OnePeaceBooks, config.HanashiMedia:
		handler = wikipediaHandler
	case config.VizMedia:
		handler = vizMediaHandler
	}

	if handler == nil {
		return seriesInfo
	}

	return sitehandlerGetSeriesVolumeInfo(seriesInfo, handler)
}

func setupHandlers() {
	if handlersInitialized {
		return
	}

	jNovelClubHandler = jnovelclub.NewJNovelClubHandler(sitehandler.SiteHandlerOptions{
		BaseURL:        jnovelclub.BaseURL,
		Verbose:        verbose,
		UserAgent:      userAgent,
		AllowedDomains: jnovelclub.AllowedDomains,
	})

	yenPressHandler = yenpress.NewYenPressHandler(sitehandler.SiteHandlerOptions{
		BaseURL:        yenpress.BaseURL,
		Verbose:        verbose,
		UserAgent:      userAgent,
		AllowedDomains: yenpress.AllowedDomains,
	})

	sevenSeasEntertainmentHandler = sevenseasentertainment.NewSevenSeasEntertainmentHandler(sitehandler.SiteHandlerOptions{
		BaseURL:        sevenseasentertainment.BaseURL,
		Verbose:        verbose,
		UserAgent:      userAgent,
		AllowedDomains: sevenseasentertainment.AllowedDomains,
	})

	vizMediaHandler = vizmedia.NewVizMediaHandler(sitehandler.SiteHandlerOptions{
		BaseURL:        vizmedia.BaseURL,
		Verbose:        verbose,
		UserAgent:      userAgent,
		AllowedDomains: vizmedia.AllowedDomains,
	})

	wikipediaHandler = wikipedia.NewWikipediaHandler(sitehandler.SiteHandlerOptions{
		BaseURL:        wikipedia.BaseURL,
		Verbose:        verbose,
		UserAgent:      userAgent,
		BuildApiPath:   wikipedia.GetWikipediaAPIUrl,
		AllowedDomains: wikipedia.AllowedDomains,
	})

	handlersInitialized = true
}

func sitehandlerGetSeriesVolumeInfo(seriesInfo config.SeriesInfo, handler sitehandler.SiteHandler) config.SeriesInfo {
	volumes, numVolumes, err := handler.GetVolumeInfo(seriesInfo.Name, sitehandler.ScrapingOptions{
		SlugOverride:          seriesInfo.SlugOverride,
		TablesToParseOverride: seriesInfo.WikipediaTablesToParseOverride,
	})
	if err != nil {
		logger.WriteError(err.Error())
	}

	if len(volumes) == -1 {
		logger.WriteErrorf("The %s light novels were not found for %q. The HTML for the site or page may have changed.\n", config.PublisherToDisplayString(seriesInfo.Publisher), seriesInfo.Name)
	}

	if numVolumes == 0 {
		logger.WriteInfof("The %s light novels do not exist for series %q.\n", config.PublisherToDisplayString(seriesInfo.Publisher), seriesInfo.Name)

		return seriesInfo
	}

	var shouldSkipGettingVolumesAndHandleExistingData bool
	switch seriesInfo.Publisher {
	case config.YenPress:
		// We cannot really trust that Yen Press release data is 100% accurate as they could have delayed the book release,
		// so we need to double check volumes any time we have an upcoming release
		shouldSkipGettingVolumesAndHandleExistingData = numVolumes == seriesInfo.TotalVolumes && len(seriesInfo.UnreleasedVolumes) == 0
	case config.JNovelClub:
		shouldSkipGettingVolumesAndHandleExistingData = len(volumes) == seriesInfo.TotalVolumes
	default:
		shouldSkipGettingVolumesAndHandleExistingData = len(volumes) == seriesInfo.TotalVolumes && (len(seriesInfo.UnreleasedVolumes) == 0 || seriesInfo.UnreleasedVolumes[0].ReleaseDate != defaultReleaseDate)
	}

	if shouldSkipGettingVolumesAndHandleExistingData {
		return handleNoChangeDisplayAndSeriesInfoUpdates(seriesInfo)
	}

	var today = time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	var unreleasedVolumes = []string{}
	var releaseDateInfo = []string{}
	for _, info := range volumes {
		if info.ReleaseDate != nil && info.ReleaseDate.Before(today) {
			break
		} else {
			var releaseDate = defaultReleaseDate
			if info.ReleaseDate != nil {
				releaseDate = info.ReleaseDate.Format(releaseDateFormat)
			}

			releaseDateInfo = append(releaseDateInfo, releaseDate)
			unreleasedVolumes = append(unreleasedVolumes, info.Name)
		}
	}

	return printReleaseInfoAndUpdateSeriesInfo(seriesInfo, unreleasedVolumes, releaseDateInfo, numVolumes, volumes[0].Name)
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
