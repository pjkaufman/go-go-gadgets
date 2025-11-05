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
	opfFolder, ncxFilename, opfFilename           string
	validationErrors                              epubcheck.ValidationErrors
	expectedErrorState                            epubcheck.ValidationErrors
	expectedFileState, validFilesToInitialContent map[string]string // filename to content
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
	//go:embed testdata/ncx-1/number-identifier.opf
	opfNumberIdentifierOriginal string
	//go:embed testdata/ncx-1/number-identifier_updated.opf
	opfNumberIdentifierExpected string
	//go:embed testdata/opf-30/missing-unique-identifier.opf
	opfMissingUniqueIdentifierOriginal string
	//go:embed testdata/opf-30/missing-unique-identifier_updated.opf
	opfMissingUniqueIdentifierExpected string
	//go:embed testdata/rsc-5/invalid-id.html
	htmlInvalidIdOriginal string
	//go:embed testdata/rsc-5/invalid-id_updated.html
	htmlInvalidIdExpected string
	//go:embed testdata/rsc-5/invalid-id.opf
	opfInvalidIdOriginal string
	//go:embed testdata/rsc-5/invalid-id_updated.opf
	opfInvalidIdExpected string
	//go:embed testdata/rsc-5/duplicate-ids.html
	htmlDuplicateIdsOriginal string
	//go:embed testdata/rsc-5/duplicate-ids_updated.html
	htmlDuplicateIdsExpected string
	//go:embed testdata/rsc-5/missing-image-alt.html
	htmlMissingImageAltOriginal string
	//go:embed testdata/rsc-5/missing-image-alt_updated.html
	htmlMissingImageAltExpected string
	//go:embed testdata/rsc-5/fixable-blockquote.html
	htmlFixableBlockquoteOriginal string
	//go:embed testdata/rsc-5/fixable-blockquote_updated.html
	htmlFixableBlockquoteExpected string
	//go:embed testdata/rsc-5/duplicate-play-order.ncx
	ncxDuplicatePlayOrderOriginal string
	//go:embed testdata/rsc-5/duplicate-play-order_updated.ncx
	ncxDuplicatePlayOrderExpected string
	//go:embed testdata/rsc-5/empty-elements.opf
	opfEmptyElementsOriginal string
	//go:embed testdata/rsc-5/empty-elements_updated.opf
	opfEmptyElementsExpected string
)

func createTestCaseFileHandlerFunction(validFilesToContent map[string]string, currentContents map[string]string) func(string) (string, error) {
	return func(s string) (string, error) {
		if content, ok := currentContents[s]; ok {
			return content, nil
		} else if content, ok := validFilesToContent[s]; ok {
			return content, nil
		}

		return "", fmt.Errorf("unexpected attempt to get file contents for file %q", s)
	}
}

