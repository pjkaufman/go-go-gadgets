package fixer

import (
	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
)

func getFilePath(opfFolder, file string) string {
	return filehandler.JoinPath(opfFolder, file)
}

func updateCssFile(addCssSectionIfMissing, addCssPageIfMissing bool, selectedCssFile, contextBreak string, handledFiles []string, getFile FileGetter, writeFile FileWriter) ([]string, error) {
	css, err := getFile(selectedCssFile)
	if err != nil {
		return nil, err
	}

	var newCssText = css
	if addCssSectionIfMissing {
		newCssText = potentiallyfixableissue.AddCssSectionBreakIfMissing(newCssText, contextBreak)
	}

	if addCssPageIfMissing {
		newCssText = potentiallyfixableissue.AddCssPageBreakIfMissing(newCssText)
	}

	if newCssText == css {
		return handledFiles, nil
	}

	err = writeFile(selectedCssFile, newCssText)
	if err != nil {
		return nil, err
	}

	return append(handledFiles, selectedCssFile), nil
}
