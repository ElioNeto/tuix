// Package css provides a CSS parser that parses stylesheets into structured rules.
//
// It supports standard CSS syntax including selectors, properties, and values.
package css

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Stylesheet represents a parsed CSS stylesheet.
type Stylesheet struct {
	Rules []*Rule
}

// Rule represents a single CSS rule (selector + declarations).
type Rule struct {
	Selectors    []Selector
	Declarations []*Declaration
	Line         int // Source line (for error reporting)
}

// Selector represents a CSS selector.
type Selector struct {
	Type     SelectorType
	Value    string // Tag name, class name, ID, etc.
	Specificity Specificity
}

// SelectorType categorizes CSS selectors.
type SelectorType int

const (
	SelectorUniversal SelectorType = iota
	SelectorTag
	SelectorClass
	SelectorID
	SelectorAttribute
	SelectorPseudoClass
	SelectorPseudoElement
	SelectorDescendant
	SelectorChild
	SelectorAdjacentSibling
	SelectorGeneralSibling
)

// Specificity represents CSS specificity (a, b, c, d).
// a - inline style, b - IDs, c - classes/attributes/pseudo-classes, d - elements/pseudo-elements
type Specificity struct {
	A, B, C, D int
}

// Less returns true if this specificity is less than other.
func (s Specificity) Less(other Specificity) bool {
	if s.A != other.A {
		return s.A < other.A
	}
	if s.B != other.B {
		return s.B < other.B
	}
	if s.C != other.C {
		return s.C < other.C
	}
	return s.D < other.D
}

// Equal returns true if the specificity values are equal.
func (s Specificity) Equal(other Specificity) bool {
	return s.A == other.A && s.B == other.B &&
		s.C == other.C && s.D == other.D
}

// Declaration represents a single CSS property declaration.
type Declaration struct {
	Property string
	Value    Value
	Important bool
}

// Value represents a CSS value.
type Value struct {
	Type  ValueType
	// Raw value variants
	Keyword   string // For keyword values
	Color     Color  // For color values
	Length    Length // For length values
	Number    float64
	Percent   float64
	String    string
	Function  string // Function name for function values
	Args      []Value // Arguments for function values
	Values    []Value // For comma-separated or space-separated lists
	Separator string // "," or " "
}

// ValueType categorizes CSS values.
type ValueType int

const (
	ValueKeyword ValueType = iota
	ValueLength
	ValueNumber
	ValuePercent
	ValueColor
	ValueString
	ValueFunction
	ValueList
	ValueInherit
	ValueInitial
	ValueUnset
	ValueAuto
	ValueNone
)

// Unit represents a CSS length unit.
type Unit int

const (
	UnitPx Unit = iota
	UnitEm
	UnitRem
	UnitPercent
	UnitVw
	UnitVh
	UnitCh
	UnitPt
	UnitCm
	UnitMm
	UnitIn
	UnitPc
)

// Length represents a CSS length value.
type Length struct {
	Value float64
	Unit  Unit
}

// Color represents a CSS color value (parsed form).
type Color struct {
	Name  string
	Hex   string
	R, G, B uint8
	A      float64
	Type   ColorType
}

// ColorType categorizes color values.
type ColorType int

const (
	ColorHex ColorType = iota
	ColorRGB
	ColorRGBA
	ColorHSL
	ColorHSLA
	ColorNamed
	ColorTransparent
	ColorCurrent
	ColorANSI
)

// ---------------------------------------------------------------------------
// Parser
// ---------------------------------------------------------------------------

// Parser parses CSS source text into a Stylesheet.
type Parser struct {
	input    string
	pos      int
	line     int
	errors   []error
}

// NewParser creates a new CSS parser.
func NewParser(input string) *Parser {
	return &Parser{input: input, line: 1}
}

// Parse parses the CSS input and returns a Stylesheet.
func (p *Parser) Parse() (*Stylesheet, []error) {
	sheet := &Stylesheet{}
	p.errors = nil
	p.skipWhitespace()

	for p.pos < len(p.input) {
		// Skip comments
		p.skipComments()

		if p.pos >= len(p.input) {
			break
		}

		// Check for @-rules (ignore for now)
		if p.input[p.pos] == '@' {
			p.skipAtRule()
			continue
		}

		// Parse selectors
		selectors, ok := p.parseSelectors()
		if !ok {
			break
		}

		// Parse declarations block
		decls, ok := p.parseBlock()
		if !ok {
			break
		}

		rule := &Rule{
			Selectors:    selectors,
			Declarations: decls,
			Line:         p.line,
		}
		sheet.Rules = append(sheet.Rules, rule)
	}

	return sheet, p.errors
}

