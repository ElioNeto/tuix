# tuix

**Build Terminal User Interfaces by writing HTML and CSS in Go.**

Tuix is a zero-dependency Go library for creating TUI (Terminal User Interface) applications. Write your UI in familiar HTML and CSS — tuix parses them into a DOM tree and stylesheet, computes layouts, and renders to the terminal with full keyboard, mouse, and resize support.

## Features

- **HTML-based UI** — Define structure with HTML tags, attributes, classes, and IDs
- **CSS styling** — Style with CSS selectors, properties, cascade, and specificity
- **Block layout engine** — Box model with margin, border, padding (inline/flex coming soon)
- **24-bit true color** — Hex, RGB, named, ANSI 16/256 color support
- **Keyboard & mouse input** — Full event handling with modifiers
- **Text alignment** — Left, center, right
- **CSS inheritance** — Inherited properties (`color`, `font-weight`, `text-align`, etc.) propagate from parent to child
- **Compound selectors** — `.class1.class2`, `div#id`, `div.class` all work
- **Alternate screen buffer** — Clean enter/exit without cluttering the terminal history
- **Resize handling** — Automatically re-layouts on terminal resize
- **Zero external dependencies** — Built entirely on Go's standard library

## Installation

```bash
go get github.com/elioneto/tuix
```

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
			padding: 1;
			background-color: #1a1a2e;
			color: #e0e0e0;
		}
		h1 {
			color: #00d4aa;
			text-align: center;
			font-weight: bold;
		}
		p {
			text-align: center;
			color: #888;
		}
	`)

	app.OnRune(func(r rune) {
		if r == 'q' {
			app.Stop()
		}
	})

	log.Fatal(app.Run())
}
```

Run it:

```bash
go run main.go
```

Press `q` to quit. The terminal will enter the alternate screen buffer and restore on exit.

---

## Table of Contents

