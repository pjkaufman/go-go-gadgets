package rulefixes

import (
	"fmt"
	"regexp"
	"strings"
)

// UpdateDuplicateIds finds and renames duplicate IDs in file contents.
// Returns the modified contents.
func UpdateDuplicateIds(contents, id string) string {
	// Pattern: id="id" or id='id'
	idPattern := fmt.Sprintf(`id=([\'"])%s[\'"]`, regexp.QuoteMeta(id))
	re := regexp.MustCompile(idPattern)

	matches := re.FindAllStringIndex(contents, -1)
	if len(matches) <= 1 {
		return contents
	}

	var sb strings.Builder
	var lastIdx int

	for i, idx := range matches {
		start, end := idx[0], idx[1]
		sb.WriteString(contents[lastIdx:start])

		// Write id= + quote + id
		quote := contents[start+3]
		sb.WriteString(contents[start : end-1]) // everything except closing quote

		if i > 0 {
			suffix := fmt.Sprintf("_%d", i+1)
			sb.WriteString(suffix)
		}
		sb.WriteByte(quote)
		lastIdx = end
	}
	sb.WriteString(contents[lastIdx:])

	return sb.String()
}
