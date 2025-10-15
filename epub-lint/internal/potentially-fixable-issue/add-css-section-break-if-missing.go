package potentiallyfixableissue

import (
	"fmt"
	"strings"
)

const (
	HrCharacter = `hr.character {
overflow: visible;
border:0;
text-align:center;
}`
	HrContentAfterTemplate = `hr.character:after {
content: %q;
display:inline-block;
position:relative;
font-size:1em;
padding:1em;
}`
)

func AddCssSectionBreakIfMissing(fileContent, contextBreak string) string {
	if strings.TrimSpace(fileContent) == "" {
		return HrCharacter + "\n" + fmt.Sprintf(HrContentAfterTemplate, contextBreak)
	}

	if strings.Contains(fileContent, HrCharacter) {
		return fileContent
	}

	return fmt.Sprintf("%s\n%s\n%s", fileContent, HrCharacter, fmt.Sprintf(HrContentAfterTemplate, contextBreak))
}
