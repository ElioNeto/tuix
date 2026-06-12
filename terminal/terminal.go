// Package terminal provides low-level terminal manipulation for TUI applications.
//
// It handles raw mode, terminal size queries, input events, and ANSI escape code
// output. All terminal I/O is done directly via syscalls, with no external dependencies.
package terminal

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Terminal manages the raw terminal state and provides input/output channels.
type Terminal struct {
	in        *os.File
	out       *os.File
	oldState  *State
	events    chan Event
	done      chan struct{}
	resize    chan struct{}
	mu        sync.Mutex
	closed    bool
	wg        sync.WaitGroup
	colorMode int // 0=16, 1=256, 2=truecolor
}

// State holds the terminal configuration for restoring later.
type State struct {
	cc [60]byte // termios structure as raw bytes
}

// Event represents a terminal event (key press, mouse, resize, etc.).
type Event struct {
	Type EventType
	// Key event fields
	Key       Key
	Rune      rune
	Modifiers Modifier
	// Mouse event fields
	MouseButton MouseButton
	MouseX      int
	MouseY      int
	// Resize event fields
	Width  int
	Height int
}

// EventType categorizes terminal events.
type EventType int

const (
	EventKey    EventType = iota // Keyboard event
	EventMouse                   // Mouse event
	EventResize                  // Terminal resize event
	EventError                   // An error occurred
)

// Key represents a keyboard key (physical or control).
type Key int

const (
	KeyNone       Key = 0
	KeyRune       Key = iota + 256 // A Unicode rune
	KeyBackspace
	KeyTab
	KeyEnter
	KeyEscape
	KeySpace
	KeyDelete
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyInsert
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyCtrlA
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlG
	KeyCtrlH
	KeyCtrlI
	KeyCtrlJ
	KeyCtrlK
	KeyCtrlL
	KeyCtrlM
	KeyCtrlN
	KeyCtrlO
	KeyCtrlP
	KeyCtrlQ
	KeyCtrlR
	KeyCtrlS
	KeyCtrlT
	KeyCtrlU
	KeyCtrlV
	KeyCtrlW
	KeyCtrlX
	KeyCtrlY
	KeyCtrlZ
	KeyCtrlOpenBracket  // ESC / Ctrl-[
	KeyCtrlBackslash    // Ctrl-\
	KeyCtrlCloseBracket // Ctrl-]
)

// Modifier represents keyboard modifiers.
type Modifier int

const (
	ModNone  Modifier = 0
	ModShift Modifier = 1 << iota
	ModAlt
	ModCtrl
)

// MouseButton represents a mouse button.
type MouseButton int

const (
	MouseNone     MouseButton = 0
	MouseLeft     MouseButton = 1
	MouseMiddle   MouseButton = 2
	MouseRight    MouseButton = 3
	MouseWheelUp  MouseButton = 4
	MouseWheelDown MouseButton = 5
)

// Open opens the terminal for raw I/O.
func Open() (*Terminal, error) {
	t := &Terminal{
		in:        os.Stdin,
		out:       os.Stdout,
		events:    make(chan Event, 64),
		done:      make(chan struct{}),
		resize:    make(chan struct{}, 1),
		colorMode: detectColorMode(),
	}

	// Save current terminal state
	state, err := makeRaw(t.in.Fd())
	if err != nil {
		return nil, err
	}
	t.oldState = state

	// Enter alternate screen buffer
	t.enterAltScreen()

	// Hide cursor
	t.hideCursor()

	// Enable mouse tracking
	t.enableMouse()

	// Start event reader goroutine
	t.wg.Add(1)
	go t.readEvents()

	// Watch for SIGWINCH
	t.wg.Add(1)
	go t.watchResize()

	return t, nil
}

// Close restores the terminal to its original state and stops event processing.
func (t *Terminal) Close() error {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return nil
	}
	t.closed = true
	close(t.done)
	t.mu.Unlock()

	// Wait for goroutines to finish
	t.wg.Wait()

	// Show cursor
	t.showCursor()

	// Disable mouse tracking
	t.disableMouse()

	// Exit alternate screen buffer
	t.exitAltScreen()

	// Restore terminal state
	if t.oldState != nil {
		restore(t.in.Fd(), t.oldState)
	}

	return nil
}

