// Package color provides color types and parsing for terminal output.
//
// Colors can be specified as named colors, hexadecimal (#RGB, #RRGGBB),
// RGB (rgb(r, g, b)), ANSI 8/16 colors, or ANSI 256-color palette indices.
// The library also supports true color (24-bit) when the terminal supports it.
package color

import (
	"fmt"
	"strconv"
	"strings"
)

// Color represents a terminal color.
// It can be a 4-bit ANSI color, an 8-bit (256) palette color,
// or a 24-bit true color.
type Color struct {
	// Type specifies the kind of color.
	Type ColorType
	// Index holds the ANSI palette index (for ANSI and 256 modes).
	Index uint8
	// R, G, B hold the true color components.
	R, G, B uint8
}

// ColorType distinguishes between different color representations.
type ColorType uint8

const (
	ColorANSI    ColorType = iota // 4-bit ANSI color (0-15)
	Color256                      // 8-bit 256-color palette (0-255)
	ColorTrue                     // 24-bit true color
)

// Predefined ANSI color indices.
const (
	ANSIBlack        uint8 = 0
	ANSIRed          uint8 = 1
	ANSIGreen        uint8 = 2
	ANSIYellow       uint8 = 3
	ANSIBlue         uint8 = 4
	ANSIMagenta      uint8 = 5
	ANSICyan         uint8 = 6
	ANSIWhite        uint8 = 7
	ANSIBrightBlack  uint8 = 8
	ANSIBrightRed    uint8 = 9
	ANSIBrightGreen  uint8 = 10
	ANSIBrightYellow uint8 = 11
	ANSIBrightBlue   uint8 = 12
	ANSIBrightMagenta uint8 = 13
	ANSIBrightCyan   uint8 = 14
	ANSIBrightWhite  uint8 = 15
)

// NamedColors maps CSS named colors to their 24-bit RGB values.
var NamedColors = map[string]Color{
	"black":   {Type: ColorTrue, R: 0, G: 0, B: 0},
	"silver":  {Type: ColorTrue, R: 192, G: 192, B: 192},
	"gray":    {Type: ColorTrue, R: 128, G: 128, B: 128},
	"white":   {Type: ColorTrue, R: 255, G: 255, B: 255},
	"maroon":  {Type: ColorTrue, R: 128, G: 0, B: 0},
	"red":     {Type: ColorTrue, R: 255, G: 0, B: 0},
	"purple":  {Type: ColorTrue, R: 128, G: 0, B: 128},
	"fuchsia": {Type: ColorTrue, R: 255, G: 0, B: 255},
	"green":   {Type: ColorTrue, R: 0, G: 128, B: 0},
	"lime":    {Type: ColorTrue, R: 0, G: 255, B: 0},
	"olive":   {Type: ColorTrue, R: 128, G: 128, B: 0},
	"yellow":  {Type: ColorTrue, R: 255, G: 255, B: 0},
	"navy":    {Type: ColorTrue, R: 0, G: 0, B: 128},
	"blue":    {Type: ColorTrue, R: 0, G: 0, B: 255},
	"teal":    {Type: ColorTrue, R: 0, G: 128, B: 128},
	"aqua":    {Type: ColorTrue, R: 0, G: 255, B: 255},
	"orange":  {Type: ColorTrue, R: 255, G: 165, B: 0},
	"pink":    {Type: ColorTrue, R: 255, G: 192, B: 203},
	"brown":   {Type: ColorTrue, R: 165, G: 42, B: 42},
	"coral":   {Type: ColorTrue, R: 255, G: 127, B: 80},
	"crimson": {Type: ColorTrue, R: 220, G: 20, B: 60},
	"darkblue":   {Type: ColorTrue, R: 0, G: 0, B: 139},
	"darkcyan":   {Type: ColorTrue, R: 0, G: 139, B: 139},
	"darkgray":   {Type: ColorTrue, R: 169, G: 169, B: 169},
	"darkgreen":  {Type: ColorTrue, R: 0, G: 100, B: 0},
	"darkorange": {Type: ColorTrue, R: 255, G: 140, B: 0},
	"darkred":    {Type: ColorTrue, R: 139, G: 0, B: 0},
	"darkviolet": {Type: ColorTrue, R: 148, G: 0, B: 211},
	"gold":       {Type: ColorTrue, R: 255, G: 215, B: 0},
	"indigo":     {Type: ColorTrue, R: 75, G: 0, B: 130},
	"ivory":      {Type: ColorTrue, R: 255, G: 255, B: 240},
	"khaki":      {Type: ColorTrue, R: 240, G: 230, B: 140},
	"lavender":   {Type: ColorTrue, R: 230, G: 230, B: 250},
	"lightblue":   {Type: ColorTrue, R: 173, G: 216, B: 230},
	"lightgray":   {Type: ColorTrue, R: 211, G: 211, B: 211},
	"lightgreen":  {Type: ColorTrue, R: 144, G: 238, B: 144},
	"lightyellow": {Type: ColorTrue, R: 255, G: 255, B: 224},
	"limegreen":   {Type: ColorTrue, R: 50, G: 205, B: 50},
	"magenta":     {Type: ColorTrue, R: 255, G: 0, B: 255},
	"peru":        {Type: ColorTrue, R: 205, G: 133, B: 63},
	"plum":        {Type: ColorTrue, R: 221, G: 160, B: 221},
	"salmon":      {Type: ColorTrue, R: 250, G: 128, B: 114},
	"sienna":      {Type: ColorTrue, R: 160, G: 82, B: 45},
	"snow":        {Type: ColorTrue, R: 255, G: 250, B: 250},
	"tan":         {Type: ColorTrue, R: 210, G: 180, B: 140},
	"tomato":      {Type: ColorTrue, R: 255, G: 99, B: 71},
	"turquoise":   {Type: ColorTrue, R: 64, G: 224, B: 208},
	"violet":      {Type: ColorTrue, R: 238, G: 130, B: 238},
	"wheat":       {Type: ColorTrue, R: 245, G: 222, B: 179},
}

