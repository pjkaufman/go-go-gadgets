//go:build unit

package converter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/stretchr/testify/assert"
)

type convertMdToCsvTestCase struct {
	inputFilePath    string
	inputFileContent string
	expectedCsv      string
}

var convertMdToCsvTestCases = map[string]convertMdToCsvTestCase{
	"a valid file should properly get turned into a csv row": {
		inputFilePath:    "He is Lord.md",
		inputFileContent: HeIsLordFileMd,
		expectedCsv:      HeIsLordFileCsv,
	},
	"make sure that no location is handled properly": {
		inputFilePath:    "Bless This House.md",
		inputFileContent: BlessThisHouseFileMd,
		expectedCsv:      BlessThisHouseFileCsv,
	},
	"make sure that blue book locations are handled properly": {
		inputFilePath:    "Bigger Than All My Problems.md",
		inputFileContent: BiggerThanAllOfMyProblemsFileMd,
		expectedCsv:      BiggerThanAllOfMyProblemsFileCsv,
	},
	"make sure that more songs we love locations are handled properly": {
		inputFilePath:    "Flow Thou River.md",
		inputFileContent: FlowThowRiverFileMd,
		expectedCsv:      FlowThowRiverFileCsv,
	},
	"make sure that copyright for authors in the church are set to 'Church'": {
		inputFilePath:    "Fill My Soul With Thy Spirit.md",
		inputFileContent: FillMySoulWithThySpiritFileMd,
		expectedCsv:      FillMySoulWithThySpiritFileCsv,
	},
	"make sure that copyright YAML property gets used when not in church and it is present2": {
		inputFilePath:    "A Glorious Church.md",
		inputFileContent: AGloriousChurchFileMd,
		expectedCsv:      AGloriousChurchFileCsv,
	},
}

func TestConvertMdToCsv(t *testing.T) {
	for name, args := range convertMdToCsvTestCases {
		t.Run(name, func(t *testing.T) {

			actual, err := converter.ConvertMdToCsv(args.inputFilePath, args.inputFilePath, args.inputFileContent)
			if err != nil {
				assert.Fail(t, "there should be no errors when parsing the YAML for the csv UTs")
			}

			assert.Equal(t, args.expectedCsv, actual)
		})
	}
}
