package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestREPL(t *testing.T) {
	input := strings.NewReader("add decision \"Is it raining?\"\nlist\nquit\n")
	var out bytes.Buffer
	Run(input, &out)

	output := out.String()
	if !strings.Contains(output, "Decision Tree CLI") {
		t.Error("missing welcome message")
	}
	if !strings.Contains(output, "Added node n1") {
		t.Errorf("missing add output in:\n%s", output)
	}
	if !strings.Contains(output, "n1 [decision]") {
		t.Errorf("missing list output in:\n%s", output)
	}
	if !strings.Contains(output, "Goodbye!") {
		t.Error("missing goodbye message")
	}
}

func TestREPLEOF(t *testing.T) {
	input := strings.NewReader("add decision test\n")
	var out bytes.Buffer
	Run(input, &out) // should not panic on EOF
}
