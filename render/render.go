// Package render provides the terminal canvas and rendering engine.
//
// It takes the layout box tree and paints it onto a character-based canvas,
// which is then output to the terminal using ANSI escape codes.
package render

import (
	"fmt"
	"sort"
	"strings"

	"github.com/elioneto/tuix/color"
	"github.com/elioneto/tuix/geometry"
	"github.com/elioneto/tuix/layout"
	"github.com/elioneto/tuix/style"
)

// Cell represents a single character cell on the terminal screen.
type Cell struct {
	Rune          rune
	Fg            color.Color
	Bg            color.Color
	FgSet, BgSet  bool
	Bold          bool
	Italic        bool
	Underline     bool
	Reverse       bool
}

// FontSizeMode indicates how a line's characters are scaled.
type FontSizeMode int

const (
	FontNormal   FontSizeMode = iota // 1×1 cell per character (default)
	FontDoubleWide                    // 2×1: DEC double-width (\x1b#6)
	FontDoubleHigh                    // 1×2: DEC double-height (\x1b#3 / \x1b#4)
	FontDoubleBoth                    // 2×2: both
)

// Canvas is a rectangular grid of character cells.
type Canvas struct {
	Cells       [][]Cell
	Width       int
	Height      int
	colorMode   int
	lineModes   []FontSizeMode // Per-line font size mode
}

// NewCanvas creates a new canvas with the given dimensions.
func NewCanvas(width, height, colorMode int) *Canvas {
	cells := make([][]Cell, height)
	for y := range cells {
		cells[y] = make([]Cell, width)
	}
	return &Canvas{
		Cells:     cells,
		Width:     width,
		Height:    height,
		colorMode: colorMode,
		lineModes: make([]FontSizeMode, height),
	}
}

// Clear resets all cells to the default state.
func (c *Canvas) Clear() {
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			c.Cells[y][x] = Cell{}
		}
		c.lineModes[y] = FontNormal
	}
}

// Set sets a cell at (x, y) with the given rune and style.
func (c *Canvas) Set(x, y int, r rune, fg, bg color.Color, attrs ...bool) {
	if x < 0 || x >= c.Width || y < 0 || y >= c.Height {
		return
	}

	bold := false
	italic := false
	underline := false
	if len(attrs) > 0 {
		bold = attrs[0]
	}
	if len(attrs) > 1 {
		italic = attrs[1]
	}
	if len(attrs) > 2 {
		underline = attrs[2]
	}

	c.Cells[y][x] = Cell{
		Rune:      r,
		Fg:        fg,
		Bg:        bg,
		FgSet:     true,
		BgSet:     true,
		Bold:      bold,
		Italic:    italic,
		Underline: underline,
	}
}

// SetRune sets only the character at (x, y) without changing style.
func (c *Canvas) SetRune(x, y int, r rune) {
	if x < 0 || x >= c.Width || y < 0 || y >= c.Height {
		return
	}
	c.Cells[y][x].Rune = r
}

// SetLineMode sets the font size mode for a row.
func (c *Canvas) SetLineMode(y int, mode FontSizeMode) {
	if y >= 0 && y < c.Height {
		c.lineModes[y] = mode
	}
}

// SetCursor sends ANSI code to position the cursor.
func SetCursor(x, y int) string {
	return fmt.Sprintf("\x1b[%d;%dH", y+1, x+1)
}

