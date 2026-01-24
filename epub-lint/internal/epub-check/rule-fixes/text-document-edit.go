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

func (te TextEdit) IsEmpty() bool {
	return te.NewText == "" && te.Range.Start.Column == 0 && te.Range.Start.Line == 0 && te.Range.End.Column == 0 && te.Range.End.Line == 0
}
