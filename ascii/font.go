package ascii

import (
	"strings"
)

// Font represents a parsed FIGlet font.
type Font struct {
	Name      string
	Height    int
	Baseline  int
	HardBlank byte
	chars     map[rune][]string
}

// ParseFIGlet parses a FIGlet font (.flf) string and returns the Font.
// The font name is used for identification only.
func ParseFIGlet(name, data string) (*Font, error) {
	// Normalize line endings: replace \r\n and \r with \n
	data = strings.ReplaceAll(data, "\r\n", "\n")
	data = strings.ReplaceAll(data, "\r", "\n")

	lines := strings.Split(data, "\n")
	if len(lines) < 2 {
		return nil, ErrMalformed("font data too short")
	}

	// Parse first line: flf2a<hardblank> <height> <baseline> <max_length> <old_layout> <comment_lines> ...
	first := lines[0]
	if len(first) < 6 || first[:5] != "flf2a" {
		return nil, ErrMalformed("missing flf2a magic")
	}

	hardBlank := first[5]

	// Parse header fields
	fields := strings.Fields(first[6:])
	if len(fields) < 5 {
		return nil, ErrMalformed("insufficient header fields")
	}

	f := &Font{
		Name:      name,
		HardBlank: hardBlank,
		chars:     make(map[rune][]string),
	}

	if err := f.parseIntField(fields[0], "height", &f.Height); err != nil {
		return nil, err
	}
	if err := f.parseIntField(fields[1], "baseline", &f.Baseline); err != nil {
		return nil, err
	}
	// max_length (fields[2]), old_layout (fields[3]), full_layout etc. ignored for rendering
	var commentLines int
	if err := f.parseIntField(fields[4], "comment_lines", &commentLines); err != nil {
		return nil, err
	}

	// Skip comment lines
	pos := 1
	for pos < len(lines) && commentLines > 0 {
		pos++
		commentLines--
	}

	if pos >= len(lines) {
		return nil, ErrMalformed("no character data after comments")
	}

	// Parse character definitions
	// Characters are in order starting from ASCII 32 (space)
	// Each character occupies f.Height lines
	// Lines end with @ (EOL marker); last line of character ends with @@ (EOL + separator)
	code := 32
	for pos < len(lines) {
		if pos+f.Height > len(lines) {
			break
		}

		charLines := make([]string, f.Height)
		maxW := 0

		for i := 0; i < f.Height; i++ {
			raw := lines[pos]
			pos++

			// Remove EOL marker(s):
			// Content is everything before the last '@' or '@@'
			cleaned := trimEOL(raw)

			// Replace hardblank character with space
			if f.HardBlank != 0 {
				cleaned = strings.ReplaceAll(cleaned, string(f.HardBlank), " ")
			}

			charLines[i] = cleaned
			if len(cleaned) > maxW {
				maxW = len(cleaned)
			}
		}

		// Normalize all lines to same width
		for i := range charLines {
			if len(charLines[i]) < maxW {
				charLines[i] += strings.Repeat(" ", maxW-len(charLines[i]))
			}
		}

		f.chars[rune(code)] = charLines
		code++

		// Stop at ASCII 126 (tilde) unless font has more
		if code > 126 {
			// Check for extended characters
			for pos < len(lines) {
				// Extended chars: each char still has f.Height lines
				if pos+f.Height > len(lines) {
					break
				}
				extLines := make([]string, f.Height)
				extMaxW := 0
				for i := 0; i < f.Height; i++ {
					raw := lines[pos]
					pos++
					cleaned := trimEOL(raw)
					if f.HardBlank != 0 {
						cleaned = strings.ReplaceAll(cleaned, string(f.HardBlank), " ")
					}
					extLines[i] = cleaned
					if len(cleaned) > extMaxW {
						extMaxW = len(cleaned)
					}
				}
				for i := range extLines {
					if len(extLines[i]) < extMaxW {
						extLines[i] += strings.Repeat(" ", extMaxW-len(extLines[i]))
					}
				}
				f.chars[rune(code)] = extLines
				code++
			}
		}
	}

	return f, nil
}

// trimEOL removes the trailing @ or @@ EOL marker(s) from a character line.
// The format is: content + "@" (normal line) or content + "@@" (last line of char).
func trimEOL(line string) string {
	// Check for @@ first (end of character marker)
	if strings.HasSuffix(line, "@") {
		// Check for double-@
		if len(line) >= 2 && line[len(line)-2] == '@' && line[len(line)-1] == '@' {
			return line[:len(line)-2]
		}
		// Single @ (EOL marker)
		return line[:len(line)-1]
	}
	// Some fonts may not use @; return as-is
	return line
}

func (f *Font) parseIntField(field, name string, out *int) error {
	n := 0
	for _, c := range field {
		if c < '0' || c > '9' {
			return ErrMalformed("non-numeric " + name + ": " + field)
		}
		n = n*10 + int(c-'0')
	}
	*out = n
	return nil
}

// ErrMalformed is returned when font data is invalid.
type ErrMalformed string

func (e ErrMalformed) Error() string { return "ascii: malformed font: " + string(e) }

// builtins holds the built-in fonts indexed by name.
var builtins = map[string]*Font{}

// registerFont adds a font to the built-in registry.
func registerFont(name, data string) error {
	f, err := ParseFIGlet(name, data)
	if err != nil {
		return err
	}
	builtins[strings.ToLower(name)] = f
	return nil
}

func init() {
	// fonts registered in data.go
}
