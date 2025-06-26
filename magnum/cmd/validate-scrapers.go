package cmd

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var validationSeries = []config.SeriesInfo{
	{ // Yen Press
		Name:         "So I'm a Spider, So What? (light novel)",
		TotalVolumes: 16,
		LatestVolume: "So I'm a Spider, So What?, Vol. 16 (light novel)",
		Publisher:    config.YenPress,
		Status:       config.Completed,
		Type:         config.LightNovel,
	},
	{ // JNovel-Club
		Name:         "Arifureta Zero",
		TotalVolumes: 6,
		LatestVolume: "Volume 6",
		Publisher:    config.JNovelClub,
		Status:       config.Completed,
		Type:         config.LightNovel,
	},
	{ // Seven Seas Entertainment
		Name:         "Berserk of Gluttony",
		TotalVolumes: 8,
		LatestVolume: "Berserk of Gluttony (Light Novel) Vol. 8",
		Publisher:    config.SevenSeasEntertainment,
		Status:       config.Completed,
		Type:         config.LightNovel,
	},
	{ // Viz Media
		Name:         "Nausicaa of the Valley of the Wind",
		TotalVolumes: 11,
		LatestVolume: "Nausicaä of the Valley of the Wind Picture Book",
		Publisher:    config.VizMedia,
		Status:       config.Completed,
		Type:         config.Manga,
	},
}

// ValidateScraperCmd represents the validate scraper command
var ValidateScraperCmd = &cobra.Command{
	Use:   "validate",
	Short: "Runs the web scraper logic for a single series on each scraper with an already known result to determine if the scraper is still functioning or if it needs an update.",
	Example: heredoc.Doc(`To test all of the scrapers used:
	magnum validate`),
	Run: func(cmd *cobra.Command, args []string) {
		logger.WriteInfo("Validating scrapers...")

		var (
			validationResult strings.Builder
			scraperName      string
		)
		validationResult.WriteString("Validation Results:\n")
		setupHandlers()
		for _, series := range validationSeries {
			output := getSeriesVolumeInfo(series)

			scraperName = config.PublisherToDisplayString(series.Publisher)
			if output.TotalVolumes == series.TotalVolumes && output.LatestVolume == series.LatestVolume {
				validationResult.WriteString(fmt.Sprintf("- %s: working as expected\n", scraperName))
			} else {
				validationResult.WriteString(fmt.Sprintf("- %s: did not parse information correctly\n", scraperName))
			}
		}

		validationResult.WriteString(fmt.Sprintf("- Wikipedia: not tested, but used for %s and %s", config.PublisherToDisplayString(config.HanashiMedia), config.PublisherToDisplayString(config.OnePeaceBooks)))

		logger.WriteInfo(validationResult.String())
	},
}

func init() {
	rootCmd.AddCommand(ValidateScraperCmd)

	ValidateScraperCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show more info about what is going on")
}
