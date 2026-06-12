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
	"fmt"
	"strings"
	"time"

	"github.com/elioneto/tuix/css"
	"github.com/elioneto/tuix/dom"
	"github.com/elioneto/tuix/layout"
	"github.com/elioneto/tuix/render"
	"github.com/elioneto/tuix/style"
	"github.com/elioneto/tuix/terminal"
)

// ToastType represents the type of a toast notification.
type ToastType int

const (
	ToastInfo    ToastType = iota // Informational toast (blue)
	ToastSuccess                  // Success toast (green)
	ToastWarning                  // Warning toast (yellow)
	ToastError                    // Error toast (red)
)

// ToastEntry represents a single toast notification in the queue.
type ToastEntry struct {
	ID       int        // Unique toast identifier
	Message  string     // Toast message text
	Type     ToastType  // Toast type (info/success/warning/error)
	Duration time.Duration // Auto-dismiss duration (0 = manual dismiss only)
	Position string     // Position hint: "top-right", "bottom-right", "top-left", "bottom-left"
}

// toastEntry is the internal mutable state for a toast.
type toastEntry struct {
	ToastEntry
}

// toastDefaultCSS provides default styling for toast notifications.
// Users can override any of these in their own CSS.
const toastDefaultCSS = `
.toast-container {
	position: fixed;
	z-index: 200;
	display: flex;
	flex-direction: column;
	gap: 1;
	pointer-events: none;
}
.toast-container-top-right {
	top: 1; right: 1;
}
.toast-container-top-left {
	top: 1; left: 1;
}
.toast-container-bottom-right {
	bottom: 1; right: 1;
}
.toast-container-bottom-left {
	bottom: 1; left: 1;
}
.toast {
	padding: 1 2;
	border: 1;
	min-width: 20;
	max-width: 40;
	background: #1a1a2e;
	color: #e0e0e0;
	border-color: #333;
}
.toast-info    { border-left: 2; border-left-color: #3498db; }
.toast-success { border-left: 2; border-left-color: #2ecc71; }
.toast-warning { border-left: 2; border-left-color: #f39c12; }
.toast-error   { border-left: 2; border-left-color: #e74c3c; }
`

var (
	// Tags that are focusable form elements
	focusableTags = map[string]bool{
		"input":    true,
		"textarea": true,
		"select":   true,
		"button":   true,
		"a":        true,
	}
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

	// Scroll state saved across re-renders (keyed by *dom.Node)
	savedScroll map[*dom.Node]struct{ X, Y int }

	// Scroll state
	scrollContainer *layout.Box // The currently focused scrollable container

	// Form control state
	formValues  map[*dom.Node]string // Current value/text for each input/textarea/select
	formCursors map[*dom.Node]int    // Cursor position for inputs/textarea
	formChecked map[*dom.Node]bool   // Checked state for checkboxes/radios

	// Focus management
	formFocusables []*dom.Node // Ordered list of focusable DOM nodes
	formFocused    int         // Index into formFocusables, or -1 if none focused

	// Focus callbacks
	onFocus func(el *dom.Node) // Called when an element receives focus
	onBlur  func(el *dom.Node) // Called when an element loses focus

	// Hover state
	hoveredNode *dom.Node // The DOM node currently under the mouse cursor
	mouseX      int        // Last known mouse X position
	mouseY      int        // Last known mouse Y position

	// Modal state
	modalHTML   string     // Raw HTML for the modal content
	modalNode   *dom.Node  // Parsed modal container node
	modalActive bool       // Whether the modal is currently shown
	onModalOpen func()     // Called when modal opens
	onModalClose func()    // Called when modal closes

	// Toast state
	toasts               []*toastEntry // Active toast notifications
	toastIDCounter       int
	toastTimerChan       chan int       // Channel receiving expired toast IDs
	onToast              func(*ToastEntry) // Called when a new toast appears

	// Alert/Confirm state
	pendingAlertCallback func(bool) // Callback for alert/confirm dialog results
}

