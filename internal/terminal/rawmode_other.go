//go:build !unix

package terminal

import "errors"

type termiosPlaceholder struct{}

func enableRawMode(fd uintptr) (termiosPlaceholder, error) {
	return termiosPlaceholder{}, errors.New("raw mode not supported on this platform")
}

func disableRawMode(fd uintptr, orig *termiosPlaceholder) error {
	return errors.New("raw mode not supported on this platform")
}
