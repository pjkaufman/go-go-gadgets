package fixer

import (
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	potentiallyfixableissue "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/potentially-fixable-issue"
)

type Fixer interface {
	InitialLog() string
	Init(epubInfo *epubhandler.EpubInfo, runAll, skipCss, runSectionBreak bool, potentiallyFixableIssues []potentiallyfixableissue.PotentiallyFixableIssue, cssFiles []string, logFile, opfFolder string, contextBreak *string, getFile FileGetter, writeFile FileWriter)
	Setup() error
	Run() error
	HandleCss() ([]string, error)
	SuccessfulLog() string
}

type FileGetter func(fileName string) (string, error)

type FileWriter func(fileName, content string) error
