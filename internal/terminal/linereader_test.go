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

func TestWordLeft(t *testing.T) {
	tests := []struct {
		buf  string
		pos  int
		want int
	}{
		{"hello world", 11, 6},
		{"hello world", 6, 0},
		{"hello world", 5, 0},
		{"hello world", 0, 0},
		{"  hello", 7, 2},
		{"hello  world", 12, 7},
		{"hello  world", 7, 0},
		{"abc", 2, 0},
	}
	for _, tt := range tests {
		got := wordLeft([]byte(tt.buf), tt.pos)
		if got != tt.want {
			t.Errorf("wordLeft(%q, %d) = %d, want %d", tt.buf, tt.pos, got, tt.want)
		}
	}
}

func TestWordRight(t *testing.T) {
	tests := []struct {
		buf  string
		pos  int
		want int
	}{
		{"hello world", 0, 6},
		{"hello world", 6, 11},
		{"hello world", 11, 11},
		{"hello  world", 0, 7},
		{"  hello", 0, 2},
		{"hello  ", 0, 7},
		{"abc", 1, 3},
	}
	for _, tt := range tests {
		got := wordRight([]byte(tt.buf), tt.pos)
		if got != tt.want {
			t.Errorf("wordRight(%q, %d) = %d, want %d", tt.buf, tt.pos, got, tt.want)
		}
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