// ANSIColors maps common ANSI color names to their indices.
var ANSIColors = map[string]uint8{
	"black":         ANSIBlack,
	"red":           ANSIRed,
	"green":         ANSIGreen,
	"yellow":        ANSIYellow,
	"blue":          ANSIBlue,
	"magenta":       ANSIMagenta,
	"cyan":          ANSICyan,
	"white":         ANSIWhite,
	"brightblack":   ANSIBrightBlack,
	"brightred":     ANSIBrightRed,
	"brightgreen":   ANSIBrightGreen,
	"brightyellow":  ANSIBrightYellow,
	"brightblue":    ANSIBrightBlue,
	"brightmagenta": ANSIBrightMagenta,
	"brightcyan":    ANSIBrightCyan,
	"brightwhite":   ANSIBrightWhite,
}

// NewANSI creates an ANSI color from a 0-15 index.
func NewANSI(index uint8) Color {
	if index > 15 {
		index = index % 16
	}
	return Color{Type: ColorANSI, Index: index}
}

// New256 creates a 256-color palette color.
func New256(index uint8) Color {
	return Color{Type: Color256, Index: index}
}

// NewTrue creates a 24-bit true color.
func NewTrue(r, g, b uint8) Color {
	return Color{Type: ColorTrue, R: r, G: g, B: b}
}

// HSL creates a Color from HSL (Hue, Saturation, Lightness) values.
// h: 0-360 degrees, s: 0-100 (percent), l: 0-100 (percent)
func HSL(h, s, l float64) Color {
	r, g, b := hslToRGB(h, s/100, l/100)
	return Color{Type: ColorTrue, R: r, G: g, B: b}
}

// HSLA creates a Color from HSLA values with alpha.
func HSLA(h, s, l, a float64) Color {
	r, g, b := hslToRGB(h, s/100, l/100)
	_ = a // alpha not stored in Color struct
	return Color{Type: ColorTrue, R: r, G: g, B: b}
}

