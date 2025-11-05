//go:build unit

package epubhandler_test

import (
	"archive/zip"
	"testing"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	"github.com/stretchr/testify/assert"
)

func strSet(strs ...string) map[string]struct{} {
	s := make(map[string]struct{}, len(strs))
	for _, v := range strs {
		s[v] = struct{}{}
	}
	return s
}

type removeUnusedFilesTestCase struct {
	name              string
	zipFiles          map[string]*zip.File
	manifestFiles     map[string]struct{}
	removableFileExts []string
	expectedHandled   []string
}

func TestRemoveUnusedFiles(t *testing.T) {
	tests := map[string]removeUnusedFilesTestCase{
		"mimetype left alone": {
			zipFiles: map[string]*zip.File{
				"mimetype": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{},
		},
		"META-INF/container.xml left alone": {
			zipFiles: map[string]*zip.File{
				"META-INF/container.xml": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{},
		},
		"onix.xml left alone": {
			zipFiles: map[string]*zip.File{
				"OEBPS/onix.xml": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{},
		},
		"encryption.xml left alone": {
			zipFiles: map[string]*zip.File{
				"META-INF/encryption.xml": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{},
		},
		"file not listed in manifest or removable exts": {
			zipFiles: map[string]*zip.File{
				"OEBPS/keepme.txt": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{},
		},
		"manifest file should be ignored": {
			zipFiles: map[string]*zip.File{
				"OEBPS/content.xhtml": nil,
			},
			manifestFiles:     strSet("OEBPS/content.xhtml"),
			removableFileExts: []string{".xhtml"},
			expectedHandled:   []string{},
		},
		".jpg file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/image.jpg": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{"OEBPS/image.jpg"},
		},
		".jpeg file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/image.jpeg": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpeg"},
			expectedHandled:   []string{"OEBPS/image.jpeg"},
		},
		".png file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/image.png": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".png"},
			expectedHandled:   []string{"OEBPS/image.png"},
		},
		".gif file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/image.gif": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".gif"},
			expectedHandled:   []string{"OEBPS/image.gif"},
		},
		".bmp file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/image.bmp": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".bmp"},
			expectedHandled:   []string{"OEBPS/image.bmp"},
		},
		".js file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/script.js": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".js"},
			expectedHandled:   []string{"OEBPS/script.js"},
		},
		".html file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/index.html": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".html"},
			expectedHandled:   []string{"OEBPS/index.html"},
		},
		".htm file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/index.htm": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".htm"},
			expectedHandled:   []string{"OEBPS/index.htm"},
		},
		".xhtml file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/content.xhtml": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".xhtml"},
			expectedHandled:   []string{"OEBPS/content.xhtml"},
		},
		".txt file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/notes.txt": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".txt"},
			expectedHandled:   []string{"OEBPS/notes.txt"},
		},
		".css file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/style.css": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".css"},
			expectedHandled:   []string{"OEBPS/style.css"},
		},
		".xml file is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/data.xml": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".xml"},
			expectedHandled:   []string{"OEBPS/data.xml"},
		},
		"complex: manifest files mixed, rights.xml should not be removed, others handled": {
			zipFiles: map[string]*zip.File{
				"OEBPS/test.html":    nil,
				"OEBPS/keep.xhtml":   nil,
				"OEBPS/style.css":    nil,
				"OEBPS/script.js":    nil,
				"OEBPS/picture.jpeg": nil,
				"OEBPS/unused.html":  nil,
				"OEBPS/image1.png":   nil,
				"OEBPS/rights.xml":   nil,
			},
			manifestFiles: strSet(
				"OEBPS/test.html", "OEBPS/keep.xhtml", "OEBPS/style.css", "OEBPS/script.js", "OEBPS/picture.jpeg",
			),
			removableFileExts: []string{".html", ".xhtml", ".css", ".js", ".jpeg", ".png", ".xml"},
			expectedHandled: []string{
				"OEBPS/unused.html",
				"OEBPS/image1.png",
				"OEBPS/rights.xml",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := epubhandler.RemoveUnusedFiles(nil, tc.zipFiles, tc.manifestFiles, tc.removableFileExts, false)
			assert.ElementsMatch(t, tc.expectedHandled, got)
		})
	}
}
