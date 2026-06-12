package ascii

import (
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strings"
)

// ImageOptions controls the image-to-ASCII conversion.
type ImageOptions struct {
	// Width is the output width in characters.
	// If 0, defaults to terminal-width detection (falls back to 80).
	Width int

	// Height is the output height in characters.
	// If 0, the aspect ratio is preserved based on Width.
	Height int

	// Charset is the character ramp used for brightness mapping.
	// If empty, CharsetStandard is used.
	Charset string

	// Color enables true-color ANSI output (each char's foreground color
	// matches the original pixel color). When false, output is grayscale.
	Color bool

	// Dither enables Floyd-Steinberg error diffusion dithering,
	// which improves perceived quality at low resolutions.
	Dither bool

	// Scale is an optional factor applied after Width/Height.
	// 1.0 = normal, 2.0 = double size, 0.5 = half size.
	Scale float64
}

// Default character ramps
const (
	// CharsetStandard: 10 levels, dense to sparse
	CharsetStandard = "@@%#*+=-:. "

	// CharsetSimple: 7 levels
	CharsetSimple = "@%#*+=-."

	// CharsetBlock: 4 block-element levels
	CharsetBlock = " ░▒▓█"

	// CharsetDetailed: 16 levels for fine gradients
	CharsetDetailed = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
)

// DefaultImageOptions returns sensible defaults.
func DefaultImageOptions() ImageOptions {
	return ImageOptions{
		Width:   80,
		Height:  0, // auto
		Charset: CharsetStandard,
		Color:   false,
		Dither:  true,
		Scale:   1.0,
	}
}

// FromFile decodes an image file (PNG, JPEG, GIF) and converts it to ASCII art.
func FromFile(path string, opts ImageOptions) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("ascii: open %q: %w", path, err)
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return "", fmt.Errorf("ascii: decode %q: %w", path, err)
	}
	_ = format

	return FromImage(img, opts), nil
}