// Events returns a channel of terminal events.
func (t *Terminal) Events() <-chan Event {
	return t.events
}

// Size returns the current terminal size in columns and rows.
func (t *Terminal) Size() (width, height int, err error) {
	return getSize(t.out.Fd())
}

// ColorMode returns the terminal's supported color mode:
// 0 = 16 colors, 1 = 256 colors, 2 = true color.
func (t *Terminal) ColorMode() int {
	return t.colorMode
}

// WriteString writes a string directly to the terminal output.
func (t *Terminal) WriteString(s string) (int, error) {
	return io.WriteString(t.out, s)
}

// Write writes bytes directly to the terminal output.
func (t *Terminal) Write(b []byte) (int, error) {
	return t.out.Write(b)
}

// readEvents reads and parses terminal input events.
func (t *Terminal) readEvents() {
	defer t.wg.Done()

	buf := make([]byte, 128)
	for {
		select {
		case <-t.done:
			return
		default:
		}

		n, err := t.in.Read(buf)
		if err != nil {
			select {
			case t.events <- Event{Type: EventError}:
			default:
			}
			return
		}

		t.parseInput(buf[:n])
	}
}

// parseInput parses raw terminal bytes into events.
func (t *Terminal) parseInput(data []byte) {
	if len(data) == 0 {
		return
	}

	// Handle escape sequences
	if data[0] == '\x1b' {
		t.parseEscapeSequence(data)
		return
	}

	// Handle regular characters and control characters
	for _, b := range data {
		switch {
		case b == '\r' || b == '\n':
			t.emitKey(KeyEnter, 0)
		case b == '\x7f' || b == '\b':
			t.emitKey(KeyBackspace, 0)
		case b == '\t':
			t.emitKey(KeyTab, 0)
		case b == '\x1b':
			t.emitKey(KeyEscape, 0)
		case b < 0x20:
			// Ctrl+letter
			t.emitKey(Key(b), ModCtrl)
		case b >= 0x20 && b <= 0x7e:
			// Printable ASCII
			t.emitRune(rune(b))
		default:
			// Multi-byte UTF-8 rune
			t.emitRune(rune(b))
		}
	}
}

// parseEscapeSequence parses ANSI escape sequences.
func (t *Terminal) parseEscapeSequence(data []byte) {
	if len(data) < 2 {
		t.emitKey(KeyEscape, 0)
		return
	}

	switch data[1] {
	case '[': // CSI sequences
		t.parseCSI(data[2:])
	case 'O': // SS3 sequences (F1-F4)
		t.parseSS3(data[2:])
	case ']': // OSC sequences
		// Ignore for now
	default:
		t.emitKey(KeyEscape, 0)
	}
}

// parseCSI parses Control Sequence Introducer (CSI) sequences.
func (t *Terminal) parseCSI(data []byte) {
	if len(data) == 0 {
		return
	}

	// Read parameters (digits separated by semicolons)
	params := make([]int, 0)
	i := 0
	for i < len(data) && data[i] >= '0' && data[i] <= '9' || data[i] == ';' {
		if data[i] == ';' {
			i++
			continue
		}
		// Read a number
		start := i
		for i < len(data) && data[i] >= '0' && data[i] <= '9' {
			i++
		}
		if start < i {
			val := 0
			for j := start; j < i; j++ {
				val = val*10 + int(data[j]-'0')
			}
			params = append(params, val)
		}
	}

	// Need at least one byte for the final character
	if i >= len(data) {
		return
	}

	final := data[i]

	// Default for missing params
	if len(params) == 0 {
		params = []int{0}
	}

	switch final {
	case 'A': // Up arrow
		t.emitKey(KeyUp, getModifier(params))
	case 'B': // Down arrow
		t.emitKey(KeyDown, getModifier(params))
	case 'C': // Right arrow
		t.emitKey(KeyRight, getModifier(params))
	case 'D': // Left arrow
		t.emitKey(KeyLeft, getModifier(params))
	case 'H': // Home
		t.emitKey(KeyHome, getModifier(params))
	case 'F': // End
		t.emitKey(KeyEnd, getModifier(params))
	case '~': // Extended keys
		if len(params) > 0 {
			t.parseExtendedKey(params)
		}
	case 'M': // Mouse event (old X10 encoding)
		t.parseMouseOld(data[i+1:])
	case '<': // SGR mouse event
		t.parseMouseSGR(data[i+1:])
	case 'Z': // Shift+Tab
		t.emitKey(KeyTab, ModShift)
	}
}

