package fixer

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/ui"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

type TuiFixer struct {
	initialModel                                ui.FixableIssuesModel
	potentiallyFixableIssues                    []potentiallyfixableissue.PotentiallyFixableIssue
	epubInfo                                    *epubhandler.EpubInfo
	getFile                                     FileGetter
	writeFile                                   FileWriter
	cssFiles, handledFiles                      []string
	logFile                                     string
	opfFolder                                   string
	selectedCssFile                             string
	contextBreak                                string
	runAll, skipCss, runSectionBreak            bool
	addCssSectionIfMissing, addCssPageIfMissing bool
}

// InitialLog is just meant to allow the CLI version to return its initial log
func (t *TuiFixer) InitialLog() string {
	return ""
}

// SuccessfulLog is just meant to allow the CLI version to return its success log
func (t *TuiFixer) SuccessfulLog() string {
	return ""
}

func (t *TuiFixer) Init(epubInfo *epubhandler.EpubInfo, runAll, skipCss, runSectionBreak bool, potentiallyFixableIssues []potentiallyfixableissue.PotentiallyFixableIssue, cssFiles []string, logFile, opfFolder string, getFile FileGetter, writeFile FileWriter) {
	t.epubInfo = epubInfo
	t.runAll = runAll
	t.skipCss = skipCss
	t.runSectionBreak = runSectionBreak
	t.potentiallyFixableIssues = potentiallyFixableIssues
	t.cssFiles = cssFiles
	t.logFile = logFile
	t.opfFolder = opfFolder
	t.getFile = getFile
	t.writeFile = writeFile
}

func (t *TuiFixer) Setup() error {
	var file *os.File
	if strings.TrimSpace(t.logFile) != "" {
		file, err := tea.LogToFile(t.logFile, "debug")
		if err != nil {
			return fmt.Errorf("failed to create TUI log file %q: %w", t.logFile, err)
		}

		defer file.Close()
	}

	t.initialModel = ui.NewFixableIssuesModel(t.runAll, t.skipCss, t.runSectionBreak, t.potentiallyFixableIssues, t.cssFiles, file, &t.contextBreak)
	var i = 0

	t.initialModel.PotentiallyFixableIssuesInfo.FileSuggestionData = make([]ui.FileSuggestionInfo, len(t.epubInfo.HtmlFiles))

	var numFixableIssues = len(t.potentiallyFixableIssues)
	// Collect file contents
	for file := range t.epubInfo.HtmlFiles {
		var filePath = getFilePath(t.opfFolder, file)
		fileText, err := t.getFile(filePath)
		if err != nil {
			return err
		}

		t.initialModel.PotentiallyFixableIssuesInfo.FileSuggestionData[i] = ui.FileSuggestionInfo{
			Name:        filePath,
			Text:        linter.CleanupHtmlSpacing(fileText),
			Suggestions: make([][]ui.SuggestionState, numFixableIssues),
		}

		i++
	}

	sort.Slice(t.initialModel.PotentiallyFixableIssuesInfo.FileSuggestionData, func(i, j int) bool {
		return t.initialModel.PotentiallyFixableIssuesInfo.FileSuggestionData[i].Name < t.initialModel.PotentiallyFixableIssuesInfo.FileSuggestionData[j].Name
	})

	return nil
}

func (t *TuiFixer) Run() error {
	p := tea.NewProgram(&t.initialModel, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	model := finalModel.(ui.FixableIssuesModel)
	if model.Err != nil {
		if errors.Is(model.Err, ui.ErrUserKilledProgram) {
			logger.WriteInfo("Quitting. User exited the program...")
			os.Exit(0)
		}

		return model.Err
	}

	// TODO: can this be done in a way that is handled by both  CLI and TUI?
	t.handledFiles = make([]string, len(model.PotentiallyFixableIssuesInfo.FileSuggestionData))
	for _, fileData := range model.PotentiallyFixableIssuesInfo.FileSuggestionData {
		t.writeFile(fileData.Name, fileData.Text)
		if err != nil {
			return err
		}

		t.handledFiles = append(t.handledFiles, fileData.Name)
	}

	t.addCssPageIfMissing = model.PotentiallyFixableIssuesInfo.AddCssPageBreakIfMissing
	t.addCssSectionIfMissing = model.PotentiallyFixableIssuesInfo.AddCssSectionBreakIfMissing
	t.selectedCssFile = model.CssSelectionInfo.SelectedCssFile

	return nil
}

func (t *TuiFixer) HandleCss() ([]string, error) {
	if !t.addCssSectionIfMissing && !t.addCssPageIfMissing {
		return t.handledFiles, nil
	}

	if strings.TrimSpace(t.selectedCssFile) == "" {
		return nil, fmt.Errorf("please select a valid css file instead of %q.\n", t.selectedCssFile)
	}

	return updateCssFile(t.addCssSectionIfMissing, t.addCssPageIfMissing, filehandler.JoinPath(t.opfFolder, t.selectedCssFile), t.contextBreak, t.handledFiles, t.getFile, t.writeFile)
}
