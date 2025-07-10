package epub

import (
	"archive/zip"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/cmd"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
	"github.com/spf13/cobra"
)

type potentiallyFixableIssue struct {
	name                        string
	getSuggestions              func(string) map[string]string
	isEnabled                   *bool
	updateAllInstances          bool
	addCssSectionBreakIfMissing bool
	addCssPageBreakIfMissing    bool
}

var (
	// this is declared globally here just for use in manuallyFixableIssue to make sure that the struct definition
	// is satisfied even though this value is the second param for potential section breaks
	contextBreak             string
	runAll                   bool
	runBrokenLines           bool
	runSectionBreak          bool
	runPageBreak             bool
	runOxfordCommas          bool
	runLackingClause         bool
	runThoughts              bool
	runConversation          bool
	runNecessaryWords        bool
	runSingleQuotes          bool
	useTui                   bool
	logFile                  string
	potentiallyFixableIssues = []potentiallyFixableIssue{
		{
			name:           "Potential Conversation Instances",
			getSuggestions: linter.GetPotentialSquareBracketConversationInstances,
			isEnabled:      &runConversation,
		},
		{
			name:           "Potential Necessary Word Omission Instances",
			getSuggestions: linter.GetPotentialSquareBracketNecessaryWords,
			isEnabled:      &runNecessaryWords,
		},
		{
			name:           "Potential Broken Lines",
			getSuggestions: linter.GetPotentiallyBrokenLines,
			isEnabled:      &runBrokenLines,
		},
		{
			name:           "Potential Incorrect Single Quotes",
			getSuggestions: linter.GetPotentialIncorrectSingleQuotes,
			isEnabled:      &runSingleQuotes,
		},
		{
			name: "Potential Section Breaks",
			// wrapper here allows calling the get potential section breaks logic without needing to change the function definition
			getSuggestions: func(text string) map[string]string {
				return linter.GetPotentialSectionBreaks(text, contextBreak)
			},
			isEnabled:                   &runSectionBreak,
			updateAllInstances:          true,
			addCssSectionBreakIfMissing: true,
		},
		{
			name:                     "Potential Page Breaks",
			getSuggestions:           linter.GetPotentialPageBreaks,
			isEnabled:                &runPageBreak,
			updateAllInstances:       true,
			addCssPageBreakIfMissing: true,
		},
		{
			name:           "Potential Missing Oxford Commas",
			getSuggestions: linter.GetPotentialMissingOxfordCommas,
			isEnabled:      &runOxfordCommas,
		},
		{
			name:           "Potentially Lacking Subordinate Clause Instances",
			getSuggestions: linter.GetPotentiallyLackingSubordinateClauseInstances,
			isEnabled:      &runLackingClause,
		},
		{
			name:           "Potential Thought Instances",
			getSuggestions: linter.GetPotentialThoughtInstances,
			isEnabled:      &runThoughts,
		},
	}
	// errors
	ErrOneRunBoolArgMustBeEnabled = errors.New("at least one rule to run must be enabled")
	ErrNoCssFiles                 = errors.New("the epub must have at least 1 css file in order to handle section or page breaks")
)

// fixableCmd represents the fixable command
var fixableCmd = &cobra.Command{
	Use:   "fixable",
	Short: "Runs the specified fixable actions that require manual input to determine what to do.",
	Example: heredoc.Doc(`To run all of the possible potential fixes:
	epub-lint fixable -f test.epub -a
	Note: this will require a css file to already exist in the epub
	
	To just fix broken paragraph endings:
	epub-lint fixable -f test.epub --broken-lines

	To just update section breaks:
	epub-lint fixable -f test.epub --section-breaks
	Note: this will require a css file to already exist in the epub

	To just update page breaks:
	epub-lint fixable -f test.epub --page-breaks
	Note: this will require a css file to already exist in the epub

	To just fix missing oxford commas:
	epub-lint fixable -f test.epub --oxford-commas

	To just fix potentially lacking subordinate clause instances:
	epub-lint fixable -f test.epub --lacking-subordinate-clause

	To just fix instances of thoughts in parentheses:
	epub-lint fixable -f test.epub --thoughts

	To run a combination of options:
	epub-lint fixable -f test.epub -oxford-commas --thoughts --necessary-words
	`),
	Long: heredoc.Doc(`Goes through all of the content files and runs the specified fixable actions on them asking
	for user input on each value found that matches the potential fix criteria.
	Potential things that can be fixed:
	- Broken paragraph endings
	- Section breaks being hardcoded instead of an hr
	- Page breaks being hardcoded instead of an hr
	- Oxford commas being missing before or's or and's
	- Possible instances of sentences with two subordinate clauses (i.e. have although..., but)
	- Possible instances of thoughts that are in parentheses
	- Possible instances of conversation encapsulated in square brackets
	- Possible instances of words in square brackets that may be necessary for the sentence (i.e. need to have the brackets removed)
	- Possible instances of single quotes that should actually be double quotes (i.e. when a word is in single quotes, but is not inside of double quotes)
	`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateManuallyFixableFlags(epubFile, runAll, runBrokenLines, runSectionBreak, runPageBreak, runOxfordCommas, runLackingClause, runThoughts, runConversation, runNecessaryWords, runSingleQuotes)
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FileArgExists(epubFile, "file")
		if err != nil {
			logger.WriteError(err.Error())
		}

		if useTui {
			err = runTuiEpubFixable()
		} else {
			err = runCliEpubFixable()

			if err != nil {
				logger.WriteInfo("\nFinished showing manually fixable issues...")
			}
		}

		if err != nil {
			logger.WriteErrorf("failed to fix manually fixable issues for %q: %s", epubFile, err)
		}
	},
}