// New creates a new tuix application.
func New() *App {
	return &App{
		layout:         layout.NewEngine(),
		toastTimerChan: make(chan int, 64),
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

// OnFocus registers a callback called when an element receives focus.
func (a *App) OnFocus(fn func(el *dom.Node)) {
	a.onFocus = fn
}

// OnBlur registers a callback called when an element loses focus.
func (a *App) OnBlur(fn func(el *dom.Node)) {
	a.onBlur = fn
}

// OnHover registers a callback called when the mouse enters an element.
func (a *App) OnHover(fn func(el *dom.Node)) {
	// TODO: store callback when needed
	_ = fn
}

// HoveredElement returns the DOM node currently under the mouse cursor, or nil.
func (a *App) HoveredElement() *dom.Node {
	return a.hoveredNode
}

// FocusedElement returns the currently focused DOM node, or nil if nothing is focused.
func (a *App) FocusedElement() *dom.Node {
	if a.formFocused >= 0 && a.formFocused < len(a.formFocusables) {
		return a.formFocusables[a.formFocused]
	}
	return nil
}

// FocusElement sets focus to the given DOM node.
// If the node is not in the focusables list, focus is cleared.
func (a *App) FocusElement(node *dom.Node) {
	a.focusNode(node)
	a.renderFrame()
}

// Blur removes focus from the currently focused element.
func (a *App) Blur() {
	if a.formFocused >= 0 {
		a.focusByIndex(-1)
		a.renderFrame()
	}
}

// ShowModal shows a modal dialog with the given HTML content.
// The modal is rendered as an overlay with backdrop and centered content.
// Escape closes the modal; Tab focus is trapped within modal elements.
// If a modal is already active, it is replaced.
func (a *App) ShowModal(html string) {
	// Clean up any existing modal state
	if a.modalActive && a.modalNode != nil {
		a.cleanFormState(a.modalNode)
	}

	a.modalHTML = html
	a.modalActive = true

	// Parse modal HTML
	parser := dom.NewParser(`<div class="modal-wrapper">` + html + `</div>`)
	a.modalNode = parser.Parse()

	// Initialize form state for modal elements
	a.initFormState(a.modalNode)

	// Rebuild focusables: when modal is active, only elements within the modal are focusable
	a.buildModalFocusables()

	if a.onModalOpen != nil {
		a.onModalOpen()
	}
	a.renderFrame()
}

// CloseModal closes the currently open modal dialog.
func (a *App) CloseModal() {
	a.modalActive = false

	// Clear any pending alert callback
	a.pendingAlertCallback = nil

	// Clean up form state for modal elements
	if a.modalNode != nil {
		a.cleanFormState(a.modalNode)
	}

	a.modalNode = nil

	// Rebuild focusables to include all focusable elements in the document
	a.initFormState(a.document)
	a.formFocused = -1
	a.renderFrame()
	if a.onModalClose != nil {
		a.onModalClose()
	}
}

// buildModalFocusables collects focusable elements within the modal node only.
func (a *App) buildModalFocusables() {
	a.formFocusables = nil
	a.formFocused = -1
	if a.modalNode == nil {
		return
	}
	for _, child := range a.modalNode.Children {
		a.collectFocusables(child)
	}
	// Auto-focus the first focusable element in the modal
	if len(a.formFocusables) > 0 {
		a.focusByIndex(0)
	}
}

// OnModalOpen registers a callback called when a modal is opened.
func (a *App) OnModalOpen(fn func()) {
	a.onModalOpen = fn
}

// OnModalClose registers a callback called when a modal is closed.
func (a *App) OnModalClose(fn func()) {
	a.onModalClose = fn
}

// --- Toast / Notification API ---

// ShowToast displays a toast notification. The toast auto-dismisses after the
// specified duration. Use ToastEntry fields to customize the message, type,
// duration, and position.
func (a *App) ShowToast(entry ToastEntry) {
	if entry.Duration == 0 {
		entry.Duration = 4 * time.Second
	}
	if entry.Position == "" {
		entry.Position = "top-right"
	}

	a.toastIDCounter++
	entry.ID = a.toastIDCounter

	e := &toastEntry{ToastEntry: entry}
	a.toasts = append(a.toasts, e)

	// Schedule auto-dismiss
	go func(id int, d time.Duration) {
		time.Sleep(d)
		select {
		case a.toastTimerChan <- id:
		default:
			// Channel full or receiver gone — drop
		}
	}(entry.ID, entry.Duration)

	if a.onToast != nil {
		a.onToast(&e.ToastEntry)
	}

	a.renderFrame()
}

// DismissToast manually dismisses a specific toast by ID.
func (a *App) DismissToast(id int) {
	a.dismissToast(id)
}

// dismissToast removes a toast from the queue and re-renders.
func (a *App) dismissToast(id int) {
	for i, t := range a.toasts {
		if t.ID == id {
			a.toasts = append(a.toasts[:i], a.toasts[i+1:]...)
			break
		}
	}
	a.renderFrame()
}

// OnToast registers a callback called when a new toast is shown.
func (a *App) OnToast(fn func(entry *ToastEntry)) {
	a.onToast = fn
}

// Alert shows a modal dialog with a message and OK button.
// The modal is dismissed when the user presses Enter or clicks OK.
func (a *App) Alert(title, message string) {
	a.alertConfirm(title, message, nil)
}

// Confirm shows a modal dialog with a message and OK/Cancel buttons.
// The callback is invoked with true if OK was pressed, false if Cancel.
func (a *App) Confirm(title, message string, cb func(ok bool)) {
	a.alertConfirm(title, message, cb)
}

func (a *App) alertConfirm(title, message string, cb func(ok bool)) {
	hasCancel := cb != nil
	buttons := `<button id="alert-ok" autofocus>OK</button>`
	if hasCancel {
		buttons = `<button id="alert-ok" autofocus>OK</button> <button id="alert-cancel">Cancel</button>`
	}

	html := `<div class="dialog alert-dialog">
		<h2>` + title + `</h2>
		<p>` + message + `</p>
		<div class="actions">` + buttons + `</div>
	</div>`

	a.pendingAlertCallback = cb
	a.ShowModal(html)
}

// buildToastHTML generates the HTML for the toast container and its children.
func (a *App) buildToastHTML() string {
	if len(a.toasts) == 0 {
		return `<div class="toast-container"></div>`
	}

	var items string
	for _, t := range a.toasts {
		typeClass := "toast-info"
		switch t.Type {
		case ToastSuccess:
			typeClass = "toast-success"
		case ToastWarning:
			typeClass = "toast-warning"
		case ToastError:
			typeClass = "toast-error"
		}

		items += `<div class="toast ` + typeClass + `" data-toast-id="` + itoa(t.ID) + `">`
		items += `<span class="toast-msg">` + t.Message + `</span>`
		items += `</div>`
	}

	pos := "top-right"
	if len(a.toasts) > 0 {
		pos = a.toasts[0].Position
	}

	return `<div class="toast-container toast-container-` + pos + `">` + items + `</div>`
}

// itoa is a simple int-to-string for building HTML attributes.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
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

	// Parse CSS — toast defaults go first so user CSS can override
	combinedCSS := toastDefaultCSS
	if a.css != "" {
		combinedCSS += "\n" + a.css
	}
	cssParser := css.NewParser(combinedCSS)
	sheet, _ := cssParser.Parse()
	a.stylesheet = sheet

	// Initialize form state by scanning the DOM
	a.initFormState(a.document)
	a.buildFocusables(a.document)

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
		case toastID := <-a.toastTimerChan:
			a.dismissToast(toastID)
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

	// Let form controls handle the event first (they consume it if focused)
	if a.handleFormEvent(event) {
		return
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

		// Handle Escape — close modal if active
		if event.Key == terminal.KeyEscape && a.modalActive {
			// Save callback before CloseModal clears it
			cb := a.pendingAlertCallback
			a.CloseModal()
			// If this was a confirm dialog, invoke callback with false (cancel)
			if cb != nil {
				cb(false)
			}
			return
		}

		// Handle Ctrl+C / Ctrl+D to quit
		if event.Key == terminal.KeyCtrlC || event.Key == terminal.KeyCtrlD {
			a.Stop()
		}

		// Handle scrolling
		a.handleScrollKey(event.Key)

	case terminal.EventResize:
		a.width = event.Width
		a.height = event.Height
		a.layout.ViewWidth = float64(a.width)
		a.layout.ViewHeight = float64(a.height)
		a.renderFrame()
		if a.onResize != nil {
			a.onResize(event.Width, event.Height)
		}

	case terminal.EventMouse:
		// Store mouse position for hover tracking
		a.mouseX = event.MouseX
		a.mouseY = event.MouseY

		// Mouse click — focus the element under the cursor
		if event.MouseButton == terminal.MouseLeft && a.rootBox != nil {
			target := a.rootBox.FindBoxAtPoint(event.MouseX, event.MouseY)

			// When modal is active, only interact with modal elements
			if a.modalActive {
				// Check if click is within the modal wrapper
				if target != nil && target.Node != nil {
					if a.isNodeInModal(target.Node) {
						a.focusNode(target.Node)
						a.renderFrame()
					}
					// Clicks on background (backdrop) are ignored — keeps focus trapped
				}
			} else {
				if target != nil && target.Node != nil {
					a.focusNode(target.Node)
				} else {
					a.formFocused = -1
				}
				a.renderFrame()
			}
		}

		// Mouse wheel scrolling
		if event.MouseButton == terminal.MouseWheelUp || event.MouseButton == terminal.MouseWheelDown {
			a.handleMouseWheel(event)
		}

		// Mouse move / drag — update hover state
		if event.MouseButton == terminal.MouseNone {
			// Pure mouse move — recompute hovered element
			oldHovered := a.hoveredNode
			if a.rootBox != nil {
				target := a.rootBox.FindBoxAtPoint(event.MouseX, event.MouseY)
				if target != nil && target.Node != nil {
					a.hoveredNode = target.Node
				} else {
					a.hoveredNode = nil
				}
			} else {
				a.hoveredNode = nil
			}
			// Re-render if hover state changed
			if oldHovered != a.hoveredNode {
				a.renderFrame()
			}
		}

		if a.onMouse != nil {
			a.onMouse(event.MouseButton, event.MouseX, event.MouseY)
		}
	}
}

// handleScrollKey processes keyboard-based scrolling.
func (a *App) handleScrollKey(key terminal.Key) {
	if a.rootBox == nil || a.modalActive {
		return
	}

	// Find a scroll container to scroll
	sc := a.scrollContainer
	if sc == nil {
		sc = a.rootBox.FindScrollContainer(nil)
	}
	if sc == nil || sc.ScrollHeight <= sc.ContentRect.Height {
		return
	}

	maxY := sc.ScrollHeight - sc.ContentRect.Height
	maxX := sc.ScrollWidth - sc.ContentRect.Width
	oldY := sc.ScrollY
	oldX := sc.ScrollX

	switch key {
	case terminal.KeyUp:
		sc.ScrollY--
	case terminal.KeyDown:
		sc.ScrollY++
	case terminal.KeyPageUp:
		sc.ScrollY -= sc.ContentRect.Height
	case terminal.KeyPageDown:
		sc.ScrollY += sc.ContentRect.Height
	case terminal.KeyHome:
		if sc.ScrollX > 0 {
			sc.ScrollX = 0
		} else {
			sc.ScrollY = 0
		}
	case terminal.KeyEnd:
		if sc.ScrollX > 0 {
			sc.ScrollX = maxX
		} else {
			sc.ScrollY = maxY
		}
	case terminal.KeyLeft:
		if sc.ScrollWidth > sc.ContentRect.Width {
			sc.ScrollX--
		}
	case terminal.KeyRight:
		if sc.ScrollWidth > sc.ContentRect.Width {
			sc.ScrollX++
		}
	default:
		return // No scroll action needed
	}

	// Clamp values
	if sc.ScrollY < 0 {
		sc.ScrollY = 0
	}
	if sc.ScrollY > maxY {
		sc.ScrollY = maxY
	}
	if sc.ScrollX < 0 {
		sc.ScrollX = 0
	}
	if sc.ScrollX > maxX {
		sc.ScrollX = maxX
	}

	if sc.ScrollY != oldY || sc.ScrollX != oldX {
		a.scrollContainer = sc
		a.renderFrame()
	}
}

// handleMouseWheel processes mouse wheel scroll events.
func (a *App) handleMouseWheel(event terminal.Event) {
	if a.rootBox == nil {
		return
	}

	// When modal is active, prevent background scrolling
	if a.modalActive {
		return
	}

	// Find the box under the mouse cursor
	target := a.rootBox.FindBoxAtPoint(event.MouseX, event.MouseY)
	if target == nil {
		return
	}

	// Find scroll container ancestor for this target
	sc := a.rootBox.FindScrollContainer(target)
	if sc == nil {
		return
	}

	maxY := sc.ScrollHeight - sc.ContentRect.Height
	oldY := sc.ScrollY

	switch event.MouseButton {
	case terminal.MouseWheelUp:
		sc.ScrollY -= 3
	case terminal.MouseWheelDown:
		sc.ScrollY += 3
	}

	if sc.ScrollY < 0 {
		sc.ScrollY = 0
	}
	if sc.ScrollY > maxY {
		sc.ScrollY = maxY
	}

	if sc.ScrollY != oldY {
		a.scrollContainer = sc
		a.renderFrame()
	}
}

// renderFrame performs layout and paints the current frame to the terminal.
// Each frame is rendered from scratch (like Bubble Tea), avoiding ANSI state
// tracking issues that come with differential rendering.
func (a *App) renderFrame() {
	if a.document == nil {
		return
	}

	// Graft modal node into the document if active
	var modalParent *dom.Node
	if a.modalActive && a.modalNode != nil {
		for _, child := range a.document.Children {
			if child.Type == dom.NodeElement {
				child.Children = append(child.Children, a.modalNode)
				a.modalNode.Parent = child
				modalParent = child
				break
			}
		}
	}

	// Graft toast container into the document if there are active toasts
	var toastParent *dom.Node
	var toastContainer *dom.Node
	if len(a.toasts) > 0 {
		toastHTML := a.buildToastHTML()
		parser := dom.NewParser(toastHTML)
		toastContainer = parser.Parse()
		for _, child := range a.document.Children {
			if child.Type == dom.NodeElement {
				child.Children = append(child.Children, toastContainer)
				toastContainer.Parent = child
				toastParent = child
				break
			}
		}
	}

	// Save scroll offsets from the old box tree (keyed by *dom.Node)
	if a.rootBox != nil {
		a.saveScrollOffsets(a.rootBox)
	}

	// Prepare DOM for form controls (sets text children, focused attributes)
	a.prepareFormDOM(a.document)

	// Resolve styles for all nodes
	styles := a.resolveStyles()

	// Perform layout
	a.rootBox = a.layout.Layout(a.document, styles)

	// Restore scroll offsets on the new box tree
	a.restoreScrollOffsets(a.rootBox)

	// Paint the layout onto a fresh canvas
	a.canvas = render.NewCanvas(a.width, a.height, a.terminal.ColorMode())
	a.painter = render.NewPainter(a.canvas, a.terminal.ColorMode())
	a.painter.Paint(a.rootBox)

	// Detach modal node if it was attached
	if modalParent != nil {
		modalParent.Children = removeFromSlice(modalParent.Children, a.modalNode)
		a.modalNode.Parent = nil
	}

	// Detach toast container if it was attached
	if toastParent != nil && toastContainer != nil {
		toastParent.Children = removeFromSlice(toastParent.Children, toastContainer)
		toastContainer.Parent = nil
	}

	// Call render callback
	if a.onRender != nil {
		a.onRender()
	}

	// Full render (pass nil to always output everything)
	output := a.canvas.Render(nil)
	a.terminal.WriteString(output)
}

// saveScrollOffsets walks the box tree and saves scroll offsets.
func (a *App) saveScrollOffsets(box *layout.Box) {
	if box == nil {
		return
	}
	if box.Node != nil && (box.ScrollX != 0 || box.ScrollY != 0) {
		if a.savedScroll == nil {
			a.savedScroll = make(map[*dom.Node]struct{ X, Y int })
		}
		a.savedScroll[box.Node] = struct{ X, Y int }{X: box.ScrollX, Y: box.ScrollY}
	}
	for _, child := range box.Children {
		a.saveScrollOffsets(child)
	}
}

// restoreScrollOffsets walks the new box tree and restores saved scroll offsets.
func (a *App) restoreScrollOffsets(box *layout.Box) {
	if box == nil {
		return
	}
	if box.Node != nil && a.savedScroll != nil {
		if saved, ok := a.savedScroll[box.Node]; ok {
			box.ScrollX = saved.X
			box.ScrollY = saved.Y
		}
	}
	for _, child := range box.Children {
		a.restoreScrollOffsets(child)
	}
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

	// Resolve this node, inheriting from parent if available
	parentStyle := styles[node.Parent]
	styles[node] = resolver.ResolveWithParent(node, parentStyle)

	for _, child := range node.Children {
		a.resolveNodeStyles(child, resolver, styles)
	}
}

// ---------------------------------------------------------------------------
// Form control state management
// ---------------------------------------------------------------------------

// initFormState scans the DOM tree and initializes form control state.
func (a *App) initFormState(node *dom.Node) {
	if node == nil || node.Type != dom.NodeElement {
		return
	}

	if a.formValues == nil {
		a.formValues = make(map[*dom.Node]string)
		a.formCursors = make(map[*dom.Node]int)
		a.formChecked = make(map[*dom.Node]bool)
	}

	tag := strings.ToLower(node.Data)
	switch tag {
	case "input":
		inputType := strings.ToLower(node.GetAttribute("type"))
		// Initialize value from attribute
		if val := node.GetAttribute("value"); val != "" {
			a.formValues[node] = val
		}
		// Initialize checked state
		if inputType == "checkbox" || inputType == "radio" {
			if node.HasAttribute("checked") {
				a.formChecked[node] = true
			}
		}
		// Register as focusable
		focusableTags[node.Data] = true // ensure input is in the set

	case "textarea":
		// Initialize value from text content
		for _, child := range node.Children {
			if child.Type == dom.NodeText {
				a.formValues[node] = child.Data
				break
			}
		}

	case "select":
		// Find initially selected option
		for _, child := range node.Children {
			if child.Type == dom.NodeElement && strings.ToLower(child.Data) == "option" {
				if child.HasAttribute("selected") || a.formValues[node] == "" {
					// Use option's text content
					for _, textChild := range child.Children {
						if textChild.Type == dom.NodeText {
							a.formValues[node] = textChild.Data
							break
						}
					}
				}
				if child.HasAttribute("selected") {
					break
				}
			}
		}

	case "option":
		// Ensure default text is set
		if a.formValues[node] == "" {
			for _, child := range node.Children {
				if child.Type == dom.NodeText {
					a.formValues[node] = child.Data
					break
				}
			}
		}

	default:
		// Check if this is a focusable tag not already handled
		if focusableTags[tag] && a.formValues[node] == "" {
			// For buttons and links, initialize with their text content
			for _, child := range node.Children {
				if child.Type == dom.NodeText {
					a.formValues[node] = child.Data
					break
				}
			}
		}
	}

	// Recurse into children
	for _, child := range node.Children {
		a.initFormState(child)
	}
}

// cleanFormState removes form state entries for a subtree of nodes.
func (a *App) cleanFormState(node *dom.Node) {
	if node == nil || node.Type != dom.NodeElement {
		return
	}
	delete(a.formValues, node)
	delete(a.formCursors, node)
	delete(a.formChecked, node)
	for _, child := range node.Children {
		a.cleanFormState(child)
	}
}

// prepareFormDOM updates DOM children for form elements to reflect current state.
// This is called before each layout pass so that layout sees the correct content.
func (a *App) prepareFormDOM(node *dom.Node) {
	if node == nil || node.Type != dom.NodeElement {
		return
	}

	tag := strings.ToLower(node.Data)

	// Apply/remove 'focused' attribute for CSS styling
	if a.formFocused >= 0 && a.formFocused < len(a.formFocusables) && a.formFocusables[a.formFocused] == node {
		node.SetAttribute("focused", "")
	} else {
		// Remove 'focused' attribute if it was set
		delete(node.Attributes, "focused")
	}

	// Apply/remove 'hovered' attribute for CSS :hover matching
	// :hover matches the hovered element and all its ancestors
	if a.hoveredNode != nil && nodeHasOrIsAncestor(node, a.hoveredNode) {
		node.SetAttribute("hovered", "")
	} else {
		delete(node.Attributes, "hovered")
	}

	switch tag {
	case "input":
		inputType := strings.ToLower(node.GetAttribute("type"))
		switch inputType {
		case "checkbox":
			text := "[ ]"
			if a.formChecked[node] {
				text = "[x]"
			}
			a.setInputTextChild(node, text)
		case "radio":
			text := "( )"
			if a.formChecked[node] {
				text = "(*) "
			}
			a.setInputTextChild(node, text)
		default: // text, password, email, search, number, etc.
			val := a.formValues[node]
			isPassword := inputType == "password"

			// Mask password values
			if isPassword && val != "" {
				val = strings.Repeat("•", len([]rune(val)))
			}

			// For search inputs, show a clear indicator (×) when non-empty
			isSearch := inputType == "search"
			showClear := isSearch && len([]rune(val)) > 0

			// For number inputs, show increment/decrement indicators when focused
			isNumber := inputType == "number"

			if val == "" {
				val = node.GetAttribute("placeholder")
				if val == "" {
					val = ""
				}
			}
			// Show value with cursor if focused
			if a.formFocused >= 0 && a.formFocused < len(a.formFocusables) &&
				a.formFocusables[a.formFocused] == node {
				runes := []rune(val)
				cursor := a.formCursors[node]
				if cursor < 0 {
					cursor = 0
				}
				if cursor > len(runes) {
					cursor = len(runes)
				}
				// Insert cursor character
				display := string(runes[:cursor]) + "▎" + string(runes[cursor:])
				if showClear {
					display += " ×"
				}
				if isNumber {
					display += " ▲▼"
				}
				a.setInputTextChild(node, display)
			} else {
				if showClear {
					val += " ×"
				}
				if isNumber {
					val += " ▲▼"
				}
				a.setInputTextChild(node, val)
			}
		}

	case "textarea":
		val := a.formValues[node]
		if val == "" {
			val = node.GetAttribute("placeholder")
		}
		a.setTextareaText(node, val)

	case "select":
		selected := a.formValues[node]
		if selected == "" {
			// Find first option text
			for _, child := range node.Children {
				if child.Type == dom.NodeElement && strings.ToLower(child.Data) == "option" {
					for _, textChild := range child.Children {
						if textChild.Type == dom.NodeText {
							selected = textChild.Data
							break
						}
					}
					if selected != "" {
						break
					}
				}
			}
		}
		a.setInputTextChild(node, " "+selected+" ")

	case "button", "a":
		// These already have text children — just ensure focused styling
		// Nothing to do here for content
	}

	// Recurse into children
	for _, child := range node.Children {
		a.prepareFormDOM(child)
	}
}

// setInputTextChild sets or replaces the text child of a node (for void elements like input).
func (a *App) setInputTextChild(node *dom.Node, text string) {
	// Find existing text child
	for _, child := range node.Children {
		if child.Type == dom.NodeText {
			child.Data = text
			return
		}
	}
	// No text child exists — create one
	textNode := &dom.Node{
		Type: dom.NodeText,
		Data: text,
		Parent: node,
	}
	node.Children = append(node.Children, textNode)
}

// setTextareaText updates the text content of a textarea.
func (a *App) setTextareaText(node *dom.Node, text string) {
	// Replace all children with a single text node
	node.Children = nil
	textNode := &dom.Node{
		Type: dom.NodeText,
		Data: text,
		Parent: node,
	}
	node.Children = append(node.Children, textNode)
}

// buildFocusables scans the DOM and builds an ordered list of focusable elements.
func (a *App) buildFocusables(node *dom.Node) {
	a.formFocusables = nil
	a.formFocused = -1
	a.collectFocusables(node)
	// Auto-focus: check for element with autofocus attribute
	for i, n := range a.formFocusables {
		if n.HasAttribute("autofocus") {
			a.focusByIndex(i)
			return
		}
	}
	// If no autofocus found, check if any element already has focused attribute
	for i, n := range a.formFocusables {
		if n.HasAttribute("focused") {
			a.focusByIndex(i)
			return
		}
	}
}

func (a *App) collectFocusables(node *dom.Node) {
	if node == nil || node.Type != dom.NodeElement {
		return
	}

	tag := strings.ToLower(node.Data)
	isFocusableTag := focusableTags[tag]
	hasTabIndex := node.HasAttribute("tabindex")

	if isFocusableTag || hasTabIndex {
		disabled := node.HasAttribute("disabled")
		if !disabled {
			tabindexStr := node.GetAttribute("tabindex")
			ti := parseTabIndex(tabindexStr)
			// tabindex < 0 means programmatically focusable only (not tab-navigable)
			if ti >= 0 {
				a.formFocusables = append(a.formFocusables, node)
			} else if isFocusableTag {
				// For focusable tags, tabindex="-1" makes them programmatically
				// focusable but NOT tab-navigable. We track them separately.
				// For now, simply don't add to focusables so Tab skips them.
			}
		}
	}

	// Recurse into children
	for _, child := range node.Children {
		a.collectFocusables(child)
	}
}

// handleFormEvent processes keyboard events for focused form controls.
// Returns true if the event was consumed by a form control.
func (a *App) handleFormEvent(event terminal.Event) bool {
	if event.Type != terminal.EventKey {
		return false
	}

	// Handle Tab key (focus navigation) even when nothing is focused
	if event.Key == terminal.KeyTab && len(a.formFocusables) > 0 {
		var nextIdx int
		if a.formFocused < 0 {
			nextIdx = 0
		} else {
			dir := 1
			if event.Modifiers&terminal.ModShift != 0 {
				dir = -1
			}
			nextIdx = a.formFocused + dir
			if nextIdx < 0 {
				nextIdx = len(a.formFocusables) - 1
			} else if nextIdx >= len(a.formFocusables) {
				nextIdx = 0
			}
		}
		a.focusByIndex(nextIdx)
		a.renderFrame()
		return true
	}

	// If nothing is focused, don't consume the event
	if a.formFocused < 0 || a.formFocused >= len(a.formFocusables) {
		return false
	}

	focused := a.formFocusables[a.formFocused]
	tag := strings.ToLower(focused.Data)

	// Handle element-specific interactions
	switch tag {
	case "input":
		inputType := strings.ToLower(focused.GetAttribute("type"))
		switch inputType {
		case "checkbox":
			if event.Key == terminal.KeyEnter || event.Key == terminal.KeySpace {
				a.formChecked[focused] = !a.formChecked[focused]
				a.renderFrame()
				return true
			}
		case "radio":
			if event.Key == terminal.KeyEnter || event.Key == terminal.KeySpace {
				a.checkRadio(focused)
				a.renderFrame()
				return true
			}
		case "search":
			// Search input: Escape clears the value, retains focus
			if event.Key == terminal.KeyEscape {
				a.formValues[focused] = ""
				a.formCursors[focused] = 0
				a.renderFrame()
				return true
			}
			return a.handleTextEdit(event, focused)
		case "number":
			// Number input: up/down arrows increment/decrement
			return a.handleNumberEdit(event, focused)
		default: // text, password, email, etc.
			return a.handleTextEdit(event, focused)
		}

	case "textarea":
		return a.handleTextEdit(event, focused)

	case "select":
		if event.Key == terminal.KeyEnter || event.Key == terminal.KeySpace {
			a.cycleSelectOption(focused)
			a.renderFrame()
			return true
		}

	case "button":
		if event.Key == terminal.KeyEnter || event.Key == terminal.KeySpace {
			// Check for alert/confirm dialog buttons
			id := focused.GetAttribute("id")
			if id == "alert-ok" {
				if a.pendingAlertCallback != nil {
					cb := a.pendingAlertCallback
					a.pendingAlertCallback = nil
					a.CloseModal()
					cb(true)
				} else {
					// Alert without callback — just close
					a.CloseModal()
				}
				return true
			}
			if id == "alert-cancel" && a.pendingAlertCallback != nil {
				cb := a.pendingAlertCallback
				a.pendingAlertCallback = nil
				a.CloseModal()
				cb(false)
				return true
			}
			// General button activation
			a.renderFrame()
			return true
		}

	case "a":
		// Links are focusable but have no special form handling yet
		return false
	}

	return false
}

// handleTextEdit processes keyboard input for text fields.
// Returns true if the event was consumed.
func (a *App) handleTextEdit(event terminal.Event, node *dom.Node) bool {
	// If readonly, don't allow editing but still allow cursor navigation
	if node.HasAttribute("readonly") {
		switch event.Key {
		case terminal.KeyLeft:
			if cursor := a.formCursors[node]; cursor > 0 {
				a.formCursors[node] = cursor - 1
				a.renderFrame()
			}
			return true
		case terminal.KeyRight:
			cursor := a.formCursors[node]
			val := []rune(a.formValues[node])
			if cursor < len(val) {
				a.formCursors[node] = cursor + 1
				a.renderFrame()
			}
			return true
		case terminal.KeyHome:
			a.formCursors[node] = 0
			a.renderFrame()
			return true
		case terminal.KeyEnd:
			a.formCursors[node] = len([]rune(a.formValues[node]))
			a.renderFrame()
			return true
		default:
			// All other keys are consumed but ignored in readonly mode
			return true
		}
	}

	val := []rune(a.formValues[node])
	cursor := a.formCursors[node]
	if cursor < 0 {
		cursor = 0
	}
	if cursor > len(val) {
		cursor = len(val)
	}

	switch event.Key {
	case terminal.KeyLeft:
		if cursor > 0 {
			a.formCursors[node] = cursor - 1
			a.renderFrame()
		}
		return true

	case terminal.KeyRight:
		if cursor < len(val) {
			a.formCursors[node] = cursor + 1
			a.renderFrame()
		}
		return true

	case terminal.KeyHome:
		a.formCursors[node] = 0
		a.renderFrame()
		return true

	case terminal.KeyEnd:
		a.formCursors[node] = len(val)
		a.renderFrame()
		return true

	case terminal.KeyBackspace:
		if cursor > 0 {
			val = append(val[:cursor-1], val[cursor:]...)
			a.formValues[node] = string(val)
			a.formCursors[node] = cursor - 1
			a.renderFrame()
		}
		return true

	case terminal.KeyDelete:
		if cursor < len(val) {
			val = append(val[:cursor], val[cursor+1:]...)
			a.formValues[node] = string(val)
			a.renderFrame()
		}
		return true

	case terminal.KeyEnter:
		// Enter in a text field — stay focused or submit
		return false

	case terminal.KeyEscape:
		// Escape — lose focus
		a.formFocused = -1
		a.renderFrame()
		return true

	default:
		// Typing a character
		if event.Key == terminal.KeyRune {
			r := event.Rune
			// Ignore control characters
			if r < 0x20 {
				return false
			}
			newRunes := make([]rune, 0, len(val)+1)
			newRunes = append(newRunes, val[:cursor]...)
			newRunes = append(newRunes, r)
			newRunes = append(newRunes, val[cursor:]...)
			a.formValues[node] = string(newRunes)
			a.formCursors[node] = cursor + 1
			a.renderFrame()
			return true
		}
	}

	return false
}

// handleNumberEdit processes keyboard events for number inputs.
// Up/down arrows increment/decrement the value.
func (a *App) handleNumberEdit(event terminal.Event, node *dom.Node) bool {
	val := a.formValues[node]

	// Parse current value, default to 0
	stepStr := node.GetAttribute("step")
	step := 1
	if stepStr != "" {
		if s, err := parseInt(stepStr); err == nil && s > 0 {
			step = s
		}
	}

	minStr := node.GetAttribute("min")
	maxStr := node.GetAttribute("max")
	hasMin := minStr != ""
	hasMax := maxStr != ""
	minVal := 0
	maxVal := 0
	if hasMin {
		minVal, _ = parseInt(minStr)
	}
	if hasMax {
		maxVal, _ = parseInt(maxStr)
	}

	switch event.Key {
	case terminal.KeyUp:
		n := 0
		if val != "" {
			n, _ = parseInt(val)
		}
		n += step
		if hasMax && n > maxVal {
			n = maxVal
		}
		a.formValues[node] = itoa(n)
		a.formCursors[node] = len([]rune(a.formValues[node]))
		a.renderFrame()
		return true

	case terminal.KeyDown:
		n := 0
		if val != "" {
			n, _ = parseInt(val)
		}
		n -= step
		if hasMin && n < minVal {
			n = minVal
		}
		a.formValues[node] = itoa(n)
		a.formCursors[node] = len([]rune(a.formValues[node]))
		a.renderFrame()
		return true

	default:
		// Fall through to text editing for typing numbers
		return a.handleTextEdit(event, node)
	}
}

// parseInt parses a string to int, returning 0 and an error on failure.
func parseInt(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}
	n := 0
	neg := false
	start := 0
	if s[0] == '-' {
		neg = true
		start = 1
	} else if s[0] == '+' {
		start = 1
	}
	for i := start; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			return 0, fmt.Errorf("invalid number: %s", s)
		}
	}
	if neg {
		return -n, nil
	}
	return n, nil
}

