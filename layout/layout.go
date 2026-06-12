// Package layout implements the box model and layout engine.
//
// It takes styled DOM nodes and produces a tree of layout boxes with
// computed positions and sizes. The layout engine supports block and
// inline formatting contexts.
package layout

import (
	"strings"

	"github.com/elioneto/tuix/dom"
	"github.com/elioneto/tuix/geometry"
	"github.com/elioneto/tuix/style"
)

// BoxType distinguishes between different types of layout boxes.
type BoxType int

const (
	BoxBlock BoxType = iota
	BoxInline
	BoxAnonymous
	BoxText
	BoxRoot
)

// Box represents a rectangular layout box in the formatting tree.
type Box struct {
	Type      BoxType
	Node      *dom.Node
	Style     style.ComputedStyle
	Rect      geometry.Rect
	Children  []*Box
	Parent    *Box

	// Content area (inside padding and border)
	ContentRect geometry.Rect

	// Box model edges
	Margin    geometry.Edges
	Border    geometry.Edges
	Padding   geometry.Edges

	// Computed dimensions
	ComputedWidth  float64
	ComputedHeight float64
}

// TotalWidth returns the full width including margin, border, and padding.
func (b *Box) TotalWidth() float64 {
	return b.Margin.Left + b.Border.Left + b.Padding.Left +
		b.ComputedWidth +
		b.Padding.Right + b.Border.Right + b.Margin.Right
}

// TotalHeight returns the full height including margin, border, and padding.
func (b *Box) TotalHeight() float64 {
	return b.Margin.Top + b.Border.Top + b.Padding.Top +
		b.ComputedHeight +
		b.Padding.Bottom + b.Border.Bottom + b.Margin.Bottom
}

// LayoutEngine takes a styled DOM tree and produces a box tree with positions.
type LayoutEngine struct {
	Root       *Box
	ViewWidth  float64
	ViewHeight float64
}

// NewEngine creates a new layout engine.
func NewEngine() *LayoutEngine {
	return &LayoutEngine{}
}

// Layout performs layout on the document and returns the root box.
func (e *LayoutEngine) Layout(doc *dom.Node, resolvers map[*dom.Node]style.ComputedStyle) *Box {
	e.Root = e.buildBoxTree(doc, resolvers, nil)
	e.Root.Type = BoxRoot

	// Set viewport size
	// The root box fills the viewport
	if e.ViewWidth > 0 {
		e.Root.ComputedWidth = e.ViewWidth
	} else {
		e.Root.ComputedWidth = 80 // Default: 80 columns
	}
	if e.ViewHeight > 0 {
		e.Root.ComputedHeight = e.ViewHeight
	} else {
		e.Root.ComputedHeight = 24 // Default: 24 rows
	}
	e.Root.Rect.Width = int(e.Root.ComputedWidth)
	e.Root.Rect.Height = int(e.Root.ComputedHeight)
	e.Root.ContentRect = geometry.Rect{
		Width:  e.Root.Rect.Width,
		Height: e.Root.Rect.Height,
	}

	// Layout children
	e.layoutBlock(e.Root)
	e.calculatePositions(e.Root, 0, 0)

	return e.Root
}

