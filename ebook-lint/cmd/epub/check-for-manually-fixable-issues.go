package epub

import (
	"archive/zip"
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	tea "github.com/charmbracelet/bubbletea"
	epubhandler "github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
	"github.com/spf13/cobra"
)

type potentiallyFixableIssue struct {
	name               string
	getSuggestions     func(string) map[string]string
	isEnabled          *bool
	updateAllInstances bool
	addCssIfMissing    bool
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
	runAlthoughBut           bool
	runThoughts              bool
	runConversation          bool
	runNecessaryWords        bool
	useTui                   bool
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
			name: "Potential Section Breaks",
			// wrapper here allows calling the get potential section breaks logic without needing to change the function definition
			getSuggestions: func(text string) map[string]string {
				return linter.GetPotentialSectionBreaks(text, contextBreak)
			},
			isEnabled:          &runSectionBreak,
			updateAllInstances: true,
			addCssIfMissing:    true,
		},
		{
			name:               "Potential Page Breaks",
			getSuggestions:     linter.GetPotentialPageBreaks,
			isEnabled:          &runPageBreak,
			updateAllInstances: true,
			addCssIfMissing:    true,
		},
		{
			name:           "Potential Missing Oxford Commas",
			getSuggestions: linter.GetPotentialMissingOxfordCommas,
			isEnabled:      &runOxfordCommas,
		},
		{
			name:           "Potential Although But Instances",
			getSuggestions: linter.GetPotentialAlthoughButInstances,
			isEnabled:      &runAlthoughBut,
		},
		{
			name:           "Potential Thought Instances",
			getSuggestions: linter.GetPotentialThoughtInstances,
			isEnabled:      &runThoughts,
		},
	}
	// errors
	ErrOneRunBoolArgMustBeEnabled   = errors.New("either run-all, run-broken-lines, run-section-breaks, run-page-breaks, run-oxford-commas, or run-although-but must be specified")
	ErrCssPathsEmptyWhenArgIsNeeded = errors.New("css-paths must have a value when including handling section or page breaks")
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

		err = filehandler.FileArgExists(epubFile, "epub-file")
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
	fixableCmd.Flags().BoolVarP(&useTui, "use-tui", "u", false, "whether to use the terminal UI for suggesting fixes")
	fixableCmd.Flags().StringVarP(&epubFile, "epub-file", "f", "", "the epub file to find manually fixable issues in")
	err := fixableCmd.MarkFlagRequired("epub-file")
	if err != nil {
		logger.WriteErrorf(`failed to mark flag "epub-file" as required on fixable command: %v\n`, err)
	}

	err = fixableCmd.MarkFlagFilename("epub-file", "epub")
	if err != nil {
		logger.WriteErrorf(`failed to mark flag "epub-file" as looking for specific file types on fixable command: %v\n`, err)
	}
}

