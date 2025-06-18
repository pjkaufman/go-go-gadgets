//go:build unit

package yenpress_test

import (
	_ "embed"
	"testing"

	sitehandler "github.com/pjkaufman/go-go-gadgets/magnum/internal/site-handler"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/yenpress"
)

var (
	//go:embed test/the-asterisk-war.golden
	theAsteriskWarMainPage string
	//go:embed test/titles-9781975369095-the-asterisk-war-vol-17-light-novel.golden
	theAsteriskWarVolume17Page string
	//go:embed test/a-certain-magical-index-light-novel.golden
	aCertainMagicalMainPage string
	//go:embed test/titles-9781975317997-a-certain-magical-index-ss-vol-2-light-novel.golden
	aCertainMagicalSSVolume2Page string

	getVolumeInfoTestSetup = sitehandler.GetVolumeInfoTestCases{
		Tests: map[string]sitehandler.GetVolumeInfoTestCase{
			"Make sure The Asterisk War volumes are correctly extracted with 17 being found and 0 volumes returned since none are going to be released": {
				SeriesName:    "The Asterisk War",
				ExpectedCount: 17,
			},
			"Make sure A Certain Magical Index volumes are correctly extracted with 1 being found and 0 volumes returned since none are going to be released and the omnibus gets skipped": {
				SeriesName:    "A Certain Magical Index (light novel)",
				ExpectedCount: 1, // there are more novels than this, but the page says it only has 1 novel...
			},
		},
		Endpoints: []sitehandler.MockedEndpoint{
			{
				Slug:     "the-asterisk-war",
				Response: theAsteriskWarMainPage,
			},
			{
				Slug:     "titles/9781975369095-the-asterisk-war-vol-17-light-novel",
				Response: theAsteriskWarVolume17Page,
			},
			{
				Slug:     "a-certain-magical-index-light-novel",
				Response: aCertainMagicalMainPage,
			},
			{
				Slug:     "titles/9781975317997-a-certain-magical-index-ss-vol-2-light-novel",
				Response: aCertainMagicalSSVolume2Page,
			},
		},
		CreateSiteHandler: func(options sitehandler.SiteHandlerOptions) sitehandler.SiteHandler {
			return yenpress.NewYenPressHandler(options)
		},
	}
)

func TestGetVolumeInfo(t *testing.T) {
	sitehandler.RunTests(t, getVolumeInfoTestSetup)
}
