//go:build unit

package converter_test

const (
	coverFileMd = `# Church Songs - E Version {{TYPE}}

Date: {{DATE_GENERATED}}
<br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/>

### Key:

#### Books

R=Red Book (Songs We Love)  
MS=More Songs section of Songs We Love  
B=Blue Book

#### Authors

CA= Cyndi Aarrestad  
EHW= Ewald Wanagas  
FTP= Frank Paterson  
GBS= Gail Shepherd  
ZW= Zelma Wanagas

<br/> <br/>

*\*When searching electronic format, punctuation & spelling will cause
non returns*

*\*Punctuation alters the alphabetical order*`
	coverFileHtmlFormat = `<div style="text-align: center%s">
<h1 id="church-songs-e-version-%s">Church Songs - E Version %s</h1>
<p>Date: %s
<br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/> <br/></p>
<h3 id="key">Key:</h3>
<h4 id="books">Books</h4>
<p>R=Red Book (Songs We Love)<br>
MS=More Songs section of Songs We Love<br>
B=Blue Book</p>
<h4 id="authors">Authors</h4>
<p>CA= Cyndi Aarrestad<br>
EHW= Ewald Wanagas<br>
FTP= Frank Paterson<br>
GBS= Gail Shepherd<br>
ZW= Zelma Wanagas</p>
<p><br/> <br/></p>
<p><em>*When searching electronic format, punctuation &amp; spelling will cause
non returns</em></p>
<p><em>*Punctuation alters the alphabetical order</em></p>
</div>
`
	AGloriousChurchFileMd = `---
melody: 
key: 
authors: Ralph E. Hudson
in-church: N
verse: 
location: (MS68)
copyright: Public Domain
type: song
tags: 🎵
---

# A Glorious Church

\~ 1 \~ Do you hear them coming, Brother, Thronging up the steeps of light,  
Clad in glorious shining garments Blood-washed garments pure and white?

CHORUS:  
\'Tis a glorious Church without spot or wrinkle, Washed in the blood of the Lamb.  
\'Tis a glorious Church, without spot or wrinkle, Washed in the blood of the Lamb.

\~ 2 \~ Do you hear the stirring anthems Filling all the earth and sky?  
\'Tis a grand victorious army. Lift its banner up on high!

\~ 3 \~ Never fear the clouds of sorrow; Never fear the storms of sin.  
Even now our joys begin.

\~ 4 \~ Wave the banner, shout His praises, For our victory is nigh!  
We shall join our conquering Savior. We shall reign with Him on high.
`
	AGloriousChurchFileCsv        = "A Glorious Church|(More Songs We Love page 68)|Ralph E. Hudson|Public Domain\n"
	FillMySoulWithThySpiritFileMd = `---
melody: 
key: Key E Flat or F or C
authors: A. Ellis
in-church: Y
verse: 
location: R199 (B15)
type: song
tags: 🎵
---

# Fill My Soul With Thy Spirit

Fill my soul with Thy Spirit, Fill my heart with Thy love;  
Let my soul be rekindled with fire from above.  
Let me drink from that fountain; Flowing boundless and free,  
Fill my soul with Thy Spirit, With love fill thou me.
`
	FillMySoulWithThySpiritFileCsv = "Fill My Soul With Thy Spirit|Red Book page 199 (Blue Book page 15)|A. Ellis|Church\n"
	FlowThowRiverFileMd            = `---
melody: 
key: Key C or E Flat
authors: M. Heartwell 
in-church: Y
verse: 
location: R32 (MS13) (B15)
type: song
tags: 🎵
---

# Flow Thou River

\~ 1 \~ Flow thou River, flow thou River, Forth from the Throne of God.  
Flow thou River, flow thou River, Forth from the Throne of God.

\~ 2 \~ \* That life may spring forth,  
That life may spring forth in a dry and thirsty land. That life may spring forth,  
That life may spring forth in a dry and thirsty land.

\~ 3 \~ \* Healing waters, healing waters, Flow from the Throne of God.  
Healing waters, healing waters, Flow from the Throne of God.

\~ 4 \~ \*\* Strength for today, strength for tomorrow, Flows from the Throne of God.  
Strength for today, strength for tomorrow, Flows from the Throne of God.

\*Differs from published songs by M. Heartwell  
\*\*Added by G.B.S.
`
	FlowThowRiverFileCsv = "Flow Thou River|Red Book page 32 (More Songs We Love page 13) (Blue Book page 15)|M. Heartwell|Church\n"
	HeIsLordFileMd       = `---
melody: 
key: Key F
authors: 
in-church: 
verse: 
location: R53
type: song
tags: 🎵
---

# He Is Lord

He is Lord, He is Lord. He is risen from the dead and He is Lord.  
Ev\'ry knee shall bow ev\'ry tongue confess, That Jesus Christ is Lord.
`
	HeIsLordFileCsv      = "He is Lord|Red Book page 53||\n"
	BlessThisHouseFileMd = `---
melody: 
key: 
authors: May H Brahe, Helen Taylor
in-church: N
verse: 
location: 
type: song
tags: 🎵
---

# Bless This House

\~ 1 \~ Bless this house, O Lord, we pray. Make it safe by night and day.  
Bless these walls so firm and stout, Keeping want and troubles out.

\~ 2 \~ Bless the roof and chimney top. Let thy love flow all about.  
Bless this house that it may prove Ever open to joy and truth.

\~ 3 \~ Bless us all that we may be Fit, O Lord, to dwell with thee.  
Bless us so that one day we, May dwell, dear Lord, with thee.
`
	BlessThisHouseFileCsv           = "Bless This House||May H Brahe, Helen Taylor|\n"
	BiggerThanAllOfMyProblemsFileMd = `---
melody: 
key: Key C
authors: Bill & Gloria Gaither 
in-church: N
verse: 
location: (B6)
type: song
tags: 🎵
---

# Bigger Than All My Problems

\~ 1 \~ Bigger than all the shadows that fall across my path  
God is bigger than any mountain that I can or cannot see;  
Bigger than my confusion, bigger than anything,  
God is bigger than any mountain that I can or cannot see.

CHORUS:  
Bigger than all my problems, bigger than all my fears;  
God is bigger than any mountain that I can or cannot see;  
Bigger than all my questions, bigger than anything,  
God is bigger than any mountain that I can or cannot see.

\~ 2 \~ Bigger than all the giants of fear and unbelief,  
God is bigger than any mountain that I can or cannot see;  
Bigger than all my hang-ups, bigger than anything,  
God is bigger than any mountain that I can or cannot see.
`
	BiggerThanAllOfMyProblemsFileCsv = "Bigger Than All My Problems|(Blue Book page 6)|Bill & Gloria Gaither|\n"
	HeIsFileMd                       = `---
melody: 
key: Key C
authors: Chris Knauf 
in-church: N
verse: 
location: (MS16) (B20)
type: song
tags: 🎵
---

# He Is

\~ 1 \~ He is fairer than the lily of the valley, He is brighter than the morning star.  
He is purer than the snow, fresher than the breeze, Lovelier by far than all of these.

\~ 2 \~ But He calms all the storms, and conquers the raging seas.  
He is the high and lofty One who inhabits eternity,  
Creator of the universe, and He\'s clothed with majesty.  
He is and ever more shall be.  
(Repeat first stanza)
`
	HeIsFileHtml = `<div class="keep-together">
<h1 id="he-is">He Is</h1>
<div><div class="metadata"><div><div class="author">Chris Knauf</div></div><div><div class="key"><b>Key C</b></div></div><div><div class="location">(MS16) (B20)</div></div></div></div><br>
<p>~ 1 ~ He is fairer than the lily of the valley, He is brighter than the morning star.<br>
He is purer than the snow, fresher than the breeze, Lovelier by far than all of these.</p>
<p>~ 2 ~ But He calms all the storms, and conquers the raging seas.<br>
He is the high and lofty One who inhabits eternity,<br>
Creator of the universe, and He&rsquo;s clothed with majesty.<br>
He is and ever more shall be.<br>
(Repeat first stanza)</p>
</div>
<br>`
	AboveItAllFileMd = `---
melody: 
key: Key G
authors: LaVerne & Edith Tripp
in-church: N
verse: 
location: (MS68) (B1)
copyright: unknown
type: song
tags: 🎵
---

# Above It All (There Stands Jesus)

Above it all, There stands Jesus. Above it all, He\'s still my King. \*He took my life,  
And He made me happy\*. Above it all, He\'s still the \*\*same\*\*.

\*This fleeting life is but a vapor\*  
\*\*King\*\*
`
	AboveItAllFileHtml = `<div class="keep-together">
<h1 id="above-it-all-there-stands-jesus">Above It All <span class="other-title">(There Stands Jesus)</span></h1>
<div><div class="metadata"><div><div class="author">LaVerne & Edith Tripp</div></div><div><div class="key"><b>Key G</b></div></div><div><div class="location">(MS68) (B1)</div></div></div></div><br>
<p>Above it all, There stands Jesus. Above it all, He&rsquo;s still my King. *He took my life,<br>
And He made me happy*. Above it all, He&rsquo;s still the **same**.</p>
<p>*This fleeting life is but a vapor*<br>
**King**</p>
</div>
<br>`
	BeholdTheHeavensFileMd = `---
melody: (tune The Kingdom of God is Not Meat and Drink)
key: 
authors: I. Amundson
in-church: Y
verse: 
location: 
type: song
tags: 🎵
---

# Behold The Heavens

\~ 1 \~ Behold the heavens are open; Behold the face of the King  
We now have a new way of walking; We now have a new song to sing.

\~ 2 \~ Behold the heavens are open; Forever the veil has been rent;  
We\'re walking together in Zion; And there are no limits in Him.
`
	BeholdTheHeavensFileHtml = `<div class="keep-together">
<h1 id="behold-the-heavens">Behold The Heavens</h1>
<div><div class="metadata row-padding"><div><div class="author"><b>I. Amundson</b></div></div><div><div class="key">&nbsp;&nbsp;&nbsp;&nbsp;</div></div><div><div class="location">&nbsp;&nbsp;&nbsp;&nbsp;</div></div></div><div class="metadata"><div><div class="melody-75"><b>(tune The Kingdom of God is Not Meat and Drink)</b></div></div></div></div><br>
<p>~ 1 ~ Behold the heavens are open; Behold the face of the King<br>
We now have a new way of walking; We now have a new song to sing.</p>
<p>~ 2 ~ Behold the heavens are open; Forever the veil has been rent;<br>
We&rsquo;re walking together in Zion; And there are no limits in Him.</p>
</div>
<br>`
	BeThouExaltedFileMd = `---
melody: 
key: Key F
authors: 
in-church: 
verse: Ps. 57:9-11
location: (MS4) (B4)
type: song
tags: 🎵
---

# Be Thou Exalted

\~ 1 \~ Be Thou exalted, oh God (x3) above the heavens.  
Let Thy glory be above the whole earth.  
(Repeat)

\~ 2 \~ I will praise Thee oh Lord among the people.  
I will sing unto Thee among the nations.  
For Thy mercy is great unto the heavens  
And Thy truth unto the clouds.
`
	BeThouExaltedFileHtml = `<div class="keep-together">
<h1 id="be-thou-exalted">Be Thou Exalted</h1>
<div><div class="metadata row-padding"><div><div class="author">&nbsp;&nbsp;&nbsp;&nbsp;</div></div><div><div class="key"><b>Key F</b></div></div><div><div class="location">(MS4) (B4)</div></div></div><div class="metadata"><div><div class="melody">&nbsp;&nbsp;&nbsp;&nbsp;</div></div><div><div class="verse">Ps. 57:9-11</div></div></div></div><br>
<p>~ 1 ~ Be Thou exalted, oh God (x3) above the heavens.<br>
Let Thy glory be above the whole earth.<br>
(Repeat)</p>
<p>~ 2 ~ I will praise Thee oh Lord among the people.<br>
I will sing unto Thee among the nations.<br>
For Thy mercy is great unto the heavens<br>
And Thy truth unto the clouds.</p>
</div>
<br>`
	BeThouExalted2FileHtml = `<div class="keep-together">
<h1 id="be-thou-exalted-2">Be Thou Exalted</h1>
<div><div class="metadata row-padding"><div><div class="author">&nbsp;&nbsp;&nbsp;&nbsp;</div></div><div><div class="key"><b>Key F</b></div></div><div><div class="location">(MS4) (B4)</div></div></div><div class="metadata"><div><div class="melody">&nbsp;&nbsp;&nbsp;&nbsp;</div></div><div><div class="verse">Ps. 57:9-11</div></div></div></div><br>
<p>~ 1 ~ Be Thou exalted, oh God (x3) above the heavens.<br>
Let Thy glory be above the whole earth.<br>
(Repeat)</p>
<p>~ 2 ~ I will praise Thee oh Lord among the people.<br>
I will sing unto Thee among the nations.<br>
For Thy mercy is great unto the heavens<br>
And Thy truth unto the clouds.</p>
</div>
<br>`

	AHymnOfPraiseFileMd = `---
melody: 
key: Key C
authors: From Hyden&nbsp;&nbsp;&nbsp;&nbsp;F.T.P.
in-church: Y
verse: 
location: R4
copyright: 
type: song
tags: 🎵
---

# A Hymn Of Praise

\~ 1 \~ A hymn of praise I sing unto Thee,  
My gracious Redeemer, my wonderful Lord.  
Thy perfect salvation, thy wonderful Love,  
Let all the creation sing praise unto Thee.

\~ 2 \~ In the beginning Thy word spoke as life;  
Unto creation, the work of Thy hands.  
Heaven rejoiceth, the morning stars sang;  
A hymn of creation and praise unto Thee.
`
	AHymnOfPraiseFileCsvCleaned = "A Hymn Of Praise|Red Book page 4|From Hyden F.T.P.|Church\n"
)
