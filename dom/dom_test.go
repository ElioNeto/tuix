package dom

import (
	"testing"
)

func TestParseSimpleElement(t *testing.T) {
	p := NewParser("<div>hello</div>")
	doc := p.Parse()
	if doc == nil {
		t.Fatal("expected document, got nil")
	}
}

func TestParseElementWithAttributes(t *testing.T) {
	p := NewParser(`<input type="text" id="name" class="input" />`)
	doc := p.Parse()
	if doc == nil {
		t.Fatal("expected document, got nil")
	}
}

func TestParseNestedElements(t *testing.T) {
	p := NewParser(`<div class="container"><span>text</span></div>`)
	doc := p.Parse()
	if doc == nil {
		t.Fatal("expected document, got nil")
	}
}

func TestParseTextContent(t *testing.T) {
	p := NewParser("<p>Hello World</p>")
	doc := p.Parse()
	// Find the text node
	var text string
	var visit func(*Node)
	visit = func(n *Node) {
		if n.Type == NodeText {
			text += n.Data
		}
		for _, c := range n.Children {
			visit(c)
		}
	}
	visit(doc)
	if text != "Hello World" {
		t.Fatalf("expected 'Hello World', got '%s'", text)
	}
}

func TestParseTags(t *testing.T) {
	tests := []struct {
		html string
		tag  string
	}{
		{"<div></div>", "div"},
		{"<br/>", "br"},
		{"<input type='text'/>", "input"},
		{"<a href='#'>link</a>", "a"},
	}
	for _, tt := range tests {
		p := NewParser(tt.html)
		doc := p.Parse()
		if doc == nil {
			t.Fatalf("parse of %q returned nil", tt.html)
		}
	}
}

func TestParseAttributes(t *testing.T) {
	p := NewParser(`<div id="main" class="container active" data-value="test">`)
	doc := p.Parse()
	var div *Node
	var visit func(*Node)
	visit = func(n *Node) {
		if n.Type == NodeElement && n.Data == "div" {
			div = n
		}
		for _, c := range n.Children {
			visit(c)
		}
	}
	visit(doc)
	if div == nil {
		t.Fatal("expected div element, not found")
	}
	if div.GetAttribute("id") != "main" {
		t.Fatalf("expected id='main', got '%s'", div.GetAttribute("id"))
	}
	if div.GetAttribute("class") != "container active" {
		t.Fatalf("expected class='container active', got '%s'", div.GetAttribute("class"))
	}
	if div.GetAttribute("data-value") != "test" {
		t.Fatalf("expected data-value='test', got '%s'", div.GetAttribute("data-value"))
	}
}

func TestParseSelfClosingTag(t *testing.T) {
	p := NewParser("<br/><hr/><img src='test.png'/>")
	doc := p.Parse()
	if doc == nil {
		t.Fatal("expected document, got nil")
	}
}

func TestParseClassAttribute(t *testing.T) {
	p := NewParser(`<span class="foo bar baz">text</span>`)
	doc := p.Parse()
	var span *Node
	var visit func(*Node)
	visit = func(n *Node) {
		if n.Type == NodeElement && n.Data == "span" {
			span = n
		}
		for _, c := range n.Children {
			visit(c)
		}
	}
	visit(doc)
	if span == nil {
		t.Fatal("expected span element, not found")
	}
	if !span.HasClass("foo") {
		t.Fatal("expected span to have class 'foo'")
	}
	if !span.HasClass("bar") {
		t.Fatal("expected span to have class 'bar'")
	}
	if !span.HasClass("baz") {
		t.Fatal("expected span to have class 'baz'")
	}
	if span.HasClass("qux") {
		t.Fatal("expected span NOT to have class 'qux'")
	}
}

func TestQuerySelectorAll(t *testing.T) {
	html := `<div class="container">
		<p>First</p>
		<p>Second</p>
		<span>Third</span>
	</div>`
	p := NewParser(html)
	doc := p.Parse()
	paragraphs := doc.QuerySelectorAll("p")
	if len(paragraphs) != 2 {
		t.Fatalf("expected 2 <p> elements, got %d", len(paragraphs))
	}
}

func TestParseEntityDecoding(t *testing.T) {
	p := NewParser("<div>&amp;&lt;&gt;&quot;</div>")
	doc := p.Parse()
	var text string
	var visit func(*Node)
	visit = func(n *Node) {
		if n.Type == NodeText {
			text += n.Data
		}
		for _, c := range n.Children {
			visit(c)
		}
	}
	visit(doc)
	expected := "&<>\""
	if text != expected {
		t.Fatalf("expected '%s', got '%s'", expected, text)
	}
}
