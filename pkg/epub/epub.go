package epub

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// EPUBPath represents an EPUB resource path using forward slashes.
type EPUBPath string

func NewEPUBPath(value string) EPUBPath {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return EPUBPath("")
	}
	cleaned := filepath.ToSlash(trimmed)
	cleaned = path.Clean(cleaned)
	if cleaned == "." {
		return EPUBPath("")
	}
	if cleaned == "" {
		return EPUBPath("")
	}
	return EPUBPath(cleaned)
}

func (p EPUBPath) String() string {
	return string(p)
}

func (p EPUBPath) Join(other string) EPUBPath {
	base := string(p)
	if base == "" {
		return NewEPUBPath(other)
	}
	return NewEPUBPath(path.Join(base, other))
}

func (p EPUBPath) RelativeTo(other EPUBPath) EPUBPath {
	base := string(other)
	if base == "" {
		base = "."
	}
	rel, err := path.Rel(base, string(p))
	if err != nil {
		return NewEPUBPath(string(p))
	}
	return NewEPUBPath(rel)
}

// SourceRange is a byte range into the original document.
type SourceRange struct {
	Start int
	End   int
}

func (r SourceRange) Len() int {
	if r.End < r.Start {
		return 0
	}
	return r.End - r.Start
}

func (r SourceRange) Valid(length int) bool {
	if length < 0 {
		return false
	}
	if r.Start < 0 || r.End < 0 {
		return false
	}
	if r.End < r.Start {
		return false
	}
	return r.End <= length
}

// Edit describes a byte-based edit to apply to a document.
type Edit struct {
	Start       int
	End         int
	Replacement []byte
}

func ApplyEdits(data []byte, edits []Edit) ([]byte, error) {
	if len(edits) == 0 {
		return append([]byte(nil), data...), nil
	}
	ordered := make([]Edit, len(edits))
	copy(ordered, edits)
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Start != ordered[j].Start {
			return ordered[i].Start < ordered[j].Start
		}
		return ordered[i].End < ordered[j].End
	})

	for i, edit := range ordered {
		if edit.Start < 0 || edit.End < 0 {
			return nil, fmt.Errorf("negative offset in edit %d: [%d,%d]", i, edit.Start, edit.End)
		}
		if edit.End < edit.Start {
			return nil, fmt.Errorf("invalid edit %d: end before start [%d,%d]", i, edit.Start, edit.End)
		}
		if edit.End > len(data) {
			return nil, fmt.Errorf("edit %d exceeds document length: [%d,%d] > %d", i, edit.Start, edit.End, len(data))
		}
		if i > 0 && ordered[i-1].End > edit.Start {
			return nil, fmt.Errorf("overlapping edits: [%d,%d] and [%d,%d]", ordered[i-1].Start, ordered[i-1].End, edit.Start, edit.End)
		}
	}

	result := make([]byte, 0, len(data)+sumReplacementSize(ordered))
	cursor := 0
	for _, edit := range ordered {
		if edit.Start > cursor {
			result = append(result, data[cursor:edit.Start]...)
		}
		result = append(result, edit.Replacement...)
		cursor = edit.End
	}
	if cursor < len(data) {
		result = append(result, data[cursor:]...)
	}
	return result, nil
}

func sumReplacementSize(edits []Edit) int {
	sum := 0
	for _, edit := range edits {
		sum += len(edit.Replacement)
	}
	return sum
}

// XMLNode is a lightweight, source-aware XML node.
type XMLNode interface {
	Range() SourceRange
}

type XMLAttribute struct {
	Name       xml.Name
	Value      string
	Range      SourceRange
	ValueRange SourceRange
}

type XMLElement struct {
	Name         xml.Name
	Range        SourceRange
	ContentRange SourceRange
	Attributes   []XMLAttribute
	Children     []XMLNode
	Parent       *XMLElement
}

func (e *XMLElement) Range() SourceRange {
	if e == nil {
		return SourceRange{}
	}
	return e.Range
}

func (e *XMLElement) Raw(data []byte) []byte {
	if e == nil || !e.Range.Valid(len(data)) {
		return nil
	}
	return data[e.Range.Start:e.Range.End]
}

func (e *XMLElement) Content(data []byte) []byte {
	if e == nil || !e.ContentRange.Valid(len(data)) {
		return nil
	}
	return data[e.ContentRange.Start:e.ContentRange.End]
}