// buildBoxTree creates a tree of boxes from the DOM tree.
func (e *LayoutEngine) buildBoxTree(node *dom.Node, resolvers map[*dom.Node]style.ComputedStyle, parent *Box) *Box {
	s := resolvers[node]
	if s.Display == style.DisplayNone {
		return nil
	}

	if node.Type == dom.NodeText {
		text := strings.TrimSpace(node.Data)
		if text == "" {
			return nil
		}
		box := &Box{
			Type:      BoxText,
			Node:      node,
			Style:     s,
			Parent:    parent,
		}
		return box
	}

	if node.Type != dom.NodeElement && node.Type != dom.NodeDocument {
		return nil
	}

	box := &Box{
		Type:   BoxBlock,
		Node:   node,
		Style:  s,
		Parent: parent,
	}

	if node.Type == dom.NodeDocument {
		box.Type = BoxRoot
	}

	if s.Display == style.DisplayInline {
		box.Type = BoxInline
	}

	// Set box model edges from computed style
	box.Margin.Top = resolveLength(s.MarginTop, e.ViewWidth)
	box.Margin.Right = resolveLength(s.MarginRight, e.ViewWidth)
	box.Margin.Bottom = resolveLength(s.MarginBottom, e.ViewWidth)
	box.Margin.Left = resolveLength(s.MarginLeft, e.ViewWidth)

	box.Border.Top = resolveLength(s.BorderTop.Width, e.ViewWidth)
	box.Border.Right = resolveLength(s.BorderRight.Width, e.ViewWidth)
	box.Border.Bottom = resolveLength(s.BorderBottom.Width, e.ViewWidth)
	box.Border.Left = resolveLength(s.BorderLeft.Width, e.ViewWidth)

	box.Padding.Top = resolveLength(s.PaddingTop, e.ViewWidth)
	box.Padding.Right = resolveLength(s.PaddingRight, e.ViewWidth)
	box.Padding.Bottom = resolveLength(s.PaddingBottom, e.ViewWidth)
	box.Padding.Left = resolveLength(s.PaddingLeft, e.ViewWidth)

	// Process children
	for _, child := range node.Children {
		childBox := e.buildBoxTree(child, resolvers, box)
		if childBox != nil {
			box.Children = append(box.Children, childBox)
		}
	}

	// If no children but has text content, might be an anonymous box
	// (handled by text nodes above)

	return box
}

// layoutBlock performs block layout on a box and its children.
func (e *LayoutEngine) layoutBlock(box *Box) {
	if box == nil {
		return
	}

	// Text boxes have their dimensions set by the parent; skip layout.
	if box.Type == BoxText {
		return
	}

	parentWidth := box.Rect.Width
	if box.Parent != nil {
		parentWidth = box.Parent.ContentRect.Width
	}
	if parentWidth <= 0 {
		parentWidth = int(e.ViewWidth)
	}

	// Resolve width
	box.ComputedWidth = resolveLength(box.Style.Width, float64(parentWidth))
	if box.ComputedWidth == 0 || box.Style.Width.Unit == style.LengthPercent {
		// Default: fill parent (auto or 100%)
		box.ComputedWidth = float64(parentWidth) -
			box.Margin.Left - box.Margin.Right -
			box.Border.Left - box.Border.Right -
			box.Padding.Left - box.Padding.Right
	} else if box.ComputedWidth > 0 {
		// Explicit width (px, em, etc): clamp to available space
		available := float64(parentWidth) -
			box.Margin.Left - box.Margin.Right -
			box.Border.Left - box.Border.Right -
			box.Padding.Left - box.Padding.Right
		if box.ComputedWidth > available {
			box.ComputedWidth = available
		}
	}
	if box.ComputedWidth < 0 {
		box.ComputedWidth = 0
	}

	// Resolve height
	box.ComputedHeight = resolveLength(box.Style.Height, e.ViewHeight)

	// Update full rect (position will be set by calculatePositions)
	box.Rect.Width = int(box.TotalWidth())
	box.Rect.Height = int(box.TotalHeight())

	// Content dimensions (position set later by calculatePositions)
	box.ContentRect.Width = int(box.ComputedWidth)
	box.ContentRect.Height = int(box.ComputedHeight)

	// Layout children (their positions are relative to this box's Rect)
	var cursorY float64
	for _, child := range box.Children {
		// Inherit parent width for percentage calculations
		if child.Type == BoxText {
			text := child.Node.Data

			// Calculate wrapped text dimensions
			availableWidth := box.ContentRect.Width

			words := strings.Fields(text)
			if len(words) == 0 {
				continue
			}

			currentLineWidth := 0.0
			maxLineWidth := 0.0
			totalHeight := 1.0

			for _, word := range words {
				wordWidth := float64(len(word))
				if currentLineWidth+wordWidth > float64(availableWidth) && currentLineWidth > 0 {
					totalHeight += 1.0
					currentLineWidth = wordWidth
				} else {
					currentLineWidth += wordWidth + 1.0
				}
				if currentLineWidth > maxLineWidth {
					maxLineWidth = currentLineWidth
				}
			}

			child.ComputedWidth = maxLineWidth
			child.ComputedHeight = totalHeight
		} else {
			child.ComputedWidth = resolveLength(child.Style.Width, box.ComputedWidth)
			if child.ComputedWidth == 0 {
				child.ComputedWidth = box.ComputedWidth -
					child.Margin.Left - child.Margin.Right -
					child.Border.Left - child.Border.Right -
					child.Padding.Left - child.Padding.Right
				if child.ComputedWidth < 0 {
					child.ComputedWidth = 0
				}
			}

			child.ComputedHeight = resolveLength(child.Style.Height, e.ViewHeight)
		}

		// Position relative to this box's Rect (top-left of border edge)
		contentOffsetX := box.Margin.Left + box.Border.Left + box.Padding.Left
		contentOffsetY := box.Margin.Top + box.Border.Top + box.Padding.Top
		child.Rect.X = int(contentOffsetX + child.Margin.Left)
		child.Rect.Y = int(contentOffsetY + cursorY + child.Margin.Top)
		child.Rect.Width = int(child.ComputedWidth + child.Margin.Left + child.Margin.Right +
			child.Border.Left + child.Border.Right +
			child.Padding.Left + child.Padding.Right)
		child.Rect.Height = int(child.ComputedHeight + child.Margin.Top + child.Margin.Bottom +
			child.Border.Top + child.Border.Bottom +
			child.Padding.Top + child.Padding.Bottom)

		// Layout children of this child (recursive)
		e.layoutBlock(child)

		cursorY += float64(child.Rect.Height)
	}

	// Update own height based on children if not explicitly set
	if box.ComputedHeight == 0 {
		box.ComputedHeight = cursorY +
			box.Padding.Top + box.Padding.Bottom +
			box.Border.Top + box.Border.Bottom
		box.Rect.Height = int(box.TotalHeight())
		box.ContentRect.Height = int(box.ComputedHeight -
			box.Padding.Top - box.Padding.Bottom -
			box.Border.Top - box.Border.Bottom)
	}
}

