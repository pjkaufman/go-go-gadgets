//go:build unit

package wikipedia_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/wikipedia"
	"github.com/stretchr/testify/assert"
)

type GetNextTableAndItsEndPositionTestCase struct {
	InputHtml         string
	ExpectedTableHtml string
	ExpectedStopIndex int
}

const (
	theWrongWayToUseHealingMagicLightNovelTable = `<table class="wikitable" width="100%" style="">


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
	</th></tr><tr style="text-align: center;"><th scope="row" id="vol1" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">1</th><td> March 25, 2016<sup id="cite_ref-6" class="reference"><a href="#cite_note-6">[6]</a></sup></td><td><style data-mw-deduplicate="TemplateStyles:r1215172403">.mw-parser-output cite.citation{font-style:inherit;word-wrap:break-word}.mw-parser-output .citation q{quotes:"\"""\"""'""'"}.mw-parser-output .citation:target{background-color:rgba(0,127,255,0.133)}.mw-parser-output .id-lock-free.id-lock-free a{background:url("//upload.wikimedia.org/wikipedia/commons/6/65/Lock-green.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-free a{background-size:contain}.mw-parser-output .id-lock-limited.id-lock-limited a,.mw-parser-output .id-lock-registration.id-lock-registration a{background:url("//upload.wikimedia.org/wikipedia/commons/d/d6/Lock-gray-alt-2.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-limited a,body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-registration a{background-size:contain}.mw-parser-output .id-lock-subscription.id-lock-subscription a{background:url("//upload.wikimedia.org/wikipedia/commons/a/aa/Lock-red-alt-2.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-subscription a{background-size:contain}.mw-parser-output .cs1-ws-icon a{background:url("//upload.wikimedia.org/wikipedia/commons/4/4c/Wikisource-logo.svg")right 0.1em center/12px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .cs1-ws-icon a{background-size:contain}.mw-parser-output .cs1-code{color:inherit;background:inherit;border:none;padding:inherit}.mw-parser-output .cs1-hidden-error{display:none;color:#d33}.mw-parser-output .cs1-visible-error{color:#d33}.mw-parser-output .cs1-maint{display:none;color:#2C882D;margin-left:0.3em}.mw-parser-output .cs1-format{font-size:95%}.mw-parser-output .cs1-kern-left{padding-left:0.2em}.mw-parser-output .cs1-kern-right{padding-right:0.2em}.mw-parser-output .citation .mw-selflink{font-weight:inherit}html.skin-theme-clientpref-night .mw-parser-output .cs1-maint{color:#18911F}html.skin-theme-clientpref-night .mw-parser-output .cs1-visible-error,html.skin-theme-clientpref-night .mw-parser-output .cs1-hidden-error{color:#f8a397}@media(prefers-color-scheme:dark){html.skin-theme-clientpref-os .mw-parser-output .cs1-visible-error,html.skin-theme-clientpref-os .mw-parser-output .cs1-hidden-error{color:#f8a397}html.skin-theme-clientpref-os .mw-parser-output .cs1-maint{color:#18911F}}</style><style class="darkreader darkreader--sync" media="screen"></style><a href="/wiki/Special:BookSources/978-4-04-068185-6" title="Special:BookSources/978-4-04-068185-6">978-4-04-068185-6</a></td><td>August 23, 2022<sup id="cite_ref-7" class="reference"><a href="#cite_note-7">[7]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-200-9" title="Special:BookSources/978-1-64273-200-9">978-1-64273-200-9</a></td></tr>
	<tr style="text-align: center;"><th scope="row" id="vol2" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">2</th><td> June 24, 2016<sup id="cite_ref-8" class="reference"><a href="#cite_note-8">[8]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-068427-7" title="Special:BookSources/978-4-04-068427-7">978-4-04-068427-7</a></td><td>May 15, 2023<sup id="cite_ref-9" class="reference"><a href="#cite_note-9">[9]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-232-0" title="Special:BookSources/978-1-64273-232-0">978-1-64273-232-0</a></td></tr>
	<tr style="text-align: center;"><th scope="row" id="vol3" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">3</th><td> September 23, 2016<sup id="cite_ref-10" class="reference"><a href="#cite_note-10">[10]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-068636-3" title="Special:BookSources/978-4-04-068636-3">978-4-04-068636-3</a></td><td>August 22, 2023</td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-286-3" title="Special:BookSources/978-1-64273-286-3">978-1-64273-286-3</a></td></tr>
	<tr style="text-align: center;"><th scope="row" id="vol4" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">4</th><td> January 25, 2017<sup id="cite_ref-11" class="reference"><a href="#cite_note-11">[11]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-069054-4" title="Special:BookSources/978-4-04-069054-4">978-4-04-069054-4</a></td><td>—</td><td>—</td></tr>
	<tr style="text-align: center;"><th scope="row" id="vol5" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">5</th><td> April 25, 2017<sup id="cite_ref-12" class="reference"><a href="#cite_note-12">[12]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-069191-6" title="Special:BookSources/978-4-04-069191-6">978-4-04-069191-6</a></td><td>—</td><td>—</td></tr>
	<tr style="text-align: center;"><th scope="row" id="vol6" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">6</th><td> September 25, 2017<sup id="cite_ref-13" class="reference"><a href="#cite_note-13">[13]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-069498-6" title="Special:BookSources/978-4-04-069498-6">978-4-04-069498-6</a></td><td>—</td><td>—</td></tr>
	<tr style="text-align: center;"><th scope="row" id="vol7" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">7</th><td> February 24, 2018<sup id="cite_ref-14" class="reference"><a href="#cite_note-14">[14]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-069727-7" title="Special:BookSources/978-4-04-069727-7">978-4-04-069727-7</a></td><td>—</td><td>—</td></tr>
	<tr style="text-align: center;"><th scope="row" id="vol8" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">8</th><td> July 25, 2018<sup id="cite_ref-15" class="reference"><a href="#cite_note-15">[15]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-065023-4" title="Special:BookSources/978-4-04-065023-4">978-4-04-065023-4</a></td><td>—</td><td>—</td></tr>
	<tr style="text-align: center;"><th scope="row" id="vol9" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">9</th><td> November 24, 2018<sup id="cite_ref-16" class="reference"><a href="#cite_note-16">[16]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-065306-8" title="Special:BookSources/978-4-04-065306-8">978-4-04-065306-8</a></td><td>—</td><td>—</td></tr>
	<tr style="text-align: center;"><th scope="row" id="vol10" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">10</th><td> April 25, 2019<sup id="cite_ref-17" class="reference"><a href="#cite_note-17">[17]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-065682-3" title="Special:BookSources/978-4-04-065682-3">978-4-04-065682-3</a></td><td>—</td><td>—</td></tr>
	<tr style="text-align: center;"><th scope="row" id="vol11" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">11</th><td> October 25, 2019<sup id="cite_ref-18" class="reference"><a href="#cite_note-18">[18]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-064060-0" title="Special:BookSources/978-4-04-064060-0">978-4-04-064060-0</a></td><td>—</td><td>—</td></tr>
	<tr style="text-align: center;"><th scope="row" id="vol12" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">12</th><td> March 25, 2020<sup id="cite_ref-19" class="reference"><a href="#cite_note-19">[19]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-064538-4" title="Special:BookSources/978-4-04-064538-4">978-4-04-064538-4</a></td><td>—</td><td>—</td></tr>
	</tbody></table>`
	theWrongWayToUseHealingMagicLightNovelSection = `<h2><span class="mw-headline" id="Media">Media</span><span class="mw-editsection"><span class="mw-editsection-bracket">[</span><a href="/w/index.php?title=The_Wrong_Way_to_Use_Healing_Magic&amp;action=edit&amp;section=3" title="Edit section: Media"><span>edit</span></a><span class="mw-editsection-bracket">]</span></span></h2>
	<h3><span class="mw-headline" id="Light_novels">Light novels</span><span class="mw-editsection"><span class="mw-editsection-bracket">[</span><a href="/w/index.php?title=The_Wrong_Way_to_Use_Healing_Magic&amp;action=edit&amp;section=4" title="Edit section: Light novels"><span>edit</span></a><span class="mw-editsection-bracket">]</span></span></h3>
	<p>The series written by Kurokata began serialization online in March 2014 on the user-generated novel publishing website <a href="/wiki/Sh%C5%8Dsetsuka_ni_Nar%C5%8D" title="Shōsetsuka ni Narō">Shōsetsuka ni Narō</a>. It was later acquired by <a href="/wiki/Media_Factory" title="Media Factory">Media Factory</a>, who have published twelve volumes with illustrations by KeG between March 25, 2016 and March 25, 2020 under their MF Books imprint. The light novel is licensed in North America by One Peace Books.<sup id="cite_ref-One-Peace_5-0" class="reference"><a href="#cite_note-One-Peace-5">[5]</a></sup>
</p>
` + theWrongWayToUseHealingMagicLightNovelTable + `
<p>A sequel light novel series by the same author and illustrator, titled <i>The Wrong Way to Use Healing Magic Returns</i>, began publication on December 25, 2023.
</p>
<table class="wikitable" width="100%" style="">


<tbody><tr style="border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom="">
<th width="4%"><abbr title="Number">No.</abbr>
</th>
<th width="48%">Japanese release date
</th>
<th width="48%">Japanese ISBN
</th></tr><tr style="text-align: center;"><th scope="row" id="vol1" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">1</th><td> December 25, 2023<sup id="cite_ref-20" class="reference"><a href="#cite_note-20">[20]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-683145-3" title="Special:BookSources/978-4-04-683145-3">978-4-04-683145-3</a></td></tr>
<tr style="text-align: center;"><th scope="row" id="vol2" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">2</th><td> March 25, 2024<sup id="cite_ref-21" class="reference"><a href="#cite_note-21">[21]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-683481-2" title="Special:BookSources/978-4-04-683481-2">978-4-04-683481-2</a></td></tr>
</tbody></table>`
	theRisingOfTheShieldHereLightNovelTable = `<table class="wikitable" width="98%" style="">


<tbody><tr style="border-bottom: 3px solid rgb(31, 58, 173); --darkreader-inline-border-bottom: #1d37a4;" data-darkreader-inline-border-bottom="">
<th width="4%"><abbr title="Number">No.</abbr>
</th>
<th width="24%">Original release date
</th>
<th width="24%">Original ISBN
</th>
<th width="24%">English release date
</th>
<th width="24%">English ISBN
</th></tr><tr style="text-align: center;"><th scope="row" id="vol1" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">1</th><td> August 22, 2013<sup id="cite_ref-3" class="reference"><a href="#cite_note-3">[3]</a></sup></td><td><style data-mw-deduplicate="TemplateStyles:r1215172403">.mw-parser-output cite.citation{font-style:inherit;word-wrap:break-word}.mw-parser-output .citation q{quotes:"\"""\"""'""'"}.mw-parser-output .citation:target{background-color:rgba(0,127,255,0.133)}.mw-parser-output .id-lock-free.id-lock-free a{background:url("//upload.wikimedia.org/wikipedia/commons/6/65/Lock-green.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-free a{background-size:contain}.mw-parser-output .id-lock-limited.id-lock-limited a,.mw-parser-output .id-lock-registration.id-lock-registration a{background:url("//upload.wikimedia.org/wikipedia/commons/d/d6/Lock-gray-alt-2.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-limited a,body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-registration a{background-size:contain}.mw-parser-output .id-lock-subscription.id-lock-subscription a{background:url("//upload.wikimedia.org/wikipedia/commons/a/aa/Lock-red-alt-2.svg")right 0.1em center/9px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .id-lock-subscription a{background-size:contain}.mw-parser-output .cs1-ws-icon a{background:url("//upload.wikimedia.org/wikipedia/commons/4/4c/Wikisource-logo.svg")right 0.1em center/12px no-repeat}body:not(.skin-timeless):not(.skin-minerva) .mw-parser-output .cs1-ws-icon a{background-size:contain}.mw-parser-output .cs1-code{color:inherit;background:inherit;border:none;padding:inherit}.mw-parser-output .cs1-hidden-error{display:none;color:#d33}.mw-parser-output .cs1-visible-error{color:#d33}.mw-parser-output .cs1-maint{display:none;color:#2C882D;margin-left:0.3em}.mw-parser-output .cs1-format{font-size:95%}.mw-parser-output .cs1-kern-left{padding-left:0.2em}.mw-parser-output .cs1-kern-right{padding-right:0.2em}.mw-parser-output .citation .mw-selflink{font-weight:inherit}html.skin-theme-clientpref-night .mw-parser-output .cs1-maint{color:#18911F}html.skin-theme-clientpref-night .mw-parser-output .cs1-visible-error,html.skin-theme-clientpref-night .mw-parser-output .cs1-hidden-error{color:#f8a397}@media(prefers-color-scheme:dark){html.skin-theme-clientpref-os .mw-parser-output .cs1-visible-error,html.skin-theme-clientpref-os .mw-parser-output .cs1-hidden-error{color:#f8a397}html.skin-theme-clientpref-os .mw-parser-output .cs1-maint{color:#18911F}}</style><style class="darkreader darkreader--sync" media="screen"></style><a href="/wiki/Special:BookSources/978-4-8401-5275-4" title="Special:BookSources/978-4-8401-5275-4">978-4-8401-5275-4</a></td><td>September 15, 2015<sup id="cite_ref-US_Book_ISBN_4-0" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-935548-72-0" title="Special:BookSources/978-1-935548-72-0">978-1-935548-72-0</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
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
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
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
<tr style="text-align: center;"><th scope="row" id="vol2" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">2</th><td> October 1, 2013<sup id="cite_ref-5" class="reference"><a href="#cite_note-5">[5]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-066049-3" title="Special:BookSources/978-4-04-066049-3">978-4-04-066049-3</a></td><td>October 20, 2015<sup id="cite_ref-US_Book_ISBN_4-1" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-935548-78-2" title="Special:BookSources/978-1-935548-78-2">978-1-935548-78-2</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
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
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
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
<tr style="text-align: center;"><th scope="row" id="vol3" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">3</th><td> December 21, 2013<sup id="cite_ref-6" class="reference"><a href="#cite_note-6">[6]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-066166-7" title="Special:BookSources/978-4-04-066166-7">978-4-04-066166-7</a></td><td>February 16, 2016<sup id="cite_ref-US_Book_ISBN_4-2" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-935548-66-9" title="Special:BookSources/978-1-935548-66-9">978-1-935548-66-9</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
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
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
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
<tr style="text-align: center;"><th scope="row" id="vol4" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">4</th><td> February 25, 2014<sup id="cite_ref-7" class="reference"><a href="#cite_note-7">[7]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-066321-0" title="Special:BookSources/978-4-04-066321-0">978-4-04-066321-0</a></td><td>June 14, 2016<sup id="cite_ref-US_Book_ISBN_4-3" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-935548-65-2" title="Special:BookSources/978-1-935548-65-2">978-1-935548-65-2</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "On The Run"</li>
<li>Chapter One: "Demi-Human Adventurer Town"</li>
<li>Chapter Two: "Noblemen"</li>
<li>Chapter Three: "Tyrant Dragon Rex"</li>
<li>Chapter Four: "The Legendary Bird God"</li>
<li>Chapter Five: "Filo vs Fitoria"</li>
<li>Chapter Six: "The Bird God's Peace"</li>
<li>Chapter Seven: "The Battle of Shield and Spear"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Eight: "Judgment"</li>
<li>Chapter Nine: "Replica"</li>
<li>Chapter Ten: "Shield of Wrath"</li>
<li>Chapter Eleven: "The Queen"</li>
<li>Chapter Twelve: "Paying the Piper"</li>
<li>Epilogue: "Friends Forever"</li>
<li>Extra Bonus Chapter: "The Fearful Filolial"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol5" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">5</th><td> April 25, 2014<sup id="cite_ref-8" class="reference"><a href="#cite_note-8">[8]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-066718-8" title="Special:BookSources/978-4-04-066718-8">978-4-04-066718-8</a></td><td>August 23, 2016<sup id="cite_ref-US_Book_ISBN_4-4" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-935548-67-6" title="Special:BookSources/978-1-935548-67-6">978-1-935548-67-6</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "Class Up"</li>
<li>Chapter One: "The Heroes' Teammates"</li>
<li>Chapter Two: "Meeting of the Heroes"</li>
<li>Chapter Three: "Power Up"</li>
<li>Chapter Four: "Weapon Copy"</li>
<li>Chapter Five: "Gravestones"</li>
<li>Chapter Six: "Cal Mira"</li>
<li>Chapter Seven: "The Tavern"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Eight: "Karma"</li>
<li>Chapter Nine: "Island Days"</li>
<li>Chapter Ten: "The Water Temple"</li>
<li>Chapter Eleven: "Inter-dimensional Whale"</li>
<li>Chapter Twelve: "L'Arc Berg"</li>
<li>Chapter Thirteen: "Soul-Healing Water"</li>
<li>Epilogue: "The Problem We Face"</li>
<li>Extra Chapter: "The Cal Mira Hot Springs"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol6" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">6</th><td> June 25, 2014<sup id="cite_ref-9" class="reference"><a href="#cite_note-9">[9]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-066790-4" title="Special:BookSources/978-4-04-066790-4">978-4-04-066790-4</a></td><td>November 22, 2016<sup id="cite_ref-US_Book_ISBN_4-5" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-935548-56-0" title="Special:BookSources/978-1-935548-56-0">978-1-935548-56-0</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "Cal Mira Superstitions"</li>
<li>Chapter One: "The Seven Star Heroes"</li>
<li>Chapter Two: "An Unhappy Girl"</li>
<li>Chapter Three: "Framed Again"</li>
<li>Chapter Four: "Custom Order"</li>
<li>Chapter Five: "Battle Advisors"</li>
<li>Chapter Six: "Hengen Muso Style"</li>
<li>Chapter Seven: "Impossible Training?"</li>
<li>Chapter Eight: "Life-Force Water"</li>
<li>Chapter Nine: "What it Means to Train"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Ten: "Kigurumi"</li>
<li>Chapter Eleven: "——'s Familiar"</li>
<li>Chapter Twelve: "Getting Ahead of the Enemy"</li>
<li>Chapter Thirteen: "Game Knowledge Bares its Fangs"</li>
<li>Chapter Fourteen: "What it Means to be a Hero"</li>
<li>Chapter Fifteen: "The Spirit Tortoise"</li>
<li>Chapter Sixteen: "The Country Above the Spirit Tortoise"</li>
<li>Epilogue: "A Disquieting Place"</li>
<li>Extra Chapter: "Trials and Tribulations of the Bow Hero"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol7" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">7</th><td> September 25, 2014<sup id="cite_ref-10" class="reference"><a href="#cite_note-10">[10]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-066996-0" title="Special:BookSources/978-4-04-066996-0">978-4-04-066996-0</a></td><td>April 18, 2017<sup id="cite_ref-US_Book_ISBN_4-6" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-944937-08-9" title="Special:BookSources/978-1-944937-08-9">978-1-944937-08-9</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "The Search"</li>
<li>Chapter One: "Helping Others"</li>
<li>Chapter Two: "Spirit Tortoise Familiar (Human Type)"</li>
<li>Chapter Three: "The Spirit Tortoise Reawakens"</li>
<li>Chapter Four: "Spirit Tortoise Tyrant"</li>
<li>Chapter Five: "Mass Destruction"</li>
<li>Chapter Six: "Versus the Spirit Tortoise, Opening Stages"</li>
<li>Chapter Seven: "Buying Time"</li>
<li>Chapter Eight: "The Search"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Nine: "The Spirit Tortoise Cave"</li>
<li>Chapter Ten: "Strangers"</li>
<li>Chapter Eleven: "The Heroes' Inscription"</li>
<li>Chapter Twelve: "The Spirit Tortoise's Heart"</li>
<li>Chapter Thirteen: "Who Pulls the Strings"</li>
<li>Chapter Fourteen: "Liberation"</li>
<li>Chapter Fifteen: "The Spirit Tortoise's Soul"</li>
<li>Epilogue: "Ost Horai"</li>
<li>Extra Chapter: "Searching for the Soul-Healing Water"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol8" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">8</th><td> November 21, 2014<sup id="cite_ref-11" class="reference"><a href="#cite_note-11">[11]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-067180-2" title="Special:BookSources/978-4-04-067180-2">978-4-04-067180-2</a></td><td>June 13, 2017<sup id="cite_ref-US_Book_ISBN_4-7" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-944937-09-6" title="Special:BookSources/978-1-944937-09-6">978-1-944937-09-6</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "The Never-Ending Labyrinth"</li>
<li>Chapter One: "The Hunting Hero"</li>
<li>Chapter Two: "Escape"</li>
<li>Chapter Three: "The Unknown World"</li>
<li>Chapter Four: "Selling Drop Items"</li>
<li>Chapter Five: "Sales Demonstration"</li>
<li>Chapter Six: "Otherworldly Equipment"</li>
<li>Chapter Seven: "The Legend of the Waves"</li>
<li>Chapter Eight: "On the Way to the Hunting Hero's House"</li>
<li>Chapter Nine: "Shikigami"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Ten: "The Katana of the Vassal Weapons:</li>
<li>Chapter Eleven: "Rescuing the Angel"</li>
<li>Chapter Twelve: "Humming Fairy"</li>
<li>Chapter Thirteen: "The Hunting Hero's Skills"</li>
<li>Chapter Fourteen: "Return Dragon Vien"</li>
<li>Chapter Fifteen: "The Katana's Choice"</li>
<li>Chapter Sixteen: "No Incantations"</li>
<li>Chapter Seventeen: "Blood Flower Strike"</li>
<li>Epilogue: "Together Again"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol9" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">9</th><td> January 23, 2015<sup id="cite_ref-12" class="reference"><a href="#cite_note-12">[12]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-067355-4" title="Special:BookSources/978-4-04-067355-4">978-4-04-067355-4</a></td><td>November 15, 2017<sup id="cite_ref-US_Book_ISBN_4-8" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-944937-25-6" title="Special:BookSources/978-1-944937-25-6">978-1-944937-25-6</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "The Waves of Another World"</li>
<li>Chapter One: "Otherworldly Techniques"</li>
<li>Chapter Two: "Quick Draw"</li>
<li>Chapter Three: "Lure"</li>
<li>Chapter Four: "Like a Charging Wild Boar"</li>
<li>Chapter Five: Together, With Conditions"</li>
<li>Chapter Six: "The Reformed"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Seven: "Barbaroi Armor"</li>
<li>Chapter Eight: "Two Swords"</li>
<li>Chapter Nine: "Kyo's Laboratory"</li>
<li>Chapter Ten: "When Trust is Lost"</li>
<li>Chapter Eleven: "Sacrifice Aura"</li>
<li>Chapter Twelve: "A Heavy Price to Pay"</li>
<li>Epilogue: "Kizuna Between Worlds"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol10" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">10</th><td> March 25, 2015<sup id="cite_ref-13" class="reference"><a href="#cite_note-13">[13]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-067485-8" title="Special:BookSources/978-4-04-067485-8">978-4-04-067485-8</a></td><td>March 20, 2018<sup id="cite_ref-US_Book_ISBN_4-9" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-944937-26-3" title="Special:BookSources/978-1-944937-26-3">978-1-944937-26-3</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "The Spirit Tortoise's Barrier"</li>
<li>Chapter One: "The Seven Star Staff Hero"</li>
<li>Chapter Two: "Whereabouts of the Slaves"</li>
<li>Chapter Three: "Acquaintances"</li>
<li>Chapter Four: "E Float Shield"</li>
<li>Chapter Five: "The Seaetto Territory"</li>
<li>Chapter Six: "Feeding the Herd"</li>
<li>Chapter Seven: "Employing the Bioplant"</li>
<li>Chapter Eight: "Children of the Sea"</li>
<li>Chapter Nine: "Hanging Out the Shield"</li>
<li>Chapter Ten: "Zeltoble"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Eleven: "Slave Hunters"</li>
<li>Chapter Twelve: "The Department Store"</li>
<li>Chapter Thirteen: "The Underground Coliseum"</li>
<li>Chapter Fourteen: "Ring Name"</li>
<li>Chapter Fifteen: "Surprise Attacks and Conspiracies"</li>
<li>Chapter Sixteen: "Nadia"</li>
<li>Chapter Seventeen: "Farce"</li>
<li>Chapter Eighteen: "Exhibition Match"</li>
<li>Chapter Nineteen: "Big Shots of the Underground"</li>
<li>Epilogue: "Come-On"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol11" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">11</th><td> June 25, 2015<sup id="cite_ref-14" class="reference"><a href="#cite_note-14">[14]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-067698-2" title="Special:BookSources/978-4-04-067698-2">978-4-04-067698-2</a></td><td>June 12, 2018<sup id="cite_ref-US_Book_ISBN_4-10" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-944937-46-1" title="Special:BookSources/978-1-944937-46-1">978-1-944937-46-1</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "To Market"</li>
<li>Chapter One: "Sacred Tree Elixir"</li>
<li>Chapter Two: "Return of the Village"</li>
<li>Chapter Three: "Alps"</li>
<li>Chapter Four: "A Shield to Protect the Shield"</li>
<li>Chapter Five: "Trash and the Hakuko"</li>
<li>Chapter Six: "The Fruits of Training"</li>
<li>Chapter Seven: "The Plan to Capture the Spear Hero"</li>
<li>Chapter Eight: "The Day the Game Ended"</li>
<li>Chapter Nine: "I Dub Thee Witch"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Ten: "New Awakening"</li>
<li>Chapter Eleven: "Loincloth Pup"</li>
<li>Chapter Twelve: "The Decision"</li>
<li>Chapter Thirteen: "Oodles of Ambushes"</li>
<li>Chapter Fourteen: "Official Request"</li>
<li>Chapter Fifteen: "The Masked Man"</li>
<li>Chapter Sixteen: "The Merits of Invading Other Worlds"</li>
<li>Chapter Seventeen: "Temptation"</li>
<li>Chapter Eighteen: "Flash"</li>
<li>Epilogue: "Making Peace with the Sword Hero"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol12" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">12</th><td> September 25, 2015<sup id="cite_ref-15" class="reference"><a href="#cite_note-15">[15]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-067787-3" title="Special:BookSources/978-4-04-067787-3">978-4-04-067787-3</a></td><td>August 18, 2018<sup id="cite_ref-US_Book_ISBN_4-11" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-944937-95-9" title="Special:BookSources/978-1-944937-95-9">978-1-944937-95-9</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "The Shield Hero's Morning"</li>
<li>Chapter One: "Instant Awakening"</li>
<li>Chapter Two: "The Alchemist"</li>
<li>Chapter Three: "Filolials and Dragons"</li>
<li>Chapter Four: "Stardust Blade"</li>
<li>Chapter Five: "Knock and Run"</li>
<li>Chapter Six: "Level Drain"</li>
<li>Chapter Seven: "Plagued Earth"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Eight: "Demon Dragon"</li>
<li>Chapter Nine: "Forced Power-Up"</li>
<li>Chapter Ten: "Purification"</li>
<li>Chapter Eleven: "Perfect Hidden Justice"</li>
<li>Chapter Twelve: "Justice vs. Justice"</li>
<li>Chapter Thirteen: "Atonement"</li>
<li>Chapter Fourteen: "Secret Base"</li>
<li>Chapter Fifteen: "Form is Emptiness"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol13" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">13</th><td> November 25, 2015<sup id="cite_ref-16" class="reference"><a href="#cite_note-16">[16]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-067965-5" title="Special:BookSources/978-4-04-067965-5">978-4-04-067965-5</a></td><td>December 18, 2018<sup id="cite_ref-US_Book_ISBN_4-12" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-944937-96-6" title="Special:BookSources/978-1-944937-96-6">978-1-944937-96-6</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "Team Assignments"</li>
<li>Chapter One: "Advance Payment"</li>
<li>Chapter Two: "Sending Word of Our Visit"</li>
<li>Chapter Three: "Arrival in Siltvelt"</li>
<li>Chapter Four: "Shield of the Beast King"</li>
<li>Chapter Five: "Harem"</li>
<li>Chapter Six: "Conspiracy"</li>
<li>Chapter Seven: "A True Siltveltian"</li>
<li>Chapter Eight: "Honor in Battle"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Nine: "Beast Transformation"</li>
<li>Chapter Ten: "Assigning the Heroes"</li>
<li>Chapter Eleven: "The Flawed Master"</li>
<li>Chapter Twelve: "A Little Help from the Water Dragon"</li>
<li>Chapter Thirteen: "Q'ten Lo Revolutionaries"</li>
<li>Chapter Fourteen: "Sakura Stone of Destiny"</li>
<li>Chapter Fifteen: "Sakura Stone of Influence"</li>
<li>Epilogue: "The Old Guy's Master"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol14" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">14</th><td> February 25, 2016<sup id="cite_ref-17" class="reference"><a href="#cite_note-17">[17]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-068121-4" title="Special:BookSources/978-4-04-068121-4">978-4-04-068121-4</a></td><td>October 15, 2019<sup id="cite_ref-US_Book_ISBN_4-13" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-018-0" title="Special:BookSources/978-1-64273-018-0">978-1-64273-018-0</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "Planning the Q'Ten Lo Invasion"</li>
<li>Chapter One: "The Sealed Orochi"</li>
<li>Chapter Two: "Sharing Power-Ups"</li>
<li>Chapter Three: "The Cursed Ama-no-Murakumo Sword"</li>
<li>Chapter Four: "Tailwind"</li>
<li>Chapter Five: "Information on the Enemy"</li>
<li>Chapter Six: "Use of Life Force"</li>
<li>Chapter Seven: "A Terrible Sense of Direction"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Eight: "Big Sister"</li>
<li>Chapter Nine: "The Miko Priestess of Carnage"</li>
<li>Chapter Ten: "Shield Power-Up Method"</li>
<li>Chapter Eleven: "A Brief Return Trip"</li>
<li>Chapter Twelve: "Past and Present"</li>
<li>Chapter Thirteen: "The Past Heavenly Emperor"</li>
<li>Chapter Fourteen: "The True Terror of Monsters"</li>
<li>Epilogue: "Dusk"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol15" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">15</th><td> September 23, 2016<sup id="cite_ref-18" class="reference"><a href="#cite_note-18">[18]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-068638-7" title="Special:BookSources/978-4-04-068638-7">978-4-04-068638-7</a></td><td>December 15, 2019<sup id="cite_ref-US_Book_ISBN_4-14" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-019-7" title="Special:BookSources/978-1-64273-019-7">978-1-64273-019-7</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "A Problem with Bandits"</li>
<li>Chapter One: "Birth of the Raph Species"</li>
<li>Chapter Two: "Territory Reform"</li>
<li>Chapter Three: "Spirit Tortoise Shell"</li>
<li>Chapter Four: "Fitoria's Request"</li>
<li>Chapter Five: "The Street Racer"</li>
<li>Chapter Six: "The Love Hunter"</li>
<li>Chapter Seven: "Filolial Terror"</li>
<li>Chapter Eight: "The Third Hero Conference"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Nine: "Siblings' Squabble"</li>
<li>Chapter Ten: "Home of the Phoenix"</li>
<li>Chapter Eleven: "The Lost Hero's Diary"</li>
<li>Chapter Twelve: "The Final Seven Star Weapon"</li>
<li>Chapter Thirteen: "The Night Before the Phoenix Battle"</li>
<li>Chapter Fourteen: "Fighting the Phoenix"</li>
<li>Chapter Fifteen: "A Forbidden Flicker"</li>
<li>Epilogue: "The Girl Who Became a Shield"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol16" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">16</th><td> January 25, 2017<sup id="cite_ref-19" class="reference"><a href="#cite_note-19">[19]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-069051-3" title="Special:BookSources/978-4-04-069051-3">978-4-04-069051-3</a></td><td>March 15, 2020<sup id="cite_ref-US_Book_ISBN_4-15" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-020-3" title="Special:BookSources/978-1-64273-020-3">978-1-64273-020-3</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "The Funeral"</li>
<li>Chapter One: "The Ocean Floor"</li>
<li>Chapter Two: "The Festival"</li>
<li>Chapter Three: "The Genius"</li>
<li>Chapter Four: "Stolen Power"</li>
<li>Chapter Five: "The Spirits"</li>
<li>Chapter Six: "The Staff Hero"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Seven: "The Wisest King of Wisdom"</li>
<li>Chapter Eight: "X"</li>
<li>Chapter Nine: "Fenrir Force"</li>
<li>Chapter Ten: "Two Regular Guys and the Strongest Seven Star Hero"</li>
<li>Chapter Eleven: "The Shield Hero Now Orders You"</li>
<li>Chapter Twelve: "The Execution"</li>
<li>Epilogue: "Vanguard of the Waves"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol17" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">17</th><td> March 25, 2017<sup id="cite_ref-20" class="reference"><a href="#cite_note-20">[20]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-069190-9" title="Special:BookSources/978-4-04-069190-9">978-4-04-069190-9</a></td><td>July 14, 2020<sup id="cite_ref-US_Book_ISBN_4-16" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-053-1" title="Special:BookSources/978-1-64273-053-1">978-1-64273-053-1</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "The Coronation"</li>
<li>Chapter One: "Talk of Love"</li>
<li>Chapter Two: "Limit break"</li>
<li>Chapter Three: "Party Selection"</li>
<li>Chapter Four: "Arrival into Conflict"</li>
<li>Chapter Five: "Inter-World Adaptation"</li>
<li>Chapter Six: "Hidden Abilities"</li>
<li>Chapter Seven: "Finding Kizuna"</li>
<li>Chapter Eight: "Subterranean Maze City"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Nine: "Outsider Theory"</li>
<li>Chapter Ten: "A Familiar Face"</li>
<li>Chapter Eleven: "The Mirror Vassal Weapon"</li>
<li>Chapter Twelve: "A Falling Out"</li>
<li>Chapter Thirteen: "Forced Possession"</li>
<li>Chapter Fourteen: "Quick Adaptation"</li>
<li>Chapter Fifteen: "Mirror"</li>
<li>Epilogue: "A Responsibility to Justice"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol18" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">18</th><td> July 25, 2017<sup id="cite_ref-21" class="reference"><a href="#cite_note-21">[21]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-069354-5" title="Special:BookSources/978-4-04-069354-5">978-4-04-069354-5</a></td><td>November 12, 2020<sup id="cite_ref-US_Book_ISBN_4-17" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-082-1" title="Special:BookSources/978-1-64273-082-1">978-1-64273-082-1</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "The Meeting to Discuss Efficient Eating Enhancement"</li>
<li>Chapter One: "Sloth"</li>
<li>Chapter Two: "Library Search"</li>
<li>Chapter Three: "Fishing Fool's Determination"</li>
<li>Chapter Four: "Sisters and Jealousy"</li>
<li>Chapter Five: "Ultimate Soup Stock"</li>
<li>Chapter Six: "Seya's Restaurant"</li>
<li>Chapter Seven: "Contentious Cooking Battle"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Eight: "Medicinal Cooking"</li>
<li>Chapter Nine: "Resolution via Violence"</li>
<li>Chapter Ten: "Dragon of Ultimate Magic"</li>
<li>Chapter Eleven: "Volunteer Soldiers"</li>
<li>Chapter Twelve: "Double Reflection"</li>
<li>Chapter Thirteen: "The Reborn"</li>
<li>Epilogue: "The Game Knowledge Pitfall"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol19" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">19</th><td> January 25, 2018<sup id="cite_ref-22" class="reference"><a href="#cite_note-22">[22]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-069665-2" title="Special:BookSources/978-4-04-069665-2">978-4-04-069665-2</a></td><td>April 27, 2021<sup id="cite_ref-US_Book_ISBN_4-18" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-104-0" title="Special:BookSources/978-1-64273-104-0">978-1-64273-104-0</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "An Exchange of Otherworldly  Information"</li>
<li>Chapter One: "A Visit to the Head Temple"</li>
<li>Chapter Two: "Holy Tool Grotto"</li>
<li>Chapter Three: "The Work of Men and Monsters"</li>
<li>Chapter Four: "A Visit to Demon Dragon Castle"</li>
<li>Chapter Five: "Appraisal Camouflage"</li>
<li>Chapter Six: "A New Heavenly King"</li>
<li>Chapter Seven: "The Demon Dragon's Treasure"</li>
<li>Chapter Eight: "An Alluring Pudding"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Nine: "Just to Make Sure"</li>
<li>Chapter Ten: "The importance of Anger"</li>
<li>Chapter Eleven: "Opposing Nullification"</li>
<li>Chapter Twelve: "Intelligence Operative"</li>
<li>Chapter Thirteen: "Insensitive Individuals"</li>
<li>Chapter Fourteen: "Megido Iron Maiden"</li>
<li>Chapter Fifteen: "Defense of the Port"</li>
<li>Epilogue: "A Visitor Late at Night"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol20" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">20</th><td> December 25, 2018<sup id="cite_ref-23" class="reference"><a href="#cite_note-23">[23]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-065134-7" title="Special:BookSources/978-4-04-065134-7">978-4-04-065134-7</a></td><td>June 22, 2021<sup id="cite_ref-US_Book_ISBN_4-19" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-105-7" title="Special:BookSources/978-1-64273-105-7">978-1-64273-105-7</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "Believing Sloth Will Save the World"</li>
<li>Chapter One: "Prisoner Transport"</li>
<li>Chapter Two: "Training for an Obstinate Man"</li>
<li>Chapter Three: "The Sword Hero’s Sense of Responsibility"</li>
<li>Chapter Four: "The Filolial Ruins"</li>
<li>Chapter Five: "Village Abnormality"</li>
<li>Chapter Six: "Encounter with Extinct Monsters"</li>
<li>Chapter Seven: "Double the Shield Heroes"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Eight: "Hero Worship"</li>
<li>Chapter Nine: "Ancient Siltran"</li>
<li>Chapter Ten: "The Evil Researcher"</li>
<li>Chapter Eleven: "Bread Trees and Bread Troubles"</li>
<li>Chapter Twelve: "A Determination to War"</li>
<li>Chapter Thirteen: "Online Trolling"</li>
<li>Epilogue: "Different Constellations"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol21" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">21</th><td> February 25, 2019<sup id="cite_ref-24" class="reference"><a href="#cite_note-24">[24]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-065546-8" title="Special:BookSources/978-4-04-065546-8">978-4-04-065546-8</a></td><td>October 26, 2021<sup id="cite_ref-US_Book_ISBN_4-20" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-132-3" title="Special:BookSources/978-1-64273-132-3">978-1-64273-132-3</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "Use of Floating Weapons"</li>
<li>Chapter One: "The Siltran Situation"</li>
<li>Chapter Two: "Wagon Travel with Keel and the Gang"</li>
<li>Chapter Three: "Wave Trauma"</li>
<li>Chapter Four: "Thanks to the Assassin"</li>
<li>Chapter Five: "Genetic Modification"</li>
<li>Chapter Six: "Mikey"</li>
<li>Chapter Seven: "The Raph Species Upgrade Plan"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Eight: "The Troubles of the Shield Hero"</li>
<li>Chapter Nine: "Confusion Target"</li>
<li>Chapter Ten: "Filolia"</li>
<li>Chapter Eleven: "The Imitations (Improved)"</li>
<li>Chapter Twelve: "The Origin of the Waves"</li>
<li>Chapter Thirteen: "How to Kill a God"</li>
<li>Epilogue: "The Fear of Those Who are Eternal"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol22" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">22</th><td> June 25, 2019<sup id="cite_ref-25" class="reference"><a href="#cite_note-25">[25]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-065839-1" title="Special:BookSources/978-4-04-065839-1">978-4-04-065839-1</a></td><td>December 21, 2021<sup id="cite_ref-US_Book_ISBN_4-21" class="reference"><a href="#cite_note-US_Book_ISBN-4">[4]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-133-0" title="Special:BookSources/978-1-64273-133-0">978-1-64273-133-0</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "The Reticent Blacksmith"</li>
<li>Chapter One: "The Progenitor"</li>
<li>Chapter Two: "Claws and Hammer Power Up Method"</li>
<li>Chapter Three: "Holn's Research Weapon"</li>
<li>Chapter Four: "Quality Check of Heroes"</li>
<li>Chapter Five: "Gathering the Power of the Raph Species"</li>
<li>Chapter Six: "Selecting the Hammer Hero"</li>
<li>Chapter Seven: "Origin of the Past Heavenly Emperor"</li>
<li>Chapter Eight: "0 Territory"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
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
	theRisingOfTheShieldHereLightNovelSection = `<h2><span class="mw-headline" id="Light_novel">Light novel</span><span class="mw-editsection"><span class="mw-editsection-bracket">[</span><a href="/w/index.php?title=List_of_The_Rising_of_the_Shield_Hero_volumes&amp;action=edit&amp;section=1" title="Edit section: Light novel"><span>edit</span></a><span class="mw-editsection-bracket">]</span></span></h2>
<h3><span class="mw-headline" id="The_Rising_of_the_Shield_Hero"><i>The Rising of the Shield Hero</i></span><span class="mw-editsection"><span class="mw-editsection-bracket">[</span><a href="/w/index.php?title=List_of_The_Rising_of_the_Shield_Hero_volumes&amp;action=edit&amp;section=2" title="Edit section: The Rising of the Shield Hero"><span>edit</span></a><span class="mw-editsection-bracket">]</span></span></h3>
` + theRisingOfTheShieldHereLightNovelTable + `
<h3><span class="mw-headline" id="Limited_Edition_The_Rising_of_the_Shield_Hero_Season_1_Light_Novel"><i>Limited Edition The Rising of the Shield Hero Season 1 Light Novel</i></span><span class="mw-editsection"><span class="mw-editsection-bracket">[</span><a href="/w/index.php?title=List_of_The_Rising_of_the_Shield_Hero_volumes&amp;action=edit&amp;section=3" title="Edit section: Limited Edition The Rising of the Shield Hero Season 1 Light Novel"><span>edit</span></a><span class="mw-editsection-bracket">]</span></span></h3>
<p>Originally released as 4 separate volumes with the Japanese Special Editions of the Season 1 anime.
</p>
<table class="wikitable" width="98%" style="">


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
</th></tr><tr style="text-align: center;"><th scope="row" id="vol" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""></th><td> April 24, 2019;<sup id="cite_ref-26" class="reference"><a href="#cite_note-26">[26]</a></sup><br> May 24, 2019;<sup id="cite_ref-27" class="reference"><a href="#cite_note-27">[27]</a></sup><br> June 26, 2019;<sup id="cite_ref-28" class="reference"><a href="#cite_note-28">[28]</a></sup><br> July 24, 2019<sup id="cite_ref-29" class="reference"><a href="#cite_note-29">[29]</a></sup></td><td>—</td><td>May 26, 2020<sup id="cite_ref-30" class="reference"><a href="#cite_note-30">[30]</a></sup></td><td>—</td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Raphtalia's Stay in Lute Village</li>
<li>Filo's Gourmet Trading Journey</li>
<li>Melty's Survival Dinner</li>
<li>Naofumi's Cal Mira Sightseeing</li></ul></td></tr></tbody></table>
</td></tr>
</tbody></table>
<h3><span class="mw-headline" id="The_Reprise_of_the_Spear_Hero"><i>The Reprise of the Spear Hero</i></span><span class="mw-editsection"><span class="mw-editsection-bracket">[</span><a href="/w/index.php?title=List_of_The_Rising_of_the_Shield_Hero_volumes&amp;action=edit&amp;section=4" title="Edit section: The Reprise of the Spear Hero"><span>edit</span></a><span class="mw-editsection-bracket">]</span></span></h3>
<table class="wikitable" width="98%" style="">


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
</th></tr><tr style="text-align: center;"><th scope="row" id="vol1" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">1</th><td> September 25, 2017<sup id="cite_ref-31" class="reference"><a href="#cite_note-31">[31]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-069502-0" title="Special:BookSources/978-4-04-069502-0">978-4-04-069502-0</a></td><td>October 16, 2018<sup id="cite_ref-US_Spear_LN_32-0" class="reference"><a href="#cite_note-US_Spear_LN-32">[32]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-003-6" title="Special:BookSources/978-1-64273-003-6">978-1-64273-003-6</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "Before Traveling to Another World"</li>
<li>Chapter One: "The Reprise of the Spear Hero"</li>
<li>Chapter Two: "Repaying Kindness"</li>
<li>Chapter Three: "Leveling"</li>
<li>Chapter Four: "Time Reversal"</li>
<li>Chapter Five: "Trap"</li>
<li>Chapter Six: "Dungeon"</li>
<li>Chapter Seven: "Gerontocracy"</li>
<li>Chapter Eight: "Aiming"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Nine: "Filolial Farmer"</li>
<li>Chapter Ten: "Hallucinations"</li>
<li>Chapter Eleven: "Camping Out"</li>
<li>Chapter Twelve: "Finesse"</li>
<li>Chapter Thirteen: "Peeping"</li>
<li>Chapter Fourteen: "Fitoria-Tan"</li>
<li>Chapter Fifteen: "Church of the Faux Heroes"</li>
<li>Epilogue: "Arrival at Siltvelt"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol2" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">2</th><td> November 25, 2017<sup id="cite_ref-33" class="reference"><a href="#cite_note-33">[33]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-069588-4" title="Special:BookSources/978-4-04-069588-4">978-4-04-069588-4</a></td><td>September 29, 2020<sup id="cite_ref-US_Spear_LN_32-1" class="reference"><a href="#cite_note-US_Spear_LN-32">[32]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-083-8" title="Special:BookSources/978-1-64273-083-8">978-1-64273-083-8</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "The Heroes' Arrival"</li>
<li>Chapter One: "Soulmate"</li>
<li>Chapter Two: "Poison"</li>
<li>Chapter Three: "Beast Spy"</li>
<li>Chapter Four: "In the Back Alley"</li>
<li>Chapter Five: "Shieldfreeden"</li>
<li>Chapter Six: "The Seven Star Whip Hero"</li>
<li>Chapter Seven: "Half-Burnt Charcoal"</li>
<li>Chapter Eight: "Remnants"</li>
<li>Chapter Nine: "Lingering Frangrance"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Ten: "Carved into the Heart"</li>
<li>Chapter Eleven: "A Guarantee of Safety"</li>
<li>Chapter Twelve: "Assassination"</li>
<li>Chapter Thirteen: "False Proof"</li>
<li>Chapter Fourteen: "Avoiding War"</li>
<li>Chapter Fifteen: "A Better World"</li>
<li>Chapter Sixteen: "Bluff"</li>
<li>Chapter Seventeen: "An Order"</li>
<li>Chapter Eighteen: "A Peddling Permit"</li>
<li>Epilogue: "Peddling Preparations"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol3" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">3</th><td> July 25, 2018<sup id="cite_ref-34" class="reference"><a href="#cite_note-34">[34]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-069805-2" title="Special:BookSources/978-4-04-069805-2">978-4-04-069805-2</a></td><td>March 26, 2021<sup id="cite_ref-US_Spear_LN_32-2" class="reference"><a href="#cite_note-US_Spear_LN-32">[32]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-1-64273-106-4" title="Special:BookSources/978-1-64273-106-4">978-1-64273-106-4</a></td></tr><tr style="vertical-align: top; border-bottom: 3px solid rgb(204, 204, 255); --darkreader-inline-border-bottom: #000075;" data-darkreader-inline-border-bottom=""><td colspan="5"><table border="0" cellspacing="0" cellpadding="0" style="width: 100%; background-color: transparent; table-layout: fixed; text-align: left; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor=""><tbody><tr style="vertical-align: top;"></tr></tbody><caption></caption><tbody><tr><td style="border: black; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Prologue: "Peak Racing"</li>
<li>Chapter One: "Carriage-Crafting"</li>
<li>Chapter Two: "A Lazy Pig"</li>
<li>Chapter Three: "A Hidden Passage"</li>
<li>Chapter Four: "Experts Know Best"</li>
<li>Chapter Five: "An Underestimation"</li>
<li>Chapter Six: "Emergency Exit"</li>
<li>Chapter Seven: "Rebirth of the Paradox"</li></ul></td>
<td style="border: black; width: 48%; --darkreader-inline-border-top: #8c8273; --darkreader-inline-border-right: #8c8273; --darkreader-inline-border-bottom: #8c8273; --darkreader-inline-border-left: #8c8273;" data-darkreader-inline-border-top="" data-darkreader-inline-border-right="" data-darkreader-inline-border-bottom="" data-darkreader-inline-border-left="">
<ul><li>Chapter Eight: "Bad Status"</li>
<li>Chapter Nine: "Unfair"</li>
<li>Chapter Ten: "The Egg Raffle"</li>
<li>Chapter Eleven: "Dressing Up the Panda"</li>
<li>Chapter Twelve: "The Fruit of Idleness"</li>
<li>Chapter Thirteen: "Tourist Attractions"</li>
<li>Chapter Fourteen: "The Assistant and the Dragon Girl"</li>
<li>Epilogue: "To Filolials, I Sound My Cry of Passion"</li></ul></td></tr></tbody></table>
</td></tr>
<tr style="text-align: center;"><th scope="row" id="vol4" style="text-align: center; font-weight: normal; background-color: transparent; --darkreader-inline-bgcolor: transparent;" data-darkreader-inline-bgcolor="">4</th><td> September 25, 2023<sup id="cite_ref-35" class="reference"><a href="#cite_note-35">[35]</a></sup></td><td><link rel="mw-deduplicated-inline-style" href="mw-data:TemplateStyles:r1215172403"><a href="/wiki/Special:BookSources/978-4-04-682894-1" title="Special:BookSources/978-4-04-682894-1">978-4-04-682894-1</a></td><td>—</td><td>—</td></tr>
</tbody></table>`
)

// Note: the math does seem off when it comes to how the expected stop index is calculated, but it logically works:
// ExpectedStopIndex = strings.Index(args.InputHtml, args.ExpectedTable) + len(args.ExpectedTable)
var GetNextTableAndItsEndPositionTestCases = map[string]GetNextTableAndItsEndPositionTestCase{
	"a section with a simple table should be properly recognized and pulled out": {
		InputHtml:         theWrongWayToUseHealingMagicLightNovelSection,
		ExpectedTableHtml: theWrongWayToUseHealingMagicLightNovelTable,
		ExpectedStopIndex: 11764,
	},
	// TODO: fix this issue
	"a section with a table with nested tables in it should be properly recognized and pulled out": {
		InputHtml:         theRisingOfTheShieldHereLightNovelSection,
		ExpectedTableHtml: theRisingOfTheShieldHereLightNovelTable,
		ExpectedStopIndex: 68247,
	},
	"a section with no table should come back with an empty string and -1 for the stop value": {
		InputHtml:         "<h1>This is a title<h1>",
		ExpectedTableHtml: "",
		ExpectedStopIndex: -1,
	},
	"a section with a table that lacks a closing table tag should come back with everything after the starting table tag": {
		InputHtml:         `<h1>This is a title<h1> <table class="wikitable"><tbody></tbody>`,
		ExpectedTableHtml: `<table class="wikitable"><tbody></tbody>`,
		ExpectedStopIndex: 64,
	},
}

func TestGetNextTableAndItsEndPosition(t *testing.T) {
	for name, args := range GetNextTableAndItsEndPositionTestCases {
		t.Run(name, func(t *testing.T) {
			actualTableHtml, actualStopPosition := wikipedia.GetNextTableAndItsEndPosition(args.InputHtml)

			assert.Equal(t, args.ExpectedTableHtml, actualTableHtml, "actual html was not the expected value")
			assert.Equal(t, args.ExpectedStopIndex, actualStopPosition, "actual stop position of the table was not the expected value")
		})
	}
}
