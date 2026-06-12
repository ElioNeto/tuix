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
	OutlineStyle   BorderStyle // CSS outline: none | solid (default solid when focused)

	// Flexbox properties
	FlexDirection  FlexDirectionType
	FlexWrap       FlexWrapType
	JustifyContent JustifyContentType
	AlignItems     AlignItemsType
	AlignContent   AlignContentType
	AlignSelf      AlignSelfType
	FlexGrow       float64
	FlexShrink     float64
	FlexBasis      Length
	Order          int
	RowGap         Length
	ColumnGap      Length
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

// FlexDirectionType represents flex-direction values.
type FlexDirectionType int

const (
	FlexDirectionRow FlexDirectionType = iota
	FlexDirectionRowReverse
	FlexDirectionColumn
	FlexDirectionColumnReverse
)

// FlexWrapType represents flex-wrap values.
type FlexWrapType int

const (
	FlexWrapNoWrap FlexWrapType = iota
	FlexWrapWrap
	FlexWrapWrapReverse
)

// JustifyContentType represents justify-content values.
type JustifyContentType int

const (
	JustifyContentFlexStart JustifyContentType = iota
	JustifyContentFlexEnd
	JustifyContentCenter
	JustifyContentSpaceBetween
	JustifyContentSpaceAround
	JustifyContentSpaceEvenly
)

// AlignItemsType represents align-items / align-self values.
type AlignItemsType int

const (
	AlignItemsStretch AlignItemsType = iota
	AlignItemsFlexStart
	AlignItemsFlexEnd
	AlignItemsCenter
	AlignItemsBaseline
)

// AlignContentType represents align-content values.
type AlignContentType int

const (
	AlignContentStretch AlignContentType = iota
	AlignContentFlexStart
	AlignContentFlexEnd
	AlignContentCenter
	AlignContentSpaceBetween
	AlignContentSpaceAround
)

// AlignSelfType represents the align-self property.
type AlignSelfType int

const (
	AlignSelfAuto AlignSelfType = iota
	AlignSelfFlexStart
	AlignSelfFlexEnd
	AlignSelfCenter
	AlignSelfStretch
	AlignSelfBaseline
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

		// Flex defaults
		FlexDirection:  FlexDirectionRow,
		FlexWrap:       FlexWrapNoWrap,
		JustifyContent: JustifyContentFlexStart,
		AlignItems:     AlignItemsStretch,
		AlignContent:   AlignContentFlexStart,
		AlignSelf:      AlignSelfAuto,
		FlexGrow:       0,
		FlexShrink:     1,
		FlexBasis:      Length{Unit: LengthAuto},
		Order:          0,

		// Default outline: solid (visible when focused)
		OutlineStyle: BorderSolid,
	}
}