// parseSS3 parses SS3 sequences (used by F1-F4 on some terminals).
func (t *Terminal) parseSS3(data []byte) {
	if len(data) == 0 {
		return
	}

	switch data[0] {
	case 'P':
		t.emitKey(KeyF1, 0)
	case 'Q':
		t.emitKey(KeyF2, 0)
	case 'R':
		t.emitKey(KeyF3, 0)
	case 'S':
		t.emitKey(KeyF4, 0)
	default:
		t.emitKey(KeyEscape, 0)
	}
}

// parseExtendedKey parses extended key sequences (ending with ~).
func (t *Terminal) parseExtendedKey(params []int) {
	if len(params) == 0 {
		return
	}

	switch params[0] {
	case 1, 7:
		t.emitKey(KeyHome, getModifier(params))
	case 2:
		t.emitKey(KeyInsert, getModifier(params))
	case 3:
		t.emitKey(KeyDelete, getModifier(params))
	case 4, 8:
		t.emitKey(KeyEnd, getModifier(params))
	case 5:
		t.emitKey(KeyPageUp, getModifier(params))
	case 6:
		t.emitKey(KeyPageDown, getModifier(params))
	case 11:
		t.emitKey(KeyF1, getModifier(params))
	case 12:
		t.emitKey(KeyF2, getModifier(params))
	case 13:
		t.emitKey(KeyF3, getModifier(params))
	case 14:
		t.emitKey(KeyF4, getModifier(params))
	case 15:
		t.emitKey(KeyF5, getModifier(params))
	case 17:
		t.emitKey(KeyF6, getModifier(params))
	case 18:
		t.emitKey(KeyF7, getModifier(params))
	case 19:
		t.emitKey(KeyF8, getModifier(params))
	case 20:
		t.emitKey(KeyF9, getModifier(params))
	case 21:
		t.emitKey(KeyF10, getModifier(params))
	case 23:
		t.emitKey(KeyF11, getModifier(params))
	case 24:
		t.emitKey(KeyF12, getModifier(params))
	}
}

// parseMouseOld parses X10 mouse encoding.
func (t *Terminal) parseMouseOld(data []byte) {
	if len(data) < 2 {
		return
	}
	btn := int(data[0]) & 0x03
	x := int(data[1]) - 32
	y := int(data[2]) - 32

	var mb MouseButton
	switch btn {
	case 0:
		mb = MouseLeft
	case 1:
		mb = MouseMiddle
	case 2:
		mb = MouseRight
	case 3:
		// Release - ignore
		return
	}

	select {
	case t.events <- Event{Type: EventMouse, MouseButton: mb, MouseX: x, MouseY: y}:
	default:
	}
}

// parseMouseSGR parses SGR-encoded mouse events.
func (t *Terminal) parseMouseSGR(data []byte) {
	// Expected format: <btn;x;y[Mm]
	str := string(data)
	var btn, x, y int
	n, err := fmt.Sscanf(str, "%d;%d;%d", &btn, &x, &y)
	if err != nil || n < 3 {
		return
	}

	// Check for trailing 'M' (press) or 'm' (release)
	isRelease := false
	for _, ch := range str {
		if ch == 'm' {
			isRelease = true
		}
	}

	var mb MouseButton
	isMotion := btn&0x20 != 0
	switch btn & 0x03 {
	case 0:
		if isMotion {
			// Motion with no button = pure mouse move
			// If bit 5 (motion) is set AND bits 0-1 are 0 with mode 1003,
			// this could be a drag if a button was previously held.
			// We treat it as MouseNone (pure move) and let the app detect drag
			// via the combination of mouse press + subsequent moves.
			mb = MouseNone
		} else {
			mb = MouseLeft
		}
	case 1:
		if isMotion {
			mb = MouseNone // middle button drag — treated as move
		} else {
			mb = MouseMiddle
		}
	case 2:
		if isMotion {
			mb = MouseNone // right button drag — treated as move
		} else {
			mb = MouseRight
		}
	}

	// Check for mouse wheel
	if btn&0x40 != 0 {
		if btn&0x01 != 0 {
			mb = MouseWheelDown
		} else {
			mb = MouseWheelUp
		}
	}

	if !isRelease {
		select {
		case t.events <- Event{Type: EventMouse, MouseButton: mb, MouseX: x, MouseY: y}:
		default:
		}
	}
}

