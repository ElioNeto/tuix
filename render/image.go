// Package render provides terminal output rendering with image protocol support.
package render

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io"
)

// ImageProtocol represents the terminal image protocol to use.
type ImageProtocol int

const (
	ImageNone  ImageProtocol = iota
	ImageKitty               // Kitty terminal graphics protocol
	ImageSixel               // Sixel graphics format
)

// KittyPlacement controls how the image is positioned.
type KittyPlacement int

const (
	KittyPlaceCursor KittyPlacement = iota // Place image at cursor position
	KittyPlaceAbsolute                     // Place at absolute coordinates
	KittyPlaceAnchored                     // Place and keep position
)

// EncodeKitty encodes an image as Kitty protocol data and writes it to w.
// The image is encoded as PNG and transmitted via Kitty's graphics protocol.
func EncodeKitty(w io.Writer, img image.Image, width, height int) error {
	// Encode image as PNG in memory
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return fmt.Errorf("png encode: %w", err)
	}

	// Kitty protocol: \x1b_G<params>;<base64>\x1b\\
	// Parameters:
	//   a=T - transmit
	//   f=100 - PNG format
	//   s=<width> - image width
	//   v=<height> - image height
	//   m=1 - more chunks follow (for large images), m=0 - last chunk

	data := buf.Bytes()
	b64 := base64.StdEncoding.EncodeToString(data)

	// Max chunk size to avoid buffer issues
	chunkSize := 4096

	for i := 0; i < len(b64); i += chunkSize {
		end := i + chunkSize
		if end > len(b64) {
			end = len(b64)
		}
		chunk := b64[i:end]

		more := 0
		if end < len(b64) {
			more = 1
		}

		params := fmt.Sprintf("a=T,f=100,s=%d,v=%d,m=%d", width, height, more)
		if _, err := fmt.Fprintf(w, "\x1b_G%s;%s\x1b\\", params, chunk); err != nil {
			return err
		}
	}

	return nil
}

// EncodeKittyPlacement outputs a Kitty protocol image placement command.
// This can be used to position an image at specific coordinates.
func EncodeKittyPlacement(w io.Writer, img image.Image, x, y, width, height int) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return fmt.Errorf("png encode: %w", err)
	}

	data := buf.Bytes()
	b64 := base64.StdEncoding.EncodeToString(data)

	params := fmt.Sprintf("a=T,f=100,s=%d,v=%d,c=%d,r=%d", width, height, x, y)
	chunkSize := 4096

	for i := 0; i < len(b64); i += chunkSize {
		end := i + chunkSize
		if end > len(b64) {
			end = len(b64)
		}
		chunk := b64[i:end]

		more := 0
		if end < len(b64) {
			more = 1
		}

		if _, err := fmt.Fprintf(w, "\x1b_G%s,m=%d;%s\x1b\\", params, more, chunk); err != nil {
			return err
		}
	}

	return nil
}

// EncodeSixel encodes an image as Sixel data and writes it to w.
// This is a basic implementation that converts the image to indexed colors.
func EncodeSixel(w io.Writer, img image.Image) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Sixel introduction sequence
	fmt.Fprint(w, "\x1bPq")

	// Set color registers (up to 256 colors)
	// For simplicity, use a basic color palette
	colors := extractColors(img, 16)
	fmt.Fprintf(w, "#0;2;0;0;0") // Black background
	for i, c := range colors {
		r, g, b := c.R, c.G, c.B
		fmt.Fprintf(w, "#%d;2;%d;%d;%d", i+1, r*100/255, g*100/255, b*100/255)
	}

	// Render image line by line using Sixel data
	// Each line of sixels represents 6 vertical pixels
	for y := 0; y < height; y += 6 {
		for x := 0; x < width; x++ {
			// Determine color for each 6-pixel vertical strip
			sixelBits := make([]byte, 6)
			for dy := 0; dy < 6 && y+dy < height; dy++ {
				r, g, b, _ := img.At(x, y+dy).RGBA()
				idx := nearestColor(colors, uint8(r>>8), uint8(g>>8), uint8(b>>8))
				_ = idx
				// For binary sixel: set bit if pixel is bright enough
				lum := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
				if lum > 128 {
					sixelBits[dy] = 1
				}
			}
			// Pack bits into sixel byte (0x3F + bits)
			var sixel byte = 0x3F // '?'
			for dy := 0; dy < 6; dy++ {
				if sixelBits[dy] != 0 {
					sixel |= 1 << uint(dy)
				}
			}
			fmt.Fprintf(w, "#0%c", sixel)
		}
		fmt.Fprint(w, "$") // Carriage return for this line
		// Move to next line
		if y+6 < height {
			fmt.Fprint(w, "-")
		}
	}

	// Sixel termination
	fmt.Fprint(w, "\x1b\\")
	return nil
}

// rgbColor holds an RGB color value.
type rgbColor struct {
	R, G, B uint8
}

// extractColors extracts up to n representative colors from the image.
func extractColors(img image.Image, n int) []rgbColor {
	bounds := img.Bounds()
	if bounds.Dx() == 0 || bounds.Dy() == 0 {
		return []rgbColor{{0, 0, 0}}
	}

	// Use a simple grid sampling to extract colors
	colors := make([]rgbColor, 0, n)
	seen := make(map[rgbColor]bool)

	stepX := max(1, bounds.Dx()/8)
	stepY := max(1, bounds.Dy()/8)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += stepY {
		for x := bounds.Min.X; x < bounds.Max.X; x += stepX {
			if len(colors) >= n {
				break
			}
			r, g, b, _ := img.At(x, y).RGBA()
			c := rgbColor{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8)}
			if !seen[c] {
				seen[c] = true
				colors = append(colors, c)
			}
		}
	}

	// Pad with black if we didn't get enough colors
	for len(colors) < n {
		colors = append(colors, rgbColor{0, 0, 0})
	}

	return colors
}

// nearestColor finds the closest color in the palette to the target.
func nearestColor(palette []rgbColor, r, g, b uint8) int {
	bestIdx := 0
	bestDist := int(^uint(0) >> 1) // max int

	for i, c := range palette {
		dr := int(r) - int(c.R)
		dg := int(g) - int(c.G)
		db := int(b) - int(c.B)
		dist := dr*dr + dg*dg + db*db
		if dist < bestDist {
			bestDist = dist
			bestIdx = i
		}
	}

	return bestIdx
}
