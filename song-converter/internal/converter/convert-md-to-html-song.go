package converter

import (
	"fmt"
	"regexp"
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

var otherTitleRegex = regexp.MustCompile(`^(<h1.*)\((.*)\)<(.*)`)

func ConvertMdToHtmlSong(filePath, fileContents string) (string, error) {
	var metadata SongMetadata
	mdContent, err := parseFrontmatter(filePath, fileContents, &metadata)
	if err != nil {
		return "", err
	}

	var metadataHtml = buildMetadataDiv(&metadata)
	html := mdToHTML([]byte(mdContent))
	html = otherTitleRegex.ReplaceAllString(html, `${1}<span class="other-title">(${2})</span><${3}`)

	// just in case we encounter this scenario where non-breaking space is encoded as its unicode value
	html = strings.ReplaceAll(html, "\u00a0\u00a0\n", "<br>\n")
	html = strings.ReplaceAll(html, "\\&", "&")
	html = strings.Replace(html, "</h1>\n", "</h1>\n"+metadataHtml, 1)
	html = strings.ReplaceAll(html, "\n\n", "\n")

	return fmt.Sprintf("<div class=\"keep-together\">\n%s</div>\n<br>", html), nil
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
			metadataHtml.WriteString(fmt.Sprintf("<div><div class=\"%s\">", class))

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
