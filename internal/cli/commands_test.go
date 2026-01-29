package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func runCommands(t *testing.T, commands ...string) (*Session, string) {
	t.Helper()
	var buf bytes.Buffer
	s := NewSession(&buf)
	for _, cmd := range commands {
		parsed := Parse(cmd)
		s.Execute(parsed)
	}
	return s, buf.String()
}

func TestCmdAdd(t *testing.T) {
	s, out := runCommands(t, `add decision "Is it raining?"`)
	if !strings.Contains(out, "Added node n1") {
		t.Errorf("output = %q", out)
	}
	if s.Tree.GetNode("n1") == nil {
		t.Error("node n1 not found")
	}
	if s.Tree.GetNode("n1").Label != "Is it raining?" {
		t.Errorf("label = %q", s.Tree.GetNode("n1").Label)
	}
}

func TestCmdAddUsage(t *testing.T) {
	_, out := runCommands(t, "add")
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected usage, got %q", out)
	}
}

func TestCmdAddBadType(t *testing.T) {
	_, out := runCommands(t, "add bogus label")
	if !strings.Contains(out, "Error:") {
		t.Errorf("expected error, got %q", out)
	}
}

func TestCmdConnect(t *testing.T) {
	s, out := runCommands(t,
		`add decision "q1"`,
		`add action "a1"`,
		`connect n1 n2 yes`,
	)
	if !strings.Contains(out, "Connected n1 -> n2") {
		t.Errorf("output = %q", out)
	}
	if !s.Tree.HasEdge("n1", "n2") {
		t.Error("edge not found")
	}
}

func TestCmdConnectUsage(t *testing.T) {
	_, out := runCommands(t, "connect n1")
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected usage, got %q", out)
	}
}

func TestCmdDisconnect(t *testing.T) {
	_, out := runCommands(t,
		`add decision "q1"`,
		`add action "a1"`,
		`connect n1 n2`,
		`disconnect n1 n2`,
	)
	if !strings.Contains(out, "Disconnected n1 -> n2") {
		t.Errorf("output = %q", out)
	}
}

func TestCmdRemove(t *testing.T) {
	s, out := runCommands(t,
		`add decision "q1"`,
		`remove n1`,
	)
	if !strings.Contains(out, "Removed node n1") {
		t.Errorf("output = %q", out)
	}
	if s.Tree.GetNode("n1") != nil {
		t.Error("node should be removed")
	}
}

func TestCmdEdit(t *testing.T) {
	s, out := runCommands(t,
		`add decision "old"`,
		`edit n1 label "new label"`,
	)
	if !strings.Contains(out, "Updated n1 label") {
		t.Errorf("output = %q", out)
	}
	if s.Tree.GetNode("n1").Label != "new label" {
		t.Errorf("label = %q", s.Tree.GetNode("n1").Label)
	}
}

func TestCmdEditType(t *testing.T) {
	s, out := runCommands(t,
		`add decision "x"`,
		`edit n1 type action`,
	)
	if !strings.Contains(out, "Updated n1 type") {
		t.Errorf("output = %q", out)
	}
	if s.Tree.GetNode("n1").Type != 1 { // Action
		t.Errorf("type = %v", s.Tree.GetNode("n1").Type)
	}
}

func TestCmdEditUsage(t *testing.T) {
	_, out := runCommands(t, "edit n1")
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected usage, got %q", out)
	}
}

func TestCmdEditBadField(t *testing.T) {
	_, out := runCommands(t,
		`add decision "x"`,
		`edit n1 color red`,
	)
	if !strings.Contains(out, "Unknown field") {
		t.Errorf("expected unknown field error, got %q", out)
	}
}

func TestCmdSetRoot(t *testing.T) {
	s, out := runCommands(t,
		`add decision "root"`,
		`set-root n1`,
	)
	if !strings.Contains(out, "Root set to n1") {
		t.Errorf("output = %q", out)
	}
	if s.Tree.RootID != "n1" {
		t.Error("root not set")
	}
}

func TestCmdList(t *testing.T) {
	_, out := runCommands(t,
		`add decision "q1"`,
		`add action "a1"`,
		`list`,
	)
	if !strings.Contains(out, "n1 [decision]") {
		t.Errorf("output = %q", out)
	}
	if !strings.Contains(out, "n2 [action]") {
		t.Errorf("output = %q", out)
	}
}

func TestCmdListEmpty(t *testing.T) {
	_, out := runCommands(t, "list")
	if !strings.Contains(out, "(no nodes)") {
		t.Errorf("expected no nodes, got %q", out)
	}
}

func TestCmdPreview(t *testing.T) {
	_, out := runCommands(t,
		`add decision "root"`,
		`add action "child"`,
		`connect n1 n2 yes`,
		`set-root n1`,
		`preview`,
	)
	if !strings.Contains(out, "<root>") {
		t.Errorf("missing root in preview: %q", out)
	}
	if !strings.Contains(out, "[child]") {
		t.Errorf("missing child in preview: %q", out)
	}
}

func TestCmdRenderDot(t *testing.T) {
	_, out := runCommands(t,
		`add decision "q1"`,
		`render dot`,
	)
	if !strings.Contains(out, "digraph") {
		t.Errorf("missing digraph in: %q", out)
	}
}

