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
			// "Make sure Mushoku Tensei volumes are correctly extracted": {
			// 	SeriesName: "Mushoku Tensei: Jobless Reincarnation (Light Novel)",
			// 	ExpectedVolumes: []*sitehandler.VolumeInfo{
			// 		{
			// 			Name:        "Arifureta Zero: Volume 6",
			// 			ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.August, 4, 0, 0, 0, 0, time.Local)),
			// 		},
			// 		{
			// 			Name:        "Arifureta Zero: Volume 5",
			// 			ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.December, 1, 0, 0, 0, 0, time.Local)),
			// 		},
			// 		{
			// 			Name:        "Arifureta Zero: Volume 4",
			// 			ReleaseDate: sitehandler.TimePtr(time.Date(2020, time.July, 6, 0, 0, 0, 0, time.Local)),
			// 		},
			// 		{
			// 			Name:        "Arifureta Zero: Volume 3",
			// 			ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.November, 16, 0, 0, 0, 0, time.Local)),
			// 		},
			// 		{
			// 			Name:        "Arifureta Zero: Volume 2",
			// 			ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.February, 4, 0, 0, 0, 0, time.Local)),
			// 		},
			// 		{
			// 			Name:        "Arifureta Zero: Volume 1",
			// 			ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.April, 11, 0, 0, 0, 0, time.Local)),
			// 		},
			// 	},
			// 	ExpectedCount: 6,
			// },
			"Make sure Berserk of Gluttony volumes are correctly extracted": {
				SeriesName: "Berserk of Gluttony (Light Novel)",
				ExpectedVolumes: []*sitehandler.VolumeInfo{
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 8",
						ReleaseDate: sitehandler.TimePtr(time.Date(2024, time.January, 11, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 7",
						ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.August, 4, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 6",
						ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.April, 7, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 5",
						ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.January, 27, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 4",
						ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.September, 30, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 3",
						ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.June, 10, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 2",
						ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.February, 11, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Berserk of Gluttony (Light Novel) Vol. 1",
						ReleaseDate: sitehandler.TimePtr(time.Date(2020, time.October, 29, 0, 0, 0, 0, time.Local)),
					},
				},
				ExpectedCount: 8,
			},
		},
		Endpoints: []sitehandler.MockedEndpoint{
			{
				Slug:     "berserk-of-gluttony-light-novel/",
				Response: berserkOfGluttonyResponse,
			},
			{
				Slug:     "mushoku-tensei-jobless-reincarnation-light-novel/",
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
