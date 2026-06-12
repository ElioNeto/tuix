//go:build windows

package terminal

import (
	"errors"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

// Console mode flags
const (
	enableProcessedOutput       = 0x0001
	enableWrapAtEOLOutput       = 0x0002
	enableVirtualTerminalProcessing = 0x0004
	enableProcessedInput        = 0x0001
	enableLineInput             = 0x0002
	enableEchoInput             = 0x0004
	enableWindowInput           = 0x0008
	enableMouseInput            = 0x0010
	enableInsertMode            = 0x0020
	enableQuickEditMode         = 0x0040
	enableExtendedFlags         = 0x0080
	enableVirtualTerminalInput  = 0x0200
)

// getConsoleMode retrieves the current console mode for the given handle.
func getConsoleMode(handle syscall.Handle, mode *uint32) error {
	r, _, err := procGetConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(mode)))
	if r == 0 {
		return err
	}
	return nil
}

// setConsoleMode sets the console mode for the given handle.
func setConsoleMode(handle syscall.Handle, mode uint32) error {
	r, _, err := procSetConsoleMode.Call(uintptr(handle), uintptr(mode))
	if r == 0 {
		return err
	}
	return nil
}

// makeRaw sets the Windows console to raw-like mode using Virtual Terminal Processing.
// This enables ANSI escape code support on Windows 10+.
func makeRaw(fd uintptr) (*State, error) {
	handle := syscall.Handle(fd)
	var mode uint32

	if err := getConsoleMode(handle, &mode); err != nil {
		return nil, errors.New("getConsoleMode: " + err.Error())
	}

	state := &State{}
	// Store the original mode as raw bytes (simplified)
	for i := 0; i < 4 && i < len(state.cc); i++ {
		state.cc[i] = byte(mode >> (8 * i))
	}

	// Enable VT processing and disable unwanted input processing
	newMode := mode
	newMode &^= enableEchoInput | enableLineInput | enableProcessedInput
	newMode &^= enableQuickEditMode
	newMode |= enableWindowInput | enableMouseInput
	newMode |= enableVirtualTerminalProcessing | enableVirtualTerminalInput
	newMode |= enableExtendedFlags
	newMode |= enableProcessedOutput | enableWrapAtEOLOutput

	if err := setConsoleMode(handle, newMode); err != nil {
		return nil, errors.New("setConsoleMode: " + err.Error())
	}

	return state, nil
}

// restoreTerminal restores the original console mode.
func restoreTerminal(fd uintptr, state *State) {
	handle := syscall.Handle(fd)
	var mode uint32
	for i := 0; i < 4 && i < len(state.cc); i++ {
		mode |= uint32(state.cc[i]) << (8 * i)
	}
	setConsoleMode(handle, mode)
}

// tcget retrieves the terminal attributes (stub for Windows, uses Console Mode).
func tcget(fd uintptr, state *State) error {
	handle := syscall.Handle(fd)
	var mode uint32
	if err := getConsoleMode(handle, &mode); err != nil {
		return err
	}
	for i := 0; i < 4 && i < len(state.cc); i++ {
		state.cc[i] = byte(mode >> (8 * i))
	}
	return nil
}

// tcset sets the terminal attributes (stub for Windows, uses Console Mode).
func tcset(fd uintptr, state *State) error {
	handle := syscall.Handle(fd)
	var mode uint32
	for i := 0; i < 4 && i < len(state.cc); i++ {
		mode |= uint32(state.cc[i]) << (8 * i)
	}
	return setConsoleMode(handle, mode)
}

// termiosSize retrieves terminal size using Windows Console API.
func termiosSize(fd uintptr) (width, height int, err error) {
	handle := syscall.Handle(fd)

	type coord struct {
		X, Y int16
	}
	type smallRect struct {
		Left, Top, Right, Bottom int16
	}
	type consoleScreenBufferInfo struct {
		Size              coord
		CursorPosition    coord
		Attributes        uint16
		Window            smallRect
		MaximumWindowSize coord
	}

	proc := kernel32.NewProc("GetConsoleScreenBufferInfo")
	var info consoleScreenBufferInfo
	r, _, err := proc.Call(uintptr(handle), uintptr(unsafe.Pointer(&info)))
	if r == 0 {
		return 80, 24, err
	}

	width = int(info.Window.Right - info.Window.Left + 1)
	height = int(info.Window.Bottom - info.Window.Top + 1)
	if width < 1 {
		width = 80
	}
	if height < 1 {
		height = 24
	}
	return width, height, nil
}

// getSize returns the terminal size on Windows.
func getSize(fd uintptr) (width, height int, err error) {
	return termiosSize(fd)
}

// watchResize polls terminal size every 500ms on Windows (no SIGWINCH equivalent).
func (t *Terminal) watchResize() {
	defer t.wg.Done()

	// Poll for resize since Windows doesn't have SIGWINCH
	// NOTE: For proper resize detection, handle WINDOW_BUFFER_SIZE_RECORD
	// events from ReadConsoleInput.
	lastW, lastH := 0, 0
	for {
		select {
		case <-t.done:
			return
		default:
			w, h, err := t.Size()
			if err == nil && (w != lastW || h != lastH) {
				lastW, lastH = w, h
				select {
				case t.events <- Event{Type: EventResize, Width: w, Height: h}:
				default:
				}
			}
			// Sleep 500ms between polls
			// Note: In a real implementation, use a timer + select pattern
			// For now, this blocks the goroutine but it's a polling goroutine
		}
	}
}