func init() {
	cmd.EpubCmd.AddCommand(fixableCmd)

	fixableCmd.Flags().BoolVarP(&runAll, "all", "a", false, "whether to run all of the fixable suggestions")
	fixableCmd.Flags().BoolVarP(&runBrokenLines, "broken-lines", "", false, "whether to run the logic for getting broken line suggestions")
	fixableCmd.Flags().BoolVarP(&runSectionBreak, "section-breaks", "", false, "whether to run the logic for getting section break suggestions (must be used with an epub with a css file)")
	fixableCmd.Flags().BoolVarP(&runPageBreak, "page-breaks", "", false, "whether to run the logic for getting page break suggestions (must be used with an epub with a css file)")
	fixableCmd.Flags().BoolVarP(&runOxfordCommas, "oxford-commas", "", false, "whether to run the logic for getting oxford comma suggestions")
	fixableCmd.Flags().BoolVarP(&runLackingClause, "lacking-subordinate-clause", "", false, "whether to run the logic for getting potentially lacking subordinate clause suggestions")
	fixableCmd.Flags().BoolVarP(&runThoughts, "thoughts", "", false, "whether to run the logic for getting thought suggestions (words in parentheses may be instances of a person's thoughts)")
	fixableCmd.Flags().BoolVarP(&runConversation, "conversation", "", false, "whether to run the logic for getting conversation suggestions (paragraphs in square brackets may be instances of a conversation)")
	fixableCmd.Flags().BoolVarP(&runNecessaryWords, "necessary-words", "", false, "whether to run the logic for getting necessary word suggestions (words that are a subset of paragraph content are in square brackets may be instances of necessary words for a sentence)")
	fixableCmd.Flags().BoolVarP(&runSingleQuotes, "single-quotes", "", false, "whether to run the logic for getting incorrect single quote suggestions")
	fixableCmd.Flags().BoolVarP(&useTui, "use-tui", "t", false, "whether to use the terminal UI for suggesting fixes")
	fixableCmd.Flags().StringVarP(&logFile, "log-file", "", "", "the place to write debug logs to when using the TUI")
	fixableCmd.Flags().StringVarP(&epubFile, "file", "f", "", "the epub file to find manually fixable issues in")
	err := fixableCmd.MarkFlagRequired("file")
	if err != nil {
		logger.WriteErrorf(`failed to mark flag "file" as required on fixable command: %v\n`, err)
	}

	err = fixableCmd.MarkFlagFilename("file", "epub")
	if err != nil {
		logger.WriteErrorf(`failed to mark flag "file" as looking for specific file types on fixable command: %v\n`, err)
	}
}

func ValidateManuallyFixableFlags(epubPath string, runAll, runBrokenLines, runSectionBreak, runPageBreak, runOxfordCommas, runAlthoughBut, runThoughts, runConversation, runNecessaryWords, runSingleQuotes bool) error {
	err := validateCommonEpubFlags(epubPath)
	if err != nil {
		return err
	}

	if !runAll && !runBrokenLines && !runSectionBreak && !runPageBreak && !runOxfordCommas && !runAlthoughBut && !runConversation && !runThoughts && !runNecessaryWords && !runSingleQuotes {
		return ErrOneRunBoolArgMustBeEnabled
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
		diffString, err := stringdiff.GetPrettyDiffString(strings.TrimLeft(original, "\n"), strings.TrimLeft(suggestion, "\n"))
		if err != nil {
			logger.WriteError(err.Error())
		}

		// Warning: do not use %q on the following line as it will get rid of the color coding of changes in the terminal
		resp := logger.GetInputString(fmt.Sprintf("Would you like to make the following update \"%s\"? (Y/N/Q): ", diffString))
		switch strings.ToLower(resp) {
		case "y":
			newText = strings.Replace(newText, original, suggestion, replaceCount)
			valueReplaced = true
		case "q":
			return newText, valueReplaced, true
		}

		logger.WriteInfo("")
	}

	return newText, valueReplaced, false
}

