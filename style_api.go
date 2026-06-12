package tuix

import (
	"fmt"
	"strings"

	"github.com/elioneto/tuix/color"
	"github.com/elioneto/tuix/geometry"
)

// BorderStyle defines the style of a border.
type BorderStyle int

const (
	BorderNone    BorderStyle = iota
	BorderNormal              // ┌─┐│└┘
	BorderRounded             // ╭─╮│╰╯
	BorderThick               // ┏━┓┃┗┛
	BorderDouble              // ╔═╗║╚╝
	BorderHidden
)

// BorderDef represents a complete border definition for Style.
type BorderDef struct {
	Style BorderStyle
	Fg    color.Color
	Bg    color.Color
}

// Style provides a programmatic API for text styling (port of Lip Gloss Style).
// It builds ANSI escape codes for terminal formatting.
type Style struct {
	fg        color.Color
	fgSet     bool
	bg        color.Color
	bgSet     bool
	bold      bool
	italic    bool
	underline bool
	width     int
	height    int
	padding   geometry.Edges
	margin    geometry.Edges
	border    BorderDef
}

// NewStyle creates a new empty Style.
func NewStyle() Style {
	return Style{}
}

// Foreground sets the foreground color.
func (s Style) Foreground(c color.Color) Style {
	s.fg = c
	s.fgSet = true
	return s
}

// Background sets the background color.
func (s Style) Background(c color.Color) Style {
	s.bg = c
	s.bgSet = true
	return s
}

// Bold sets bold attribute.
func (s Style) Bold(v bool) Style {
	s.bold = v
	return s
}

// Italic sets italic attribute.
func (s Style) Italic(v bool) Style {
	s.italic = v
	return s
}

// Underline sets underline attribute.
func (s Style) Underline(v bool) Style {
	s.underline = v
	return s
}

// Width sets the minimum width of the styled output.
func (s Style) Width(w int) Style {
	s.width = w
	return s
}

// Height sets the minimum height of the styled output.
func (s Style) Height(h int) Style {
	s.height = h
	return s
}

// Padding sets padding (top, right, bottom, left) or all sides if one value.
func (s Style) Padding(v ...int) Style {
	switch len(v) {
	case 1:
		s.padding = geometry.Edges{Top: float64(v[0]), Right: float64(v[0]), Bottom: float64(v[0]), Left: float64(v[0])}
	case 2:
		s.padding = geometry.Edges{Top: float64(v[0]), Right: float64(v[1]), Bottom: float64(v[0]), Left: float64(v[1])}
	case 4:
		s.padding = geometry.Edges{Top: float64(v[0]), Right: float64(v[1]), Bottom: float64(v[2]), Left: float64(v[3])}
	}
	return s
}

// Margin sets margin (top, right, bottom, left) or all sides if one value.
func (s Style) Margin(v ...int) Style {
	switch len(v) {
	case 1:
		s.margin = geometry.Edges{Top: float64(v[0]), Right: float64(v[0]), Bottom: float64(v[0]), Left: float64(v[0])}
	case 2:
		s.margin = geometry.Edges{Top: float64(v[0]), Right: float64(v[1]), Bottom: float64(v[0]), Left: float64(v[1])}
	case 4:
		s.margin = geometry.Edges{Top: float64(v[0]), Right: float64(v[1]), Bottom: float64(v[2]), Left: float64(v[3])}
	}
	return s
}

// Border sets the border style and optional colors.
func (s Style) Border(bs BorderStyle, colors ...color.Color) Style {
	s.border.Style = bs
	if len(colors) > 0 {
		s.border.Fg = colors[0]
	}
	if len(colors) > 1 {
		s.border.Bg = colors[1]
	}
	return s
}

// BorderForeground sets the border foreground color.
func (s Style) BorderForeground(c color.Color) Style {
	s.border.Fg = c
	return s
}

// BorderBackground sets the border background color.
func (s Style) BorderBackground(c color.Color) Style {
	s.border.Bg = c
	return s
}

