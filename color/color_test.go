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
