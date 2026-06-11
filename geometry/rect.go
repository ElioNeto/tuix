// Package geometry provides fundamental geometric types used throughout tuix.
package geometry

// Rect represents a rectangle with position (X, Y) and size (Width, Height).
type Rect struct {
	X, Y          int
	Width, Height int
}

// Contains checks whether the point (px, py) lies within the rectangle.
func (r Rect) Contains(px, py int) bool {
	return px >= r.X && px < r.X+r.Width &&
		py >= r.Y && py < r.Y+r.Height
}

// Intersect returns the overlapping region of two rectangles.
// If they do not overlap, the returned rectangle has zero width or height.
func (r Rect) Intersect(other Rect) Rect {
	x1 := max(r.X, other.X)
	y1 := max(r.Y, other.Y)
	x2 := min(r.X+r.Width, other.X+other.Width)
	y2 := min(r.Y+r.Height, other.Y+other.Height)

	if x1 < x2 && y1 < y2 {
		return Rect{X: x1, Y: y1, Width: x2 - x1, Height: y2 - y1}
	}
	return Rect{}
}

// Offset returns a copy of the rectangle shifted by (dx, dy).
func (r Rect) Offset(dx, dy int) Rect {
	return Rect{X: r.X + dx, Y: r.Y + dy, Width: r.Width, Height: r.Height}
}

// Right returns the X coordinate of the right edge.
func (r Rect) Right() int { return r.X + r.Width }

// Bottom returns the Y coordinate of the bottom edge.
func (r Rect) Bottom() int { return r.Y + r.Height }

// IsEmpty returns true if the rectangle has zero width or height.
func (r Rect) IsEmpty() bool { return r.Width <= 0 || r.Height <= 0 }
