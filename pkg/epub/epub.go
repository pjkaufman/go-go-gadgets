package epub

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// EPUBPath represents an EPUB resource path, using forward slashes as the canonical format.
type EPUBPath string

func NewEPUBPath(value string) EPUBPath {
	if value == "" {
		return EPUBPath("")
	}
	cleaned := filepath.ToSlash(value)
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		return EPUBPath("")
	}
	cleaned = path.Clean(cleaned)
	if cleaned == "." {
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
	joined := path.Join(base, other)
	return NewEPUBPath(joined)
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

// Edit describes a byte-based replacement to apply to a document.
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
		result = append(result, data[cursor:edit.Start]...)
		result = append(result, edit.Replacement...)
		cursor = edit.End
	}
	result = append(result, data[cursor:]...)
	return result, nil
}

func sumReplacementSize(edits []Edit) int {
	sum := 0
	for _, e := range edits {
		sum += len(e.Replacement)
	}
	return sum
}

// XMLNode is a lightweight source-aware XML node.
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

func (e *XMLElement) Range() SourceRange { if e == nil { return SourceRange{} }; return e.Range }

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
	for _, a := range e.Attributes {
		if a.Name.Local == "id" && a.Value == id {
			return e
		}
	}
	for _, c := range e.Children {
		if childEl, ok := c.(*XMLElement); ok {
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
	var matches []*XMLElement
	if e.Name.Local == name {
		matches = append(matches, e)
	}
	for _, child := range e.Children {
		if childEl, ok := child.(*XMLElement); ok {
			matches = append(matches, childEl.FindElements(name)...)
		}
	}
	return matches
}

type XMLText struct {
	Range SourceRange
}

func (t *XMLText) Range() SourceRange { if t == nil { return SourceRange{} }; return t.Range }

func (t *XMLText) Raw(data []byte) []byte {
	if t == nil || !t.Range.Valid(len(data)) {
		return nil
	}
	return data[t.Range.Start:t.Range.End]
}

type XMLComment struct {
	Range SourceRange
}

func (c *XMLComment) Range() SourceRange { if c == nil { return SourceRange{} }; return c.Range }

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
	if data == nil {
		return &XMLDocument{Data: nil}, nil
	}

	doc := &XMLDocument{Data: append([]byte(nil), data...)}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var stack []*XMLElement
	searchPos := 0
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			start, end := findTagRange(data, searchPos)
			if start == -1 || end == -1 {
				start = searchPos
				end = len(data)
			}
			sel := &XMLElement{
				Name:   t.Name,
				Range:  SourceRange{Start: start, End: end},
				Parent: parentFromStack(stack),
			}
			sel.Attributes = attributesFromToken(data, t.Attr, start, end)
			if sel.Parent != nil {
				sel.Parent.Children = append(sel.Parent.Children, sel)
			} else {
				doc.Root = sel
			}
			stack = append(stack, sel)
			searchPos = end
		case xml.EndElement:
			start, end := findClosingTagRange(data, searchPos)
			if start == -1 || end == -1 {
				start = searchPos
				end = searchPos
			}
			if len(stack) > 0 {
				current := stack[len(stack)-1]
				current.ContentRange = SourceRange{Start: current.Range.End, End: start}
				if current.ContentRange.Start > current.ContentRange.End {
					current.ContentRange = SourceRange{Start: current.Range.End, End: current.Range.End}
				}
				stack = stack[:len(stack)-1]
			}
			searchPos = end
		case xml.CharData:
			start, end := nextTextRange(data, searchPos)
			if start < end && len(stack) > 0 {
				value := &XMLText{Range: SourceRange{Start: start, End: end}}
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, value)
				searchPos = end
			} else {
				searchPos = nextTagPosition(data, searchPos)
			}
		case xml.Comment:
			start, end := findCommentRange(data, searchPos)
			if start != -1 && end != -1 {
				comment := &XMLComment{Range: SourceRange{Start: start, End: end}}
				if len(stack) > 0 {
					stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, comment)
				}
				searchPos = end
			}
		case xml.ProcInst, xml.Directive:
			start, end := findTagRange(data, searchPos)
			if start != -1 && end != -1 {
				searchPos = end
			}
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

func findTagRange(data []byte, from int) (int, int) {
	start := bytes.Index(data[from:], []byte("<"))
	if start == -1 {
		return -1, -1
	}
	start += from
	end := start
	inQuote := byte(0)
	for end < len(data) {
		b := data[end]
		if inQuote != 0 {
			if b == inQuote {
				inQuote = 0
			}
			end++
			continue
		}
		if b == '\'' || b == '"' {
			inQuote = b
			end++
			continue
		}
		if b == '>' {
			return start, end + 1
		}
		end++
	}
	return -1, -1
}

func findClosingTagRange(data []byte, from int) (int, int) {
	start := bytes.Index(data[from:], []byte("</"))
	if start == -1 {
		return -1, -1
	}
	start += from
	end := start + 2
	for end < len(data) {
		if data[end] == '>' {
			return start, end + 1
		}
		end++
	}
	return -1, -1
}

func nextTagPosition(data []byte, from int) int {
	idx := bytes.Index(data[from:], []byte("<"))
	if idx == -1 {
		return len(data)
	}
	return from + idx
}

func nextTextRange(data []byte, from int) (int, int) {
	start := from
	end := nextTagPosition(data, from)
	if end < start {
		return start, start
	}
	if end > len(data) {
		end = len(data)
	}
	if start == end {
		return start, start
	}
	return start, end
}

func findCommentRange(data []byte, from int) (int, int) {
	start := bytes.Index(data[from:], []byte("<!--"))
	if start == -1 {
		return -1, -1
	}
	start += from
	end := bytes.Index(data[start+4:], []byte("-->"))
	if end == -1 {
		return -1, -1
	}
	end += start + 4
	return start, end + 3
}