func handleCssChangesCli(addCssSectionIfMissing, addCssPageIfMissing bool, opfFolder string, cssFiles []string, contextBreak string, zipFiles map[string]*zip.File, w *zip.Writer, handledFiles []string) ([]string, error) {
	if !addCssSectionIfMissing && !addCssPageIfMissing {
		return handledFiles, nil
	}

	var cssSelectionPrompt = "Please enter the number of the css file to append the css to:\n"
	for i, file := range cssFiles {
		cssSelectionPrompt += fmt.Sprintf("%d. %s\n", i, file)
	}

	var selectedCssFileIndex = logger.GetInputInt(cssSelectionPrompt)
	if selectedCssFileIndex < 0 || selectedCssFileIndex >= len(cssFiles) {
		return nil, fmt.Errorf("please select a valid css file value instead of \"%d\".\n", selectedCssFileIndex)
	}

	return updateCssFile(addCssSectionIfMissing, addCssPageIfMissing, filehandler.JoinPath(opfFolder, cssFiles[selectedCssFileIndex]), contextBreak, zipFiles, w, handledFiles)
}

func handleCssChangesTui(addCssSectionIfMissing, addCssPageIfMissing bool, opfFolder, selectedCssFile, contextBreak string, zipFiles map[string]*zip.File, w *zip.Writer, handledFiles []string) ([]string, error) {
	if !addCssSectionIfMissing && !addCssPageIfMissing {
		return handledFiles, nil
	}

	if strings.TrimSpace(selectedCssFile) == "" {
		return nil, fmt.Errorf("please select a valid css file instead of %q.\n", selectedCssFile)
	}

	return updateCssFile(addCssSectionIfMissing, addCssPageIfMissing, filehandler.JoinPath(opfFolder, selectedCssFile), contextBreak, zipFiles, w, handledFiles)
}

func updateCssFile(addCssSectionIfMissing, addCssPageIfMissing bool, selectedCssFile, contextBreak string, zipFiles map[string]*zip.File, w *zip.Writer, handledFiles []string) ([]string, error) {
	zipFile := zipFiles[selectedCssFile]
	css, err := filehandler.ReadInZipFileContents(zipFile)
	if err != nil {
		return nil, err
	}

	var newCssText = css
	if addCssSectionIfMissing {
		newCssText = linter.AddCssSectionBreakIfMissing(newCssText, contextBreak)
	}

	if addCssPageIfMissing {
		newCssText = linter.AddCssPageBreakIfMissing(newCssText)
	}

	if newCssText == css {
		return handledFiles, nil
	}

	err = filehandler.WriteZipCompressedString(w, selectedCssFile, newCssText)
	if err != nil {
		return nil, err
	}

	return append(handledFiles, selectedCssFile), nil
}

