# tuix

**Build Terminal User Interfaces by writing HTML and CSS in Go.**

Tuix is a zero-dependency Go library for creating TUI (Terminal User Interface) applications. Write your UI in familiar HTML and CSS — tuix parses them into a DOM tree and stylesheet, computes layouts, and renders to the terminal with full keyboard, mouse, and resize support.

## Features

- **HTML-based UI** — Define structure with HTML tags, attributes, classes, and IDs
- **CSS styling** — Style with CSS selectors, properties, cascade, and specificity
- **Block layout engine** — Box model with margin, border, padding
- **Flexbox layout** — `display: flex` with `flex-direction`, `flex-wrap`, `justify-content`, `align-items`, `flex-grow/shrink/basis`, `gap`, `order`
- **Inline formatting** — Inline elements (`span`, `a`, `strong`, `em`, etc.) flow within text lines
- **Scrolling & overflow** — `overflow: scroll/auto/hidden` with scrollbar, keyboard and mouse wheel scrolling
- **Interactive form controls** — Text input, checkboxes, radio buttons, textarea, select dropdowns, buttons with Tab navigation and focus
- **ASCII art text generation** — Built-in FIGlet font support: Graffiti, Standard, Big, Block, Shadow (and any `.flf` font you load)
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
├── ascii/
│   ├── ascii.go         # Generate(), Render(), Font type, Must helper
│   ├── font.go          # FIGlet font parser (.flf format)
│   └── data.go          # Built-in fonts (Graffiti, Standard, Big, Block, Shadow)
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

- **Tags**: `div`, `h1`..`h6`, `p`, `span`, `button`, `a`, `ul`, `ol`, `li`, `header`, `footer`, `section`, `main`, `article`, `aside`, `nav`, `label`, `input`, `textarea`, `select`, `option`, `table`, `tr`, `td`, `th`, `img`, `br`, `hr`, `strong`, `em`, `b`, `i`, `u`, `code`, `pre`, `blockquote`, `cite`
- **Attributes**: `id`, `class`, `style`, `type`, `name`, `value`, `placeholder`, `disabled`, `checked`
- **Self-closing tags**: `br`, `hr`, `img`, `input`, `meta`
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
| `outline`, `outline-style` | `none`, `solid` | ❌ |

#### Flexbox

| Property | Values | Applies to |
|---|---|---|
| `flex-direction` | `row`, `column`, `row-reverse`, `column-reverse` | flex container |
| `flex-wrap` | `nowrap`, `wrap`, `wrap-reverse` | flex container |
| `justify-content` | `flex-start`, `flex-end`, `center`, `space-between`, `space-around`, `space-evenly` | flex container |
| `align-items` | `flex-start`, `flex-end`, `center`, `stretch`, `baseline` | flex container |
| `align-content` | `flex-start`, `flex-end`, `center`, `stretch`, `space-between`, `space-around` | flex container (multi-line) |
| `gap`, `row-gap`, `column-gap` | `<length>` | flex container |
| `flex-grow` | `<number>` (default 0) | flex item |
| `flex-shrink` | `<number>` (default 1) | flex item |
| `flex-basis` | `<length>`, `auto` (default `auto`) | flex item |
| `order` | `<integer>` (default 0) | flex item |
| `align-self` | `auto`, `flex-start`, `flex-end`, `center`, `stretch`, `baseline` | flex item |

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
| Pseudo-class | `:hover`, `:focus`, `:focus-visible`, `:focus-within`, `:disabled`, `:enabled`, `:required`, `:optional`, `:read-only`, `:read-write`, `:placeholder-shown`, `:valid`, `:invalid` | Matches dynamic element state |
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

app.OnFocus(func(el *dom.Node) {
    // Called when an element receives focus
})

app.OnBlur(func(el *dom.Node) {
    // Called when an element loses focus
})
```

### Focus Management API

```go
// FocusElement sets focus to a specific DOM node
app.FocusElement(node)

// Blur removes focus from the currently focused element
app.Blur()