// Render produces the ANSI escape sequence string to paint the canvas
// onto the terminal. It uses differential rendering if old is provided.
func (c *Canvas) Render(old *Canvas) string {
	var buf strings.Builder

	// If old canvas has different dimensions, or this is the first render,
	// clear the entire screen first to remove leftover content.
	if old == nil || old.Width != c.Width || old.Height != c.Height {
		buf.WriteString("\x1b[2J") // Clear entire screen
	}

	// Move to home position
	buf.WriteString("\x1b[H")

	currentFg := color.Color{}
	currentBg := color.Color{}
	fgSet := false
	bgSet := false

	for y := 0; y < c.Height; y++ {
		// Output DEC line attribute for font size mode
		switch c.lineModes[y] {
		case FontDoubleWide:
			buf.WriteString("\x1b#6")
		case FontDoubleHigh:
			buf.WriteString("\x1b#3")
		case FontDoubleBoth:
			buf.WriteString("\x1b#6\x1b#3")
		}

		for x := 0; x < c.Width; x++ {
			cell := c.Cells[y][x]

			// Skip if unchanged (differential rendering)
			if old != nil && y < old.Height && x < old.Width {
				oldCell := old.Cells[y][x]
				if oldCell == cell {
					continue
				}
			}

			// Apply bold/normal BEFORE color to avoid terminals that reset
			// foreground on SGR 22 (normal intensity)
			if cell.Bold {
				buf.WriteString("\x1b[1m")
			} else {
				buf.WriteString("\x1b[22m")
			}

			// Apply styles
			if cell.FgSet && (cell.Fg != currentFg || !fgSet) {
				buf.WriteString(cell.Fg.ANSI(c.colorMode))
				currentFg = cell.Fg
				fgSet = true
			}
			if cell.BgSet && (cell.Bg != currentBg || !bgSet) {
				buf.WriteString(cell.Bg.ANSIBackground(c.colorMode))
				currentBg = cell.Bg
				bgSet = true
			}
			if cell.BgSet && (cell.Bg != currentBg || !bgSet) {
				buf.WriteString(cell.Bg.ANSIBackground(c.colorMode))
				currentBg = cell.Bg
				bgSet = true
			}

			if cell.Bold {
				buf.WriteString("\x1b[1m")
			} else {
				buf.WriteString("\x1b[22m")
				// Re-emit foreground color as \x1b[22m may reset it in some terminals
				if fgSet {
					buf.WriteString(currentFg.ANSI(c.colorMode))
				}
				if bgSet {
					buf.WriteString(currentBg.ANSIBackground(c.colorMode))
				}
			}
			if cell.Italic {
				buf.WriteString("\x1b[3m")
			}
			if cell.Underline {
				buf.WriteString("\x1b[4m")
			}

			// Write the character
			if cell.Rune == 0 {
				buf.WriteRune(' ')
			} else {
				buf.WriteRune(cell.Rune)
			}
		}

		// Newline at end of row (except last row)
		if y < c.Height-1 {
			buf.WriteString("\r\n")
		}
	}

	// Reset styles
	buf.WriteString("\x1b[0m")

	return buf.String()
}

// RenderFull renders the complete canvas without differential optimization.
func (c *Canvas) RenderFull() string {
	return c.Render(nil)
}

// ---------------------------------------------------------------------------
// Painter: converts layout boxes to canvas cells
// ---------------------------------------------------------------------------

// Painter takes a layout box tree and paints it onto a canvas.
type Painter struct {
	Canvas     *Canvas
	colorMode  int
	fgColor    color.Color
	bgColor    color.Color
	clipStack  []geometry.Rect // Stack of clip rectangles for overflow clipping
}

// NewPainter creates a new painter.
func NewPainter(canvas *Canvas, colorMode int) *Painter {
	return &Painter{
		Canvas:    canvas,
		colorMode: colorMode,
		fgColor:   color.NewTrue(200, 200, 200), // Default light gray
		bgColor:   color.NewTrue(0, 0, 0),       // Default black
	}
}

// Paint draws a layout box tree onto the canvas.
func (p *Painter) Paint(box *layout.Box) {
	if box == nil {
		return
	}
	p.paintBox(box)
}

