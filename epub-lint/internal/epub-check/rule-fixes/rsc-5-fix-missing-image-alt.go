package rulefixes

func FixMissingImageAlt(line, column int, contents string) (string, int) {
	if line < 1 {
		return contents, 0
	}

	// column is the index of the `>` in `/>`
	offset := getColumnOffset(contents, line, column)
	if offset == -1 {
		return contents, 0
	}

	var emptyAlt = "alt=\"\""
	if contents[offset-3] != ' ' && contents[offset-3] != '\t' {
		emptyAlt = " " + emptyAlt
	}

	return contents[:offset-2] + emptyAlt + contents[offset-2:], len(emptyAlt)
}
