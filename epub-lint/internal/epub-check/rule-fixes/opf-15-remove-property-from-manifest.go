package rulefixes

import (
	"fmt"
	"strings"
)

func RemovePropertyFromManifest(opfContents, fileName, property string) (string, error) {
	startIndex, endIndex, manifestContent, err := getManifestContents(opfContents)
	if err != nil {
		return "", err
	}

	var (
		lines    = strings.Split(manifestContent, "\n")
		fileHref = fmt.Sprintf(`href="%s"`, fileName)
	)
	for i, line := range lines {
		if strings.Contains(line, fileHref) {
			propStart := strings.Index(line, `properties="`)
			if propStart != -1 {
				valStart := propStart + len(`properties="`)
				propEnd := strings.Index(line[valStart:], `"`) + valStart
				properties := line[valStart:propEnd]

				if strings.Contains(properties, property) {
					updatedProperties := strings.Replace(properties, property, "", 1)
					updatedProperties = strings.TrimSpace(updatedProperties)

					if updatedProperties == "" {
						// Remove the properties attribute if it's empty plus the space right before it
						lines[i] = line[:propStart-1] + line[propEnd+1:]
					} else {
						lines[i] = line[:propStart] + `properties="` + updatedProperties + line[propEnd:]
					}
				}
			}
			break
		}
	}

	updatedManifestContent := strings.Join(lines, "\n")
	updatedOpfContents := opfContents[:startIndex+len(manifestStartTag)] + updatedManifestContent + opfContents[endIndex:]

	return updatedOpfContents, nil
}
