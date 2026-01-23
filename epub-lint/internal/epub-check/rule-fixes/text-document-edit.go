package rulefixes

type Position struct {
	Line   int
	Column int
}

type Range struct {
	Start Position
	End   Position // end is exclusive, not inclusive
}

type TextEdit struct {
	Range   Range
	NewText string
}

type TextDocumentEdit struct {
	FilePath string
	Edits    []TextEdit
}