// hslToRGB converts HSL (hue 0-360, s 0-1, l 0-1) to RGB 0-255.
func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	// Normalize hue to 0-360
	h = normalizeHue(h)
	if s < 0 { s = 0 }
	if s > 1 { s = 1 }
	if l < 0 { l = 0 }
	if l > 1 { l = 1 }

	c := (1 - abs(2*l-1)) * s
	x := c * (1 - abs(mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return uint8((r + m) * 255), uint8((g + m) * 255), uint8((b + m) * 255)
}

// abs returns the absolute value of a float64.
func abs(x float64) float64 {
	if x < 0 { return -x }
	return x
}

// mod implements floating-point modulo (always positive result).
func mod(x, m float64) float64 {
	if m == 0 { return 0 }
	x = x - m*float64(int(x/m))
	if x < 0 { x += m }
	return x
}

// normalizeHue normalizes hue to 0-360 range.
func normalizeHue(h float64) float64 {
	h = mod(h, 360)
	if h < 0 { h += 360 }
	return h
}

// Lerp linearly interpolates between two colors by factor t (0=color a, 1=color b).
func Lerp(a, b Color, t float64) Color {
	if t <= 0 { return a }
	if t >= 1 { return b }
	return Color{Type: ColorTrue,
		R: uint8(float64(a.R)*(1-t) + float64(b.R)*t),
		G: uint8(float64(a.G)*(1-t) + float64(b.G)*t),
		B: uint8(float64(a.B)*(1-t) + float64(b.B)*t),
	}
}

// Darken reduces the brightness of a color by the given percentage (0-100).
func Darken(c Color, amount float64) Color {
	if amount <= 0 { return c }
	if amount > 100 { amount = 100 }
	h, s, l := rgbToHSL(c.R, c.G, c.B)
	l -= l * amount / 100
	if l < 0 { l = 0 }
	return HSL(h, s, l)
}

// Lighten increases the brightness of a color by the given percentage (0-100).
func Lighten(c Color, amount float64) Color {
	if amount <= 0 { return c }
	if amount > 100 { amount = 100 }
	h, s, l := rgbToHSL(c.R, c.G, c.B)
	l += (100 - l) * amount / 100
	if l > 100 { l = 100 }
	return HSL(h, s, l)
}

// Saturate increases color saturation by percentage (0-100).
func Saturate(c Color, amount float64) Color {
	if amount <= 0 { return c }
	h, s, l := rgbToHSL(c.R, c.G, c.B)
	s += (100 - s) * amount / 100
	if s > 100 { s = 100 }
	return HSL(h, s, l)
}

// Desaturate reduces color saturation by percentage (0-100).
func Desaturate(c Color, amount float64) Color {
	if amount <= 0 { return c }
	h, s, l := rgbToHSL(c.R, c.G, c.B)
	s -= s * amount / 100
	if s < 0 { s = 0 }
	return HSL(h, s, l)
}

// RotateHue shifts the hue by the given degrees (can be negative).
func RotateHue(c Color, degrees float64) Color {
	h, s, l := rgbToHSL(c.R, c.G, c.B)
	h += degrees
	return HSL(h, s, l)
}

// Grayscale converts a color to grayscale using luminance weights.
func Grayscale(c Color) Color {
	lum := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
	g := uint8(lum)
	return Color{Type: ColorTrue, R: g, G: g, B: g}
}

// Invert inverts all RGB channels (255 - value).
func Invert(c Color) Color {
	return Color{Type: ColorTrue, R: 255 - c.R, G: 255 - c.G, B: 255 - c.B}
}

// Luminance returns the relative luminance of a color (WCAG formula).
func Luminance(c Color) float64 {
	r := linearize(float64(c.R) / 255)
	g := linearize(float64(c.G) / 255)
	b := linearize(float64(c.B) / 255)
	return 0.2126*r + 0.7152*g + 0.0722*b
}

// linearize applies the sRGB linearization curve for luminance calculation.
func linearize(ch float64) float64 {
	if ch <= 0.04045 {
		return ch / 12.92
	}
	return ((ch + 0.055) / 1.055) * ((ch + 0.055) / 1.055)
}

// ContrastRatio calculates the WCAG contrast ratio between two colors.
// Returns a value from 1:1 (no contrast) to 21:1 (maximum).
func ContrastRatio(fg, bg Color) float64 {
	l1 := Luminance(fg)
	l2 := Luminance(bg)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

// rgbToHSL converts RGB (0-255) to HSL (hue 0-360, sat 0-100, light 0-100).
func rgbToHSL(r, g, b uint8) (float64, float64, float64) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := maxVal(rf, gf, bf)
	min := minVal(rf, gf, bf)
	l := (max + min) / 2

	if max == min {
		return 0, 0, l * 100
	}

	d := max - min
	var s float64
	if l > 0.5 {
		s = d / (2 - max - min)
	} else {
		s = d / (max + min)
	}

	var h float64
	switch {
	case max == rf:
		h = (gf - bf) / d
		if gf < bf {
			h += 6
		}
	case max == gf:
		h = (bf-rf)/d + 2
	default:
		h = (rf-gf)/d + 4
	}
	h *= 60

	return h, s * 100, l * 100
}

func maxVal(a, b, c float64) float64 {
	if a > b {
		if a > c { return a }
		return c
	}
	if b > c { return b }
	return c
}

func minVal(a, b, c float64) float64 {
	if a < b {
		if a < c { return a }
		return c
	}
	if b < c { return b }
	return c
}

// Predefined common colors as exported constants for convenience.
var (
	Black       = Color{Type: ColorTrue, R: 0, G: 0, B: 0}
	White       = Color{Type: ColorTrue, R: 255, G: 255, B: 255}
	Red         = Color{Type: ColorTrue, R: 255, G: 0, B: 0}
	Green       = Color{Type: ColorTrue, R: 0, G: 128, B: 0}
	Blue        = Color{Type: ColorTrue, R: 0, G: 0, B: 255}
	Yellow      = Color{Type: ColorTrue, R: 255, G: 255, B: 0}
	Cyan        = Color{Type: ColorTrue, R: 0, G: 255, B: 255}
	Magenta     = Color{Type: ColorTrue, R: 255, G: 0, B: 255}
	Orange      = Color{Type: ColorTrue, R: 255, G: 165, B: 0}
	Lime        = Color{Type: ColorTrue, R: 0, G: 255, B: 0}
	Pink        = Color{Type: ColorTrue, R: 255, G: 192, B: 203}
	Purple      = Color{Type: ColorTrue, R: 128, G: 0, B: 128}
	Navy        = Color{Type: ColorTrue, R: 0, G: 0, B: 128}
	Teal        = Color{Type: ColorTrue, R: 0, G: 128, B: 128}
	Gray        = Color{Type: ColorTrue, R: 128, G: 128, B: 128}
	Silver      = Color{Type: ColorTrue, R: 192, G: 192, B: 192}
	Transparent = Color{Type: ColorTrue, R: 0, G: 0, B: 0}
)

// Hex parses a hex color string (#RGB, #RRGGBB, #RRGGBBAA) and returns the Color.
// Returns Black if the string cannot be parsed.
func Hex(s string) Color {
	c, ok := ParseColor(s)
	if !ok {
		return Black
	}
	return c
}

// ParseColor parses a CSS color string and returns the corresponding Color.
// Supported formats:
//   - Named colors: "red", "blue", etc.
//   - Hexadecimal: "#RGB", "#RRGGBB", "#RRGGBBAA"
//   - RGB: "rgb(r, g, b)"
//   - RGBA: "rgba(r, g, b, a)"
//   - ANSI: "ansi(0)" through "ansi(15)"
//   - 256: "color(0)" through "color(255)"
//   - Transparent: "transparent"
//   - CurrentColor: "currentcolor" (returns a special marker)
func ParseColor(s string) (Color, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Color{}, false
	}

	// Handle special keywords
	if strings.EqualFold(s, "transparent") {
		return NewTrue(0, 0, 0), true
	}
	if strings.EqualFold(s, "currentcolor") {
		return NewTrue(0, 0, 0), true
	}

	// Named colors (CSS named colors take precedence over ANSI names for true color)
	if c, ok := NamedColors[strings.ToLower(s)]; ok {
		return c, true
	}

	// ANSI named colors (map to ANSI indices)
	if idx, ok := ANSIColors[strings.ToLower(s)]; ok {
		return NewANSI(idx), true
	}

	// Hex color
	if s[0] == '#' {
		return parseHex(s)
	}

	// rgb() function
	if strings.HasPrefix(strings.ToLower(s), "rgb(") {
		return parseRGB(s)
	}

	// rgba() function
	if strings.HasPrefix(strings.ToLower(s), "rgba(") {
		return parseRGBA(s)
	}

	// hsl() function
	if strings.HasPrefix(strings.ToLower(s), "hsl(") {
		return parseHSL(s)
	}

	// hsla() function
	if strings.HasPrefix(strings.ToLower(s), "hsla(") {
		return parseHSLA(s)
	}

	// ansi() function
	if strings.HasPrefix(strings.ToLower(s), "ansi(") {
		return parseANSIFunc(s)
	}

	// color() function (256 color)
	if strings.HasPrefix(strings.ToLower(s), "color(") {
		return parseColorFunc(s)
	}

	return Color{}, false
}

