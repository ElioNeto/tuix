// Package dom provides a Document Object Model for HTML content.
//
// It includes a simple HTML parser that builds a tree of Nodes,
// which can then be styled and rendered by the tuix engine.
package dom

import (
	"strings"
	"unicode"
)

// NodeType distinguishes between different kinds of DOM nodes.
type NodeType int

const (
	NodeElement NodeType = iota
	NodeText
	NodeDocument
	NodeComment
)

// Node represents a single node in the DOM tree.
type Node struct {
	Type       NodeType
	Data       string            // Tag name for elements, text content for text nodes
	Attributes map[string]string // Element attributes
	Children   []*Node
	Parent     *Node
}

// Document creates a new document node.
func Document() *Node {
	return &Node{Type: NodeDocument, Data: "document"}
}

// Element creates a new element node.
func Element(tag string) *Node {
	return &Node{
		Type:       NodeElement,
		Data:       tag,
		Attributes: make(map[string]string),
	}
}

// Text creates a new text node.
func Text(content string) *Node {
	return &Node{
		Type: NodeText,
		Data: content,
	}
}

// AppendChild adds a child node to this node.
func (n *Node) AppendChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

// GetAttribute returns the value of an attribute, or empty string if not found.
func (n *Node) GetAttribute(name string) string {
	if n.Attributes == nil {
		return ""
	}
	return n.Attributes[strings.ToLower(name)]
}

// SetAttribute sets an attribute on the element.
func (n *Node) SetAttribute(name, value string) {
	if n.Attributes == nil {
		n.Attributes = make(map[string]string)
	}
	n.Attributes[strings.ToLower(name)] = value
}

// HasAttribute returns true if the element has the specified attribute.
func (n *Node) HasAttribute(name string) bool {
	_, ok := n.Attributes[strings.ToLower(name)]
	return ok
}

// HasClass checks whether the element has the given CSS class.
func (n *Node) HasClass(class string) bool {
	classes := strings.Fields(n.GetAttribute("class"))
	for _, c := range classes {
		if c == class {
			return true
		}
	}
	return false
}

// ID returns the element's id attribute.
func (n *Node) ID() string {
	return n.GetAttribute("id")
}

// TagName returns the tag name (empty for non-element nodes).
func (n *Node) TagName() string {
	if n.Type != NodeElement {
		return ""
	}
	return n.Data
}

// TextContent returns the concatenated text content of all descendant text nodes.
func (n *Node) TextContent() string {
	var buf strings.Builder
	n.collectText(&buf)
	return buf.String()
}