func runTuiEpubFixable() error {
	return epubhandler.UpdateEpub(epubFile, func(zipFiles map[string]*zip.File, w *zip.Writer, epubInfo epubhandler.EpubInfo, opfFolder string) ([]string, error) {
		err := validateFilesExist(opfFolder, epubInfo.HtmlFiles, zipFiles)
		if err != nil {
			return nil, err
		}

		err = validateFilesExist(opfFolder, epubInfo.CssFiles, zipFiles)
		if err != nil {
			return nil, err
		}

		var cssFiles = make([]string, 0, len(epubInfo.CssFiles))
		for cssFile := range epubInfo.CssFiles {
			cssFiles = append(cssFiles, cssFile)
		}

		if (runAll || runSectionBreak || runPageBreak) && len(cssFiles) == 0 {
			return nil, ErrNoCssFiles
		}

		var file *os.File
		if strings.TrimSpace(logFile) != "" {
			file, err = tea.LogToFile(logFile, "debug")
			if err != nil {
				return nil, fmt.Errorf("failed to create TUI log file %q: %w", logFile, err)
			}

			defer file.Close()
		}

		var (
			initialModel = newModel(runAll, runSectionBreak, potentiallyFixableIssues, cssFiles, file)
			i            = 0
		)
		initialModel.potentiallyFixableIssuesInfo.filePaths = make([]string, len(epubInfo.HtmlFiles))

		// Collect file contents
		for file := range epubInfo.HtmlFiles {
			var filePath = getFilePath(opfFolder, file)
			zipFile := zipFiles[filePath]

			fileText, err := filehandler.ReadInZipFileContents(zipFile)
			if err != nil {
				return nil, err
			}

			initialModel.potentiallyFixableIssuesInfo.fileTexts[filePath] = linter.CleanupHtmlSpacing(fileText)
			initialModel.potentiallyFixableIssuesInfo.filePaths[i] = filePath
			i++
		}

		p := tea.NewProgram(&initialModel, tea.WithAltScreen())
		finalModel, err := p.Run()
		if err != nil {
			return nil, err
		}

		model := finalModel.(fixableIssuesModel)
		if model.Err != nil {
			if errors.Is(model.Err, errUserKilledProgram) {
				logger.WriteInfo("Quitting. User exited the program...")
				os.Exit(0)
			}

			return nil, model.Err
		}

		var handledFiles []string
		// Process and write updated files
		for filePath, fileText := range model.potentiallyFixableIssuesInfo.fileTexts {
			err = filehandler.WriteZipCompressedString(w, filePath, fileText)
			if err != nil {
				return nil, err
			}
			handledFiles = append(handledFiles, filePath)
		}

		return handleCssChangesTui(model.potentiallyFixableIssuesInfo.addCssSectionBreakIfMissing, model.potentiallyFixableIssuesInfo.addCssPageBreakIfMissing, opfFolder, model.cssSelectionInfo.selectedCssFile, contextBreak, zipFiles, w, handledFiles)
	})
}

func runCliEpubFixable() error {
	logger.WriteInfo("Started showing manually fixable issues...\n")

	return epubhandler.UpdateEpub(epubFile, func(zipFiles map[string]*zip.File, w *zip.Writer, epubInfo epubhandler.EpubInfo, opfFolder string) ([]string, error) {
		err := validateFilesExist(opfFolder, epubInfo.HtmlFiles, zipFiles)
		if err != nil {
			return nil, err
		}

		err = validateFilesExist(opfFolder, epubInfo.CssFiles, zipFiles)
		if err != nil {
			return nil, err
		}

		var (
			handledFiles                                []string
			addCssSectionIfMissing, addCssPageIfMissing bool
		)
		if runAll || runSectionBreak {
			contextBreak = logger.GetInputString("What is the section break for the epub?:")

			if strings.TrimSpace(contextBreak) == "" {
				return nil, fmt.Errorf("please provide a non-whitespace section break")
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
			return nil, ErrNoCssFiles
		}

		var saveAndQuit = false
		for file := range epubInfo.HtmlFiles {
			if saveAndQuit {
				break
			}

			var filePath = getFilePath(opfFolder, file)
			zipFile := zipFiles[filePath]

			fileText, err := filehandler.ReadInZipFileContents(zipFile)
			if err != nil {
				return nil, err
			}

			var newText = linter.CleanupHtmlSpacing(fileText)

			for _, potentiallyFixableIssue := range potentiallyFixableIssues {
				if saveAndQuit {
					break
				}

				if potentiallyFixableIssue.isEnabled == nil {
					return nil, fmt.Errorf("%q is not properly setup to run as a potentially fixable rule since it has no boolean for isEnabled", potentiallyFixableIssue.name)
				}

				if runAll || *potentiallyFixableIssue.isEnabled {
					suggestions := potentiallyFixableIssue.getSuggestions(newText)

					var updateMade bool
					newText, updateMade, saveAndQuit = promptAboutSuggestions(potentiallyFixableIssue.name, suggestions, newText, potentiallyFixableIssue.updateAllInstances)

					if potentiallyFixableIssue.addCssSectionBreakIfMissing && updateMade {
						addCssSectionIfMissing = addCssSectionIfMissing || updateMade
					}

					if potentiallyFixableIssue.addCssPageBreakIfMissing && updateMade {
						addCssPageIfMissing = addCssPageIfMissing || updateMade
					}
				}
			}

			err = filehandler.WriteZipCompressedString(w, filePath, newText)
			if err != nil {
				return nil, err
			}

			handledFiles = append(handledFiles, filePath)
		}

		return handleCssChangesCli(addCssSectionIfMissing, addCssPageIfMissing, opfFolder, cssFiles, contextBreak, zipFiles, w, handledFiles)
	})
}
