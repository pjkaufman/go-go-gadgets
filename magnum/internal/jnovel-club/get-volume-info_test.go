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
						ReleaseDate: sitehandler.TimePtr(time.Unix(1659589200, 0)),
					},
					{
						Name:        "Arifureta Zero: Volume 5",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1638338400, 0)),
					},
					{
						Name:        "Arifureta Zero: Volume 4",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1594011600, 0)),
					},
					{
						Name:        "Arifureta Zero: Volume 3",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1573884000, 0)),
					},
					{
						Name:        "Arifureta Zero: Volume 2",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1549256400, 0)),
					},
					{
						Name:        "Arifureta Zero: Volume 1",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1523433600, 0)),
					},
				},
				ExpectedCount: 6,
			},
			"Make sure How a Realist Hero Rebuilt the Kingdom volumes are correctly extracted": {
				SeriesName: "How a Realist Hero Rebuilt the Kingdom",
				ExpectedVolumes: []*sitehandler.VolumeInfo{
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Short Story Chronicles",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1780549200, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 20",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1766988000, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 19",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1742792400, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 18",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1701237600, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 17",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1667800800, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 16",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1654578000, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 15",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1642399200, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 14",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1635138000, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 13",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1615010400, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 12",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1598072400, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 11",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1587877200, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 10",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1571547600, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 9",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1563598800, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 8",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1550210400, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 7",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1537610400, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 6",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1527811200, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 5",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1517464800, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 4",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1508342400, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 3",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1501686000, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 2",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1494687600, 0)),
					},
					{
						Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 1",
						ReleaseDate: sitehandler.TimePtr(time.Unix(1487862000, 0)),
					},
				},
				ExpectedCount: 21,
			},
		},
		Endpoints: []sitehandler.MockedEndpoint{
			{
				Slug:     "series/arifureta-zero",
				Response: arifuretaZeroResponse,
			},
			{
				Slug:     "series/how-a-realist-hero-rebuilt-the-kingdom",
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
