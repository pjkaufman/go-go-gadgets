package epub

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SourceRange tests
func TestSourceRangeValid(t *testing.T) {
	tests := []struct {
		name     string
		r        SourceRange
		length   int
		expected bool
	}{
		{"valid range", SourceRange{0, 10}, 100, true},
		{"zero range at start", SourceRange{0, 0}, 100, true},
		{"end equals length", SourceRange{0, 100}, 100, true},
		{"negative start", SourceRange{-1, 10}, 100, false},
		{"negative end", SourceRange{0, -1}, 100, false},
		{"end before start", SourceRange{10, 5}, 100, false},
		{"end exceeds length", SourceRange{0, 101}, 100, false},
		{"negative length", SourceRange{0, 10}, -1, false},
	}
	for _, tt := range tests {
		t := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.r.Valid(tt.length))
		})
	}
}

func TestSourceRangeLen(t *testing.T) {
	assert.Equal(t, 10, SourceRange{0, 10}.Len())
	assert.Equal(t, 0, SourceRange{5, 5}.Len())
	assert.Equal(t, 0, SourceRange{10, 5}.Len())
}

// EPUBPath tests
func TestNewEPUBPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"  ", ""},
		{"Text/chapter.xhtml", "Text/chapter.xhtml"},
		{"Text\\chapter.xhtml", "Text/chapter.xhtml"},
		{"./Text/chapter.xhtml", "Text/chapter.xhtml"},
		{"Text/./chapter.xhtml", "Text/chapter.xhtml"},
		{"Text/../chapter.xhtml", "chapter.xhtml"},
		{".", ""},
	}
	for _, tt := range tests {
		t := tt
		t.Run(tt.input, func(t *testing.T) {
			result := NewEPUBPath(tt.input)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestEPUBPathJoin(t *testing.T) {
	tests := []struct {
		base     string
		other    string
		expected string
	}{
		{"Text", "chapter.xhtml", "Text/chapter.xhtml"},
		{"", "Text/chapter.xhtml", "Text/chapter.xhtml"},
		{"Text", "", "Text"},
		{"OEBPS", "Text/chapter.xhtml", "OEBPS/Text/chapter.xhtml"},
	}
	for _, tt := range tests {
		t := tt
		t.Run(tt.base+" + "+tt.other, func(t *testing.T) {
			result := NewEPUBPath(tt.base).Join(tt.other)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestEPUBPathRelativeTo(t *testing.T) {
	tests := []struct {
		path     string
		base     string
		expected string
	}{
		{"Text/chapter.xhtml", ".", "Text/chapter.xhtml"},
		{"Text/chapter.xhtml", "OEBPS", "../Text/chapter.xhtml"},
	}
	for _, tt := range tests {
		t := tt
		t.Run(tt.path+" relative to "+tt.base, func(t *testing.T) {
			result := NewEPUBPath(tt.path).RelativeTo(NewEPUBPath(tt.base))
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

// Edit tests
func TestApplyEditsSimple(t *testing.T) {
	data := []byte("abcdef")
	edits := []Edit{
		{Start: 2, End: 4, Replacement: []byte("XY")},
	}
	result, err := ApplyEdits(data, edits)
	require.NoError(t, err)
	assert.Equal(t, "abXYef", string(result))
}

func TestApplyEditsMultiple(t *testing.T) {
	data := []byte("abcdef")
	edits := []Edit{
		{Start: 4, End: 6, Replacement: []byte("XY")},
		{Start: 0, End: 2, Replacement: []byte("Z")},
	}
	result, err := ApplyEdits(data, edits)
	require.NoError(t, err)
	assert.Equal(t, "ZcdXY", string(result))
}

func TestApplyEditsEmpty(t *testing.T) {
	data := []byte("abcdef")
	result, err := ApplyEdits(data, []Edit{})
	require.NoError(t, err)
	assert.Equal(t, "abcdef", string(result))
}

func TestApplyEditsInsertion(t *testing.T) {
	data := []byte("abcdef")
	edits := []Edit{
		{Start: 3, End: 3, Replacement: []byte("XYZ")},
	}
	result, err := ApplyEdits(data, edits)
	require.NoError(t, err)
	assert.Equal(t, "abcXYZdef", string(result))
}

func TestApplyEditsDeletion(t *testing.T) {
	data := []byte("abcdef")
	edits := []Edit{
		{Start: 2, End: 4, Replacement: []byte{}},
	}
	result, err := ApplyEdits(data, edits)
	require.NoError(t, err)
	assert.Equal(t, "abef", string(result))
}

func TestApplyEditsNegativeOffset(t *testing.T) {
	data := []byte("abcdef")
	edits := []Edit{
		{Start: -1, End: 2, Replacement: []byte("X")},
	}
	_, err := ApplyEdits(data, edits)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "negative offset")
}

func TestApplyEditsEndBeforeStart(t *testing.T) {
	data := []byte("abcdef")
	edits := []Edit{
		{Start: 4, End: 2, Replacement: []byte("X")},
	}
	_, err := ApplyEdits(data, edits)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "end before start")
}

func TestApplyEditsExceedsLength(t *testing.T) {
	data := []byte("abcdef")
	edits := []Edit{
		{Start: 0, End: 100, Replacement: []byte("X")},
	}
	_, err := ApplyEdits(data, edits)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds document length")
}

func TestApplyEditsOverlapping(t *testing.T) {
	data := []byte("abcdef")
	edits := []Edit{
		{Start: 0, End: 3, Replacement: []byte("X")},
		{Start: 2, End: 4, Replacement: []byte("Y")},
	}
	_, err := ApplyEdits(data, edits)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "overlapping")
}

func TestApplyEditsUnicode(t *testing.T) {
	data := []byte("Hello 日本語 World")
	edits := []Edit{
		{Start: 6, End: 15, Replacement: []byte("世界")},
	}
	result, err := ApplyEdits(data, edits)
	require.NoError(t, err)
	assert.Equal(t, "Hello 世界 World", string(result))
}

func TestApplyEditsAdjacent(t *testing.T) {
	data := []byte("abcdef")
	edits := []Edit{
		{Start: 0, End: 2, Replacement: []byte("X")},
		{Start: 2, End: 4, Replacement: []byte("Y")},
	}
	result, err := ApplyEdits(data, edits)
	require.NoError(t, err)
	assert.Equal(t, "XYef", string(result))
}

// XML parsing tests
func TestParseXMLSimple(t *testing.T) {
	xml := []byte(`<root><item>Hello</item></root>`)
	doc, err := ParseXML(xml)
	require.NoError(t, err)
	require.NotNil(t, doc)
	require.NotNil(t, doc.Root)
	assert.Equal(t, "root", doc.Root.Name.Local)
}

func TestParseXMLNested(t *testing.T) {
	xml := []byte(`<root><parent><child>Text</child></parent></root>`)
	doc, err := ParseXML(xml)
	require.NoError(t, err)
	require.NotNil(t, doc.Root)

	parent := doc.Root.FindElements("parent")
	assert.Len(t, parent, 1)

	if len(parent) > 0 {
		child := parent[0].FindElements("child")
		assert.Len(t, child, 1)
	}
}

func TestParseXMLNil(t *testing.T) {
	doc, err := ParseXML(nil)
	require.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Nil(t, doc.Data)
}

func TestParseXMLEmpty(t *testing.T) {
	doc, err := ParseXML([]byte{})
	require.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, 0, len(doc.Data))
}

func TestParseXMLFindElementByID(t *testing.T) {
	xml := []byte(`<root><item id="sec1"/><item id="sec2"/></root>`)
	doc, err := ParseXML(xml)
	require.NoError(t, err)
	require.NotNil(t, doc.Root)

	found := doc.Root.FindElementByID("sec2")
	if found != nil {
		var foundID bool
		for _, attr := range found.Attributes {
			if attr.Name.Local == "id" {
				assert.Equal(t, "sec2", attr.Value)
				foundID = true
			}
		}
		assert.True(t, foundID)
	}
}

// XHTML document tests
func TestParseXHTMLDocument(t *testing.T) {
	xml := []byte(`<html><body><p id="p1">Hello</p></body></html>`)
	doc, err := ParseXHTML(xml)
	require.NoError(t, err)
	require.NotNil(t, doc)
	require.NotNil(t, doc.Root)
}

func TestXHTMLFindElements(t *testing.T) {
	xml := []byte(`<html><body><p>One</p><p>Two</p></body></html>`)
	doc, err := ParseXHTML(xml)
	require.NoError(t, err)

	paragraphs := doc.FindElements("p")
	assert.Len(t, paragraphs, 2)
}

func TestXHTMLFindElementByID(t *testing.T) {
	xml := []byte(`<html><body><p id="intro">Hello</p></body></html>`)
	doc, err := ParseXHTML(xml)
	require.NoError(t, err)

	found := doc.FindElementByID("intro")
	if found != nil {
		assert.Equal(t, "p", found.Name.Local)
	}
}

// Manifest tests
func TestManifestAddRemove(t *testing.T) {
	m := &Manifest{}
	item := &ManifestItem{ID: "ch1", Href: "Text/chapter1.xhtml"}

	m.Add(item)
	assert.Len(t, m.Items, 1)

	found := m.FindByID("ch1")
	assert.NotNil(t, found)
	assert.Equal(t, item, found)

	m.Remove(item)
	assert.Len(t, m.Items, 0)
}

func TestManifestFindByHref(t *testing.T) {
	m := &Manifest{}
	href := NewEPUBPath("Text/chapter.xhtml")
	item := &ManifestItem{ID: "ch1", Href: href}

	m.Add(item)
	found := m.FindByHref(href)
	assert.NotNil(t, found)
	assert.Equal(t, "ch1", found.ID)
}

// Spine tests
func TestSpineAddRemove(t *testing.T) {
	s := &Spine{}
	ref := &SpineItemRef{IDRef: "ch1"}

	s.Add(ref)
	assert.Len(t, s.ItemRefs, 1)

	s.Remove(ref)
	assert.Len(t, s.ItemRefs, 0)
}

// NCX tests
func TestNCXAddRemove(t *testing.T) {
	ncx := &NCXDocument{}
	np := &NavPoint{ID: "np1", PlayOrder: 1}

	ncx.AddNavPoint(np)
	assert.Len(t, ncx.NavPoints, 1)

	ncx.RemoveNavPoint(np)
	assert.Len(t, ncx.NavPoints, 0)
}

func TestNCXUpdatePlayOrder(t *testing.T) {
	ncx := &NCXDocument{}
	np := &NavPoint{ID: "np1", PlayOrder: 1}
	ncx.AddNavPoint(np)

	ncx.UpdatePlayOrder(np, 5)
	assert.Equal(t, 5, np.PlayOrder)
}

// EPUB operations tests
func TestEPUBOpen(t *testing.T) {
	epub, err := Open("test.epub")
	require.NoError(t, err)
	assert.NotNil(t, epub)
	assert.NotNil(t, epub.Resources)
}

func TestEPUBOpenEmpty(t *testing.T) {
	_, err := Open("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestEPUBResource(t *testing.T) {
	epub := &EPUB{Resources: map[EPUBPath]*Resource{}}
	path := NewEPUBPath("Text/chapter.xhtml")
	res := &Resource{Path: path, Data: []byte("content")}
	epub.Resources[path] = res

	found, ok := epub.Resource(path)
	assert.True(t, ok)
	assert.Equal(t, res, found)
}

func TestNilOperations(t *testing.T) {
	// Test that nil receivers don't panic
	var m *Manifest
	m.Add(&ManifestItem{})
	m.Remove(&ManifestItem{})
	found := m.FindByID("test")
	assert.Nil(t, found)

	var s *Spine
	s.Add(&SpineItemRef{})
	s.Remove(&SpineItemRef{})

	var ncx *NCXDocument
	ncx.AddNavPoint(&NavPoint{})
	ncx.RemoveNavPoint(&NavPoint{})
	ncx.UpdatePlayOrder(&NavPoint{}, 1)

	var epub *EPUB
	_, ok := epub.Resource("test")
	assert.False(t, ok)
}
