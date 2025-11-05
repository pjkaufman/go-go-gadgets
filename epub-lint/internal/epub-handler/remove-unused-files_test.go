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
		"When dealing with mimetype, it should be left alone": {
			zipFiles: map[string]*zip.File{
				"mimetype": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{},
		},
		"When dealing with META-INF/container.xml, it should be left alone": {
			zipFiles: map[string]*zip.File{
				"META-INF/container.xml": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{},
		},
		"When dealing with an onix.xml, it should be left alone": {
			zipFiles: map[string]*zip.File{
				"OEBPS/onix.xml": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{},
		},
		"When dealing with an encryption.xml, it should be left alone": {
			zipFiles: map[string]*zip.File{
				"META-INF/encryption.xml": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{},
		},
		"When a file is not listed in the manifest or removable exts, it should be left alone": {
			zipFiles: map[string]*zip.File{
				"OEBPS/keepme.txt": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{},
		},
		"When a file is in the manifest, it should not be removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/content.xhtml": nil,
			},
			manifestFiles:     strSet("OEBPS/content.xhtml"),
			removableFileExts: []string{".xhtml"},
			expectedHandled:   []string{},
		},
		"When `.jpg` is in the file types to remove and a jpg file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/image.jpg": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpg"},
			expectedHandled:   []string{"OEBPS/image.jpg"},
		},
		"When `.jpeg` is in the file types to remove and a jpeg file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/image.jpeg": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".jpeg"},
			expectedHandled:   []string{"OEBPS/image.jpeg"},
		},
		"When `.png` is in the file types to remove and a png file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/image.png": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".png"},
			expectedHandled:   []string{"OEBPS/image.png"},
		},
		"When `.gif` is in the file types to remove and a gif file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/image.gif": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".gif"},
			expectedHandled:   []string{"OEBPS/image.gif"},
		},
		"When `.bmp` is in the file types to remove and a bmp file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/image.bmp": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".bmp"},
			expectedHandled:   []string{"OEBPS/image.bmp"},
		},
		"When `.js` is in the file types to remove and a js file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/script.js": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".js"},
			expectedHandled:   []string{"OEBPS/script.js"},
		},
		"When `.html` is in the file types to remove and an html file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/index.html": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".html"},
			expectedHandled:   []string{"OEBPS/index.html"},
		},
		"When `.htm` is in the file types to remove and an htm file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/index.htm": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".htm"},
			expectedHandled:   []string{"OEBPS/index.htm"},
		},
		"When `.xhtml` is in the file types to remove and an xhtml file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/content.xhtml": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".xhtml"},
			expectedHandled:   []string{"OEBPS/content.xhtml"},
		},
		"When `.txt` is in the file types to remove and a txt file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/notes.txt": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".txt"},
			expectedHandled:   []string{"OEBPS/notes.txt"},
		},
		"When `.css` is in the file types to remove and a css file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/style.css": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".css"},
			expectedHandled:   []string{"OEBPS/style.css"},
		},
		"When `.xml` is in the file types to remove and an xml file is not in the manifest, it is removed": {
			zipFiles: map[string]*zip.File{
				"OEBPS/data.xml": nil,
			},
			manifestFiles:     strSet(),
			removableFileExts: []string{".xml"},
			expectedHandled:   []string{"OEBPS/data.xml"},
		},
		"When some files are present in the manifest and others are not, then the proper files should be removed based on the remove file types": {
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