func (p *Painter) paintBox(box *layout.Box) {
	if box == nil {
		return
	}

	// Don't paint boxes with no visible area
	if box.Rect.Width <= 0 || box.Rect.Height <= 0 {
		return
	}

	// Determine foreground and background colors
	fg := p.resolveForeground(box.Style)
	bg := p.resolveBackground(box.Style)

	// Paint background
	p.paintBackground(box, bg)

	switch box.Type {
	case layout.BoxText:
		p.paintText(box, fg, bg)
	case layout.BoxBlock, layout.BoxInline, layout.BoxRoot, layout.BoxFlex:
		p.paintElement(box, fg, bg)
	}

	// Check if this box is a scroll container
	isScrollY := box.Style.OverflowY == style.OverflowScroll || box.Style.OverflowY == style.OverflowAuto
	isScrollX := box.Style.OverflowX == style.OverflowScroll || box.Style.OverflowX == style.OverflowAuto
	shouldClipY := isScrollY || box.Style.OverflowY == style.OverflowHidden
	shouldClipX := isScrollX || box.Style.OverflowX == style.OverflowHidden
	clip := shouldClipX || shouldClipY

	// Only clip when content actually overflows
	if clip {
		contentH := box.ContentRect.Height
		contentW := box.ContentRect.Width
		if isScrollY && box.ScrollHeight <= contentH {
			isScrollY = false
		}
		if isScrollX && box.ScrollWidth <= contentW {
			isScrollX = false
		}
		clip = isScrollX || isScrollY
	}

	if clip {
		// Push clip rect = content area in canvas coordinates
		clipRect := geometry.Rect{
			X:      box.ContentRect.X,
			Y:      box.ContentRect.Y,
			Width:  box.ContentRect.Width,
			Height: box.ContentRect.Height,
		}
		p.clipStack = append(p.clipStack, clipRect)

		// Paint children with scroll offset, sorted by z-index
		children := zSortedChildren(box.Children)
		for _, child := range children {
			origX := child.Rect.X
			origY := child.Rect.Y
			origCX := child.ContentRect.X
			origCY := child.ContentRect.Y

			// Apply scroll offset
			child.Rect.X -= box.ScrollX
			child.Rect.Y -= box.ScrollY
			child.ContentRect.X -= box.ScrollX
			child.ContentRect.Y -= box.ScrollY

			p.paintBox(child)

			// Restore original positions
			child.Rect.X = origX
			child.Rect.Y = origY
			child.ContentRect.X = origCX
			child.ContentRect.Y = origCY
		}

		p.clipStack = p.clipStack[:len(p.clipStack)-1]

		// Paint scrollbar indicators
		if isScrollY && box.ScrollHeight > box.ContentRect.Height {
			p.paintScrollbarY(box)
		}
		if isScrollX && box.ScrollWidth > box.ContentRect.Width {
			p.paintScrollbarX(box)
		}
	} else {
		// Paint children sorted by z-index (no clipping)
		for _, child := range zSortedChildren(box.Children) {
			p.paintBox(child)
		}
	}
}

// zSortedChildren returns a copy of children sorted by ascending ZIndex.
// Elements with ZIndex == 0 preserve their original relative order (stable sort).
// Negative ZIndex elements paint first (behind), positive paint last (on top).
func zSortedChildren(children []*layout.Box) []*layout.Box {
	// First pass: check if any child has a non-zero ZIndex
	hasZIndex := false
	for _, c := range children {
		if c.ZIndex != 0 {
			hasZIndex = true
			break
		}
	}
	if !hasZIndex {
		return children
	}

	sorted := make([]*layout.Box, len(children))
	copy(sorted, children)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].ZIndex < sorted[j].ZIndex
	})
	return sorted
}

// paintScrollbarY paints a vertical scrollbar on the right side of the content area.
func (p *Painter) paintScrollbarY(box *layout.Box) {
	contentH := box.ContentRect.Height
	if contentH <= 0 {
		return
	}
	scrollH := box.ScrollHeight
	if scrollH <= contentH {
		return
	}

	barX := box.ContentRect.X + box.ContentRect.Width - 1
	if barX < 0 || barX >= p.Canvas.Width {
		return
	}

	// Calculate thumb position
	thumbSize := contentH * contentH / scrollH
	if thumbSize < 1 {
		thumbSize = 1
	}
	maxScroll := scrollH - contentH
	if maxScroll <= 0 {
		return
	}
	thumbPos := box.ScrollY * (contentH - thumbSize) / maxScroll

	scrollbarColor := color.NewTrue(80, 80, 80)
	thumbColor := color.NewTrue(140, 140, 140)

	for y := 0; y < contentH; y++ {
		cy := box.ContentRect.Y + y
		if cy < 0 || cy >= p.Canvas.Height {
			continue
		}

		if y >= thumbPos && y < thumbPos+thumbSize {
			p.Canvas.Set(barX, cy, '▓', thumbColor, color.Color{})
		} else {
			p.Canvas.Set(barX, cy, '│', scrollbarColor, color.Color{})
		}
	}
}

