//go:build unit

package converter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type ConvertMdToCsvTestCase struct {
	InputFilePath    string
	InputFileContent string
	ExpectedCsv      string
}

var ConvertMdToCsvTestCases = map[string]ConvertMdToCsvTestCase{
	"a valid file should properly get turned into a csv row": {
		InputFilePath:    "He is Lord.md",
		InputFileContent: HeIsLordFileMd,
		ExpectedCsv:      HeIsLordFileCsv,
	},
	"make sure that no location is handled properly": {
		InputFilePath:    "Bless This House.md",
		InputFileContent: BlessThisHouseFileMd,
		ExpectedCsv:      BlessThisHouseFileCsv,
	},
	"make sure that blue book locations are handled properly": {
		InputFilePath:    "Bigger Than All My Problems.md",
		InputFileContent: BiggerThanAllOfMyProblemsFileMd,
		ExpectedCsv:      BiggerThanAllOfMyProblemsFileCsv,
	},
	"make sure that more songs we love locations are handled properly": {
		InputFilePath:    "Flow Thou River.md",
		InputFileContent: FlowThowRiverFileMd,
		ExpectedCsv:      FlowThowRiverFileCsv,
	},
	"make sure that copyright for authors in the church are set to 'Church'": {
		InputFilePath:    "Fill My Soul With Thy Spirit.md",
		InputFileContent: FillMySoulWithThySpiritFileMd,
		ExpectedCsv:      FillMySoulWithThySpiritFileCsv,
	},
	"make sure that copyright YAML property gets used when not in church and it is present2": {
		InputFilePath:    "A Glorious Church.md",
		InputFileContent: AGloriousChurchFileMd,
		ExpectedCsv:      AGloriousChurchFileCsv,
	},
}

func TestConvertMdToCsv(t *testing.T) {
	for name, args := range ConvertMdToCsvTestCases {
		t.Run(name, func(t *testing.T) {

			actual, err := converter.ConvertMdToCsv(args.InputFilePath, args.InputFilePath, args.InputFileContent)
			if err != nil {
				assert.Fail(t, "there should be no errors when parsing the YAML for the csv UTs")
			}

			assert.Equal(t, args.ExpectedCsv, actual)
		})
	}
}
