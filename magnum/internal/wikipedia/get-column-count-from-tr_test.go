//go:build unit

package wikipedia_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/wikipedia"
	"github.com/stretchr/testify/assert"
)

type GetColumnCountFromTrTestCase struct {
	InputTr             string
	ExpectedTdCount     int
	ExpectedActualCount int
}

var GetColumnCountFromTrTestCases = map[string]GetColumnCountFromTrTestCase{
	"a simple table row should get the correct amount of rows returned": {
		InputTr:             `<tr style="text-align: center;"><th scope="row" id="vol1" style="text-align: center; font-weight: normal; background-color: transparent;">1</th><td> August 22, 2013<sup id="cite_ref-3" class="reference"><a href="#cite_note-3">[3]</a></sup></td><td><style data-mw-deduplicate="TemplateStyles:r1215172403">.mw-parser-output cite.citation{font-style:inherit;word-wrap:break-word}.mw-parser-output .citation q{quotes:"\"""\"""'""'"}.mw-parser-output .citation:target{background-color:rgba(0,127,255,0.133)}.mw-parser-output .id-lock-free.id-lock-free a{background:url("//upload.wikimedia.org/wikipedia/commons/6/65/Lock-green.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-free a{background-size:contain}.mw-parser-output .id-lock-limited.id-lock-limited a,.mw-parser-output .id-lock-registration.id-lock-registration a{background:url("//upload.wikimedia.org/wikipedia/commons/d/d6/Lock-gray-alt-2.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-limited a,body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-registration a{background-size:contain}.mw-parser-output .id-lock-subscription.id-lock-subscription a{background:url("//upload.wikimedia.org/wikipedia/commons/a/aa/Lock-red-alt-2.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-subscription a{background-size:contain}.mw-parser-output .cs1-ws-icon a{background:url("//upload.wikimedia.org/wikipedia/commons/4/4c/Wikisource-logo.svg")right 0.1em center/12px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .cs1-ws-icon a{background-size:contain}.mw-parser-output .cs1-code{color:inherit;background:inherit;border:none;padding:inherit}.mw-parser-output .cs1-hidden-error{display:none;color:#d33}.mw-parser-output .cs1-visible-error{color:#d33}.mw-parser-output .cs1-maint{display:none;color:#2C882D;margin-left:0.3em}.mw-parser-output .cs1-format{font-size:95%}.mw-parser-output .cs1-kern-left{padding-left:0.2em}.mw-parser-output .cs1-kern-right{padding-right:0.2em}.mw-parser-output .citation .mw-selflink{font-weight:inherit}html.skin-theme-clientpref-night .mw-parser-output .cs1-maint{color:#18911F}html.skin-theme-clientpref-night .mw-parser-output .cs1-visible-error,html.skin-theme-clientpref-night .mw-parser-output .cs1-hidden-error{color:#f8a397}@media(prefers-color-scheme:dark){html.skin-theme-clientpref-os .mw-parser-output .cs1-visible-error,html.skin-theme-clientpref-os .mw-parser-output .cs1-hidden-error{color:#f8a397}html.skin-theme-clientpref-os .mw-parser-output .cs1-maint{color:#18911F}}</style><a href="/wiki/Special:BookSources/978-4-8401-5275-4" title="Special:BookSources/978-4-8401-5275-4">978-4-8401-5275-4</a></td><td>September 15, 2015<sup id="cite_ref-US_Book_ISBN_4-0" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-935548-72-0" title="Special:BookSources/978-1-935548-72-0">978-1-935548-72-0</a></td></tr>`,
		ExpectedTdCount:     4,
		ExpectedActualCount: 4,
	},
	"a row with a colspan should be handled properly": {
		InputTr: `  <tr>
		<td></td>
    <td colspan="2">Sum: $180</td>
		<td></td>
  </tr>`,
		ExpectedTdCount:     3,
		ExpectedActualCount: 4,
	},
	"a row with an empty colspan should be handled properly": {
		InputTr: `  <tr>
		<td></td>
    <td colspan="">Sum: $180</td>
		<td></td>
  </tr>`,
		ExpectedTdCount:     3,
		ExpectedActualCount: 3,
	},
	"a row with multiple colspans should be handled properly": {
		InputTr: `  <tr>
		<td></td>
    <td colspan="2">Sum: $180</td>
		<td colspan="4"></td>
		<td></td>
  </tr>`,
		ExpectedTdCount:     4,
		ExpectedActualCount: 8,
	},
}

func TestGetColumnCountFromTr(t *testing.T) {
	for name, args := range GetColumnCountFromTrTestCases {
		t.Run(name, func(t *testing.T) {
			actualNumTds, actualColumnNum, err := wikipedia.GetColumnCountFromTr(args.InputTr)

			assert.Nil(t, err)
			assert.Equal(t, args.ExpectedActualCount, actualColumnNum, "actual column value was not the expected value")
			assert.Equal(t, args.ExpectedTdCount, actualNumTds, "actual number of tds was not the expected value")
		})
	}
}