// FocusedElement returns the currently focused node (or nil)
focused := app.FocusedElement()
```

#### `tabindex` Attribute

| Value | Behavior |
|-------|----------|
| `tabindex="0"` | Element is focusable and Tab-navigable in DOM order |
| `tabindex="-1"` | Element is programmatically focusable (via `FocusElement()`) but skipped during Tab navigation |
| `tabindex="5"` | Positive values are focusable and Tab-navigable (in numerical order, then DOM order) |
| *(no attribute)* | Native focusable tags (`<input>`, `<button>`, `<select>`, `<textarea>`, `<a>`) are Tab-navigable |

#### `autofocus` Attribute

Add `autofocus` to any focusable element to give it focus when the app starts:

```html
<input type="text" autofocus />
<div tabindex="0" autofocus>Focused on start</div>
```
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

## ASCII Art Generator

The `ascii` package provides FIGlet-compatible ASCII art text generation with built-in fonts.

### Functions

```go
// Generate creates ASCII art text using the named font.
// Supported fonts: "graffiti", "standard", "big", "block", "shadow"
// Font names are case-insensitive.
art, err := ascii.Generate("Hello", "graffiti")

// AvailableFonts returns the list of built-in font names.
fonts := ascii.AvailableFonts() // ["big", "block", "graffiti", "shadow", "standard"]

// Must panics on error — useful for one-liners.
art := ascii.Must(ascii.Generate("TUIX", "big"))
```

### Built-in Fonts

| Font      | Description                    | Height | Source |
|-----------|--------------------------------|--------|--------|
| Graffiti  | Spray-painted street art style | 6      | FIGlet |
| Standard  | Classic FIGlet font            | 6      | FIGlet |
| Big       | Large, bold lettering          | 8      | FIGlet |
| Block     | Clean, geometric block letters | 8      | FIGlet |
| Shadow    | Letters with 3D shadow effect  | 5      | FIGlet |

### Loading Custom Fonts

Load any FIGlet `.flf` font file:

```go
data, _ := os.ReadFile("path/to/font.flf")
font, err := ascii.ParseFIGlet("myfont", string(data))
art := font.Render("Hello")
```

### Example

```go
package main

import (
    "fmt"
    "github.com/elioneto/tuix/ascii"
)

func main() {
    art, _ := ascii.Generate("Hello!", "graffiti")
    fmt.Println(art)
}
```

Output:

```
  ___ ___          .__   .__           ._.
 /   |   \   ____  |  |  |  |    ____  | |
/    ~    \_/ __ \ |  |  |  |   /  _ \ | |
\    Y    /\  ___/ |  |__|  |__(  <_> ) \|
 \___|_  /  \___  > |____|____/ \____/  \|
        \/       \/
```

### Image-to-ASCII Conversion

Convert images (PNG, JPEG, GIF) to ASCII art with optional dithering and true-color output:

```go
// Basic grayscale ASCII from file
art, err := ascii.FromFile("photo.png", ascii.DefaultImageOptions())

// With custom options
art := ascii.FromImage(img, ascii.ImageOptions{
    Width:   80,
    Height:  0,        // 0 = auto (preserves aspect ratio)
    Charset: ascii.CharsetBlock, // or Standard, Simple, Detailed
    Color:   true,     // true-color ANSI foreground per pixel
    Dither:  true,     // Floyd-Steinberg error diffusion
    Scale:   1.0,
})

// Animated GIF support
frames, err := ascii.FromFileGIF("animation.gif", opts)
for _, frame := range frames {
    fmt.Print(frame)
    time.Sleep(100 * time.Millisecond)
}
```

**Built-in character sets:**

| Constant | Levels | Description |
|----------|--------|-------------|
| `CharsetStandard` | 10 | `@@%#*+=-:. ` — balanced ramp |
| `CharsetSimple` | 7 | `@%#*+=-.` — bold, compact |
| `CharsetBlock` | 4 | ` ░▒▓█` — block elements |
| `CharsetDetailed` | 16 | Full detailed ramp for fine gradients |
  ___ ___          .__   .__           ._.
 /   |   \   ____  |  |  |  |    ____  | |
