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
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/elioneto/tuix/color"
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

// designSystemCSS provides pre-built component classes for a consistent UI.
// These use placeholder hex values that get overridden by themeToCSS.
// The themeToCSS() output is applied on top so the actual theme colors win.
const designSystemCSS = `
/* Button variants */
.btn {
	padding: 0 2;
	border: solid;
	cursor: pointer;
	text-align: center;
}
.btn-primary {
	background-color: #00d4aa;
	color: #1a1a2e;
	font-weight: bold;
}
.btn-secondary {
	border-color: #e94560;
	color: #e94560;
}
.btn-danger {
	border-color: #e74c3c;
	color: #e74c3c;
}
.btn-ghost {
	border: none;
	background: transparent;
}
.btn-sm { padding: 0 1; }
.btn-lg { padding: 0 4; }

/* Input variants */
.input {
	padding: 0 1;
	border: solid;
}
.input-error {
	border-color: #e74c3c;
}
.input-sm { padding: 0 1; }
.input-lg { padding: 0 3; }

/* Card component */
.card {
	border: solid #0f3460;
	background-color: #16213e;
	padding: 1;
}

/* Badge / Tag */
.badge {
	padding: 0 1;
	border: solid;
	font-weight: bold;
}
.badge-primary {
	background-color: #00d4aa;
	color: #1a1a2e;
}
.badge-success {
	background-color: #2ecc71;
	color: #1a1a2e;
}
.badge-warning {
	background-color: #f39c12;
	color: #1a1a2e;
}
.badge-error {
	background-color: #e74c3c;
	color: #1a1a2e;
}

/* Layout utilities */
.flex { display: flex; }
.flex-col { display: flex; flex-direction: column; }
.flex-wrap { flex-wrap: wrap; }
.items-center { align-items: center; }
.justify-center { justify-content: center; }
.justify-between { justify-content: space-between; }
.gap-1 { gap: 1; }
.gap-2 { gap: 2; }
.gap-4 { gap: 4; }
.w-full { width: 100%; }
.text-center { text-align: center; }
.text-bold { font-weight: bold; }

/* Typography scale */
.text-xs { font-size: 8; }
.text-sm { font-size: 12; }
.text-base { font-size: 16; }
.text-lg { font-size: 20; }
.text-xl { font-size: 24; }
.text-xxl { font-size: 32; }

/* Color utilities */
.text-primary   { color: #00d4aa; }
.text-secondary { color: #0f3460; }
.text-accent    { color: #e94560; }
.text-success   { color: #2ecc71; }
.text-warning   { color: #f39c12; }
.text-error     { color: #e74c3c; }
.text-muted     { color: #555; }
.bg-primary     { background-color: #00d4aa; }
.bg-secondary   { background-color: #0f3460; }
.bg-accent      { background-color: #e94560; }
.bg-success     { background-color: #2ecc71; }
.bg-warning     { background-color: #f39c12; }
.bg-error       { background-color: #e74c3c; }

/* Spacing utilities */
.p-1  { padding: 1; }
.p-2  { padding: 2; }
.p-4  { padding: 4; }
.px-1 { padding-left: 1; padding-right: 1; }
.px-2 { padding-left: 2; padding-right: 2; }
.py-1 { padding-top: 1; padding-bottom: 1; }
.py-2 { padding-top: 2; padding-bottom: 2; }
.m-1  { margin: 1; }
.m-2  { margin: 2; }
.m-4  { margin: 4; }
.mx-1 { margin-left: 1; margin-right: 1; }
.mx-2 { margin-left: 2; margin-right: 2; }
.my-1 { margin-top: 1; margin-bottom: 1; }
.my-2 { margin-top: 2; margin-bottom: 2; }

/* Animations — these apply dynamic effects via the animation system */
.animate-spin   { font-weight: bold; }
.animate-pulse  { }
.animate-blink  { }

/* Navbar */
.navbar {
	display: flex;
	align-items: center;
	padding: 0 1;
	background-color: #0f3460;
	border: solid #0f3460;
}
.navbar .nav-brand {
	font-weight: bold;
	color: #00d4aa;
	margin-right: 2;
}
.navbar .nav-item {
	padding: 0 1;
	color: #c0c0c0;
}
.navbar .nav-item:hover {
	color: #00d4aa;
}

/* List group */
.list {
	border: solid #0f3460;
}
.list-item {
	padding: 0 1;
	border-bottom: solid #0f3460;
}
.list-item:last-child {
	border-bottom: none;
}
.list-item:hover {
	background-color: #0f3460;
}

/* Tabs */
.tabs {
	display: flex;
	border-bottom: solid #0f3460;
}
.tab {
	padding: 0 2;
	color: #c0c0c0;
}
.tab-active {
	color: #00d4aa;
	font-weight: bold;
	border-bottom: solid #00d4aa;
}
.tab:hover {
	color: #00d4aa;
}

/* Table */
.table {
	width: 100%;
	border: solid #0f3460;
}
.table-header {
	font-weight: bold;
	background-color: #0f3460;
	padding: 0 1;
	border-bottom: solid #0f3460;
}
.table-row {
	padding: 0 1;
}
.table-row:nth-child(even) {
	background-color: #16213e;
}

/* Grid */
.grid-2 { display: flex; gap: 1; }
.grid-2 > * { flex: 1; }
.grid-3 { display: flex; gap: 1; }
.grid-3 > * { flex: 1; }
`

// themeToCSS generates CSS rules from a Theme.
func themeToCSS(t Theme) string {
	return fmt.Sprintf(`
/* Theme colors */
.bg-primary    { background-color: %s; }
.bg-secondary  { background-color: %s; }
.bg-accent     { background-color: %s; }
.bg-success    { background-color: %s; }
.bg-warning    { background-color: %s; }
.bg-error      { background-color: %s; }
.bg-surface    { background-color: %s; }
.bg-background { background-color: %s; }
.text-primary    { color: %s; }
.text-secondary  { color: %s; }
.text-accent     { color: %s; }
.text-success    { color: %s; }
.text-warning    { color: %s; }
.text-error      { color: %s; }
.text-muted      { color: %s; }
.border-primary    { border-color: %s; }
.border-secondary  { border-color: %s; }
.border-accent     { border-color: %s; }
.border-success    { border-color: %s; }
.border-warning    { border-color: %s; }
.border-error      { border-color: %s; }
`,
		colorToHex(t.Primary), colorToHex(t.Secondary), colorToHex(t.Accent),
		colorToHex(t.Success), colorToHex(t.Warning), colorToHex(t.Error),
		colorToHex(t.Surface), colorToHex(t.Background),
		colorToHex(t.Primary), colorToHex(t.Secondary), colorToHex(t.Accent),
		colorToHex(t.Success), colorToHex(t.Warning), colorToHex(t.Error),
		colorToHex(t.Muted),
		colorToHex(t.Primary), colorToHex(t.Secondary), colorToHex(t.Accent),
		colorToHex(t.Success), colorToHex(t.Warning), colorToHex(t.Error),
	)
}