func (n *Node) collectText(buf *strings.Builder) {
	switch n.Type {
	case NodeText:
		buf.WriteString(n.Data)
	case NodeElement:
		if n.Data != "script" && n.Data != "style" {
			for _, child := range n.Children {
				child.collectText(buf)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// HTML Parser
// ---------------------------------------------------------------------------

// Parser parses HTML source into a DOM tree.
type Parser struct {
	input string
	pos   int
}

// NewParser creates a new HTML parser for the given input string.
func NewParser(input string) *Parser {
	return &Parser{input: input}
}

// Parse parses the HTML input and returns the root document node.
func (p *Parser) Parse() *Node {
	doc := Document()
	p.parseNodes(doc)
	return doc
}

func (p *Parser) parseNodes(parent *Node) {
	for p.pos < len(p.input) {
		p.skipWhitespace()
		if p.pos >= len(p.input) {
			break
		}

		if p.input[p.pos] == '<' {
			if p.pos+1 < len(p.input) {
				if p.input[p.pos+1] == '/' {
					// Closing tag
					break
				} else if p.input[p.pos+1] == '!' {
					// Comment or doctype
					p.parseCommentOrDoctype()
				} else {
					// Opening tag
					elem := p.parseElement()
					if elem != nil {
						parent.AppendChild(elem)
					}
				}
			}
		} else {
			// Text node
			text := p.parseText()
			if text != "" {
				parent.AppendChild(Text(text))
			}
		}
	}
}

func (p *Parser) parseCommentOrDoctype() {
	if p.pos+3 < len(p.input) && p.input[p.pos:p.pos+4] == "<!--" {
		// Comment
		end := strings.Index(p.input[p.pos:], "-->")
		if end == -1 {
			p.pos = len(p.input)
			return
		}
		p.pos += end + 3
	} else if p.pos+8 < len(p.input) && strings.HasPrefix(strings.ToUpper(p.input[p.pos:]), "<!DOCTYPE") {
		// DOCTYPE
		end := strings.Index(p.input[p.pos:], ">")
		if end == -1 {
			p.pos = len(p.input)
			return
		}
		p.pos += end + 1
	} else {
		p.pos++
	}
}

func (p *Parser) parseElement() *Node {
	if p.pos >= len(p.input) || p.input[p.pos] != '<' {
		return nil
	}
	p.pos++ // skip '<'

	// Parse tag name
	tagName := p.parseTagName()
	if tagName == "" {
		return nil
	}

	elem := Element(tagName)

	// Parse attributes
	p.skipWhitespace()
	for p.pos < len(p.input) && p.input[p.pos] != '>' && p.input[p.pos] != '/' {
		attrName, attrValue := p.parseAttribute()
		if attrName != "" {
			elem.SetAttribute(attrName, attrValue)
		}
		p.skipWhitespace()
	}

	// Self-closing tag or void element
	isVoid := p.pos < len(p.input) && p.input[p.pos] == '/'
	if isVoid {
		p.pos++ // skip '/'
	}

	if p.pos < len(p.input) && p.input[p.pos] == '>' {
		p.pos++ // skip '>'
	}

	if isVoid || isVoidElement(tagName) {
		return elem
	}

	// Parse children
	p.parseNodes(elem)

	// Expect closing tag
	p.skipWhitespace()
	if p.pos+2+len(tagName) <= len(p.input) &&
		p.input[p.pos] == '<' && p.input[p.pos+1] == '/' {
		end := strings.Index(p.input[p.pos:], ">")
		if end != -1 {
			p.pos += end + 1
		}
	}

	return elem
}

func (p *Parser) parseTagName() string {
	start := p.pos
	for p.pos < len(p.input) && isTagChar(p.input[p.pos]) {
		p.pos++
	}
	return strings.ToLower(p.input[start:p.pos])
}

func (p *Parser) parseAttribute() (name, value string) {
	// Parse attribute name
	start := p.pos
	for p.pos < len(p.input) && isAttributeChar(p.input[p.pos]) {
		p.pos++
	}
	if p.pos == start {
		return "", ""
	}
	name = strings.ToLower(p.input[start:p.pos])

	p.skipWhitespace()
	if p.pos >= len(p.input) || p.input[p.pos] != '=' {
		return name, ""
	}
	p.pos++ // skip '='

	p.skipWhitespace()
	if p.pos >= len(p.input) {
		return name, ""
	}

	// Parse value
	quote := byte(0)
	if p.input[p.pos] == '"' || p.input[p.pos] == '\'' {
		quote = p.input[p.pos]
		p.pos++ // skip opening quote
	}

	start = p.pos
	if quote != 0 {
		end := strings.IndexByte(p.input[p.pos:], quote)
		if end == -1 {
			p.pos = len(p.input)
			value = p.input[start:]
		} else {
			value = p.input[start : p.pos+end]
			p.pos += end + 1 // skip closing quote
		}
	} else {
		for p.pos < len(p.input) && !unicode.IsSpace(rune(p.input[p.pos])) &&
			p.input[p.pos] != '>' && p.input[p.pos] != '/' {
			p.pos++
		}
		value = p.input[start:p.pos]
	}

	return name, value
}

func (p *Parser) parseText() string {
	start := p.pos
	for p.pos < len(p.input) && p.input[p.pos] != '<' {
		p.pos++
	}
	text := p.input[start:p.pos]
	// Decode HTML entities (basic)
	text = decodeEntities(text)
	return text
}

func (p *Parser) skipWhitespace() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
}

// isTagChar returns true if the byte is valid in a tag name.
func isTagChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') || b == '-' || b == '_' || b == ':'
}

func isAttributeChar(b byte) bool {
	return isTagChar(b)
}

// isVoidElement returns true for HTML elements that cannot have children.
func isVoidElement(tag string) bool {
	switch tag {
	case "area", "base", "br", "col", "embed", "hr", "img", "input",
		"link", "meta", "param", "source", "track", "wbr":
		return true
	}
	return false
}

// decodeEntities decodes common HTML entities in a string.
func decodeEntities(s string) string {
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&apos;", "'")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&#x27;", "'")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	return s
}

// QuerySelectorAll finds all descendant elements matching the given CSS selector.
// This is a simple implementation supporting tag, class, and ID selectors.
func (n *Node) QuerySelectorAll(selector string) []*Node {
	var results []*Node
	n.querySelectorAll(selector, &results)
	return results
}

func (n *Node) querySelectorAll(selector string, results *[]*Node) {
	if n.Type == NodeElement && n.matchesSelector(selector) {
		*results = append(*results, n)
	}
	for _, child := range n.Children {
		child.querySelectorAll(selector, results)
	}
}

// matchesSelector checks if this element matches a simple CSS selector.
// Supported: "tag", ".class", "#id", "tag.class", "tag#id"
func (n *Node) matchesSelector(selector string) bool {
	if n.Type != NodeElement {
		return false
	}

	selector = strings.TrimSpace(selector)
	parts := strings.Split(selector, " ")
	// Only handle the last part (no combinators yet)
	sel := parts[len(parts)-1]

	return matchSimpleSelector(n, sel)
}

func matchSimpleSelector(n *Node, sel string) bool {
	if sel == "*" {
		return true
	}

	// Split by . and #
	var tag, id, class string

	// Check for ID selector (#)
	if idx := strings.IndexByte(sel, '#'); idx >= 0 {
		id = sel[idx+1:]
		sel = sel[:idx]
	}

	// Check for class selector (.)
	if idx := strings.IndexByte(sel, '.'); idx >= 0 {
		class = sel[idx+1:]
		sel = sel[:idx]
	}

	tag = sel

	// Match tag
	if tag != "" && n.Data != tag {
		return false
	}

	// Match ID
	if id != "" && n.ID() != id {
		return false
	}

	// Match class
	if class != "" && !n.HasClass(class) {
		return false
	}

	return true
}