// Render applies the style to the given text and returns a styled string.
func (s Style) Render(text string) string {
	var ansi strings.Builder

	// Open ANSI escape codes
	ansi.WriteString("\x1b[0m") // Reset first
	needsReset := false

	if s.bold {
		ansi.WriteString("\x1b[1m")
	}
	if s.italic {
		ansi.WriteString("\x1b[3m")
	}
	if s.underline {
		ansi.WriteString("\x1b[4m")
	}

	if s.fgSet {
		fmt.Fprintf(&ansi, "\x1b[38;2;%d;%d;%dm", s.fg.R, s.fg.G, s.fg.B)
	}
	if s.bgSet {
		fmt.Fprintf(&ansi, "\x1b[48;2;%d;%d;%dm", s.bg.R, s.bg.G, s.bg.B)
	}

	if s.bold || s.italic || s.underline || s.fgSet || s.bgSet {
		needsReset = true
	}

	// Build the content with padding
	content := text

	// Apply left padding
	if int(s.padding.Left) > 0 {
		content = strings.Repeat(" ", int(s.padding.Left)) + content
	}
	// Apply right padding
	if int(s.padding.Right) > 0 {
		content += strings.Repeat(" ", int(s.padding.Right))
	}

	// Apply width (pad or truncate)
	runes := []rune(content)
	if s.width > 0 && len(runes) < s.width {
		content = content + strings.Repeat(" ", s.width-len(runes))
	}

	// Apply top/bottom padding (vertical) via newlines
	topPad := strings.Repeat("\n", int(s.padding.Top))
	bottomPad := strings.Repeat("\n", int(s.padding.Bottom))

	// Apply border
	borderStr := ""
	if s.border.Style != BorderNone {
		borderStr = s.renderBorder(content)
	}

	result := strings.Repeat("\n", int(s.margin.Top))
	result += topPad
	if borderStr != "" {
		result += borderStr
	} else {
		if needsReset {
			result += ansi.String() + content + "\x1b[0m"
		} else {
			result += content
		}
	}
	result += bottomPad
	result += strings.Repeat("\n", int(s.margin.Bottom))

	return result
}

// renderBorder creates a bordered box around the content.
func (s Style) renderBorder(content string) string {
	// Determine the border characters based on style
	type borderChars struct {
		tl, tr, bl, br, h, v string
	}
	var bc borderChars

	switch s.border.Style {
	case BorderNormal:
		bc = borderChars{"┌", "┐", "└", "┘", "─", "│"}
	case BorderRounded:
		bc = borderChars{"╭", "╮", "╰", "╯", "─", "│"}
	case BorderThick:
		bc = borderChars{"┏", "┓", "┗", "┛", "━", "┃"}
	case BorderDouble:
		bc = borderChars{"╔", "╗", "╚", "╝", "═", "║"}
	default:
		return content
	}

	// Get content width
	lines := strings.Split(content, "\n")
	maxWidth := 0
	for _, line := range lines {
		w := len([]rune(line))
		if w > maxWidth {
			maxWidth = w
		}
	}

	// Build border ANSI
	openFg := ""
	closeFg := ""
	if s.border.Fg.Type == color.ColorTrue || s.border.Fg.Type == color.ColorANSI || s.border.Fg.Type == color.Color256 {
		openFg = fmt.Sprintf("\x1b[38;2;%d;%d;%dm", s.border.Fg.R, s.border.Fg.G, s.border.Fg.B)
		closeFg = "\x1b[39m"
	}

	// Build the bordered box
	var b strings.Builder
	// Top border
	b.WriteString(openFg)
	b.WriteString(bc.tl)
	for i := 0; i < maxWidth; i++ {
		b.WriteString(bc.h)
	}
	b.WriteString(bc.tr)
	b.WriteString(closeFg)
	b.WriteString("\n")

	// Content lines with side borders
	for _, line := range lines {
		b.WriteString(openFg)
		b.WriteString(bc.v)
		b.WriteString(closeFg)
		b.WriteString(line)
		padding := maxWidth - len([]rune(line))
		if padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
		}
		b.WriteString(openFg)
		b.WriteString(bc.v)
		b.WriteString(closeFg)
		b.WriteString("\n")
	}

	// Bottom border
	b.WriteString(openFg)
	b.WriteString(bc.bl)
	for i := 0; i < maxWidth; i++ {
		b.WriteString(bc.h)
	}
	b.WriteString(bc.br)
	b.WriteString(closeFg)

	return b.String()
}

// RenderHTML renders HTML content with the style applied.
// This is a simple wrapper that applies the style to the rendered HTML.
func (s Style) RenderHTML(html string) string {
	// For now, just apply the style as ANSI codes around the HTML.
	// In a more complete implementation, this would parse inline HTML.
	return s.Render(html)
}

// JoinHorizontal joins styled strings horizontally with optional separator.
func JoinHorizontal(strs ...string) string {
	var b strings.Builder
	for i, s := range strs {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(s)
	}
	return b.String()
}

// JoinVertical joins styled strings vertically (each on a new line).
func JoinVertical(strs ...string) string {
	return strings.Join(strs, "\n")
}
