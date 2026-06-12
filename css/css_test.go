package css

import (
	"testing"
)

func TestParseSimpleSelector(t *testing.T) {
	p := NewParser("div { color: red; }")
	sheet, _ := p.Parse()
	if len(sheet.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(sheet.Rules))
	}
	rule := sheet.Rules[0]
	if len(rule.Selectors) != 1 {
		t.Fatalf("expected 1 selector, got %d", len(rule.Selectors))
	}
	if rule.Selectors[0].Type != SelectorTag {
		t.Fatalf("expected SelectorTag, got %v", rule.Selectors[0].Type)
	}
	if rule.Selectors[0].Value != "div" {
		t.Fatalf("expected 'div', got '%s'", rule.Selectors[0].Value)
	}
	if len(rule.Declarations) != 1 {
		t.Fatalf("expected 1 declaration, got %d", len(rule.Declarations))
	}
	if rule.Declarations[0].Property != "color" {
		t.Fatalf("expected 'color', got '%s'", rule.Declarations[0].Property)
	}
	// 'red' is parsed as ValueColor because it's a named CSS color
	val := rule.Declarations[0].Value
	if val.Type != ValueColor {
		t.Fatalf("expected ValueColor for 'red', got %v", val.Type)
	}
}

func TestParseClassSelector(t *testing.T) {
	p := NewParser(".btn { padding: 1; }")
	sheet, _ := p.Parse()
	if len(sheet.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(sheet.Rules))
	}
	sel := sheet.Rules[0].Selectors[0]
	if sel.Type != SelectorClass {
		t.Fatalf("expected SelectorClass, got %v", sel.Type)
	}
	if sel.Value != "btn" {
		t.Fatalf("expected 'btn', got '%s'", sel.Value)
	}
}

func TestParseIDSelector(t *testing.T) {
	p := NewParser("#app { width: 100%; }")
	sheet, _ := p.Parse()
	sel := sheet.Rules[0].Selectors[0]
	if sel.Type != SelectorID {
		t.Fatalf("expected SelectorID, got %v", sel.Type)
	}
	if sel.Value != "app" {
		t.Fatalf("expected 'app', got '%s'", sel.Value)
	}
}

func TestParseCompoundSelector(t *testing.T) {
	p := NewParser("div.btn.primary { color: blue; }")
	sheet, _ := p.Parse()
	sel := sheet.Rules[0].Selectors[0]
	if sel.Type != SelectorTag {
		t.Fatalf("expected SelectorTag, got %v", sel.Type)
	}
	if sel.Value != "div" {
		t.Fatalf("expected 'div', got '%s'", sel.Value)
	}
	if len(sel.Classes) != 2 {
		t.Fatalf("expected 2 classes, got %d: %v", len(sel.Classes), sel.Classes)
	}
}

func TestParseMultipleRules(t *testing.T) {
	css := `
		h1 { color: red; }
		h2 { color: blue; font-size: 16; }
	`
	p := NewParser(css)
	sheet, _ := p.Parse()
	if len(sheet.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(sheet.Rules))
	}
}

func TestSpecificity(t *testing.T) {
	tests := []struct {
		selector string
		b        int // IDs
		c        int // classes
		d        int // elements
	}{
		{"div", 0, 0, 1},
		{".btn", 0, 1, 0},
		{"#app", 1, 0, 0},
		{"div.btn", 0, 1, 1},
		{"#app .btn", 1, 1, 0},
	}

	for _, tt := range tests {
		p := NewParser(tt.selector + " { color: red; }")
		sheet, errs := p.Parse()
		if len(errs) > 0 {
			t.Fatalf("selector %q: unexpected errors: %v", tt.selector, errs)
		}
		spec := sheet.Rules[0].Selectors[0].Specificity
		if spec.B != tt.b || spec.C != tt.c || spec.D != tt.d {
			t.Errorf("selector %q: expected (%d,%d,%d), got (%d,%d,%d)",
				tt.selector, tt.b, tt.c, tt.d, spec.B, spec.C, spec.D)
		}
	}
}

func TestParseUniversalSelector(t *testing.T) {
	p := NewParser("* { margin: 0; }")
	sheet, _ := p.Parse()
	if sheet.Rules[0].Selectors[0].Type != SelectorUniversal {
		t.Fatalf("expected SelectorUniversal, got %v", sheet.Rules[0].Selectors[0].Type)
	}
}

func TestParsePseudoClassSelector(t *testing.T) {
	// Standalone pseudo-class selector
	p := NewParser(":focus { border-color: cyan; }")
	sheet, _ := p.Parse()
	sel := sheet.Rules[0].Selectors[0]
	if sel.Type != SelectorPseudoClass {
		t.Fatalf("expected SelectorPseudoClass, got %v", sel.Type)
	}
	if sel.Value != "focus" {
		t.Fatalf("expected 'focus', got '%s'", sel.Value)
	}
}

func TestParseCompoundPseudoClass(t *testing.T) {
	// Compound: input:focus — the :focus is merged into the tag selector
	// (Note: current implementation merges specificity but doesn't store pseudo-class info)
	p := NewParser("input:focus { border-color: cyan; }")
	sheet, _ := p.Parse()
	sel := sheet.Rules[0].Selectors[0]
	// The compound selector becomes a SelectorTag
	if sel.Type != SelectorTag && sel.Type != SelectorPseudoClass {
		t.Fatalf("expected compound selector type, got %v", sel.Type)
	}
}

func TestParseHexColorValue(t *testing.T) {
	p := NewParser("div { color: #ff6600; }")
	sheet, _ := p.Parse()
	val := sheet.Rules[0].Declarations[0].Value
	if val.Type != ValueColor {
		t.Fatalf("expected ValueColor, got %v", val.Type)
	}
}

func TestParseEmptyStylesheet(t *testing.T) {
	p := NewParser("")
	sheet, _ := p.Parse()
	if len(sheet.Rules) != 0 {
		t.Fatalf("expected 0 rules, got %d", len(sheet.Rules))
	}
}
