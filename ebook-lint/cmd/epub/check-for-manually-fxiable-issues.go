package epub

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
	"github.com/spf13/cobra"
)

var (
	runAll            bool
	runBrokenLines    bool
	runSectionBreak   bool
	runPageBreak      bool
	runOxfordCommas   bool
	runAlthoughBut    bool
	runThoughts       bool
	runConversation   bool
	runNecessaryWords bool
)

const (
	OneRunBoolArgMustBeEnabled   = "either run-all, run-broken-lines, run-section-breaks, run-page-breaks, run-oxford-commas, or run-although-but must be specified"
	CssPathsEmptyWhenArgIsNeeded = "css-paths must have a value when including handling section or page breaks"
)

// fixableCmd represents the fixable command
var fixableCmd = &cobra.Command{
	Use:   "fixable",
	Short: "Runs the specified fixable actions that require manual input to determine what to do.",
	Example: heredoc.Doc(`To run all of the possible potential fixes:
	ebook-lint epub fixable -f test.epub -a
	Note: this will require a css file to already exist in the epub
	
	To just fix broken paragraph endings:
	ebook-lint epub fixable -f test.epub -b

	To just update section breaks:
	ebook-lint epub fixable -f test.epub -s
	Note: this will require a css file to already exist in the epub

	To just update page breaks:
	ebook-lint epub fixable -f test.epub -p
	Note: this will require a css file to already exist in the epub

	To just fix missing oxford commas:
	ebook-lint epub fixable -f test.epub -o

	To just fix although but instances:
	ebook-lint epub fixable -f test.epub -n

	To just fix instances of thoughts in parentheses:
	ebook-lint epub fixable -f test.epub -t

	To run a combination of options:
	ebook-lint epub fixable -f test.epub -otn
	`),
	Long: heredoc.Doc(`Goes through all of the content files and runs the specified fixable actions on them asking
	for user input on each value found that matches the potential fix criteria.
	Potential things that can be fixed:
	- Broken paragraph endings
	- Section breaks being hardcoded instead of an hr
	- Page breaks being hardcoded instead of an hr
	- Oxford commas being missing before or's or and's
	- Possible instances of sentences that have although ..., but in them
	- Possible instances of thoughts that are in parentheses
	- Possible instances of conversation encapsulated in square brackets
	- Possible instances of words in square brackets that may be necessary for the sentence (i.e. need to have the brackets removed)
	`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateManuallyFixableFlags(epubFile, runAll, runBrokenLines, runSectionBreak, runPageBreak, runOxfordCommas, runAlthoughBut, runThoughts, runConversation, runNecessaryWords)
		if err != nil {
			logger.WriteError(err.Error())
		}

		filehandler.FileMustExist(epubFile, "epub-file")

		logger.WriteInfo("Started showing manually fixable issues...\n")

		var epubFolder = filehandler.GetFileFolder(epubFile)
		var dest = filehandler.JoinPath(epubFolder, "epub")
		filehandler.UnzipRunOperationAndRezip(epubFile, dest, func() {
			opfFolder, epubInfo := getEpubInfo(dest, epubFile)
			validateFilesExist(opfFolder, epubInfo.HtmlFiles)
			validateFilesExist(opfFolder, epubInfo.CssFiles)

			var addCssSectionIfMissing bool = false
			var addCssPageIfMissing bool = false
			var contextBreak string
			if runAll || runSectionBreak {
				contextBreak = logger.GetInputString("What is the section break for the epub?:")

				if strings.TrimSpace(contextBreak) == "" {
					logger.WriteError("Please provide a non-whitespace section break")
				}

				/**
				TODO: handle the scenario where the section break is an image

				Image Context Breaks
				To use an image:

				In the CSS:
				hr.image {
				display:block;
				background: transparent url("images/sectionBreakImage.png") no-repeat center;
				height:2em;
				border:0;
				}

				In the HTML:
				<hr class="image" />
				**/
			}

			var cssFiles = make([]string, len(epubInfo.CssFiles))
			var i = 0
			for cssFile := range epubInfo.CssFiles {
				cssFiles[i] = cssFile
				i++
			}

			if (runAll || runSectionBreak || runPageBreak) && len(cssFiles) == 0 {
				logger.WriteError(CssPathsEmptyWhenArgIsNeeded)
			}

			var saveAndQuit = false
			for file := range epubInfo.HtmlFiles {
				if saveAndQuit {
					break
				}

				var filePath = getFilePath(opfFolder, file)
				fileText := filehandler.ReadInFileContents(filePath)

				var newText = linter.CleanupHtmlSpacing(fileText)
				if runAll || runBrokenLines {
					var brokenLineFixSuggestions = linter.GetPotentiallyBrokenLines(newText)
					newText, _, saveAndQuit = promptAboutSuggestions("Potential Broken Lines", brokenLineFixSuggestions, newText, false)
				}

				if (runAll || runSectionBreak) && !saveAndQuit {
					var contextBreakSuggestions = linter.GetPotentialSectionBreaks(newText, contextBreak)

					var contextBreakUpdated bool
					newText, contextBreakUpdated, saveAndQuit = promptAboutSuggestions("Potential Section Breaks", contextBreakSuggestions, newText, true)
					addCssSectionIfMissing = addCssSectionIfMissing || contextBreakUpdated
				}

				if (runAll || runPageBreak) && !saveAndQuit {
					var pageBreakSuggestions = linter.GetPotentialPageBreaks(newText)

					var pageBreakUpdated bool
					newText, pageBreakUpdated, saveAndQuit = promptAboutSuggestions("Potential Page Breaks", pageBreakSuggestions, newText, true)
					addCssPageIfMissing = addCssPageIfMissing || pageBreakUpdated
				}

				if (runAll || runOxfordCommas) && !saveAndQuit {
					var oxfordCommaSuggestions = linter.GetPotentialMissingOxfordCommas(newText)
					newText, _, saveAndQuit = promptAboutSuggestions("Potential Missing Oxford Commas", oxfordCommaSuggestions, newText, false)
				}

				if (runAll || runAlthoughBut) && !saveAndQuit {
					var althoughButSuggestions = linter.GetPotentialAlthoughButInstances(newText)
					newText, _, saveAndQuit = promptAboutSuggestions("Potential Although But Instances", althoughButSuggestions, newText, false)
				}

				if (runAll || runThoughts) && !saveAndQuit {
					var thoughtSuggestions = linter.GetPotentialThoughtInstances(newText)
					newText, _, saveAndQuit = promptAboutSuggestions("Potential Thought Instances", thoughtSuggestions, newText, false)
				}

				if runAll || runConversation && !saveAndQuit {
					var conversationSuggestions = linter.GetPotentialSquareBracketConversationInstances(newText)
					newText, _, saveAndQuit = promptAboutSuggestions("Potential Conversation Instances", conversationSuggestions, newText, false)
				}

				if runAll || runNecessaryWords && !saveAndQuit {
					var necessaryWordSuggestions = linter.GetPotentialSquareBracketNecessaryWords(newText)
					newText, _, saveAndQuit = promptAboutSuggestions("Potential Necessary Word Omission Instances", necessaryWordSuggestions, newText, false)
				}

				if fileText == newText {
					continue
				}

				filehandler.WriteFileContents(filePath, newText)
			}

			handleCssChanges(addCssSectionIfMissing, addCssPageIfMissing, opfFolder, cssFiles, contextBreak)
		})

		logger.WriteInfo("\nFinished showing manually fixable issues...")
	},
}

