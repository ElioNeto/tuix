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