func parseHex(s string) (Color, bool) {
	s = strings.TrimPrefix(s, "#")
	if len(s) == 0 {
		return Color{}, false
	}

	switch len(s) {
	case 3:
		// #RGB -> #RRGGBB
		r, err1 := strconv.ParseUint(s[0:1]+s[0:1], 16, 8)
		g, err2 := strconv.ParseUint(s[1:2]+s[1:2], 16, 8)
		b, err3 := strconv.ParseUint(s[2:3]+s[2:3], 16, 8)
		if err1 != nil || err2 != nil || err3 != nil {
			return Color{}, false
		}
		return NewTrue(uint8(r), uint8(g), uint8(b)), true

	case 6:
		r, err1 := strconv.ParseUint(s[0:2], 16, 8)
		g, err2 := strconv.ParseUint(s[2:4], 16, 8)
		b, err3 := strconv.ParseUint(s[4:6], 16, 8)
		if err1 != nil || err2 != nil || err3 != nil {
			return Color{}, false
		}
		return NewTrue(uint8(r), uint8(g), uint8(b)), true

	case 8:
		// #RRGGBBAA
		r, err1 := strconv.ParseUint(s[0:2], 16, 8)
		g, err2 := strconv.ParseUint(s[2:4], 16, 8)
		b, err3 := strconv.ParseUint(s[4:6], 16, 8)
		_, err4 := strconv.ParseUint(s[6:8], 16, 8) // alpha is parsed but we don't store it
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return Color{}, false
		}
		return NewTrue(uint8(r), uint8(g), uint8(b)), true

	default:
		return Color{}, false
	}
}

