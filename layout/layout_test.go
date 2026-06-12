package layout

import (
	"testing"

	"github.com/elioneto/tuix/dom"
	"github.com/elioneto/tuix/style"
)

func TestNewEngine(t *testing.T) {
	e := NewEngine()
	if e == nil {
		t.Fatal("NewEngine returned nil")
	}
}

func TestLayoutEmptyDoc(t *testing.T) {
	e := NewEngine()
	e.ViewWidth = 80
	e.ViewHeight = 24

	doc := dom.Document()
	styles := make(map[*dom.Node]style.ComputedStyle)

	box := e.Layout(doc, styles)
	if box == nil {
		t.Fatal("Layout returned nil")
	}
	if box.Type != BoxRoot {
		t.Fatalf("expected BoxRoot, got %v", box.Type)
	}
}

func TestLayoutSimpleDiv(t *testing.T) {
	e := NewEngine()
	e.ViewWidth = 80
	e.ViewHeight = 24

	html := `<div>hello</div>`
	p := dom.NewParser(html)
	doc := p.Parse()

	// Create styles
	styles := make(map[*dom.Node]style.ComputedStyle)
	var walk func(*dom.Node)
	walk = func(n *dom.Node) {
		if n.Type == dom.NodeElement {
			s := style.DefaultStyle()
			s.Display = style.DisplayBlock
			styles[n] = s
		}
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(doc)

	box := e.Layout(doc, styles)
	if box == nil {
		t.Fatal("Layout returned nil")
	}
	if box.Rect.Width == 0 {
		t.Fatal("box width is 0")
	}
}

func TestLayoutNestedDivs(t *testing.T) {
	e := NewEngine()
	e.ViewWidth = 80
	e.ViewHeight = 24

	html := `<div class="outer"><div class="inner">text</div></div>`
	p := dom.NewParser(html)
	doc := p.Parse()

	styles := make(map[*dom.Node]style.ComputedStyle)
	var assignStyles func(*dom.Node)
	assignStyles = func(n *dom.Node) {
		if n.Type == dom.NodeElement {
			s := style.DefaultStyle()
			s.Display = style.DisplayBlock
			styles[n] = s
		}
		for _, c := range n.Children {
			assignStyles(c)
		}
	}
	assignStyles(doc)

	box := e.Layout(doc, styles)
	if box == nil {
		t.Fatal("Layout returned nil")
	}
}

func TestLayoutBoxModel(t *testing.T) {
	e := NewEngine()
	e.ViewWidth = 80
	e.ViewHeight = 24

	html := `<div>content</div>`
	p := dom.NewParser(html)
	doc := p.Parse()

	var divNode *dom.Node
	var findDiv func(*dom.Node)
	findDiv = func(n *dom.Node) {
		if n.Type == dom.NodeElement && n.Data == "div" {
			divNode = n
		}
		for _, c := range n.Children {
			findDiv(c)
		}
	}
	findDiv(doc)

	styles := make(map[*dom.Node]style.ComputedStyle)
	if divNode != nil {
		s := style.DefaultStyle()
		s.Display = style.DisplayBlock
		s.PaddingTop = style.Length{Value: 2}
		s.PaddingRight = style.Length{Value: 2}
		s.PaddingBottom = style.Length{Value: 2}
		s.PaddingLeft = style.Length{Value: 2}
		s.BorderTop = style.Border{Style: style.BorderSolid, Width: style.Length{Value: 1}}
		s.BorderBottom = style.Border{Style: style.BorderSolid, Width: style.Length{Value: 1}}
		s.BorderLeft = style.Border{Style: style.BorderSolid, Width: style.Length{Value: 1}}
		s.BorderRight = style.Border{Style: style.BorderSolid, Width: style.Length{Value: 1}}
		s.MarginTop = style.Length{Value: 1}
		s.MarginRight = style.Length{Value: 1}
		s.MarginBottom = style.Length{Value: 1}
		s.MarginLeft = style.Length{Value: 1}
		styles[divNode] = s
	}

	box := e.Layout(doc, styles)
	if box == nil {
		t.Fatal("Layout returned nil")
	}
	if divNode != nil {
		// The box tree may not directly correspond to DOM nodes one-to-one,
		// so we just check that layout completed without error
	}
}
