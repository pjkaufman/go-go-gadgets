//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/linter"
	"github.com/stretchr/testify/assert"
)

type GetPotentiallyBrokenLinesTestCase struct {
	InputText           string
	ExpectedSuggestions map[string]string
}

var GetPotentiallyBrokenLinesTestCases = map[string]GetPotentiallyBrokenLinesTestCase{
	"make sure that a file with no potentially broken paragraphs gives no suggestions": {
		InputText: `<p>Here is some content.</p>
	<p>Here is some more content</p>`,
		ExpectedSuggestions: map[string]string{},
	},
	"make sure that a file with paragraphs that end in a letter get picked up as potentially needing a change": {
		InputText: `<p>Here is some content.</p>
			<p class="calibre1"><a id="p169"></a>If there are objects with a simple structure and the same properties, then they can be recognized as a single "set" allowing decomposition of the </p>
			<p class="calibre1">"set" rather than each object separately. </p>`,
		ExpectedSuggestions: map[string]string{
			`
			<p class="calibre1"><a id="p169"></a>If there are objects with a simple structure and the same properties, then they can be recognized as a single "set" allowing decomposition of the </p>
			<p class="calibre1">"set" rather than each object separately. </p>`: `
			<p class="calibre1"><a id="p169"></a>If there are objects with a simple structure and the same properties, then they can be recognized as a single "set" allowing decomposition of the "set" rather than each object separately. </p>`,
		},
	},
	"make sure that a file with paragraphs that end in a comma get picked up as potentially needing a change": {
		InputText: `<p class="calibre1">The information provided by Edward Clark had a brief project plan. </p>
			<p class="calibre1">Tatsuya ran his eyes through the original document, and Miyuki read the translated text. </p>
			<p class="calibre1">Minami placed a cup of freshly brewed tea in front of them. As if on a signal, Tatsuya and Miyuki simultaneously looked up from the electronic paper with the details of the project. </p>
			<p class="calibre1">"…The so-called Dione, it does not seem to be Saturn's companion </p>
			<p class="calibre1">'Dione', but the goddess of Greek myths." </p>
			<p class="calibre1">"Right. Wife of Zeus, who gave birth to Aphrodite. From that version of the myth, where Aphrodite is born from sea foam." </p>`,
		ExpectedSuggestions: map[string]string{
			`
			<p class="calibre1">"…The so-called Dione, it does not seem to be Saturn's companion </p>
			<p class="calibre1">'Dione', but the goddess of Greek myths." </p>`: `
			<p class="calibre1">"…The so-called Dione, it does not seem to be Saturn's companion 'Dione', but the goddess of Greek myths." </p>`,
		},
	},
	"make sure that a file with paragraphs that end in a number get picked up as potentially needing a change": {
		InputText: `<p class="calibre1">The Deputy Director showed interest in these words and encouraged him to continue. The Director of the Intelligence Department was absent at this meeting. The Deputy Director who was present was also a person not disclosed to the public. </p>
			<p class="calibre1">"I think you all already know that the Yotsuba family, to which Shiba Tatsuya belongs, is in a cooperative relationship with the 1-0-1 </p>
			<p class="calibre1">Battalion." </p>
			<p class="calibre1">After Onda's words, those sitting at the table nodded. </p>`,
		ExpectedSuggestions: map[string]string{
			`
			<p class="calibre1">"I think you all already know that the Yotsuba family, to which Shiba Tatsuya belongs, is in a cooperative relationship with the 1-0-1 </p>
			<p class="calibre1">Battalion." </p>`: `
			<p class="calibre1">"I think you all already know that the Yotsuba family, to which Shiba Tatsuya belongs, is in a cooperative relationship with the 1-0-1 Battalion." </p>`,
		},
	},
	"make sure that a file with multiple paragraphs back to back that potentially are broken, get condensed down into 1 suggestions": {
		InputText: `<p>Some content here.</p>
			<p>Here is a list, </p>
			<p>a set of todos,</p>
			<p>and its own sentence. </p>`,
		ExpectedSuggestions: map[string]string{
			`
			<p>Here is a list, </p>
			<p>a set of todos,</p>
			<p>and its own sentence. </p>`: `
			<p>Here is a list, a set of todos, and its own sentence. </p>`,
		},
	},
	"make sure that a file with a line with an odd number of double quotes gets a suggestion": {
		InputText: `<p class="calibre1">Saeki and Kazama sighed at the same time. And they exchanged tense smiles. They thought it was funny how they seriously discussed the concepts, like the final boss and the hero. </p>
			<p class="calibre1">"We will pass your opinion to the Intelligence Department through Major Onda. I don't know how much it will help… Thanks for the help. </p>
			<p class="calibre1">Lieutenant-Colonel, you are free." </p>
			<p class="calibre1">"Understood." </p>
			<p class="calibre1">Kazama saluted Saeki and left her office. </p><hr class="character" />`,
		ExpectedSuggestions: map[string]string{
			`
			<p class="calibre1">"We will pass your opinion to the Intelligence Department through Major Onda. I don't know how much it will help… Thanks for the help. </p>
			<p class="calibre1">Lieutenant-Colonel, you are free." </p>`: `
			<p class="calibre1">"We will pass your opinion to the Intelligence Department through Major Onda. I don't know how much it will help… Thanks for the help. Lieutenant-Colonel, you are free." </p>`,
		},
	},
	"make sure that we properly determine the original issue for missing double quotes": {
		InputText: `<p class="calibre1">This requires only a "spell" of auto-suggestion, which stimulates the Magic Calculation Area to create a magic sequence. It just takes more time to activate the magic without a CAD. </p>
		<p class="calibre1">"Set: decreasing entropy • density control • phase transition • </p>
		<p class="calibre1">condensation • transformation of the form of energy • acceleration • </p>
		<p class="calibre1">sublimation: input! Perform a modification of the phenomenon! Magic 'Dry meteor'!" </p>
		<p class="calibre1">If you are able to clearly define a concept with words, and introduce this concept into yourself, then you don't need to say it out loud. </p>
		<p class="calibre1">But when you're in front of the enemy, it's too slow to do so. This is equivalent to being under attack all this time. Modern magic discarded casting tools and chose CADs to avoid this. </p>`,
		ExpectedSuggestions: map[string]string{
			`
		<p class="calibre1">"Set: decreasing entropy • density control • phase transition • </p>
		<p class="calibre1">condensation • transformation of the form of energy • acceleration • </p>
		<p class="calibre1">sublimation: input! Perform a modification of the phenomenon! Magic 'Dry meteor'!" </p>`: `
		<p class="calibre1">"Set: decreasing entropy • density control • phase transition • condensation • transformation of the form of energy • acceleration • sublimation: input! Perform a modification of the phenomenon! Magic 'Dry meteor'!" </p>`,
		},
	},
	"make sure that we properly handle Mr. at the end of a sentence": {
		InputText: `<p class="calibre1">Momoyama took out a white envelope from the drawer of the table and laid it on the table. </p>
		<p class="calibre1">"Here is written the request "To make sure that Taurus Silver, Mr. </p>
		<p class="calibre1">Tatsuya Shiba was able to take part in Project Dione." The National Security Agency of the USNA has come to the conclusion that you are Taurus Silver, and requests your participation in the project." </p>
		<p class="calibre1">"Principle. I'm still a high school student at this school. I'm not going to interrupt my studies halfway." </p>
		<p class="calibre1">Tatsuya didn't answer the question "Are you Taurus Silver?". He deliberately ignored this part, and fundamentally refused to participate in the project, or rather rejected it. </p>`,
		ExpectedSuggestions: map[string]string{
			`
		<p class="calibre1">"Here is written the request "To make sure that Taurus Silver, Mr. </p>
		<p class="calibre1">Tatsuya Shiba was able to take part in Project Dione." The National Security Agency of the USNA has come to the conclusion that you are Taurus Silver, and requests your participation in the project." </p>`: `
		<p class="calibre1">"Here is written the request "To make sure that Taurus Silver, Mr. Tatsuya Shiba was able to take part in Project Dione." The National Security Agency of the USNA has come to the conclusion that you are Taurus Silver, and requests your participation in the project." </p>`,
		},
	},
	"make sure that we properly handle % at the end of a sentence": {
		InputText: `<p class="calibre1">On Monday, when Raymond Clark in the form of "the first Sage" performed on TV, the population's interest in Taurus Silver rose sharply. But the next day the interest began to fade away, and today, on Wednesday, it has ceased to be a topic of discussion among ordinary people. </p>
		<p class="calibre1"><a id="p57"></a>There are no people in the world of magic who do not know the famous Taurus Silver. However, those who can use magic make up only one thousandth of the adult population. But this does not mean that 99.99% </p>
		<p class="calibre1">of people live without magic, because some people are related to magic in the role of engineers, managers, politicians, soldiers, and other civil servants, even without the ability to use magic in practice. </p>
		<p class="calibre1">In recent years, even people who are engaged in magic began to appear, while speaking of the Anti-magician movement. </p>
		<p class="calibre1">Quite a few citizens indirectly benefit from the use of magic for public order, national defense and disaster response. However, most people still live without a direct relationship to magic. </p>`,
		ExpectedSuggestions: map[string]string{
			`
		<p class="calibre1"><a id="p57"></a>There are no people in the world of magic who do not know the famous Taurus Silver. However, those who can use magic make up only one thousandth of the adult population. But this does not mean that 99.99% </p>
		<p class="calibre1">of people live without magic, because some people are related to magic in the role of engineers, managers, politicians, soldiers, and other civil servants, even without the ability to use magic in practice. </p>`: `
		<p class="calibre1"><a id="p57"></a>There are no people in the world of magic who do not know the famous Taurus Silver. However, those who can use magic make up only one thousandth of the adult population. But this does not mean that 99.99% of people live without magic, because some people are related to magic in the role of engineers, managers, politicians, soldiers, and other civil servants, even without the ability to use magic in practice. </p>`,
		},
	},
	"make sure that we properly handle – at the end of a sentence": {
		InputText: `<p class="calibre1">Here is some text.</p>
		<p class="calibre1">Once again, Minoru was at a loss for words. Miyuki was telling the truth – Minoru was forced to admit he had miscalculated. Leading Tatsuya away was not enough – </p>
		<p class="calibre1">when planning a diversion, it was more important to lead Miyuki away. </p>`,
		ExpectedSuggestions: map[string]string{
			`
		<p class="calibre1">Once again, Minoru was at a loss for words. Miyuki was telling the truth – Minoru was forced to admit he had miscalculated. Leading Tatsuya away was not enough – </p>
		<p class="calibre1">when planning a diversion, it was more important to lead Miyuki away. </p>`: `
		<p class="calibre1">Once again, Minoru was at a loss for words. Miyuki was telling the truth – Minoru was forced to admit he had miscalculated. Leading Tatsuya away was not enough – when planning a diversion, it was more important to lead Miyuki away. </p>`,
		},
	},
	"make sure that we properly handle — at the end of a sentence": {
		InputText: `<p>Text here.</p>
		<p class="calibre1">Blood vessels, nerves, and other tissues found in the body that met in a straight line —</p>
		<p class="calibre1">everything was also decomposed. </p>
		<p class="calibre1">A hole was drilled in the joint of the right shoulder. </p>`,
		ExpectedSuggestions: map[string]string{
			`
		<p class="calibre1">Blood vessels, nerves, and other tissues found in the body that met in a straight line —</p>
		<p class="calibre1">everything was also decomposed. </p>`: `
		<p class="calibre1">Blood vessels, nerves, and other tissues found in the body that met in a straight line — everything was also decomposed. </p>`,
		},
	},
	"make sure that we properly handle Mt. at the end of a sentence": {
		InputText: `<p>Text here.</p>
		<p class="calibre1">How many mountains are the size of Mt.</p>
		<p class="calibre1">Rogers? </p>`,
		ExpectedSuggestions: map[string]string{
			`
		<p class="calibre1">How many mountains are the size of Mt.</p>
		<p class="calibre1">Rogers? </p>`: `
		<p class="calibre1">How many mountains are the size of Mt. Rogers? </p>`,
		},
	},
	"make sure that we properly handle a line starting with a lowercase letter": {
		InputText: `<p>First line.</p>
		<p>Text here...</p>
		<p class="calibre1">here is a continuation of the previous line.</p>`,
		ExpectedSuggestions: map[string]string{
			`
		<p>Text here...</p>
		<p class="calibre1">here is a continuation of the previous line.</p>`: `
		<p>Text here... here is a continuation of the previous line.</p>`,
		},
	},
}

func TestGetPotentiallyBrokenLines(t *testing.T) {
	for name, args := range GetPotentiallyBrokenLinesTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentiallyBrokenLines(args.InputText)

			assert.Equal(t, args.ExpectedSuggestions, actual)
		})
	}
}
