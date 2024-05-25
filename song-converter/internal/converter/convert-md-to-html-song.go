package converter

import (
	"fmt"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type SongMetadata struct {
	Melody         string `yaml:"melody"`
	SongKey        string `yaml:"key"`
	Authors        string `yaml:"authors"`
	InChurch       string `yaml:"in-church"`
	VerseReference string `yaml:"verse"`
	BookLocation   string `yaml:"location"`
	Copyright      string `yaml:"copyright"`
}

const (
	emptyColumnContent = "&nbsp;&nbsp;&nbsp;&nbsp;"
	closeMetadata      = "</div><br>"
)

var commonOrPotentialIssueFixer = strings.NewReplacer("\u00a0\u00a0\n", "<br>\n", "\\&", "&", "\n\n", "\n")

func ConvertMdToHtmlSong(filePath, fileContents string) (string, error) {
	var metadata SongMetadata
	mdContent, err := parseFrontmatter(filePath, fileContents, &metadata)
	if err != nil {
		return "", err
	}

	var metadataHtml = buildMetadataDiv(&metadata)
	html := mdToHTML([]byte(mdContent))
	html = replaceOtherTitle(html)
	html = strings.Replace(html, "</h1>\n", "</h1>\n"+metadataHtml, 1)

	return fmt.Sprintf("<div class=\"keep-together\">\n%s</div>\n<br>", commonOrPotentialIssueFixer.Replace(html)), nil
}

func mdToHTML(md []byte) string {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func buildMetadataDiv(metadata *SongMetadata) string {
	if metadata == nil {
		return ""
	}

	var metadataCount = 0
	var row1 = 0
	var row2 = 0

	metadataCount, row2 = updateCountsIfMetatdataExists(metadata.Melody, metadataCount, row2)
	metadataCount, row2 = updateCountsIfMetatdataExists(metadata.VerseReference, metadataCount, row2)
	metadataCount, row1 = updateCountsIfMetatdataExists(metadata.Authors, metadataCount, row1)
	metadataCount, row1 = updateCountsIfMetatdataExists(metadata.SongKey, metadataCount, row1)
	metadataCount, row1 = updateCountsIfMetatdataExists(metadata.BookLocation, metadataCount, row1)

	if metadataCount == 0 {
		return ""
	}

	var metadataHtml = strings.Builder{}
	metadataHtml.WriteString("<div>")

	var (
		addRowEntry = func(value, class, nonEmptyValue string) {
			metadataHtml.WriteString(fmt.Sprintf("<div><div class=%q>", class))

			if strings.Trim(value, "") != "" {
				metadataHtml.WriteString(nonEmptyValue)
			} else {
				metadataHtml.WriteString(emptyColumnContent)
			}

			metadataHtml.WriteString("</div></div>")
		}
		addBoldRowEntry = func(value, class string) {
			addRowEntry(value, class, fmt.Sprintf("<b>%s</b>", value))
		}
		addRegularRowEntry = func(value, class string) {
			addRowEntry(value, class, value)
		}
	)

	if row1 != 0 {
		if row2 != 0 {
			metadataHtml.WriteString("<div class=\"metadata row-padding\">")
		} else {
			metadataHtml.WriteString("<div class=\"metadata\">")
		}

		if strings.EqualFold(metadata.InChurch, "Y") {
			addBoldRowEntry(metadata.Authors, "author")
		} else {
			addRegularRowEntry(metadata.Authors, "author")
		}

		addBoldRowEntry(metadata.SongKey, "key")
		addRegularRowEntry(metadata.BookLocation, "location")

		metadataHtml.WriteString("</div>")
	}

	if row2 == 0 {
		metadataHtml.WriteString(closeMetadata)

		return metadataHtml.String()
	}

	metadataHtml.WriteString("<div class=\"metadata\">")
	if row2 == 1 && metadata.Melody != "" {
		addBoldRowEntry(metadata.Melody, "melody-75")
	} else {
		addBoldRowEntry(metadata.Melody, "melody")
		addRegularRowEntry(metadata.VerseReference, "verse")
	}

	metadataHtml.WriteString("</div>")
	metadataHtml.WriteString(closeMetadata)

	return metadataHtml.String()
}

func updateCountsIfMetatdataExists(value string, metadataElements, rowElements int) (int, int) {
	if (strings.Trim(value, "")) == "" {
		return metadataElements, rowElements
	}

	return metadataElements + 1, rowElements + 1
}

func parseFrontmatter(filePath, fileContents string, metadata *SongMetadata) (string, error) {
	restOfContent, err := frontmatter.Parse(strings.NewReader(fileContents), metadata)
	if err != nil {
		return "", fmt.Errorf(`there was an error getting the frontmatter for file '%s': %w`, filePath, err)
	}

	return string(restOfContent), nil
}

func replaceOtherTitle(html string) string {
	openH1 := strings.Index(html, "<h1")
	if openH1 == -1 {
		return html
	}

	var textAfterH1 = html[openH1:]
	openPara := strings.Index(textAfterH1, "(")
	if openPara == -1 {
		return html
	}

	closePara := strings.Index(textAfterH1, ")</h1")
	if closePara == -1 {
		return html
	}

	return html[:openPara] + "<span class=\"other-title\">" + html[openPara:closePara+1] + "</span>" + html[closePara+1:]
}
