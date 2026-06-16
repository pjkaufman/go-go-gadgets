//go:build unit

package yenpress_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/yenpress"
	"github.com/stretchr/testify/assert"
)

type parseVolumeInfoTestCase struct {
	inputHtml          string
	inputSeriesName    string
	expectedVolumeInfo *yenpress.VolumeInfo
	expectError        bool
}

const (
	danmachiColectorsEdition = `                        
        	<a href="/titles/9798855411362-is-it-wrong-to-try-to-pick-up-girls-in-a-dungeon-collector-s-edition-vol-1" class="hovered-shadow">
            <img class="four-swiper-img img-box-shadow b-lazy b-error" src="https://images.yenpress.com/imgs/9798855411362.jpg?w=285&amp;h=422&amp;type=books&amp;s=63b8273a1011bd02f508abddf011e7d0" alt="Is It Wrong to Try to Pick Up Girls in a Dungeon? Collector's Edition, Vol. 1">
            <p class="paragraph"><span>Is It Wrong to Try to Pick Up Girls in a Dungeon? Collector&#39;s Edition, Vol. 1</span></p>
			</a>`
	danmachi19 = `
			<a href="/titles/9781975393403-is-it-wrong-to-try-to-pick-up-girls-in-a-dungeon-vol-19-light-novel" class="hovered-shadow">
            <img class="four-swiper-img img-box-shadow b-lazy b-error" src="https://images.yenpress.com/imgs/9781975393403.jpg?w=285&amp;h=422&amp;type=books&amp;s=17a153d45d9271596cda3e94cc70342a" alt="Is It Wrong to Try to Pick Up Girls in a Dungeon?, Vol. 19 (light novel)">
            <p class="paragraph"><span>Is It Wrong to Try to Pick Up Girls in a Dungeon?, Vol. 19 (light novel)</span></p>
       		</a>`
	noTitle = `                        
            <a href="/titles/9781975393403-is-it-wrong-to-try-to-pick-up-girls-in-a-dungeon-vol-19-light-novel" class="hovered-shadow">
            <img class="four-swiper-img img-box-shadow b-lazy b-error" src="https://images.yenpress.com/imgs/9781975393403.jpg?w=285&amp;h=422&amp;type=books&amp;s=17a153d45d9271596cda3e94cc70342a" alt="Is It Wrong to Try to Pick Up Girls in a Dungeon?, Vol. 19 (light novel)">
			</a>`
	noRelativeLink                     = `<p class="paragraph"><span>Is It Wrong to Try to Pick Up Girls in a Dungeon?, Vol. 19 (light novel)</span></p>`
	aCertainMagicalIndexOmnibusEdition = `
        <a href="/titles/9781975351588-a-certain-magical-index-the-old-testament-omnibus-edition" class="hovered-shadow">
            <img class="four-swiper-img img-box-shadow b-lazy b-error" src="https://images.yenpress.com/imgs/9781975351588.jpg?w=285&amp;h=422&amp;type=books&amp;s=f15dea06671803ef8922897ea2d41857" alt="A Certain Magical Index: The Old Testament Omnibus Edition">
            <p class="paragraph"><span>A Certain Magical Index: The Old Testament Omnibus Edition</span></p>
        </a>`
)

var parseVolumeInfoTestCases = map[string]parseVolumeInfoTestCase{
	"a volume that has \"collector's edition\" in the name should be ignored": {
		inputHtml:          danmachiColectorsEdition,
		inputSeriesName:    "Is It Wrong to Try to Pick Up Girls in a Dungeon? (light novel)",
		expectedVolumeInfo: nil,
	},
	"a simple volume with a relative link and name is properly parsed": {
		inputHtml:       danmachi19,
		inputSeriesName: "Is It Wrong to Try to Pick Up Girls in a Dungeon? (light novel)",
		expectedVolumeInfo: &yenpress.VolumeInfo{
			Name:         "Is It Wrong to Try to Pick Up Girls in a Dungeon?, Vol. 19 (light novel)",
			RelativeLink: "/titles/9781975393403-is-it-wrong-to-try-to-pick-up-girls-in-a-dungeon-vol-19-light-novel",
		},
	},
	"if a volume is missing the name, an error will be thrown": {
		inputHtml:          noTitle,
		inputSeriesName:    "Is It Wrong to Try to Pick Up Girls in a Dungeon? (light novel)",
		expectedVolumeInfo: nil,
		expectError:        true,
	},
	"if a volume has is missing the relative link, an error will be thrown": {
		inputHtml:          noRelativeLink,
		inputSeriesName:    "Is It Wrong to Try to Pick Up Girls in a Dungeon? (light novel)",
		expectedVolumeInfo: nil,
		expectError:        true,
	},
	"a volume that has \"omnibus edition\" in the name should be ignored": {
		inputHtml:          aCertainMagicalIndexOmnibusEdition,
		inputSeriesName:    "A Certain Magical Index (light novel)",
		expectedVolumeInfo: nil,
	},
}

func TestParseWikipediaTableToVolumeInfo(t *testing.T) {
	t.Parallel()

	for name, args := range parseVolumeInfoTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actualVolumeInfo, err := yenpress.ParseVolumeInfo(args.inputSeriesName, args.inputHtml)

			assert.Equal(t, err != nil, args.expectError)
			assert.Equal(t, args.expectedVolumeInfo != nil, actualVolumeInfo != nil)

			if !args.expectError && args.expectedVolumeInfo != nil {
				assert.Equal(t, args.expectedVolumeInfo.Name, actualVolumeInfo.Name)
				assert.Equal(t, args.expectedVolumeInfo.RelativeLink, actualVolumeInfo.RelativeLink)
			}
		})
	}
}
