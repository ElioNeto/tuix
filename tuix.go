// Package tuix provides a Go library for building Terminal User Interfaces
// by writing HTML and CSS. It parses HTML into a DOM tree, applies CSS styles,
// performs layout, and renders to the terminal with full mouse and keyboard support.
//
// Usage:
//
//	app := tuix.New()
//	app.SetHTML(`<div class="hello">Hello, World!</div>`)
//	app.SetCSS(`
//	  .hello { color: green; font-weight: bold; }
//	`)
//	app.Run()
package tuix

import (
	"github.com/elioneto/tuix/css"
	"github.com/elioneto/tuix/dom"
	"github.com/elioneto/tuix/layout"
	"github.com/elioneto/tuix/render"
	"github.com/elioneto/tuix/style"
	"github.com/elioneto/tuix/terminal"
)

// App is the main application type. Create one with New().
type App struct {
	html       string
	css        string
	document   *dom.Node
	stylesheet *css.Stylesheet
	terminal   *terminal.Terminal
	layout     *layout.LayoutEngine
	rootBox    *layout.Box
	canvas     *render.Canvas
	oldCanvas  *render.Canvas
	painter    *render.Painter
	width      int
	height     int
	running    bool

	// Callbacks
	onInit    func()
	onRender  func()
	onEvent   func(terminal.Event)
	onKey     func(terminal.Key)
	onRune    func(rune)
	onResize  func(w, h int)
	onMouse   func(terminal.MouseButton, int, int)
	onClose   func()

	// State
	stylesResolved bool
}

// New creates a new tuix application.
func New() *App {
	return &App{
		layout: layout.NewEngine(),
	}
}

// SetHTML sets the HTML content for the application.
func (a *App) SetHTML(html string) {
	a.html = html
	a.stylesResolved = false
}

// SetCSS sets the CSS stylesheet for the application.
func (a *App) SetCSS(cssContent string) {
	a.css = cssContent
	a.stylesResolved = false
}

// OnInit registers a callback called once when the app starts.
func (a *App) OnInit(fn func()) {
	a.onInit = fn
}

// OnRender registers a callback called before each frame render.
func (a *App) OnRender(fn func()) {
	a.onRender = fn
}

// OnEvent registers a callback for every terminal event.
func (a *App) OnEvent(fn func(terminal.Event)) {
	a.onEvent = fn
}

// OnKey registers a callback for key press events.
func (a *App) OnKey(fn func(terminal.Key)) {
	a.onKey = fn
}

// OnRune registers a callback for character input events.
func (a *App) OnRune(fn func(rune)) {
	a.onRune = fn
}

// OnResize registers a callback for terminal resize events.
func (a *App) OnResize(fn func(w, h int)) {
	a.onResize = fn
}

// OnMouse registers a callback for mouse events.
func (a *App) OnMouse(fn func(terminal.MouseButton, int, int)) {
	a.onMouse = fn
}

// OnClose registers a callback for when the app exits.
func (a *App) OnClose(fn func()) {
	a.onClose = fn
}

// Run starts the application's main loop.
// It opens the terminal in raw mode, parses HTML/CSS, and enters the event loop.
func (a *App) Run() error {
	var err error
	a.terminal, err = terminal.Open()
	if err != nil {
		return err
	}
	defer a.terminal.Close()

	// Get initial terminal size
	a.width, a.height, err = a.terminal.Size()
	if err != nil {
		a.width, a.height = 80, 24 // sensible defaults
	}

	// Update viewport in layout engine
	a.layout.ViewWidth = float64(a.width)
	a.layout.ViewHeight = float64(a.height)

	// Parse HTML
	if a.html != "" {
		parser := dom.NewParser(a.html)
		a.document = parser.Parse()
	} else {
		a.document = dom.Document()
	}

	// Parse CSS
	if a.css != "" {
		cssParser := css.NewParser(a.css)
		sheet, _ := cssParser.Parse()
		a.stylesheet = sheet
	} else {
		a.stylesheet = &css.Stylesheet{}
	}

	// Create canvas
	a.canvas = render.NewCanvas(a.width, a.height, a.terminal.ColorMode())

	// Call init callback
	if a.onInit != nil {
		a.onInit()
	}

	// Initial render
	a.renderFrame()

	// Enter event loop
	a.running = true
	defer func() { a.running = false }()

	for a.running {
		select {
		case event, ok := <-a.terminal.Events():
			if !ok {
				return nil
			}
			a.handleEvent(event)
		}
	}

	if a.onClose != nil {
		a.onClose()
	}

	return nil
}

