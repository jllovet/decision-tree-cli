package terminal

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestLineReaderPipeMode(t *testing.T) {
	in := strings.NewReader("hello\nworld\n")
	var out bytes.Buffer
	lr := NewLineReader(in, &out)

	line, err := lr.ReadLine("> ")
	if err != nil {
		t.Fatal(err)
	}
	if line != "hello" {
		t.Errorf("line = %q, want hello", line)
	}

	line, err = lr.ReadLine("> ")
	if err != nil {
		t.Fatal(err)
	}
	if line != "world" {
		t.Errorf("line = %q, want world", line)
	}
}

func TestLineReaderPipeEOF(t *testing.T) {
	in := strings.NewReader("")
	var out bytes.Buffer
	lr := NewLineReader(in, &out)

	_, err := lr.ReadLine("> ")
	if err != io.EOF {
		t.Errorf("err = %v, want io.EOF", err)
	}
}

func TestLineReaderPipePrompt(t *testing.T) {
	in := strings.NewReader("cmd\n")
	var out bytes.Buffer
	lr := NewLineReader(in, &out)

	lr.ReadLine("dt> ")
	if !strings.Contains(out.String(), "dt> ") {
		t.Errorf("output %q should contain prompt", out.String())
	}
}

func TestLineReaderNonTTY(t *testing.T) {
	// strings.Reader is not *os.File, so should use pipe mode
	in := strings.NewReader("test\n")
	var out bytes.Buffer
	lr := NewLineReader(in, &out)
	if lr.isTTY {
		t.Error("strings.Reader should not be detected as TTY")
	}
}