var handleValidationErrorTestCases = map[string]handleValidationErrorTestCase{
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
		validFilesToInitialContent: map[string]string{
			"OPS/content.opf": opfAddPropertiesExpected,
		},
	},
	"OPF 15: Removing properties from different files should work without issue": {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/content.opf": opfRemovePropertiesExpected},
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
		validFilesToInitialContent: map[string]string{
			"OPS/content.opf": opfRemovePropertiesOriginal,
		},
	},
	"OPF 30: When the unique-identifier property does not match any existing identifiers, add the unique identifier to the first identifier without an id": {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/content.opf": opfMissingUniqueIdentifierExpected},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "OPF-030",
					FilePath: "OOPS/content.opf",
					Message:  `The unique-identifier "BookId" was not found`,
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "OPF-030",
					FilePath: "OOPS/content.opf",
					Message:  `The unique-identifier "BookId" was not found`,
				},
			},
		},
		validFilesToInitialContent: map[string]string{
			"OPS/content.opf": opfMissingUniqueIdentifierOriginal,
		},
	},
	"NCX 1: When no identifier is present in the OPF, add the one from the NCX file": {
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
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "NCX-001",
					FilePath: "OPS/toc.ncx",
					Message:  `NCX identifier ("urn:uuid:1da9fa05e-dd8b-4be3-85ab-455656cc14f2") does not match OPF identifier ("").`,
				},
			},
		},
		validFilesToInitialContent: map[string]string{
			"OPS/content.opf": opfNoIdentifierOriginal,
			"OPS/toc.ncx":     ncxUuidIdentifier,
		},
	},
	"NCX 1: When an identifier is present in the OPF and differs from the one in the NCX, add the one from the NCX file": {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/content.opf": opfNumberIdentifierExpected},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "NCX-001",
					FilePath: "OPS/toc.ncx",
					Message:  `NCX identifier ("urn:uuid:1da9fa05e-dd8b-4be3-85ab-455656cc14f2") does not match OPF identifier ("1234").`,
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "NCX-001",
					FilePath: "OPS/toc.ncx",
					Message:  `NCX identifier ("urn:uuid:1da9fa05e-dd8b-4be3-85ab-455656cc14f2") does not match OPF identifier ("1234").`,
				},
			},
		},
		validFilesToInitialContent: map[string]string{
			"OPS/content.opf": opfNumberIdentifierOriginal,
			"OPS/toc.ncx":     ncxUuidIdentifier,
		},
	},
	"RSC 5: When there is an invalid id in a html/xhtml file it should get fixed": {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/Text/chapter1.html": htmlInvalidIdExpected},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/chapter1.html",
					// The error is the one for epub 3, but it should be fine handling it the same as an epub 2 one since epub 2 is more restrictive
					Message: `Error while parsing file: value of attribute "id" is invalid; must be a string matching the regular expression "[^\s]+"`,
					Location: &epubcheck.Position{
						Line:   16,
						Column: 40,
					},
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/chapter1.html",
					Message:  `Error while parsing file: value of attribute "id" is invalid; must be a string matching the regular expression "[^\s]+"`,
					Location: &epubcheck.Position{
						Line:   16,
						Column: 40,
					},
				},
			},
		},
		validFilesToInitialContent: map[string]string{
			"OPS/Text/chapter1.html": htmlInvalidIdOriginal,
		},
	},
	"RSC 5: When there is an invalid id in an opf file it should get fixed": {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/content.opf": opfInvalidIdExpected},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/content.opf",
					// The error is the one for epub 3, but it should be fine handling it the same as an epub 2 one since epub 2 is more restrictive
					Message: `Error while parsing file: value of attribute "idref" is invalid; must be a string matching the regular expression "[^\s]+"`,
					Location: &epubcheck.Position{
						Line:   77,
						Column: 29,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/content.opf",
					// The error is the one for epub 3, but it should be fine handling it the same as an epub 2 one since epub 2 is more restrictive
					Message: `Error while parsing file: value of attribute "id" is invalid; must be a string matching the regular expression "[^\s]+"`,
					Location: &epubcheck.Position{
						Line:   15,
						Column: 21,
					},
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/content.opf",
					// The error is the one for epub 3, but it should be fine handling it the same as an epub 2 one since epub 2 is more restrictive
					Message: `Error while parsing file: value of attribute "idref" is invalid; must be a string matching the regular expression "[^\s]+"`,
					Location: &epubcheck.Position{
						Line:   77,
						Column: 29,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/content.opf",
					// The error is the one for epub 3, but it should be fine handling it the same as an epub 2 one since epub 2 is more restrictive
					Message: `Error while parsing file: value of attribute "id" is invalid; must be a string matching the regular expression "[^\s]+"`,
					Location: &epubcheck.Position{
						Line:   15,
						Column: 21,
					},
				},
			},
		},
		validFilesToInitialContent: map[string]string{
			"OPS/content.opf": opfInvalidIdOriginal,
		},
	},
	"RSC 5: When there are duplicate ids in a file, they should be updated accordingly": {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/Text/prologue.html": htmlDuplicateIdsExpected},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/prologue.html",
					Message:  `Error while parsing file: Duplicate ID "line"`,
					Location: &epubcheck.Position{
						Line:   17,
						Column: 41,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/prologue.html",
					Message:  `Error while parsing file: Duplicate ID "auto_bookmark_toc_9"`,
					Location: &epubcheck.Position{
						Line:   14,
						Column: 29,
					},
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/prologue.html",
					Message:  `Error while parsing file: Duplicate ID "line"`,
					Location: &epubcheck.Position{
						Line:   17,
						Column: 41,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/prologue.html",
					Message:  `Error while parsing file: Duplicate ID "auto_bookmark_toc_9"`,
					Location: &epubcheck.Position{
						Line:   14,
						Column: 29,
					},
				},
			},
		},
		validFilesToInitialContent: map[string]string{
			"OPS/Text/prologue.html": htmlDuplicateIdsOriginal,
		},
	},
	`RSC 5: When an image is missing its "alt" attribute, an empty one should be added`: {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/Text/frontmatter.html": htmlMissingImageAltExpected},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/frontmatter.html",
					Message:  `Error while parsing file: element "img" missing required attribute "alt"`,
					Location: &epubcheck.Position{
						Line:   15,
						Column: 29,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/frontmatter.html",
					Message:  `Error while parsing file: element "img" missing required attribute "alt"`,
					Location: &epubcheck.Position{
						Line:   14,
						Column: 28,
					},
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/frontmatter.html",
					Message:  `Error while parsing file: element "img" missing required attribute "alt"`,
					Location: &epubcheck.Position{
						Line:   15,
						Column: 29,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/frontmatter.html",
					Message:  `Error while parsing file: element "img" missing required attribute "alt"`,
					Location: &epubcheck.Position{
						Line:   14,
						Column: 28,
					},
				},
			},
		},
		validFilesToInitialContent: map[string]string{
			"OPS/Text/frontmatter.html": htmlMissingImageAltOriginal,
		},
	},
	`RSC 5: When a blockquote fails to be parsed, if the blockquote ends in a closing span element, a self-closing tag, or just text, then it should be updated to have a paragraph tag inserted at the start and end`: {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/Text/content.html": htmlFixableBlockquoteExpected},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/content.html",
					Message:  `Error while parsing file: element "blockquote" incomplete;`,
					Location: &epubcheck.Position{
						Line:   16,
						Column: 58,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/content.html",
					Message:  `Error while parsing file: element "blockquote" incomplete;`,
					Location: &epubcheck.Position{
						Line:   15,
						Column: 51,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/content.html",
					Message:  `Error while parsing file: element "blockquote" incomplete;`,
					Location: &epubcheck.Position{
						Line:   14,
						Column: 60,
					},
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/content.html",
					Message:  `Error while parsing file: element "blockquote" incomplete;`,
					Location: &epubcheck.Position{
						Line:   16,
						Column: 58,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/content.html",
					Message:  `Error while parsing file: element "blockquote" incomplete;`,
					Location: &epubcheck.Position{
						Line:   15,
						Column: 51,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/Text/content.html",
					Message:  `Error while parsing file: element "blockquote" incomplete;`,
					Location: &epubcheck.Position{
						Line:   14,
						Column: 60,
					},
				},
			},
		},
		validFilesToInitialContent: map[string]string{
			"OPS/Text/content.html": htmlFixableBlockquoteOriginal,
		},
	},
	`RSC 5: When a play order is found to be identical/incorrect, the play order should be updated to be back in order`: {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/toc.ncx": ncxDuplicatePlayOrderExpected},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/toc.ncx",
					Message:  `Error while parsing file: identical playOrder values for navPoint/navTarget/pageTarget that do not refer to same target`,
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/toc.ncx",
					Message:  `Error while parsing file: identical playOrder values for navPoint/navTarget/pageTarget that do not refer to same target`,
				},
			},
		},
		validFilesToInitialContent: map[string]string{
			"OPS/toc.ncx": ncxDuplicatePlayOrderOriginal,
		},
	},
	`RSC 5: When an OPF element is empty it should be removed and all errors on the same line should be removed and other lines after it decremented`: {
		opfFolder:         "OPS",
		opfFilename:       "OPS/content.opf",
		ncxFilename:       "OPS/toc.ncx",
		expectedFileState: map[string]string{"OPS/content.opf": opfEmptyElementsExpected},
		validationErrors: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-005",
					FilePath: "OPS/content.opf",
					Message:  `Error while parsing file: character content of element "dc:identifier"`,
					Location: &epubcheck.Position{
						Line:   7,
						Column: 36,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/content.opf",
					Message:  `Error while parsing file: character content of element "dc:creator"`,
					Location: &epubcheck.Position{
						Line:   12,
						Column: 45,
					},
				},
				{
					Code:     "RSC-005",
					FilePath: "OPS/content.opf",
					Message:  `Error while parsing file: attribute "opf:role"`,
					Location: &epubcheck.Position{
						Line:   12,
						Column: 32,
					},
				},
				{
					Code:     "RSC-999",
					FilePath: "OPS/content.opf",
					Message:  `Some error here..."`,
					Location: &epubcheck.Position{
						Line:   50,
						Column: 35,
					},
				},
			},
		},
		expectedErrorState: epubcheck.ValidationErrors{
			ValidationIssues: []epubcheck.ValidationError{
				{
					Code:     "RSC-999",
					FilePath: "OPS/content.opf",
					Message:  `Some error here..."`,
					Location: &epubcheck.Position{
						Line:   48,
						Column: 35,
					},
				},
			},
		},
		validFilesToInitialContent: map[string]string{
			"OPS/content.opf": opfEmptyElementsOriginal,
		},
	},
}

func TestHandleValidationErrors(t *testing.T) {
	for name, tc := range handleValidationErrorTestCases {
		t.Run(name, func(t *testing.T) {
			var nameToUpdatedFileContents = map[string]string{}
			err := epubcheck.HandleValidationErrors(tc.opfFolder, tc.ncxFilename, tc.opfFilename, nameToUpdatedFileContents, &tc.validationErrors, createTestCaseFileHandlerFunction(tc.validFilesToInitialContent, nameToUpdatedFileContents))
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
