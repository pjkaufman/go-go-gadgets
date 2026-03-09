//go:build unit

package converter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type htmlToTextTestCase struct {
	inputText     string
	expectedLines []string
}

var (
	expectedGloriousDayLines = []string{
		"A Glorious Church",
		"Ralph E. Hudson (MS68)",
		"~ 1 ~ Do you hear them coming, Brother, Thronging up the steeps of light,",
		"Clad in glorious shining garments Blood-washed garments pure and white?",
		"CHORUS:",
		"‘Tis a glorious Church without spot or wrinkle, Washed in the blood of the Lamb.",
		"‘Tis a glorious Church, without spot or wrinkle, Washed in the blood of the Lamb.",
		"~ 2 ~ Do you hear the stirring anthems Filling all the earth and sky?",
		"‘Tis a grand victorious army. Lift its banner up on high!",
		"~ 3 ~ Never fear the clouds of sorrow; Never fear the storms of sin.",
		"Even now our joys begin.",
		"~ 4 ~ Wave the banner, shout His praises, For our victory is nigh!",
		"We shall join our conquering Savior. We shall reign with Him on high.",
	}
	aGloriousDayBrHtml = `<div class="keep-together">
<h1 id="a-glorious-church">A Glorious Church</h1>
<div><div class="metadata"><div><div class="author">Ralph E. Hudson</div></div><div><div class="key">&nbsp;&nbsp;&nbsp;&nbsp;</div></div><div><div class="location">(MS68)</div></div></div></div><br>
<p>~ 1 ~ Do you hear them coming, Brother, Thronging up the steeps of light,<br>
Clad in glorious shining garments Blood-washed garments pure and white?</p>
<p>CHORUS:<br>
&lsquo;Tis a glorious Church without spot or wrinkle, Washed in the blood of the Lamb.<br>
&lsquo;Tis a glorious Church, without spot or wrinkle, Washed in the blood of the Lamb.</p>
<p>~ 2 ~ Do you hear the stirring anthems Filling all the earth and sky?<br>
&lsquo;Tis a grand victorious army. Lift its banner up on high!</p>
<p>~ 3 ~ Never fear the clouds of sorrow; Never fear the storms of sin.<br>
Even now our joys begin.</p>
<p>~ 4 ~ Wave the banner, shout His praises, For our victory is nigh!<br>
We shall join our conquering Savior. We shall reign with Him on high.</p>
</div>`
	aGloriousDayBrSlashHtml = `<div class="keep-together">
<h1 id="a-glorious-church">A Glorious Church</h1>
<div><div class="metadata"><div><div class="author">Ralph E. Hudson</div></div><div><div class="key">&nbsp;&nbsp;&nbsp;&nbsp;</div></div><div><div class="location">(MS68)</div></div></div></div><br/>
<p>~ 1 ~ Do you hear them coming, Brother, Thronging up the steeps of light,<br/>
Clad in glorious shining garments Blood-washed garments pure and white?</p>
<p>CHORUS:<br>
&lsquo;Tis a glorious Church without spot or wrinkle, Washed in the blood of the Lamb.<br/>
&lsquo;Tis a glorious Church, without spot or wrinkle, Washed in the blood of the Lamb.</p>
<p>~ 2 ~ Do you hear the stirring anthems Filling all the earth and sky?<br/>
&lsquo;Tis a grand victorious army. Lift its banner up on high!</p>
<p>~ 3 ~ Never fear the clouds of sorrow; Never fear the storms of sin.<br/>
Even now our joys begin.</p>
<p>~ 4 ~ Wave the banner, shout His praises, For our victory is nigh!<br/>
We shall join our conquering Savior. We shall reign with Him on high.</p>
</div>`
)

var htmlToTextTestCases = map[string]htmlToTextTestCase{
	"A sample song input should get properly converted back to lines based on the lines in the html": {
		inputText:     aGloriousDayBrHtml,
		expectedLines: expectedGloriousDayLines,
	},
	"A song should have the same output for <br> and <br/> for a line break in the HTML": {
		inputText:     aGloriousDayBrHtml + "\n" + aGloriousDayBrSlashHtml,
		expectedLines: append(expectedGloriousDayLines, expectedGloriousDayLines...),
	},
	"A file with an anchor tag without text, but a valid href to an id should get the text for the id in question": {
		inputText:     `<li><a href="#a-glorious-church"></a></li>` + "\n" + aGloriousDayBrHtml,
		expectedLines: append([]string{"A Glorious Church"}, expectedGloriousDayLines...),
	},
	"A file with an anchor tag without text, but a valid href to an id that has another title should get the text for the id in question removing the corresponding span tags": {
		inputText: `<li><a href="#test-song"></a></li>
		<div>
			<h1 id="test-song">Song&rsquo;s Name 1 <span class="other-title">(Other Name)</span></h1>
			<p>Line 1<br>
			Line 2<br>
			Line 3</p>
		</div>`,
		expectedLines: []string{
			"Song’s Name 1 (Other Name)",
			"Song’s Name 1 (Other Name)",
			"Line 1",
			"Line 2",
			"Line 3",
		},
	},
}

func TestHtmlToText(t *testing.T) {
	for name, args := range htmlToTextTestCases {
		t.Run(name, func(t *testing.T) {
			actual := converter.HtmlToText(args.inputText)

			assert.Equal(t, args.expectedLines, actual)
		})
	}
}