func init() {
	EpubCmd.AddCommand(fixableCmd)

	fixableCmd.Flags().BoolVarP(&runAll, "run-all", "a", false, "whether to run all of the fixable suggestions")
	fixableCmd.Flags().BoolVarP(&runBrokenLines, "run-broken-lines", "b", false, "whether to run the logic for getting broken line suggestions")
	fixableCmd.Flags().BoolVarP(&runSectionBreak, "run-section-breaks", "s", false, "whether to run the logic for getting section break suggestions (must be used with css-paths)")
	fixableCmd.Flags().BoolVarP(&runPageBreak, "run-page-breaks", "p", false, "whether to run the logic for getting page break suggestions (must be used with css-paths)")
	fixableCmd.Flags().BoolVarP(&runOxfordCommas, "run-oxford-commas", "o", false, "whether to run the logic for getting oxford comma suggestions")
	fixableCmd.Flags().BoolVarP(&runAlthoughBut, "run-although-but", "n", false, "whether to run the logic for getting although but suggestions")
	fixableCmd.Flags().BoolVarP(&runThoughts, "run-thoughts", "t", false, "whether to run the logic for getting thought suggestions (words in parentheses may be instances of a person's thoughts)")
	fixableCmd.Flags().BoolVarP(&runConversation, "run-conversation", "c", false, "whether to run the logic for getting conversation suggestions (paragraphs in square brackets may be instances of a conversation)")
	fixableCmd.Flags().BoolVarP(&runNecessaryWords, "run-necessary-words", "w", false, "whether to run the logic for getting necessary word suggestions (words that are a subset of paragraph content are in square brackets may be instances of necessary words for a sentence)")
	fixableCmd.Flags().StringVarP(&epubFile, "epub-file", "f", "", "the epub file to find manually fixable issues in")
	err := fixableCmd.MarkFlagRequired("epub-file")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "epub-file" as required on fixable command: %v`, err))
	}

	err = fixableCmd.MarkFlagFilename("epub-file", "epub")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "epub-file" as looking for specific file types on fixable command: %v`, err))
	}
}

