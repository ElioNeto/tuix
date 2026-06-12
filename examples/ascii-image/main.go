// Example: Image-to-ASCII Art Converter
//
// This example demonstrates converting images to ASCII art
// using the ascii package's image conversion feature.
//
// Usage:
//   go run ./examples/ascii-image/                    # uses built-in test pattern
//   go run ./examples/ascii-image/ --file photo.jpg   # converts a specific image
//   go run ./examples/ascii-image/ --file photo.jpg --width 60 --color --dither
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/elioneto/tuix/ascii"
)

func main() {
	filePath := flag.String("file", "", "path to image file (PNG/JPEG/GIF)")
	width := flag.Int("width", 70, "output width in characters")
	color := flag.Bool("color", false, "enable true-color ANSI output")
	dither := flag.Bool("dither", true, "enable Floyd-Steinberg dithering")
	charset := flag.String("charset", ascii.CharsetStandard, "character ramp")
	_ = flag.Bool("test", false, "generate a built-in test pattern image")
	flag.Parse()

	opts := ascii.ImageOptions{
		Width:   *width,
		Height:  0, // auto aspect ratio
		Charset: *charset,
		Color:   *color,
		Dither:  *dither,
	}

	var img image.Image

	if *filePath != "" {
		// Load from file
		f, err := os.Open(*filePath)
		if err != nil {
			log.Fatalf("Failed to open %q: %v", *filePath, err)
		}
		defer f.Close()
		img, _, err = image.Decode(f)
		if err != nil {
			log.Fatalf("Failed to decode %q: %v", *filePath, err)
		}
		fmt.Printf("Input: %s (%dx%d)\n", *filePath, img.Bounds().Dx(), img.Bounds().Dy())
	} else {
		// Generate a built-in test pattern
		img = generateTestPattern(100, 60)
		fmt.Printf("Input: built-in test pattern (%dx%d)\n", img.Bounds().Dx(), img.Bounds().Dy())
	}

	fmt.Printf("Options: width=%d color=%v dither=%v charset=%q\n",
		opts.Width, opts.Color, opts.Dither, opts.Charset)

	// Convert to ASCII
	art := ascii.FromImage(img, opts)

	// Print result
	fmt.Println("\nOutput:")
	fmt.Println(strings.Repeat("─", *width))
	fmt.Print(art)
	fmt.Println(strings.Repeat("─", *width))

	// Print stats
	lines := strings.Count(art, "\n")
	fmt.Printf("Dimensions: %d × %d characters\n", *width, lines-1)
}

// generateTestPattern creates a synthetic test image with gradients,
// shapes, and text-like patterns for visual evaluation.
func generateTestPattern(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Fill with a gradient
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := uint8(float64(x) / float64(w) * 255)
			g := uint8(float64(y) / float64(h) * 255)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// Draw rectangles
	drawRect(img, 10, 10, 30, 25, color.RGBA{255, 80, 80, 255})
	drawRect(img, 50, 10, 30, 25, color.RGBA{80, 255, 80, 255})
	drawRect(img, 10, 35, 30, 20, color.RGBA{80, 80, 255, 255})
	drawRect(img, 50, 35, 40, 20, color.RGBA{255, 255, 80, 255})

	// Draw a circle
	drawCircle(img, 85, 45, 12, color.RGBA{255, 80, 255, 255})

	// Draw a diagonal line
	for i := 0; i < 40; i++ {
		img.Set(i, i, color.RGBA{255, 255, 255, 255})
	}

	return img
}

func drawRect(img *image.RGBA, x, y, w, h int, c color.RGBA) {
	for dy := 0; dy < h && y+dy < img.Bounds().Dy(); dy++ {
		for dx := 0; dx < w && x+dx < img.Bounds().Dx(); dx++ {
			img.Set(x+dx, y+dy, c)
		}
	}
	// Border in white
	border := color.RGBA{255, 255, 255, 255}
	for dx := 0; dx < w && x+dx < img.Bounds().Dx(); dx++ {
		if y >= 0 && y < img.Bounds().Dy() {
			img.Set(x+dx, y, border)
		}
		if y+h-1 >= 0 && y+h-1 < img.Bounds().Dy() {
			img.Set(x+dx, y+h-1, border)
		}
	}
	for dy := 0; dy < h && y+dy < img.Bounds().Dy(); dy++ {
		if x >= 0 && x < img.Bounds().Dx() {
			img.Set(x, y+dy, border)
		}
		if x+w-1 >= 0 && x+w-1 < img.Bounds().Dx() {
			img.Set(x+w-1, y+dy, border)
		}
	}
}

func drawCircle(img *image.RGBA, cx, cy, r int, c color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			dx := x - cx
			dy := y - cy
			if dx*dx+dy*dy <= r*r {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, c)
				}
			}
		}
	}
	// Draw outline
	border := color.RGBA{255, 255, 255, 255}
	for angle := 0; angle < 360; angle++ {
		rad := float64(angle) * 3.14159 / 180.0
		x := cx + int(float64(r)*cos(rad))
		y := cy + int(float64(r)*sin(rad))
		if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
			img.Set(x, y, border)
		}
	}
}

func cos(a float64) float64 {
	// Simple Taylor approximation for demonstration
	if a < 0 {
		a = -a
	}
	a = a - 2*3.14159*float64(int(a/(2*3.14159)))
	if a < 0 {
		a += 2 * 3.14159
	}
	// Basic cosine via Chebyshev-like approximation
	if a > 3.14159 {
		return -cos(2*3.14159 - a)
	}
	a2 := a * a
	return 1 - a2/2 + a2*a2/24 - a2*a2*a2/720
}

func sin(a float64) float64 {
	return cos(a - 3.14159/2)
}
