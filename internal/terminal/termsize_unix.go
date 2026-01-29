//go:build unix

package terminal

import "golang.org/x/sys/unix"

// TermSize returns the terminal dimensions (rows, cols) for the given fd.
// Falls back to 24x80 on error.
func TermSize(fd uintptr) (rows, cols int) {
	ws, err := unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ)
	if err != nil || ws.Row == 0 || ws.Col == 0 {
		return 24, 80
	}
	return int(ws.Row), int(ws.Col)
}