// parseSelectors parses a comma-separated list of selectors.
func (p *Parser) parseSelectors() ([]Selector, bool) {
	selectors := make([]Selector, 0)
	p.skipWhitespace()

	for {
		sel, ok := p.parseSelector()
		if !ok {
			p.error("invalid selector")
			return nil, false
		}
		selectors = append(selectors, sel)

		p.skipWhitespace()
		if p.pos >= len(p.input) {
			break
		}

		if p.input[p.pos] == ',' {
			p.pos++
			p.skipWhitespace()
			continue
		}
		break
	}

	// Expect '{' after selectors
	p.skipWhitespace()
	if p.pos >= len(p.input) || p.input[p.pos] != '{' {
		p.error("expected '{' after selector")
		return nil, false
	}

	return selectors, true
}

// parseSelector parses a single compound selector.
func (p *Parser) parseSelector() (Selector, bool) {
	parts := make([]Selector, 0)
	simple, ok := p.parseSimpleSelector()
	if !ok {
		return Selector{}, false
	}
	parts = append(parts, simple)

	for {
		p.skipWhitespace()
		if p.pos >= len(p.input) {
			break
		}

		ch := p.input[p.pos]
		if ch == ',' || ch == '{' {
			break
		}

		if ch == '>' || ch == '+' || ch == '~' {
			// Combinator
			combinator := ch
			p.pos++
			p.skipWhitespace()

			var ct SelectorType
			switch combinator {
			case '>':
				ct = SelectorChild
			case '+':
				ct = SelectorAdjacentSibling
			case '~':
				ct = SelectorGeneralSibling
			}

			parts = append(parts, Selector{Type: ct})

			next, ok := p.parseSimpleSelector()
			if !ok {
				return Selector{}, false
			}
			parts = append(parts, next)
		} else if ch == '.' || ch == '#' || ch == ':' || ch == '[' ||
			(ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			ch == '*' || ch == '|' {
			// Descendant combinator (implied by whitespace)
			if len(parts) > 0 {
				parts = append(parts, Selector{Type: SelectorDescendant})
			}
			next, ok := p.parseSimpleSelector()
			if !ok {
				return Selector{}, false
			}
			parts = append(parts, next)
		} else {
			break
		}
	}

	if len(parts) == 1 {
		return parts[0], true
	}

	// Combine into a sequence (we handle this via the selector list)
	// For now, return the first selector with full specificity
	return Selector{Type: SelectorTag, Value: "*"}, true
}

// parseSimpleSelector parses a single simple selector (tag, .class, #id, etc.).
func (p *Parser) parseSimpleSelector() (Selector, bool) {
	p.skipWhitespace()

	if p.pos >= len(p.input) {
		return Selector{}, false
	}

	ch := p.input[p.pos]

	switch {
	case ch == '*':
		p.pos++
		return Selector{Type: SelectorUniversal, Value: "*",
			Specificity: Specificity{D: 1}}, true

	case ch == '.':
		p.pos++
		name := p.parseIdent()
		if name == "" {
			return Selector{}, false
		}
		return Selector{Type: SelectorClass, Value: name,
			Specificity: Specificity{C: 1}}, true

	case ch == '#':
		p.pos++
		name := p.parseIdent()
		if name == "" {
			return Selector{}, false
		}
		return Selector{Type: SelectorID, Value: name,
			Specificity: Specificity{B: 1}}, true

	case ch == ':':
		p.pos++
		isPseudoElement := false
		if p.pos < len(p.input) && p.input[p.pos] == ':' {
			isPseudoElement = true
			p.pos++
		}
		name := p.parseIdent()
		if name == "" {
			return Selector{}, false
		}

		if isPseudoElement {
			return Selector{Type: SelectorPseudoElement, Value: name}, true
		}
		return Selector{Type: SelectorPseudoClass, Value: name,
			Specificity: Specificity{C: 1}}, true

	case ch == '[':
		return p.parseAttributeSelector()

	default:
		if isIdentStart(ch) {
			name := p.parseIdent()
			if name == "" {
				return Selector{}, false
			}
			// Check for namespace (tag|name)
			if p.pos < len(p.input) && p.input[p.pos] == '|' {
				p.pos++
				local := p.parseIdent()
				name = name + "|" + local
			}
			return Selector{Type: SelectorTag, Value: name,
				Specificity: Specificity{D: 1}}, true
		}
	}

	return Selector{}, false
}

// parseAttributeSelector parses [attr] or [attr=value] style selectors.
func (p *Parser) parseAttributeSelector() (Selector, bool) {
	p.pos++ // skip '['
	p.skipWhitespace()

	attr := p.parseIdent()
	if attr == "" {
		return Selector{}, false
	}

	s := Selector{Type: SelectorAttribute, Value: attr,
		Specificity: Specificity{C: 1}}

	p.skipWhitespace()

	if p.pos < len(p.input) && p.input[p.pos] == '=' {
		p.pos++
		// Check for prefix operator
		// (already consumed the '=')
		p.skipWhitespace()
		val := p.parseAttributeValue()
		s.Value = attr + "=" + val
	} else if p.pos+1 < len(p.input) &&
		(p.input[p.pos] == '~' && p.input[p.pos+1] == '=') {
		p.pos += 2
		p.skipWhitespace()
		val := p.parseAttributeValue()
		s.Value = attr + "~=" + val
	} else if p.pos+1 < len(p.input) &&
		(p.input[p.pos] == '|' && p.input[p.pos+1] == '=') {
		p.pos += 2
		p.skipWhitespace()
		val := p.parseAttributeValue()
		s.Value = attr + "|=" + val
	} else if p.pos+1 < len(p.input) &&
		(p.input[p.pos] == '^' && p.input[p.pos+1] == '=') {
		p.pos += 2
		p.skipWhitespace()
		val := p.parseAttributeValue()
		s.Value = attr + "^=" + val
	} else if p.pos+1 < len(p.input) &&
		(p.input[p.pos] == '$' && p.input[p.pos+1] == '=') {
		p.pos += 2
		p.skipWhitespace()
		val := p.parseAttributeValue()
		s.Value = attr + "$=" + val
	} else if p.pos+1 < len(p.input) &&
		(p.input[p.pos] == '*' && p.input[p.pos+1] == '=') {
		p.pos += 2
		p.skipWhitespace()
		val := p.parseAttributeValue()
		s.Value = attr + "*=" + val
	}

	p.skipWhitespace()
	if p.pos >= len(p.input) || p.input[p.pos] != ']' {
		return Selector{}, false
	}
	p.pos++ // skip ']'

	return s, true
}

func (p *Parser) parseAttributeValue() string {
	if p.pos >= len(p.input) {
		return ""
	}
	if p.input[p.pos] == '"' || p.input[p.pos] == '\'' {
		quote := p.input[p.pos]
		p.pos++
		start := p.pos
		for p.pos < len(p.input) && p.input[p.pos] != quote {
			p.pos++
		}
		val := p.input[start:p.pos]
		if p.pos < len(p.input) {
			p.pos++ // skip closing quote
		}
		return val
	}
	return p.parseIdent()
}

// parseBlock parses a { ... } declaration block.
func (p *Parser) parseBlock() ([]*Declaration, bool) {
	if p.pos >= len(p.input) || p.input[p.pos] != '{' {
		return nil, false
	}
	p.pos++ // skip '{'

	decls := make([]*Declaration, 0)
	p.skipWhitespace()
	p.skipComments()

	for p.pos < len(p.input) && p.input[p.pos] != '}' {
		decl, ok := p.parseDeclaration()
		if ok {
			decls = append(decls, decl)
		}
		p.skipWhitespace()
		p.skipComments()

		// Skip semicolons
		for p.pos < len(p.input) && p.input[p.pos] == ';' {
			p.pos++
			p.skipWhitespace()
			p.skipComments()
		}
	}

	if p.pos >= len(p.input) {
		p.error("unexpected end of input in declaration block")
		return decls, false
	}
	p.pos++ // skip '}'
	return decls, true
}

// parseDeclaration parses a single property: value; declaration.
func (p *Parser) parseDeclaration() (*Declaration, bool) {
	p.skipWhitespace()
	prop := p.parseIdent()
	if prop == "" {
		return nil, false
	}

	p.skipWhitespace()
	if p.pos >= len(p.input) || p.input[p.pos] != ':' {
		return nil, false
	}
	p.pos++ // skip ':'
	p.skipWhitespace()

	// Parse the value
	val, ok := p.parseValue()
	if !ok {
		// Skip until ; or }
		for p.pos < len(p.input) && p.input[p.pos] != ';' && p.input[p.pos] != '}' {
			p.pos++
		}
		return nil, false
	}

	// Check for !important
	p.skipWhitespace()
	important := false
	if p.pos+9 < len(p.input) &&
		strings.EqualFold(p.input[p.pos:p.pos+9], "!important") {
		important = true
		p.pos += 10 // skip "!important" and whitespace
	}

	return &Declaration{
		Property:  prop,
		Value:     val,
		Important: important,
	}, true
}

// parseValue parses a CSS value.
func (p *Parser) parseValue() (Value, bool) {
	p.skipWhitespace()
	if p.pos >= len(p.input) {
		return Value{}, false
	}

	values := make([]Value, 0)

	for {
		p.skipWhitespace()
		if p.pos >= len(p.input) {
			break
		}

		ch := p.input[p.pos]

		// Check for end of value
		if ch == ';' || ch == '}' || ch == ')' || ch == ']' {
			break
		}

		// Check for comma (separator in lists)
		if ch == ',' {
			break
		}

		v, ok := p.parseSingleValue()
		if !ok {
			break
		}
		values = append(values, v)

		// Check for comma
		p.skipWhitespace()
		if p.pos < len(p.input) && p.input[p.pos] == ',' {
			// Comma-separated list; collect all values with this separator
			allValues := values
			allValues = append(allValues, Value{})
			p.pos++
			for {
				p.skipWhitespace()
				v, ok := p.parseSingleValue()
				if !ok {
					break
				}
				allValues = append(allValues, v)
				p.skipWhitespace()
				if p.pos < len(p.input) && p.input[p.pos] == ',' {
					p.pos++
					allValues = append(allValues, Value{})
				} else {
					break
				}
			}
			values = []Value{{
				Type:      ValueKeyword,
				Values:    allValues,
				Separator: ",",
			}}
			break
		}
	}

	if len(values) == 0 {
		return Value{}, false
	}

	if len(values) == 1 {
		return values[0], true
	}

	return Value{
		Type:      ValueKeyword,
		Values:    values,
		Separator: " ",
	}, true
}

// parseSingleValue parses a single CSS value token.
func (p *Parser) parseSingleValue() (Value, bool) {
	p.skipWhitespace()

	if p.pos >= len(p.input) {
		return Value{}, false
	}

	ch := p.input[p.pos]

	// Number (including lengths and percentages)
	if ch == '-' || ch == '+' || (ch >= '0' && ch <= '9') || ch == '.' {
		return p.parseNumericValue()
	}

	// Hash color
	if ch == '#' {
		return p.parseHashColor()
	}

	// Function
	if isIdentStart(ch) {
		ident := p.parseIdent()

		if p.pos < len(p.input) && p.input[p.pos] == '(' {
			return p.parseFunctionValue(ident)
		}

		// Keyword or special value
		switch strings.ToLower(ident) {
		case "inherit":
			return Value{Type: ValueInherit, Keyword: ident}, true
		case "initial":
			return Value{Type: ValueInitial, Keyword: ident}, true
		case "unset":
			return Value{Type: ValueUnset, Keyword: ident}, true
		case "auto":
			return Value{Type: ValueAuto, Keyword: ident}, true
		case "none":
			return Value{Type: ValueNone, Keyword: ident}, true
		default:
			return Value{Type: ValueKeyword, Keyword: ident}, true
		}
	}

	// String
	if ch == '"' || ch == '\'' {
		return p.parseStringValue()
	}

	// URL (url(...))
	if ch == 'u' || ch == 'U' {
		saved := p.pos
		ident := p.parseIdent()
		if strings.ToLower(ident) == "url" && p.pos < len(p.input) && p.input[p.pos] == '(' {
			return p.parseURLValue()
		}
		p.pos = saved
	}

	return Value{}, false
}

func (p *Parser) parseNumericValue() (Value, bool) {
	saved := p.pos

	// Sign
	if p.pos < len(p.input) && (p.input[p.pos] == '-' || p.input[p.pos] == '+') {
		p.pos++
	}

	// Integer part
	hasDigits := false
	for p.pos < len(p.input) && p.input[p.pos] >= '0' && p.input[p.pos] <= '9' {
		p.pos++
		hasDigits = true
	}

	// Decimal part
	if p.pos < len(p.input) && p.input[p.pos] == '.' {
		p.pos++
		for p.pos < len(p.input) && p.input[p.pos] >= '0' && p.input[p.pos] <= '9' {
			p.pos++
			hasDigits = true
		}
	}

	if !hasDigits {
		p.pos = saved
		return Value{}, false
	}

	numStr := p.input[saved:p.pos]
	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		p.pos = saved
		return Value{}, false
	}

	// Check for unit
	if p.pos < len(p.input) && isIdentStart(p.input[p.pos]) {
		unitStr := p.parseIdent()
		unit, ok := parseUnit(unitStr)
		if ok {
			return Value{
				Type:   ValueLength,
				Length: Length{Value: val, Unit: unit},
				Number: val,
			}, true
		}
		// Unknown unit, treat as keyword
		p.pos = saved
		return Value{}, false
	}

	// Check for percentage
	if p.pos < len(p.input) && p.input[p.pos] == '%' {
		p.pos++
		return Value{Type: ValuePercent, Percent: val, Number: val}, true
	}

	// Bare number
	return Value{Type: ValueNumber, Number: val}, true
}

