//go:build unit

package slug_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/slug"
	"github.com/stretchr/testify/assert"
)

type GetSeriesSlugFromNameTestCase struct {
	SeriesName   string
	ExpectedSlug string
}

var GetSeriesSlugFromNameTestCases = map[string]GetSeriesSlugFromNameTestCase{
	"a simple name with no characters to remove from the title should properly have the spaces converted to dashes": {
		SeriesName:   "Berserk of Gluttony",
		ExpectedSlug: "berserk-of-gluttony",
	},
	"a name with an exclamation point and comma has the exclamation point and comma properly removed and dashes added": {
		SeriesName:   "No Game No Life, Please!",
		ExpectedSlug: "no-game-no-life-please",
	},
	"a name with a question mark, apostrophe, and parens has them properly removed and dashes added": {
		SeriesName:   "So I'm a Spider, So What? (light novel)",
		ExpectedSlug: "so-i-m-a-spider-so-what-light-novel",
	},
	"a name with a dash has it properly removed and dashes added": {
		SeriesName:   "The Devil Is a Part-Timer!",
		ExpectedSlug: "the-devil-is-a-part-timer",
	},
	"a name with a double dash has it properly condensed down to a single dash and dashes added": {
		SeriesName:   "86--EIGHTY-SIX (light novel)",
		ExpectedSlug: "86-eighty-six-light-novel",
	},
	"a name with a colon has it properly removed and dashes added": {
		SeriesName:   "Is It Wrong to Try to Pick Up Girls in a Dungeon? On the Side: Sword Oratoria (light novel)",
		ExpectedSlug: "is-it-wrong-to-try-to-pick-up-girls-in-a-dungeon-on-the-side-sword-oratoria-light-novel",
	},
}

func TestGetSeriesSlugFromName(t *testing.T) {
	for name, args := range GetSeriesSlugFromNameTestCases {
		t.Run(name, func(t *testing.T) {
			actualSlug := slug.GetSeriesSlugFromName(args.SeriesName)

			assert.Equal(t, args.ExpectedSlug, actualSlug)
		})
	}
}
