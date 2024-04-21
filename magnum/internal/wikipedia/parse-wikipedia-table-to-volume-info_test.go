//go:build unit

package wikipedia_test

import (
	"testing"
	"time"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/wikipedia"
	"github.com/stretchr/testify/assert"
)

type ParseWikipediaTableToVolumeInfoTestCase struct {
	InputTableHtml     string
	InputNamePrefix    string
	ExpectedVolumeInfo []wikipedia.VolumeInfo
}

const (
	mushokuTensieTable = `<table class="wikitable" style="text-align: center;">

<tbody><tr>
<th rowspan="2" scope="col" width="3%">Volume no.
</th>
<th rowspan="2" scope="col" width="20%">Content
</th>
<th rowspan="2" scope="col" width="10%">Japanese release date
</th>
<th rowspan="2" scope="col" width="10%">Japanese ISBN
</th>
<th colspan="2" scope="col" width="15%">English release date
</th>
<th rowspan="2" scope="col" width="10%">English ISBN
</th></tr>
<tr>
<th scope="col">Digital
</th>
<th scope="col">Physical
</th></tr>
<tr>
<th scope="row">1
</th>
<td style="text-align: left;">Web Novel 1
</td>
<td><span data-sort-value="000000002014-01-24-0000" style="white-space:nowrap">January 24, 2014</span><sup id="cite_ref-LN_1_8-0" class="reference"><a href="#cite_note-LN_1-8">[6]</a></sup>
</td>
<td><style data-mw-deduplicate="TemplateStyles:r1133582631">.mw-parser-output cite.citation{font-style:inherit;word-wrap:break-word}.mw-parser-output .citation q{quotes:"\"""\"""'""'"}.mw-parser-output .citation:target{background-color:rgba(0,127,255,0.133)}.mw-parser-output .id-lock-free a,.mw-parser-output .citation .cs1-lock-free a{background:url("//upload.wikimedia.org/wikipedia/commons/6/65/Lock-green.svg")right 0.1em center/9px no-repeat}.mw-parser-output .id-lock-limited a,.mw-parser-output .id-lock-registration a,.mw-parser-output .citation .cs1-lock-limited a,.mw-parser-output .citation .cs1-lock-registration a{background:url("//upload.wikimedia.org/wikipedia/commons/d/d6/Lock-gray-alt-2.svg")right 0.1em center/9px no-repeat}.mw-parser-output .id-lock-subscription a,.mw-parser-output .citation .cs1-lock-subscription a{background:url("//upload.wikimedia.org/wikipedia/commons/a/aa/Lock-red-alt-2.svg")right 0.1em center/9px no-repeat}.mw-parser-output .cs1-ws-icon a{background:url("//upload.wikimedia.org/wikipedia/commons/4/4c/Wikisource-logo.svg")right 0.1em center/12px no-repeat}.mw-parser-output .cs1-code{color:inherit;background:inherit;border:none;padding:inherit}.mw-parser-output .cs1-hidden-error{display:none;color:#d33}.mw-parser-output .cs1-visible-error{color:#d33}.mw-parser-output .cs1-maint{display:none;color:#3a3;margin-left:0.3em}.mw-parser-output .cs1-format{font-size:95%}.mw-parser-output .cs1-kern-left{padding-left:0.2em}.mw-parser-output .cs1-kern-right{padding-right:0.2em}.mw-parser-output .citation .mw-selflink{font-weight:inherit}</style><a href="/wiki/Special:BookSources/978-4-04-066220-6" title="Special:BookSources/978-4-04-066220-6">978-4-04-066220-6</a>
</td>
<td><span data-sort-value="000000002019-04-04-0000" style="white-space:nowrap">April 4, 2019</span>
</td>
<td><span data-sort-value="000000002019-05-21-0000" style="white-space:nowrap">May 21, 2019</span><sup id="cite_ref-9" class="reference"><a href="#cite_note-9">[7]</a></sup>
</td>
<td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-1-64275-138-3" title="Special:BookSources/978-1-64275-138-3">978-1-64275-138-3</a>
</td></tr>
<tr>
<th scope="row">2
</th>
<td style="text-align: left;">Web Novel 2
</td>
<td><span data-sort-value="000000002014-03-25-0000" style="white-space:nowrap">March 25, 2014</span><sup id="cite_ref-LN_2_10-0" class="reference"><a href="#cite_note-LN_2-10">[8]</a></sup>
</td>
<td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-066393-7" title="Special:BookSources/978-4-04-066393-7">978-4-04-066393-7</a>
</td>
<td><span data-sort-value="000000002019-05-23-0000" style="white-space:nowrap">May 23, 2019</span>
</td>
<td><span data-sort-value="000000002019-07-30-0000" style="white-space:nowrap">July 30, 2019</span><sup id="cite_ref-11" class="reference"><a href="#cite_note-11">[9]</a></sup>
</td>
<td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-1-64275-140-6" title="Special:BookSources/978-1-64275-140-6">978-1-64275-140-6</a>
</td></tr>
<tr>
<tr>
<th scope="row">26
</th>
<td style="text-align: left;">End of Web Novel 23 and 24
</td>
<td><span data-sort-value="000000002022-11-25-0000" style="white-space:nowrap">November 25, 2022</span><sup id="cite_ref-58" class="reference"><a href="#cite_note-58">[56]</a></sup>
</td>
<td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-681933-8" title="Special:BookSources/978-4-04-681933-8">978-4-04-681933-8</a>
</td>
<td>
</td>
<td><span data-sort-value="000000002024-03-12-0000" style="white-space:nowrap">March 12, 2024</span><sup id="cite_ref-59" class="reference"><a href="#cite_note-59">[57]</a></sup>
</td>
<td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/979-8-88-843435-2" title="Special:BookSources/979-8-88-843435-2">979-8-88-843435-2</a>
</td></tr>
</tbody></table>`
	asteriskWarTable = `<table class="wikitable" width="100%" style="">


<tbody><tr style="border-bottom: 3px solid #CCF">
<th width="4%"><abbr title="Number">No.</abbr>
</th>
<th width="24%">Original release date
</th>
<th width="24%">Original ISBN
</th>
<th width="24%">English release date
</th>
<th width="24%">English ISBN
</th></tr>


<tr style="text-align: center;"><th scope="row" id="vol1" style="text-align: center; font-weight: normal; background-color: transparent;">1</th><td> September 25, 2012<sup id="cite_ref-First_14-1" class="reference"><a href="#cite_note-First-14">[7]</a></sup></td><td><style data-mw-deduplicate="TemplateStyles:r1133582631">.mw-parser-output cite.citation{font-style:inherit;word-wrap:break-word}.mw-parser-output .citation q{quotes:"\"""\"""'""'"}.mw-parser-output .citation:target{background-color:rgba(0,127,255,0.133)}.mw-parser-output .id-lock-free a,.mw-parser-output .citation .cs1-lock-free a{background:url("//upload.wikimedia.org/wikipedia/commons/6/65/Lock-green.svg")right 0.1em center/9px no-repeat}.mw-parser-output .id-lock-limited a,.mw-parser-output .id-lock-registration a,.mw-parser-output .citation .cs1-lock-limited a,.mw-parser-output .citation .cs1-lock-registration a{background:url("//upload.wikimedia.org/wikipedia/commons/d/d6/Lock-gray-alt-2.svg")right 0.1em center/9px no-repeat}.mw-parser-output .id-lock-subscription a,.mw-parser-output .citation .cs1-lock-subscription a{background:url("//upload.wikimedia.org/wikipedia/commons/a/aa/Lock-red-alt-2.svg")right 0.1em center/9px no-repeat}.mw-parser-output .cs1-ws-icon a{background:url("//upload.wikimedia.org/wikipedia/commons/4/4c/Wikisource-logo.svg")right 0.1em center/12px no-repeat}.mw-parser-output .cs1-code{color:inherit;background:inherit;border:none;padding:inherit}.mw-parser-output .cs1-hidden-error{display:none;color:#d33}.mw-parser-output .cs1-visible-error{color:#d33}.mw-parser-output .cs1-maint{display:none;color:#3a3;margin-left:0.3em}.mw-parser-output .cs1-format{font-size:95%}.mw-parser-output .cs1-kern-left{padding-left:0.2em}.mw-parser-output .cs1-kern-right{padding-right:0.2em}.mw-parser-output .citation .mw-selflink{font-weight:inherit}</style><a href="/wiki/Special:BookSources/978-4-04-066697-6" title="Special:BookSources/978-4-04-066697-6">978-4-04-066697-6</a></td><td>August 30, 2016<sup id="cite_ref-Eng_ver1_16-1" class="reference"><a href="#cite_note-Eng_ver1-16">[9]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-0-31-631527-2" title="Special:BookSources/978-0-31-631527-2">978-0-31-631527-2</a></td></tr>
<tr style="text-align: center;"><th scope="row" id="vol2" style="text-align: center; font-weight: normal; background-color: transparent;">2</th><td> January 25, 2013<sup id="cite_ref-18" class="reference"><a href="#cite_note-18">[11]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-066698-3" title="Special:BookSources/978-4-04-066698-3">978-4-04-066698-3</a></td><td>December 20, 2016<sup id="cite_ref-19" class="reference"><a href="#cite_note-19">[12]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-0-31-639858-9" title="Special:BookSources/978-0-31-639858-9">978-0-31-639858-9</a></td></tr>
<tr style="text-align: center;"><th scope="row" id="vol3" style="text-align: center; font-weight: normal; background-color: transparent;">3</th><td> May 24, 2013<sup id="cite_ref-20" class="reference"><a href="#cite_note-20">[13]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-067093-5" title="Special:BookSources/978-4-04-067093-5">978-4-04-067093-5</a></td><td>April 18, 2017<sup id="cite_ref-21" class="reference"><a href="#cite_note-21">[14]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-0-31-639860-2" title="Special:BookSources/978-0-31-639860-2">978-0-31-639860-2</a></td></tr>
<tr style="text-align: center;"><th scope="row" id="vol4" style="text-align: center; font-weight: normal; background-color: transparent;">4</th><td> September 25, 2013<sup id="cite_ref-22" class="reference"><a href="#cite_note-22">[15]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-84-015417-8" title="Special:BookSources/978-4-84-015417-8">978-4-84-015417-8</a></td><td>August 22, 2017<sup id="cite_ref-23" class="reference"><a href="#cite_note-23">[16]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-0-31-639862-6" title="Special:BookSources/978-0-31-639862-6">978-0-31-639862-6</a></td></tr>
</tbody></table>`
	wrongWayToHealTable = `<table class="wikitable" width="100%" style="">


<tbody><tr style="border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom="">
<th width="4%"><abbr title="Number">No.</abbr>
</th>
<th width="24%">Original release date
</th>
<th width="24%">Original ISBN
</th>
<th width="24%">English release date
</th>
<th width="24%">English ISBN
</th></tr>


<tr style="text-align: center;"><th scope="row" id="vol1" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">1</th><td> March 25, 2016<sup id="cite_ref-8" class="reference"><a href="#cite_note-8">[8]</a></sup></td><td><style data-mw-deduplicate="TemplateStyles:r1133582631">.mw-parser-output cite.citation{font-style:inherit;word-wrap:break-word}.mw-parser-output .citation q{quotes:"\"""\"""'""'"}.mw-parser-output .citation:target{background-color:rgba(0,127,255,0.133)}.mw-parser-output .id-lock-free a,.mw-parser-output .citation .cs1-lock-free a{background:url("//upload.wikimedia.org/wikipedia/commons/6/65/Lock-green.svg")right 0.1em center/9px no-repeat}.mw-parser-output .id-lock-limited a,.mw-parser-output .id-lock-registration a,.mw-parser-output .citation .cs1-lock-limited a,.mw-parser-output .citation .cs1-lock-registration a{background:url("//upload.wikimedia.org/wikipedia/commons/d/d6/Lock-gray-alt-2.svg")right 0.1em center/9px no-repeat}.mw-parser-output .id-lock-subscription a,.mw-parser-output .citation .cs1-lock-subscription a{background:url("//upload.wikimedia.org/wikipedia/commons/a/aa/Lock-red-alt-2.svg")right 0.1em center/9px no-repeat}.mw-parser-output .cs1-ws-icon a{background:url("//upload.wikimedia.org/wikipedia/commons/4/4c/Wikisource-logo.svg")right 0.1em center/12px no-repeat}.mw-parser-output .cs1-code{color:inherit;background:inherit;border:none;padding:inherit}.mw-parser-output .cs1-hidden-error{display:none;color:#d33}.mw-parser-output .cs1-visible-error{color:#d33}.mw-parser-output .cs1-maint{display:none;color:#3a3;margin-left:0.3em}.mw-parser-output .cs1-format{font-size:95%}.mw-parser-output .cs1-kern-left{padding-left:0.2em}.mw-parser-output .cs1-kern-right{padding-right:0.2em}.mw-parser-output .citation .mw-selflink{font-weight:inherit}</style><style class="darkreader darkreader--sync" media="screen"></style><a href="/wiki/Special:BookSources/978-4-04-068185-6" title="Special:BookSources/978-4-04-068185-6">978-4-04-068185-6</a></td><td>August 23, 2022<sup id="cite_ref-9" class="reference"><a href="#cite_note-9">[9]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-1-64273-200-9" title="Special:BookSources/978-1-64273-200-9">978-1-64273-200-9</a></td></tr>
<tr style="text-align: center;"><th scope="row" id="vol2" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">2</th><td> June 24, 2016<sup id="cite_ref-10" class="reference"><a href="#cite_note-10">[10]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-068427-7" title="Special:BookSources/978-4-04-068427-7">978-4-04-068427-7</a></td><td>May 15, 2023<sup id="cite_ref-11" class="reference"><a href="#cite_note-11">[11]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-1-64273-232-0" title="Special:BookSources/978-1-64273-232-0">978-1-64273-232-0</a></td></tr>
<tr style="text-align: center;"><th scope="row" id="vol3" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">3</th><td> September 23, 2016<sup id="cite_ref-12" class="reference"><a href="#cite_note-12">[12]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-068636-3" title="Special:BookSources/978-4-04-068636-3">978-4-04-068636-3</a></td><td>August 22, 2023</td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-1-64273-286-3" title="Special:BookSources/978-1-64273-286-3">978-1-64273-286-3</a></td></tr>
<tr style="text-align: center;"><th scope="row" id="vol4" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">4</th><td> January 25, 2017<sup id="cite_ref-13" class="reference"><a href="#cite_note-13">[13]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-069054-4" title="Special:BookSources/978-4-04-069054-4">978-4-04-069054-4</a></td><td>—</td><td>—</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol5" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">5</th><td> April 25, 2017<sup id="cite_ref-14" class="reference"><a href="#cite_note-14">[14]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-069191-6" title="Special:BookSources/978-4-04-069191-6">978-4-04-069191-6</a></td><td>—</td><td>—</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol6" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">6</th><td> September 25, 2017<sup id="cite_ref-15" class="reference"><a href="#cite_note-15">[15]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-069498-6" title="Special:BookSources/978-4-04-069498-6">978-4-04-069498-6</a></td><td>—</td><td>—</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol7" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">7</th><td> February 24, 2018<sup id="cite_ref-16" class="reference"><a href="#cite_note-16">[16]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-069727-7" title="Special:BookSources/978-4-04-069727-7">978-4-04-069727-7</a></td><td>—</td><td>—</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol8" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">8</th><td> July 25, 2018<sup id="cite_ref-17" class="reference"><a href="#cite_note-17">[17]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-065023-4" title="Special:BookSources/978-4-04-065023-4">978-4-04-065023-4</a></td><td>—</td><td>—</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol9" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">9</th><td> November 24, 2018<sup id="cite_ref-18" class="reference"><a href="#cite_note-18">[18]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-065306-8" title="Special:BookSources/978-4-04-065306-8">978-4-04-065306-8</a></td><td>—</td><td>—</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol10" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">10</th><td> April 25, 2019<sup id="cite_ref-19" class="reference"><a href="#cite_note-19">[19]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-065682-3" title="Special:BookSources/978-4-04-065682-3">978-4-04-065682-3</a></td><td>—</td><td>—</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol11" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">11</th><td> October 25, 2019<sup id="cite_ref-20" class="reference"><a href="#cite_note-20">[20]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-064060-0" title="Special:BookSources/978-4-04-064060-0">978-4-04-064060-0</a></td><td>—</td><td>—</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol12" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">12</th><td> March 25, 2020<sup id="cite_ref-21" class="reference"><a href="#cite_note-21">[21]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1133582631"><a href="/wiki/Special:BookSources/978-4-04-064538-4" title="Special:BookSources/978-4-04-064538-4">978-4-04-064538-4</a></td><td>—</td><td>—</td></tr>
</tbody></table>`
	risingOfTheShieldHeroTable = `<table class="wikitable" width="98%" style="">


<tbody><tr style="border-bottom: 3px solid #1F3AAD">
<th width="4%"><abbr title="Number">No.</abbr>
</th>
<th width="24%">Original release date
</th>
<th width="24%">Original ISBN
</th>
<th width="24%">English release date
</th>
<th width="24%">English ISBN
</th></tr><tr style="text-align: center;"><th scope="row" id="vol1" style="text-align: center; font-weight: normal; background-color: transparent;">1</th><td> August 22, 2013<sup id="cite_ref-3" class="reference"><a href="#cite_note-3">[3]</a></sup></td><td><style data-mw-deduplicate="TemplateStyles:r1215172403">.mw-parser-output cite.citation{font-style:inherit;word-wrap:break-word}.mw-parser-output .citation q{quotes:"\"""\"""'""'"}.mw-parser-output .citation:target{background-color:rgba(0,127,255,0.133)}.mw-parser-output .id-lock-free.id-lock-free a{background:url("//upload.wikimedia.org/wikipedia/commons/6/65/Lock-green.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-free a{background-size:contain}.mw-parser-output .id-lock-limited.id-lock-limited a,.mw-parser-output .id-lock-registration.id-lock-registration a{background:url("//upload.wikimedia.org/wikipedia/commons/d/d6/Lock-gray-alt-2.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-limited a,body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-registration a{background-size:contain}.mw-parser-output .id-lock-subscription.id-lock-subscription a{background:url("//upload.wikimedia.org/wikipedia/commons/a/aa/Lock-red-alt-2.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-subscription a{background-size:contain}.mw-parser-output .cs1-ws-icon a{background:url("//upload.wikimedia.org/wikipedia/commons/4/4c/Wikisource-logo.svg")right 0.1em center/12px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .cs1-ws-icon a{background-size:contain}.mw-parser-output .cs1-code{color:inherit;background:inherit;border:none;padding:inherit}.mw-parser-output .cs1-hidden-error{display:none;color:#d33}.mw-parser-output .cs1-visible-error{color:#d33}.mw-parser-output .cs1-maint{display:none;color:#2C882D;margin-left:0.3em}.mw-parser-output .cs1-format{font-size:95%}.mw-parser-output .cs1-kern-left{padding-left:0.2em}.mw-parser-output .cs1-kern-right{padding-right:0.2em}.mw-parser-output .citation .mw-selflink{font-weight:inherit}html.skin-theme-clientpref-night .mw-parser-output .cs1-maint{color:#18911F}html.skin-theme-clientpref-night .mw-parser-output .cs1-visible-error,html.skin-theme-clientpref-night .mw-parser-output .cs1-hidden-error{color:#f8a397}@media(prefers-color-scheme:dark){html.skin-theme-clientpref-os .mw-parser-output .cs1-visible-error,html.skin-theme-clientpref-os .mw-parser-output .cs1-hidden-error{color:#f8a397}html.skin-theme-clientpref-os .mw-parser-output .cs1-maint{color:#18911F}}</style><a href="/wiki/Special:BookSources/978-4-8401-5275-4" title="Special:BookSources/978-4-8401-5275-4">978-4-8401-5275-4</a></td><td>September 15, 2015<sup id="cite_ref-US_Book_ISBN_4-0" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-935548-72-0" title="Special:BookSources/978-1-935548-72-0">978-1-935548-72-0</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid #CCF;"><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left;"><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black;">
<ul><li>Chapter One: "A Royal Summons"</li>
<li>Chapter Two: "The Heroes"</li>
<li>Chapter Three: "A Heroic Discussion"</li>
<li>Chapter Four: "Specially-Arranged Funding"</li>
<li>Chapter Five: "The Reality of the Shield"</li>
<li>Chapter Six: "A Backstabber Named Landmine"</li>
<li>Chapter Seven: "False Charges"</li>
<li>Chapter Eight: "Ruined Reputation"</li>
<li>Chapter Nine: "They Call it a Slave"</li>
<li>Chapter Ten: "Kids' Menu"</li>
<li>Chapter Eleven: "The Fruits of Slavery"</li>
<li>Chapter Twelve: "What's Yours is Mine"</li>
<li>Chapter Thirteen: "Medicine"</li></ul></td>
<td style="border: black; width: 48%;">
<ul><li>Chapter Fourteen: "To Take a Life"</li>
<li>Chapter Fifteen: "Demi-Humans"</li>
<li>Chapter Sixteen: "The Two-Headed, Black Dog"</li>
<li>Chapter Seventeen: "Preparing for the Wave"</li>
<li>Chapter Eighteen: "Barbarian Armor"</li>
<li>Chapter Nineteen: "The Dragon Hourglass"</li>
<li>Chapter Twenty: "The Sword"</li>
<li>Chapter Twenty-One: "The Wave of Destruction"</li>
<li>Chapter Twenty-Two: "The Clash of Spear and Shield"</li>
<li>Chapter Twenty-Three: "All I'd Wanted To Hear"</li>
<li>Epilogue</li>
<li>Special Extra Chapter One: "The Spear Hero's Buffoonery"</li>
<li>Special Extra Chapter Two: "The Flag on the Kid's Meal"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol2" style="text-align: center; font-weight: normal; background-color: transparent;">2</th><td> October 1, 2013<sup id="cite_ref-5" class="reference"><a href="#cite_note-5">[5]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-066049-3" title="Special:BookSources/978-4-04-066049-3">978-4-04-066049-3</a></td><td>October 20, 2015<sup id="cite_ref-US_Book_ISBN_4-1" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-935548-78-2" title="Special:BookSources/978-1-935548-78-2">978-1-935548-78-2</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid #CCF;"><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left;"><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black;">
<ul><li>Prologue: "Shared Pain"</li>
<li>Chapter One: "Egg Machine"</li>
<li>Chapter Two: "Gratitude for Life"</li>
<li>Chapter Three: "Filo"</li>
<li>Chapter Four: "Growth"</li>
<li>Chapter Five: "Kick and Run"</li>
<li>Chapter Six: "Wings"</li>
<li>Chapter Seven: "Transformation"</li>
<li>Chapter Eight: "Carrot and Stick"</li>
<li>Chapter Nine: "Rewards"</li>
<li>Chapter Ten: "Traveling Merchant"</li>
<li>Chapter Eleven: "Travel by Carriage"</li></ul></td>
<td style="border: black; width: 48%;">
<ul><li>Chapter Twelve: "Rumors of the Heroes"</li>
<li>Chapter Thirteen: "Take Everything but Life"</li>
<li>Chapter Fourteen: "Magic Practice"</li>
<li>Chapter Fifteen: "Why it Was Sealed"</li>
<li>Chapter Sixteen: "Invading Vines"</li>
<li>Chapter Seventeen: "Improving the Product Line"</li>
<li>Chapter Eighteen: "Diseased Village"</li>
<li>Chapter Nineteen: "Curse Series"</li>
<li>Chapter Twenty: "The Shield of Rage"</li>
<li>Epilogue: "As a Shield..."</li>
<li>Special Extra Chapter: "Presents"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol3" style="text-align: center; font-weight: normal; background-color: transparent;">3</th><td> December 21, 2013<sup id="cite_ref-6" class="reference"><a href="#cite_note-6">[6]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-066166-7" title="Special:BookSources/978-4-04-066166-7">978-4-04-066166-7</a></td><td>February 16, 2016<sup id="cite_ref-US_Book_ISBN_4-2" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-935548-66-9" title="Special:BookSources/978-1-935548-66-9">978-1-935548-66-9</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid #CCF;"><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left;"><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black;">
<ul><li>Prologue</li>
<li>Chapter One: "Filo's Friend"</li>
<li>Chapter Two: "The Fruits of Peddling"</li>
<li>Chapter Three: "Everyone Loves Angels"</li>
<li>Chapter Four: "The Volunteer"</li>
<li>Chapter Five: "A Royal Order"</li>
<li>Chapter Six: "Welcome"</li>
<li>Chapter Seven: "General Commander"</li>
<li>Chapter Eight: "Before the Storm"</li>
<li>Chapter Nine: "Framed Again?"</li>
<li>Chapter Ten: "The Third Wave"</li>
<li>Chapter Eleven: "Grow Up"</li></ul></td>
<td style="border: black; width: 48%;">
<ul><li>Chapter Twelve: "Iron Maiden"</li>
<li>Chapter Thirteen: "Parting Ways"</li>
<li>Chapter Fourteen: "On the Road Again"</li>
<li>Chapter Fifteen: "The Shield Demon"</li>
<li>Chapter Sixteen: "Appointment Arrangements"</li>
<li>Chapter Seventeen: "The Princess's True Strength"</li>
<li>Chapter Eighteen: "Persuasion"</li>
<li>Chapter Nineteen: "The Tools"</li>
<li>Chapter Twenty: "Shadow"</li>
<li>Epilogue: "Name"</li>
<li>Extra Chapter: "Before I Met My Best Friend"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol22" style="text-align: center; font-weight: normal; background-color: transparent;">22</th><td> June 25, 2019<sup id="cite_ref-25" class="reference"><a href="#cite_note-25">[25]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-065839-1" title="Special:BookSources/978-4-04-065839-1">978-4-04-065839-1</a></td><td>December 21, 2021<sup id="cite_ref-US_Book_ISBN_4-21" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-133-0" title="Special:BookSources/978-1-64273-133-0">978-1-64273-133-0</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid #CCF;"><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left;"><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black;">
<ul><li>Prologue: "The Reticent Blacksmith"</li>
<li>Chapter One: "The Progenitor"</li>
<li>Chapter Two: "Claws and Hammer Power Up Method"</li>
<li>Chapter Three: "Holn's Research Weapon"</li>
<li>Chapter Four: "Quality Check of Heroes"</li>
<li>Chapter Five: "Gathering the Power of the Raph Species"</li>
<li>Chapter Six: "Selecting the Hammer Hero"</li>
<li>Chapter Seven: "Origin of the Past Heavenly Emperor"</li>
<li>Chapter Eight: "0 Territory"</li></ul></td>
<td style="border: black; width: 48%;">
<ul><li>Chapter Nine: "Heroic Presence"</li>
<li>Chapter Ten: "Exclusive Equipment"</li>
<li>Chapter Eleven: "Dinosaur Hunting"</li>
<li>Chapter Twelve: "Unexpected Visitors from Piensa"</li>
<li>Chapter Thirteen: "Beast Transformation of Rage"</li>
<li>Chapter Fourteen: "Known Vassal Weapon Heroes"</li>
<li>Chapter Fifteen: "In the Sanctuary"</li>
<li>Epilogue: "Dragon Slayer"</li></ul></td></tr></tbody></table>
</td></tr>
</tbody></table>`
)