// colorToHex returns the hex string representation of a color (e.g. "#ff6600").
func colorToHex(c color.Color) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

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

// Theme defines semantic colors for a design system.
// Each slot maps to a specific visual role in the UI.
type Theme struct {
	Primary   color.Color
	Secondary color.Color
	Accent    color.Color
	Success   color.Color
	Warning   color.Color
	Error     color.Color
	Surface   color.Color
	Background color.Color
	Text      color.Color
	Muted     color.Color
	Border    color.Color
	Focus     color.Color
}

// DefaultDarkTheme is the built-in dark color scheme.
var DefaultDarkTheme = Theme{
	Primary:    color.Color{Type: color.ColorTrue, R: 0x00, G: 0xD4, B: 0xAA},
	Secondary:  color.Color{Type: color.ColorTrue, R: 0x0F, G: 0x34, B: 0x60},
	Accent:     color.Color{Type: color.ColorTrue, R: 0xE9, G: 0x45, B: 0x60},
	Success:    color.Color{Type: color.ColorTrue, R: 0x2E, G: 0xCC, B: 0x71},
	Warning:    color.Color{Type: color.ColorTrue, R: 0xF3, G: 0x9C, B: 0x12},
	Error:      color.Color{Type: color.ColorTrue, R: 0xE7, G: 0x4C, B: 0x3C},
	Surface:    color.Color{Type: color.ColorTrue, R: 0x16, G: 0x21, B: 0x3E},
	Background: color.Color{Type: color.ColorTrue, R: 0x1A, G: 0x1A, B: 0x2E},
	Text:       color.Color{Type: color.ColorTrue, R: 0xC0, G: 0xC0, B: 0xC0},
	Muted:      color.Color{Type: color.ColorTrue, R: 0x55, G: 0x55, B: 0x55},
	Border:     color.Color{Type: color.ColorTrue, R: 0x0F, G: 0x34, B: 0x60},
	Focus:      color.Color{Type: color.ColorTrue, R: 0x00, G: 0xD4, B: 0xAA},
}

// DefaultLightTheme is the built-in light color scheme.
var DefaultLightTheme = Theme{
	Primary:    color.Color{Type: color.ColorTrue, R: 0x00, G: 0x7B, B: 0xFF},
	Secondary:  color.Color{Type: color.ColorTrue, R: 0x6C, G: 0x75, B: 0x7D},
	Accent:     color.Color{Type: color.ColorTrue, R: 0xE9, G: 0x45, B: 0x60},
	Success:    color.Color{Type: color.ColorTrue, R: 0x2E, G: 0xCC, B: 0x71},
	Warning:    color.Color{Type: color.ColorTrue, R: 0xF3, G: 0x9C, B: 0x12},
	Error:      color.Color{Type: color.ColorTrue, R: 0xE7, G: 0x4C, B: 0x3C},
	Surface:    color.Color{Type: color.ColorTrue, R: 0xFF, G: 0xFF, B: 0xFF},
	Background: color.Color{Type: color.ColorTrue, R: 0xF5, G: 0xF5, B: 0xF5},
	Text:       color.Color{Type: color.ColorTrue, R: 0x33, G: 0x33, B: 0x33},
	Muted:      color.Color{Type: color.ColorTrue, R: 0x99, G: 0x99, B: 0x99},
	Border:     color.Color{Type: color.ColorTrue, R: 0xDD, G: 0xDD, B: 0xDD},
	Focus:      color.Color{Type: color.ColorTrue, R: 0x00, G: 0x7B, B: 0xFF},
}

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
	scrollContainer *layout.Box // The currently focused scroll container

	// Mouse drag state
	mouseDown bool // Whether a mouse button is currently held down

	// Dragged element for slider drag tracking
	dragTarget *dom.Node // The node being dragged (e.g., range slider thumb)

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

	// Tooltip state
	tooltipTimerChan   chan struct{} // Signals when tooltip delay elapses
	tooltipPendingNode *dom.Node    // Node whose tooltip timer is pending
	tooltipNode        *dom.Node    // Node whose tooltip is currently shown
	tooltipText        string       // Text to display in tooltip
	tooltipX, tooltipY int          // Position to render tooltip (near cursor)

	// Design system / theme
	theme            Theme
	useDesignSystem  bool
	themeCSS         string // Generated CSS from the active theme
	themeChanged     bool   // Flag to regenerate stylesheet on next render

	// Animation state
	animFrame      int           // Current animation frame counter
	animTicker     *time.Ticker  // Periodic ticker for animations
	animRunning    bool          // Whether animations are active
	animTickerChan chan struct{} // Signal channel for animation ticks

	// Datalist state
	datalistMap      map[string][]string // datalist id → option texts
	datalistInput    *dom.Node           // Input currently showing a datalist dropdown
	datalistFiltered []string            // Filtered suggestion list
	datalistHighlight int                // Index of highlighted suggestion (-1 = none)
	datalistOpen     bool                // Whether the suggestion dropdown is open
}

