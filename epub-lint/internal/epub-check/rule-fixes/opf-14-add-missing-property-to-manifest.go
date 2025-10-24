package rulefixes

import (
	"fmt"
	"strings"
)

const (
	manifestStartTag = "<manifest>"
	manifestEndTag   = "</manifest>"
)

var ErrNoManifest = fmt.Errorf("manifest tag not found in OPF contents")

func AddScriptedToManifest(opfContents, fileName, property string) (string, error) {
	startIndex, endIndex, manifestContent, err := getManifestContents(opfContents)
	if err != nil {
		return "", err
	}

	lines := strings.Split(manifestContent, "\n")
	for i, line := range lines {
		if strings.Contains(line, fmt.Sprintf(`href="%s"`, fileName)) {

			if strings.Contains(line, "properties=\"\"") {
				lines[i] = strings.Replace(line, `properties="`, `properties="`+property, 1)
			} else if strings.Contains(line, `properties="`) {
				if !strings.Contains(line, `scripted`) {
					lines[i] = strings.Replace(line, `properties="`, `properties="`+property+" ", 1)
				}
			} else {
				lines[i] = strings.Replace(line, `/>`, ` properties="`+property+`"/>`, 1)
			}

			break
		}
	}

	updatedManifestContent := strings.Join(lines, "\n")
	updatedOpfContents := opfContents[:startIndex+len(manifestStartTag)] + updatedManifestContent + opfContents[endIndex:]

	return updatedOpfContents, nil
}

func getManifestContents(opfContents string) (int, int, string, error) {
	startIndex := strings.Index(opfContents, manifestStartTag)
	endIndex := strings.Index(opfContents, manifestEndTag)

	if startIndex == -1 || endIndex == -1 {
		return 0, 0, "", ErrNoManifest
	}

	return startIndex, endIndex, opfContents[startIndex+len(manifestStartTag) : endIndex], nil
}