func (p *Parser) parseHashColor() (Value, bool) {
	p.pos++ // skip '#'
	start := p.pos
	for p.pos < len(p.input) && isHexChar(p.input[p.pos]) {
		p.pos++
	}
	hex := p.input[start:p.pos]
	if len(hex) != 3 && len(hex) != 6 {
		// Not a valid hex color, treat as unknown
		p.pos = start
		return Value{}, false
	}

	return Value{
		Type:  ValueColor,
		Color: Color{Type: ColorHex, Hex: "#" + hex},
	}, true
}

func (p *Parser) parseFunctionValue(name string) (Value, bool) {
	p.pos++ // skip '('

	args := make([]Value, 0)
	p.skipWhitespace()
	for p.pos < len(p.input) && p.input[p.pos] != ')' {
		val, ok := p.parseValue()
		if ok {
			args = append(args, val)
		}
		p.skipWhitespace()
		if p.pos < len(p.input) && p.input[p.pos] == ',' {
			p.pos++
			p.skipWhitespace()
		}
	}

	if p.pos >= len(p.input) {
		return Value{}, false
	}
	p.pos++ // skip ')'

	nameLower := strings.ToLower(name)

	// Handle color functions
	switch nameLower {
	case "rgb", "rgba":
		if len(args) >= 3 {
			r := parseNumberArg(args[0])
			g := parseNumberArg(args[1])
			b := parseNumberArg(args[2])
			a := 1.0
			if len(args) >= 4 {
				a = parseNumberArg(args[3])
			}
			if r >= 0 && g >= 0 && b >= 0 {
				return Value{
					Type: ValueColor,
					Color: Color{
						Type: ColorRGB,
						R:    uint8(clampF(r, 0, 255)),
						G:    uint8(clampF(g, 0, 255)),
						B:    uint8(clampF(b, 0, 255)),
						A:    clampF(a, 0, 1),
					},
				}, true
			}
		}
	case "hsl", "hsla":
		// Basic handling
		if len(args) >= 3 {
			return Value{
				Type: ValueColor,
				Color: Color{
					Type: ColorHSL,
				},
			}, true
		}
	}

	return Value{
		Type:     ValueFunction,
		Function: nameLower,
		Args:     args,
	}, true
}

