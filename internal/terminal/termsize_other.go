//go:build !unix

package terminal

// TermSize returns a default 24x80 on non-unix platforms.
func TermSize(fd uintptr) (rows, cols int) {
	return 24, 80
}