/    ~    \_/ __ \ |  |  |  |   /  _ \ | |
\    Y    /\  ___/ |  |__|  |__(  <_> ) \|
 \___|_  /  \___  >|____/|____/ \____/  __
       \/       \/                      \/
```

---

## Examples

All examples are in the `examples/` directory:

| Example | Description | Run command |
|---------|-------------|-------------|
| **Counter** | Interactive counter with buttons, mouse click, and keyboard shortcuts | `go run ./examples/counter/` |
| **Flexbox** | Flexbox layout demo with centering, wrapping, gap, and ordering | `go run ./examples/flexbox/` |
| **Inline** | Inline text formatting with `span`, `strong`, `em`, nested elements | `go run ./examples/inline/` |
| **Scrolling** | Scrollable content areas with scrollbars, keyboard & mouse wheel | `go run ./examples/scrolling/` |
| **Forms** | Full form demo — text, search, number, range, color, date inputs, checkbox, radio, select dropdown, textarea, buttons, progress bar, meter gauge, datalist autocomplete | `go run ./examples/forms/` |
| **ASCII** | ASCII art text generation with 5 built-in FIGlet fonts | `go run ./examples/ascii/` |
| **ASCII Image** | Image-to-ASCII conversion — convert PNG/JPEG/GIF to ASCII art with dithering and color | `go run ./examples/ascii-image/` |
| **Focus** | Focus management with Tab/Shift+Tab, focus ring, auto-focus, callbacks, tabindex | `go run ./examples/focus/` |
| **Hover** | Hover effects with `:hover` pseudo-class, real-time mouse tracking, enter/leave callbacks | `go run ./examples/hover/` |
| **Modal** | Modal dialog overlay with focus trap, Escape-to-close, backdrop | `go run ./examples/modal/` |
| **Toast** | Toast notifications with 4 types (info/success/warning/error), auto-dismiss, stacking | `go run ./examples/toast/` |
| **Z-Index** | Stacking order with z-index for overlapping elements | `go run ./examples/zindex/` |
| **Design System** | Pre-built components (buttons, badges, cards, navbar, tabs, lists, tables) with interactive theme switcher | `go run ./examples/design-system/` |

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

## Design System

Tuix includes a built-in design system with pre-built CSS components and theme support.

### Using the Design System

```go
app := tuix.New()
app.UseDesignSystem() // Enables all pre-built components
app.Run()
```

`UseDesignSystem()` automatically applies `DefaultDarkTheme` and the built-in component CSS.

### Theme API

```go
import "github.com/elioneto/tuix"

// Built-in themes
app.SetTheme(tuix.DefaultDarkTheme)
app.SetTheme(tuix.DefaultLightTheme)