// focusNode finds the given node in formFocusables and sets focus to it.
func (a *App) focusNode(node *dom.Node) {
	if node == nil {
		a.focusByIndex(-1)
		return
	}
	for i, n := range a.formFocusables {
		if n == node {
			a.focusByIndex(i)
			return
		}
	}
	// Node not found in focusables — walk up to find a focusable ancestor
	for p := node.Parent; p != nil; p = p.Parent {
		for i, n := range a.formFocusables {
			if n == p {
				a.focusByIndex(i)
				return
			}
		}
	}
	a.focusByIndex(-1)
}

// focusByIndex sets focus to the element at the given index.
// Pass -1 to clear focus.
func (a *App) focusByIndex(idx int) {
	old := a.formFocused
	if old >= 0 && old < len(a.formFocusables) && a.onBlur != nil {
		a.onBlur(a.formFocusables[old])
	}
	a.formFocused = idx
	if idx >= 0 && idx < len(a.formFocusables) && a.onFocus != nil {
		a.onFocus(a.formFocusables[idx])
	}
}

// checkRadio toggles a radio button and unchecks others in the same group.
func (a *App) checkRadio(node *dom.Node) {
	name := node.GetAttribute("name")
	a.formChecked[node] = true
	if name != "" {
		// Uncheck all other radios with the same name
		for _, focusable := range a.formFocusables {
			if focusable != node &&
				strings.ToLower(focusable.Data) == "input" &&
				strings.ToLower(focusable.GetAttribute("type")) == "radio" &&
				focusable.GetAttribute("name") == name {
				a.formChecked[focusable] = false
			}
		}
	}
}

