package converter

import (
	"fmt"
	"strings"
)

func ConvertMdToCsv(fileName, filePath, fileContents string) (string, error) {
	var metadata SongMetadata
	_, err := parseFrontmatter(filePath, fileContents, &metadata)
	if err != nil {
		return "", err
	}

	return strings.Replace(fileName, ".md", "", 1) + "|" + buildMetadataCsv(&metadata) + "\n", nil
}

func buildMetadataCsv(metadata *SongMetadata) string {
	if metadata == nil {
		return "||"
	}

	var copyright = metadata.Copyright
	if strings.EqualFold(metadata.InChurch, "Y") {
		copyright = "Church"
	}

	return fmt.Sprintf("%s|%s|%s", updateBookLocationInfo(metadata.BookLocation), metadata.Authors, copyright)
}

func updateBookLocationInfo(bookLocation string) string {
	if bookLocation == "" {
		return ""
	}

	var newBookLocation = strings.ReplaceAll(bookLocation, "B", "Blue Book page ")
	newBookLocation = strings.ReplaceAll(newBookLocation, "R", "Red Book page ")
	newBookLocation = strings.ReplaceAll(newBookLocation, "MS", "More Songs We Love page ")

	return newBookLocation
}
