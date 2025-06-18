//go:build unit

package vizmedia_test

import (
	_ "embed"
	"testing"
	"time"

	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/vizmedia"
)

var (
	alreadyReleasedDate = time.Now().Add(-1 * 24 * time.Hour)
	//go:embed test/nausicaa-of-the-valley-of-the-wind.golden
	nausicaaOfTheValleyOfTheWindMainPage string
	//go:embed test/manga-books-nausicaa-of-the-valley-of-the-wind-section-115444-more.golden
	nausicaaOfTheValleyOfTheWindVolumesPage string

	getVolumeInfoTestSetup = sitehandler.GetVolumeInfoTestCases{
		Tests: map[string]sitehandler.GetVolumeInfoTestCase{
			"Make sure Nausicaä of the Valley of the Wind volumes are correctly extracted": {
				SeriesName:   "Nausicaä of the Valley of the Wind",
				SlugOverride: sitehandler.StringPtr("nausicaa-of-the-valley-of-the-wind"),
				ExpectedVolumes: []*sitehandler.VolumeInfo{
					{
						Name:        "Nausicaä of the Valley of the Wind Picture Book",
						ReleaseDate: sitehandler.TimePtr(alreadyReleasedDate),
					},
					{
						Name:        "The Art of Nausicaä of the Valley of the Wind",
						ReleaseDate: sitehandler.TimePtr(alreadyReleasedDate),
					},
					{
						Name:        "Nausicaä of the Valley of the Wind Box Set",
						ReleaseDate: sitehandler.TimePtr(alreadyReleasedDate),
					},
					{
						Name:        "Nausicaä of the Valley of the Wind: Watercolor Impressions",
						ReleaseDate: sitehandler.TimePtr(alreadyReleasedDate),
					},
					{
						Name:        "Nausicaä of the Valley of the Wind, Vol. 7",
						ReleaseDate: sitehandler.TimePtr(alreadyReleasedDate),
					},
					{
						Name:        "Nausicaä of the Valley of the Wind, Vol. 6",
						ReleaseDate: sitehandler.TimePtr(alreadyReleasedDate),
					},
					{
						Name:        "Nausicaä of the Valley of the Wind, Vol. 5",
						ReleaseDate: sitehandler.TimePtr(alreadyReleasedDate),
					},
					{
						Name:        "Nausicaä of the Valley of the Wind, Vol. 4",
						ReleaseDate: sitehandler.TimePtr(alreadyReleasedDate),
					},
					{
						Name:        "Nausicaä of the Valley of the Wind, Vol. 3",
						ReleaseDate: sitehandler.TimePtr(alreadyReleasedDate),
					},
					{
						Name:        "Nausicaä of the Valley of the Wind, Vol. 2",
						ReleaseDate: sitehandler.TimePtr(alreadyReleasedDate),
					},
					{
						Name:        "Nausicaä of the Valley of the Wind, Vol. 1",
						ReleaseDate: sitehandler.TimePtr(alreadyReleasedDate),
					},
				},
				ExpectedCount: 11,
			},
		},
		Endpoints: []sitehandler.MockedEndpoint{
			{
				Slug:     "nausicaa-of-the-valley-of-the-wind",
				Response: nausicaaOfTheValleyOfTheWindMainPage,
			},
			{
				Slug:     "manga-books/nausicaa-of-the-valley-of-the-wind/section/115444/more",
				Response: nausicaaOfTheValleyOfTheWindVolumesPage,
			},
		},
		CreateSiteHandler: func(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
			return vizmedia.NewVizMediaHandler(options)
		},
	}
)

func TestGetVolumeInfo(t *testing.T) {
	sitehandler.RunTests(t, getVolumeInfoTestSetup)
}
