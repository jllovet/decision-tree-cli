//go:build !unix

package terminal

func isTerminalFd(fd uintptr) bool {
	return false
}
