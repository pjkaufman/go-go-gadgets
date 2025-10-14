package potentiallyfixableissue

import (
	"fmt"
	"strings"
)

const HrBlankSpace = `hr.blankSpace {
border:0;
height:2em;
}`

func AddCssPageBreakIfMissing(fileContent string) string {
	if strings.TrimSpace(fileContent) == "" {
		return HrBlankSpace + "\n"
	}

	if strings.Contains(fileContent, HrBlankSpace) {
		return fileContent
	}

	return fmt.Sprintf("%s\n%s", fileContent, HrBlankSpace)
}
