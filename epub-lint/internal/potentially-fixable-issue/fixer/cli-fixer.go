package fixer

import (
	"fmt"
	"strings"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	stringdiff "github.com/pjkaufman/go-go-gadgets/pkg/string-diff"
)

const cliLineSeparator = "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"

type CliFixer struct {
	potentiallyFixableIssues                    []potentiallyfixableissue.PotentiallyFixableIssue
	epubInfo                                    *epubhandler.EpubInfo
	getFile                                     FileGetter
	writeFile                                   FileWriter
	cssFiles, handledFiles                      []string
	opfFolder                                   string
	contextBreak                                *string
	runAll, skipCss, runSectionBreak            bool
	addCssSectionIfMissing, addCssPageIfMissing bool
}

func (t *CliFixer) InitialLog() string {
	return "Started showing manually fixable issues...\n"
}

func (t *CliFixer) SuccessfulLog() string {
	return "\nFinished showing manually fixable issues..."
}

func (t *CliFixer) Init(epubInfo *epubhandler.EpubInfo, runAll, skipCss, runSectionBreak bool, potentiallyFixableIssues []potentiallyfixableissue.PotentiallyFixableIssue, cssFiles []string, logFile, opfFolder string, contextBreak *string, getFile FileGetter, writeFile FileWriter) {
	t.epubInfo = epubInfo
	t.runAll = runAll
	t.skipCss = skipCss
	t.runSectionBreak = runSectionBreak
	t.potentiallyFixableIssues = potentiallyFixableIssues
	t.cssFiles = cssFiles
	t.opfFolder = opfFolder
	t.getFile = getFile
	t.writeFile = writeFile
	t.contextBreak = contextBreak
}

func (t *CliFixer) Setup() error {
	if !t.skipCss && (t.runAll || t.runSectionBreak) {
		*t.contextBreak = logger.GetInputString("What is the section break for the epub?:")

		if strings.TrimSpace(*t.contextBreak) == "" {
			return fmt.Errorf("please provide a non-whitespace section break")
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

	return nil
}

func (t *CliFixer) Run() error {
	var saveAndQuit = false
	for file := range t.epubInfo.HtmlFiles {
		if saveAndQuit {
			break
		}

		var filePath = getFilePath(t.opfFolder, file)
		fileText, err := t.getFile(filePath)
		if err != nil {
			return err
		}

		var newText = linter.CleanupHtmlSpacing(fileText)

		for _, potentiallyFixableIssue := range t.potentiallyFixableIssues {
			if saveAndQuit {
				break
			}

			if t.skipCss && (potentiallyFixableIssue.AddCssPageBreakIfMissing || potentiallyFixableIssue.AddCssSectionBreakIfMissing) {
				continue
			}

			if potentiallyFixableIssue.IsEnabled == nil {
				return fmt.Errorf("%q is not properly setup to run as a potentially fixable rule since it has no boolean for isEnabled", potentiallyFixableIssue.Name)
			}

			if t.runAll || *potentiallyFixableIssue.IsEnabled {
				suggestions, err := potentiallyFixableIssue.GetSuggestions(newText)
				if err != nil {
					return err
				}

				var updateMade bool
				newText, updateMade, saveAndQuit = promptAboutSuggestions(potentiallyFixableIssue.Name, suggestions, newText, potentiallyFixableIssue.UpdateAllInstances)

				if potentiallyFixableIssue.AddCssSectionBreakIfMissing && updateMade {
					t.addCssSectionIfMissing = t.addCssSectionIfMissing || updateMade
				}

				if potentiallyFixableIssue.AddCssPageBreakIfMissing && updateMade {
					t.addCssPageIfMissing = t.addCssPageIfMissing || updateMade
				}
			}
		}

		err = t.writeFile(filePath, newText)
		if err != nil {
			return err
		}

		t.handledFiles = append(t.handledFiles, filePath)
	}

	return nil
}

func (t *CliFixer) HandleCss() ([]string, error) {
	if !t.addCssSectionIfMissing && !t.addCssPageIfMissing {
		return t.handledFiles, nil
	}

	var cssSelectionPrompt = "Please enter the number of the css file to append the css to:\n"
	for i, file := range t.cssFiles {
		cssSelectionPrompt += fmt.Sprintf("%d. %s\n", i, file)
	}

	var selectedCssFileIndex = logger.GetInputInt(cssSelectionPrompt)
	if selectedCssFileIndex < 0 || selectedCssFileIndex >= len(t.cssFiles) {
		return nil, fmt.Errorf("please select a valid css file value instead of \"%d\".\n", selectedCssFileIndex)
	}

	return updateCssFile(t.addCssSectionIfMissing, t.addCssPageIfMissing, filehandler.JoinPath(t.opfFolder, t.cssFiles[selectedCssFileIndex]), *t.contextBreak, t.handledFiles, t.getFile, t.writeFile)
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
