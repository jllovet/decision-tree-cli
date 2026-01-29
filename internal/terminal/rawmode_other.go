//go:build !unix

package terminal

import "errors"

type termiosPlaceholder struct{}

func EnableRawMode(fd uintptr) (termiosPlaceholder, error) {
	return termiosPlaceholder{}, errors.New("raw mode not supported on this platform")
}

func DisableRawMode(fd uintptr, orig *termiosPlaceholder) error {
	return errors.New("raw mode not supported on this platform")
}