// New creates a new tuix application.
func New() *App {
	return &App{
		layout:         layout.NewEngine(),
		toastTimerChan: make(chan int, 64),
		tooltipTimerChan: make(chan struct{}, 8),
		datalistMap:    make(map[string][]string),
		animTickerChan: make(chan struct{}, 8),
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
}

// SetTheme sets the active color theme and regenerates theme CSS rules.
func (a *App) SetTheme(t Theme) {
	a.theme = t
	a.themeCSS = themeToCSS(t)
	a.themeChanged = true
}

// UseDesignSystem enables the built-in design system CSS with component classes.
// Call SetTheme() before this to apply a custom theme.
func (a *App) UseDesignSystem() {
	a.useDesignSystem = true
	if a.themeCSS == "" {
		a.SetTheme(DefaultDarkTheme)
	}
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

	// Parse CSS — uses rebuildStylesheet to combine toast + design system + theme + user styles
	a.rebuildStylesheet()

	// Initialize form state by scanning the DOM
	a.initFormState(a.document)
	a.buildFocusables(a.document)
	a.buildDatalistMap(a.document)

	// Create canvas
	a.canvas = render.NewCanvas(a.width, a.height, a.terminal.ColorMode())

	// Call init callback
	if a.onInit != nil {
		a.onInit()
	}

	// Initial render
	a.renderFrame()

	// Start animation ticker (ticks every 200ms for animations)
	a.animTicker = time.NewTicker(200 * time.Millisecond)
	defer a.animTicker.Stop()
	a.animRunning = true

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
		case <-a.tooltipTimerChan:
			if a.tooltipPendingNode != nil && a.tooltipPendingNode == a.hoveredNode {
				// Still hovering the same node — show its tooltip
				a.tooltipNode = a.tooltipPendingNode
				a.tooltipText = a.tooltipPendingNode.GetAttribute("title")
				a.renderFrame()
			}
		case <-a.animTicker.C:
			// Animation frame tick
			a.animFrame++
			a.renderFrame()
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
						a.handleRangeClick(event.MouseX, target)
						a.renderFrame()
					}
					// Clicks on background (backdrop) are ignored — keeps focus trapped
				}
			} else {
				if target != nil && target.Node != nil {
					a.focusNode(target.Node)
					a.handleRangeClick(event.MouseX, target)
				} else {
					a.closeDatalist()
					a.formFocused = -1
				}
				a.renderFrame()
			}
		}

		// Mouse wheel scrolling
		if event.MouseButton == terminal.MouseWheelUp || event.MouseButton == terminal.MouseWheelDown {
			a.handleMouseWheel(event)
		}

		// Mouse move / drag — update hover state and handle range drag
		if event.MouseButton == terminal.MouseNone {
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

			// Handle range slider drag
			if a.dragTarget != nil {
				a.dragRangeByMouse(event.MouseX)
			}

			// Tooltip management: start/dismiss on hover change
			if oldHovered != a.hoveredNode {
				// Dismiss any active or pending tooltip
				a.tooltipNode = nil
				a.tooltipText = ""
				a.tooltipPendingNode = nil

				// If new hovered node has a title attribute, start tooltip timer
				if a.hoveredNode != nil {
					if title := a.hoveredNode.GetAttribute("title"); title != "" {
						a.tooltipPendingNode = a.hoveredNode
						a.tooltipX = event.MouseX
						a.tooltipY = event.MouseY
						go func(node *dom.Node, tx, ty int) {
							time.Sleep(500 * time.Millisecond)
							select {
							case a.tooltipTimerChan <- struct{}{}:
							default:
							}
						}(a.hoveredNode, event.MouseX, event.MouseY)
					}
				}
			}

			// Re-render if hover state changed, drag was active, or tooltip state changed
			if oldHovered != a.hoveredNode || a.dragTarget != nil {
				a.renderFrame()
			}
		} else {
			// Any button event (including release) ends drag
			a.dragTarget = nil
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

	// Rebuild stylesheet if theme changed
	if a.themeChanged {
		a.rebuildStylesheet()
		a.themeChanged = false
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

	// Paint tooltip if active
	if a.tooltipNode != nil && a.tooltipText != "" {
		a.paintTooltip()
	}

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

	// Render images using Kitty/Sixel protocol if supported
	if a.terminal.SupportsImages() && a.rootBox != nil {
		a.renderImages(a.rootBox)
	}
}

// rebuildStylesheet re-parses the combined CSS (toast + design system + theme + user).
func (a *App) rebuildStylesheet() {
	combinedCSS := toastDefaultCSS
	if a.useDesignSystem {
		combinedCSS = designSystemCSS + "\n" + combinedCSS
	}
	if a.themeCSS != "" {
		combinedCSS = a.themeCSS + "\n" + combinedCSS
	}
	if a.css != "" {
		combinedCSS += "\n" + a.css
	}
	cssParser := css.NewParser(combinedCSS)
	sheet, _ := cssParser.Parse()
	a.stylesheet = sheet
}

// paintTooltip renders a tooltip popup near the mouse cursor.
func (a *App) paintTooltip() {
	if a.canvas == nil || a.tooltipText == "" {
		return
	}

	// Compute tooltip position: 2 cells right, 1 cell down from mouse
	tipX := a.tooltipX + 2
	tipY := a.tooltipY + 1

	// Measure tooltip width
	runes := []rune(a.tooltipText)
	if len(runes) == 0 {
		return
	}
	tipW := len(runes) + 2 // +2 for padding on each side
	tipH := 1 // single line

	// Clamp to screen
	if tipX+tipW >= a.width {
		tipX = a.width - tipW - 1
	}
	if tipX < 0 {
		tipX = 0
	}
	if tipY+tipH >= a.height {
		tipY = a.height - tipH - 1
	}
	if tipY < 0 {
		tipY = 0
	}

	// Colors
	bg := color.Color{Type: color.ColorTrue, R: 0x33, G: 0x33, B: 0x33}
	fg := color.Color{Type: color.ColorTrue, R: 0xFF, G: 0xFF, B: 0xFF}

	// Draw background and text
	for x := tipX; x < tipX+tipW; x++ {
		a.canvas.Set(x, tipY, ' ', fg, bg)
	}
	// Draw text
	for i, r := range runes {
		if tipX+1+i < tipX+tipW {
			a.canvas.Set(tipX+1+i, tipY, r, fg, bg)
		}
	}
}

// renderImages walks the box tree and renders <img> elements using the
// terminal's image protocol (Kitty or Sixel).
func (a *App) renderImages(box *layout.Box) {
	if box == nil {
		return
	}

	// Check if this box is an <img> element
	if box.Node != nil && strings.ToLower(box.Node.Data) == "img" {
		src := box.Node.GetAttribute("src")
		if src == "" {
			return
		}

		// Determine image dimensions
		width := box.ContentRect.Width
		height := box.ContentRect.Height
		if width <= 0 || height <= 0 {
			// Default size: 20x10 cells
			width = 20
			height = 10
		}

		// Load image from file
		imgFile, err := os.Open(src)
		if err != nil {
			return
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			return
		}

		// Render using the appropriate protocol
		imgWidth := width
		imgHeight := height
		// Position at the content area of the box
		posX := box.ContentRect.X
		posY := box.ContentRect.Y

		switch a.terminal.ImageMode() {
		case 1: // Kitty protocol
			// Position cursor to the image location
			fmt.Fprintf(a.terminal, "\x1b[%d;%dH", posY+1, posX+1)
			render.EncodeKitty(a.terminal, img, imgWidth, imgHeight)
		case 2: // Sixel protocol
			// Position cursor to the image location
			fmt.Fprintf(a.terminal, "\x1b[%d;%dH", posY+1, posX+1)
			render.EncodeSixel(a.terminal, img)
		}
	}

	// Recurse into children
	for _, child := range box.Children {
		a.renderImages(child)
	}
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
		isMultiple := node.HasAttribute("multiple")
		if isMultiple {
			// Collect all initially selected options (including inside optgroup)
			var selectedOptions []string
			var walkSelected func(*dom.Node)
			walkSelected = func(n *dom.Node) {
				if n == nil || n.Type != dom.NodeElement {
					return
				}
				tag := strings.ToLower(n.Data)
				if tag == "option" && n.HasAttribute("selected") {
					for _, tc := range n.Children {
						if tc.Type == dom.NodeText {
							selectedOptions = append(selectedOptions, tc.Data)
							break
						}
					}
				}
				if tag == "optgroup" {
					for _, child := range n.Children {
						walkSelected(child)
					}
				}
			}
			for _, child := range node.Children {
				walkSelected(child)
			}
			if len(selectedOptions) > 0 {
				a.formValues[node] = strings.Join(selectedOptions, "|")
			}
		} else {
			// Find initially selected option (including inside optgroup)
			var findSelected func(*dom.Node) bool
			findSelected = func(n *dom.Node) bool {
				if n == nil || n.Type != dom.NodeElement {
					return false
				}
				tag := strings.ToLower(n.Data)
				if tag == "option" {
					if n.HasAttribute("selected") || a.formValues[node] == "" {
						for _, tc := range n.Children {
							if tc.Type == dom.NodeText {
								a.formValues[node] = tc.Data
								return n.HasAttribute("selected")
							}
						}
					}
				}
				if tag == "optgroup" {
					for _, child := range n.Children {
						if findSelected(child) {
							return true
						}
					}
				}
				return false
			}
			for _, child := range node.Children {
				if findSelected(child) {
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
		case "search":
			val := a.formValues[node]
			showClear := val != ""
			if val == "" {
				val = node.GetAttribute("placeholder")
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
				a.setInputTextChild(node, display)
			} else {
				if showClear {
					val += " ×"
				}
				a.setInputTextChild(node, val)
			}

		case "range":
			// Range slider handled separately
			val := a.formValues[node]
			a.prepareRangeDOM(node, val)

		case "color":
			// Color picker handled separately
			val := a.formValues[node]
			if val == "" {
				val = node.GetAttribute("value")
				if val == "" {
					val = "#000000"
				}
				a.formValues[node] = val
			}
			a.prepareColorDOM(node, val)

		case "date":
			// Date picker handled separately
			val := a.formValues[node]
			if val == "" {
				val = node.GetAttribute("value")
				if val == "" {
					// Default to today's date
					now := time.Now()
					val = fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
				}
				a.formValues[node] = val
			}
			a.prepareDateDOM(node, val)

		case "file":
			// File input
			val := a.formValues[node]
			if val == "" {
				val = "Choose file..."
			}
			display := "📎 " + val
			isFocused := a.formFocused >= 0 && a.formFocused < len(a.formFocusables) && a.formFocusables[a.formFocused] == node
			if isFocused {
				display += " [Enter to browse]"
			}
			a.setInputTextChild(node, display)

		default: // text, password, email, number, etc.
			val := a.formValues[node]
			isPassword := inputType == "password"
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
				if isPassword {
					// Replace all characters with password mask
					runes = []rune(strings.Repeat("•", len(runes)))
					if cursor < 0 {
						cursor = 0
					}
					if cursor > len(runes) {
						cursor = len(runes)
					}
					display := string(runes[:cursor]) + "▎" + string(runes[cursor:])
					a.setInputTextChild(node, display)
				} else {
					// Insert cursor character
					display := string(runes[:cursor]) + "▎" + string(runes[cursor:])
					if isNumber {
						display += " ▲▼"
					}
					a.setInputTextChild(node, display)
				}
			} else {
				if isPassword {
					val = strings.Repeat("•", len([]rune(val)))
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
		isMultiple := node.HasAttribute("multiple")
		if isMultiple && a.formValues[node] == "" {
			a.formValues[node] = ""
		}

		selected := a.formValues[node]
		if !isMultiple && selected == "" {
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

		// Collect all option texts
		options := collectSelectOptions(node)

		// Check if focused — show dropdown
		isFocused := a.formFocused >= 0 && a.formFocused < len(a.formFocusables) && a.formFocusables[a.formFocused] == node
		if isFocused && len(options) > 0 {
			var display string
			if isMultiple {
				display += " " + selected + " [multi]\n"
				selectedSet := make(map[string]bool)
				if selected != "" {
					for _, s := range strings.Split(selected, "|") {
						selectedSet[strings.TrimSpace(s)] = true
					}
				}
				for _, opt := range options {
					if selectedSet[opt] {
						display += "  [✓] " + opt + "\n"
					} else {
						display += "  [ ] " + opt + "\n"
					}
				}
			} else {
				display += " " + selected + " ▲▼\n"
				for _, opt := range options {
					prefix := "  ○ "
					if opt == selected {
						prefix = "  ● "
					}
					display += prefix + opt + "\n"
				}
			}
			// Remove trailing newline
			display = strings.TrimRight(display, "\n")
			a.setInputTextChild(node, display)
		} else {
			if isMultiple {
				// Show count of selected items
				count := 0
				if selected != "" {
					count = len(strings.Split(selected, "|"))
				}
				display := fmt.Sprintf(" %d selected", count)
				if count == 0 {
					display = " (none selected)"
				}
				a.setInputTextChild(node, display)
			} else {
				a.setInputTextChild(node, " "+selected+" ")
			}
		}

	case "button", "a":
		// These already have text children — just ensure focused styling
		// Nothing to do here for content

	case "progress":
		a.prepareProgressDOM(node)

	case "meter":
		a.prepareMeterDOM(node)
	}

	// Render datalist suggestion dropdown for focused inputs with an open datalist
	if a.datalistOpen && a.datalistInput == node {
		// Re-filter suggestions based on current value
		options := a.getDatalistForInput(node)
		if len(options) > 0 {
			a.datalistFiltered = nil
			val := strings.ToLower(a.formValues[node])
			for _, opt := range options {
				if val == "" || strings.HasPrefix(strings.ToLower(opt), val) {
					a.datalistFiltered = append(a.datalistFiltered, opt)
				}
			}
			if len(a.datalistFiltered) == 0 {
				a.datalistFiltered = nil
			}
		}

		if len(a.datalistFiltered) > 0 {
			// Clamp highlight
			if a.datalistHighlight < 0 {
				a.datalistHighlight = 0
			}
			if a.datalistHighlight >= len(a.datalistFiltered) {
				a.datalistHighlight = len(a.datalistFiltered) - 1
			}

			currentText := ""
			for _, child := range node.Children {
				if child.Type == dom.NodeText {
					currentText = child.Data
					break
				}
			}
			// Append suggestions below the current value
			var sb strings.Builder
			sb.WriteString(currentText)
			for i, opt := range a.datalistFiltered {
				sb.WriteString("\n")
				if i == a.datalistHighlight {
					sb.WriteString("  ● ")
				} else {
					sb.WriteString("  ○ ")
				}
				sb.WriteString(opt)
			}
			// Update the text child
			for _, child := range node.Children {
				if child.Type == dom.NodeText {
					child.Data = sb.String()
					break
				}
			}
		} else {
			// No suggestions match — close the datalist
			a.closeDatalist()
		}
	}

	// Set placeholder-shown attribute for :placeholder-shown pseudo-class matching
	if tag == "input" || tag == "textarea" {
		hasPlaceholder := node.GetAttribute("placeholder") != ""
		val := a.formValues[node]
		if hasPlaceholder && val == "" {
			node.SetAttribute("placeholder-shown", "")
		} else {
			delete(node.Attributes, "placeholder-shown")
		}
	}

	// Set invalid attribute for :valid/:invalid pseudo-class matching
	if tag == "input" || tag == "textarea" || tag == "select" {
		isRequired := node.HasAttribute("required")
		val := a.formValues[node]
		hasValue := val != ""
		if tag == "select" {
			// Select always has a value (first option by default)
			hasValue = true
		}
		if isRequired && !hasValue {
			node.SetAttribute("invalid", "")
		} else {
			delete(node.Attributes, "invalid")
		}
	}

	// Apply animation classes — update element state based on animation frame
	if node.HasClass("animate-spin") {
		spinners := []string{"|", "/", "-", "\\"}
		ch := spinners[a.animFrame%len(spinners)]
		for _, child := range node.Children {
			if child.Type == dom.NodeText {
				if len(child.Data) > 0 {
					runes := []rune(child.Data)
					runes[len(runes)-1] = []rune(ch)[0]
					child.Data = string(runes)
				} else {
					child.Data = ch
				}
				break
			}
		}
	}
	if node.HasClass("animate-pulse") {
		if a.animFrame%4 < 2 {
			delete(node.Attributes, "pulsing")
		} else {
			node.SetAttribute("pulsing", "")
		}
	}
	if node.HasClass("animate-blink") {
		if a.animFrame%2 == 0 {
			node.SetAttribute("blinking", "")
		} else {
			delete(node.Attributes, "blinking")
		}
	}

	// Recurse into children
	for _, child := range node.Children {
		a.prepareFormDOM(child)
	}
}

// prepareProgressDOM renders a <progress> element as an ASCII progress bar.
func (a *App) prepareProgressDOM(node *dom.Node) {
	valueStr := node.GetAttribute("value")
	maxStr := node.GetAttribute("max")

	max := 100.0
	if maxStr != "" {
		if m, err := parseInt(maxStr); err == nil && m > 0 {
			max = float64(m)
		}
	}

	// Indeterminate: no value attribute
	if valueStr == "" {
		// Show indeterminate state
		a.setInputTextChild(node, " [∘∘∘∘∘∘∘∘∘∘]")
		return
	}

	value := 0.0
	if v, err := parseInt(valueStr); err == nil {
		value = float64(v)
	}
	if value < 0 {
		value = 0
	}
	if value > max {
		value = max
	}

	// Bar width is 10 characters
	barWidth := 10
	filled := int((value / max) * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "="
		} else if i == filled && filled < barWidth {
			bar += ">"
		} else {
			bar += " "
		}
	}

	pct := int((value / max) * 100)
	bar += "] "
	if value == max {
		bar += "100%"
	} else if pct < 10 {
		bar += "  " + itoa(pct) + "%"
	} else if pct < 100 {
		bar += " " + itoa(pct) + "%"
	} else {
		bar += itoa(pct) + "%"
	}

	a.setInputTextChild(node, bar)
}

// prepareMeterDOM renders a <meter> element as an ASCII gauge.
func (a *App) prepareMeterDOM(node *dom.Node) {
	valueStr := node.GetAttribute("value")
	minStr := node.GetAttribute("min")
	maxStr := node.GetAttribute("max")
	lowStr := node.GetAttribute("low")
	highStr := node.GetAttribute("high")
	optimumStr := node.GetAttribute("optimum")

	min := 0.0
	if minStr != "" {
		if m, err := parseInt(minStr); err == nil {
			min = float64(m)
		}
	}
	max := 1.0
	if maxStr != "" {
		if m, err := parseInt(maxStr); err == nil && m > 0 {
			max = float64(m)
		}
	}
	value := 0.0
	if valueStr != "" {
		if v, err := parseInt(valueStr); err == nil {
			value = float64(v)
		}
	}
	_ = lowStr
	_ = highStr
	_ = optimumStr

	if value < min {
		value = min
	}
	if value > max {
		value = max
	}

	range_ := max - min
	if range_ <= 0 {
		range_ = 1
	}

	barWidth := 10
	filled := int(((value - min) / range_) * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "]"

	a.setInputTextChild(node, bar)
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

// buildDatalistMap scans the DOM for <datalist> elements and builds a map of
// datalist id → option texts for use with <input list="id">.
func (a *App) buildDatalistMap(node *dom.Node) {
	if node == nil || node.Type != dom.NodeElement {
		return
	}
	tag := strings.ToLower(node.Data)
	if tag == "datalist" {
		id := node.GetAttribute("id")
		if id != "" {
			options := make([]string, 0)
			for _, child := range node.Children {
				if child.Type == dom.NodeElement && strings.ToLower(child.Data) == "option" {
					for _, tc := range child.Children {
						if tc.Type == dom.NodeText && strings.TrimSpace(tc.Data) != "" {
							options = append(options, strings.TrimSpace(tc.Data))
							break
						}
					}
				}
			}
			if len(options) > 0 {
				a.datalistMap[id] = options
			}
		}
	}
	for _, child := range node.Children {
		a.buildDatalistMap(child)
	}
}

// getDatalistForInput returns the datalist options linked to a given input,
// or nil if no datalist is linked.
func (a *App) getDatalistForInput(node *dom.Node) []string {
	listAttr := node.GetAttribute("list")
	if listAttr == "" {
		return nil
	}
	options, ok := a.datalistMap[listAttr]
	if !ok {
		return nil
	}
	return options
}

// handleDatalistEdit processes keyboard events for datalist autocomplete.
// It returns true if the event was consumed, false to fall through to normal text editing.
func (a *App) handleDatalistEdit(event terminal.Event, node *dom.Node) bool {
	options := a.getDatalistForInput(node)
	if options == nil {
		return false
	}

	switch event.Key {
	case terminal.KeyDown:
		// Open or navigate down in the datalist
		if !a.datalistOpen {
			// Open the dropdown with filtered suggestions
			a.openDatalist(node)
		} else {
			// Move highlight down
			a.datalistHighlight++
			if a.datalistHighlight >= len(a.datalistFiltered) {
				a.datalistHighlight = 0
			}
		}
		a.renderFrame()
		return true

	case terminal.KeyUp:
		if a.datalistOpen {
			if a.datalistHighlight <= 0 {
				// Close the dropdown if at top
				a.closeDatalist()
			} else {
				a.datalistHighlight--
			}
			a.renderFrame()
			return true
		}
		return false // fall through to normal cursor movement

	case terminal.KeyEnter:
		if a.datalistOpen && a.datalistHighlight >= 0 && a.datalistHighlight < len(a.datalistFiltered) {
			// Select the highlighted suggestion
			a.formValues[node] = a.datalistFiltered[a.datalistHighlight]
			a.formCursors[node] = len([]rune(a.formValues[node]))
			a.closeDatalist()
			a.renderFrame()
			return true
		}
		// Enter while datalist is closed — treat normally (fall through)
		return false

	case terminal.KeyEscape:
		if a.datalistOpen {
			a.closeDatalist()
			a.renderFrame()
			return true
		}
		return false

	default:
		// For typing keys, don't close the datalist — it will be re-filtered
		// in prepareFormDOM on the next render. Fall through to handleTextEdit.
		return false
	}
}

// openDatalist filters the datalist options by the current input value and shows them.
func (a *App) openDatalist(node *dom.Node) {
	options := a.getDatalistForInput(node)
	if options == nil {
		return
	}

	a.datalistInput = node
	a.datalistFiltered = nil

	// Filter by prefix match (case-insensitive)
	val := strings.ToLower(a.formValues[node])
	for _, opt := range options {
		if val == "" || strings.HasPrefix(strings.ToLower(opt), val) {
			a.datalistFiltered = append(a.datalistFiltered, opt)
		}
	}

	if len(a.datalistFiltered) == 0 {
		a.datalistFiltered = options
	}

	a.datalistHighlight = 0
	a.datalistOpen = true
}

// closeDatalist closes the datalist suggestion dropdown.
func (a *App) closeDatalist() {
	a.datalistOpen = false
	a.datalistInput = nil
	a.datalistFiltered = nil
	a.datalistHighlight = -1
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
		case "range":
			// Range slider: left/right arrows adjust, mouse drag adjusts
			return a.handleRangeEdit(event, focused)
		case "color":
			// Color input: left/right arrows cycle colors
			return a.handleColorEdit(event, focused)
		case "date":
			// Date input: arrow keys navigate segments, up/down adjust
			return a.handleDateEdit(event, focused)
		case "file":
			// File input: Enter/Space opens file picker
			return a.handleFileEdit(event, focused)
		default: // text, password, email, etc.
			return a.handleTextEdit(event, focused)
		}

	case "textarea":
		return a.handleTextEdit(event, focused)

	case "select":
		isMultiple := focused.HasAttribute("multiple")
		if event.Key == terminal.KeyUp || event.Key == terminal.KeyLeft {
			a.cycleSelectOption(focused, -1)
			a.renderFrame()
			return true
		}
		if event.Key == terminal.KeyDown || event.Key == terminal.KeyRight {
			a.cycleSelectOption(focused, 1)
			a.renderFrame()
			return true
		}
		if event.Key == terminal.KeyEnter || event.Key == terminal.KeySpace {
			if isMultiple {
				a.toggleMultiSelectOption(focused)
			}
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
	// Check for datalist autocomplete first
	if a.handleDatalistEdit(event, node) {
		return true
	}

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
		// Escape — lose focus (datalist closing is handled by handleDatalistEdit)
		a.closeDatalist()
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

// handleRangeEdit handles keyboard input for <input type="range">.
// Left/Right arrows adjust value by step; Home/End go to min/max.
func (a *App) handleRangeEdit(event terminal.Event, node *dom.Node) bool {
	val := a.formValues[node]
	if val == "" {
		val = "0"
	}

	stepStr := node.GetAttribute("step")
	step := 1
	if stepStr != "" {
		if s, err := parseInt(stepStr); err == nil && s > 0 {
			step = s
		}
	}

	minStr := node.GetAttribute("min")
	maxStr := node.GetAttribute("max")
	var minVal, maxVal int
	hasMin := false
	hasMax := false
	if minStr != "" {
		if v, err := parseInt(minStr); err == nil {
			minVal = v
			hasMin = true
		}
	}
	if maxStr != "" {
		if v, err := parseInt(maxStr); err == nil {
			maxVal = v
			hasMax = true
		}
	}
	// Default range if neither min nor max specified
	if !hasMin && !hasMax {
		hasMin = true
		hasMax = true
		minVal = 0
		maxVal = 100
	} else if !hasMin {
		minVal = maxVal - 100
		hasMin = true
	} else if !hasMax {
		maxVal = minVal + 100
		hasMax = true
	}

	n := 0
	if val != "" {
		n, _ = parseInt(val)
	}

	switch event.Key {
	case terminal.KeyLeft:
		n -= step
		if n < minVal {
			n = minVal
		}
		a.formValues[node] = itoa(n)
		a.renderFrame()
		return true

	case terminal.KeyRight:
		n += step
		if n > maxVal {
			n = maxVal
		}
		a.formValues[node] = itoa(n)
		a.renderFrame()
		return true

	case terminal.KeyHome:
		a.formValues[node] = itoa(minVal)
		a.renderFrame()
		return true

	case terminal.KeyEnd:
		a.formValues[node] = itoa(maxVal)
		a.renderFrame()
		return true

	default:
		// All other keys are ignored (no text editing for range sliders)
		return true
	}
}

// prepareRangeDOM renders <input type="range"> as a visual slider:
//   [===============-------------------] 45
func (a *App) prepareRangeDOM(node *dom.Node, val string) {
	if val == "" {
		val = "0"
	}

	// Parse range attributes
	minStr := node.GetAttribute("min")
	maxStr := node.GetAttribute("max")
	var minVal, maxVal float64
	var hasMin, hasMax bool
	if minStr != "" {
		if v, err := parseInt(minStr); err == nil {
			minVal = float64(v)
			hasMin = true
		}
	}
	if maxStr != "" {
		if v, err := parseInt(maxStr); err == nil {
			maxVal = float64(v)
			hasMax = true
		}
	}
	if !hasMin && !hasMax {
		minVal = 0
		maxVal = 100
	} else if !hasMin {
		minVal = maxVal - 100
	} else if !hasMax {
		maxVal = minVal + 100
	}

	// Parse current value
	curVal := 0.0
	if v, err := parseInt(val); err == nil {
		curVal = float64(v)
	}
	// Clamp
	if curVal < minVal {
		curVal = minVal
	}
	if curVal > maxVal {
		curVal = maxVal
	}

	// Compute display
	rangeSize := maxVal - minVal
	var fraction float64
	if rangeSize > 0 {
		fraction = (curVal - minVal) / rangeSize
	} else {
		fraction = 0
	}

	// Build the track: 20 chars wide
	trackWidth := 20
	filled := int(fraction * float64(trackWidth))
	// Clamp thumb position to valid range [0, trackWidth-1]
	if filled >= trackWidth {
		filled = trackWidth - 1
	}
	if filled < 0 {
		filled = 0
	}

	track := "["
	for i := 0; i < trackWidth; i++ {
		if i == filled {
			track += "o"
		} else if i < filled {
			track += "="
		} else {
			track += "-"
		}
	}
	track += "]"

	// Show value
	display := track + " " + itoa(int(curVal))
	a.setInputTextChild(node, display)
}

// prepareColorDOM renders <input type="color"> as the current hex value.
func (a *App) prepareColorDOM(node *dom.Node, val string) {
	if val == "" {
		val = "#000000"
	}
	// Validate hex format
	if len(val) < 7 || val[0] != '#' {
		val = "#000000"
	}
	// Show the color value
	a.setInputTextChild(node, val)
}

// handleColorEdit handles keyboard input for <input type="color">.
// Left/Right cycle through a curated palette of colors.
func (a *App) handleColorEdit(event terminal.Event, node *dom.Node) bool {
	// Curated palette of common web colors
	palette := []string{
		"#000000", "#444444", "#888888", "#BBBBBB", "#FFFFFF",
		"#FF0000", "#FF4444", "#FF8888",
		"#00FF00", "#44FF44", "#88FF88",
		"#0000FF", "#4444FF", "#8888FF",
		"#FFFF00", "#FFAA00", "#FF6600",
		"#00FFFF", "#00FFAA", "#00AAFF",
		"#FF00FF", "#AA00FF", "#FF0088",
		"#800000", "#008000", "#000080",
		"#808000", "#008080", "#800080",
		"#C0C0C0", "#A52A2A", "#2E8B57",
	}

	val := a.formValues[node]
	if val == "" {
		val = "#000000"
	}

	// Find current color index
	idx := -1
	for i, c := range palette {
		if strings.EqualFold(c, val) {
			idx = i
			break
		}
	}
	if idx < 0 {
		idx = 0
	}

	switch event.Key {
	case terminal.KeyLeft:
		idx--
		if idx < 0 {
			idx = len(palette) - 1
		}
		a.formValues[node] = palette[idx]
		a.renderFrame()
		return true

	case terminal.KeyRight:
		idx++
		if idx >= len(palette) {
			idx = 0
		}
		a.formValues[node] = palette[idx]
		a.renderFrame()
		return true

	case terminal.KeyHome:
		a.formValues[node] = palette[0]
		a.renderFrame()
		return true

	case terminal.KeyEnd:
		a.formValues[node] = palette[len(palette)-1]
		a.renderFrame()
		return true

	default:
		// All other keys are ignored
		return true
	}
}

// handleFileEdit handles keyboard input for <input type="file">.
// Enter/Space opens a file picker dialog; typing edits the path directly.
func (a *App) handleFileEdit(event terminal.Event, node *dom.Node) bool {
	switch event.Key {
	case terminal.KeyEnter, terminal.KeySpace:
		a.openFilePicker(node)
		return true
	default:
		return a.handleTextEdit(event, node)
	}
}

// openFilePicker opens a modal file browser dialog.
func (a *App) openFilePicker(node *dom.Node) {
	// Get current directory (from value or cwd)
	dir := a.formValues[node]
	if dir != "" {
		if fi, err := os.Stat(dir); err == nil && !fi.IsDir() {
			dir = filepath.Dir(dir)
		}
	} else {
		dir, _ = os.Getwd()
	}

	// Read directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		dir, _ = os.Getwd()
		entries, _ = os.ReadDir(dir)
	}

	// Build file list HTML
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<div id="file-picker" style="padding: 1; background-color: #1a1a2e; border: solid #0f3460;">
		<div style="color: #00d4aa; font-weight: bold; margin-bottom: 1;">📁 File Browser</div>
		<div style="color: #555; margin-bottom: 1;">%s</div>
		<div class="file-list">`, dir))

	// Parent directory
	sb.WriteString(`<div class="file-item" data-path=".." style="padding: 0 1; color: #00d4aa;">📁 ..</div>`)

	// Directories first
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			sb.WriteString(fmt.Sprintf(`<div class="file-item" data-path="%s" style="padding: 0 1; color: #00d4aa;">📁 %s</div>`, filepath.Join(dir, name), name))
		}
	}
	// Files
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			sb.WriteString(fmt.Sprintf(`<div class="file-item" data-path="%s" style="padding: 0 1;">📄 %s</div>`, filepath.Join(dir, name), name))
		}
	}

	sb.WriteString(`</div>
		<div style="margin-top: 1; color: #555;">↑↓ navigate · Enter select · Esc cancel</div>
	</div>`)

	// Show as modal
	a.ShowModal(sb.String())
	// On modal close, re-focus the form
	modalNode := a.modalNode
	_ = modalNode
	a.onModalClose = func() {
		a.formFocused = -1
		a.renderFrame()
	}

	// Hook into form events for the file picker
	a.formFocused = -1 // Temporarily unfocus to avoid conflicting with modal
}

// prepareDateDOM renders <input type="date"> with segment highlighting.
// Format: "2026-06-12" with year/month/day segments.
func (a *App) prepareDateDOM(node *dom.Node, val string) {
	if val == "" {
		now := time.Now()
		val = fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	}
	// Validate format YYYY-MM-DD
	if len(val) != 10 || val[4] != '-' || val[7] != '-' {
		now := time.Now()
		val = fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	}
	// Show date with segment indicators when focused
	isFocused := a.formFocused >= 0 && a.formFocused < len(a.formFocusables) && a.formFocusables[a.formFocused] == node
	if isFocused {
		// Get the active segment from cursor or default to 0
		segment := a.formCursors[node]
		if segment < 0 || segment > 2 {
			segment = 0
		}
		// Parse segments
		parts := strings.Split(val, "-")
		if len(parts) != 3 {
			a.setInputTextChild(node, val)
			return
		}
		// Bold brackets around active segment
		display := ""
		for i, part := range parts {
			if i > 0 {
				display += "-"
			}
			if i == segment {
				display += "[" + part + "]"
			} else {
				display += part
			}
		}
		a.setInputTextChild(node, display)
	} else {
		a.setInputTextChild(node, val)
	}
}

// handleDateEdit handles keyboard input for <input type="date">.
// Navigation between year/month/day segments, up/down to adjust values.
func (a *App) handleDateEdit(event terminal.Event, node *dom.Node) bool {
	val := a.formValues[node]
	if val == "" {
		now := time.Now()
		val = fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
		a.formValues[node] = val
	}

	// Parse current date
	parts := strings.Split(val, "-")
	if len(parts) != 3 {
		return true
	}
	year, _ := parseInt(parts[0])
	month, _ := parseInt(parts[1])
	day, _ := parseInt(parts[2])

	// Get current segment (0=year, 1=month, 2=day)
	segment := a.formCursors[node]
	if segment < 0 || segment > 2 {
		segment = 0
	}

	switch event.Key {
	case terminal.KeyLeft:
		segment--
		if segment < 0 {
			segment = 2
		}
		a.formCursors[node] = segment
		a.renderFrame()
		return true

	case terminal.KeyRight:
		segment++
		if segment > 2 {
			segment = 0
		}
		a.formCursors[node] = segment
		a.renderFrame()
		return true

	case terminal.KeyUp:
		switch segment {
		case 0: // year
			year++
		case 1: // month
			month++
			if month > 12 {
				month = 1
			}
		case 2: // day
			day++
			maxDay := daysInMonth(year, month)
			if day > maxDay {
				day = 1
			}
		}
		a.formValues[node] = fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		a.renderFrame()
		return true

	case terminal.KeyDown:
		switch segment {
		case 0: // year
			year--
		case 1: // month
			month--
			if month < 1 {
				month = 12
			}
		case 2: // day
			day--
			if day < 1 {
				maxDay := daysInMonth(year, month)
				day = maxDay
			}
		}
		a.formValues[node] = fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		a.renderFrame()
		return true

	default:
		// All other keys are ignored
		return true
	}
}

// daysInMonth returns the number of days in a given month/year.
func daysInMonth(year, month int) int {
	switch time.Month(month) {
	case time.January, time.March, time.May, time.July,
		time.August, time.October, time.December:
		return 31
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		if isLeapYear(year) {
			return 29
		}
		return 28
	}
	return 30
}

// isLeapYear returns true if the given year is a leap year.
func isLeapYear(year int) bool {
	return year%400 == 0 || (year%4 == 0 && year%100 != 0)
}

// handleRangeClick handles a mouse click on a range input.
// It computes the value from the mouse X position relative to the slider track.
func (a *App) handleRangeClick(mouseX int, box *layout.Box) {
	if box == nil || box.Node == nil {
		return
	}
	tag := strings.ToLower(box.Node.Data)
	if tag != "input" {
		return
	}
	inputType := strings.ToLower(box.Node.GetAttribute("type"))
	if inputType != "range" {
		return
	}

	val := a.computeRangeValueFromX(box.Node, box, mouseX)
	a.formValues[box.Node] = itoa(val)
	a.dragTarget = box.Node
}

// dragRangeByMouse handles mouse drag on a range input.
func (a *App) dragRangeByMouse(mouseX int) {
	if a.dragTarget == nil || a.rootBox == nil {
		return
	}
	val := a.computeRangeValueFromX(a.dragTarget, a.findBoxForNode(a.dragTarget), mouseX)
	a.formValues[a.dragTarget] = itoa(val)
}

// computeRangeValueFromX computes the range value from a mouse X position on the track.
func (a *App) computeRangeValueFromX(node *dom.Node, box *layout.Box, mouseX int) int {
	if box == nil {
		return 0
	}
	// Parse range attributes
	minStr := node.GetAttribute("min")
	maxStr := node.GetAttribute("max")
	minVal := 0
	maxVal := 100
	if minStr != "" {
		if v, err := parseInt(minStr); err == nil {
			minVal = v
		}
	}
	if maxStr != "" {
		if v, err := parseInt(maxStr); err == nil {
			maxVal = v
		}
	}

	// The track starts inside the content area, 1 cell after the left edge (for the "[" char)
	contentX := box.ContentRect.X
	// The track display is "[" + trackWidth chars + "]" + " " + valueText
	trackWidth := 20 // must match prepareRangeDOM
	trackStart := contentX + 1 // after the "["

	// Compute fraction based on mouse X within the track
	relX := mouseX - trackStart
	fraction := float64(relX) / float64(trackWidth)
	if fraction < 0 {
		fraction = 0
	}
	if fraction > 1 {
		fraction = 1
	}

	rangeSize := maxVal - minVal
	val := minVal + int(fraction*float64(rangeSize)+0.5) // round
	if val < minVal {
		val = minVal
	}
	if val > maxVal {
		val = maxVal
	}
	return val
}

// findBoxForNode finds the layout box that corresponds to a given DOM node.
func (a *App) findBoxForNode(node *dom.Node) *layout.Box {
	if a.rootBox == nil || node == nil {
		return nil
	}
	return a.findBoxInTree(a.rootBox, node)
}

func (a *App) findBoxInTree(box *layout.Box, node *dom.Node) *layout.Box {
	if box == nil {
		return nil
	}
	if box.Node == node {
		return box
	}
	for _, child := range box.Children {
		if found := a.findBoxInTree(child, node); found != nil {
			return found
		}
	}
	return nil
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
	// Close any open datalist when focus changes
	a.closeDatalist()
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

// collectSelectOptions collects all <option> text values from a <select> element,
// including options nested inside <optgroup> elements.
func collectSelectOptions(node *dom.Node) []string {
	var options []string
	var walk func(*dom.Node)
	walk = func(n *dom.Node) {
		if n == nil || n.Type != dom.NodeElement {
			return
		}
		tag := strings.ToLower(n.Data)
		if tag == "option" {
			for _, tc := range n.Children {
				if tc.Type == dom.NodeText {
					options = append(options, tc.Data)
					break
				}
			}
		} else if tag == "optgroup" {
			// Recurse into optgroup for its option children
			for _, child := range n.Children {
				walk(child)
			}
		}
	}
	for _, child := range node.Children {
		walk(child)
	}
	return options
}

// cycleSelectOption moves to the next or previous option in a select element.
// dir should be 1 for next, -1 for previous.
func (a *App) cycleSelectOption(node *dom.Node, dir int) {
	options := collectSelectOptions(node)
	if len(options) == 0 {
		return
	}

	current := a.formValues[node]
	for i, opt := range options {
		if opt == current {
			next := (i + dir) % len(options)
			if next < 0 {
				next += len(options)
			}
			a.formValues[node] = options[next]
			return
		}
	}
	// Current value not found — set to first option
	a.formValues[node] = options[0]
}

// toggleMultiSelectOption toggles the selection state of the currently highlighted
// option in a <select multiple> element.
func (a *App) toggleMultiSelectOption(node *dom.Node) {
	current := a.formValues[node]

	// Find the current option and toggle it
	var toggledOpt string
	options := collectSelectOptions(node)
	for _, opt := range options {
		if opt == current {
			toggledOpt = opt
			break
		}
	}

	if toggledOpt == "" && len(options) > 0 {
		// No current selection — toggle first option
		toggledOpt = options[0]
	}

	if toggledOpt == "" {
		return
	}

	// Toggle the option in the pipe-separated list
	selected := make(map[string]bool)
	if a.formValues[node] != "" {
		for _, s := range strings.Split(a.formValues[node], "|") {
			selected[strings.TrimSpace(s)] = true
		}
	}

	if selected[toggledOpt] {
		delete(selected, toggledOpt)
	} else {
		selected[toggledOpt] = true
	}

	// Rebuild the value string
	var sb strings.Builder
	first := true
	for _, opt := range options {
		if selected[opt] {
			if !first {
				sb.WriteString("|")
			}
			sb.WriteString(opt)
			first = false
		}
	}
	a.formValues[node] = sb.String()
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
