// Package ascii provides ASCII art text generation using FIGlet fonts.
//
// It supports multiple built-in font styles and can load external FIGlet
// font files (.flf). The main entry point is [Generate], which takes a
// string and a font name and returns the ASCII art representation.
//
// Example:
//
//	art, err := ascii.Generate("Hello", "graffiti")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(art)
package ascii

import (
	"fmt"
	"strings"
)

// Generate creates ASCII art text using the named font.
//
// Supported built-in fonts: "graffiti", "standard", "big", "block", "shadow".
// Font names are case-insensitive.
func Generate(text, fontName string) (string, error) {
	font, ok := builtins[strings.ToLower(fontName)]
	if !ok {
		names := make([]string, 0, len(builtins))
		for n := range builtins {
			names = append(names, n)
		}
		return "", fmt.Errorf("ascii: unknown font %q; available: %v", fontName, names)
	}
	return font.Render(text), nil
}

// AvailableFonts returns the names of all built-in fonts.
func AvailableFonts() []string {
	names := make([]string, 0, len(builtins))
	for n := range builtins {
		names = append(names, n)
	}
	return names
}

// Render generates ASCII art text using the given Font.
func (f *Font) Render(text string) string {
	if f == nil || f.Height == 0 || len(text) == 0 {
		return ""
	}

	// Build output lines
	lines := make([]string, f.Height)
	for _, r := range text {
		ch, ok := f.chars[r]
		if !ok {
			// Use space for unknown characters
			ch = f.chars[32]
		}
		for i := 0; i < f.Height; i++ {
			if i < len(ch) {
				lines[i] += ch[i]
			} else {
				lines[i] += strings.Repeat(" ", f.width(r))
			}
		}
	}

	// Trim trailing spaces from each line
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " ")
	}

	return strings.Join(lines, "\n")
}

// width returns the width of a character in the font.
func (f *Font) width(r rune) int {
	ch, ok := f.chars[r]
	if !ok {
		return 0
	}
	max := 0
	for _, line := range ch {
		if len(line) > max {
			max = len(line)
		}
	}
	return max
}

// Must panics if err is non-nil, otherwise returns s.
// Useful for one-liners:
//
//	art := ascii.Must(ascii.Generate("Hello", "graffiti"))
func Must(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}
