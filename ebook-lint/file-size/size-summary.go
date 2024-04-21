package filesize

import "fmt"

const (
	CliLineSeparator    = "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	FileSummaryTemplate = `
%[1]s
Before:
%s %s
After:
%s %s
%[1]s
`
	FilesSummaryTemplate = `
%[1]s
Before:
%s
After:
%s
%[1]s
`
	kilobytesInAMegabyte float64 = 1024
	kilobytesInAGigabyte float64 = 1000000
)

func FileSizeSummary(originalFile, newFile string, oldKbSize, newKbSize float64) string {
	return fmt.Sprintf(FileSummaryTemplate, CliLineSeparator, originalFile, kbSizeToString(oldKbSize), newFile, kbSizeToString(newKbSize))
}

func FilesSizeSummary(oldKbSizeSum, newKbSizeSum float64) string {
	return fmt.Sprintf(FilesSummaryTemplate, CliLineSeparator, kbSizeToString(oldKbSizeSum), kbSizeToString(newKbSizeSum))
}

func kbSizeToString(size float64) string {
	if size > kilobytesInAGigabyte {
		return fmt.Sprintf("%.2f GB", size/kilobytesInAGigabyte)
	} else if size > kilobytesInAMegabyte {
		return fmt.Sprintf("%.2f MB", size/kilobytesInAMegabyte)
	}

	return fmt.Sprintf("%.2f KB", size)
}
