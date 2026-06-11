// Package style resolves computed styles by applying CSS rules to DOM nodes.
package style

import (
	"strings"

	"github.com/elioneto/tuix/css"
	"github.com/elioneto/tuix/dom"
)

// ComputedStyle holds the final computed values for a DOM element.
type ComputedStyle struct {
	Display        DisplayType
	Position       PositionType
	Width          Length
	Height         Length
	MinWidth       Length
	MinHeight      Length
	MaxWidth       Length
	MaxHeight      Length
	MarginTop      Length
	MarginRight    Length
	MarginBottom   Length
	MarginLeft     Length
	PaddingTop     Length
	PaddingRight   Length
	PaddingBottom  Length
	PaddingLeft    Length
	BorderTop      Border
	BorderRight    Border
	BorderBottom   Border
	BorderLeft     Border
	Color          ColorValue
	Background     BackgroundValue
	FontSize       Length
	FontWeight     int
	TextAlign      TextAlignType
	OverflowX      OverflowType
	OverflowY      OverflowType
	Cursor         CursorType
	Opacity        float64
	ZIndex         int
	Visibility     VisibilityType
	WhiteSpace     WhiteSpaceType
	LineHeight     Length
}

// DisplayType represents the CSS display property.
type DisplayType int

const (
	DisplayInline DisplayType = iota
	DisplayBlock
	DisplayInlineBlock
	DisplayNone
	DisplayFlex
	DisplayInlineFlex
	DisplayGrid
)

// PositionType represents the CSS position property.
type PositionType int

const (
	PositionStatic PositionType = iota
	PositionRelative
	PositionAbsolute
	PositionFixed
	PositionSticky
)

// Length represents a length value in the computed style.
type Length struct {
	Value float64
	Unit  LengthUnit
}

// LengthUnit represents the unit of a length.
type LengthUnit int

const (
	LengthPx LengthUnit = iota
	LengthEm
	LengthRem
	LengthPercent
	LengthAuto
	LengthNone
)

// Border represents a CSS border.
type Border struct {
	Width Length
	Style BorderStyle
	Color ColorValue
}

// BorderStyle represents the CSS border-style property.
type BorderStyle int

const (
	BorderNone BorderStyle = iota
	BorderSolid
	BorderDashed
	BorderDotted
	BorderDouble
	BorderGroove
	BorderRidge
	BorderInset
	BorderOutset
)

// ColorValue represents a resolved color value.
type ColorValue struct {
	Defined bool
	R, G, B uint8
	A       float64
}

// BackgroundValue represents background properties.
type BackgroundValue struct {
	Color ColorValue
}

// TextAlignType represents text alignment.
type TextAlignType int

const (
	TextAlignLeft TextAlignType = iota
	TextAlignCenter
	TextAlignRight
	TextAlignJustify
)

// OverflowType represents overflow behavior.
type OverflowType int

const (
	OverflowVisible OverflowType = iota
	OverflowHidden
	OverflowScroll
	OverflowAuto
)

// CursorType represents cursor style.
type CursorType int

const (
	CursorAuto CursorType = iota
	CursorDefault
	CursorPointer
	CursorText
	CursorNone
	CursorHelp
)

// VisibilityType represents visibility.
type VisibilityType int

const (
	VisibilityVisible VisibilityType = iota
	VisibilityHidden
	VisibilityCollapse
)

// WhiteSpaceType represents whitespace handling.
type WhiteSpaceType int

const (
	WhiteSpaceNormal WhiteSpaceType = iota
	WhiteSpaceNowrap
	WhiteSpacePre
	WhiteSpacePreWrap
	WhiteSpacePreLine
)

