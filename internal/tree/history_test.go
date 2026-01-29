package tree

import (
	"testing"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

func TestHistoryUndoRedo(t *testing.T) {
	tr := model.NewTree("test")
	h := NewHistory()

	// Add a node
	cmd := NewAddNodeCmd(model.Decision, "q1")
	if err := h.Execute(tr, cmd); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if len(tr.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(tr.Nodes))
	}

	// Undo should remove it
	if err := h.Undo(tr); err != nil {
		t.Fatalf("Undo: %v", err)
	}
	if len(tr.Nodes) != 0 {
		t.Fatalf("expected 0 nodes after undo, got %d", len(tr.Nodes))
	}

	// Redo should restore it
	if err := h.Redo(tr); err != nil {
		t.Fatalf("Redo: %v", err)
	}
	if len(tr.Nodes) != 1 {
		t.Fatalf("expected 1 node after redo, got %d", len(tr.Nodes))
	}
}

func TestHistoryUndoEmpty(t *testing.T) {
	tr := model.NewTree("test")
	h := NewHistory()
	if err := h.Undo(tr); err == nil {
		t.Error("expected error for empty undo")
	}
}

func TestHistoryRedoEmpty(t *testing.T) {
	tr := model.NewTree("test")
	h := NewHistory()
	if err := h.Redo(tr); err == nil {
		t.Error("expected error for empty redo")
	}
}

func TestHistoryRedoClearedOnNewAction(t *testing.T) {
	tr := model.NewTree("test")
	h := NewHistory()

	h.Execute(tr, NewAddNodeCmd(model.Decision, "a"))
	h.Undo(tr)
	if !h.CanRedo() {
		t.Fatal("should be able to redo")
	}

	// New action clears redo
	h.Execute(tr, NewAddNodeCmd(model.Action, "b"))
	if h.CanRedo() {
		t.Error("redo should be cleared after new action")
	}
}

func TestConnectDisconnectCommands(t *testing.T) {
	tr := model.NewTree("test")
	h := NewHistory()

	h.Execute(tr, NewAddNodeCmd(model.Decision, "a"))
	h.Execute(tr, NewAddNodeCmd(model.Action, "b"))
	h.Execute(tr, NewConnectCmd("n1", "n2", "yes"))

	if !tr.HasEdge("n1", "n2") {
		t.Fatal("edge should exist")
	}

	h.Undo(tr)
	if tr.HasEdge("n1", "n2") {
		t.Error("edge should be removed after undo")
	}

	h.Redo(tr)
	if !tr.HasEdge("n1", "n2") {
		t.Error("edge should be restored after redo")
	}
}

func TestDisconnectCommand(t *testing.T) {
	tr := model.NewTree("test")
	h := NewHistory()

	h.Execute(tr, NewAddNodeCmd(model.Decision, "a"))
	h.Execute(tr, NewAddNodeCmd(model.Action, "b"))
	h.Execute(tr, NewConnectCmd("n1", "n2", "labeled"))
	h.Execute(tr, NewDisconnectCmd("n1", "n2"))

	if tr.HasEdge("n1", "n2") {
		t.Error("edge should be removed")
	}

	h.Undo(tr)
	if !tr.HasEdge("n1", "n2") {
		t.Error("edge should be restored")
	}
	// Check label preserved
	for _, e := range tr.Edges {
		if e.FromID == "n1" && e.ToID == "n2" && e.Label != "labeled" {
			t.Errorf("label = %q, want %q", e.Label, "labeled")
		}
	}
}

func TestEditLabelCommand(t *testing.T) {
	tr := model.NewTree("test")
	h := NewHistory()

	h.Execute(tr, NewAddNodeCmd(model.Decision, "old"))
	h.Execute(tr, NewEditLabelCmd("n1", "new"))

	if tr.GetNode("n1").Label != "new" {
		t.Error("label should be updated")
	}

	h.Undo(tr)
	if tr.GetNode("n1").Label != "old" {
		t.Error("label should be restored")
	}
}

func TestEditTypeCommand(t *testing.T) {
	tr := model.NewTree("test")
	h := NewHistory()

	h.Execute(tr, NewAddNodeCmd(model.Decision, "x"))
	h.Execute(tr, NewEditTypeCmd("n1", model.Action))

	if tr.GetNode("n1").Type != model.Action {
		t.Error("type should be updated")
	}

	h.Undo(tr)
	if tr.GetNode("n1").Type != model.Decision {
		t.Error("type should be restored")
	}
}

func TestSetRootCommand(t *testing.T) {
	tr := model.NewTree("test")
	h := NewHistory()

	h.Execute(tr, NewAddNodeCmd(model.Decision, "a"))
	h.Execute(tr, NewAddNodeCmd(model.Decision, "b"))
	h.Execute(tr, NewSetRootCmd("n1"))

	if tr.RootID != "n1" {
		t.Error("root should be n1")
	}

	h.Execute(tr, NewSetRootCmd("n2"))
	if tr.RootID != "n2" {
		t.Error("root should be n2")
	}

	h.Undo(tr)
	if tr.RootID != "n1" {
		t.Error("root should be restored to n1")
	}
}

func TestPasteSubtreeCommand(t *testing.T) {
	tr := model.NewTree("test")
	h := NewHistory()

	h.Execute(tr, NewAddNodeCmd(model.Decision, "root"))
	h.Execute(tr, NewAddNodeCmd(model.Action, "child"))
	h.Execute(tr, NewConnectCmd("n1", "n2", "yes"))

	cb, _ := CopySubtree(tr, "n1")
	pasteCmd := NewPasteSubtreeCmd(cb)
	h.Execute(tr, pasteCmd)

	if len(tr.Nodes) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(tr.Nodes))
	}

	h.Undo(tr)
	if len(tr.Nodes) != 2 {
		t.Fatalf("expected 2 nodes after undo, got %d", len(tr.Nodes))
	}
}

func TestRemoveNodeCommand(t *testing.T) {
	tr := model.NewTree("test")
	h := NewHistory()

	h.Execute(tr, NewAddNodeCmd(model.Decision, "root"))
	h.Execute(tr, NewAddNodeCmd(model.Action, "child"))
	h.Execute(tr, NewConnectCmd("n1", "n2", "yes"))
	h.Execute(tr, NewSetRootCmd("n1"))

	h.Execute(tr, NewRemoveNodeCmd("n1"))
	if tr.GetNode("n1") != nil {
		t.Error("n1 should be removed")
	}
	if tr.RootID != "" {
		t.Error("root should be cleared")
	}

	h.Undo(tr)
	if tr.GetNode("n1") == nil {
		t.Error("n1 should be restored")
	}
	if tr.RootID != "n1" {
		t.Error("root should be restored")
	}
	if !tr.HasEdge("n1", "n2") {
		t.Error("edge should be restored")
	}
}
