//go:build unix

package terminal

import (
	"golang.org/x/sys/unix"
)

func EnableRawMode(fd uintptr) (unix.Termios, error) {
	orig, err := unix.IoctlGetTermios(int(fd), ioctlGetTermios)
	if err != nil {
		return unix.Termios{}, err
	}

	raw := *orig
	// Input: disable BRKINT, ICRNL, INPCK, ISTRIP, IXON
	raw.Iflag &^= unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON
	// Output: disable OPOST
	raw.Oflag &^= unix.OPOST
	// Control: set CS8
	raw.Cflag |= unix.CS8
	// Local: disable ECHO, ICANON, IEXTEN, ISIG
	raw.Lflag &^= unix.ECHO | unix.ICANON | unix.IEXTEN | unix.ISIG
	// Read returns after 1 byte
	raw.Cc[unix.VMIN] = 1
	raw.Cc[unix.VTIME] = 0

	if err := unix.IoctlSetTermios(int(fd), ioctlSetTermios, &raw); err != nil {
		return unix.Termios{}, err
	}
	return *orig, nil
}

func DisableRawMode(fd uintptr, orig *unix.Termios) error {
	return unix.IoctlSetTermios(int(fd), ioctlSetTermios, orig)
}
