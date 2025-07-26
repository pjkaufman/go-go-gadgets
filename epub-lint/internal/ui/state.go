package ui

type State struct {
	// general
	CurrentStage          int
	BodyHeight, BodyWidth int
	Ready                 bool
	RunAll                bool
	// file data
	FilePaths []string
	FileTexts []string
	// body data
	ContextBreak string
}