- [Architecture](#architecture)
- [HTML Guide](#html-guide)
- [CSS Guide](#css-guide)
  - [Supported Properties](#supported-properties)
  - [Selectors](#selectors)
  - [Colors](#colors)
  - [Values & Units](#values--units)
  - [Inheritance & Cascade](#inheritance--cascade)
- [Events](#events)
- [API Reference](#api-reference)
- [Examples](#examples)
- [Terminal Compatibility](#terminal-compatibility)
- [Roadmap](#roadmap)

---

## Architecture

```
tuix/
├── tuix.go              # Public API — App struct, event loop, callbacks
├── dom/
│   └── dom.go           # HTML parser, DOM tree (Node, Element, Text)
├── css/
│   └── css.go           # CSS parser (selectors, rules, declarations, values)
├── style/
│   └── style.go         # Style resolver (cascade, specificity, inheritance)
├── layout/
│   └── layout.go        # Layout engine (block formatting, box model)
├── render/
│   └── render.go        # Canvas, painter, ANSI output, borders, text
├── terminal/
│   ├── terminal.go      # Raw mode, input parsing, ANSI helpers
│   └── term_unix.go     # Unix syscalls (termios, ioctl)
├── color/
│   └── color.go         # Color types, parsing, ANSI sequence generation
└── geometry/
    ├── rect.go           # Rectangle with intersection, containment
    ├── point.go          # 2D point
    ├── size.go           # 2D dimensions
    └── edges.go          # Box model edges (top, right, bottom, left)
```

### Data Flow

```
HTML String  ──▶  DOM Parser  ──▶  DOM Tree
                                         │
CSS String   ──▶  CSS Parser  ──▶  Stylesheet
                                         │
                                 Style Resolver (cascade + inheritance)
                                         │
                                 Layout Engine (box model)
                                         │
                                 Painter → Canvas
                                         │
                                 Terminal Output (ANSI)
```

---

## HTML Guide

Tuix parses a subset of HTML. Tags, attributes, classes, and IDs work as expected.

```html
<div id="app" class="container">
    <h1>Title</h1>
    <p class="highlight">Content text</p>
    <button id="submit">Click</button>
</div>
```

- **Tags**: `div`, `h1`..`h6`, `p`, `span`, `button`, `a`, `ul`, `li`, `header`, `footer`, `section`, `main`, `article`, `aside`, `nav`, `label`, `input`, `textarea`, `select`, `option`, `table`, `tr`, `td`, `th`, `img`, `br`, `hr`, `strong`, `em`, `b`, `i`, `u`, `code`, `pre`, `blockquote`, `cite`
- **Attributes**: `id`, `class`, `style` (not yet processed)
- **Self-closing tags**: `br`, `hr`, `img`, `input`
- **Text nodes**: Content between tags is parsed as text

---

## CSS Guide

### Supported Properties

#### Box Model

| Property | Values | Example |
|---|---|---|
| `width` | `<length>`, `<percentage>`, `auto` | `width: 100%` |
| `height` | `<length>`, `<percentage>`, `auto` | `height: 10` |
| `min-width`, `min-height` | `<length>` | `min-width: 5` |
| `max-width`, `max-height` | `<length>`, `none` | `max-width: 80` |
| `margin` | 1–4 `<length>` values | `margin: 2` or `margin: 1 2` |
| `margin-top`, `margin-right`, `margin-bottom`, `margin-left` | `<length>` | `margin-top: 1` |
| `padding` | 1–4 `<length>` values | `padding: 2` or `padding: 0 1` |
| `padding-top`, `padding-right`, `padding-bottom`, `padding-left` | `<length>` | `padding-left: 1` |
| `border` | `<width> <style> <color>` | `border: solid` or `border: 1px solid red` |
| `border-top`, `border-right`, `border-bottom`, `border-left` | `<width> <style> <color>` | `border-bottom: solid` |
| `border-width` | 1–4 `<length>` | `border-width: 1` |
| `border-style` | `none`, `solid`, `dashed`, `dotted`, `double` | `border-style: solid` |
| `border-color` | 1–4 `<color>` | `border-color: #e94560` |

#### Typography & Color

| Property | Values | Inherited |
|---|---|---|
| `color` | `<color>` | ✅ |
| `background`, `background-color` | `<color>` | ❌ |
| `font-size` | `<number>`, `<length>` | ✅ |
| `font-weight` | `normal`, `bold`, 100–900 | ✅ |
| `text-align` | `left`, `center`, `right`, `justify` | ✅ |
| `line-height` | `<length>` | ✅ |
| `white-space` | `normal`, `nowrap`, `pre`, `pre-wrap`, `pre-line` | ✅ |

#### Layout & Display

| Property | Values | Inherited |
|---|---|---|
| `display` | `block`, `inline`, `inline-block`, `none`, `flex` | ❌ |
| `position` | `static`, `relative`, `absolute`, `fixed` | ❌ |
| `overflow`, `overflow-x`, `overflow-y` | `visible`, `hidden`, `scroll`, `auto` | ❌ |
| `opacity` | `<number>` (0–1) | ❌ |
| `z-index` | `<integer>` | ❌ |
| `visibility` | `visible`, `hidden`, `collapse` | ✅ |
| `cursor` | `auto`, `default`, `pointer`, `text`, `none`, `help` | ✅ |

> **Note:** Units are in **character cells** (columns × rows), not CSS pixels.  
> `margin: 2` means 2 character cells of margin.  
> `font-size: 32` makes line spacing 2 rows tall.
> This is because terminal characters are monospaced and each cell is one unit.

### Selectors

| Selector | Example | Description |
|---|---|---|
| Universal | `*` | Matches any element |
| Tag | `h1` | Matches elements by tag name |
| Class | `.btn` | Matches elements by class |
| ID | `#app` | Matches element by ID |
| Attribute | `[disabled]`, `[type=text]` | Matches by attribute (presence or value) |
| Compound | `.btn.primary` | Matches element matching ALL conditions |
| Descendant | `div p` | Matches `p` inside `div` |
| Child | `div > p` | Matches direct child |
| Adjacent | `h1 + p` | Matches `p` immediately after `h1` |
| Pseudo-class | `:hover` | Parsed but not yet active |
| Comma list | `h1, h2, h3` | Multiple selectors share the same declarations |

### Colors

| Format | Example |
|---|---|
| Named | `red`, `blue`, `green`, `orange`, `coral`, `crimson`… |
| Hex 3-digit | `#f00` (→ `#ff0000`) |
| Hex 6-digit | `#00d4aa` |
| Hex 8-digit | `#c91c9eff` (RRGGBBAA, alpha is parsed) |
| RGB | `rgb(255, 0, 0)` |
| ANSI | `ansi(1)` (0–15) |
| 256-color | `color(196)` (0–255) |
| Transparent | `transparent` |

Font-weight `bold` renders as ANSI bold (`\x1b[1m`).

### Values & Units

| Unit | Example | Description |
|---|---|---|
| (bare number) | `16` | Character cells |
| `px` | `2px` | Character cells (1px = 1 cell) |
| `%` | `100%` | Percentage of parent's content width |
| `em`, `rem` | `2em` | Relative to base font size (16 cells) |
| `auto` | `auto` | Default fill behavior |

### Inheritance & Cascade

Tuix implements the **CSS cascade**:

1. **Specificity** — More specific selectors override less specific ones
   - Inline style (not yet implemented) > ID > Class/Attribute/Pseudo-class > Element
2. **Source order** — When specificity is equal, the rule declared **later** in the stylesheet wins
3. **Inheritance** — Properties like `color`, `font-weight`, `text-align`, `font-size`, `visibility`, `cursor`, `white-space` are inherited from parent elements

```css
/* This has specificity (0,0,1,0) */
.btn { color: #00d4aa; }

/* This has specificity (0,0,2,0) — higher, so it overrides */
.btn.secondary { color: #e94560; }
```

---

## Events

### App Callbacks

```go
app.OnInit(func() {
    // Called once when the app starts (after terminal setup)
})

app.OnRender(func() {
    // Called before each frame render
})

app.OnClose(func() {
    // Called when the app is about to exit
})

app.OnEvent(func(ev terminal.Event) {
    // Called for every event
})

app.OnKey(func(key terminal.Key) {
    // Called on key press
})

app.OnRune(func(r rune) {
    // Called on character input
})

app.OnResize(func(w, h int) {
    // Called when terminal is resized
})

app.OnMouse(func(btn terminal.MouseButton, x, y int) {
    // Called on mouse events
})
```

### Key Constants

```go
terminal.KeyUp / KeyDown / KeyLeft / KeyRight
terminal.KeyEnter / KeyEscape / KeyBackspace / KeyTab
terminal.KeyHome / KeyEnd / KeyPageUp / KeyPageDown
terminal.KeyDelete / KeyInsert
terminal.KeyF1  …  KeyF12
terminal.KeyCtrlA … KeyCtrlZ
```

### Modifiers

```go
terminal.ModShift
terminal.ModAlt
terminal.ModCtrl
```

### Mouse Buttons

```go
terminal.MouseLeft
terminal.MouseMiddle
terminal.MouseRight
terminal.MouseWheelUp
terminal.MouseWheelDown
```

---

## API Reference

### `tuix.App`

| Method | Description |
|---|---|
| `New() *App` | Creates a new tuix application |
| `SetHTML(html string)` | Sets the HTML content |
| `SetCSS(css string)` | Sets the CSS stylesheet |
| `Run() error` | Starts the app (opens terminal, enters alt screen, event loop) |
| `Stop()` | Stops the app and exits |
| `Rebuild()` | Forces a full re-render (call after DOM changes) |
| `SetTitle(title string)` | Sets the terminal window title |
| `Width() int` | Current terminal width (columns) |
| `Height() int` | Current terminal height (rows) |
| `Document() *dom.Node` | The parsed DOM document |
| `Stylesheet() *css.Stylesheet` | The parsed CSS stylesheet |
| `RootBox() *layout.Box` | The root layout box after layout computation |
| `Terminal() *terminal.Terminal` | The underlying terminal handle |

### `dom.Node`

| Method | Description |
|---|---|
| `TagName() string` | Element tag name |
| `ID() string` | `id` attribute |
| `HasClass(name string) bool` | Checks class membership |
| `GetAttribute(name string) string` | Returns attribute value |
| `SetAttribute(name, value string)` | Sets an attribute |
| `TextContent() string` | Concatenated text of all descendants |
| `AppendChild(child *Node)` | Adds a child node |
| `QuerySelectorAll(sel string) []*Node` | Finds matching descendants |

### `terminal.Event`

```go
type Event struct {
    Type        EventType      // EventKey, EventMouse, EventResize
    Key         Key            // For keyboard events
    Rune        rune           // For character input
    Modifiers   Modifier       // Shift, Alt, Ctrl
    MouseButton MouseButton    // For mouse events
    MouseX, MouseY int         // Cell coordinates
    Width, Height int          // For resize events
}
```

---

## Examples

### Counter (interactive)

Full interactive example with buttons, click handling, and keyboard shortcuts:

```bash
go run ./examples/counter/
```

Features:
- `+1` / `-1` buttons via mouse click
- Keyboard shortcuts: `i` / `d` to increment/decrement, `r` to reset
- `q` or `Ctrl+C` to quit
- Real-time DOM manipulation and re-render

---

## Terminal Compatibility

Tuix uses standard ANSI escape codes and should work in any modern terminal emulator:

| Terminal | Status |
|---|---|
| GNOME Terminal (PopOS, Ubuntu) | ✅ Fully supported |
| xterm / xterm-256color | ✅ |
| Alacritty | ✅ |
| Kitty | ✅ |
| WezTerm | ✅ |
| Windows Terminal | ✅ (via WSL) |
| iTerm2 (macOS) | ✅ |
| tmux / screen | ✅ |

### Requirements

- Go 1.22+
- Linux, macOS, or other Unix-like OS (uses `termios` syscalls)
- Terminal must support ANSI escape codes (all modern terminals do)

---

## Roadmap

- [ ] **Flexbox layout** — Support for `display: flex`
- [ ] **Inline formatting** — Text flowing horizontally within a line
- [ ] **Scrolling / overflow** — Scrollable content areas
- [ ] **Input fields** — Text input, form controls
- [ ] **Focus management** — Tab order, `:focus` pseudo-class
- [ ] **Hover effects** — `:hover` pseudo-class trigger
- [ ] **Animation** — CSS transitions and animations
- [ ] **Images** — Image rendering via sixel or kitty protocol
- [ ] **Test suite** — Comprehensive unit and integration tests
- [ ] **Windows native support** — Windows console API fallback

---

## License

MIT © Elio Neto