// paintScrollbarX paints a horizontal scrollbar at the bottom of the content area.
func (p *Painter) paintScrollbarX(box *layout.Box) {
	contentW := box.ContentRect.Width
	if contentW <= 1 {
		return
	}
	scrollW := box.ScrollWidth
	if scrollW <= contentW {
		return
	}

	barY := box.ContentRect.Y + box.ContentRect.Height - 1
	if barY < 0 || barY >= p.Canvas.Height {
		return
	}

	// Calculate thumb position
	thumbSize := contentW * contentW / scrollW
	if thumbSize < 1 {
		thumbSize = 1
	}
	maxScroll := scrollW - contentW
	if maxScroll <= 0 {
		return
	}
	thumbPos := box.ScrollX * (contentW - thumbSize) / maxScroll

	scrollbarColor := color.NewTrue(80, 80, 80)
	thumbColor := color.NewTrue(140, 140, 140)

	for x := 0; x < contentW; x++ {
		cx := box.ContentRect.X + x
		if cx < 0 || cx >= p.Canvas.Width {
			continue
		}

		if x >= thumbPos && x < thumbPos+thumbSize {
			p.Canvas.Set(cx, barY, '▓', thumbColor, color.Color{})
		} else {
			p.Canvas.Set(cx, barY, '─', scrollbarColor, color.Color{})
		}
	}
}

// clipRect intersects r with all active clip rects.
func (p *Painter) clipRect(r geometry.Rect) geometry.Rect {
	for i := len(p.clipStack) - 1; i >= 0; i-- {
		r = r.Intersect(p.clipStack[i])
		if r.IsEmpty() {
			break
		}
	}
	return r
}

// isVisible returns true if the point (x, y) is within all clip rects.
func (p *Painter) isVisible(x, y int) bool {
	for i := range p.clipStack {
		if !p.clipStack[i].Contains(x, y) {
			return false
		}
	}
	return true
}

func (p *Painter) paintBackground(box *layout.Box, bg color.Color) {
	if !box.Style.Background.Color.Defined {
		return
	}

	rect := p.clipRect(box.Rect)
	if rect.IsEmpty() {
		return
	}

	// Clip to canvas
	if rect.X < 0 {
		rect.Width += rect.X
		rect.X = 0
	}
	if rect.Y < 0 {
		rect.Height += rect.Y
		rect.Y = 0
	}
	if rect.X+rect.Width > p.Canvas.Width {
		rect.Width = p.Canvas.Width - rect.X
	}
	if rect.Y+rect.Height > p.Canvas.Height {
		rect.Height = p.Canvas.Height - rect.Y
	}

	for y := rect.Y; y < rect.Y+rect.Height; y++ {
		for x := rect.X; x < rect.X+rect.Width; x++ {
			if y >= 0 && y < p.Canvas.Height && x >= 0 && x < p.Canvas.Width {
				p.Canvas.Cells[y][x].Bg = bg
				p.Canvas.Cells[y][x].BgSet = true
			}
		}
	}
}

