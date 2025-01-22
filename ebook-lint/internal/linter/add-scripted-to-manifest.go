package linter

import (
	"fmt"
	"strings"
)

const (
	startTag = "<manifest>"
	endTag   = "</manifest>"
)

func AddScriptedToManifest(opfContents string, fileName string) (string, error) {
	startIndex, endIndex, manifestContent, err := getManifestContents(opfContents)
	if err != nil {
		return "", err
	}

	lines := strings.Split(manifestContent, "\n")
	for i, line := range lines {
		if strings.Contains(line, fmt.Sprintf(`href="%s"`, fileName)) {

			if strings.Contains(line, "properties=\"\"") {
				lines[i] = strings.Replace(line, `properties="`, `properties="scripted`, 1)
			} else if strings.Contains(line, `properties="`) {
				if !strings.Contains(line, `scripted`) {
					lines[i] = strings.Replace(line, `properties="`, `properties="scripted `, 1)
				}
			} else {
				lines[i] = strings.Replace(line, `/>`, ` properties="scripted"/>`, 1)
			}

			break
		}
	}

	updatedManifestContent := strings.Join(lines, "\n")
	updatedOpfContents := opfContents[:startIndex+len(startTag)] + updatedManifestContent + opfContents[endIndex:]

	return updatedOpfContents, nil
}

func getManifestContents(opfContents string) (int, int, string, error) {
	startIndex := strings.Index(opfContents, startTag)
	endIndex := strings.Index(opfContents, endTag)

	if startIndex == -1 || endIndex == -1 {
		return 0, 0, "", fmt.Errorf("manifest tag not found in OPF contents")
	}

	return startIndex, endIndex, opfContents[startIndex+len(startTag) : endIndex], nil
}
