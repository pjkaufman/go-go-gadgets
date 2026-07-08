package fixer

import (
	"errors"
	"fmt"
	"strings"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	suggestionmanager "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/suggestion-manager"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
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
	suggestionManager                           *suggestionmanager.SuggestionManager
}

func (c *CliFixer) InitialLog() string {
	return "Started showing manually fixable issues...\n"
}

func (c *CliFixer) SuccessfulLog() string {
	return "\nFinished showing manually fixable issues..."
}

func (c *CliFixer) Init(epubInfo *epubhandler.EpubInfo, runAll, skipCss, runSectionBreak bool, potentiallyFixableIssues []potentiallyfixableissue.PotentiallyFixableIssue, cssFiles []string, logFile, opfFolder string, contextBreak *string, getFile FileGetter, writeFile FileWriter) {
	c.epubInfo = epubInfo
	c.runAll = runAll
	c.skipCss = skipCss
	c.runSectionBreak = runSectionBreak
	c.potentiallyFixableIssues = potentiallyFixableIssues
	c.cssFiles = cssFiles
	c.opfFolder = opfFolder
	c.getFile = getFile
	c.writeFile = writeFile
	c.contextBreak = contextBreak
}

func (c *CliFixer) Setup() error {
	if !c.skipCss && (c.runAll || c.runSectionBreak) {
		*c.contextBreak = logger.GetInputString("What is the section break for the epub?:")

		if strings.TrimSpace(*c.contextBreak) == "" {
			return errors.New("please provide a non-whitespace section break")
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

	var filePathToText = make(map[string]string, len(c.epubInfo.HtmlFiles))
	// Collect file contents
	for file := range c.epubInfo.HtmlFiles {
		var filePath = getFilePath(c.opfFolder, file)
		fileText, err := c.getFile(filePath)
		if err != nil {
			return err
		}

		filePathToText[filePath] = linter.CleanupHtmlSpacing(fileText)
	}

	c.suggestionManager = suggestionmanager.NewSuggestionManager(c.potentiallyFixableIssues, filePathToText, c.runAll, c.skipCss, nil)

	return nil
}

func (c *CliFixer) Run() error {
	var saveAndQuit = false

	foundSuggestion, err := c.suggestionManager.SetupForNextSuggestions()
	if err != nil {
		return err
	}

	for foundSuggestion {
		if saveAndQuit {
			break
		}

		var updateMade bool
		updateMade, saveAndQuit = promptAboutSuggestions(c.suggestionManager)

		if c.suggestionManager.CurrentSuggestion.AddCssSectionBreakIfMissing && updateMade {
			c.addCssSectionIfMissing = c.addCssSectionIfMissing || updateMade
		}

		if c.suggestionManager.CurrentSuggestion.AddCssPageBreakIfMissing && updateMade {
			c.addCssPageIfMissing = c.addCssPageIfMissing || updateMade
		}

		foundSuggestion, err = c.suggestionManager.SetupForNextSuggestions()
		if err != nil {
			return err
		}
	}

	c.handledFiles = make([]string, len(c.suggestionManager.FileSuggestionData))
	for _, fileData := range c.suggestionManager.FileSuggestionData {
		err = c.writeFile(fileData.Name, fileData.Text)
		if err != nil {
			return err
		}

		c.handledFiles = append(c.handledFiles, fileData.Name)
	}

	return nil
}

func (c *CliFixer) HandleCss() ([]string, error) {
	if !c.addCssSectionIfMissing && !c.addCssPageIfMissing {
		return c.handledFiles, nil
	}

	var cssSelectionPrompt strings.Builder
	cssSelectionPrompt.WriteString("Please enter the number of the css file to append the css to:\n")

	for i, file := range c.cssFiles {
		fmt.Fprintf(&cssSelectionPrompt, "%d. %s\n", i, file)
	}

	var selectedCssFileIndex = logger.GetInputInt(cssSelectionPrompt.String())
	if selectedCssFileIndex < 0 || selectedCssFileIndex >= len(c.cssFiles) {
		return nil, fmt.Errorf("please select a valid css file value instead of \"%d\".\n", selectedCssFileIndex)
	}

	return updateCssFile(c.addCssSectionIfMissing, c.addCssPageIfMissing, filehandler.JoinPath(c.opfFolder, c.cssFiles[selectedCssFileIndex]), *c.contextBreak, c.handledFiles, c.getFile, c.writeFile)
}

func promptAboutSuggestions(suggestionManager *suggestionmanager.SuggestionManager) (bool, bool) {
	var valueReplaced = false

	if suggestionManager.CurrentSuggestionState == nil {
		return valueReplaced, false
	}

	logger.WriteInfo(cliLineSeparator)
	logger.WriteInfo(suggestionManager.CurrentSuggestionName + ":")
	logger.WriteInfo(cliLineSeparator + "\n")

	var hasSuggestion = true
	for hasSuggestion {
		err := suggestionManager.CurrentSuggestionState.GetStringDiffAsDisplay()
		if err != nil {
			logger.WriteFatal(err.Error())
		}

		//nolint:gocritic // Warning: do not use %q on the following line as it will get rid of the color coding of changes in the terminal
		resp := logger.GetInputString(fmt.Sprintf("Would you like to make the following update \"%s\"? (Y/N/Q): ", suggestionManager.CurrentSuggestionState.Display))
		switch strings.ToLower(resp) {
		case "y":
			err = suggestionManager.AcceptSuggestion()
			if err != nil {
				logger.WriteFatal(err.Error())
			}

			valueReplaced = true
		case "q":
			return valueReplaced, true
		}

		logger.WriteInfo("")

		hasSuggestion = suggestionManager.MoveToNextSuggestion()
	}

	return valueReplaced, false
}

// Cleanup is empty for now as there is not really anything to cleanup here
func (c *CliFixer) Cleanup() {
}