// DefaultDisplayForTag returns the default display value for a given HTML tag name.
func DefaultDisplayForTag(tag string) DisplayType {
	switch strings.ToLower(tag) {
	case "a", "abbr", "acronym", "b", "bdi", "bdo", "big", "cite", "code",
		"dfn", "em", "i", "kbd", "label", "map", "mark", "output", "q",
		"ruby", "s", "samp", "small", "span", "strong", "sub", "sup",
		"time", "tt", "u", "var", "wbr":
		return DisplayInline
	case "input", "textarea", "select", "button", "option", "progress", "meter":
		return DisplayBlock
	default:
		return DisplayBlock
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

// Resolve computes the style for a DOM node, optionally inheriting from a parent style.
func (r *Resolver) Resolve(node *dom.Node) ComputedStyle {
	return r.ResolveWithParent(node, ComputedStyle{})
}

// ResolveWithParent computes the style, inheriting from the given parent style.
func (r *Resolver) ResolveWithParent(node *dom.Node, parent ComputedStyle) ComputedStyle {
	var style ComputedStyle

	if node.Type != dom.NodeElement {
		// Non-element nodes (text) inherit from parent
		style = parent
		// Ensure DisplayNone from parent is not inherited
		if style.Display == DisplayNone {
			style = DefaultStyle()
		}
		return style
	}

	// Start from default and apply tag-specific display
	style = DefaultStyle()
	style.Display = DefaultDisplayForTag(node.Data)

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

	// Sort by specificity and source order (ascending — lower specificity first,
	// so later rules override earlier ones).
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			si := matches[i].rule.Selectors[0].Specificity
			sj := matches[j].rule.Selectors[0].Specificity
			if sj.Less(si) || (si.Equal(sj) && matches[j].index < matches[i].index) {
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

	// CSS inheritance: inherit inherited properties from parent
	inheritStyle(&style, parent)

	return style
}

// matchesSelector checks if a DOM node matches a CSS selector.
func matchesSelector(node *dom.Node, sel css.Selector) bool {
	if node.Type != dom.NodeElement {
		return false
	}

	switch sel.Type {
	case css.SelectorUniversal:
		// May also have compound conditions
		if sel.Tag != "" && node.TagName() != sel.Tag {
			return false
		}
		if sel.ID != "" && node.ID() != sel.ID {
			return false
		}
		for _, cls := range sel.Classes {
			if !node.HasClass(cls) {
				return false
			}
		}
		return true

	case css.SelectorTag:
		if node.TagName() != sel.Value {
			return false
		}
		// Check compound conditions
		if sel.ID != "" && node.ID() != sel.ID {
			return false
		}
		for _, cls := range sel.Classes {
			if !node.HasClass(cls) {
				return false
			}
		}
		return true

	case css.SelectorClass:
		// Primary class must match
		if !node.HasClass(sel.Value) {
			return false
		}
		// Additional classes must also match
		for _, cls := range sel.Classes {
			if cls == sel.Value {
				continue
			}
			if !node.HasClass(cls) {
				return false
			}
		}
		// Check ID if present
		if sel.ID != "" && node.ID() != sel.ID {
			return false
		}
		// Check tag if present
		if sel.Tag != "" && node.TagName() != sel.Tag {
			return false
		}
		return true

	case css.SelectorID:
		if node.ID() != sel.Value {
			return false
		}
		// Check compound conditions
		for _, cls := range sel.Classes {
			if !node.HasClass(cls) {
				return false
			}
		}
		if sel.Tag != "" && node.TagName() != sel.Tag {
			return false
		}
		return true

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

	case css.SelectorPseudoClass:
		// Match pseudo-classes based on dynamic state
		// The 'focused' and 'hovered' attributes are set by tuix before each layout pass
		switch sel.Value {
		case "focus", "focus-visible":
			return node.HasAttribute("focused") || node.GetAttribute("focused") != ""
		case "focus-within":
			// Check if this node or any descendant has focused attribute
			return hasFocusWithin(node)
		case "hover":
			return node.HasAttribute("hovered") || node.GetAttribute("hovered") != ""
		default:
			return false
		}
	}

	return false
}

// hasFocusWithin checks if the node or any of its descendants has the "focused" attribute set.
func hasFocusWithin(node *dom.Node) bool {
	if node == nil {
		return false
	}
	if node.HasAttribute("focused") {
		return true
	}
	for _, child := range node.Children {
		if hasFocusWithin(child) {
			return true
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
		vals := borderValues(decl.Value)
		if len(vals) > 0 {
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
		styles := borderStyleValues(decl.Value)
		if len(styles) > 0 {
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
		} else {
			bs := parseBorderStyleKeyword(decl.Value.Keyword)
			style.BorderTop.Style = bs
			style.BorderRight.Style = bs
			style.BorderBottom.Style = bs
			style.BorderLeft.Style = bs
		}

	case "border-color":
		colors := borderColorValues(decl.Value)
		if len(colors) > 0 {
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
		} else {
			cv := cssColorToColorValue(decl.Value)
			style.BorderTop.Color = cv
			style.BorderRight.Color = cv
			style.BorderBottom.Color = cv
			style.BorderLeft.Color = cv
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

	case "outline":
		switch strings.ToLower(decl.Value.Keyword) {
		case "none":
			style.OutlineStyle = BorderNone
		case "solid":
			style.OutlineStyle = BorderSolid
		default:
			style.OutlineStyle = BorderSolid
		}

	case "outline-style":
		switch strings.ToLower(decl.Value.Keyword) {
		case "none":
			style.OutlineStyle = BorderNone
		case "solid":
			style.OutlineStyle = BorderSolid
		default:
			style.OutlineStyle = BorderSolid
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

	// --- Flexbox properties ---
	case "flex-direction":
		switch strings.ToLower(decl.Value.Keyword) {
		case "row":
			style.FlexDirection = FlexDirectionRow
		case "row-reverse":
			style.FlexDirection = FlexDirectionRowReverse
		case "column":
			style.FlexDirection = FlexDirectionColumn
		case "column-reverse":
			style.FlexDirection = FlexDirectionColumnReverse
		}

	case "flex-wrap":
		switch strings.ToLower(decl.Value.Keyword) {
		case "nowrap":
			style.FlexWrap = FlexWrapNoWrap
		case "wrap":
			style.FlexWrap = FlexWrapWrap
		case "wrap-reverse":
			style.FlexWrap = FlexWrapWrapReverse
		}

	case "flex-flow":
		// Shorthand for flex-direction and flex-wrap
		parts := valueList(decl.Value)
		for _, p := range parts {
			switch strings.ToLower(p.Keyword) {
			case "row":
				style.FlexDirection = FlexDirectionRow
			case "row-reverse":
				style.FlexDirection = FlexDirectionRowReverse
			case "column":
				style.FlexDirection = FlexDirectionColumn
			case "column-reverse":
				style.FlexDirection = FlexDirectionColumnReverse
			case "nowrap":
				style.FlexWrap = FlexWrapNoWrap
			case "wrap":
				style.FlexWrap = FlexWrapWrap
			case "wrap-reverse":
				style.FlexWrap = FlexWrapWrapReverse
			}
		}

	case "justify-content":
		switch strings.ToLower(decl.Value.Keyword) {
		case "flex-start":
			style.JustifyContent = JustifyContentFlexStart
		case "flex-end":
			style.JustifyContent = JustifyContentFlexEnd
		case "center":
			style.JustifyContent = JustifyContentCenter
		case "space-between":
			style.JustifyContent = JustifyContentSpaceBetween
		case "space-around":
			style.JustifyContent = JustifyContentSpaceAround
		case "space-evenly":
			style.JustifyContent = JustifyContentSpaceEvenly
		}

	case "align-items":
		switch strings.ToLower(decl.Value.Keyword) {
		case "stretch":
			style.AlignItems = AlignItemsStretch
		case "flex-start":
			style.AlignItems = AlignItemsFlexStart
		case "flex-end":
			style.AlignItems = AlignItemsFlexEnd
		case "center":
			style.AlignItems = AlignItemsCenter
		case "baseline":
			style.AlignItems = AlignItemsBaseline
		}

	case "align-content":
		switch strings.ToLower(decl.Value.Keyword) {
		case "stretch":
			style.AlignContent = AlignContentStretch
		case "flex-start":
			style.AlignContent = AlignContentFlexStart
		case "flex-end":
			style.AlignContent = AlignContentFlexEnd
		case "center":
			style.AlignContent = AlignContentCenter
		case "space-between":
			style.AlignContent = AlignContentSpaceBetween
		case "space-around":
			style.AlignContent = AlignContentSpaceAround
		}

	case "align-self":
		switch strings.ToLower(decl.Value.Keyword) {
		case "auto":
			style.AlignSelf = AlignSelfAuto
		case "flex-start":
			style.AlignSelf = AlignSelfFlexStart
		case "flex-end":
			style.AlignSelf = AlignSelfFlexEnd
		case "center":
			style.AlignSelf = AlignSelfCenter
		case "stretch":
			style.AlignSelf = AlignSelfStretch
		case "baseline":
			style.AlignSelf = AlignSelfBaseline
		}

	case "flex-grow":
		if decl.Value.Type == css.ValueNumber {
			style.FlexGrow = decl.Value.Number
		}

	case "flex-shrink":
		if decl.Value.Type == css.ValueNumber {
			style.FlexShrink = decl.Value.Number
		}

	case "flex-basis":
		style.FlexBasis = cssLengthToStyleLength(decl.Value)

	case "flex":
		// Shorthand: flex-grow flex-shrink flex-basis
		parts := valueList(decl.Value)
		if len(parts) > 0 {
			// First value: flex-grow
			if parts[0].Type == css.ValueNumber {
				style.FlexGrow = parts[0].Number
			}
		}
		if len(parts) > 1 {
			// Second value: flex-shrink
			if parts[1].Type == css.ValueNumber {
				style.FlexShrink = parts[1].Number
			}
		}
		if len(parts) > 2 {
			// Third value: flex-basis
			style.FlexBasis = cssLengthToStyleLength(parts[2])
		}

	case "order":
		if decl.Value.Type == css.ValueNumber {
			style.Order = int(decl.Value.Number)
		}

	case "gap":
		// Gap shorthand: row-gap column-gap
		vals := lengthValues(decl.Value)
		if len(vals) > 0 {
			style.RowGap = vals[0]
			style.ColumnGap = vals[0]
		}
		if len(vals) > 1 {
			style.ColumnGap = vals[1]
		}

	case "row-gap":
		style.RowGap = cssLengthToStyleLength(decl.Value)

	case "column-gap":
		style.ColumnGap = cssLengthToStyleLength(decl.Value)
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
	// Handle single value (e.g., border: solid, border: 1px)
	if v.Type == css.ValueKeyword && len(v.Values) == 0 {
		bs := parseBorderStyleKeyword(v.Keyword)
		if bs != BorderNone {
			style.BorderTop.Style = bs
			style.BorderRight.Style = bs
			style.BorderBottom.Style = bs
			style.BorderLeft.Style = bs
			// Default width if none set
			if style.BorderTop.Width.Unit == LengthAuto {
				w := Length{Value: 1, Unit: LengthPx}
				style.BorderTop.Width = w
				style.BorderRight.Width = w
				style.BorderBottom.Width = w
				style.BorderLeft.Width = w
			}
		}
		return
	}

	if v.Type == css.ValueColor {
		cv := cssColorToColorValue(v)
		style.BorderTop.Color = cv
		style.BorderRight.Color = cv
		style.BorderBottom.Color = cv
		style.BorderLeft.Color = cv
		return
	}

	if v.Type == css.ValueLength || v.Type == css.ValueNumber {
		l := cssLengthToStyleLength(v)
		style.BorderTop.Width = l
		style.BorderRight.Width = l
		style.BorderBottom.Width = l
		style.BorderLeft.Width = l
		return
	}

	// Multiple values (e.g., border: 1px solid red)
	parts := valueList(v)
	for _, sv := range parts {
			if sv.Type == css.ValueLength || sv.Type == css.ValueKeyword || sv.Type == css.ValueColor {
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
					} else if sv.Type == css.ValueLength || sv.Type == css.ValueNumber {
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

func parseBorderValue(v css.Value) Border {
	b := Border{Width: Length{Unit: LengthAuto}, Style: BorderNone, Color: ColorValue{Defined: false}}
	parts := valueList(v)
	for _, sv := range parts {
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

// valueList extracts sub-values from a Value that may contain multiple parts
// (space-separated or ValueList). Returns the sub-values or a slice with v itself.
func valueList(v css.Value) []css.Value {
	if v.Type == css.ValueList {
		return v.Values
	}
	if v.Type == css.ValueKeyword && len(v.Values) > 0 {
		return v.Values
	}
	return []css.Value{v}
}

func lengthValues(v css.Value) []Length {
	parts := valueList(v)
	rv := make([]Length, 0, len(parts))
	for _, sv := range parts {
		rv = append(rv, cssLengthToStyleLength(sv))
	}
	return rv
}

func borderStyleValues(v css.Value) []BorderStyle {
	parts := valueList(v)
	rv := make([]BorderStyle, 0, len(parts))
	for _, sv := range parts {
		rv = append(rv, parseBorderStyleKeyword(sv.Keyword))
	}
	return rv
}

func borderColorValues(v css.Value) []ColorValue {
	parts := valueList(v)
	rv := make([]ColorValue, 0, len(parts))
	for _, sv := range parts {
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

// inheritStyle applies CSS inheritance: for inherited properties that were not
// explicitly set on the child, copy the parent's value.
// CSS inherited properties: color, font-weight, font-size, text-align,
// visibility, cursor, white-space, line-height, opacity (partial).
func inheritStyle(child *ComputedStyle, parent ComputedStyle) {
	if !child.Color.Defined {
		child.Color = parent.Color
	}
	if child.FontWeight == 400 && parent.FontWeight != 400 {
		child.FontWeight = parent.FontWeight
	}
	if child.FontSize.Unit == LengthPx && child.FontSize.Value == 16 &&
		parent.FontSize.Unit != LengthPx || parent.FontSize.Value != 16 {
		child.FontSize = parent.FontSize
	}
	if child.TextAlign == TextAlignLeft && parent.TextAlign != TextAlignLeft {
		child.TextAlign = parent.TextAlign
	}
	if child.Visibility == VisibilityVisible && parent.Visibility != VisibilityVisible {
		child.Visibility = parent.Visibility
	}
	if child.Cursor == CursorAuto && parent.Cursor != CursorAuto {
		child.Cursor = parent.Cursor
	}
	if child.WhiteSpace == WhiteSpaceNormal && parent.WhiteSpace != WhiteSpaceNormal {
		child.WhiteSpace = parent.WhiteSpace
	}
}