func (e *XMLElement) FindElementByID(id string) *XMLElement {
	if e == nil {
		return nil
	}
	for _, attr := range e.Attributes {
		if attr.Name.Local == "id" && attr.Value == id {
			return e
		}
	}
	for _, child := range e.Children {
		if childEl, ok := child.(*XMLElement); ok {
			if found := childEl.FindElementByID(id); found != nil {
				return found
			}
		}
	}
	return nil
}

func (e *XMLElement) FindElements(name string) []*XMLElement {
	if e == nil {
		return nil
	}
	var out []*XMLElement
	if e.Name.Local == name {
		out = append(out, e)
	}
	for _, child := range e.Children {
		if childEl, ok := child.(*XMLElement); ok {
			out = append(out, childEl.FindElements(name)...)
		}
	}
	return out
}

type XMLText struct {
	Range SourceRange
}

func (t *XMLText) Range() SourceRange {
	if t == nil {
		return SourceRange{}
	}
	return t.Range
}

func (t *XMLText) Raw(data []byte) []byte {
	if t == nil || !t.Range.Valid(len(data)) {
		return nil
	}
	return data[t.Range.Start:t.Range.End]
}

type XMLComment struct {
	Range SourceRange
}

func (c *XMLComment) Range() SourceRange {
	if c == nil {
		return SourceRange{}
	}
	return c.Range
}

func (c *XMLComment) Raw(data []byte) []byte {
	if c == nil || !c.Range.Valid(len(data)) {
		return nil
	}
	return data[c.Range.Start:c.Range.End]
}

// XMLDocument is a byte-preserving XML document with a source-aware tree.
type XMLDocument struct {
	Data []byte
	Root *XMLElement
}

func ParseXML(data []byte) (*XMLDocument, error) {
	doc := &XMLDocument{Data: append([]byte(nil), data...)}
	if len(data) == 0 {
		return doc, nil
	}

	cursor := 0
	stack := []*XMLElement{}

	for cursor < len(data) {
		nextTag := bytes.Index(data[cursor:], []byte("<"))
		if nextTag < 0 {
			break
		}
		absoluteStart := cursor + nextTag
		if nextTag > 0 {
			textRange := SourceRange{Start: cursor, End: absoluteStart}
			if len(stack) > 0 {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, &XMLText{Range: textRange})
			}
			cursor = absoluteStart
		}
		if cursor >= len(data) {
			break
		}

		switch {
		case bytes.HasPrefix(data[cursor:], []byte("<!--")):
			end := bytes.Index(data[cursor+4:], []byte("-->"))
			if end < 0 {
				break
			}
			commentEnd := cursor + 4 + end + 3
			if len(stack) > 0 {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, &XMLComment{Range: SourceRange{Start: cursor, End: commentEnd}})
			}
			cursor = commentEnd
		case bytes.HasPrefix(data[cursor:], []byte("<?")) || bytes.HasPrefix(data[cursor:], []byte("<!")):
			end := findTagEnd(data, cursor)
			if end < 0 {
				break
			}
			cursor = end
		case bytes.HasPrefix(data[cursor:], []byte("</")):
			end := findTagEnd(data, cursor)
			if end < 0 {
				break
			}
			if len(stack) > 0 {
				current := stack[len(stack)-1]
				current.ContentRange = SourceRange{Start: current.Range.End, End: cursor}
				if current.ContentRange.Start > current.ContentRange.End {
					current.ContentRange = SourceRange{Start: current.Range.End, End: current.Range.End}
				}
				stack = stack[:len(stack)-1]
			}
			cursor = end
		default:
			end := findTagEnd(data, cursor)
			if end < 0 {
				break
			}
			raw := string(data[cursor:end])
			dec := xml.NewDecoder(strings.NewReader(raw))
			tok, err := dec.Token()
			if err != nil {
				cursor = end
				continue
			}
			se, ok := tok.(xml.StartElement)
			if !ok {
				cursor = end
				continue
			}

			attrs := make([]XMLAttribute, 0, len(se.Attr))
			for _, attr := range se.Attr {
				attrRange, valueRange := findAttributeRanges(raw, attr.Name.Local, attr.Value)
				attrs = append(attrs, XMLAttribute{
					Name:       attr.Name,
					Value:      attr.Value,
					Range:      SourceRange{Start: cursor + attrRange.Start, End: cursor + attrRange.End},
					ValueRange: SourceRange{Start: cursor + valueRange.Start, End: cursor + valueRange.End},
				})
			}

			el := &XMLElement{
				Name:       se.Name,
				Range:      SourceRange{Start: cursor, End: end},
				Attributes: attrs,
				Parent:     parentFromStack(stack),
			}
			if el.Parent != nil {
				el.Parent.Children = append(el.Parent.Children, el)
			} else {
				doc.Root = el
			}
			stack = append(stack, el)
			if strings.HasSuffix(raw, "/>") {
				stack = stack[:len(stack)-1]
			}
			cursor = end
		}
	}
	return doc, nil
}

