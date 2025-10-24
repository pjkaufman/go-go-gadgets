//go:build unit

package epubcheck_test

import (
	"fmt"
	"testing"

	_ "embed"

	epubcheck "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check"
	"github.com/stretchr/testify/assert"
)

type handleValidationErrorTestCase struct {
	opfFolder, ncxFilename, opfFilename string
	validationErrors                    epubcheck.ValidationErrors
	expectedErrorState                  epubcheck.ValidationErrors
	expectedFileState                   map[string]string // filename to content
	getContentByFileName                func(string) (string, error)
}

var (
	//go:embed testdata/opf-15/remove-properties.opf
	opfRemovePropertiesOriginal string
	//go:embed testdata/opf-15/remove-properties_updated.opf
	opfRemovePropertiesExpected string
	//go:embed testdata/opf-14/add-properties.opf
	opfAddPropertiesOriginal string
	//go:embed testdata/opf-14/add-properties_updated.opf
	opfAddPropertiesExpected string
	//go:embed testdata/ncx-1/no-identifier.opf
	opfNoIdentifierOriginal string
	//go:embed testdata/ncx-1/no-identifier_updated.opf
	opfNoIdentifierExpected string
	//go:embed testdata/ncx-1/uuid-identifier.ncx
	ncxUuidIdentifier string
)

func createTestCaseFileHandlerFunction(validFilesToContent map[string]string) func(string) (string, error) {
	return func(s string) (string, error) {
		if content, ok := validFilesToContent[s]; ok {
			return content, nil
		}

		return "", fmt.Errorf("unexpected attempt to get file contents for file %q", s)
	}
}

var handleValidationErrorTestCases = map[string]handleValidationErrorTestCase{
	"OPF 15: Removing properties from different files should work without issue": {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/content.opf": opfRemovePropertiesOriginal},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0001.html",
					Message:  `The property "svg" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0003.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0004.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0005.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0006.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0007.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0008.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0001.html",
					Message:  `The property "svg" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0003.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0004.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0005.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0006.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0007.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
				{
					Code:     "OPF-015",
					FilePath: "OPS/section-0008.html",
					Message:  `The property "scripted" should not be declared in the OPF file.`,
				},
			},
		},
		getContentByFileName: createTestCaseFileHandlerFunction(map[string]string{
			"OPS/content.opf": opfRemovePropertiesOriginal,
		}),
	},
	"OPF 14: Adding properties to different files should work without issue": {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/content.opf": opfAddPropertiesExpected},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0003.html",
					Message:  `The property "svg" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0003.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0004.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0005.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0006.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0007.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0008.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0003.html",
					Message:  `The property "svg" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0003.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0004.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0005.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0006.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0007.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
				{
					Code:     "OPF-014",
					FilePath: "OPS/section-0008.html",
					Message:  `The property "scripted" should be declared in the OPF file.`,
				},
			},
		},
		getContentByFileName: createTestCaseFileHandlerFunction(map[string]string{
			"OPS/content.opf": opfAddPropertiesExpected,
		}),
	},
	"NCX 1: When no identifier is present in the OPF, add the one from the NCX file and make sure that all OPF updates that are for after that line get incremented a line": {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/content.opf": opfNoIdentifierExpected},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "NCX-001",
					FilePath: "OPS/toc.ncx",
					Message:  `NCX identifier ("urn:uuid:1da9fa05e-dd8b-4be3-85ab-455656cc14f2") does not match OPF identifier ("").`,
				},
				{
					Code:     "RSC-500",
					FilePath: "OPS/toc.ncx",
					Message:  `Some placeholder error`,
					Location: &epubcheck.Position{Line: 35, Column: 10},
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "NCX-001",
					FilePath: "OPS/toc.ncx",
					Message:  `NCX identifier ("urn:uuid:1da9fa05e-dd8b-4be3-85ab-455656cc14f2") does not match OPF identifier ("").`,
				},
				{
					Code:     "RSC-500",
					FilePath: "OPS/toc.ncx",
					Message:  `Some placeholder error`,
					Location: &epubcheck.Position{Line: 36, Column: 10},
				},
			},
		},
		getContentByFileName: createTestCaseFileHandlerFunction(map[string]string{
			"OPS/content.opf": opfNoIdentifierOriginal,
			"OPS/toc.ncx":     ncxUuidIdentifier,
		}),
	},
}

/**
What should be tested here:
- Make sure each of the expected rules runs and makes the expected update (these will be single tests with edge cases allotted as well)
	- NCX-001: fix book id discrepancy between ncx and opf files
	- RSC 5:
	  - Invalid id:
		  - Should work on opf files
			- Should work on ncx files
			- Should work on html/xhtml files
			- Make sure that it increments or decrements anything that comes after it on the same line based on the change in characters from the rule
		- Invalid attribute:
		  - Should swap to the refines syntax in the opf file and work on multiple entries
		- Empty metadata property:
		  - Should remove the empty metadata property in opf files, remove any errors for that line, and decrement subsequent line number references on errors
		-
- Make sure that multiple rules play well together
*/

func TestHandleValidationErrors(t *testing.T) {
	for name, tc := range handleValidationErrorTestCases {
		t.Run(name, func(t *testing.T) {
			var nameToUpdatedFileContents = map[string]string{}
			err := epubcheck.HandleValidationErrors(tc.opfFolder, tc.ncxFilename, tc.opfFilename, nameToUpdatedFileContents, &tc.validationErrors, tc.getContentByFileName)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedErrorState, tc.validationErrors)

			for name, expectedContents := range tc.expectedFileState {
				actualContents, ok := nameToUpdatedFileContents[name]
				assert.True(t, ok, fmt.Sprintf("expected file %q to be updated, but it was not", name))
				assert.Equal(t, expectedContents, actualContents, fmt.Sprintf("expected file contents for %q did not match actual contents", name))
			}
		})
	}
}
