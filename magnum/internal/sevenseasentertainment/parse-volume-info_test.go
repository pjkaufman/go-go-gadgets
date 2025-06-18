//go:build unit

package sevenseasentertainment_test

import (
	"testing"
	"time"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/sevenseasentertainment"
	"github.com/stretchr/testify/assert"
)

type ParseVolumeInfoTestCase struct {
	InputHtml          string
	InputSeriesName    string
	InputVolumeNum     int
	ExpectedVolumeInfo *sevenseasentertainment.VolumeInfo
	ExpectError        bool
}

const (
	mushokuTenseiVolume1 = `                        
            <img src="https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT.jpg" alt="">
            
            <h3>Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 1</h3>

            <b>Release Date</b>: May 21, 2019<br>

                        <b>Early Digital:</b> Apr 04, 2019<br>
            
            <b>Price:</b> $15.99<br>
            <b>Format:</b> Light Novel<br>
            <b>ISBN:</b> 978-1-64275-138-3`
	mushokuTenseiAudiobook1 = `                          <img src="https://sevenseasentertainment.com/wp-content/uploads/2023/07/MushokuTensei1_audiobook_cover-site.jpg" alt="">
                <h3>Mushoku Tensei: Jobless Reincarnation (Audiobook) Vol. 1</h3>
                <b>Release Date</b>: 2023-09-28<br>
                <b>Length:</b> 7 hrs 18 min            `
	noTitle = `                        
            <img src="https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT.jpg" alt="">
            
            <h3></h3>

            <b>Release Date</b>: May 21, 2019<br>

                        <b>Early Digital:</b> Apr 04, 2019<br>
            
            <b>Price:</b> $15.99<br>
            <b>Format:</b> Light Novel<br>
            <b>ISBN:</b> 978-1-64275-138-3`
	noRelease = `                        
            <img src="https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT.jpg" alt="">
            
            <h3>Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 1</h3>

            
            <b>Price:</b> $15.99<br>
            <b>Format:</b> Light Novel<br>
            <b>ISBN:</b> 978-1-64275-138-3`
	mushokuTenseiVolume14 = `                      
            <img src="https://sevenseasentertainment.com/wp-content/uploads/2022/01/mushokuLN14_site-resize.jpg" alt="">
            
            <h3>Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 14</h3>

            <b>Release Date</b>: Jan 18, 2022<br>
            
            <b>Price:</b> $15.99<br>
            <b>Format:</b> Light Novel<br>
            <b>ISBN:</b> 978-1-64827-360-5`
)

var ParseVolumeInfoTestCases = map[string]ParseVolumeInfoTestCase{
	"a simple volume with a release date and early digital date should properly parse the early digital date as the release date and pull the proper name": {
		InputHtml:       mushokuTenseiVolume1,
		InputSeriesName: "Mushoku Tensei",
		ExpectedVolumeInfo: &sevenseasentertainment.VolumeInfo{
			Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 1",
			ReleaseDate: getDatePointer(2019, 4, time.April),
		},
	},
	"if a volume is missing the early digital release date info the release date should be used instead, the release date is left as nil": {
		InputHtml:       mushokuTenseiVolume14,
		InputSeriesName: "Mushoku Tensei",
		ExpectedVolumeInfo: &sevenseasentertainment.VolumeInfo{
			Name:        "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 14",
			ReleaseDate: getDatePointer(2022, 18, time.January),
		},
	},
	"if a volume is missing the release date info (both early digital and regular release), the release date is left as nil": {
		InputHtml:       noRelease,
		InputSeriesName: "Mushoku Tensei",
		ExpectedVolumeInfo: &sevenseasentertainment.VolumeInfo{
			Name: "Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 1",
		},
	},
	"if a volume has audiobook in the title it should be skipped with a nil value returned": {
		InputHtml:          mushokuTenseiAudiobook1,
		InputSeriesName:    "Mushoku Tensei",
		ExpectedVolumeInfo: nil,
	},
	"if a volume is missing the title, an error occurs": {
		InputHtml:          noTitle,
		InputSeriesName:    "Mushoku Tensei",
		ExpectedVolumeInfo: nil,
		ExpectError:        true,
	},
}

func TestParseWikipediaTableToVolumeInfo(t *testing.T) {
	for name, args := range ParseVolumeInfoTestCases {
		t.Run(name, func(t *testing.T) {
			actualVolumeInfo, err := sevenseasentertainment.ParseVolumeInfo(args.InputSeriesName, args.InputHtml, args.InputVolumeNum)

			assert.Equal(t, err != nil, args.ExpectError)
			assert.Equal(t, args.ExpectedVolumeInfo != nil, actualVolumeInfo != nil)

			if !args.ExpectError && args.ExpectedVolumeInfo != nil {
				assert.Equal(t, args.ExpectedVolumeInfo.Name, actualVolumeInfo.Name)
				assert.Equal(t, args.ExpectedVolumeInfo.ReleaseDate, actualVolumeInfo.ReleaseDate)
			}
		})
	}
}

func getDatePointer(year, day int, month time.Month) *time.Time {
	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	return &date
}