// DefaultStyle returns a computed style with default values.
func DefaultStyle() ComputedStyle {
	return ComputedStyle{
		Display:    DisplayBlock,
		Position:   PositionStatic,
		Color:      ColorValue{Defined: false},
		Background: BackgroundValue{Color: ColorValue{Defined: false}},
		FontSize:   Length{Value: 16, Unit: LengthPx},
		FontWeight: 400,
		TextAlign:  TextAlignLeft,
		Visibility: VisibilityVisible,
		WhiteSpace: WhiteSpaceNormal,
		Opacity:    1.0,
		MinWidth:   Length{Unit: LengthAuto},
		MinHeight:  Length{Unit: LengthAuto},
		MaxWidth:   Length{Unit: LengthNone},
		MaxHeight:  Length{Unit: LengthNone},
	}
}

// Resolver resolves CSS styles for DOM nodes.
type Resolver struct {
	stylesheet *css.Stylesheet
}

// NewResolver creates a new style resolver.
func NewResolver(sheet *css.Stylesheet) *Resolver {
	return &Resolver{stylesheet: sheet}
}

// Resolve computes the style for a DOM node.
// For now, it uses a simplified approach: only matches selectors at the
// element level and applies property values directly.
func (r *Resolver) Resolve(node *dom.Node) ComputedStyle {
	style := DefaultStyle()

	if node.Type != dom.NodeElement {
		return style
	}

	// Collect matching rules
	type match struct {
		rule   *css.Rule
		index  int // Order of appearance (source order)
	}

	var matches []match

	for _, rule := range r.stylesheet.Rules {
		for _, sel := range rule.Selectors {
			if matchesSelector(node, sel) {
				matches = append(matches, match{rule: rule, index: len(matches)})
				break
			}
		}
	}

	// Sort by specificity and source order (simple bubble sort)
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			less := false
			si := matches[i].rule.Selectors[0].Specificity
			sj := matches[j].rule.Selectors[0].Specificity
			if si.Less(sj) {
				less = true
			} else if si.Equal(sj) {
				less = matches[i].index < matches[j].index
			}
			if less {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	// Apply declarations
	for _, m := range matches {
		for _, decl := range m.rule.Declarations {
			applyDeclaration(&style, decl)
		}
	}

	return style
}

// matchesSelector checks if a DOM node matches a CSS selector.
func matchesSelector(node *dom.Node, sel css.Selector) bool {
	if node.Type != dom.NodeElement {
		return false
	}

	switch sel.Type {
	case css.SelectorUniversal:
		return true

	case css.SelectorTag:
		return node.TagName() == sel.Value

	case css.SelectorClass:
		return node.HasClass(sel.Value)

	case css.SelectorID:
		return node.ID() == sel.Value

	case css.SelectorAttribute:
		// Parse attr=value or just attr
		if idx := strings.Index(sel.Value, "="); idx >= 0 {
			attr := sel.Value[:idx]
			op := ""
			val := sel.Value[idx+1:]
			// Check for operator prefix (like ~=, |=, etc.)
			if strings.HasSuffix(attr, "~") {
				op = "~="
				attr = attr[:len(attr)-1]
			} else if strings.HasSuffix(attr, "|") {
				op = "|="
				attr = attr[:len(attr)-1]
			} else if strings.HasSuffix(attr, "^") {
				op = "^="
				attr = attr[:len(attr)-1]
			} else if strings.HasSuffix(attr, "$") {
				op = "$="
				attr = attr[:len(attr)-1]
			} else if strings.HasSuffix(attr, "*") {
				op = "*="
				attr = attr[:len(attr)-1]
			} else {
				op = "="
			}

			nodeVal := node.GetAttribute(attr)
			switch op {
			case "=":
				return nodeVal == val
			case "~=":
				words := strings.Fields(nodeVal)
				for _, w := range words {
					if w == val {
						return true
					}
				}
				return false
			case "|=":
				return nodeVal == val || strings.HasPrefix(nodeVal, val+"-")
			case "^=":
				return strings.HasPrefix(nodeVal, val)
			case "$=":
				return strings.HasSuffix(nodeVal, val)
			case "*=":
				return strings.Contains(nodeVal, val)
			}
		} else {
			// Just check for attribute presence
			return node.GetAttribute(sel.Value) != "" || node.HasAttribute(sel.Value)
		}
	}

	return false
}

// applyDeclaration applies a CSS declaration to a computed style.
func applyDeclaration(style *ComputedStyle, decl *css.Declaration) {
	switch decl.Property {
	case "display":
		switch strings.ToLower(decl.Value.Keyword) {
		case "block":
			style.Display = DisplayBlock
		case "inline":
			style.Display = DisplayInline
		case "inline-block":
			style.Display = DisplayInlineBlock
		case "none":
			style.Display = DisplayNone
		case "flex":
			style.Display = DisplayFlex
		case "inline-flex":
			style.Display = DisplayInlineFlex
		case "grid":
			style.Display = DisplayGrid
		}

	case "position":
		switch strings.ToLower(decl.Value.Keyword) {
		case "static":
			style.Position = PositionStatic
		case "relative":
			style.Position = PositionRelative
		case "absolute":
			style.Position = PositionAbsolute
		case "fixed":
			style.Position = PositionFixed
		case "sticky":
			style.Position = PositionSticky
		}

	case "width":
		style.Width = cssLengthToStyleLength(decl.Value)

	case "height":
		style.Height = cssLengthToStyleLength(decl.Value)

	case "min-width":
		style.MinWidth = cssLengthToStyleLength(decl.Value)

	case "min-height":
		style.MinHeight = cssLengthToStyleLength(decl.Value)

	case "max-width":
		style.MaxWidth = cssLengthToStyleLength(decl.Value)

	case "max-height":
		style.MaxHeight = cssLengthToStyleLength(decl.Value)

	case "margin":
		applyMargin(style, decl.Value)

	case "margin-top":
		style.MarginTop = cssLengthToStyleLength(decl.Value)

	case "margin-right":
		style.MarginRight = cssLengthToStyleLength(decl.Value)

	case "margin-bottom":
		style.MarginBottom = cssLengthToStyleLength(decl.Value)

	case "margin-left":
		style.MarginLeft = cssLengthToStyleLength(decl.Value)

	case "padding":
		applyPadding(style, decl.Value)

	case "padding-top":
		style.PaddingTop = cssLengthToStyleLength(decl.Value)

	case "padding-right":
		style.PaddingRight = cssLengthToStyleLength(decl.Value)

	case "padding-bottom":
		style.PaddingBottom = cssLengthToStyleLength(decl.Value)

	case "padding-left":
		style.PaddingLeft = cssLengthToStyleLength(decl.Value)

	case "border":
		applyBorder(style, decl.Value)

	case "border-top":
		style.BorderTop = parseBorderValue(decl.Value)

	case "border-right":
		style.BorderRight = parseBorderValue(decl.Value)

	case "border-bottom":
		style.BorderBottom = parseBorderValue(decl.Value)

	case "border-left":
		style.BorderLeft = parseBorderValue(decl.Value)

	case "border-width":
		if decl.Value.Type == css.ValueList {
			vals := borderValues(decl.Value)
			if len(vals) >= 1 {
				style.BorderTop.Width = vals[0]
				style.BorderRight.Width = vals[0]
				style.BorderBottom.Width = vals[0]
				style.BorderLeft.Width = vals[0]
			}
			if len(vals) >= 2 {
				style.BorderRight.Width = vals[1]
				style.BorderLeft.Width = vals[1]
			}
			if len(vals) >= 3 {
				style.BorderBottom.Width = vals[2]
			}
			if len(vals) >= 4 {
				style.BorderLeft.Width = vals[3]
			}
		} else {
			l := cssLengthToStyleLength(decl.Value)
			style.BorderTop.Width = l
			style.BorderRight.Width = l
			style.BorderBottom.Width = l
			style.BorderLeft.Width = l
		}

	case "border-style":
		if decl.Value.Type == css.ValueList {
			styles := borderStyleValues(decl.Value)
			if len(styles) >= 1 {
				style.BorderTop.Style = styles[0]
				style.BorderRight.Style = styles[0]
				style.BorderBottom.Style = styles[0]
				style.BorderLeft.Style = styles[0]
			}
			if len(styles) >= 2 {
				style.BorderRight.Style = styles[1]
				style.BorderLeft.Style = styles[1]
			}
			if len(styles) >= 3 {
				style.BorderBottom.Style = styles[2]
			}
			if len(styles) >= 4 {
				style.BorderLeft.Style = styles[3]
			}
		}

	case "border-color":
		if decl.Value.Type == css.ValueList {
			colors := borderColorValues(decl.Value)
			if len(colors) >= 1 {
				style.BorderTop.Color = colors[0]
				style.BorderRight.Color = colors[0]
				style.BorderBottom.Color = colors[0]
				style.BorderLeft.Color = colors[0]
			}
			if len(colors) >= 2 {
				style.BorderRight.Color = colors[1]
				style.BorderLeft.Color = colors[1]
			}
			if len(colors) >= 3 {
				style.BorderBottom.Color = colors[2]
			}
			if len(colors) >= 4 {
				style.BorderLeft.Color = colors[3]
			}
		}

	case "color":
		style.Color = cssColorToColorValue(decl.Value)

	case "background", "background-color":
		style.Background.Color = cssColorToColorValue(decl.Value)

	case "font-size":
		style.FontSize = cssLengthToStyleLength(decl.Value)

	case "font-weight":
		switch strings.ToLower(decl.Value.Keyword) {
		case "normal":
			style.FontWeight = 400
		case "bold":
			style.FontWeight = 700
		case "lighter":
			if style.FontWeight > 100 {
				style.FontWeight -= 100
			}
		case "bolder":
			if style.FontWeight < 900 {
				style.FontWeight += 100
			}
		default:
			if decl.Value.Type == css.ValueNumber {
				style.FontWeight = int(decl.Value.Number)
			}
		}

	case "text-align":
		switch strings.ToLower(decl.Value.Keyword) {
		case "left":
			style.TextAlign = TextAlignLeft
		case "center":
			style.TextAlign = TextAlignCenter
		case "right":
			style.TextAlign = TextAlignRight
		case "justify":
			style.TextAlign = TextAlignJustify
		}

	case "overflow":
		o := parseOverflow(strings.ToLower(decl.Value.Keyword))
		style.OverflowX = o
		style.OverflowY = o

	case "overflow-x":
		style.OverflowX = parseOverflow(strings.ToLower(decl.Value.Keyword))

	case "overflow-y":
		style.OverflowY = parseOverflow(strings.ToLower(decl.Value.Keyword))

	case "cursor":
		switch strings.ToLower(decl.Value.Keyword) {
		case "auto":
			style.Cursor = CursorAuto
		case "default":
			style.Cursor = CursorDefault
		case "pointer":
			style.Cursor = CursorPointer
		case "text":
			style.Cursor = CursorText
		case "none":
			style.Cursor = CursorNone
		case "help":
			style.Cursor = CursorHelp
		}

	case "opacity":
		if decl.Value.Type == css.ValueNumber {
			style.Opacity = clampF(decl.Value.Number, 0, 1)
		}

	case "z-index":
		if decl.Value.Type == css.ValueNumber {
			style.ZIndex = int(decl.Value.Number)
		}

	case "visibility":
		switch strings.ToLower(decl.Value.Keyword) {
		case "visible":
			style.Visibility = VisibilityVisible
		case "hidden":
			style.Visibility = VisibilityHidden
		case "collapse":
			style.Visibility = VisibilityCollapse
		}

	case "white-space":
		switch strings.ToLower(decl.Value.Keyword) {
		case "normal":
			style.WhiteSpace = WhiteSpaceNormal
		case "nowrap":
			style.WhiteSpace = WhiteSpaceNowrap
		case "pre":
			style.WhiteSpace = WhiteSpacePre
		case "pre-wrap":
			style.WhiteSpace = WhiteSpacePreWrap
		case "pre-line":
			style.WhiteSpace = WhiteSpacePreLine
		}

	case "line-height":
		style.LineHeight = cssLengthToStyleLength(decl.Value)
	}
}

func applyMargin(style *ComputedStyle, v css.Value) {
	vals := lengthValues(v)
	switch len(vals) {
	case 1:
		style.MarginTop = vals[0]
		style.MarginRight = vals[0]
		style.MarginBottom = vals[0]
		style.MarginLeft = vals[0]
	case 2:
		style.MarginTop = vals[0]
		style.MarginRight = vals[1]
		style.MarginBottom = vals[0]
		style.MarginLeft = vals[1]
	case 3:
		style.MarginTop = vals[0]
		style.MarginRight = vals[1]
		style.MarginBottom = vals[2]
		style.MarginLeft = vals[1]
	case 4:
		style.MarginTop = vals[0]
		style.MarginRight = vals[1]
		style.MarginBottom = vals[2]
		style.MarginLeft = vals[3]
	}
}

func applyPadding(style *ComputedStyle, v css.Value) {
	vals := lengthValues(v)
	switch len(vals) {
	case 1:
		style.PaddingTop = vals[0]
		style.PaddingRight = vals[0]
		style.PaddingBottom = vals[0]
		style.PaddingLeft = vals[0]
	case 2:
		style.PaddingTop = vals[0]
		style.PaddingRight = vals[1]
		style.PaddingBottom = vals[0]
		style.PaddingLeft = vals[1]
	case 3:
		style.PaddingTop = vals[0]
		style.PaddingRight = vals[1]
		style.PaddingBottom = vals[2]
		style.PaddingLeft = vals[1]
	case 4:
		style.PaddingTop = vals[0]
		style.PaddingRight = vals[1]
		style.PaddingBottom = vals[2]
		style.PaddingLeft = vals[3]
	}
}

func applyBorder(style *ComputedStyle, v css.Value) {
	if v.Type == css.ValueList {
		for _, sv := range v.Values {
			if sv.Type == css.ValueLength || sv.Type == css.ValueKeyword {
				// Check if it's a width, style, or color
				switch strings.ToLower(sv.Keyword) {
				case "solid", "dashed", "dotted", "double", "groove", "ridge", "inset", "outset", "none":
					bs := parseBorderStyleKeyword(sv.Keyword)
					style.BorderTop.Style = bs
					style.BorderRight.Style = bs
					style.BorderBottom.Style = bs
					style.BorderLeft.Style = bs
				case "transparent", "currentcolor":
					cv := cssColorToColorValue(sv)
					style.BorderTop.Color = cv
					style.BorderRight.Color = cv
					style.BorderBottom.Color = cv
					style.BorderLeft.Color = cv
				default:
					if sv.Type == css.ValueColor {
						cv := cssColorToColorValue(sv)
						style.BorderTop.Color = cv
						style.BorderRight.Color = cv
						style.BorderBottom.Color = cv
						style.BorderLeft.Color = cv
					} else if sv.Type == css.ValueLength {
						l := cssLengthToStyleLength(sv)
						style.BorderTop.Width = l
						style.BorderRight.Width = l
						style.BorderBottom.Width = l
						style.BorderLeft.Width = l
					}
				}
			}
		}
	}
}

func parseBorderValue(v css.Value) Border {
	b := Border{Width: Length{Unit: LengthAuto}, Style: BorderNone, Color: ColorValue{Defined: false}}
	if v.Type == css.ValueList {
		for _, sv := range v.Values {
			switch {
			case sv.Type == css.ValueColor || sv.Color.Type != css.ColorHex || sv.Color.Hex != "":
				b.Color = cssColorToColorValue(sv)
			case sv.Type == css.ValueLength || sv.Keyword == "thin" || sv.Keyword == "medium" || sv.Keyword == "thick":
				b.Width = cssLengthToStyleLength(sv)
			case sv.Type == css.ValueKeyword:
				bs := parseBorderStyleKeyword(sv.Keyword)
				if bs != BorderNone {
					b.Style = bs
				}
			}
		}
	}
	return b
}

func parseBorderStyleKeyword(s string) BorderStyle {
	switch strings.ToLower(s) {
	case "none":
		return BorderNone
	case "solid":
		return BorderSolid
	case "dashed":
		return BorderDashed
	case "dotted":
		return BorderDotted
	case "double":
		return BorderDouble
	case "groove":
		return BorderGroove
	case "ridge":
		return BorderRidge
	case "inset":
		return BorderInset
	case "outset":
		return BorderOutset
	}
	return BorderNone
}

func parseOverflow(s string) OverflowType {
	switch s {
	case "visible":
		return OverflowVisible
	case "hidden":
		return OverflowHidden
	case "scroll":
		return OverflowScroll
	case "auto":
		return OverflowAuto
	}
	return OverflowVisible
}

// Helper functions for value extraction.
func cssLengthToStyleLength(v css.Value) Length {
	switch v.Type {
	case css.ValueLength:
		switch v.Length.Unit {
		case css.UnitPx:
			return Length{Value: v.Length.Value, Unit: LengthPx}
		case css.UnitEm:
			return Length{Value: v.Length.Value, Unit: LengthEm}
		case css.UnitRem:
			return Length{Value: v.Length.Value, Unit: LengthRem}
		case css.UnitPercent:
			return Length{Value: v.Length.Value, Unit: LengthPercent}
		}
	case css.ValuePercent:
		return Length{Value: v.Percent, Unit: LengthPercent}
	case css.ValueNumber:
		return Length{Value: v.Number, Unit: LengthPx}
	case css.ValueAuto:
		return Length{Unit: LengthAuto}
	case css.ValueKeyword:
		switch strings.ToLower(v.Keyword) {
		case "auto":
			return Length{Unit: LengthAuto}
		case "none":
			return Length{Unit: LengthNone}
		case "thin":
			return Length{Value: 1, Unit: LengthPx}
		case "medium":
			return Length{Value: 3, Unit: LengthPx}
		case "thick":
			return Length{Value: 5, Unit: LengthPx}
		}
	}
	return Length{Unit: LengthAuto}
}

func cssColorToColorValue(v css.Value) ColorValue {
	if v.Type == css.ValueColor {
		switch v.Color.Type {
		case css.ColorHex, css.ColorRGB, css.ColorNamed:
			return ColorValue{
				Defined: true,
				R:       v.Color.R,
				G:       v.Color.G,
				B:       v.Color.B,
				A:       v.Color.A,
			}
		case css.ColorTransparent:
			return ColorValue{Defined: true}
		}
	}
	return ColorValue{Defined: false}
}

func lengthValues(v css.Value) []Length {
	if v.Type == css.ValueList {
		rv := make([]Length, 0, len(v.Values))
		for _, sv := range v.Values {
			if sv.Type != css.ValueList {
				rv = append(rv, cssLengthToStyleLength(sv))
			}
		}
		return rv
	}
	return []Length{cssLengthToStyleLength(v)}
}

func borderStyleValues(v css.Value) []BorderStyle {
	if v.Type != css.ValueList {
		return nil
	}
	rv := make([]BorderStyle, 0, len(v.Values))
	for _, sv := range v.Values {
		rv = append(rv, parseBorderStyleKeyword(sv.Keyword))
	}
	return rv
}

func borderColorValues(v css.Value) []ColorValue {
	if v.Type != css.ValueList {
		return nil
	}
	rv := make([]ColorValue, 0, len(v.Values))
	for _, sv := range v.Values {
		rv = append(rv, cssColorToColorValue(sv))
	}
	return rv
}

func borderValues(v css.Value) []Length {
	return lengthValues(v)
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
