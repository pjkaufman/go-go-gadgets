//go:build unit

package sevenseasentertainment_test

// import (
// 	_ "embed"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	jnovelclub "github.com/pjkaufman/go-go-gadgets/magnum/internal/jnovel-club"
// 	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
// 	"github.com/stretchr/testify/assert"
// )

// var (
// 	//go:embed test/arifureta-zero.golden
// 	arifuretaZeroResponse string
// 	//go:embed test/how-a-realist-hero-rebuilt-the-kingdom.golden
// 	howARealisHeroRebuiltTheKingdom string
// 	getVolumeInfoTestCases          = map[string]getVolumeInfoTestCase{
// 		"Make sure Arifureta Zero volumes are correctly extracted": {
// 			SeriesName: "Arifureta Zero",
// 			ExpectedVolumes: []*sitehandler.VolumeInfo{
// 				{
// 					Name:        "Arifureta Zero: Volume 6",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.August, 4, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "Arifureta Zero: Volume 5",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.December, 1, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "Arifureta Zero: Volume 4",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2020, time.July, 6, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "Arifureta Zero: Volume 3",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.November, 16, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "Arifureta Zero: Volume 2",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.February, 4, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "Arifureta Zero: Volume 1",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.April, 11, 0, 0, 0, 0, time.Local)),
// 				},
// 			},
// 			ExpectedCount: 6,
// 		},
// 		"Make sure How a Realist Hero Rebuilt the Kingdom volumes are correctly extracted": {
// 			SeriesName: "How a Realist Hero Rebuilt the Kingdom",
// 			ExpectedVolumes: []*sitehandler.VolumeInfo{
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 19",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2025, time.March, 24, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 18",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2023, time.November, 29, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 17",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.November, 7, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 16",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.June, 7, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 15",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2022, time.January, 17, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 14",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.October, 25, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 13",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2021, time.March, 6, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 12",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2020, time.August, 22, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 11",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2020, time.April, 26, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 10",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.October, 20, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 9",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.July, 20, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 8",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2019, time.February, 15, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 7",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.September, 22, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 6",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.May, 31, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 5",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2018, time.February, 1, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 4",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.October, 18, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 3",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.August, 2, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 2",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.May, 13, 0, 0, 0, 0, time.Local)),
// 				},
// 				{
// 					Name:        "How a Realist Hero Rebuilt the Kingdom: Volume 1",
// 					ReleaseDate: sitehandler.TimePtr(time.Date(2017, time.February, 23, 0, 0, 0, 0, time.Local)),
// 				},
// 			},
// 			ExpectedCount: 19,
// 		},
// 	}
// 	robotsFile = `
// 	User-agent: *
// 	Allow: /allowed
// 	Disallow: /disallowed
// 	Disallow: /allowed*q=
// 	`
// 	releaseDateFormat = "January 2, 2006"
// )

// type getVolumeInfoTestCase struct {
// 	SeriesName      string
// 	SlugOverride    *string
// 	ExpectedVolumes []*sitehandler.VolumeInfo
// 	ExpectedCount   int
// }

// func createJnovelServerInstance() *httptest.Server {
// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/arifureta-zero", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "text/html")
// 		fmt.Fprint(w, arifuretaZeroResponse)
// 	})

// 	mux.HandleFunc("/how-a-realist-hero-rebuilt-the-kingdom", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "text/html")
// 		fmt.Fprint(w, howARealisHeroRebuiltTheKingdom)
// 	})

// 	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(200)
// 		w.Write([]byte(robotsFile))
// 	})

// 	return httptest.NewUnstartedServer(mux)
// }

// func TestGetVolumeInfo(t *testing.T) {
// 	srv := createJnovelServerInstance()
// 	srv.Start()

// 	options := sitehandler.SiteHandlerOptions{
// 		BaseURL: srv.URL + "/",
// 		Verbose: false,
// 	}

// 	jnc := jnovelclub.NewJNovelClubHandler(options)

// 	for name, args := range getVolumeInfoTestCases {
// 		t.Run(name, func(t *testing.T) {
// 			scrapingOptions := sitehandler.ScrapingOptions{}
// 			if args.SlugOverride != nil {
// 				scrapingOptions.SlugOverride = args.SlugOverride
// 			}

// 			actualVolumes, actualCount, err := jnc.GetVolumeInfo(args.SeriesName, scrapingOptions)

// 			assert.Nil(t, err)
// 			assert.Equal(t, args.ExpectedCount, actualCount)
// 			assert.Equal(t, len(args.ExpectedVolumes), len(actualVolumes))

// 			if args.ExpectedVolumes != nil {
// 				for i, expectedVolume := range args.ExpectedVolumes {
// 					assert.Equal(t, expectedVolume.Name, actualVolumes[i].Name)
// 					if expectedVolume.ReleaseDate != nil && actualVolumes[i].ReleaseDate != nil {
// 						assert.Equal(t, expectedVolume.ReleaseDate.Format(releaseDateFormat), actualVolumes[i].ReleaseDate.Format(releaseDateFormat))
// 					} else {
// 						assert.Equal(t, expectedVolume.ReleaseDate, actualVolumes[i].ReleaseDate)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// func sitehandler.TimePtr(t time.Time) *time.Time {
// 	return &t
// }
