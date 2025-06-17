//go:build unit

package jnovelclub_test

import (
	_ "embed"
	"testing"
	"time"

	jnovelclub "github.com/pjkaufman/go-go-gadgets/magnum/internal/jnovel-club"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
)

var (
	//go:embed test/arifureta-zero.golden
	arifuretaZeroResponse string
	//go:embed test/how-a-realist-hero-rebuilt-the-kingdom.golden
	howARealisHeroRebuiltTheKingdom string

	getVolumeInfoTestSetup = sitehandler.GetVolumeInfoTestCases{
		Tests: map[string]sitehandler.GetVolumeInfoTestCase{
			"Make sure Arifureta Zero volumes are correctly extracted": {
				SeriesName: "Arifureta Zero",
				ExpectedVolumes: []*sitehandler.VolumeInfo{
					{
						Name:        "Arifureta Zero: Volume 6",
						ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.August, 4, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Arifureta Zero: Volume 5",
						ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.December, 1, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Arifureta Zero: Volume 4",
						ReleaseDate: sitehandler.TimePtr(time.Date(2020, time.July, 6, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Arifureta Zero: Volume 3",
						ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.November, 16, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Arifureta Zero: Volume 2",
						ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.February, 4, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Arifureta Zero: Volume 1",
						ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.April, 11, 0, 0, 0, 0, time.Local)),
					},
				},
				ExpectedCount: 6,
			},
			"Make sure How a Realist Hero Rebuilt the Kingdom volumes are correctly extracted": {
				SeriesName: "How a Realist Hero Rebuilt the Kingdom",
				ExpectedVolumes: []*sitehandler.VolumeInfo{
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 19",
						ReleaseDate: sitehandler.TimePtr(time.Date(2025, time.March, 24, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 18",
						ReleaseDate: sitehandler.TimePtr(time.Date(2023, time.November, 29, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 17",
						ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.November, 7, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 16",
						ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.June, 7, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 15",
						ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.January, 17, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 14",
						ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.October, 25, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 13",
						ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.March, 6, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 12",
						ReleaseDate: sitehandler.TimePtr(time.Date(2020, time.August, 22, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 11",
						ReleaseDate: sitehandler.TimePtr(time.Date(2020, time.April, 26, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 10",
						ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.October, 20, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 9",
						ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.July, 20, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 8",
						ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.February, 15, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 7",
						ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.September, 22, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 6",
						ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.May, 31, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 5",
						ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.February, 1, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 4",
						ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.October, 18, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 3",
						ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.August, 2, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 2",
						ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.May, 13, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 1",
						ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.February, 23, 0, 0, 0, 0, time.Local)),
					},
				},
				ExpectedCount: 19,
			},
		},
		Endpoints: []sitehandler.MockedEndpoint{
			{
				Slug:     "arifureta-zero",
				Response: arifuretaZeroResponse,
			},
			{
				Slug:     "how-a-realist-hero-rebuilt-the-kingdom",
				Response: howARealisHeroRebuiltTheKingdom,
			},
		},
		CreateSiteHandler: func(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
			return jnovelclub.NewJNovelClubHandler(options)
		},
	}
)

func TestGetVolumeInfo(t *testing.T) {
	sitehandler.RunTests(t, getVolumeInfoTestSetup)
}
