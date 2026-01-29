package terminal

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// LineReader provides line editing with history support. When the input is a
// TTY it enables raw mode for byte-at-a-time reading with arrow key handling.
// Otherwise it falls back to bufio.Scanner for pipe/test compatibility.
type LineReader struct {
	in      io.Reader
	out     io.Writer
	history *InputHistory
	isTTY   bool
	fd      uintptr
	scanner *bufio.Scanner // non-TTY fallback
}

// NewLineReader creates a LineReader. It auto-detects whether in is a terminal.
func NewLineReader(in io.Reader, out io.Writer) *LineReader {
	lr := &LineReader{
		in:      in,
		out:     out,
		history: NewInputHistory(500),
	}
	if f, ok := in.(*os.File); ok {
		if isTerminal(f) {
			lr.isTTY = true
			lr.fd = f.Fd()
		}
	}
	if !lr.isTTY {
		lr.scanner = bufio.NewScanner(in)
	}
	return lr
}

// wordLeft returns the position at the start of the word to the left of pos.
// It skips spaces first, then skips non-space characters.
func wordLeft(buf []byte, pos int) int {
	for pos > 0 && buf[pos-1] == ' ' {
		pos--
	}
	for pos > 0 && buf[pos-1] != ' ' {
		pos--
	}
	return pos
}

// wordRight returns the position at the end of the word to the right of pos.
// It skips non-space characters first, then skips spaces.
func wordRight(buf []byte, pos int) int {
	n := len(buf)
	for pos < n && buf[pos] != ' ' {
		pos++
	}
	for pos < n && buf[pos] == ' ' {
		pos++
	}
	return pos
}

// isTerminal checks if f is a terminal by attempting to get termios settings.
func isTerminal(f *os.File) bool {
	return isTerminalFd(f.Fd())
}

// ReadLine displays the prompt and reads one line of input. Returns io.EOF
// when the input stream ends or the user presses Ctrl+D.
func (lr *LineReader) ReadLine(prompt string) (string, error) {
	if !lr.isTTY {
		return lr.readLinePipe(prompt)
	}
	return lr.readLineTTY(prompt)
}

// Close is a no-op but satisfies resource-cleanup patterns.
func (lr *LineReader) Close() error {
	return nil
}

func (lr *LineReader) readLinePipe(prompt string) (string, error) {
	fmt.Fprint(lr.out, prompt)
	if lr.scanner.Scan() {
		return lr.scanner.Text(), nil
	}
	if err := lr.scanner.Err(); err != nil {
		return "", err
	}
	return "", io.EOF
}

func (lr *LineReader) readLineTTY(prompt string) (string, error) {
	orig, err := EnableRawMode(lr.fd)
	if err != nil {
		// Fall back to pipe mode if raw mode fails
		if lr.scanner == nil {
			lr.scanner = bufio.NewScanner(lr.in)
		}
		lr.isTTY = false
		return lr.readLinePipe(prompt)
	}
	defer DisableRawMode(lr.fd, &orig)

	lr.history.Reset()

	buf := make([]byte, 0, 256)
	pos := 0 // cursor position within buf

	writePromptAndBuf := func() {
		fmt.Fprintf(lr.out, "\r\x1b[K%s%s", prompt, string(buf))
		// Move cursor to correct position
		if pos < len(buf) {
			fmt.Fprintf(lr.out, "\x1b[%dD", len(buf)-pos)
		}
	}

	fmt.Fprint(lr.out, prompt)

	b := make([]byte, 1)
	for {
		_, err := lr.in.Read(b)
		if err != nil {
			return "", err
		}

		switch {
		case b[0] == '\r' || b[0] == '\n':
			// Submit line
			fmt.Fprint(lr.out, "\r\n")
			line := string(buf)
			lr.history.Add(line)
			return line, nil

		case b[0] == 0x7f || b[0] == 0x08:
			// Backspace
			if pos > 0 {
				buf = append(buf[:pos-1], buf[pos:]...)
				pos--
				writePromptAndBuf()
			}

		case b[0] == 0x03:
			// Ctrl+C: clear line
			buf = buf[:0]
			pos = 0
			writePromptAndBuf()

		case b[0] == 0x01:
			// Ctrl+A: move to start of line
			pos = 0
			writePromptAndBuf()

		case b[0] == 0x05:
			// Ctrl+E: move to end of line
			pos = len(buf)
			writePromptAndBuf()

		case b[0] == 0x0b:
			// Ctrl+K: delete from cursor to end of line
			buf = buf[:pos]
			writePromptAndBuf()

		case b[0] == 0x15:
			// Ctrl+U: delete from cursor to start of line
			buf = buf[pos:]
			pos = 0
			writePromptAndBuf()

		case b[0] == 0x04:
			// Ctrl+D: EOF
			if len(buf) == 0 {
				fmt.Fprint(lr.out, "\r\n")
				return "", io.EOF
			}

		case b[0] == 0x1b:
			// Escape sequence
			seq := make([]byte, 2)
			lr.in.Read(seq)
			if seq[0] == '[' {
				switch seq[1] {
				case 'A': // Up
					if s, ok := lr.history.Prev(); ok {
						buf = []byte(s)
						pos = len(buf)
						writePromptAndBuf()
					}
				case 'B': // Down
					if s, ok := lr.history.Next(); ok {
						buf = []byte(s)
						pos = len(buf)
						writePromptAndBuf()
					}
				case 'C': // Right
					if pos < len(buf) {
						pos++
						fmt.Fprint(lr.out, "\x1b[C")
					}
				case 'D': // Left
					if pos > 0 {
						pos--
						fmt.Fprint(lr.out, "\x1b[D")
					}
				case '1': // Extended sequence: \x1b[1;3C / \x1b[1;3D
					ext := make([]byte, 3) // ;3C or ;3D
					lr.in.Read(ext)
					if ext[0] == ';' && ext[1] == '3' {
						switch ext[2] {
						case 'C': // Opt+Right: word right
							pos = wordRight(buf, pos)
							writePromptAndBuf()
						case 'D': // Opt+Left: word left
							pos = wordLeft(buf, pos)
							writePromptAndBuf()
						}
					}
				}
			} else if seq[0] == 0x7f {
				// Opt+Backspace: delete word before cursor
				newPos := wordLeft(buf, pos)
				buf = append(buf[:newPos], buf[pos:]...)
				pos = newPos
				writePromptAndBuf()
			} else if seq[0] == 'b' {
				// Alt+B: word left
				pos = wordLeft(buf, pos)
				writePromptAndBuf()
			} else if seq[0] == 'f' {
				// Alt+F: word right
				pos = wordRight(buf, pos)
				writePromptAndBuf()
			}

		case b[0] >= 0x20:
			// Printable character: insert at cursor position
			buf = append(buf, 0)
			copy(buf[pos+1:], buf[pos:])
			buf[pos] = b[0]
			pos++
			writePromptAndBuf()
		}
	}
}
