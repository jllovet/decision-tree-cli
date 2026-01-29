package cli

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		wantName string
		wantArgs []string
	}{
		{"add decision \"Is it raining?\"", "add", []string{"decision", "Is it raining?"}},
		{"connect n1 n2 yes", "connect", []string{"n1", "n2", "yes"}},
		{"  QUIT  ", "quit", nil},
		{"", "", nil},
		{"edit n1 label 'new label'", "edit", []string{"n1", "label", "new label"}},
		{"render dot", "render", []string{"dot"}},
		{"add action Take umbrella", "add", []string{"action", "Take", "umbrella"}},
	}
	for _, tc := range tests {
		got := Parse(tc.input)
		if got.Name != tc.wantName {
			t.Errorf("Parse(%q).Name = %q, want %q", tc.input, got.Name, tc.wantName)
		}
		if len(got.Args) != len(tc.wantArgs) {
			t.Errorf("Parse(%q).Args = %v (len %d), want %v (len %d)", tc.input, got.Args, len(got.Args), tc.wantArgs, len(tc.wantArgs))
			continue
		}
		for i := range got.Args {
			if got.Args[i] != tc.wantArgs[i] {
				t.Errorf("Parse(%q).Args[%d] = %q, want %q", tc.input, i, got.Args[i], tc.wantArgs[i])
			}
		}
	}
}

func TestTokenize(t *testing.T) {
	tokens := tokenize(`add decision "Is it raining?"`)
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d: %v", len(tokens), tokens)
	}
	if tokens[2] != "Is it raining?" {
		t.Errorf("tokens[2] = %q", tokens[2])
	}
}
