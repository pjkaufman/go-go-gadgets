package linter

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"

	epubhandler "github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/epub-handler"
)

func RemoveScriptedFromManifest(opfContents string, zipPath string) (string, error) {
	opfInfo, err := epubhandler.GetOpfXml(opfContents)
	if err != nil {
		return "", err
	}

	var fileName = filepath.Base(zipPath)
	for _, item := range opfInfo.Manifest.Items {
		if strings.HasSuffix(item.Href, fileName) {
			if item.Properties != nil && strings.Contains(*item.Properties, "scripted") {
				updatedVal := strings.TrimSpace(strings.Replace(*item.Properties, "scripted", "", 1))
				if updatedVal != "" {
					*item.Properties = updatedVal
				} else {
					item.Properties = nil
				}
			}

			break
		}
	}

	updatedOpfContents, err := xml.MarshalIndent(opfInfo, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated OPF contents: %v", err)
	}

	return string(updatedOpfContents), nil
}