// getModifier extracts the modifier from CSI parameters.
// In many terminals, the second param or a param >= 27 holds the modifier.
func getModifier(params []int) Modifier {
	if len(params) < 2 {
		return ModNone
	}
	mod := params[1]
	switch mod {
	case 2:
		return ModShift
	case 3:
		return ModAlt
	case 4:
		return ModShift | ModAlt
	case 5:
		return ModCtrl
	case 6:
		return ModShift | ModCtrl
	case 7:
		return ModAlt | ModCtrl
	case 8:
		return ModShift | ModAlt | ModCtrl
	}
	return ModNone
}

func (t *Terminal) emitKey(key Key, mod Modifier) {
	select {
	case t.events <- Event{Type: EventKey, Key: key, Modifiers: mod}:
	default:
	}
}

func (t *Terminal) emitRune(r rune) {
	select {
	case t.events <- Event{Type: EventKey, Key: KeyRune, Rune: r}:
	default:
	}
}

// watchResize listens for SIGWINCH signals.
func (t *Terminal) watchResize() {
	defer t.wg.Done()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)
	defer signal.Stop(sig)

	for {
		select {
		case <-t.done:
			return
		case <-sig:
			w, h, err := t.Size()
			if err == nil {
				select {
				case t.events <- Event{Type: EventResize, Width: w, Height: h}:
				default:
				}
			}
		}
	}
}

// enterAltScreen switches to the alternate screen buffer.
func (t *Terminal) enterAltScreen() {
	t.WriteString("\x1b[?1049h")
}

// exitAltScreen returns to the main screen buffer.
func (t *Terminal) exitAltScreen() {
	t.WriteString("\x1b[?1049l")
}

// hideCursor hides the terminal cursor.
func (t *Terminal) hideCursor() {
	t.WriteString("\x1b[?25l")
}

// showCursor shows the terminal cursor.
func (t *Terminal) showCursor() {
	t.WriteString("\x1b[?25h")
}

// clearScreen clears the entire screen and moves cursor home.
func (t *Terminal) clearScreen() {
	t.WriteString("\x1b[2J\x1b[H")
}

// enableMouse enables mouse tracking (X10 + SGR + motion).
func (t *Terminal) enableMouse() {
	t.WriteString("\x1b[?1000h") // Enable mouse tracking
	t.WriteString("\x1b[?1002h") // Enable button event tracking
	t.WriteString("\x1b[?1003h") // Enable motion event tracking (hover)
	t.WriteString("\x1b[?1006h") // Enable SGR mouse
}

// disableMouse disables mouse tracking.
func (t *Terminal) disableMouse() {
	t.WriteString("\x1b[?1006l")
	t.WriteString("\x1b[?1003l")
	t.WriteString("\x1b[?1002l")
	t.WriteString("\x1b[?1000l")
}

// detectColorMode detects the terminal's color capabilities by checking
// environment variables and terminfo database.
func detectColorMode() int {
	term := os.Getenv("TERM")
	colorterm := os.Getenv("COLORTERM")

	// True color support
	if colorterm == "truecolor" || colorterm == "24bit" || colorterm == "yes" {
		return 2
	}

	// 256 color support
	switch term {
	case "xterm-256color", "screen-256color", "tmux-256color",
		"xterm-kitty", "alacritty", "wezterm":
		return 1
	}

	// Check COLORTERM fallback
	if colorterm != "" {
		return 1
	}

	return 0
}
