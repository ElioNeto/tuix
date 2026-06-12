package color

import (
	"testing"
)

func TestParseHex6(t *testing.T) {
	c, ok := ParseColor("#ff6600")
	if !ok {
		t.Fatal("failed to parse #ff6600")
	}
	if c.R != 0xFF || c.G != 0x66 || c.B != 0x00 {
		t.Fatalf("expected (255,102,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
	if c.Type != ColorTrue {
		t.Fatalf("expected ColorTrue, got %v", c.Type)
	}
}

func TestParseHex3(t *testing.T) {
	c, ok := ParseColor("#f60")
	if !ok {
		t.Fatal("failed to parse #f60")
	}
	if c.R != 0xFF || c.G != 0x66 || c.B != 0x00 {
		t.Fatalf("expected (255,102,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestParseNamedColor(t *testing.T) {
	c, ok := ParseColor("red")
	if !ok {
		t.Fatal("failed to parse 'red'")
	}
	if c.R != 0xFF || c.G != 0x00 || c.B != 0x00 {
		t.Fatalf("expected (255,0,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestParseNamedGreen(t *testing.T) {
	c, ok := ParseColor("lime")
	if !ok {
		t.Fatal("failed to parse 'lime'")
	}
	if c.R != 0x00 || c.G != 0xFF || c.B != 0x00 {
		t.Fatalf("expected (0,255,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestParseNamedBlue(t *testing.T) {
	c, ok := ParseColor("blue")
	if !ok {
		t.Fatal("failed to parse 'blue'")
	}
	if c.R != 0x00 || c.G != 0x00 || c.B != 0xFF {
		t.Fatalf("expected (0,0,255), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestParseInvalidColor(t *testing.T) {
	_, ok := ParseColor("notacolor")
	if ok {
		t.Fatal("expected false for invalid color")
	}
}

func TestParseEmptyColor(t *testing.T) {
	_, ok := ParseColor("")
	if ok {
		t.Fatal("expected false for empty string")
	}
}

func TestNewTrue(t *testing.T) {
	c := NewTrue(100, 150, 200)
	if c.R != 100 || c.G != 150 || c.B != 200 {
		t.Fatalf("expected (100,150,200), got (%d,%d,%d)", c.R, c.G, c.B)
	}
	if c.Type != ColorTrue {
		t.Fatalf("expected ColorTrue, got %v", c.Type)
	}
}

func TestParseHex8(t *testing.T) {
	c, ok := ParseColor("#ff660080")
	if !ok {
		t.Fatal("failed to parse #ff660080")
	}
	if c.R != 0xFF || c.G != 0x66 || c.B != 0x00 {
		t.Fatalf("expected (255,102,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestParseRGBA(t *testing.T) {
	c, ok := ParseColor("rgba(255, 102, 0, 0.5)")
	if !ok {
		t.Fatal("failed to parse rgba(255, 102, 0, 0.5)")
	}
	if c.R != 0xFF || c.G != 0x66 || c.B != 0x00 {
		t.Fatalf("expected (255,102,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestColorConstants(t *testing.T) {
	if Red.R != 255 || Red.G != 0 || Red.B != 0 {
		t.Fatal("Red constant has wrong value")
	}
	if Green.R != 0 || Green.G != 128 || Green.B != 0 {
		t.Fatal("Green constant has wrong value")
	}
	if Blue.R != 0 || Blue.G != 0 || Blue.B != 255 {
		t.Fatal("Blue constant has wrong value")
	}
	if Black.R != 0 || Black.G != 0 || Black.B != 0 {
		t.Fatal("Black constant has wrong value")
	}
	if White.R != 255 || White.G != 255 || White.B != 255 {
		t.Fatal("White constant has wrong value")
	}
}

func TestHexHelper(t *testing.T) {
	c := Hex("#ff6600")
	if c.R != 0xFF || c.G != 0x66 || c.B != 0x00 {
		t.Fatalf("Hex('#ff6600') expected (255,102,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestHexHelperInvalid(t *testing.T) {
	c := Hex("not a color")
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatal("Hex with invalid input should return Black")
	}
}

func TestParseRGBFunction(t *testing.T) {
	c, ok := ParseColor("rgb(100, 200, 50)")
	if !ok {
		t.Fatal("failed to parse rgb(100, 200, 50)")
	}
	if c.R != 100 || c.G != 200 || c.B != 50 {
		t.Fatalf("expected (100,200,50), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestParseANSIFunction(t *testing.T) {
	c, ok := ParseColor("ansi(1)")
	if !ok {
		t.Fatal("failed to parse ansi(1)")
	}
	if c.Type != ColorANSI || c.Index != 1 {
		t.Fatalf("expected ANSI 1, got type=%v index=%d", c.Type, c.Index)
	}
}

func TestParse256Function(t *testing.T) {
	c, ok := ParseColor("color(42)")
	if !ok {
		t.Fatal("failed to parse color(42)")
	}
	if c.Type != Color256 || c.Index != 42 {
		t.Fatalf("expected Color256 42, got type=%v index=%d", c.Type, c.Index)
	}
}

func TestHSLRed(t *testing.T) {
	c := HSL(0, 100, 50)
	if c.R != 255 || c.G != 0 || c.B != 0 {
		t.Fatalf("HSL(0,100,50) expected (255,0,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestHSLGreen(t *testing.T) {
	c := HSL(120, 100, 50)
	if c.R != 0 || c.G != 255 || c.B != 0 {
		t.Fatalf("HSL(120,100,50) expected (0,255,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestHSLBlue(t *testing.T) {
	c := HSL(240, 100, 50)
	if c.R != 0 || c.G != 0 || c.B != 255 {
		t.Fatalf("HSL(240,100,50) expected (0,0,255), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestHSLBlack(t *testing.T) {
	c := HSL(0, 0, 0)
	if c.R != 0 || c.G != 0 || c.B != 0 {
		t.Fatalf("HSL(0,0,0) expected (0,0,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestHSLWhite(t *testing.T) {
	c := HSL(0, 0, 100)
	if c.R != 255 || c.G != 255 || c.B != 255 {
		t.Fatalf("HSL(0,0,100) expected (255,255,255), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestHSLOrange(t *testing.T) {
	c := HSL(30, 100, 50)
	// Orange is approximately (255, 127, 0)
	if c.R != 255 || c.G != 127 || c.B != 0 {
		t.Fatalf("HSL(30,100,50) expected (255,127,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestParseHSLFunction(t *testing.T) {
	c, ok := ParseColor("hsl(120, 100%, 50%)")
	if !ok {
		t.Fatal("failed to parse hsl(120, 100%, 50%)")
	}
	if c.R != 0 || c.G != 255 || c.B != 0 {
		t.Fatalf("hsl(120,100,50) expected (0,255,0), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestParseHSLAFunction(t *testing.T) {
	c, ok := ParseColor("hsla(240, 100%, 50%, 0.5)")
	if !ok {
		t.Fatal("failed to parse hsla(240, 100%, 50%, 0.5)")
	}
	if c.R != 0 || c.G != 0 || c.B != 255 {
		t.Fatalf("hsla(240,100,50,0.5) expected (0,0,255), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestHSLTeal(t *testing.T) {
	c := HSL(180, 100, 50)
	if c.R != 0 || c.G != 255 || c.B != 255 {
		t.Fatalf("HSL(180,100,50) expected (0,255,255), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestHSLMagenta(t *testing.T) {
	c := HSL(300, 100, 50)
	if c.R != 255 || c.G != 0 || c.B != 255 {
		t.Fatalf("HSL(300,100,50) expected (255,0,255), got (%d,%d,%d)", c.R, c.G, c.B)
	}
}

func TestHSLLightness(t *testing.T) {
	dark := HSL(0, 100, 25)
	light := HSL(0, 100, 75)
	if dark.R >= light.R {
		t.Fatalf("expected dark R < light R: %d >= %d", dark.R, light.R)
	}
}

func TestLerp(t *testing.T) {
	black := Color{Type: ColorTrue, R: 0, G: 0, B: 0}
	white := Color{Type: ColorTrue, R: 255, G: 255, B: 255}
	mid := Lerp(black, white, 0.5)
	if mid.R != 127 || mid.G != 127 || mid.B != 127 {
		t.Fatalf("Lerp(black,white,0.5) expected (127,127,127), got (%d,%d,%d)", mid.R, mid.G, mid.B)
	}
}

func TestLerpEdgeCases(t *testing.T) {
	a := Color{Type: ColorTrue, R: 100, G: 100, B: 100}
	b := Color{Type: ColorTrue, R: 200, G: 200, B: 200}
	if Lerp(a, b, 0).R != 100 { t.Fatal("Lerp(a,b,0) should return a") }
	if Lerp(a, b, 1).R != 200 { t.Fatal("Lerp(a,b,1) should return b") }
}

func TestDarken(t *testing.T) {
	c := Color{Type: ColorTrue, R: 200, G: 100, B: 50}
	dark := Darken(c, 50)
	// Darken should change the color (not produce same values)
	if dark.R == c.R && dark.G == c.G && dark.B == c.B {
		t.Fatal("Darken(50%) should change the color")
	}
	// Darkened color should have lower luminance
	if Luminance(dark) >= Luminance(c) {
		t.Fatal("Darken(50%) should reduce luminance")
	}
}

func TestLighten(t *testing.T) {
	c := Color{Type: ColorTrue, R: 100, G: 50, B: 20}
	light := Lighten(c, 50)
	if light.R <= c.R || light.G <= c.G || light.B <= c.B {
		t.Fatal("Lighten(50%) should increase all channels")
	}
}

func TestGrayscale(t *testing.T) {
	c := Color{Type: ColorTrue, R: 100, G: 150, B: 200}
	g := Grayscale(c)
	if g.R != g.G || g.G != g.B {
		t.Fatal("Grayscale should produce equal R,G,B")
	}
}

func TestInvert(t *testing.T) {
	c := Color{Type: ColorTrue, R: 100, G: 150, B: 200}
	inv := Invert(c)
	if inv.R != 155 || inv.G != 105 || inv.B != 55 {
		t.Fatalf("Invert(100,150,200) expected (155,105,55), got (%d,%d,%d)", inv.R, inv.G, inv.B)
	}
}

func TestContrastRatio(t *testing.T) {
	black := Color{Type: ColorTrue, R: 0, G: 0, B: 0}
	white := Color{Type: ColorTrue, R: 255, G: 255, B: 255}
	ratio := ContrastRatio(black, white)
	if ratio < 20 || ratio > 22 {
		t.Fatalf("ContrastRatio(black,white) expected ~21, got %f", ratio)
	}
}

func TestSaturate(t *testing.T) {
	gray := Color{Type: ColorTrue, R: 128, G: 128, B: 128}
	sat := Saturate(gray, 100)
	// A gray with 100% saturation should become a pure color
	if sat.R == sat.G && sat.G == sat.B {
		t.Fatal("Saturate(gray,100) should produce a non-gray color")
	}
}

func TestDesaturate(t *testing.T) {
	red := Color{Type: ColorTrue, R: 255, G: 0, B: 0}
	desat := Desaturate(red, 100)
	// 100% desaturated red should be gray
	if desat.R != desat.G || desat.G != desat.B {
		t.Fatal("Desaturate(red,100) should produce gray")
	}
}

func TestLuminance(t *testing.T) {
	black := Color{Type: ColorTrue, R: 0, G: 0, B: 0}
	white := Color{Type: ColorTrue, R: 255, G: 255, B: 255}
	if Luminance(black) != 0 {
		t.Fatal("Luminance(black) should be 0")
	}
	if Luminance(white) != 1 {
		t.Fatal("Luminance(white) should be 1")
	}
}

func TestRotateHue(t *testing.T) {
	red := Color{Type: ColorTrue, R: 255, G: 0, B: 0}
	green := RotateHue(red, 120)
	if green.R != 0 || green.G != 255 || green.B != 0 {
		t.Fatalf("RotateHue(red,120) expected (0,255,0), got (%d,%d,%d)", green.R, green.G, green.B)
	}
}

func TestHSLToRGBRoundtrip(t *testing.T) {
	// Test that converting RGB→HSL→RGB preserves the color
	colors := []Color{
		{Type: ColorTrue, R: 255, G: 0, B: 0},
		{Type: ColorTrue, R: 0, G: 255, B: 0},
		{Type: ColorTrue, R: 0, G: 0, B: 255},
		{Type: ColorTrue, R: 128, G: 128, B: 128},
		{Type: ColorTrue, R: 255, G: 255, B: 0},
	}
	for _, c := range colors {
		h, s, l := rgbToHSL(c.R, c.G, c.B)
		back := HSL(h, s, l)
		if absDiff(back.R, c.R) > 2 || absDiff(back.G, c.G) > 2 || absDiff(back.B, c.B) > 2 {
			t.Fatalf("roundtrip RGB→HSL→RGB failed for (%d,%d,%d): got (%d,%d,%d)",
				c.R, c.G, c.B, back.R, back.G, back.B)
		}
	}
}

func absDiff(a, b uint8) int {
	if a > b { return int(a) - int(b) }
	return int(b) - int(a)
}

func TestANSI256PaletteSize(t *testing.T) {
	if len(ANSI256Palette) != 256 {
		t.Fatalf("expected 256 colors, got %d", len(ANSI256Palette))
	}
}

func TestANSI256PaletteFirst(t *testing.T) {
	if ANSI256Palette[0].R != 0 || ANSI256Palette[0].G != 0 || ANSI256Palette[0].B != 0 {
		t.Fatal("ANSI256Palette[0] should be black")
	}
}

func TestANSI256PaletteLast(t *testing.T) {
	if ANSI256Palette[255].R != 238 || ANSI256Palette[255].G != 238 || ANSI256Palette[255].B != 238 {
		t.Fatalf("ANSI256Palette[255] expected (238,238,238), got (%d,%d,%d)",
			ANSI256Palette[255].R, ANSI256Palette[255].G, ANSI256Palette[255].B)
	}
}

func TestTo256TrueColor(t *testing.T) {
	red := Color{Type: ColorTrue, R: 255, G: 0, B: 0}
	idx := To256(red)
	// Red should map to index 9 (bright red) or nearby
	if idx != 9 && idx != 1 && idx != 196 {
		t.Fatalf("To256(red) expected 9, 1, or 196, got %d", idx)
	}
}

func TestTo256Already(t *testing.T) {
	c := Color{Type: Color256, Index: 42}
	if To256(c) != 42 {
		t.Fatal("To256 of Color256 should return its index")
	}
}

func TestToANSI(t *testing.T) {
	red := Color{Type: ColorTrue, R: 255, G: 0, B: 0}
	idx := ToANSI(red)
	if idx != 9 { // bright red
		t.Fatalf("ToANSI(red) expected 9, got %d", idx)
	}
}

func TestToANSIBlack(t *testing.T) {
	black := Color{Type: ColorTrue, R: 0, G: 0, B: 0}
	if ToANSI(black) != 0 {
		t.Fatalf("ToANSI(black) expected 0, got %d", ToANSI(black))
	}
}

func TestToANSIWhite(t *testing.T) {
	white := Color{Type: ColorTrue, R: 255, G: 255, B: 255}
	if ToANSI(white) != 15 {
		t.Fatalf("ToANSI(white) expected 15, got %d", ToANSI(white))
	}
}

func TestToTrue(t *testing.T) {
	c := Color{Type: ColorANSI, Index: 9} // bright red
	trueC := ToTrue(c)
	if trueC.R != 255 || trueC.G != 0 || trueC.B != 0 {
		t.Fatalf("ToTrue(ANSI 9) expected (255,0,0), got (%d,%d,%d)", trueC.R, trueC.G, trueC.B)
	}
}

func TestToTrueAlready(t *testing.T) {
	c := Color{Type: ColorTrue, R: 100, G: 150, B: 200}
	trueC := ToTrue(c)
	if trueC.R != 100 || trueC.G != 150 || trueC.B != 200 {
		t.Fatal("ToTrue of true color should return unchanged")
	}
}

func TestClosestPalette(t *testing.T) {
	palette := []Color{
		{Type: ColorTrue, R: 255, G: 0, B: 0},
		{Type: ColorTrue, R: 0, G: 255, B: 0},
		{Type: ColorTrue, R: 0, G: 0, B: 255},
	}
	red := Color{Type: ColorTrue, R: 200, G: 0, B: 0}
	idx := ClosestPalette(red, palette)
	if idx != 0 {
		t.Fatalf("ClosestPalette(expected 0 for red), got %d", idx)
	}
}

func TestClosestPaletteEmpty(t *testing.T) {
	if ClosestPalette(Color{}, nil) != -1 {
		t.Fatal("ClosestPalette with empty palette should return -1")
	}
}

func TestANSI256Grayscale(t *testing.T) {
	// Index 232 should be very dark gray
	if ANSI256Palette[232].R != 8 {
		t.Fatalf("ANSI256Palette[232].R expected 8, got %d", ANSI256Palette[232].R)
	}
	// Index 255 should be near white (238)
	if ANSI256Palette[255].R != 238 {
		t.Fatalf("ANSI256Palette[255].R expected 238, got %d", ANSI256Palette[255].R)
	}
}

func TestTo256Gray(t *testing.T) {
	gray := Color{Type: ColorTrue, R: 128, G: 128, B: 128}
	idx := To256(gray)
	// Gray should map to grayscale range (232-255)
	if idx < 232 || idx > 255 {
		t.Fatalf("To256(gray) expected 232-255, got %d", idx)
	}
}