// FromImage converts any image.Image to ASCII art using the given options.
func FromImage(img image.Image, opts ImageOptions) string {
	if opts.Width <= 0 {
		opts.Width = 80
	}
	if opts.Charset == "" {
		opts.Charset = CharsetStandard
	}
	if opts.Scale == 0 {
		opts.Scale = 1.0
	}

	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	// Compute output dimensions preserving aspect ratio
	outW := opts.Width
	if opts.Scale != 1.0 {
		outW = int(float64(outW) * opts.Scale)
	}
	if outW < 1 {
		outW = 1
	}

	outH := opts.Height
	if outH <= 0 {
		// Preserve aspect ratio.
		// Characters are roughly 2:1 height:width in terminals,
		// so we adjust by 0.5 to approximate square pixels.
		ratio := float64(srcH) / float64(srcW) * 0.5
		outH = int(float64(outW) * ratio)
	}
	if opts.Scale != 1.0 {
		outH = int(float64(outH) * opts.Scale)
	}
	if outH < 1 {
		outH = 1
	}

	// Resize image to output dimensions using nearest-neighbor sampling
	pixels := resizeNearest(img, srcW, srcH, outW, outH)

	if opts.Dither {
		pixels = floydSteinbergDither(pixels, outW, outH)
	}

	// Map pixels to characters
	charset := []rune(opts.Charset)
	charsetLen := len(charset)

	var sb strings.Builder
	sb.Grow((outW + 1) * outH)

	for y := 0; y < outH; y++ {
		for x := 0; x < outW; x++ {
			p := pixels[y*outW+x]

			if opts.Color {
				// True-color output: use ANSI 24-bit foreground color
				r, g, b := p.R, p.G, p.B
				// Brightness determines which character to use
				lum := luminance(r, g, b)
				idx := clamp(luminanceToIndex(lum, charsetLen), 0, charsetLen-1)
				ch := charset[idx]
				sb.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s", r, g, b, string(ch)))
			} else {
				lum := p.R // grayscale
				idx := clamp(luminanceToIndex(lum, charsetLen), 0, charsetLen-1)
				sb.WriteRune(charset[idx])
			}
		}
		if opts.Color {
			sb.WriteString("\x1b[0m")
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

// px represents a single pixel with RGBA values.
type px struct {
	R, G, B, A uint32
}

// resizeNearest resizes an image to the given dimensions using nearest-neighbor sampling.
func resizeNearest(img image.Image, srcW, srcH, dstW, dstH int) []px {
	pixels := make([]px, dstW*dstH)
	for y := 0; y < dstH; y++ {
		for x := 0; x < dstW; x++ {
			// Map destination pixel to source coordinates
			srcX := int(float64(x) * float64(srcW) / float64(dstW))
			srcY := int(float64(y) * float64(srcH) / float64(dstH))
			if srcX >= srcW {
				srcX = srcW - 1
			}
			if srcY >= srcH {
				srcY = srcH - 1
			}
			r, g, b, a := img.At(srcX, srcY).RGBA()
			// Premultiply alpha into color (assume opaque background)
			if a > 0 {
				r = r * 0xffff / a
				g = g * 0xffff / a
				b = b * 0xffff / a
			}
			pixels[y*dstW+x] = px{
				R: uint32(uint8(r >> 8)),
				G: uint32(uint8(g >> 8)),
				B: uint32(uint8(b >> 8)),
				A: uint32(uint8(a >> 8)),
			}
		}
	}
	return pixels
}

// luminance computes relative luminance from RGB (0–255 each).
// Uses standard NTSC/Rec.601 weights.
func luminance(r, g, b uint32) uint32 {
	return uint32(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
}

// luminanceToIndex maps a luminance value (0–255) to a charset index.
func luminanceToIndex(lum uint32, charsetLen int) int {
	// Invert: high luminance (bright) → sparse characters, low luminance (dark) → dense
	// Charset goes from dense to sparse, so bright = high index
	if charsetLen <= 1 {
		return 0
	}
	inv := 255 - lum
	return int(inv) * (charsetLen - 1) / 255
}

// clamp restricts n to [lo, hi].
func clamp(n, lo, hi int) int {
	if n < lo {
		return lo
	}
	if n > hi {
		return hi
	}
	return n
}

// floydSteinbergDither applies Floyd-Steinberg error diffusion dithering.
// It diffuses quantization error to neighboring pixels.
func floydSteinbergDither(pixels []px, w, h int) []px {
	// Work on a copy so we don't modify the original
	work := make([]float64, w*h)
	for i, p := range pixels {
		work[i] = float64(luminance(p.R, p.G, p.B))
	}

	numLevels := 16 // quantize to 16 levels for dithering
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*w + x
			old := work[idx]
			newVal := math.Round(old/255.0*float64(numLevels-1)) / float64(numLevels-1) * 255.0
			err := old - newVal
			work[idx] = newVal

			// Diffuse error to neighbors (Floyd-Steinberg weights)
			if x+1 < w {
				work[y*w+(x+1)] += err * 7.0 / 16.0
			}
			if y+1 < h {
				if x-1 >= 0 {
					work[(y+1)*w+(x-1)] += err * 3.0 / 16.0
				}
				work[(y+1)*w+x] += err * 5.0 / 16.0
				if x+1 < w {
					work[(y+1)*w+(x+1)] += err * 1.0 / 16.0
				}
			}
		}
	}

	// Convert quantized luminance back to grayscale pixels
	out := make([]px, w*h)
	for i, l := range work {
		v := uint32(clamp(int(math.Round(l)), 0, 255))
		out[i] = px{R: v, G: v, B: v, A: 255}
	}
	return out
}

// FromFileGIF extracts frames from an animated GIF and returns each as ASCII art.
// Returns a slice of strings, one per frame.
func FromFileGIF(path string, opts ImageOptions) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ascii: open %q: %w", path, err)
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, fmt.Errorf("ascii: decode gif %q: %w", path, err)
	}

	frames := make([]string, len(g.Image))
	for i, frame := range g.Image {
		frames[i] = FromImage(frame, opts)
	}
	return frames, nil
}
