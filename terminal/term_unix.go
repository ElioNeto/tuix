//go:build linux || darwin || freebsd || openbsd || netbsd

package terminal

import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

// makeRaw sets the terminal to raw mode and returns the previous state.
func makeRaw(fd uintptr) (*State, error) {
	var oldState syscall.Termios
	if err := tcget(fd, &oldState); err != nil {
		return nil, err
	}

	newState := oldState

	// Apply raw mode settings
	newState.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK |
		syscall.ISTRIP | syscall.INLCR | syscall.IGNCR |
		syscall.ICRNL | syscall.IXON
	newState.Oflag &^= syscall.OPOST
	newState.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON |
		syscall.ISIG | syscall.IEXTEN
	newState.Cflag &^= syscall.CSIZE | syscall.PARENB
	newState.Cflag |= syscall.CS8
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

	if err := tcset(fd, &newState); err != nil {
		return nil, err
	}

	// Store old state
	buf := make([]byte, unsafe.Sizeof(oldState))
	copy(buf, (*[1 << 30]byte)(unsafe.Pointer(&oldState))[:unsafe.Sizeof(oldState)])
	var stateBuf [60]byte
	copy(stateBuf[:], buf)
	return &State{cc: stateBuf}, nil
}

// restore resets the terminal to a previous state.
func restore(fd uintptr, state *State) {
	var oldState syscall.Termios
	copy((*[1 << 30]byte)(unsafe.Pointer(&oldState))[:unsafe.Sizeof(oldState)], state.cc[:])
	tcset(fd, &oldState)
}

// getSize returns the terminal width and height.
func getSize(fd uintptr) (width, height int, err error) {
	ws, err := ioctlWinsize(fd)
	if err != nil {
		return 80, 24, err
	}
	return int(ws.Col), int(ws.Row), nil
}

// tcget retrieves terminal parameters.
func tcget(fd uintptr, termios *syscall.Termios) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd,
		syscall.TCGETS, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}

// tcset sets terminal parameters.
func tcset(fd uintptr, termios *syscall.Termios) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd,
		syscall.TCSETS, uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}

// Winsize holds terminal window size information.
type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// ioctlWinsize retrieves terminal window size.
func ioctlWinsize(fd uintptr) (*Winsize, error) {
	ws := &Winsize{}
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd,
		syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(ws)), 0, 0, 0)
	if err != 0 {
		return nil, err
	}
	return ws, nil
}

// watchResize listens for SIGWINCH signals (Unix-only).
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
