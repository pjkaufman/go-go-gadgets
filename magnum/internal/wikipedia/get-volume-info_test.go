//go:build unit

package wikipedia_test

import (
	_ "embed"
	"net/http"
	"testing"
	"time"

	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/wikipedia"
)

const testApiPath = "w/json"

var (
	//go:embed test/wiki-list-of-the-rising-of-the-shield-hero-volumes.golden
	risingOfTheShieldHeroResponse string
	//go:embed test/api-list-of-the-rising-of-the-shield-hero-volumes.golden
	risingOfTheShieldHeroApiResponse string
	//go:embed test/wiki-rokka-braves-of-the-six-flowers.golden
	rokkaBravesOfTheSixFlowersResponse string
	//go:embed test/api-rokka-braves-of-the-six-flowers.golden
	rokkaBravesOfTheSixFlowersApiResponse string

	getVolumeInfoTestSetup = sitehandler.GetVolumeInfoTestCases{
		Tests: map[string]sitehandler.GetVolumeInfoTestCase{
			"Make sure Rokka: Braves of the Six Flowers volumes are correctly extracted": {
				SeriesName: "Rokka: Braves of the Six Flowers",
				// SlugOverride: sitehandler.StringPtr("")
				ExpectedVolumes: []*sitehandler.VolumeInfo{
					{
						Name:        "Rokka: Braves of the Six Flowers Vol. Archive 1",
						ReleaseDate: nil,
					},
					{
						Name:        "Rokka: Braves of the Six Flowers Vol. 6",
						ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.December, 11, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Rokka: Braves of the Six Flowers Vol. 5",
						ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.August, 21, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Rokka: Braves of the Six Flowers Vol. 4",
						ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.April, 24, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Rokka: Braves of the Six Flowers Vol. 3",
						ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.December, 12, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Rokka: Braves of the Six Flowers Vol. 2",
						ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.August, 22, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Rokka: Braves of the Six Flowers Vol. 1",
						ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.April, 18, 0, 0, 0, 0, time.Local)),
					},
				},
				ExpectedCount: 7,
			},
			"Make sure The Rising of the Shield Hero volumes are correctly extracted": {
				SeriesName:            "The Rising of the Shield Hero",
				SlugOverride:          sitehandler.StringPtr("List_of_The_Rising_of_the_Shield_Hero_volumes"),
				TablesToParseOverride: sitehandler.IntPtr(1),
				ExpectedVolumes: []*sitehandler.VolumeInfo{
					{
						Name:        "The Rising of the Shield Hero Vol. 22",
						ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.December, 21, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 21",
						ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.October, 26, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 20",
						ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.June, 22, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 19",
						ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.April, 27, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 18",
						ReleaseDate: sitehandler.TimePtr(time.Date(2020, time.November, 12, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 17",
						ReleaseDate: sitehandler.TimePtr(time.Date(2020, time.July, 14, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 16",
						ReleaseDate: sitehandler.TimePtr(time.Date(2020, time.March, 15, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 15",
						ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.December, 15, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 14",
						ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.October, 15, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 13",
						ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.December, 18, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 12",
						ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.August, 18, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 11",
						ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.June, 12, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 10",
						ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.March, 20, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 9",
						ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.November, 15, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 8",
						ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.June, 13, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 7",
						ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.April, 18, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 6",
						ReleaseDate: sitehandler.TimePtr(time.Date(2016, time.November, 22, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 5",
						ReleaseDate: sitehandler.TimePtr(time.Date(2016, time.August, 23, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 4",
						ReleaseDate: sitehandler.TimePtr(time.Date(2016, time.June, 14, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 3",
						ReleaseDate: sitehandler.TimePtr(time.Date(2016, time.February, 16, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 2",
						ReleaseDate: sitehandler.TimePtr(time.Date(2015, time.October, 20, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "The Rising of the Shield Hero Vol. 1",
						ReleaseDate: sitehandler.TimePtr(time.Date(2015, time.September, 15, 0, 0, 0, 0, time.Local)),
					},
				},
				ExpectedCount: 22,
			},
		},
		Endpoints: []sitehandler.MockedEndpoint{
			{
				Slug:     "wiki/List_of_The_Rising_of_the_Shield_Hero_volumes",
				Response: risingOfTheShieldHeroResponse,
			},
			{
				Slug:     "wiki/Rokka:_Braves_of_the_Six_Flowers",
				Response: rokkaBravesOfTheSixFlowersResponse,
			},
			{
				Slug: testApiPath,
				CustomHandler: func(w http.ResponseWriter, r *http.Request) {
					pageTitle := r.URL.Query().Get("page")

					switch pageTitle {
					case "List_of_The_Rising_of_the_Shield_Hero_volumes":
						w.Write([]byte(risingOfTheShieldHeroApiResponse))
					case "Rokka:_Braves_of_the_Six_Flowers":
						w.Write([]byte(rokkaBravesOfTheSixFlowersApiResponse))
					default:
						http.Error(w, "Not found", 404)
					}
				},
			},
		},
		CreateSiteHandler: func(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
			options.ApiPath = testApiPath

			return wikipedia.NewWikipediaHandler(options)
		},
	}
)

func TestGetVolumeInfo(t *testing.T) {
	sitehandler.RunTests(t, getVolumeInfoTestSetup)
}