func TestCmdRenderMermaid(t *testing.T) {
	_, out := runCommands(t,
		`add decision "q1"`,
		`render mermaid`,
	)
	if !strings.Contains(out, "flowchart TB") {
		t.Errorf("missing flowchart in: %q", out)
	}
}

func TestCmdRenderUsage(t *testing.T) {
	_, out := runCommands(t, "render")
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected usage, got %q", out)
	}
}

func TestCmdRenderBadFormat(t *testing.T) {
	_, out := runCommands(t, "render svg")
	if !strings.Contains(out, "Unknown format") {
		t.Errorf("expected unknown format, got %q", out)
	}
}

func TestCmdCopyPaste(t *testing.T) {
	s, out := runCommands(t,
		`add decision "root"`,
		`add action "child"`,
		`connect n1 n2 yes`,
		`copy n1`,
		`paste`,
	)
	if !strings.Contains(out, "Copied subtree from n1 (2 nodes)") {
		t.Errorf("copy output: %q", out)
	}
	if !strings.Contains(out, "Pasted 2 nodes") {
		t.Errorf("paste output: %q", out)
	}
	if len(s.Tree.Nodes) != 4 {
		t.Errorf("expected 4 nodes, got %d", len(s.Tree.Nodes))
	}
}

func TestCmdPasteEmpty(t *testing.T) {
	_, out := runCommands(t, "paste")
	if !strings.Contains(out, "Clipboard is empty") {
		t.Errorf("expected clipboard empty, got %q", out)
	}
}

func TestCmdSaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")

	// Save
	_, out := runCommands(t,
		`add decision "q1"`,
		`add action "a1"`,
		`connect n1 n2 yes`,
		`set-root n1`,
		"save "+path,
	)
	if !strings.Contains(out, "Saved to") {
		t.Errorf("save output: %q", out)
	}

	// Load in new session
	s, out2 := runCommands(t, "load "+path)
	if !strings.Contains(out2, "Loaded") {
		t.Errorf("load output: %q", out2)
	}
	if len(s.Tree.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(s.Tree.Nodes))
	}
}

func TestCmdUndoRedo(t *testing.T) {
	s, out := runCommands(t,
		`add decision "q1"`,
		`undo`,
	)
	if !strings.Contains(out, "Undone") {
		t.Errorf("output = %q", out)
	}
	if len(s.Tree.Nodes) != 0 {
		t.Error("node should be undone")
	}

	// Redo
	var buf bytes.Buffer
	s.Out = &buf
	s.Execute(Parse("redo"))
	if !strings.Contains(buf.String(), "Redone") {
		t.Errorf("redo output = %q", buf.String())
	}
	if len(s.Tree.Nodes) != 1 {
		t.Error("node should be restored")
	}
}

func TestCmdUndoEmpty(t *testing.T) {
	_, out := runCommands(t, "undo")
	if !strings.Contains(out, "Error:") {
		t.Errorf("expected error, got %q", out)
	}
}

func TestCmdHelp(t *testing.T) {
	_, out := runCommands(t, "help")
	if !strings.Contains(out, "Commands:") {
		t.Errorf("missing help text in: %q", out)
	}
}

func TestCmdQuit(t *testing.T) {
	var buf bytes.Buffer
	s := NewSession(&buf)
	cont := s.Execute(Parse("quit"))
	if cont {
		t.Error("quit should return false")
	}
}

func TestCmdExit(t *testing.T) {
	var buf bytes.Buffer
	s := NewSession(&buf)
	cont := s.Execute(Parse("exit"))
	if cont {
		t.Error("exit should return false")
	}
}

func TestCmdUnknown(t *testing.T) {
	_, out := runCommands(t, "bogus")
	if !strings.Contains(out, "Unknown command") {
		t.Errorf("expected unknown command, got %q", out)
	}
}

func TestCmdEmptyLine(t *testing.T) {
	var buf bytes.Buffer
	s := NewSession(&buf)
	cont := s.Execute(Parse(""))
	if !cont {
		t.Error("empty line should continue")
	}
}

func TestCmdInitNoArgs(t *testing.T) {
	_, out := runCommands(t, "init")
	if !strings.Contains(out, "Available templates:") {
		t.Errorf("expected template list, got %q", out)
	}
	if !strings.Contains(out, "auth-flow") {
		t.Errorf("expected auth-flow in list, got %q", out)
	}
}

func TestCmdInitWithName(t *testing.T) {
	s, out := runCommands(t, "init auth-flow")
	if !strings.Contains(out, "Initialized tree from template") {
		t.Errorf("expected confirmation, got %q", out)
	}
	if len(s.Tree.Nodes) == 0 {
		t.Error("tree should have nodes after init")
	}
	if s.Tree.RootID == "" {
		t.Error("tree should have root after init")
	}
}

func TestCmdInitUnknown(t *testing.T) {
	_, out := runCommands(t, "init bogus")
	if !strings.Contains(out, "Unknown template") {
		t.Errorf("expected unknown template error, got %q", out)
	}
	if !strings.Contains(out, "Available templates:") {
		t.Errorf("expected template list in error output, got %q", out)
	}
}