func (p *Parser) parseStringValue() (Value, bool) {
	quote := p.input[p.pos]
	p.pos++
	start := p.pos
	for p.pos < len(p.input) && p.input[p.pos] != quote {
		if p.input[p.pos] == '\\' {
			p.pos++
			if p.pos < len(p.input) {
				p.pos++
			}
		} else {
			p.pos++
		}
	}
	str := p.input[start:p.pos]
	if p.pos < len(p.input) {
		p.pos++ // skip closing quote
	}
	return Value{Type: ValueString, String: str}, true
}

func (p *Parser) parseURLValue() (Value, bool) {
	// Skip 'url('
	p.pos += 4 // "url("

	p.skipWhitespace()
	url := ""
	if p.pos < len(p.input) && (p.input[p.pos] == '"' || p.input[p.pos] == '\'') {
		// Quoted URL
		v, _ := p.parseStringValue()
		url = v.String
	} else {
		// Unquoted URL
		start := p.pos
		for p.pos < len(p.input) && p.input[p.pos] != ')' && !unicode.IsSpace(rune(p.input[p.pos])) {
			p.pos++
		}
		url = p.input[start:p.pos]
	}

	p.skipWhitespace()
	if p.pos < len(p.input) && p.input[p.pos] == ')' {
		p.pos++
	}

	return Value{
		Type:   ValueKeyword,
		Keyword: "url",
		Values: []Value{{Type: ValueString, String: url}},
	}, true
}