func parseRGB(s string) (Color, bool) {
	// rgb(r, g, b)
	inner := strings.TrimSpace(s[4 : len(s)-1])
	parts := strings.Split(inner, ",")
	if len(parts) != 3 {
		return Color{}, false
	}

	r, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	g, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	b, err3 := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err1 != nil || err2 != nil || err3 != nil {
		return Color{}, false
	}

	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return Color{}, false
	}

	return NewTrue(uint8(r), uint8(g), uint8(b)), true
}

// parseRGBA parses an rgba(r, g, b, a) color string.
// The alpha value is parsed but currently not stored (Color has no alpha field).
func parseRGBA(s string) (Color, bool) {
	// rgba(r, g, b, a)
	inner := strings.TrimSpace(s[5 : len(s)-1])
	parts := strings.Split(inner, ",")
	if len(parts) != 4 {
		return Color{}, false
	}

	r, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	g, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	b, err3 := strconv.Atoi(strings.TrimSpace(parts[2]))
	_, err4 := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64) // alpha parsed but not stored
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return Color{}, false
	}

	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return Color{}, false
	}

	return NewTrue(uint8(r), uint8(g), uint8(b)), true
}

// parseHSL parses an hsl(h, s%, l%) color string.
func parseHSL(s string) (Color, bool) {
	// hsl(h, s%, l%)
	inner := strings.TrimSpace(s[4 : len(s)-1])
	parts := strings.Split(inner, ",")
	if len(parts) != 3 {
		return Color{}, false
	}

	h, err1 := parseHSLValue(strings.TrimSpace(parts[0]))
	sat, err2 := parseHSLPercent(strings.TrimSpace(parts[1]))
	l, err3 := parseHSLPercent(strings.TrimSpace(parts[2]))
	if err1 != nil || err2 != nil || err3 != nil {
		return Color{}, false
	}

	return HSL(h, sat, l), true
}

