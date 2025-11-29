package cmd

import (
	"archive/zip"
	"errors"

	"github.com/MakeNowJust/heredoc"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue/fixer"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

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
	potentiallyFixableIssues = []potentiallyfixableissue.PotentiallyFixableIssue{
		{
			Name:           "Potential Conversation Instances",
			GetSuggestions: potentiallyfixableissue.GetPotentialSquareBracketConversationInstances,
			IsEnabled:      &runConversation,
		},
		{
			Name:           "Potential Necessary Word Omission Instances",
			GetSuggestions: potentiallyfixableissue.GetPotentialSquareBracketNecessaryWords,
			IsEnabled:      &runNecessaryWords,
		},
		{
			Name:           "Potential Broken Lines",
			GetSuggestions: potentiallyfixableissue.GetPotentiallyBrokenLines,
			IsEnabled:      &runBrokenLines,
		},
		{
			Name:           "Potential Incorrect Single Quotes",
			GetSuggestions: potentiallyfixableissue.GetPotentialIncorrectSingleQuotes,
			IsEnabled:      &runSingleQuotes,
		},
		{
			Name: "Potential Section Breaks",
			// wrapper here allows calling the get potential section breaks logic without needing to change the function definition
			GetSuggestions: func(text string) (map[string]string, error) {
				return potentiallyfixableissue.GetPotentialSectionBreaks(text, contextBreak)
			},
			IsEnabled:                   &runSectionBreak,
			UpdateAllInstances:          true,
			AddCssSectionBreakIfMissing: true,
		},
		{
			Name:                     "Potential Page Breaks",
			GetSuggestions:           potentiallyfixableissue.GetPotentialPageBreaks,
			IsEnabled:                &runPageBreak,
			UpdateAllInstances:       true,
			AddCssPageBreakIfMissing: true,
		},
		{
			Name:           "Potential Missing Oxford Commas",
			GetSuggestions: potentiallyfixableissue.GetPotentialMissingOxfordCommas,
			IsEnabled:      &runOxfordCommas,
		},
		{
			Name:           "Potentially Lacking Subordinate Clause Instances",
			GetSuggestions: potentiallyfixableissue.GetPotentiallyLackingSubordinateClauseInstances,
			IsEnabled:      &runLackingClause,
		},
		{
			Name:           "Potential Thought Instances",
			GetSuggestions: potentiallyfixableissue.GetPotentialThoughtInstances,
			IsEnabled:      &runThoughts,
		},
	}
	// errors
	ErrOneRunBoolArgMustBeEnabled = errors.New("at least one rule to run must be enabled")
	ErrNoCssFiles                 = errors.New("the epub must have at least 1 css file in order to handle section or page breaks")
)

