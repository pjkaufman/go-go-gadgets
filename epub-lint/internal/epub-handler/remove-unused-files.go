package epubhandler

import (
	"archive/zip"
	"path"
	"slices"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func RemoveUnusedFiles(handledFiles []string, zipFiles map[string]*zip.File, manifestFiles map[string]struct{}, removableFileExts []string, verbose bool) []string {
	for filePath := range zipFiles {
		if _, exists := manifestFiles[filePath]; exists {
			continue
		}

		if strings.HasSuffix(filePath, "META-INF/container.xml") {
			continue
		}

		if hasFilename(filePath, "onix.xml") || hasFilename(filePath, "encryption.xml") {
			continue
		}

		if hasExt(removableFileExts, filePath) {
			// label file as handled despite not saving it to the destination
			handledFiles = append(handledFiles, filePath)

			if verbose {
				logger.WriteInfof("Removed file %q from the epub since it is not in the manifest.\n", filePath)
			}
		}
	}

	return handledFiles
}

func hasExt(slice []string, file string) bool {
	var ext = path.Ext(file)
	return slices.Contains(slice, ext)
}

func hasFilename(filePath, file string) bool {
	if !strings.HasSuffix(filePath, file) {
		return false
	}

	return path.Base(filePath) == file
}