// parseHSLA parses an hsla(h, s%, l%, a) color string.
func parseHSLA(s string) (Color, bool) {
	// hsla(h, s%, l%, a)
	inner := strings.TrimSpace(s[5 : len(s)-1])
	parts := strings.Split(inner, ",")
	if len(parts) != 4 {
		return Color{}, false
	}

	h, err1 := parseHSLValue(strings.TrimSpace(parts[0]))
	sat, err2 := parseHSLPercent(strings.TrimSpace(parts[1]))
	l, err3 := parseHSLPercent(strings.TrimSpace(parts[2]))
	_, err4 := parseHSLAlpha(strings.TrimSpace(parts[3]))
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return Color{}, false
	}

	return HSL(h, sat, l), true
}

// parseHSLValue parses a hue value (0-360, can be a float).
func parseHSLValue(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSuffix(s, "deg"), 64)
}

// parseHSLPercent parses a percentage value (0-100) or a 0-1 float.
func parseHSLPercent(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		return strconv.ParseFloat(s[:len(s)-1], 64)
	}
	// Treat as 0-1 value, convert to percentage
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return v * 100, nil
}

// parseHSLAlpha parses an alpha value for hsla().
func parseHSLAlpha(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

func parseANSIFunc(s string) (Color, bool) {
	inner := strings.TrimSpace(s[5 : len(s)-1])
	idx, err := strconv.Atoi(inner)
	if err != nil || idx < 0 || idx > 15 {
		return Color{}, false
	}
	return NewANSI(uint8(idx)), true
}

func parseColorFunc(s string) (Color, bool) {
	inner := strings.TrimSpace(s[6 : len(s)-1])
	idx, err := strconv.Atoi(inner)
	if err != nil || idx < 0 || idx > 255 {
		return Color{}, false
	}
	return New256(uint8(idx)), true
}

// ANSIOutput returns the ANSI escape sequence for setting the foreground
// color to this color, assuming the terminal supports the given mode.
// mode: 0 = 16 colors (ANSI), 1 = 256 colors, 2 = true color
func (c Color) ANSI(mode int) string {
	switch c.Type {
	case ColorANSI:
		if c.Index <= 7 {
			return fmt.Sprintf("\x1b[%dm", 30+c.Index)
		}
		return fmt.Sprintf("\x1b[%dm", 82+c.Index)
	case Color256:
		return fmt.Sprintf("\x1b[38;5;%dm", c.Index)
	case ColorTrue:
		if mode >= 2 {
			return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
		}
		// Fallback: approximate to 256-color palette
		idx := approximateTo256(c.R, c.G, c.B)
		return fmt.Sprintf("\x1b[38;5;%dm", idx)
	}
	return ""
}

// ANSIBackground returns the ANSI escape sequence for setting the background
// color.
func (c Color) ANSIBackground(mode int) string {
	switch c.Type {
	case ColorANSI:
		if c.Index <= 7 {
			return fmt.Sprintf("\x1b[%dm", 40+c.Index)
		}
		return fmt.Sprintf("\x1b[%dm", 92+c.Index)
	case Color256:
		return fmt.Sprintf("\x1b[48;5;%dm", c.Index)
	case ColorTrue:
		if mode >= 2 {
			return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
		}
		idx := approximateTo256(c.R, c.G, c.B)
		return fmt.Sprintf("\x1b[48;5;%dm", idx)
	}
	return ""
}

// approximateTo256 converts an RGB color to the nearest 256-color palette index.
func approximateTo256(r, g, b uint8) uint8 {
	// 6x6x6 cube (216 colors) + grayscale ramp (24 colors) + 16 basic colors
	// For simplicity, use the web-safe cube approach.
	if r == g && g == b {
		// Grayscale
		if r < 8 {
			return 16
		}
		idx := (int(r)-8)/10 + 232
		if idx > 255 {
			idx = 255
		}
		return uint8(idx)
	}

	// 6x6x6 cube starting at index 16
	ir := int(r) * 5 / 255
	ig := int(g) * 5 / 255
	ib := int(b) * 5 / 255
	return uint8(16 + ir*36 + ig*6 + ib)
}

// String returns a human-readable representation of the color.
func (c Color) String() string {
	switch c.Type {
	case ColorANSI:
		return fmt.Sprintf("ansi(%d)", c.Index)
	case Color256:
		return fmt.Sprintf("color(%d)", c.Index)
	case ColorTrue:
		return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
	}
	return "unknown"
}

// ANSI256Palette is the complete 256-color ANSI palette.
// Indices 0-15: standard ANSI colors
// Indices 16-231: 6x6x6 RGB cube
// Indices 232-255: grayscale ramp
var ANSI256Palette [256]Color

func init() {
	// Initialize the 256-color palette
	// Standard 16 ANSI colors
	ansiColors := []struct{ r, g, b uint8 }{
		{0, 0, 0},       // 0: Black
		{128, 0, 0},     // 1: Red
		{0, 128, 0},     // 2: Green
		{128, 128, 0},   // 3: Yellow
		{0, 0, 128},     // 4: Blue
		{128, 0, 128},   // 5: Magenta
		{0, 128, 128},   // 6: Cyan
		{192, 192, 192}, // 7: White
		{128, 128, 128}, // 8: Bright Black (Gray)
		{255, 0, 0},     // 9: Bright Red
		{0, 255, 0},     // 10: Bright Green
		{255, 255, 0},   // 11: Bright Yellow
		{0, 0, 255},     // 12: Bright Blue
		{255, 0, 255},   // 13: Bright Magenta
		{0, 255, 255},   // 14: Bright Cyan
		{255, 255, 255}, // 15: Bright White
	}
	for i, c := range ansiColors {
		ANSI256Palette[i] = Color{Type: ColorTrue, R: c.r, G: c.g, B: c.b}
	}

	// 6x6x6 color cube (indices 16-231)
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				idx := 16 + r*36 + g*6 + b
				rr := uint8(r * 255 / 5)
				gg := uint8(g * 255 / 5)
				bb := uint8(b * 255 / 5)
				ANSI256Palette[idx] = Color{Type: ColorTrue, R: rr, G: gg, B: bb}
			}
		}
	}

	// Grayscale ramp (indices 232-255)
	for i := 0; i < 24; i++ {
		v := uint8(i*10 + 8)
		ANSI256Palette[232+i] = Color{Type: ColorTrue, R: v, G: v, B: v}
	}
}

