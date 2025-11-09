//go:build unit

package potentiallyfixableissue_test

import (
	"testing"

	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	"github.com/stretchr/testify/assert"
)

type getPotentialIncorrectSingleQuotesTestCase struct {
	inputText           string
	expectedSuggestions map[string]string
}

var getPotentialIncorrectSingleQuotesTestCases = map[string]getPotentialIncorrectSingleQuotesTestCase{
	"make sure that a file with single quoted words outside of double quotes should have its values updated": {
		inputText: `<p>He said 'hello' and 'goodbye'</p>`,
		expectedSuggestions: map[string]string{
			"<p>He said 'hello' and 'goodbye'</p>": `<p>He said "hello" and "goodbye"</p>`,
		},
	},
	"make sure that a file with a contraction does not get a suggestion": {
		inputText:           `<p>It's a wonderful day today isn't it?</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a single quoted word inside of double quotes does not get a suggestion": {
		inputText:           `<p>She said "don't touch the 'red' button"</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a possesive use of a single quote does not get a suggestion": {
		inputText:           `<p>Charles' was not happy. Twas harder on Bob's faculties though.</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a word ending with s that is single quoted gets a suggestion and is not considered a possesive": {
		inputText: `<p>Charles' 'Great Expectations' is long</p>`,
		expectedSuggestions: map[string]string{
			"<p>Charles' 'Great Expectations' is long</p>": `<p>Charles' "Great Expectations" is long</p>`,
		},
	},
	"make sure that a file with a contraction and single quoted word outside of double quotes should get a suggestion for just the single quoted word": {
		inputText: `<p>It's a 'beautiful' day</p>`,
		expectedSuggestions: map[string]string{
			"<p>It's a 'beautiful' day</p>": `<p>It's a "beautiful" day</p>`,
		},
	},
	"make sure that a file with a single quoted word in between two sets of quotes gets a suggestion with it double quoted": {
		inputText: `<p>He had to be ready as he came to defeat the enemy.</p>
<p>He said, "Hello. Can you hear that?". He stalled to get read to use his 'special move'. He called out again, "Are you there?"</p>`,
		expectedSuggestions: map[string]string{
			`<p>He said, "Hello. Can you hear that?". He stalled to get read to use his 'special move'. He called out again, "Are you there?"</p>`: `<p>He said, "Hello. Can you hear that?". He stalled to get read to use his "special move". He called out again, "Are you there?"</p>`,
		},
	},
	"make sure that a file with a pluralized number is not considered to be a single quote and does not cause a suggestion by itself": {
		inputText: `	<p>Normally, simply drawing a weapon in front of the Emperor in court was a grave offense, to say nothing of actually offering violence to a member of the Imperial family. However, the court was currently paralyzed in the wake of the earthquake. The Praetorians who should have defended the Emperor and his household were missing. Since there was nobody to maintain order, the area before the throne was a sea of chaos.</p>
  <p>Tomita, watching from the side, flicked his Type 64's fire selector to レ(automatic fire), while Kuribayashi inspected Tyuule and the black-haired girl on the ground.</p>
  <p>"Are you alright?"</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with the contract \"'Cause\" does not cause a suggestion to be made on its own": {
		inputText:           `	<p>He said, "'Cause there ain't enough room for the both of us."</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with the contract \"'em\" does not cause a suggestion to be made on its own": {
		inputText:           `	<p>He said, "Get 'em!"</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a single quote usage with plural decade (1960's) does not cause a suggestion": {
		inputText:           `<p>Music from the 1960's is classic.</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a single quote usage with omitted numerals at start of year ('99) does not cause a suggestion": {
		inputText:           `<p>I graduated in '99.</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a single quote usage with a number possessive (90's influence) does not cause a suggestion": {
		inputText:           `<p>The 90's influence on fashion is unmistakable.</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a single quote usage with colloquial omission and decade ('80s) does not cause a suggestion": {
		inputText:           `<p>Rocking in the '80s was the best!</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a single quoted word ending with a double quote is handled properly not getting a suggestion": {
		inputText:           `<p>"First, take this," Giraud said, holding the tool he had assembled over to Dean. It appeared to be something like a thick mask. "That is a make-up tool called a 'facemaker'"</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a single quoted sentence that has a contraction in it is considered fine": {
		inputText:           `<p><strong>"Victorique! Are you there? What took you so long? I bet you were reading a thick Latin book again, eating macaroons, saying 'Who's Kujou?'. Hello?"</strong></p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that when a contraction that includes the word `I` in it starts a quote inside a quote that it is recognized as a contraction": {
		inputText:           `<p><strong>"He said 'I'm going to the store', before he left"</strong></p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a simple version of a possesive in a single quote works": {
		inputText:           `<p><strong>"He said 'Hello there. Boss' work is not done yet.', before he left"</strong></p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure tha a possessive that starts a single quote works": {
		inputText:           `<p><strong>"He said 'Boss' work is not done yet.', before he left"</strong></p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that multiple possessives in a single quotes works": {
		inputText:           `<p><strong>"He said 'Boss' blabberers' work is not done yet.', before he left"</strong></p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a single quote that ends with the letter s and has a possesive works": {
		inputText:           `<p><strong>"He said 'Boss' work is', before he left"</strong></p>`,
		expectedSuggestions: map[string]string{},
	},
}

func TestGetPotentialIncorrectSingleQuotes(t *testing.T) {
	for name, args := range getPotentialIncorrectSingleQuotesTestCases {
		t.Run(name, func(t *testing.T) {
			actual, err := potentiallyfixableissue.GetPotentialIncorrectSingleQuotes(args.inputText)

			assert.Nil(t, err)
			assert.Equal(t, args.expectedSuggestions, actual)
		})
	}
}

const testFile = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>Quote Conversion Test Cases</title>
    <style type="text/css">
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        .section {
            margin-bottom: 30px;
            padding: 20px;
            border: 1px solid #ccc;
            border-radius: 5px;
        }
        .quote-example {
            margin: 10px 0;
            padding: 10px;
            background-color: #f5f5f5;
        }
        h2 {
            color: #333;
            border-bottom: 2px solid #666;
            padding-bottom: 5px;
        }
        .note {
            font-style: italic;
            color: #666;
        }
    </style>
</head>
<body>
    <h1>Quote Conversion Test Cases</h1>
    
    <div class="section">
        <h2>Literary Quotes</h2>
        <p class="quote-example">The professor began, 'It's important to understand that Shakespeare's 'Hamlet' represents one of the finest examples of tragic literature.'</p>
        <p class="quote-example">'Don't you think,' she asked, 'that Dickens' 'Oliver Twist' perfectly captures the Victorian era's social issues?'</p>
        <p class="quote-example">The critic wrote, 'Jane Austen's 'Pride and Prejudice' remains a masterpiece of social commentary.'</p>
        <p class="quote-example">'I've always felt that Poe's 'The Raven' exemplifies Gothic poetry at its finest,' the student remarked.</p>
    </div>

    <div class="section">
        <h2>Dialogue Examples</h2>
        <p class="quote-example">'Hello,' she said, 'I'm glad you could make it to today's meeting.'</p>
        <p class="quote-example">The manager announced, 'We'll be implementing new policies starting next week's meeting.'</p>
        <p class="quote-example">'That's interesting,' he replied, 'but I don't quite understand the implications.'</p>
        <p class="quote-example">'It's time,' the captain declared, 'to set sail for tomorrow's adventure.'</p>
    </div>

    <div class="section">
        <h2>Nested Quotations</h2>
        <p class="quote-example">She remembered the teacher saying, 'When you read 'The Great Gatsby,' pay attention to the green light's symbolism.'</p>
        <p class="quote-example">The historian noted, 'Churchill's famous 'We shall fight them on the beaches' speech rallied the nation's spirit.'</p>
        <p class="quote-example">'I recall my grandmother's favorite saying: 'A stitch in time saves nine,'' he mentioned with a smile.</p>
        <p class="quote-example">The guide explained, 'The inscription reads 'Per aspera ad astra,' which means 'through hardships to the stars.''</p>
    </div>

    <div class="section">
        <h2>Academic Examples</h2>
        <p class="quote-example">The researcher stated, 'It's crucial to note that Smith's 'Theory of Economic Growth' contradicts earlier models.'</p>
        <p class="quote-example">'In today's lecture,' the professor began, 'we'll examine Einstein's 'Theory of Relativity.''</p>
        <p class="quote-example">The paper concluded, 'Darwin's 'Origin of Species' revolutionized our understanding of evolution.'</p>
        <p class="quote-example">'According to recent studies,' the scientist explained, 'we're seeing unprecedented changes in climate patterns.'</p>
    </div>

    <div class="section">
        <h2>Historical Quotes</h2>
        <p class="quote-example">The textbook stated, 'Lincoln's 'Gettysburg Address' remains one of America's most significant speeches.'</p>
        <p class="quote-example">'In examining Kennedy's 'Ask not what your country can do for you' speech,' the historian began, 'we see a call to civic duty.'</p>
        <p class="quote-example">'Let's consider Gandhi's famous saying: 'Be the change you wish to see in the world,'' the teacher suggested.</p>
        <p class="quote-example">The documentary narrated, 'Martin Luther King Jr.'s 'I Have a Dream' speech transformed the civil rights movement.'</p>
    </div>

    <div class="section">
        <h2>Business Contexts</h2>
        <p class="quote-example">'We're implementing what we call 'Project Phoenix,'' the CEO announced at today's meeting.</p>
        <p class="quote-example">The consultant explained, 'It's essential to understand that 'customer-first' isn't just a slogan.'</p>
        <p class="quote-example">'In today's market,' the analyst began, 'we're seeing what I call 'aggressive growth patterns.''</p>
        <p class="quote-example">The report stated, 'Companies that don't embrace 'digital transformation' risk becoming obsolete.'</p>
    </div>

    <div class="section">
        <h2>Technical Documentation</h2>
        <p class="quote-example">'It's important to note,' the manual stated, 'that the 'sleep mode' function conserves battery life.'</p>
        <p class="quote-example">The documentation explained, 'When the system displays 'Error 404,' it means the page wasn't found.'</p>
        <p class="quote-example">'In this tutorial,' the instructor began, 'we'll explore what's called 'responsive design.''</p>
        <p class="quote-example">The guide noted, 'If you see the message 'Connection refused,' check your network settings.'</p>
    </div>

    <div class="section">
        <h2>Cultural References</h2>
        <p class="quote-example">The critic wrote, 'The band's 'Symphony of Change' represents a departure from their usual style.'</p>
        <p class="quote-example">'In this exhibition,' the curator explained, 'we're featuring what's called 'abstract expressionism.''</p>
        <p class="quote-example">The review stated, 'The director's 'Vision of Tomorrow' challenges conventional storytelling.'</p>
        <p class="quote-example">'Let's examine the artist's 'Blue Period,'' the lecturer suggested, 'and its influence on modern art.'</p>
    </div>

    <div class="section">
        <h2>Educational Contexts</h2>
        <p class="quote-example">'Today's lesson,' the teacher announced, 'will focus on what's called 'critical thinking.''</p>
        <p class="quote-example">The textbook explained, 'It's essential to understand the concept of 'photosynthesis' in biology.'</p>
        <p class="quote-example">'In mathematics,' the instructor noted, 'we're introducing the concept of 'variable expressions.''</p>
        <p class="quote-example">The guide stated, 'Let's explore what we call 'the scientific method.''</p>
    </div>

    <div class="section">
        <h2>Mixed Complex Examples</h2>
        <p class="quote-example">'It's fascinating,' the professor remarked, 'how Shakespeare's 'To be or not to be' soliloquy captures Hamlet's existential crisis.'</p>
        <p class="quote-example">The journalist wrote, 'During yesterday's press conference, the diplomat stated, 'We're committed to maintaining what's called 'strategic patience.'''</p>
        <p class="quote-example">'I'm reminded of my father's favorite saying: 'It's not about the destination, it's about the journey,'' she reflected during today's speech.</p>
        <p class="quote-example">The philosopher mused, 'In considering Descartes' 'I think, therefore I am,' we're confronting the essence of consciousness.'</p>
    </div>

    <div class="note">
        <p>Note: All examples contain valid quote patterns that should be properly handled by the quote conversion algorithm. Each example includes a mix of contractions (it's, don't, we're) and quotations that should be converted.</p>
    </div>
</body>
</html>`

func BenchmarkGetPotentialIncorrectSingleQuotes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		potentiallyfixableissue.GetPotentialIncorrectSingleQuotes(testFile)
	}
}
