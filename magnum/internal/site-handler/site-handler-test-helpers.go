//go:build unit

package sitehandler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	robotsFile = `
	User-agent: *
	Allow: /allowed
	Disallow: /disallowed
	Disallow: /allowed*q=
	`
	releaseDateFormat = "January 2, 2006"
)

type GetVolumeInfoTestCase struct {
	SeriesName      string
	SlugOverride    *string
	ExpectedVolumes []*VolumeInfo
	ExpectedCount   int
}

type MockedEndpoint struct {
	Slug     string
	Response string
}

type GetVolumeInfoTestCases struct {
	Tests             map[string]GetVolumeInfoTestCase
	Endpoints         []MockedEndpoint
	CreateSiteHandler func(SiteHandlerOptions) SiteHandler
}

func createMockServerInstance(endpoints []MockedEndpoint) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(robotsFile))
	})

	for _, endpoint := range endpoints {
		mux.HandleFunc(fmt.Sprintf("/%s", endpoint.Slug), func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, endpoint.Response)
		})
	}

	return httptest.NewUnstartedServer(mux)
}

func RunTests(t *testing.T, cases GetVolumeInfoTestCases) {
	srv := createMockServerInstance(cases.Endpoints)
	srv.Start()

	options := SiteHandlerOptions{
		BaseURL: srv.URL + "/",
		Verbose: false,
	}

	handler := cases.CreateSiteHandler(options)

	for name, args := range cases.Tests {
		t.Run(name, func(t *testing.T) {
			scrapingOptions := ScrapingOptions{}
			if args.SlugOverride != nil {
				scrapingOptions.SlugOverride = args.SlugOverride
			}

			actualVolumes, actualCount, err := handler.GetVolumeInfo(args.SeriesName, scrapingOptions)

			assert.Nil(t, err)
			assert.Equal(t, args.ExpectedCount, actualCount)
			assert.Equal(t, len(args.ExpectedVolumes), len(actualVolumes))

			if args.ExpectedVolumes != nil {
				for i, expectedVolume := range args.ExpectedVolumes {
					assert.Equal(t, expectedVolume.Name, actualVolumes[i].Name)
					if expectedVolume.ReleaseDate != nil && actualVolumes[i].ReleaseDate != nil {
						assert.Equal(t, expectedVolume.ReleaseDate.Format(releaseDateFormat), actualVolumes[i].ReleaseDate.Format(releaseDateFormat))
					} else {
						assert.Equal(t, expectedVolume.ReleaseDate, actualVolumes[i].ReleaseDate)
					}
				}
			}
		})
	}
}

func TimePtr(t time.Time) *time.Time {
	return &t
}