func (p *Painter) paintText(box *layout.Box, fg, bg color.Color) {
	text := strings.TrimSpace(box.Node.Data)
	if text == "" {
		return
	}

	// Get the clip stack to constrain drawing
	hasClip := len(p.clipStack) > 0

	// Determine the containing block for alignment.
	// The text box itself may be narrower than its parent; alignment should
	// happen relative to the parent's content area.
	alignBox := box
	if box.Parent != nil && (box.Parent.Type == layout.BoxBlock || box.Parent.Type == layout.BoxRoot) {
		alignBox = box.Parent
	}
	contentX := alignBox.ContentRect.X
	availableWidth := alignBox.ContentRect.Width
	if availableWidth <= 0 {
		contentX = box.Rect.X
		availableWidth = box.Rect.Width
	}

	// Determine text alignment — inherit from parent if not set
	align := box.Style.TextAlign
	if align == style.TextAlignLeft && box.Parent != nil && box.Parent.Style.TextAlign != style.TextAlignLeft {
		align = box.Parent.Style.TextAlign
	}

	contentY := box.ContentRect.Y

	// Apply font-size as line spacing: each "row" of text uses
	// ceil(fontSize / 16) actual terminal rows.
	lineSpan := 1
	if fs := box.Style.FontSize; fs.Value > 20 {
		lineSpan = int(fs.Value / 16)
		if lineSpan < 1 {
			lineSpan = 1
		}
		if lineSpan > 4 {
			lineSpan = 4
		}
	}

	// Split text by newlines to support forced line breaks.
	// Each segment is a "paragraph" laid out separately.
	segments := strings.Split(text, "\n")
	globalRow := 0

	for _, seg := range segments {
		words := strings.Fields(seg)
		if len(words) == 0 {
			globalRow += lineSpan
			continue
		}

		wordIndex := 0
		linesInSegment := 0
		for rowOffset := 0; wordIndex < len(words); rowOffset += lineSpan {
			// Calculate how many words fit on this line and the line width
			lineWords := make([]string, 0)
			lineWidth := 0
			for i := wordIndex; i < len(words); i++ {
				w := words[i]
				wlen := len(w)
				if len(lineWords) > 0 && lineWidth+1+wlen > availableWidth {
					break
				}
				lineWords = append(lineWords, w)
				if lineWidth > 0 {
					lineWidth++ // space
				}
				lineWidth += wlen
			}

			if len(lineWords) == 0 {
				wordIndex++
				continue
			}

			// Calculate start X based on alignment
			startX := contentX
			switch align {
			case style.TextAlignCenter:
				startX = contentX + (availableWidth-lineWidth)/2
			case style.TextAlignRight:
				startX = contentX + availableWidth - lineWidth
			}

			y := contentY + globalRow + rowOffset

			// Skip this line if it's not visible (check against clip stack)
			if hasClip && !p.isVisible(startX, y) {
				wordIndex += len(lineWords)
				linesInSegment++
				continue
			}

			currentX := startX

			for i, word := range lineWords {
				if i > 0 {
					// Space between words
					if !hasClip || p.isVisible(currentX, y) {
						if currentX >= 0 && currentX < p.Canvas.Width && y >= 0 && y < p.Canvas.Height {
							p.Canvas.Set(currentX, y, ' ', fg, bg)
						}
					}
					currentX++
				}
				for _, ch := range word {
					if !hasClip || p.isVisible(currentX, y) {
						if currentX >= 0 && currentX < p.Canvas.Width && y >= 0 && y < p.Canvas.Height {
							p.Canvas.Set(currentX, y, ch, fg, bg, box.Style.FontWeight >= 700, false, false)
						}
					}
					currentX++
				}
			}

			wordIndex += len(lineWords)
			linesInSegment++
		}

		globalRow += linesInSegment * lineSpan
	}
}

func (p *Painter) paintElement(box *layout.Box, fg, bg color.Color) {
	// Paint borders
	if box.Style.BorderTop.Style != style.BorderNone ||
		box.Style.BorderBottom.Style != style.BorderNone ||
		box.Style.BorderLeft.Style != style.BorderNone ||
		box.Style.BorderRight.Style != style.BorderNone {
		p.paintBorders(box)
	}

	// Paint focus ring (outline) for focused elements
	if box.Style.OutlineStyle != style.BorderNone && box.Node != nil && box.Node.HasAttribute("focused") {
		p.paintFocusRing(box)
	}
}