func attributesFromToken(data []byte, attrs []xml.Attr, tagStart, tagEnd int) []XMLAttribute {
	if len(attrs) == 0 {
		return nil
	}
	result := make([]XMLAttribute, 0, len(attrs))
	raw := string(data[tagStart:tagEnd])
	for _, attr := range attrs {
		attrRange, valueRange := findAttributeRange(raw, attr.Name.Local, attr.Value)
		result = append(result, XMLAttribute{
			Name:       attr.Name,
			Value:      attr.Value,
			Range:      SourceRange{Start: tagStart + attrRange.Start, End: tagStart + attrRange.End},
			ValueRange: SourceRange{Start: tagStart + valueRange.Start, End: tagStart + valueRange.End},
		})
	}
	return result
}

func findAttributeRange(raw, name, value string) (SourceRange, SourceRange) {
	attrName := name
	if attrName == "" {
		return SourceRange{}, SourceRange{}
	}
	namePattern := attrName + "="
	start := strings.Index(raw, namePattern)
	if start == -1 {
		namePattern = attrName + " ="
		start = strings.Index(raw, namePattern)
	}
	if start == -1 {
		return SourceRange{}, SourceRange{}
	}
	valueStart := start + len(namePattern)
	if valueStart >= len(raw) {
		return SourceRange{Start: start, End: len(raw)}, SourceRange{Start: start, End: len(raw)}
	}
	v := raw[valueStart:]
	quote := byte(0)
	if len(v) > 0 && (v[0] == '\'' || v[0] == '"') {
		quote = v[0]
		valueStart++
		valueEnd := valueStart
		for valueEnd < len(raw) && raw[valueEnd] != quote {
			valueEnd++
		}
		return SourceRange{Start: start, End: valueEnd + 1}, SourceRange{Start: valueStart, End: valueEnd}
	}
	valueEnd := valueStart
	for valueEnd < len(raw) && !strings.ContainsRune(" \t\r\n/>", rune(raw[valueEnd])) {
		valueEnd++
	}
	return SourceRange{Start: start, End: valueEnd}, SourceRange{Start: valueStart, End: valueEnd}
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
			if strings.Contains(string(text.Raw(node.Raw(nil))), query) {
				matches = append(matches, text)
			}
		}
		if childEl, ok := child.(*XMLElement); ok {
			matches = append(matches, findTextNodes(childEl, query)...)
		}
	}
	return matches
}

// Metadata / OPF placeholder data model.

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

func (m *Manifest) FindByID(id string) *ManifestItem {
	for _, item := range m.Items {
		if item != nil && item.ID == id {
			return item
		}
	}
	return nil
}

func (m *Manifest) FindByHref(href EPUBPath) *ManifestItem {
	for _, item := range m.Items {
		if item != nil && item.Href == href {
			return item
		}
	}
	return nil
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

// Navigation list structures.
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
	if e == nil {
		return nil
	}
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
		frag := ""
		if idx := strings.Index(string(resolved), "#"); idx >= 0 {
			frag = string(resolved)[idx+1:]
		}
		return &ResourceRef{Resource: res, Fragment: frag}, nil
	}
	return nil, fmt.Errorf("resource %q not found", resolved)
}

func (e *EPUB) ResolveHrefString(from *Resource, href string) (*ResourceRef, error) {
	return e.ResolveHref(from, NewEPUBPath(href))
}

// Small helper for path/string comparisons.
func normalizeEPUBPathForComparison(value string) string {
	if value == "" {
		return ""
	}
	cleaned := filepath.ToSlash(value)
	return path.Clean(cleaned)
}

func init() {
	// no-op
}

// Ensure file names remain usable in tests with the package name consistent.
func checksum(data []byte) uint32 {
	var sum uint32
	for _, b := range data {
		sum += uint32(b)
	}
	return sum
}

func findTextNodeByValue(node *XMLElement, value string) *XMLText {
	if node == nil {
		return nil
	}
	for _, child := range node.Children {
		if text, ok := child.(*XMLText); ok {
			if string(text.Raw(node.Raw(nil))) == value {
				return text
			}
		}
		if childEl, ok := child.(*XMLElement); ok {
			if found := findTextNodeByValue(childEl, value); found != nil {
				return found
			}
		}
	}
	return nil
}

func findElementByName(node *XMLElement, name string) *XMLElement {
	if node == nil {
		return nil
	}
	if node.Name.Local == name {
		return node
	}
	for _, child := range node.Children {
		if childEl, ok := child.(*XMLElement); ok {
			if found := findElementByName(childEl, name); found != nil {
				return found
			}
		}
	}
	return nil
}

func findElementByIDInTree(node *XMLElement, id string) *XMLElement {
	if node == nil {
		return nil
	}
	for _, attr := range node.Attributes {
		if attr.Name.Local == "id" && attr.Value == id {
			return node
		}
	}
	for _, child := range node.Children {
		if childEl, ok := child.(*XMLElement); ok {
			if found := findElementByIDInTree(childEl, id); found != nil {
				return found
			}
		}
	}
	return nil
}

func findTextNodesUnder(node *XMLElement, predicate func(*XMLText) bool) []*XMLText {
	if node == nil {
		return nil
	}
	var matches []*XMLText
	for _, child := range node.Children {
		if text, ok := child.(*XMLText); ok && predicate(text) {
			matches = append(matches, text)
		}
		if childEl, ok := child.(*XMLElement); ok {
			matches = append(matches, findTextNodesUnder(childEl, predicate)...)
		}
	}
	return matches
}