// To256 finds the nearest 256-color palette index for a given color.
// This is the public equivalent of approximateTo256.
func To256(c Color) uint8 {
	if c.Type == Color256 {
		return c.Index
	}
	if c.Type == ColorANSI {
		return uint8(c.Index)
	}
	return approximateTo256(c.R, c.G, c.B)
}

// ToANSI finds the nearest standard ANSI 16-color index for a given color.
func ToANSI(c Color) uint8 {
	if c.Type == ColorANSI {
		return c.Index
	}
	if c.Type == Color256 {
		c = ToTrue(c)
	}

	// Find nearest of the 16 standard ANSI colors
	bestIdx := uint8(0)
	bestDist := int(^uint(0) >> 1)
	for i := 0; i < 16; i++ {
		p := ANSI256Palette[i]
		dr := int(c.R) - int(p.R)
		dg := int(c.G) - int(p.G)
		db := int(c.B) - int(p.B)
		dist := dr*dr + dg*dg + db*db
		if dist < bestDist {
			bestDist = dist
			bestIdx = uint8(i)
		}
	}
	return bestIdx
}

// ToTrue converts an ANSI or 256-color to its true color RGB equivalent.
// If the color is already true color, returns it unchanged.
func ToTrue(c Color) Color {
	if c.Type == ColorTrue {
		return c
	}
	if c.Type == ColorANSI || c.Type == Color256 {
		return ANSI256Palette[c.Index]
	}
	return c
}

// ClosestPalette finds the index of the closest color in a palette
// using Euclidean distance in RGB space.
func ClosestPalette(c Color, palette []Color) int {
	if len(palette) == 0 {
		return -1
	}
	bestIdx := 0
	bestDist := int(^uint(0) >> 1)
	for i, p := range palette {
		dr := int(c.R) - int(p.R)
		dg := int(c.G) - int(p.G)
		db := int(c.B) - int(p.B)
		dist := dr*dr + dg*dg + db*db
		if dist < bestDist {
			bestDist = dist
			bestIdx = i
		}
	}
	return bestIdx
}
