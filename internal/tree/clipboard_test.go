package tree

import (
	"testing"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

func TestCopySubtree(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "root")
	AddNode(tr, model.Action, "child1")
	AddNode(tr, model.Action, "child2")
	ConnectNodes(tr, "n1", "n2", "yes")
	ConnectNodes(tr, "n1", "n3", "no")

	cb, err := CopySubtree(tr, "n1")
	if err != nil {
		t.Fatalf("CopySubtree: %v", err)
	}
	if len(cb.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(cb.Nodes))
	}
	if len(cb.Edges) != 2 {
		t.Errorf("expected 2 edges, got %d", len(cb.Edges))
	}
	if cb.Root != "n1" {
		t.Errorf("Root = %q, want %q", cb.Root, "n1")
	}
}

func TestCopySubtreeNotFound(t *testing.T) {
	tr := model.NewTree("test")
	_, err := CopySubtree(tr, "missing")
	if err == nil {
		t.Error("expected error")
	}
}

func TestPasteSubtree(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "root")
	AddNode(tr, model.Action, "child")
	ConnectNodes(tr, "n1", "n2", "yes")

	cb, _ := CopySubtree(tr, "n1")
	idMap := PasteSubtree(tr, cb)

	if len(idMap) != 2 {
		t.Fatalf("expected 2 ID mappings, got %d", len(idMap))
	}

	// Original nodes should still exist
	if tr.GetNode("n1") == nil || tr.GetNode("n2") == nil {
		t.Error("original nodes should still exist")
	}

	// New nodes should exist with remapped IDs
	newRoot := idMap["n1"]
	newChild := idMap["n2"]
	if tr.GetNode(newRoot) == nil {
		t.Errorf("new root %q not found", newRoot)
	}
	if tr.GetNode(newChild) == nil {
		t.Errorf("new child %q not found", newChild)
	}

	// New nodes should have correct labels
	if tr.GetNode(newRoot).Label != "root" {
		t.Error("pasted root label mismatch")
	}

	// New edge should exist
	if !tr.HasEdge(newRoot, newChild) {
		t.Error("pasted edge not found")
	}
}
