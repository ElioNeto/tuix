# tuix

**Build Terminal User Interfaces by writing HTML and CSS in Go.**

Tuix is a zero-dependency Go library for creating TUI applications. Write your UI in familiar HTML and CSS, and tuix handles the rest — parsing, styling, layout, and rendering to the terminal with full keyboard and mouse support.

## Features

- **HTML parsing** — Write UI structure in HTML (tags, attributes, classes, IDs)
- **CSS styling** — Style your UI with CSS (selectors, properties, cascade, specificity)
- **Layout engine** — Block layout with box model (margin, border, padding)
- **Terminal rendering** — ANSI escape code output with 24-bit color support
- **Input handling** — Keyboard, mouse, and resize events
- **Zero dependencies** — Built entirely on Go's standard library
- **Canvas-based rendering** — Pixel-perfect character grid with differential updates

## Quick Start

```go
package main

import (
	"log"
	"github.com/elioneto/tuix"
	"github.com/elioneto/tuix/terminal"
)

func main() {
	app := tuix.New()

	app.SetHTML(`
		<div id="app">
			<h1>Hello, tuix!</h1>
			<p>Press q to quit.</p>
		</div>
	`)

	app.SetCSS(`
		#app {
			padding: 10px;
			background-color: #1a1a2e;
			color: #e0e0e0;
		}
		h1 { color: #00d4aa; text-align: center; }
		p  { text-align: center; color: #888; }
	`)

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	log.Fatal(app.Run())
}
```

## Architecture

```
tuix/
├── tuix.go          # Public API - App struct, event loop
├── dom/             # HTML parser & DOM tree
│   ├── dom.go       # Node types, parser, querySelectorAll
├── css/             # CSS parser & selector engine
│   ├── css.go       # Stylesheet, rules, declarations, values
├── style/           # Computed style resolution
│   ├── style.go     # Resolver, cascade, property application
├── layout/          # Box model & layout engine
│   ├── layout.go    # Block/inline layout, box tree
├── render/          # Terminal canvas & painting
│   ├── render.go    # Canvas, painter, borders, text
├── terminal/        # Raw terminal I/O
│   ├── terminal.go  # Raw mode, event parsing, ANSI output
│   ├── term_unix.go # Unix-specific syscalls (ioctl, termios)
├── color/           # Color parsing & ANSI conversion
│   ├── color.go     # Hex, RGB, named, ANSI, 256-color
└── geometry/        # Geometric primitives
    ├── rect.go      # Rectangle operations
    ├── size.go      # 2D dimensions
    ├── point.go     # 2D coordinates
    └── edges.go     # Box model edges (margin, border, padding)
```

## Requirements

- Go 1.22+
- Linux, macOS, or other Unix-like OS (terminal in raw mode)

## License

MIT
