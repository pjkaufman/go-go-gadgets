//go:build unit

package sevenseasentertainment_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/sevenseasentertainment"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
)

var (
	//go:embed test/mushoku-tensei-jobless-reincarnation-light-novel.golden
	mushokuTenseiResponse string
	//go:embed test/berserk-of-gluttony-light-novel.golden
	berserkOfGluttonyResponse string

	getVolumeInfoTestSetup = sitehandler.GetVolumeInfoTestCases{
		Tests: map[string]sitehandler.GetVolumeInfoTestCase{
			"Make sure Mushoku Tensei volumes are correctly extracted without the omnibus version": {
				SeriesName: "Mushoku Tensei: Jobless Reincarnation (Light Novel)",
				ExpectedVolumes: []*sitehandler.VolumeInfo{
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation – A Journey of Two Lifetimes",
						ReleaseDate: new(time.Date(2024, time.September, 17, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation – Recollections (Light Novel) [Ebook]",
						ReleaseDate: new(time.Date(2024, time.May, 30, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 26",
						ReleaseDate: new(time.Date(2024, time.February, 1, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 25",
						ReleaseDate: new(time.Date(2023, time.October, 26, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 24",
						ReleaseDate: new(time.Date(2023, time.August, 24, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 23",
						ReleaseDate: new(time.Date(2023, time.July, 6, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 22",
						ReleaseDate: new(time.Date(2023, time.May, 4, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 21",
						ReleaseDate: new(time.Date(2023, time.February, 23, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 20",
						ReleaseDate: new(time.Date(2022, time.December, 22, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 19",
						ReleaseDate: new(time.Date(2022, time.September, 29, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 18",
						ReleaseDate: new(time.Date(2022, time.July, 14, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 17",
						ReleaseDate: new(time.Date(2022, time.June, 16, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 16",
						ReleaseDate: new(time.Date(2022, time.April, 21, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 15",
						ReleaseDate: new(time.Date(2022, time.February, 17, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 14",
						ReleaseDate: new(time.Date(2021, time.November, 25, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 13",
						ReleaseDate: new(time.Date(2021, time.September, 23, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 12",
						ReleaseDate: new(time.Date(2021, time.July, 15, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 11",
						ReleaseDate: new(time.Date(2021, time.May, 6, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 10",
						ReleaseDate: new(time.Date(2021, time.March, 25, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 9",
						ReleaseDate: new(time.Date(2021, time.January, 21, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 8",
						ReleaseDate: new(time.Date(2020, time.October, 1, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 7",
						ReleaseDate: new(time.Date(2020, time.July, 9, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 6",
						ReleaseDate: new(time.Date(2020, time.April, 2, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 5",
						ReleaseDate: new(time.Date(2020, time.January, 16, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 4",
						ReleaseDate: new(time.Date(2019, time.October, 10, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 3",
						ReleaseDate: new(time.Date(2019, time.August, 1, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 2",
						ReleaseDate: new(time.Date(2019, time.May, 23, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 1",
						ReleaseDate: new(time.Date(2019, time.April, 4, 0, 0, 0, 0, time.Local)),
					},
				},
				ExpectedCount: 28,
			},
			"Make sure Berserk of Gluttony volumes are correctly extracted": {
				SeriesName: "Berserk of Gluttony (Light Novel)",
				ExpectedVolumes: []*sitehandler.VolumeInfo{
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 8",
						ReleaseDate: new(time.Date(2024, time.January, 11, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 7",
						ReleaseDate: new(time.Date(2022, time.August, 4, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 6",
						ReleaseDate: new(time.Date(2022, time.April, 7, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 5",
						ReleaseDate: new(time.Date(2022, time.January, 27, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 4",
						ReleaseDate: new(time.Date(2021, time.September, 30, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 3",
						ReleaseDate: new(time.Date(2021, time.June, 10, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 2",
						ReleaseDate: new(time.Date(2021, time.February, 11, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 1",
						ReleaseDate: new(time.Date(2020, time.October, 29, 0, 0, 0, 0, time.Local)),
					},
				},
				ExpectedCount: 8,
			},
		},
		Endpoints: []sitehandler.MockedEndpoint{
			{
				Slug:     "series/berserk-of-gluttony-light-novel/",
				Response: berserkOfGluttonyResponse,
			},
			{
				Slug:     "series/mushoku-tensei-jobless-reincarnation-light-novel/",
				Response: mushokuTenseiResponse,
			},
		},
		CreateSiteHandler: func(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
			return sevenseasentertainment.NewSevenSeasEntertainmentHandler(options)
		},
	}
)

func TestGetVolumeInfo(t *testing.T) {
	sitehandler.RunTests(t, getVolumeInfoTestSetup)
}