// cycleSelectOption moves to the next option in a select element.
func (a *App) cycleSelectOption(node *dom.Node) {
	options := make([]string, 0)
	for _, child := range node.Children {
		if child.Type == dom.NodeElement && strings.ToLower(child.Data) == "option" {
			for _, textChild := range child.Children {
				if textChild.Type == dom.NodeText {
					options = append(options, textChild.Data)
					break
				}
			}
		}
	}
	if len(options) == 0 {
		return
	}

	current := a.formValues[node]
	for i, opt := range options {
		if opt == current {
			next := (i + 1) % len(options)
			a.formValues[node] = options[next]
			return
		}
	}
	// Current value not found — set to first option
	a.formValues[node] = options[0]
}

// parseTabIndex parses a tabindex attribute value string and returns the integer value.
// Returns 0 if the string is empty or unparseable.
func parseTabIndex(s string) int {
	if s == "" {
		return 0
	}
	neg := false
	start := 0
	if s[0] == '-' {
		neg = true
		start = 1
	} else if s[0] == '+' {
		start = 1
	}
	n := 0
	for i := start; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			return 0 // invalid, treat as 0
		}
	}
	if neg {
		return -n
	}
	return n
}

// nodeHasOrIsAncestor returns true if 'node' is 'target' or an ancestor of 'target'.
func nodeHasOrIsAncestor(node, target *dom.Node) bool {
	for p := target; p != nil; p = p.Parent {
		if p == node {
			return true
		}
	}
	return false
}

// removeFromSlice removes the first occurrence of target from slice and returns the new slice.
func removeFromSlice(slice []*dom.Node, target *dom.Node) []*dom.Node {
	for i, n := range slice {
		if n == target {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// isNodeInModal returns true if the given node is a descendant of the modal wrapper.
func (a *App) isNodeInModal(node *dom.Node) bool {
	if a.modalNode == nil || node == nil {
		return false
	}
	for p := node; p != nil; p = p.Parent {
		if p == a.modalNode {
			return true
		}
	}
	return false
}