func (p *Painter) paintFocusRing(box *layout.Box) {
	rect := p.clipRect(box.Rect)
	if rect.IsEmpty() {
		return
	}

	// Use a bright color for the focus ring
	ringColor := color.NewTrue(0, 212, 170) // Cyan/teal

	top := rect.Y
	bottom := rect.Y + rect.Height - 1
	left := rect.X
	right := rect.X + rect.Width - 1

	// Top and bottom edges
	for x := left; x <= right; x++ {
		if x >= 0 && x < p.Canvas.Width && top >= 0 && top < p.Canvas.Height {
			p.Canvas.Set(x, top, '─', ringColor, color.Color{})
		}
		if x >= 0 && x < p.Canvas.Width && bottom >= 0 && bottom < p.Canvas.Height && bottom != top {
			p.Canvas.Set(x, bottom, '─', ringColor, color.Color{})
		}
	}
	// Left and right edges
	for y := top; y <= bottom; y++ {
		if left >= 0 && left < p.Canvas.Width && y >= 0 && y < p.Canvas.Height {
			p.Canvas.Set(left, y, '│', ringColor, color.Color{})
		}
		if right >= 0 && right < p.Canvas.Width && y >= 0 && y < p.Canvas.Height && right != left {
			p.Canvas.Set(right, y, '│', ringColor, color.Color{})
		}
	}
	// Corners
	if left >= 0 && left < p.Canvas.Width && top >= 0 && top < p.Canvas.Height {
		p.Canvas.Set(left, top, '┌', ringColor, color.Color{})
	}
	if right >= 0 && right < p.Canvas.Width && top >= 0 && top < p.Canvas.Height {
		p.Canvas.Set(right, top, '┐', ringColor, color.Color{})
	}
	if left >= 0 && left < p.Canvas.Width && bottom >= 0 && bottom < p.Canvas.Height {
		p.Canvas.Set(left, bottom, '└', ringColor, color.Color{})
	}
	if right >= 0 && right < p.Canvas.Width && bottom >= 0 && bottom < p.Canvas.Height {
		p.Canvas.Set(right, bottom, '┘', ringColor, color.Color{})
	}
}

func (p *Painter) paintBorders(box *layout.Box) {
	rect := p.clipRect(box.Rect)
	if rect.IsEmpty() {
		return
	}

	borderColor := color.NewTrue(100, 100, 100) // Default border color

	if box.Style.BorderTop.Color.Defined {
		borderColor = styleColorToColor(box.Style.BorderTop.Color)
	}

	top := rect.Y
	bottom := rect.Y + rect.Height - 1
	left := rect.X
	right := rect.X + rect.Width - 1

	if box.Style.BorderTop.Style != style.BorderNone {
		for x := left; x <= right && x < p.Canvas.Width; x++ {
			if top >= 0 && top < p.Canvas.Height && x >= 0 {
				p.Canvas.Set(x, top, '─', borderColor, color.Color{})
			}
		}
	}

	if box.Style.BorderBottom.Style != style.BorderNone {
		for x := left; x <= right && x < p.Canvas.Width; x++ {
			if bottom >= 0 && bottom < p.Canvas.Height && x >= 0 {
				p.Canvas.Set(x, bottom, '─', borderColor, color.Color{})
			}
		}
	}

	if box.Style.BorderLeft.Style != style.BorderNone {
		for y := top; y <= bottom && y < p.Canvas.Height; y++ {
			if left >= 0 && left < p.Canvas.Width && y >= 0 {
				p.Canvas.Set(left, y, '│', borderColor, color.Color{})
			}
		}
	}

	if box.Style.BorderRight.Style != style.BorderNone {
		for y := top; y <= bottom && y < p.Canvas.Height; y++ {
			if right >= 0 && right < p.Canvas.Width && y >= 0 {
				p.Canvas.Set(right, y, '│', borderColor, color.Color{})
			}
		}
	}

	// Draw corners
	if box.Style.BorderTop.Style != style.BorderNone &&
		box.Style.BorderLeft.Style != style.BorderNone {
		p.Canvas.Set(left, top, '┌', borderColor, color.Color{})
	}
	if box.Style.BorderTop.Style != style.BorderNone &&
		box.Style.BorderRight.Style != style.BorderNone {
		p.Canvas.Set(right, top, '┐', borderColor, color.Color{})
	}
	if box.Style.BorderBottom.Style != style.BorderNone &&
		box.Style.BorderLeft.Style != style.BorderNone {
		p.Canvas.Set(left, bottom, '└', borderColor, color.Color{})
	}
	if box.Style.BorderBottom.Style != style.BorderNone &&
		box.Style.BorderRight.Style != style.BorderNone {
		p.Canvas.Set(right, bottom, '┘', borderColor, color.Color{})
	}
}

func (p *Painter) resolveForeground(s style.ComputedStyle) color.Color {
	if s.Color.Defined {
		return styleColorToColor(s.Color)
	}
	return p.fgColor
}

func (p *Painter) resolveBackground(s style.ComputedStyle) color.Color {
	if s.Background.Color.Defined {
		return styleColorToColor(s.Background.Color)
	}
	return p.bgColor
}

func styleColorToColor(c style.ColorValue) color.Color {
	return color.NewTrue(c.R, c.G, c.B)
}
