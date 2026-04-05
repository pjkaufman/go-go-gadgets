//go:build unit

package sitehandler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	SeriesName            string
	SlugOverride          *string
	TablesToParseOverride *int // for Wikipedia
	ExpectedVolumes       []*VolumeInfo
	ExpectedCount         int
}

type MockedEndpoint struct {
	Slug          string
	Response      string
	IsJson        bool
	CustomHandler func(w http.ResponseWriter, r *http.Request)
}

type GetVolumeInfoTestCases struct {
	Tests             map[string]GetVolumeInfoTestCase
	Endpoints         []MockedEndpoint
	CreateSiteHandler func(SiteHandlerOptions) SiteHandler
}

func createMockServerInstance(endpoints []MockedEndpoint) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(robotsFile))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	for _, endpoint := range endpoints {
		var handler func(w http.ResponseWriter, r *http.Request)
		if endpoint.CustomHandler != nil {
			handler = endpoint.CustomHandler
		} else {
			handler = func(w http.ResponseWriter, r *http.Request) {
				var contentType = "text/html"
				if endpoint.IsJson {
					contentType = "application/json"
				}

				w.Header().Set("Content-Type", contentType)

				fmt.Fprint(w, endpoint.Response)
			}
		}

		mux.HandleFunc("/"+endpoint.Slug, handler)
	}

	return httptest.NewUnstartedServer(mux)
}

func RunTests(t *testing.T, cases GetVolumeInfoTestCases) {
	srv := createMockServerInstance(cases.Endpoints)
	srv.Start()

	options := SiteHandlerOptions{
		BaseURL:        srv.URL + "/",
		Verbose:        true,
		AllowedDomains: []string{"127.0.0.1"}, // the server should run on localhost
	}

	handler := cases.CreateSiteHandler(options)

	for name, args := range cases.Tests {
		t.Run(name, func(t *testing.T) {
			scrapingOptions := ScrapingOptions{
				SlugOverride:          args.SlugOverride,
				TablesToParseOverride: args.TablesToParseOverride,
			}

			actualVolumes, actualCount, err := handler.GetVolumeInfo(args.SeriesName, scrapingOptions)

			require.NoError(t, err)
			assert.Equal(t, args.ExpectedCount, actualCount)
			assert.Len(t, actualVolumes, len(args.ExpectedVolumes))

			if len(args.ExpectedVolumes) == len(actualVolumes) {
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

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}
