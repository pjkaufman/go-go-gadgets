package jnovels

import (
	"fmt"
	"path/filepath"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
)

type JNovelsCleanupContext struct {
	EpubInfo            epubhandler.EpubInfo
	OpfFolder           string
	OpfFileName         string
	NcxFileName         string
	FileBasenameMap     map[string][]string
	UpdatedFileContents map[string]string
	GetFileContents     func(string) (string, error)
}

func CleanupJNovelsFiles(ctx JNovelsCleanupContext) ([]string, error) {
	htmlHandledFiles, err := cleanupJNovelsXhtml(ctx)
	if err != nil {
		return nil, err
	}

	var imageHandledFiles []string
	imageHandledFiles, err = cleanupJNovelsImage(ctx)
	if err != nil {
		return nil, err
	}

	return append(htmlHandledFiles, imageHandledFiles...), nil
}

// Removes jnovels.xhtml from OPF, NAV, and TOC files
func cleanupJNovelsXhtml(ctx JNovelsCleanupContext) ([]string, error) {
	var handledFiles []string

	for _, filename := range ctx.FileBasenameMap[JnovelsFile] {
		handledFiles = append(handledFiles, filename)

		updatedOpfContents, err := epubhandler.RemoveFileFromOpf(
			ctx.UpdatedFileContents[ctx.OpfFileName],
			JnovelsFile,
		)
		if err != nil {
			return handledFiles, fmt.Errorf("Failed to remove file %q from opf: %w", filename, err)
		}

		ctx.UpdatedFileContents[ctx.OpfFileName] = updatedOpfContents

		err = removeFromNavAndTocFiles(ctx, JnovelsFile)
		if err != nil {
			return handledFiles, err
		}

		ncxFolderPath := filepath.Dir(ctx.NcxFileName) // used instead of the file path as that results in an additional "../" being added
		var relativeFilePath string
		relativeFilePath, err = filepath.Rel(ncxFolderPath, filename)
		if err != nil {
			return handledFiles, fmt.Errorf("Failed to determine relative path between ncx file %q and file %q: %w", ctx.NcxFileName, filename, err)
		}

		var (
			priorNcx   = ctx.UpdatedFileContents[ctx.NcxFileName]
			updatedNcx = epubhandler.RemoveFileFromNcx(ctx.UpdatedFileContents[ctx.NcxFileName], relativeFilePath)
		)
		if priorNcx != updatedNcx {
			edits := rulefixes.FixPlayOrder(updatedNcx)
			if len(edits) != 0 {
				updatedNcx, err = positions.ApplyEdits(ctx.NcxFileName, updatedNcx, edits)
				if err != nil {
					return handledFiles, err
				}
			}
		}

		ctx.UpdatedFileContents[ctx.NcxFileName] = updatedNcx
	}

	return handledFiles, nil
}

// Removes the image file and updates landmarks in NAV file
func cleanupJNovelsImage(ctx JNovelsCleanupContext) ([]string, error) {
	var handledFiles []string

	for _, filename := range ctx.FileBasenameMap[JnovelsImage] {
		handledFiles = append(handledFiles, filename)

		updatedOpfContents, err := epubhandler.RemoveFileFromOpf(
			ctx.UpdatedFileContents[ctx.OpfFileName],
			JnovelsImage,
		)
		if err != nil {
			return handledFiles, fmt.Errorf("Failed to remove file %q from opf: %w", filename, err)
		}

		ctx.UpdatedFileContents[ctx.OpfFileName] = updatedOpfContents

		if ctx.EpubInfo.NavFile != "" && (ctx.EpubInfo.TocFile != "" || ctx.EpubInfo.CoverFile != "") {
			var (
				filePath      = filehandler.JoinPath(ctx.OpfFolder, ctx.EpubInfo.NavFile)
				navFolderPath = filepath.Dir(filePath) // used instead of the file path as that results in an additional "../" being added
			)
			contents, err := ctx.GetFileContents(filePath)
			if err != nil {
				return handledFiles, err
			}

			var relativeImagePath string
			relativeImagePath, err = filepath.Rel(navFolderPath, filename)
			if err != nil {
				return handledFiles, fmt.Errorf("Failed to determine relative path between nav file %q and file %q: %w", filePath, filename, err)
			}

			var (
				relativeCoverPath string
				coverPath         string
			)
			if ctx.EpubInfo.CoverFile != "" {
				coverPath = filehandler.JoinPath(ctx.OpfFolder, ctx.EpubInfo.CoverFile)

				relativeCoverPath, err = filepath.Rel(navFolderPath, coverPath)
				if err != nil {
					return handledFiles, fmt.Errorf("Failed to determine relative path between nav file %q and file %q: %w", filePath, coverPath, err)
				}
			}

			var (
				relativeTocPath string
				tocPath         string
			)
			if ctx.EpubInfo.TocFile != "" {
				tocPath = filehandler.JoinPath(ctx.OpfFolder, ctx.EpubInfo.TocFile)

				relativeTocPath, err = filepath.Rel(navFolderPath, tocPath)
				if err != nil {
					return handledFiles, fmt.Errorf("Failed to determine relative path between nav file %q and file %q: %w", filePath, tocPath, err)
				}
			}

			updatedNavContents := epubhandler.UpdateLandmarks(contents, relativeImagePath, relativeCoverPath, relativeTocPath)
			ctx.UpdatedFileContents[filePath] = updatedNavContents
		}
	}

	return handledFiles, nil
}

func removeFromNavAndTocFiles(ctx JNovelsCleanupContext, fileName string) error {
	if ctx.EpubInfo.NavFile != "" {
		filePath := filehandler.JoinPath(ctx.OpfFolder, ctx.EpubInfo.NavFile)
		contents, err := ctx.GetFileContents(filePath)
		if err != nil {
			return err
		}

		ctx.UpdatedFileContents[filePath] = epubhandler.RemoveFileFromNav(contents, fileName)
	}

	if ctx.EpubInfo.TocFile != "" && ctx.EpubInfo.NavFile != ctx.EpubInfo.TocFile {
		filePath := filehandler.JoinPath(ctx.OpfFolder, ctx.EpubInfo.TocFile)
		contents, err := ctx.GetFileContents(filePath)
		if err != nil {
			return err
		}

		ctx.UpdatedFileContents[filePath] = epubhandler.RemoveFileFromNav(contents, fileName)
	}

	return nil
}