// Stop terminates the application's main loop.
func (a *App) Stop() {
	a.running = false
}

// Document returns the parsed DOM document.
func (a *App) Document() *dom.Node {
	return a.document
}

// Stylesheet returns the parsed CSS stylesheet.
func (a *App) Stylesheet() *css.Stylesheet {
	return a.stylesheet
}

// RootBox returns the root layout box after layout computation.
func (a *App) RootBox() *layout.Box {
	return a.rootBox
}

// Canvas returns the current render canvas.
func (a *App) Canvas() *render.Canvas {
	return a.canvas
}

// Terminal returns the terminal handle.
func (a *App) Terminal() *terminal.Terminal {
	return a.terminal
}

// Width returns the current terminal width in columns.
func (a *App) Width() int {
	return a.width
}

// Height returns the current terminal height in rows.
func (a *App) Height() int {
	return a.height
}

// SetTitle sets the terminal window title.
func (a *App) SetTitle(title string) {
	if a.terminal != nil {
		a.terminal.WriteString("\x1b]0;" + title + "\x07")
	}
}

// Rebuild forces a full re-render of the application.
func (a *App) Rebuild() {
	a.stylesResolved = false
	a.renderFrame()
}

// handleEvent processes a single terminal event.
func (a *App) handleEvent(event terminal.Event) {
	// Global event callback
	if a.onEvent != nil {
		a.onEvent(event)
	}

	switch event.Type {
	case terminal.EventKey:
		if event.Key == terminal.KeyRune {
			if a.onRune != nil {
				a.onRune(event.Rune)
			}
		}
		if a.onKey != nil {
			a.onKey(event.Key)
		}

		// Handle Ctrl+C / Ctrl+D to quit
		if event.Key == terminal.KeyCtrlC || event.Key == terminal.KeyCtrlD {
			a.Stop()
		}

	case terminal.EventResize:
		a.width = event.Width
		a.height = event.Height
		a.layout.ViewWidth = float64(a.width)
		a.layout.ViewHeight = float64(a.height)
		a.canvas = render.NewCanvas(a.width, a.height, a.terminal.ColorMode())
		a.renderFrame()
		if a.onResize != nil {
			a.onResize(event.Width, event.Height)
		}

	case terminal.EventMouse:
		if a.onMouse != nil {
			a.onMouse(event.MouseButton, event.MouseX, event.MouseY)
		}
	}
}

// renderFrame performs layout and paints the current frame to the terminal.
func (a *App) renderFrame() {
	if a.document == nil {
		return
	}

	// Resolve styles for all nodes
	styles := a.resolveStyles()

	// Perform layout
	a.rootBox = a.layout.Layout(a.document, styles)

	// Clear canvas
	a.canvas.Clear()

	// Paint the layout
	a.painter = render.NewPainter(a.canvas, a.terminal.ColorMode())
	a.painter.Paint(a.rootBox)

	// Call render callback
	if a.onRender != nil {
		a.onRender()
	}

	// Output to terminal
	output := a.canvas.Render(a.oldCanvas)
	a.terminal.WriteString(output)

	// Save canvas for differential rendering
	oldCanvas := a.canvas
	a.oldCanvas = oldCanvas

	// Create a fresh canvas for next frame
	a.canvas = render.NewCanvas(a.width, a.height, a.terminal.ColorMode())
}

// resolveStyles walks the DOM tree and resolves styles for each node.
func (a *App) resolveStyles() map[*dom.Node]style.ComputedStyle {
	resolver := style.NewResolver(a.stylesheet)
	styles := make(map[*dom.Node]style.ComputedStyle)
	a.resolveNodeStyles(a.document, resolver, styles)
	return styles
}

func (a *App) resolveNodeStyles(node *dom.Node, resolver *style.Resolver,
	styles map[*dom.Node]style.ComputedStyle) {
	if node == nil {
		return
	}
	styles[node] = resolver.Resolve(node)
	for _, child := range node.Children {
		a.resolveNodeStyles(child, resolver, styles)
	}
}