func parentFromStack(stack []*XMLElement) *XMLElement {
	if len(stack) == 0 {
		return nil
	}
	return stack[len(stack)-1]
}

func findTagEnd(data []byte, start int) int {
	if start < 0 || start >= len(data) {
		return -1
	}
	inQuote := byte(0)
	for i := start; i < len(data); i++ {
		b := data[i]
		if inQuote != 0 {
			if b == inQuote {
				inQuote = 0
			}
			continue
		}
		if b == '\'' || b == '"' {
			inQuote = b
			continue
		}
		if b == '>' {
			return i + 1
		}
	}
	return -1
}

func findAttributeRanges(raw, name, value string) (SourceRange, SourceRange) {
	if name == "" {
		return SourceRange{}, SourceRange{}
	}
	start := strings.Index(raw, name)
	if start < 0 {
		return SourceRange{}, SourceRange{}
	}
	valueIdx := strings.Index(raw[start:], "=")
	if valueIdx < 0 {
		return SourceRange{Start: start, End: len(raw)}, SourceRange{Start: start, End: len(raw)}
	}
	valuePos := start + valueIdx + 1
	for valuePos < len(raw) && (raw[valuePos] == ' ' || raw[valuePos] == '\t' || raw[valuePos] == '\n' || raw[valuePos] == '\r') {
		valuePos++
	}
	if valuePos < len(raw) && (raw[valuePos] == '\'' || raw[valuePos] == '"') {
		quote := raw[valuePos]
		valueStart := valuePos + 1
		valueEnd := valueStart
		for valueEnd < len(raw) && raw[valueEnd] != quote {
			valueEnd++
		}
		endOfAttr := valueEnd + 1
		if endOfAttr > len(raw) {
			endOfAttr = len(raw)
		}
		return SourceRange{Start: start, End: endOfAttr}, SourceRange{Start: valueStart, End: valueEnd}
	}
	valueEnd := valuePos
	for valueEnd < len(raw) && raw[valueEnd] != '>' && raw[valueEnd] != ' ' && raw[valueEnd] != '\t' && raw[valueEnd] != '\n' && raw[valueEnd] != '\r' {
		valueEnd++
	}
	return SourceRange{Start: start, End: valueEnd}, SourceRange{Start: valuePos, End: valueEnd}
}

// XHTMLDocument is a source-aware XHTML document.
type XHTMLDocument struct {
	*XMLDocument
}

func ParseXHTML(data []byte) (*XHTMLDocument, error) {
	doc, err := ParseXML(data)
	if err != nil {
		return nil, err
	}
	return &XHTMLDocument{XMLDocument: doc}, nil
}

func (d *XHTMLDocument) FindElements(name string) []*XMLElement {
	if d == nil || d.Root == nil {
		return nil
	}
	return d.Root.FindElements(name)
}

func (d *XHTMLDocument) FindElementByID(id string) *XMLElement {
	if d == nil || d.Root == nil {
		return nil
	}
	return d.Root.FindElementByID(id)
}

func (d *XHTMLDocument) FindText(query string) []*XMLText {
	if d == nil || d.Root == nil {
		return nil
	}
	return findTextNodes(d.Root, query)
}

func findTextNodes(node *XMLElement, query string) []*XMLText {
	if node == nil {
		return nil
	}
	var matches []*XMLText
	for _, child := range node.Children {
		if text, ok := child.(*XMLText); ok {
			if text != nil && text.Range.Valid(len(node.Raw(nil))) {
				if strings.Contains(string(text.Raw(nil)), query) {
					matches = append(matches, text)
				}
			}
		}
		if childEl, ok := child.(*XMLElement); ok {
			matches = append(matches, findTextNodes(childEl, query)...)
		}
	}
	return matches
}

// Metadata / OPF placeholder types.
type Metadata struct {
	Identifier string
}

type ManifestItem struct {
	ID        string
	Href      EPUBPath
	MediaType string
	Properties []string
	Element   *XMLElement
}

type Manifest struct {
	Items []*ManifestItem
}

func (m *Manifest) Add(item *ManifestItem) {
	if m == nil {
		return
	}
	m.Items = append(m.Items, item)
}

