package hanashimedia_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/hanashimedia"
	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
)

var (
	//go:embed test/tsukimichi.golden
	moonlitFantasyJsonResponse string
	//go:embed test/gate-thus-the-jsdf-fought-there.golden
	gateJsonResponse       string
	getVolumeInfoTestSetup = sitehandler.GetVolumeInfoTestCases{
		Tests: map[string]sitehandler.GetVolumeInfoTestCase{
			"Make sure that Moonlit Fantasy volumes are correctly extracted": {
				SeriesName:   "Moonlit Fantasy",
				SlugOverride: new("tsukimichi"),
				ExpectedVolumes: []*sitehandler.VolumeInfo{
					{
						Name:        "Tsukimichi - Vol. 19",
						ReleaseDate: new(time.Date(2026, time.August, 30, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 18",
						ReleaseDate: new(time.Date(2026, time.July, 30, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 17",
						ReleaseDate: new(time.Date(2026, time.June, 30, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 16",
						ReleaseDate: new(time.Date(2026, time.May, 30, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 15",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 14",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 13",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 12",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 11",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 10",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 9",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 8.5",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 8",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 7",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 6",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 5",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 4",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 3",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 2",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "Tsukimichi - Vol. 1",
						ReleaseDate: new(time.Date(2026, time.March, 28, 0, 0, 0, 0, time.Local)),
					},
				},
				ExpectedCount: 20,
			},
			"Make sure Gate Volumes are correctly extracted": {
				SeriesName: "Gate: Thus the JSDF Fought There!",
				ExpectedVolumes: []*sitehandler.VolumeInfo{
					{
						Name:        "GATE: Vol. 2 - Part I",
						ReleaseDate: new(time.Date(2026, time.September, 30, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "GATE: Vol. 1 - Part II",
						ReleaseDate: new(time.Date(2026, time.June, 30, 0, 0, 0, 0, time.Local)),
					},
					{
						Name:        "GATE: Vol. 1 - Part I",
						ReleaseDate: new(time.Date(2026, time.April, 6, 0, 0, 0, 0, time.Local)),
					},
				},
				ExpectedCount: 3,
			},
		},
		Endpoints: []sitehandler.MockedEndpoint{
			{
				Slug: "series",
				CustomHandler: func(w http.ResponseWriter, r *http.Request) {
					body, _ := io.ReadAll(r.Body)
					var (
						seriesAt hanashimedia.SeriesRequest
						err      = json.Unmarshal(body, &seriesAt)
					)
					if err != nil {
						http.Error(w, fmt.Sprintf("Internal Server Error: %s: %s", body, err), http.StatusInternalServerError)

						return
					}

					switch strings.ToLower(seriesAt.At) {
					case "tsukimichi":
						_, err := w.Write([]byte(moonlitFantasyJsonResponse))
						if err != nil {
							http.Error(w, "Internal Error", http.StatusInternalServerError)

							return
						}
					case "gate-thus-the-jsdf-fought-there":
						_, err := w.Write([]byte(gateJsonResponse))
						if err != nil {
							http.Error(w, "Internal Server Error", http.StatusInternalServerError)

							return
						}
					default:
						http.Error(w, "Not found", http.StatusNotFound)
					}
				},
			},
		},
		CreateSiteHandler: func(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
			return hanashimedia.NewHanashiMediaHandler(options)
		},
	}
)

func TestGetVolumeInfo(t *testing.T) {
	sitehandler.RunTests(t, getVolumeInfoTestSetup)
}