// parseIdent parses a CSS identifier.
func (p *Parser) parseIdent() string {
	start := p.pos
	if p.pos < len(p.input) && isIdentStart(p.input[p.pos]) {
		p.pos++
		for p.pos < len(p.input) && isIdentChar(p.input[p.pos]) {
			p.pos++
		}
		return p.input[start:p.pos]
	}
	return ""
}

// skipWhitespace skips whitespace characters and tracks line numbers.
func (p *Parser) skipWhitespace() {
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == '\n' {
			p.line++
			p.pos++
		} else if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\f' {
			p.pos++
		} else {
			break
		}
	}
}

// skipComments skips /* ... */ comments.
func (p *Parser) skipComments() {
	for p.pos+1 < len(p.input) && p.input[p.pos] == '/' && p.input[p.pos+1] == '*' {
		p.pos += 2
		for p.pos < len(p.input) {
			if p.input[p.pos] == '*' && p.pos+1 < len(p.input) && p.input[p.pos+1] == '/' {
				p.pos += 2
				break
			}
			if p.input[p.pos] == '\n' {
				p.line++
			}
			p.pos++
		}
	}
}

// skipAtRule skips @-rules (like @media, @keyframes, etc.).
func (p *Parser) skipAtRule() {
	for p.pos < len(p.input) && p.input[p.pos] != '{' && p.input[p.pos] != ';' {
		if p.input[p.pos] == '\n' {
			p.line++
		}
		p.pos++
	}
	if p.pos < len(p.input) && p.input[p.pos] == '{' {
		depth := 1
		p.pos++
		for p.pos < len(p.input) && depth > 0 {
			if p.input[p.pos] == '{' {
				depth++
			} else if p.input[p.pos] == '}' {
				depth--
			} else if p.input[p.pos] == '\n' {
				p.line++
			}
			p.pos++
		}
	} else if p.pos < len(p.input) && p.input[p.pos] == ';' {
		p.pos++
	}
}

func (p *Parser) error(msg string) {
	p.errors = append(p.errors, fmt.Errorf("line %d: %s", p.line, msg))
}

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

func isIdentStart(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') ||
		b == '_' || b == '-' || b > 0x7f
}

func isIdentChar(b byte) bool {
	return isIdentStart(b) || (b >= '0' && b <= '9')
}

func isHexChar(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
}

func parseUnit(s string) (Unit, bool) {
	switch strings.ToLower(s) {
	case "px":
		return UnitPx, true
	case "em":
		return UnitEm, true
	case "rem":
		return UnitRem, true
	case "%":
		return UnitPercent, true
	case "vw":
		return UnitVw, true
	case "vh":
		return UnitVh, true
	case "ch":
		return UnitCh, true
	case "pt":
		return UnitPt, true
	case "cm":
		return UnitCm, true
	case "mm":
		return UnitMm, true
	case "in":
		return UnitIn, true
	case "pc":
		return UnitPc, true
	}
	return UnitPx, false
}

func parseNumberArg(v Value) float64 {
	switch v.Type {
	case ValueNumber:
		return v.Number
	case ValueLength:
		return v.Length.Value
	case ValuePercent:
		return v.Percent
	}
	return 0
}

func clampF(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
