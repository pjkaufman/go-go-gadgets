package rulefixes

import (
	"fmt"
	"strings"

	epubhandler "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-handler"
)

var ErrNoManifest = fmt.Errorf("manifest tag not found in OPF contents")

func AddPropertyToManifest(opfContents, fileName, property string) (string, error) {
	startIndex, endIndex, manifestContent, err := epubhandler.GetManifestContents(opfContents)
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
	updatedOpfContents := opfContents[:startIndex+len(epubhandler.ManifestStartTag)] + updatedManifestContent + opfContents[endIndex:]

	return updatedOpfContents, nil
}