var ParseWikipediaTableToVolumeInfoTestCases = map[string]ParseWikipediaTableToVolumeInfoTestCase{
	"a simple file with 6 columns and an unreleased volume with no announced date is handled correctly": {
		InputTableHtml:  mushokuTensieTable,
		InputNamePrefix: "Mushoku Tensei",
		ExpectedVolumeInfo: []wikipedia.VolumeInfo{
			{
				Name:        "Mushoku Tensei Vol. 1",
				ReleaseDate: getDatePointer(2019, 4, time.April),
			},
			{
				Name:        "Mushoku Tensei Vol. 2",
				ReleaseDate: getDatePointer(2019, 23, time.May),
			},
			{
				Name:        "Mushoku Tensei Vol. 26",
				ReleaseDate: nil,
			},
		},
	},
	"a simple table with 4 rows gets properly handled": {
		InputTableHtml:  asteriskWarTable,
		InputNamePrefix: "Asterisk War",
		ExpectedVolumeInfo: []wikipedia.VolumeInfo{
			{
				Name:        "Asterisk War Vol. 1",
				ReleaseDate: getDatePointer(2016, 30, time.August),
			},
			{
				Name:        "Asterisk War Vol. 2",
				ReleaseDate: getDatePointer(2016, 20, time.December),
			},
			{
				Name:        "Asterisk War Vol. 3",
				ReleaseDate: getDatePointer(2017, 18, time.April),
			},
			{
				Name:        "Asterisk War Vol. 4",
				ReleaseDate: getDatePointer(2017, 22, time.August),
			},
		},
	},
	"a table with 5 columns should be properly parsed": {
		InputTableHtml:  wrongWayToHealTable,
		InputNamePrefix: "The Wrong Way to Use Healing Magic",
		ExpectedVolumeInfo: []wikipedia.VolumeInfo{
			{
				Name:        "The Wrong Way to Use Healing Magic Vol. 1",
				ReleaseDate: getDatePointer(2022, 23, time.August),
			},
			{
				Name:        "The Wrong Way to Use Healing Magic Vol. 2",
				ReleaseDate: getDatePointer(2023, 15, time.May),
			},
			{
				Name:        "The Wrong Way to Use Healing Magic Vol. 3",
				ReleaseDate: getDatePointer(2023, 22, time.August),
			},
			{
				Name: "The Wrong Way to Use Healing Magic Vol. 4",
			},
			{
				Name: "The Wrong Way to Use Healing Magic Vol. 5",
			},
			{
				Name: "The Wrong Way to Use Healing Magic Vol. 6",
			},
			{
				Name: "The Wrong Way to Use Healing Magic Vol. 7",
			},
			{
				Name: "The Wrong Way to Use Healing Magic Vol. 8",
			},
			{
				Name: "The Wrong Way to Use Healing Magic Vol. 9",
			},
			{
				Name: "The Wrong Way to Use Healing Magic Vol. 10",
			},
			{
				Name: "The Wrong Way to Use Healing Magic Vol. 11",
			},
			{
				Name: "The Wrong Way to Use Healing Magic Vol. 12",
			},
		},
	},
	"a table with 5 columns and some rows with less full columns should be properly parsed": {
		InputTableHtml:  risingOfTheShieldHeroTable,
		InputNamePrefix: "The Rising of the Shield Hero",
		ExpectedVolumeInfo: []wikipedia.VolumeInfo{
			{
				Name:        "The Rising of the Shield Hero Vol. 1",
				ReleaseDate: getDatePointer(2015, 15, time.September),
			},
			{
				Name:        "The Rising of the Shield Hero Vol. 2",
				ReleaseDate: getDatePointer(2015, 20, time.October),
			},
			{
				Name:        "The Rising of the Shield Hero Vol. 3",
				ReleaseDate: getDatePointer(2016, 16, time.February),
			},
			{
				Name:        "The Rising of the Shield Hero Vol. 22",
				ReleaseDate: getDatePointer(2021, 21, time.December),
			},
		},
	},
}

func TestParseWikipediaTableToVolumeInfo(t *testing.T) {
	for name, args := range ParseWikipediaTableToVolumeInfoTestCases {
		t.Run(name, func(t *testing.T) {
			actualVolumeInfo := wikipedia.ParseWikipediaTableToVolumeInfo(args.InputNamePrefix, args.InputTableHtml)
			assert.Equal(t, len(args.ExpectedVolumeInfo), len(actualVolumeInfo))

			for i, volume := range args.ExpectedVolumeInfo {
				assert.Equal(t, volume.Name, actualVolumeInfo[i].Name)
				assert.Equal(t, volume.ReleaseDate, actualVolumeInfo[i].ReleaseDate)
			}
		})
	}
}

func getDatePointer(year, day int, month time.Month) *time.Time {
	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	return &date
}