// contentCmd represents the fix content command
var contentCmd = &cobra.Command{
	Use:   "content",
	Short: "Runs the specified fixable actions that require manual input to determine what to do.",
	Example: heredoc.Doc(`To run all of the possible potential fixes:
	epub-lint fix content -f test.epub -a
	Note: this will require a css file to already exist in the epub
	
	To just fix broken paragraph endings:
	epub-lint fix content -f test.epub --broken-lines

	To just update section breaks:
	epub-lint fix content -f test.epub --section-breaks
	Note: this will require a css file to already exist in the epub

	To just update page breaks:
	epub-lint fix content -f test.epub --page-breaks
	Note: this will require a css file to already exist in the epub

	To just fix missing oxford commas:
	epub-lint fix content -f test.epub --oxford-commas

	To just fix potentially lacking subordinate clause instances:
	epub-lint fix content -f test.epub --lacking-subordinate-clause

	To just fix instances of thoughts in parentheses:
	epub-lint fix content -f test.epub --thoughts

	To run a combination of options:
	epub-lint fix content -f test.epub -oxford-commas --thoughts --necessary-words
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

		var handler fixer.Fixer
		if useTui {
			handler = &fixer.TuiFixer{}
		} else {
			handler = &fixer.CliFixer{}
		}

		var initialLog = handler.InitialLog()
		if initialLog != "" {
			logger.WriteInfo(initialLog)
		}

		err = epubhandler.UpdateEpub(epubFile, func(zipFiles map[string]*zip.File, w *zip.Writer, epubInfo epubhandler.EpubInfo, opfFolder string) ([]string, error) {
			err = validateFilesExist(opfFolder, epubInfo.HtmlFiles, zipFiles)
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

			var skipCss = runAll && len(cssFiles) == 0
			if (runSectionBreak || runPageBreak) && len(cssFiles) == 0 {
				return nil, ErrNoCssFiles
			}

			handler.Init(&epubInfo, runAll, skipCss, runSectionBreak, potentiallyFixableIssues, cssFiles, logFile, opfFolder, &contextBreak, func(fileName string) (string, error) {
				zipFile := zipFiles[fileName]

				fileText, err := filehandler.ReadInZipFileContents(zipFile)
				if err != nil {
					return "", err
				}

				return fileText, nil
			}, func(fileName, content string) error {
				err = filehandler.WriteZipCompressedString(w, fileName, content)
				if err != nil {
					return err
				}

				return nil
			})

			err = handler.Setup()
			if err != nil {
				return nil, err
			}

			err = handler.Run()
			if err != nil {
				return nil, err
			}

			return handler.HandleCss()
		})

		if err != nil {
			successLog := handler.SuccessfulLog()

			if successLog != "" {
				logger.WriteInfo(successLog)
			}
		}

		if err != nil {
			logger.WriteErrorf("failed to fix manually fixable content issues for %q: %s", epubFile, err)
		}
	},
}

func init() {
	fixCmd.AddCommand(contentCmd)

	contentCmd.Flags().BoolVarP(&runAll, "all", "a", false, "whether to run all of the fixable suggestions")
	contentCmd.Flags().BoolVarP(&runBrokenLines, "broken-lines", "", false, "whether to run the logic for getting broken line suggestions")
	contentCmd.Flags().BoolVarP(&runSectionBreak, "section-breaks", "", false, "whether to run the logic for getting section break suggestions (must be used with an epub with a css file)")
	contentCmd.Flags().BoolVarP(&runPageBreak, "page-breaks", "", false, "whether to run the logic for getting page break suggestions (must be used with an epub with a css file)")
	contentCmd.Flags().BoolVarP(&runOxfordCommas, "oxford-commas", "", false, "whether to run the logic for getting oxford comma suggestions")
	contentCmd.Flags().BoolVarP(&runLackingClause, "lacking-subordinate-clause", "", false, "whether to run the logic for getting potentially lacking subordinate clause suggestions")
	contentCmd.Flags().BoolVarP(&runThoughts, "thoughts", "", false, "whether to run the logic for getting thought suggestions (words in parentheses may be instances of a person's thoughts)")
	contentCmd.Flags().BoolVarP(&runConversation, "conversation", "", false, "whether to run the logic for getting conversation suggestions (paragraphs in square brackets may be instances of a conversation)")
	contentCmd.Flags().BoolVarP(&runNecessaryWords, "necessary-words", "", false, "whether to run the logic for getting necessary word suggestions (words that are a subset of paragraph content are in square brackets may be instances of necessary words for a sentence)")
	contentCmd.Flags().BoolVarP(&runSingleQuotes, "single-quotes", "", false, "whether to run the logic for getting incorrect single quote suggestions")
	contentCmd.Flags().BoolVarP(&useTui, "use-tui", "t", false, "whether to use the terminal UI for suggesting fixes")
	contentCmd.Flags().StringVarP(&logFile, "log-file", "", "", "the place to write debug logs to when using the TUI")
	contentCmd.Flags().StringVarP(&epubFile, "file", "f", "", "the epub file to find manually fixable issues in")
	err := contentCmd.MarkFlagRequired("file")
	if err != nil {
		logger.WriteErrorf(`failed to mark flag "file" as required on fix content command: %v\n`, err)
	}

	err = contentCmd.MarkFlagFilename("file", "epub")
	if err != nil {
		logger.WriteErrorf(`failed to mark flag "file" as looking for specific file types on fix content command: %v\n`, err)
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
