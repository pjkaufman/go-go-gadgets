package compare

type DifferenceType = int

const (
	LikelyMismatch DifferenceType = iota
	DefiniteMismatch
	WrappedLine
	PartiallyWrappedLine
	Whitespace
	Line
)
