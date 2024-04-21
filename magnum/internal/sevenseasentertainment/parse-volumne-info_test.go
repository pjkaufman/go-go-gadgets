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
	mushokuTenseiVolume1    = `<div class="series-volume"> <a href="https://sevenseasentertainment.com/books/mushoku-tensei-jobless-reincarnation-light-novel-vol-1/"><img width="135" height="190" src="https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT-135x190.jpg" class="attachment-thumbnail size-thumbnail" alt="" decoding="async" srcset="https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT-135x190.jpg 135w, https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT-213x300.jpg 213w, https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT.jpg 320w" sizes="(max-width: 135px) 100vw, 135px"></a><h3><a href="https://sevenseasentertainment.com/books/mushoku-tensei-jobless-reincarnation-light-novel-vol-1/"></a><a href="https://sevenseasentertainment.com/books/mushoku-tensei-jobless-reincarnation-light-novel-vol-1/">Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 1</a></h3> <b>Release Date</b>: 2019/05/21<br> <b>Early Digital:</b> 2019/04/04<br> <b>Price:</b> $15.99<br> <b>Format:</b> Light Novel<br> <b>ISBN:</b> 978-1-64275-138-3</div>`
	mushokuTenseiAudiobook1 = `<div class="series-volume" style="min-height: 200px;"> <a href="https://sevenseasentertainment.com/audio/mushoku-tensei-jobless-reincarnation-audiobook-vol-1/"><img width="135" height="135" src="https://sevenseasentertainment.com/wp-content/uploads/2023/07/MushokuTensei1_audiobook_cover-site-135x135.jpg" class="attachment-thumbnail size-thumbnail" alt="" decoding="async" loading="lazy" srcset="https://sevenseasentertainment.com/wp-content/uploads/2023/07/MushokuTensei1_audiobook_cover-site-135x135.jpg 135w, https://sevenseasentertainment.com/wp-content/uploads/2023/07/MushokuTensei1_audiobook_cover-site.jpg 450w" sizes="(max-width: 135px) 100vw, 135px"></a><h3><a href="https://sevenseasentertainment.com/audio/mushoku-tensei-jobless-reincarnation-audiobook-vol-1/"></a><a href="https://sevenseasentertainment.com/audio/mushoku-tensei-jobless-reincarnation-audiobook-vol-1/">Mushoku Tensei: Jobless Reincarnation (Audiobook) Vol. 1</a></h3> <b>Release Date</b>: 2023/09/28<br> <b>Length:</b> 7 hrs 18 min</div>`
	noTitle                 = `<div class="series-volume"> <a href="https://sevenseasentertainment.com/books/mushoku-tensei-jobless-reincarnation-light-novel-vol-1/"><img width="135" height="190" src="https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT-135x190.jpg" class="attachment-thumbnail size-thumbnail" alt="" decoding="async" srcset="https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT-135x190.jpg 135w, https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT-213x300.jpg 213w, https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT.jpg 320w" sizes="(max-width: 135px) 100vw, 135px"></a><h3></h3> <b>Release Date</b>: 2019/05/21<br> <b>Early Digital:</b> 2019/04/04<br> <b>Price:</b> $15.99<br> <b>Format:</b> Light Novel<br> <b>ISBN:</b> 978-1-64275-138-3</div>`
	noRelease               = `<div class="series-volume"> <a href="https://sevenseasentertainment.com/books/mushoku-tensei-jobless-reincarnation-light-novel-vol-1/"><img width="135" height="190" src="https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT-135x190.jpg" class="attachment-thumbnail size-thumbnail" alt="" decoding="async" srcset="https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT-135x190.jpg 135w, https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT-213x300.jpg 213w, https://sevenseasentertainment.com/wp-content/uploads/2018/11/MUSHOKU-TENSEI-LN-1-cover-FRONT.jpg 320w" sizes="(max-width: 135px) 100vw, 135px"></a><h3><a href="https://sevenseasentertainment.com/books/mushoku-tensei-jobless-reincarnation-light-novel-vol-1/"></a><a href="https://sevenseasentertainment.com/books/mushoku-tensei-jobless-reincarnation-light-novel-vol-1/">Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 1</a></h3> <b>Price:</b> $15.99<br> <b>Format:</b> Light Novel<br> <b>ISBN:</b> 978-1-64275-138-3</div>`
	mushokuTenseiVolume14   = `<div class="series-volume"> <a href="https://sevenseasentertainment.com/books/mushoku-tensei-jobless-reincarnation-light-novel-vol-14/"><img width="135" height="190" src="https://sevenseasentertainment.com/wp-content/uploads/2022/01/mushokuLN14_site-resize-135x190.jpg" class="attachment-thumbnail size-thumbnail" alt="" decoding="async" loading="lazy" srcset="https://sevenseasentertainment.com/wp-content/uploads/2022/01/mushokuLN14_site-resize-135x190.jpg 135w, https://sevenseasentertainment.com/wp-content/uploads/2022/01/mushokuLN14_site-resize.jpg 320w" sizes="(max-width: 135px) 100vw, 135px"></a><h3><a href="https://sevenseasentertainment.com/books/mushoku-tensei-jobless-reincarnation-light-novel-vol-14/"></a><a href="https://sevenseasentertainment.com/books/mushoku-tensei-jobless-reincarnation-light-novel-vol-14/">Mushoku Tensei: Jobless Reincarnation (Light Novel) Vol. 14</a></h3> <b>Release Date</b>: 2022/01/18<br> <b>Price:</b> $15.99<br> <b>Format:</b> Light Novel<br> <b>ISBN:</b> 978-1-64827-360-5</div>`
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
