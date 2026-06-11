package geometry

// Size represents a 2D dimension.
type Size struct {
	Width, Height int
}

// Zero checks whether the size has zero width and height.
func (s Size) IsZero() bool { return s.Width == 0 && s.Height == 0 }
