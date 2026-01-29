//go:build unix

package terminal

import "golang.org/x/sys/unix"

func isTerminalFd(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), ioctlGetTermios)
	return err == nil
}
