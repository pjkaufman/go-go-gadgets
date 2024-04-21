package converter

import "strings"

const CsvHeading = "Song|Location|Author|Copyright\n"

func BuildCsv(mdInfo []MdFileInfo) (string, error) {
	var csvContents = strings.Builder{}
	csvContents.WriteString(CsvHeading)

	for _, mdData := range mdInfo {
		csvString, err := ConvertMdToCsv(mdData.FileName, mdData.FilePath, mdData.FileContents)
		if err != nil {
			return "", err
		}

		csvContents.WriteString(csvString)
	}

	return strings.ReplaceAll(csvContents.String(), "&nbsp;&nbsp;&nbsp;&nbsp;", " "), nil
}
