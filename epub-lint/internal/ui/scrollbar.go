// pulled from https://github.com/charmbracelet/bubbles/pull/536

package ui

// Msg signals that scrollbar parameters must be updated.
type Msg struct {
	Total   int
	Visible int
	Offset  int
}

// HeightMsg signals that scrollbar height must be updated.
type HeightMsg int