// calculatePositions sets the absolute position for a box and its children.
// box.Rect is relative to parent's Rect; after this call it becomes absolute.
func (e *LayoutEngine) calculatePositions(box *Box, parentX, parentY float64) {
	// Make this box's position absolute
	box.Rect.X += int(parentX)
	box.Rect.Y += int(parentY)

	// Content area is now absolute too
	box.ContentRect.X = box.Rect.X + int(box.Margin.Left+box.Border.Left+box.Padding.Left)
	box.ContentRect.Y = box.Rect.Y + int(box.Margin.Top+box.Border.Top+box.Padding.Top)
	box.ContentRect.Width = int(box.ComputedWidth)
	box.ContentRect.Height = int(box.ComputedHeight)

	// Recurse into children with this box's absolute position
	for _, child := range box.Children {
		e.calculatePositions(child, float64(box.Rect.X), float64(box.Rect.Y))
	}
}

// resolveLength resolves a length value to a concrete pixel value.
func resolveLength(l style.Length, parent float64) float64 {
	switch l.Unit {
	case style.LengthPx:
		return l.Value
	case style.LengthEm, style.LengthRem:
		return l.Value * 16 // Base font size
	case style.LengthPercent:
		return l.Value * parent / 100.0
	case style.LengthAuto:
		return 0
	case style.LengthNone:
		return 0
	}
	return 0
}

// FindBoxAtPoint returns the deepest box at the given point.
func (b *Box) FindBoxAtPoint(x, y int) *Box {
	if !b.Rect.Contains(x, y) {
		return nil
	}

	// Check children in reverse order (last child is visually on top)
	for i := len(b.Children) - 1; i >= 0; i-- {
		if child := b.Children[i].FindBoxAtPoint(x, y); child != nil {
			return child
		}
	}

	return b
}
