package rulefixes

func FixMissingImageAlt(line, column int, contents string) string {
	if line < 1 {
		return contents
	}

	// column is the index of the `>` in `/>`
	offset := GetPositionOffset(contents, line, column)
	if offset == -1 {
		return contents
	}

	var emptyAlt = "alt=\"\""
	if contents[offset-3] != ' ' && contents[offset-3] != '\t' {
		emptyAlt = " " + emptyAlt
	}

	return contents[:offset-2] + emptyAlt + contents[offset-2:]
}