func ValidateManuallyFixableFlags(epubPath string, runAll, runBrokenLines, runSectionBreak, runPageBreak, runOxfordCommas, runAlthoughBut, runThoughts, runConversation, runNecessaryWords bool) error {
	err := validateCommonEpubFlags(epubPath)
	if err != nil {
		return err
	}

	if !runAll && !runBrokenLines && !runSectionBreak && !runPageBreak && !runOxfordCommas && !runAlthoughBut && !runConversation && !runThoughts && !runNecessaryWords {
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

func handleCssChanges(addCssSectionIfMissing, addCssPageIfMissing bool, opfFolder string, cssFiles []string, contextBreak string, zipFiles map[string]*zip.File, w *zip.Writer, handledFiles []string) ([]string, error) {
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

	var cssFile = cssFiles[selectedCssFileIndex]
	var cssFilePath = filehandler.JoinPath(opfFolder, cssFile)
	zipFile := zipFiles[cssFilePath]
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

	err = filehandler.WriteZipCompressedString(w, cssFilePath, newCssText)
	if err != nil {
		return nil, err
	}

	return append(handledFiles, cssFilePath), nil
}

// type fixableTuiModel struct {
// 	stage                    fixableStage
// 	currentSuggestionGroup   string
// 	suggestions              map[string]string
// 	originalText             string
// 	currentSuggestionKey     string
// 	editMode                 bool
// 	editedText               string
// 	contextBreakInput        string
// 	currentFile              string
// 	cssFiles                 []string
// 	selectedCssFileIndex     int
// 	potentiallyFixableIssues []potentiallyFixableIssue
// 	currentIssueIndex        int
// 	runAll                   bool
// 	fileTexts                map[string]string
// 	handledFiles             []string
// 	saveAndQuit              bool
// 	currentFileIndex         int
// 	fileNames                []string
// }

// type fixableStage int

// const (
// 	stageContextBreak fixableStage = iota
// 	stageCssSelection
// 	suggestionsProcessing
// 	finalStage
// )

// // Style variables for TUI
// var (
// 	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)
// 	subtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
// 	activeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("190"))
// 	// inactiveStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
// 	// diffAddStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
// 	// diffRemoveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
// )

// func (m *fixableTuiModel) initializeFileProcessing() {
// 	m.fileNames = make([]string, 0, len(m.fileTexts))
// 	for filename := range m.fileTexts {
// 		m.fileNames = append(m.fileNames, filename)
// 	}
// 	m.currentFileIndex = 0
// 	m.currentFile = m.fileNames[0]
// 	m.originalText = m.fileTexts[m.currentFile]
// }

// func (m *fixableTuiModel) processNextFile() bool {
// 	m.currentFileIndex++
// 	if m.currentFileIndex >= len(m.fileNames) {
// 		return false
// 	}
// 	m.currentFile = m.fileNames[m.currentFileIndex]
// 	m.originalText = m.fileTexts[m.currentFile]
// 	m.currentIssueIndex = 0
// 	return true
// }

// func (m *fixableTuiModel) retrieveSuggestions() {
// 	// Reset suggestions for the current issue
// 	m.currentSuggestionGroup = ""
// 	m.suggestions = make(map[string]string)

// 	// Find next enabled issue
// 	for m.currentIssueIndex < len(m.potentiallyFixableIssues) {
// 		issue := m.potentiallyFixableIssues[m.currentIssueIndex]

// 		// Special handling for section breaks
// 		if issue.name == "Potential Section Breaks" && m.contextBreakInput == "" {
// 			m.currentIssueIndex++
// 			continue
// 		}

// 		if m.runAll || *issue.isEnabled {
// 			// For section breaks, pass the context break
// 			var suggestions map[string]string = issue.getSuggestions(m.originalText)

// 			if len(suggestions) > 0 {
// 				m.currentSuggestionGroup = issue.name
// 				m.suggestions = suggestions

// 				// Select the first suggestion
// 				for key := range suggestions {
// 					m.currentSuggestionKey = key
// 					return
// 				}
// 			}
// 		}
// 		m.currentIssueIndex++
// 	}

// 	// If no suggestions found, try next file
// 	if m.processNextFile() {
// 		m.retrieveSuggestions()
// 	} else {
// 		// No more files or suggestions
// 		m.stage = finalStage
// 	}
// }

// func (m fixableTuiModel) Init() tea.Cmd {
// 	return nil
// }

// func (m fixableTuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch m.stage {
// 	case stageContextBreak:
// 		return m.updateContextBreak(msg)
// 	case stageCssSelection:
// 		return m.updateCssSelection(msg)
// 	case suggestionsProcessing:
// 		return m.updateSuggestions(msg)
// 	case finalStage:
// 		return m, tea.Quit
// 	}
// 	return m, nil
// }

// func (m *fixableTuiModel) updateContextBreak(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		case tea.KeyEnter:
// 			if strings.TrimSpace(m.contextBreakInput) != "" {
// 				m.stage = stageCssSelection
// 				return m, nil
// 			}
// 		case tea.KeyBackspace:
// 			if len(m.contextBreakInput) > 0 {
// 				m.contextBreakInput = m.contextBreakInput[:len(m.contextBreakInput)-1]
// 			}
// 		case tea.KeyRunes:
// 			m.contextBreakInput += string(msg.Runes)
// 		}
// 	}
// 	return m, nil
// }

// func (m *fixableTuiModel) updateCssSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "enter":
// 			if m.selectedCssFileIndex >= 0 && m.selectedCssFileIndex < len(m.cssFiles) {
// 				m.stage = suggestionsProcessing
// 				m.initializeFileProcessing()
// 				m.retrieveSuggestions()
// 				return m, nil
// 			}
// 		case "up":
// 			if m.selectedCssFileIndex > 0 {
// 				m.selectedCssFileIndex--
// 			}
// 		case "down":
// 			if m.selectedCssFileIndex < len(m.cssFiles)-1 {
// 				m.selectedCssFileIndex++
// 			}
// 		}
// 	}
// 	return m, nil
// }

// func (m *fixableTuiModel) updateSuggestions(msg tea.Msg) (tea.Model, tea.Cmd) {
// switch msg := msg.(type) {
// case tea.KeyMsg:
// 	switch msg.String() {
// 	case "ctrl+c", "q":
// 		return m, tea.Quit

// 	case "e":
// 		if !m.editMode && m.currentSuggestionKey != "" {
// 			m.editMode = true
// 			m.editedText = m.suggestions[m.currentSuggestionKey]
// 			return m, nil
// 		}

// 	case "enter":
// 		if m.editMode {
// 			// Save edited text
// 			m.suggestions[m.currentSuggestionKey] = m.editedText
// 			m.editMode = false
// 			return m, nil
// 		}

// 		// Accept current suggestion
// 		if m.currentSuggestionKey != "" {
// 			// TODO: the replace count should be -1 in some instances
// 			m.originalText = strings.Replace(m.originalText, m.currentSuggestionKey, m.suggestions[m.currentSuggestionKey], 1)
// 			m.fileTexts[m.currentFile] = m.originalText
// 			// delete(m.suggestions, m.currentSuggestionKey)

// 			if len(m.suggestions) > 0 {
// 				// Select the first remaining suggestion
// 				for key := range m.suggestions {
// 					m.currentSuggestionKey = key
// 					break
// 				}
// 			} else {
// 				// Move to next issue or file
// 				m.currentIssueIndex++
// 				m.retrieveSuggestions()
// 			}
// 		} else {
// 			// Move to next issue or file
// 			m.currentIssueIndex++
// 			m.retrieveSuggestions()
// 		}
// 		return m, nil

// 	case "c":
// 		// Copy current suggestion to clipboard
// 		if m.currentSuggestionKey != "" {
// 			clipboard.WriteAll(m.suggestions[m.currentSuggestionKey])
// 		}
// 		return m, nil

// 	case "right", "l":
// 		if m.editMode {
// 			return m, nil
// 		}
// 		// Move to next suggestion
// 		if m.currentSuggestionKey != "" && len(m.suggestions) > 1 {
// 			var keys []string
// 			for key := range m.suggestions {
// 				keys = append(keys, key)
// 			}

// 			// Find current key index
// 			currentIndex := -1
// 			for i, key := range keys {
// 				if key == m.currentSuggestionKey {
// 					currentIndex = i
// 					break
// 				}
// 			}

// 			// Select next key
// 			nextIndex := (currentIndex + 1) % len(keys)
// 			m.currentSuggestionKey = keys[nextIndex]
// 		}
// 		return m, nil

// 	case "left", "h":
// 		if m.editMode {
// 			return m, nil
// 		}
// 		// Move to previous suggestion
// 		if m.currentSuggestionKey != "" && len(m.suggestions) > 1 {
// 			var keys []string
// 			for key := range m.suggestions {
// 				keys = append(keys, key)
// 			}

// 			// Find current key index
// 			currentIndex := -1
// 			for i, key := range keys {
// 				if key == m.currentSuggestionKey {
// 					currentIndex = i
// 					break
// 				}
// 			}

// 			// Select previous key
// 			prevIndex := (currentIndex - 1 + len(keys)) % len(keys)
// 			m.currentSuggestionKey = keys[prevIndex]
// 		}
// 		return m, nil
// 	}

// 	// Handle edit mode text input
// 	if m.editMode {
// 		switch msg.Type {
// 		case tea.KeyBackspace:
// 			if len(m.editedText) > 0 {
// 				m.editedText = m.editedText[:len(m.editedText)-1]
// 			}
// 		case tea.KeyRunes:
// 			m.editedText += string(msg.Runes)
// 		}
// 		return m, nil
// 	}
// }
// return m, nil
// }

// func (m fixableTuiModel) View() string {
// 	switch m.stage {
// 	case stageContextBreak:
// 		return fmt.Sprintf("Enter section break context:\n\n> %s", m.contextBreakInput)

// 	case stageCssSelection:
// 		var s strings.Builder
// 		s.WriteString("Select CSS file to modify:\n\n")
// 		for i, file := range m.cssFiles {
// 			if i == m.selectedCssFileIndex {
// 				s.WriteString(fmt.Sprintf("> %s\n", file))
// 			} else {
// 				s.WriteString(fmt.Sprintf("  %s\n", file))
// 			}
// 		}
// 		return s.String()

// 	case suggestionsProcessing:
// 		if m.saveAndQuit {
// 			return "Finished processing files. Saving changes...\n"
// 		}

// 		// No more suggestions
// 		if len(m.suggestions) == 0 {
// 			return "No more suggestions. Press 'q' to quit.\n"
// 		}

// 		var s strings.Builder
// s.WriteString(titleStyle.Render(fmt.Sprintf("Current File: %s", m.currentFile)) + "\n")
// s.WriteString(subtitleStyle.Render(fmt.Sprintf("Issue Group: %s", m.currentSuggestionGroup)) + "\n\n")

// 		// Show current suggestion
// 		if m.currentSuggestionKey != "" {
// 			if m.editMode {
// 				s.WriteString("Edit Mode (press enter to confirm):\n")
// 				s.WriteString(activeStyle.Render(m.editedText) + "\n\n")
// 			} else {
// 				// Generate diff view
// 				diffString, err := stringdiff.GetPrettyDiffString(
// 					strings.TrimLeft(m.currentSuggestionKey, "\n"),
// 					strings.TrimLeft(m.suggestions[m.currentSuggestionKey], "\n"),
// 				)
// 				if err != nil {
// 					s.WriteString("Error generating diff: " + err.Error() + "\n")
// 				} else {
// 					s.WriteString(diffString + "\n\n")
// 				}
// 			}
// 		}

// 		// Controls help
// 		s.WriteString(subtitleStyle.Render("Controls:") + "\n")
// s.WriteString("← / → : Previous/Next Suggestion   ")
// s.WriteString("Enter: Accept   ")
// s.WriteString("E: Edit   ")
// s.WriteString("C: Copy   ")
// s.WriteString("Q: Quit\n")

// 		// Suggestion progress
// 		s.WriteString(fmt.Sprintf("\nSuggestion %d of %d", 1, len(m.suggestions)) + "\n")

// 		return s.String()

// 	case finalStage:
// 		return "Processing complete. Press 'q' to quit.\n"

// 	default:
// 		return "Unexpected stage\n"
// 	}
// }

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
			return nil, ErrCssPathsEmptyWhenArgIsNeeded
		}

		var (
			initialModel = NewFixableTuiModel(runAll, runSectionBreak, potentiallyFixableIssues, cssFiles)
			i            = 0
		)
		initialModel.filePaths = make([]string, len(epubInfo.HtmlFiles))

		// Collect file contents
		for file := range epubInfo.HtmlFiles {
			var filePath = getFilePath(opfFolder, file)
			zipFile := zipFiles[filePath]

			fileText, err := filehandler.ReadInZipFileContents(zipFile)
			if err != nil {
				return nil, err
			}

			initialModel.fileTexts[filePath] = linter.CleanupHtmlSpacing(fileText)
			initialModel.filePaths[i] = filePath
			i++
		}

		p := tea.NewProgram(&initialModel, tea.WithAltScreen())
		finalModel, err := p.Run()
		if err != nil {
			return nil, err
		}

		model := finalModel.(FixableTuiModel)
		if model.Err != nil {
			return nil, model.Err
		}

		// Process and write updated files
		// for filePath, fileText := range model.fileTexts {
		// 	err = filehandler.WriteZipCompressedString(w, filePath, fileText)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	model.handledFiles = append(model.handledFiles, filePath)
		// }

		// // Handle CSS changes
		// return handleCssChanges(false, false, opfFolder, cssFiles, model.contextBreakInput, zipFiles, w, model.handledFiles)

		return nil, nil
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
			return nil, ErrCssPathsEmptyWhenArgIsNeeded
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

					if potentiallyFixableIssue.addCssIfMissing && updateMade {
						addCssSectionIfMissing = addCssSectionIfMissing || updateMade
					}
				}
			}

			err = filehandler.WriteZipCompressedString(w, filePath, newText)
			if err != nil {
				return nil, err
			}

			handledFiles = append(handledFiles, filePath)
		}

		return handleCssChanges(addCssSectionIfMissing, addCssPageIfMissing, opfFolder, cssFiles, contextBreak, zipFiles, w, handledFiles)
	})
	// if err != nil {
	// 	logger.WriteErrorf("failed to fix manually fixable issues for %q: %s", epubFile, err)
	// }

	// logger.WriteInfo("\nFinished showing manually fixable issues...")
}