func ValidateManuallyFixableFlags(epubPath string, runAll, runBrokenLines, runSectionBreak, runPageBreak, runOxfordCommas, runAlthoughBut, runThoughts, runConversation, runNecessaryWords bool) error {
	err := validateCommonEpubFlags(epubPath)
	if err != nil {
		return err
	}

	if !runAll && !runBrokenLines && !runSectionBreak && !runPageBreak && !runOxfordCommas && !runAlthoughBut && !runConversation && !runThoughts && !runNecessaryWords {
		return errors.New(OneRunBoolArgMustBeEnabled)
	}

	return nil
}

func promptAboutSuggestions(suggestionsTitle string, suggestions map[string]string, fileText string, replaceAllInstances bool) (string, bool, bool) {
	var valueReplaced = false
	var newText = fileText

	if len(suggestions) == 0 {
		return newText, valueReplaced, false
	}

	// replace count was added to make sure that if we have a case where the original and suggested value
	// may exist more than once in a file we want to go ahead and replace all instances of the original
	// with the suggested
	var replaceCount = 1
	if replaceAllInstances {
		replaceCount = -1
	}

	logger.WriteInfo(cliLineSeparator)
	logger.WriteInfo(suggestionsTitle + ":")
	logger.WriteInfo(cliLineSeparator + "\n")

	for original, suggestion := range suggestions {
		resp := logger.GetInputString(fmt.Sprintf("Would you like to make the following update \"%s\"? (Y/N): ", stringdiff.GetPrettyDiffString(strings.TrimLeft(original, "\n"), strings.TrimLeft(suggestion, "\n"))))
		if strings.EqualFold(resp, "Y") {
			newText = strings.Replace(newText, original, suggestion, replaceCount)
			valueReplaced = true
		} else if strings.EqualFold(resp, "Q") {
			return newText, valueReplaced, true
		}

		logger.WriteInfo("")
	}

	return newText, valueReplaced, false
}

func handleCssChanges(addCssSectionIfMissing, addCssPageIfMissing bool, opfFolder string, cssFiles []string, contextBreak string) {
	if !addCssSectionIfMissing && !addCssPageIfMissing {
		return
	}

	var cssSelectionPrompt = "Please enter the number of the css file to append the css to:\n"
	for i, file := range cssFiles {
		cssSelectionPrompt += fmt.Sprintf("%d. %s\n", i, file)
	}

	var selectedCssFileIndex = logger.GetInputInt(cssSelectionPrompt)
	if selectedCssFileIndex < 0 || selectedCssFileIndex >= len(cssFiles) {
		logger.WriteError(fmt.Sprintf("Please select a valid css file value instead of \"%d\".", selectedCssFileIndex))
	}

	var cssFile = cssFiles[selectedCssFileIndex]
	var cssFilePath = filehandler.JoinPath(opfFolder, cssFile)
	var css = filehandler.ReadInFileContents(cssFilePath)
	var newCssText = css

	if addCssSectionIfMissing {
		newCssText = linter.AddCssSectionBreakIfMissing(newCssText, contextBreak)
	}

	if addCssPageIfMissing {
		newCssText = linter.AddCssPageBreakIfMissing(newCssText)
	}

	if newCssText != css {
		filehandler.WriteFileContents(cssFilePath, newCssText)
	}
}
