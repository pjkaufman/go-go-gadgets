//go:build unit

package converter_test

import (
	"fmt"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type BuildCsvTestCase struct {
	InputMdInfo []converter.MdFileInfo
	ExpectedCsv string
}

var BuildCsvTestCases = map[string]BuildCsvTestCase{
	"no files provided should just result in the heading line plus a new line character": {
		ExpectedCsv: converter.CsvHeading,
	},
	"multiple files should be added together with a new line character after it": {
		InputMdInfo: []converter.MdFileInfo{
			{
				FilePath:     "A Glorious Church.md",
				FileName:     "A Glorious Church.md",
				FileContents: AGloriousChurchFileMd,
			},
			{
				FilePath:     "A Hymn Of Praise.md",
				FileName:     "A Hymn Of Praise.md",
				FileContents: AHymnOfPraiseFileMd,
			},
			{
				FilePath:     "Bless This House.md",
				FileName:     "Bless This House.md",
				FileContents: BlessThisHouseFileMd,
			},
			{
				FilePath:     "Bigger Than All My Problems.md",
				FileName:     "Bigger Than All My Problems.md",
				FileContents: BiggerThanAllOfMyProblemsFileMd,
			},
			{
				FilePath:     "Fill My Soul With Thy Spirit.md",
				FileName:     "Fill My Soul With Thy Spirit.md",
				FileContents: FillMySoulWithThySpiritFileMd,
			},
			{
				FilePath:     "Flow Thou River.md",
				FileName:     "Flow Thou River.md",
				FileContents: FlowThowRiverFileMd,
			},
			{
				FilePath:     "He is Lord.md",
				FileName:     "He is Lord.md",
				FileContents: HeIsLordFileMd,
			},
		},
		ExpectedCsv: fmt.Sprintf("%s%s%s%s%s%s%s%s", converter.CsvHeading, AGloriousChurchFileCsv, AHymnOfPraiseFileCsvCleaned, BlessThisHouseFileCsv, BiggerThanAllOfMyProblemsFileCsv, FillMySoulWithThySpiritFileCsv, FlowThowRiverFileCsv, HeIsLordFileCsv),
	},
}

func TestBuildCsv(t *testing.T) {
	for name, args := range BuildCsvTestCases {
		t.Run(name, func(t *testing.T) {

			actual, err := converter.BuildCsv(args.InputMdInfo)
			if err != nil {
				assert.Fail(t, "there should be no errors when parsing the YAML for the html UTs")
			}

			assert.Equal(t, args.ExpectedCsv, actual)
		})
	}
}