func (m *Manifest) Remove(item *ManifestItem) {
	if m == nil || item == nil {
		return
	}
	filtered := make([]*ManifestItem, 0, len(m.Items))
	for _, existing := range m.Items {
		if existing != item {
			filtered = append(filtered, existing)
		}
	}
	m.Items = filtered
}

func (m *Manifest) FindByID(id string) *ManifestItem {
	if m == nil {
		return nil
	}
	for _, item := range m.Items {
		if item != nil && item.ID == id {
			return item
		}
	}
	return nil
}

func (m *Manifest) FindByHref(href EPUBPath) *ManifestItem {
	if m == nil {
		return nil
	}
	for _, item := range m.Items {
		if item != nil && item.Href == href {
			return item
		}
	}
	return nil
}

type SpineItemRef struct {
	IDRef   string
	Element *XMLElement
}

type Spine struct {
	ItemRefs []*SpineItemRef
}

func (s *Spine) Add(ref *SpineItemRef) {
	if s == nil {
		return
	}
	s.ItemRefs = append(s.ItemRefs, ref)
}

func (s *Spine) Remove(ref *SpineItemRef) {
	if s == nil || ref == nil {
		return
	}
	filtered := make([]*SpineItemRef, 0, len(s.ItemRefs))
	for _, item := range s.ItemRefs {
		if item != ref {
			filtered = append(filtered, item)
		}
	}
	s.ItemRefs = filtered
}

type PackageDocument struct {
	*XMLDocument
	Metadata *Metadata
	Manifest *Manifest
	Spine    *Spine
}

type NavigationItem struct {
	Label    string
	Href     EPUBPath
	Children []*NavigationItem
	Element  *XMLElement
}

type NavigationList struct {
	Items []*NavigationItem
}

type NavigationDocument struct {
	*XMLDocument
	TOC       *NavigationList
	Landmarks *NavigationList
	PageList  *NavigationList
}

type NCXDocument struct {
	*XMLDocument
	NavPoints []*NavPoint
}

type NavPoint struct {
	ID        string
	PlayOrder int
	Label     string
	Src       EPUBPath
	Element   *XMLElement
}

func (n *NCXDocument) AddNavPoint(np *NavPoint) {
	if n == nil || np == nil {
		return
	}
	n.NavPoints = append(n.NavPoints, np)
}

func (n *NCXDocument) RemoveNavPoint(np *NavPoint) {
	if n == nil || np == nil {
		return
	}
	filtered := make([]*NavPoint, 0, len(n.NavPoints))
	for _, item := range n.NavPoints {
		if item != np {
			filtered = append(filtered, item)
		}
	}
	n.NavPoints = filtered
}

func (n *NCXDocument) UpdatePlayOrder(np *NavPoint, playOrder int) {
	if n == nil || np == nil {
		return
	}
	for _, item := range n.NavPoints {
		if item == np {
			item.PlayOrder = playOrder
			return
		}
	}
}

type Resource struct {
	Path      EPUBPath
	MediaType string
	Properties []string
	Data      []byte
}

type ResourceRef struct {
	Resource *Resource
	Fragment string
}

type Container struct{}

type EPUB struct {
	Resources map[EPUBPath]*Resource
	Container *Container
	Package   *PackageDocument
	Nav       *NavigationDocument
	NCX       *NCXDocument
}

func Open(pathname string) (*EPUB, error) {
	if strings.TrimSpace(pathname) == "" {
		return nil, fmt.Errorf("epub path must not be empty")
	}
	return &EPUB{Resources: map[EPUBPath]*Resource{}}, nil
}

func (e *EPUB) Resource(pathname EPUBPath) (*Resource, bool) {
	if e == nil {
		return nil, false
	}
	res, ok := e.Resources[pathname]
	return res, ok
}

func (e *EPUB) Save(pathname string) error {
	if strings.TrimSpace(pathname) == "" {
		return fmt.Errorf("epub save path must not be empty")
	}
	return nil
}

func (e *EPUB) ResolveHref(from *Resource, href EPUBPath) (*ResourceRef, error) {
	if e == nil {
		return nil, fmt.Errorf("epub has not been initialized")
	}
	resolved := NewEPUBPath(string(href))
	if from != nil && from.Path != "" {
		resolved = EPUBPath(path.Clean(path.Join(string(from.Path), string(resolved))))
	}
	if res, ok := e.Resources[resolved]; ok {
		return &ResourceRef{Resource: res, Fragment: ""}, nil
	}
	return nil, fmt.Errorf("resource %q not found", resolved)
}

func normalizeEPUBPathForComparison(value string) string {
	if value == "" {
		return ""
	}
	return path.Clean(filepath.ToSlash(value))
}
