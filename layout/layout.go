// Package layout implements the box model and layout engine.
//
// It takes styled DOM nodes and produces a tree of layout boxes with
// computed positions and sizes. The layout engine supports block and
// inline formatting contexts.
package layout

import (
	"sort"
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
	BoxFlex
)

// Box represents a rectangular layout box in the formatting tree.
//
// Rect represents the BORDER BOX (content + padding + border).
// Margins are outside the border box — they add spacing between boxes
// but are not included in Rect.Width/Height.
type Box struct {
	Type      BoxType
	Node      *dom.Node
	Style     style.ComputedStyle
	Rect      geometry.Rect
	Children  []*Box
	Parent    *Box

	// Content area (inside padding and border)
	ContentRect geometry.Rect

	// Box model edges (outside the border box / Rect)
	Margin    geometry.Edges
	Border    geometry.Edges
	Padding   geometry.Edges

	// Computed dimensions (content area only, not including padding/border/margin)
	ComputedWidth  float64
	ComputedHeight float64

	// Scroll offset (0,0 when not scrolled)
	ScrollX int
	ScrollY int

	// Total scrollable content dimensions (used for scrollbar max)
	ScrollWidth  int
	ScrollHeight int

	// ZIndex for stacking order (0 = auto/DOM order)
	ZIndex int
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

// BorderWidth returns just the border+padding+content (no margins).
func (b *Box) BorderWidth() float64 {
	return b.Border.Left + b.Padding.Left + b.ComputedWidth +
		b.Padding.Right + b.Border.Right
}

// BorderHeight returns just the border+padding+content (no margins).
func (b *Box) BorderHeight() float64 {
	return b.Border.Top + b.Padding.Top + b.ComputedHeight +
		b.Padding.Bottom + b.Border.Bottom
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

	// Layout children using the correct algorithm per box type
	e.layoutBox(e.Root)
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
			ZIndex:    s.ZIndex,
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
		ZIndex: s.ZIndex,
		Parent: parent,
	}

	if node.Type == dom.NodeDocument {
		box.Type = BoxRoot
	}

	if s.Display == style.DisplayInline {
		box.Type = BoxInline
	}

	if s.Display == style.DisplayFlex {
		box.Type = BoxFlex
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

// layoutBox dispatches layout to the correct algorithm based on box type.
func (e *LayoutEngine) layoutBox(box *Box) {
	if box == nil {
		return
	}

	// Text boxes and inline boxes are sized by their parent; skip direct layout.
	if box.Type == BoxText || box.Type == BoxInline {
		return
	}

	switch box.Type {
	case BoxFlex:
		e.layoutFlex(box)
	default:
		e.layoutBlock(box)
	}
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
	// Rect.Width/Height is the BORDER BOX (content + padding + border), no margins.
	box.Rect.Width = int(box.BorderWidth())
	box.Rect.Height = int(box.BorderHeight())

	// Content dimensions (position set later by calculatePositions)
	box.ContentRect.Width = int(box.ComputedWidth)
	box.ContentRect.Height = int(box.ComputedHeight)

	// Layout children (their positions are relative to this box's Rect)
	var cursorY float64

	// First, check if this box has inline-level children. Inline children
	// (BoxInline, BoxText) participate in an inline formatting context and
	// flow into line boxes rather than stacking vertically.
	hasInline := false
	for _, child := range box.Children {
		if child.Type == BoxInline || child.Type == BoxText {
			hasInline = true
			break
		}
	}

	if hasInline {
		// Use inline formatting context for this block's children.
		// Extract inline runs (consecutive inline/text children) and
		// lay them out as line boxes.
		var i int
		for i < len(box.Children) {
			// Skip block children — they get stacked vertically as usual.
			if box.Children[i].Type != BoxInline && box.Children[i].Type != BoxText {
				// Block child: lay out normally
				child := box.Children[i]
				e.layoutBlockChild(child, box, &cursorY)
				i++
				continue
			}

			// Collect consecutive inline children
			inlineStart := i
			for i < len(box.Children) &&
				(box.Children[i].Type == BoxInline || box.Children[i].Type == BoxText) {
				i++
			}
			inlineGroup := box.Children[inlineStart:i]

			// Lay out this inline group
			e.layoutInline(box, inlineGroup, &cursorY)
		}
	} else {
		// Pure block layout: stack children vertically
		for _, child := range box.Children {
			e.layoutBlockChild(child, box, &cursorY)
		}
	}

	// Update own height based on children if not explicitly set
	if box.ComputedHeight == 0 {
		box.ComputedHeight = cursorY +
			box.Padding.Top + box.Padding.Bottom +
			box.Border.Top + box.Border.Bottom
		box.Rect.Height = int(box.BorderHeight())
		box.ContentRect.Height = int(box.ComputedHeight -
			box.Padding.Top - box.Padding.Bottom -
			box.Border.Top - box.Border.Bottom)
	}

	// Compute scrollable content extents
	box.ScrollWidth = box.ContentRect.Width
	box.ScrollHeight = int(cursorY) // total height of children in the content area
	if box.ScrollHeight < box.ContentRect.Height {
		box.ScrollHeight = box.ContentRect.Height
	}
	if box.ScrollWidth < box.ContentRect.Width {
		box.ScrollWidth = box.ContentRect.Width
	}
}

// layoutBlockChild lays out a single block-level child within its parent.
func (e *LayoutEngine) layoutBlockChild(child *Box, parent *Box, cursorY *float64) {
	if child.Type == BoxText {
		text := child.Node.Data

		// Calculate wrapped text dimensions
		availableWidth := parent.ContentRect.Width

		words := strings.Fields(text)
		if len(words) == 0 {
			return
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
		child.ComputedWidth = resolveLength(child.Style.Width, parent.ComputedWidth)
		if child.ComputedWidth == 0 {
			child.ComputedWidth = parent.ComputedWidth -
				child.Margin.Left - child.Margin.Right -
				child.Border.Left - child.Border.Right -
				child.Padding.Left - child.Padding.Right
			if child.ComputedWidth < 0 {
				child.ComputedWidth = 0
			}
		}

		child.ComputedHeight = resolveLength(child.Style.Height, e.ViewHeight)
	}

	// Position relative to this box's Rect (border box)
	// Margins are outside the border box — margin.Left shifts X, margin.Top shifts Y.
	contentOffsetX := parent.Border.Left + parent.Padding.Left
	contentOffsetY := parent.Border.Top + parent.Padding.Top
	child.Rect.X = int(contentOffsetX + child.Margin.Left)
	child.Rect.Y = int(contentOffsetY + *cursorY + child.Margin.Top)
	// Rect.Width/Height is the border box (no margins)
	child.Rect.Width = int(child.ComputedWidth + child.Border.Left + child.Border.Right +
		child.Padding.Left + child.Padding.Right)
	child.Rect.Height = int(child.ComputedHeight + child.Border.Top + child.Border.Bottom +
		child.Padding.Top + child.Padding.Bottom)

	// Layout children of this child (recursive)
	e.layoutBox(child)

	*cursorY += float64(child.Rect.Height)
}

// layoutInline performs inline formatting on a group of inline-level children.
// It generates line boxes and positions each child within them.
func (e *LayoutEngine) layoutInline(container *Box, inlineChildren []*Box, cursorY *float64) {
	availableWidth := container.ComputedWidth
	contentX := container.Border.Left + container.Padding.Left

	// Phase 1: Split inline children into "runs" that fit on lines.
	type run struct {
		box       *Box
		text      string // for text boxes
		wordWidth float64
	}
	var currentLineRuns []*run
	currentLineWidth := 0.0

	// lineRuns stores all runs grouped by line
	var lines [][]*run
	lineHeights := []float64{}

	flushLine := func() {
		if len(currentLineRuns) > 0 {
			lines = append(lines, currentLineRuns)
			// Calculate line height: max of all item heights
			maxH := 1.0
			for _, r := range currentLineRuns {
				h := float64(r.box.Rect.Height)
				if r.box.ComputedHeight > h {
					h = r.box.ComputedHeight
				}
				if h < 1 {
					h = 1
				}
				if h > maxH {
					maxH = h
				}
			}
			lineHeights = append(lineHeights, maxH)
			currentLineRuns = nil
			currentLineWidth = 0
		}
	}

		for _, child := range inlineChildren {
		if child.Type == BoxInline {
			// Add space before this inline element if not at start of line
			if currentLineWidth > 0 {
				currentLineRuns = append(currentLineRuns, &run{
					box:       child,
					text:      " ",
					wordWidth: 1,
				})
				currentLineWidth++
			}

			// Inline element: determine its content size
			child.ComputedWidth = resolveLength(child.Style.Width, availableWidth)
			if child.ComputedWidth == 0 {
				// For inline elements without explicit width, compute from content
				child.ComputedWidth = textContentWidth(child)
			}
			child.ComputedHeight = resolveLength(child.Style.Height, e.ViewHeight)
			if child.ComputedHeight == 0 {
				child.ComputedHeight = 1
			}

			// Border-box width
			bbw := child.ComputedWidth + child.Border.Left + child.Border.Right +
				child.Padding.Left + child.Padding.Right

			if currentLineWidth+bbw > availableWidth && currentLineWidth > 0 {
				flushLine()
			}

			currentLineRuns = append(currentLineRuns, &run{
				box:       child,
				wordWidth: bbw,
			})
			currentLineWidth += bbw
		} else if child.Type == BoxText {
			text := child.Node.Data
			if strings.TrimSpace(text) == "" {
				continue
			}
			words := strings.Fields(text)

			for _, word := range words {
				wordWidth := float64(len(word))
				needsSpace := currentLineWidth > 0
				totalWordWidth := wordWidth
				if needsSpace {
					totalWordWidth++ // space before word
				}

				if currentLineWidth+totalWordWidth > availableWidth && currentLineWidth > 0 {
					flushLine()
					needsSpace = false
				}

				if needsSpace {
					// Add a space run
					currentLineRuns = append(currentLineRuns, &run{
						box:       child,
						text:      " ",
						wordWidth: 1,
					})
					currentLineWidth++
				}

				currentLineRuns = append(currentLineRuns, &run{
					box:       child,
					text:      word,
					wordWidth: wordWidth,
				})
				currentLineWidth += wordWidth
			}
		}
	}
	flushLine()

	// Phase 2: Position items within each line
	for li, line := range lines {
		lineHeight := lineHeights[li]
		cursorX := 0.0
		positioned := make(map[*Box]bool) // track boxes already positioned

		for _, r := range line {
			box := r.box

			// Only position if not already positioned (multiple runs may share same box)
			if !positioned[box] {
				positioned[box] = true
				box.Rect.X = int(contentX + cursorX)
				box.Rect.Y = int(*cursorY)
			}

			if r.text != "" && box.Type == BoxText {
				// For text runs, set the box to represent this word
				box.ComputedWidth = r.wordWidth
				box.ComputedHeight = lineHeight
			} else if box.Type == BoxInline {
				box.ComputedHeight = lineHeight
				// Propagate position and size to child text boxes
				for _, child := range box.Children {
					if child.Type == BoxText {
						child.Rect.X = box.Rect.X + int(box.Border.Left+box.Padding.Left)
						child.Rect.Y = box.Rect.Y + int(box.Border.Top+box.Padding.Top)
						child.Rect.Width = int(r.wordWidth)
						child.Rect.Height = int(lineHeight)
						child.ContentRect.X = child.Rect.X
						child.ContentRect.Y = child.Rect.Y
						child.ContentRect.Width = child.Rect.Width
						child.ContentRect.Height = child.Rect.Height
					}
				}
			}

			// Border box dimensions
			box.Rect.Width = int(r.wordWidth)
			box.Rect.Height = int(lineHeight)

			// Content rect
			box.ContentRect.X = box.Rect.X + int(box.Border.Left+box.Padding.Left)
			box.ContentRect.Y = box.Rect.Y + int(box.Border.Top+box.Padding.Top)
			box.ContentRect.Width = int(box.ComputedWidth)
			box.ContentRect.Height = int(box.ComputedHeight)

			cursorX += r.wordWidth
		}

		*cursorY += lineHeight
	}
}

// layoutFlex performs flex layout on a flex container and its children.
// It implements a simplified flexbox algorithm that handles:
// - flex-direction, flex-wrap, justify-content
// - align-items, align-self (stretch, flex-start, flex-end, center)
// - flex-grow, flex-shrink, flex-basis
// - order, gap
func (e *LayoutEngine) layoutFlex(container *Box) {
	if container == nil {
		return
	}

	parentWidth := float64(container.Rect.Width)
	if container.Parent != nil {
		parentWidth = float64(container.Parent.ContentRect.Width)
	}
	if parentWidth <= 0 {
		parentWidth = e.ViewWidth
	}

	// Resolve container width
	container.ComputedWidth = resolveLength(container.Style.Width, parentWidth)
	if container.ComputedWidth == 0 || container.Style.Width.Unit == style.LengthPercent {
		container.ComputedWidth = parentWidth -
			container.Margin.Left - container.Margin.Right -
			container.Border.Left - container.Border.Right -
			container.Padding.Left - container.Padding.Right
	}
	if container.ComputedWidth < 0 {
		container.ComputedWidth = 0
	}

	// Resolve container height (may be 0 = auto, adjusted later)
	container.ComputedHeight = resolveLength(container.Style.Height, e.ViewHeight)

	container.Rect.Width = int(container.BorderWidth())
	container.Rect.Height = int(container.BorderHeight())
	container.ContentRect.Width = int(container.ComputedWidth)
	container.ContentRect.Height = int(container.ComputedHeight)

	// Determine main and cross axis
	mainIsX := container.Style.FlexDirection == style.FlexDirectionRow ||
		container.Style.FlexDirection == style.FlexDirectionRowReverse
	reverse := container.Style.FlexDirection == style.FlexDirectionRowReverse ||
		container.Style.FlexDirection == style.FlexDirectionColumnReverse

	// Content area dimensions (inside padding+border)
	contentX := container.Border.Left + container.Padding.Left
	contentY := container.Border.Top + container.Padding.Top
	contentWidth := container.ComputedWidth
	contentHeight := container.ComputedHeight

	// Resolve gaps
	gapMain := resolveLength(container.Style.ColumnGap, contentWidth)
	gapCross := resolveLength(container.Style.RowGap, contentHeight)
	if !mainIsX {
		gapMain = resolveLength(container.Style.RowGap, contentHeight)
		gapCross = resolveLength(container.Style.ColumnGap, contentWidth)
	}

	// Collect and sort flex items
	items := make([]*Box, 0, len(container.Children))
	for _, child := range container.Children {
		if child.Style.Display == style.DisplayNone {
			continue
		}
		items = append(items, child)
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Style.Order < items[j].Style.Order
	})

	// Helper: get/set main size and position for a box
	getMainSize := func(b *Box) float64 {
		if mainIsX {
			return b.ComputedWidth
		}
		return b.ComputedHeight
	}
	setMainSize := func(b *Box, v float64) {
		if mainIsX {
			b.ComputedWidth = v
		} else {
			b.ComputedHeight = v
		}
	}
	getCrossSize := func(b *Box) float64 {
		if !mainIsX {
			return b.ComputedWidth
		}
		return b.ComputedHeight
	}
	setCrossSize := func(b *Box, v float64) {
		if !mainIsX {
			b.ComputedWidth = v
		} else {
			b.ComputedHeight = v
		}
	}
	getMainMarginStart := func(b *Box) float64 {
		if mainIsX {
			return b.Margin.Left
		}
		return b.Margin.Top
	}
	getMainMarginEnd := func(b *Box) float64 {
		if mainIsX {
			return b.Margin.Right
		}
		return b.Margin.Bottom
	}
	getCrossMarginStart := func(b *Box) float64 {
		if !mainIsX {
			return b.Margin.Left
		}
		return b.Margin.Top
	}
	getCrossMarginEnd := func(b *Box) float64 {
		if !mainIsX {
			return b.Margin.Right
		}
		return b.Margin.Bottom
	}

	// Determine the container's main size (available for flex items)
	containerMainSize := contentWidth
	if !mainIsX {
		containerMainSize = contentHeight
	}

	// Phase 1: Determine flex bases and hypothetical main sizes
	hypoMainSizes := make([]float64, len(items))
	for i, item := range items {
		// Resolve width and height if explicitly set
		item.ComputedWidth = resolveLength(item.Style.Width, contentWidth)
		item.ComputedHeight = resolveLength(item.Style.Height, contentHeight)

		// Determine flex basis — the initial main size before growing/shrinking
		var flexBasis float64
		if item.Style.FlexBasis.Unit != style.LengthAuto {
			flexBasis = resolveLength(item.Style.FlexBasis, containerMainSize)
		} else {
			// When flex-basis is auto, use the item's main-size property (width or height)
			flexBasis = getMainSize(item)
		}

		// Apply min/max constraints to flex basis
		// TODO: handle min-width/min-height, max-width/max-height

		hypoMainSizes[i] = flexBasis
	}

	// Calculate total hypothetical main size (including gaps)
	totalHypoMain := 0.0
	for i, item := range items {
		// Include margins in the hypothetical size
		totalHypoMain += hypoMainSizes[i] + getMainMarginStart(item) + getMainMarginEnd(item)
		if i > 0 {
			totalHypoMain += gapMain
		}
	}

	// Phase 2: Distribute free space (flex-grow / flex-shrink)
	freeSpace := containerMainSize - totalHypoMain

	var totalGrow, totalShrink float64
	for _, item := range items {
		totalGrow += item.Style.FlexGrow
		totalShrink += item.Style.FlexShrink
	}

	mainSizes := make([]float64, len(items))
	copy(mainSizes, hypoMainSizes)

	if freeSpace > 0 && totalGrow > 0 {
		// Distribute positive free space via flex-grow
		extraPerUnit := freeSpace / totalGrow
		for i, item := range items {
			mainSizes[i] += extraPerUnit * item.Style.FlexGrow
			if mainSizes[i] < 0 {
				mainSizes[i] = 0
			}
		}
	} else if freeSpace < 0 && totalShrink > 0 {
		// Distribute negative free space via flex-shrink
		// Shrink proportionally to flex-shrink * flex-basis
		var scaledShrinkSum float64
		for i, item := range items {
			if hypoMainSizes[i] > 0 {
				scaledShrinkSum += item.Style.FlexShrink * hypoMainSizes[i]
			}
		}
		if scaledShrinkSum > 0 {
			shrinkNeeded := -freeSpace
			for i, item := range items {
				if hypoMainSizes[i] > 0 {
					shrink := shrinkNeeded * (item.Style.FlexShrink * hypoMainSizes[i]) / scaledShrinkSum
					mainSizes[i] -= shrink
					if mainSizes[i] < 0 {
						mainSizes[i] = 0
					}
				}
			}
		}
	}

	// Apply final main sizes to items
	for i, item := range items {
		setMainSize(item, mainSizes[i])
	}

	// Phase 3: Calculate border-box dimensions for each item
	// (needed for positioning)
	type itemDim struct {
		mainSize       float64 // content main size
		crossSize      float64 // content cross size
		borderMainSize float64 // border-box main size
	}
	itemDims := make([]itemDim, len(items))
	for i, item := range items {
		// Clamp to min/max
		if mainSizes[i] < 0 {
			mainSizes[i] = 0
		}

		// For the cross size, initial estimate is 0 or explicitly set
		crossSize := getCrossSize(item)
		if crossSize == 0 {
			// If no explicit cross size, we'll set it based on align-items: stretch
			// or leave as 0 and use the container's cross size later
		}

		bl := item.Border.Left
		br := item.Border.Right
		bt := item.Border.Top
		bb := item.Border.Bottom
		pl := item.Padding.Left
		pr := item.Padding.Right
		pt := item.Padding.Top
		pb := item.Padding.Bottom

		var borderMainSize float64
		if mainIsX {
			borderMainSize = mainSizes[i] + bl + br + pl + pr
			crossSize = item.ComputedHeight // was set earlier
		} else {
			borderMainSize = mainSizes[i] + bt + bb + pt + pb
			crossSize = item.ComputedWidth
		}
		itemDims[i] = itemDim{
			mainSize:       mainSizes[i],
			crossSize:      crossSize,
			borderMainSize: borderMainSize,
		}
	}

	// Phase 4: If wrapping, break items into lines
	type flexLine struct {
		items    []int // indices into items
		mainSize float64
	}
	lines := []flexLine{{}}

	if container.Style.FlexWrap == style.FlexWrapWrap || container.Style.FlexWrap == style.FlexWrapWrapReverse {
		currentLine := 0
		currentLineMain := 0.0
		for i := range items {
			itemMain := itemDims[i].borderMainSize + getMainMarginStart(items[i]) + getMainMarginEnd(items[i])
			if currentLineMain+itemMain > containerMainSize && len(lines[currentLine].items) > 0 {
				lines = append(lines, flexLine{})
				currentLine++
				currentLineMain = 0
			}
			lines[currentLine].items = append(lines[currentLine].items, i)
			lines[currentLine].mainSize += itemMain
		}
	} else {
		// No wrap — all items on one line
		for i := range items {
			lines[0].items = append(lines[0].items, i)
		}
	}

	// Phase 5: Position items on main axis per line
	lineOffset := 0.0
	wrapReverse := container.Style.FlexWrap == style.FlexWrapWrapReverse

	for _, line := range lines {
		// Calculate the total main size of items on this line (including gaps)
		lineMainUsed := 0.0
		for idx, i := range line.items {
			item := items[i]
			itemMain := itemDims[i].borderMainSize + getMainMarginStart(item) + getMainMarginEnd(item)
			lineMainUsed += itemMain
			if idx > 0 {
				lineMainUsed += gapMain
			}
		}

		// Available main space within the container for this line
		availableMain := containerMainSize

		// Position items along main axis based on justify-content
		var startOffset float64
		var spacing float64 // extra space between items
		remainingSpace := availableMain - lineMainUsed

		switch container.Style.JustifyContent {
		case style.JustifyContentFlexStart, style.JustifyContentFlexEnd:
			if (container.Style.JustifyContent == style.JustifyContentFlexEnd) != reverse {
				startOffset = availableMain - lineMainUsed
			}
		case style.JustifyContentCenter:
			startOffset = remainingSpace / 2
		case style.JustifyContentSpaceBetween:
			if len(line.items) > 1 {
				spacing = remainingSpace / float64(len(line.items)-1)
			} else {
				startOffset = 0
			}
		case style.JustifyContentSpaceAround:
			if len(line.items) > 0 {
				spacing = remainingSpace / float64(len(line.items))
				startOffset = spacing / 2
			}
		case style.JustifyContentSpaceEvenly:
			if len(line.items) > 0 {
				spacing = remainingSpace / float64(len(line.items)+1)
				startOffset = spacing
			}
		}

		cursor := startOffset
		for _, idx := range line.items {
			item := items[idx]

			// Position this item along main axis
			var mainPos float64
			if mainIsX {
				mainPos = contentX + cursor + getMainMarginStart(item)
			} else {
				mainPos = contentY + cursor + getMainMarginStart(item)
			}

			// Calculate cross size for this item
			// For stretch: use line's cross size
			// For other alignments: use item's content cross size
			var itemCrossSize float64
			alignSelf := resolveAlignSelf(item.Style.AlignSelf, container.Style.AlignItems)

			// Determine the available cross space for this line
			var availableCrossForLine float64
			if mainIsX {
				availableCrossForLine = contentHeight - lineOffset
			} else {
				availableCrossForLine = contentWidth - lineOffset
			}
			if availableCrossForLine < 0 {
				availableCrossForLine = 0
			}

			if alignSelf == style.AlignItemsStretch && availableCrossForLine > 0 {
				// Stretch to fill available cross space
				itemCrossSize = availableCrossForLine
			} else {
				itemCrossSize = itemDims[idx].crossSize
				if itemCrossSize <= 0 {
					// If cross size not determined, use 1 row as default
					itemCrossSize = 1
				}
			}

			// Apply the cross size to the item
			setCrossSize(item, itemCrossSize)

			// Set item position and size (border-box)
			var crossStart float64
			if mainIsX {
				// main axis = X, cross axis = Y
				lineCrossStart := contentY + lineOffset
				availableCross := availableCrossForLine

				switch alignSelf {
				case style.AlignItemsFlexStart, style.AlignItemsStretch:
					crossStart = lineCrossStart
				case style.AlignItemsFlexEnd:
					crossStart = lineCrossStart + availableCross - itemCrossSize
				case style.AlignItemsCenter:
					crossStart = lineCrossStart + (availableCross-itemCrossSize)/2
				default:
					crossStart = lineCrossStart
				}

				item.Rect.X = int(mainPos)
				item.Rect.Y = int(crossStart + getCrossMarginStart(item))
			} else {
				// main axis = Y, cross axis = X
				lineCrossStart := contentX + lineOffset
				availableCross := availableCrossForLine

				switch alignSelf {
				case style.AlignItemsFlexStart, style.AlignItemsStretch:
					crossStart = lineCrossStart
				case style.AlignItemsFlexEnd:
					crossStart = lineCrossStart + availableCross - itemCrossSize
				case style.AlignItemsCenter:
					crossStart = lineCrossStart + (availableCross-itemCrossSize)/2
				default:
					crossStart = lineCrossStart
				}

				item.Rect.X = int(crossStart + getCrossMarginStart(item))
				item.Rect.Y = int(mainPos)
			}

			// Set border box dimensions
			bw := item.BorderWidth()
			bh := item.BorderHeight()
			item.Rect.Width = int(bw)
			item.Rect.Height = int(bh)

			// Update content rect
			item.ContentRect.X = item.Rect.X + int(item.Border.Left+item.Padding.Left)
			item.ContentRect.Y = item.Rect.Y + int(item.Border.Top+item.Padding.Top)
			item.ContentRect.Width = int(item.ComputedWidth)
			item.ContentRect.Height = int(item.ComputedHeight)

			// Recurse into children of this flex item
			e.layoutBox(item)

			cursor += itemDims[idx].borderMainSize + gapMain + spacing + getMainMarginStart(item) + getMainMarginEnd(item)
		}

		// Update the container's cross size based on lines
		// Find max cross size in this line
		maxCrossInLine := 0.0
		for _, idx := range line.items {
			item := items[idx]
			cs := getCrossSize(item) + getCrossMarginStart(item) + getCrossMarginEnd(item)
			if cs > maxCrossInLine {
				maxCrossInLine = cs
			}
		}
		lineOffset += maxCrossInLine + gapCross
	}

	// Update container height based on content if auto
	containerHeightFromContent := lineOffset
	if mainIsX {
		// Container height = max of line cross-axis usage + padding + border
		if container.ComputedHeight == 0 {
			containerHeightFromContent += container.Border.Top + container.Border.Bottom +
				container.Padding.Top + container.Padding.Bottom
			container.ComputedHeight = containerHeightFromContent
			container.Rect.Height = int(container.BorderHeight())
			container.ContentRect.Height = int(container.ComputedHeight -
				container.Padding.Top - container.Padding.Bottom -
				container.Border.Top - container.Border.Bottom)
		}
	} else {
		// Container width = max of line cross-axis usage + padding + border
		// (already set above)
	}

	// If wrap-reverse, flip the cross-axis positions
	if wrapReverse && mainIsX {
		for _, item := range items {
			item.Rect.Y = int(contentHeight) - item.Rect.Y - item.Rect.Height
		}
	} else if wrapReverse && !mainIsX {
		for _, item := range items {
			item.Rect.X = int(contentWidth) - item.Rect.X - item.Rect.Width
		}
	}

	// Compute scrollable content extents (relative to content area origin)
	maxX := 0
	maxY := 0
	contentOriginX := container.ContentRect.X
	contentOriginY := container.ContentRect.Y
	for _, item := range items {
		itemEndX := item.Rect.X + item.Rect.Width - contentOriginX
		itemEndY := item.Rect.Y + item.Rect.Height - contentOriginY
		if itemEndX > maxX {
			maxX = itemEndX
		}
		if itemEndY > maxY {
			maxY = itemEndY
		}
	}
	container.ScrollWidth = container.ContentRect.Width
	container.ScrollHeight = container.ContentRect.Height
	if maxX > container.ContentRect.Width {
		container.ScrollWidth = maxX
	}
	if maxY > container.ContentRect.Height {
		container.ScrollHeight = maxY
	}
}

// resolveAlignSelf resolves the effective align-self for an item,
// considering the container's align-items if align-self is auto.
func resolveAlignSelf(as style.AlignSelfType, containerAlign style.AlignItemsType) style.AlignItemsType {
	switch as {
	case style.AlignSelfAuto:
		return containerAlign
	case style.AlignSelfFlexStart:
		return style.AlignItemsFlexStart
	case style.AlignSelfFlexEnd:
		return style.AlignItemsFlexEnd
	case style.AlignSelfCenter:
		return style.AlignItemsCenter
	case style.AlignSelfStretch:
		return style.AlignItemsStretch
	case style.AlignSelfBaseline:
		return style.AlignItemsBaseline
	}
	return containerAlign
}

// calculatePositions sets the absolute position for a box and its children.
// box.Rect is the border box (no margins) relative to parent's border box.
// After this call, positions become absolute.
func (e *LayoutEngine) calculatePositions(box *Box, parentX, parentY float64) {
	// Make this box's position absolute (border box position)
	box.Rect.X += int(parentX)
	box.Rect.Y += int(parentY)

	// Content area is absolute too — starts inside padding+border
	box.ContentRect.X = box.Rect.X + int(box.Border.Left+box.Padding.Left)
	box.ContentRect.Y = box.Rect.Y + int(box.Border.Top+box.Padding.Top)
	box.ContentRect.Width = int(box.ComputedWidth)
	box.ContentRect.Height = int(box.ComputedHeight)

	// Recurse into children with this box's absolute position.
	// Children are positioned relative to this box's border box origin.
	for _, child := range box.Children {
		e.calculatePositions(child, float64(box.Rect.X), float64(box.Rect.Y))
	}
}

// resolveLength resolves a length value to a concrete pixel value.
func resolveLength(l style.Length, parent float64) float64 {
	if l.Unit == style.LengthAuto || l.Value == 0 {
		return 0
	}
	if l.Unit == style.LengthPercent {
		return l.Value / 100.0 * parent
	}
	return l.Value
}

// textContentWidth computes the total character width of text content
// in an inline element's child text nodes.
func textContentWidth(box *Box) float64 {
	var w float64
	for _, child := range box.Children {
		if child.Type == BoxText && child.Node != nil {
			for _, word := range strings.Fields(child.Node.Data) {
				w += float64(len(word))
				if w > 0 {
					w++ // space between words
				}
			}
		}
	}
	if w < 1 {
		w = 1
	}
	return w
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

// FindScrollContainer returns the first scrollable ancestor of a box,
// or the first scrollable container in the tree if no specific box is given.
func (b *Box) FindScrollContainer(from *Box) *Box {
	// If we have a starting box, walk up the tree looking for a scroll container
	if from != nil {
		for cur := from; cur != nil; cur = cur.Parent {
			if cur.ScrollHeight > cur.ContentRect.Height &&
				(cur.Style.OverflowY == style.OverflowScroll || cur.Style.OverflowY == style.OverflowAuto) {
				return cur
			}
			if cur.ScrollWidth > cur.ContentRect.Width &&
				(cur.Style.OverflowX == style.OverflowScroll || cur.Style.OverflowX == style.OverflowAuto) {
				return cur
			}
		}
	}

	// Otherwise, DFS from this box
	var found *Box
	b.findScrollDFS(&found)
	return found
}

func (b *Box) findScrollDFS(found **Box) {
	if *found != nil {
		return
	}
	if b.ScrollHeight > b.ContentRect.Height &&
		(b.Style.OverflowY == style.OverflowScroll || b.Style.OverflowY == style.OverflowAuto) {
		*found = b
		return
	}
	if b.ScrollWidth > b.ContentRect.Width &&
		(b.Style.OverflowX == style.OverflowScroll || b.Style.OverflowX == style.OverflowAuto) {
		*found = b
		return
	}
	for _, child := range b.Children {
		child.findScrollDFS(found)
	}
}
