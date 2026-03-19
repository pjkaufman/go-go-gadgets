//go:build unit

package jnovels_test

import (
	"fmt"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/jnovels"
	"github.com/stretchr/testify/assert"
)

type cleanupJNovelsFilesTestCase struct {
	input                jnovels.JNovelsCleanupContext
	expectedFileContent  map[string]string
	expectedHandledFiles []string
}

var cleanupJNovelsFilesTestCases = map[string]cleanupJNovelsFilesTestCase{
	// When no JNovels files are present, no changes are made and no files are handled
	// When there is a JNovels html file present, it should be handled and removed from the OPF file
	// When there is a JNovels image file present, it should be handled and removed from the OPF file
	// When there are the JNovels html and image files present, they should be handled and removed from the OPF file
	// When there is a JNovels html file present and it is referenced from the nav file, it should be handled and removed from the OPF and nav files
	// When there is a JNovels html file present and it is referenced from the nav and NCX files, it should be handled and removed from the OPF, NCX, and nav files
	// When there is a JNovels image file present and it is referenced as a part of the landmarks as the cover and there is a cover file, it should be handled, removed from the OPF file, and update to the cover file in the nav file
	// When there is a JNovels image file present and it is referenced as a part of the landmarks as the toc and there is a toc file, it should be handled, removed from the OPF file, and update to the toc file in the nav file
	// When there is a JNovels image file present and it is referenced as a part of the landmarks as the cover and toc and there is a cover and toc file, it should be handled, removed from the OPF file, and update to the cover and toc files in the nav file
}

func TestCleanupJNovelsFiles(t *testing.T) {
	for name, tc := range cleanupJNovelsFilesTestCases {
		t.Run(name, func(t *testing.T) {
			handledFiles, err := jnovels.CleanupJNovelsFiles(tc.input)

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedHandledFiles, handledFiles)
			for name, fileContent := range tc.expectedFileContent {
				updatedContents, found := tc.input.UpdatedFileContents[name]
				if !found {
					assert.Fail(t, fmt.Sprintf("expected %q to be updated, but it was not", name))
				} else {
					assert.Equal(t, fileContent, updatedContents, "%q does not match expected content")
				}
			}
		})
	}
}