// Custom theme
app.SetTheme(tuix.Theme{
    Primary:   color.Color{R: 0x00, G: 0xD4, B: 0xAA},
    Secondary: color.Color{R: 0x0F, G: 0x34, B: 0x60},
    Accent:    color.Color{R: 0xE9, G: 0x45, B: 0x60},
    Success:   color.Color{R: 0x2E, G: 0xCC, B: 0x71},
    Warning:   color.Color{R: 0xF3, G: 0x9C, B: 0x12},
    Error:     color.Color{R: 0xE7, G: 0x4C, B: 0x3C},
    Surface:   color.Color{R: 0x16, G: 0x21, B: 0x3E},
    Background: color.Color{R: 0x1A, G: 0x1A, B: 0x2E},
    Text:      color.Color{R: 0xC0, G: 0xC0, B: 0xC0},
    Muted:     color.Color{R: 0x55, G: 0x55, B: 0x55},
    Border:    color.Color{R: 0x0F, G: 0x34, B: 0x60},
    Focus:     color.Color{R: 0x00, G: 0xD4, B: 0xAA},
})
```

Theme colors are applied as CSS utility classes — `.bg-primary`, `.text-muted`, `.border-error`, etc. — and can be used alongside regular CSS.

### Pre-built Components

| Class | Description |
|-------|-------------|
| `.btn` | Base button with border and padding |
| `.btn-primary` | Primary action button (filled) |
| `.btn-secondary` | Secondary action button (outlined) |
| `.btn-danger` | Destructive action button |
| `.btn-ghost` | Ghost button (no border/background) |
| `.btn-sm` / `.btn-lg` | Small / Large button sizes |
| `.input` | Text input with consistent styling |
| `.input-error` | Input in error state |
| `.input-sm` / `.input-lg` | Small / Large input sizes |
| `.badge` | Badge/tag base |
| `.badge-primary` / `.badge-success` / `.badge-warning` / `.badge-error` | Colored badges |
| `.card` | Card container with border and padding |
| `.navbar` | Top navigation bar with brand and items |
| `.nav-brand` | Brand/logo text in navbar |
| `.nav-item` | Navigation item in navbar |
| `.list` | List group container |
| `.list-item` | List group item |
| `.tabs` | Tab navigation container |
| `.tab` | Individual tab item |
| `.tab-active` | Active tab state |
| `.table` | Table container |
| `.table-header` | Table header row |
| `.table-row` | Table data row |

### Layout Utilities

| Class | Description |
|-------|-------------|
| `.flex` | `display: flex` |
| `.flex-col` | `flex-direction: column` |
| `.flex-wrap` | `flex-wrap: wrap` |
| `.items-center` | `align-items: center` |
| `.justify-center` | `justify-content: center` |
| `.justify-between` | `justify-content: space-between` |
| `.gap-1` / `.gap-2` / `.gap-4` | `gap` spacing |
| `.w-full` | `width: 100%` |
| `.text-center` | `text-align: center` |
| `.text-bold` | `font-weight: bold` |
| `.grid-2` / `.grid-3` | Equal-width grid columns |

---

## Roadmap

### ✅ Completed

- [x] **Block layout** — Box model with margin, border, padding, word-wrap
- [x] **Flexbox layout** — `display: flex` with full property support
- [x] **Inline formatting** — Text and inline elements flowing within lines
- [x] **Scrolling / overflow** — Scrollable containers with scrollbars and mouse/keyboard input
- [x] **Input fields & form controls** — Text input, checkbox, radio, select, textarea, button with Tab navigation
- [x] **ASCII art text generator** — FIGlet font support (Graffiti, Standard, Big, Block, Shadow) via `ascii` package

### 📋 Upcoming

- [x] **Focus management** — `:focus` pseudo-class, focus ring, `tabindex`, auto-focus, OnFocus/OnBlur callbacks
- [x] **Hover effects** — `:hover` pseudo-class trigger, mouse enter/leave events
- [x] **ASCII art image converter** — PNG/JPEG/GIF → ASCII art with dithering, true-color output, GIF frame support
- [ ] **CSS animations & transitions** — Animated property changes
- [x] **Modal / Dialog** — Overlay modal with backdrop, focus trap, `Esc` to close
- [x] **Alert / Toast notifications** — Non-blocking notification popups with auto-dismiss, alert/confirm dialogs
- [x] **Z-index / stacking contexts** — Proper layering of overlapping elements
- [x] **Enhanced form controls** — Search, number, range, color, date inputs, pseudo-classes (:disabled/:enabled/:required/:optional/:read-only/:read-write/:placeholder-shown), progress/meter elements, tooltips, datalist autocomplete, select dropdown
- [x] **Design System** — Theme engine (light/dark), pre-built components (buttons, badges, cards, navbar, tabs, lists, tables), layout utilities, `UseDesignSystem()` API
- [ ] **Image rendering** — Sixel and Kitty image protocols
- [ ] **Comprehensive test suite** — Unit and integration tests
- [ ] **Windows native support** — Windows console API fallback
- [ ] **CSS animations & transitions** — Animated property changes

---

## License

MIT © Elio Neto
