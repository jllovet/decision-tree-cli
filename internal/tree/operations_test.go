package tree

import (
	"testing"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

func TestAddNode(t *testing.T) {
	tr := model.NewTree("test")
	id := AddNode(tr, model.Decision, "Is it raining?")
	if id != "n1" {
		t.Errorf("ID = %q, want %q", id, "n1")
	}
	n := tr.GetNode(id)
	if n == nil {
		t.Fatal("node not found")
	}
	if n.Label != "Is it raining?" {
		t.Errorf("Label = %q, want %q", n.Label, "Is it raining?")
	}
	if n.Type != model.Decision {
		t.Errorf("Type = %v, want %v", n.Type, model.Decision)
	}
}

func TestRemoveNode(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "root")
	AddNode(tr, model.Action, "child")
	tr.Edges = append(tr.Edges, model.Edge{FromID: "n1", ToID: "n2"})
	tr.RootID = "n1"

	if err := RemoveNode(tr, "n1"); err != nil {
		t.Fatalf("RemoveNode: %v", err)
	}
	if tr.GetNode("n1") != nil {
		t.Error("node n1 should be removed")
	}
	if len(tr.Edges) != 0 {
		t.Errorf("edges should be empty, got %d", len(tr.Edges))
	}
	if tr.RootID != "" {
		t.Errorf("root should be cleared, got %q", tr.RootID)
	}
}

func TestRemoveNodeNotFound(t *testing.T) {
	tr := model.NewTree("test")
	if err := RemoveNode(tr, "missing"); err == nil {
		t.Error("expected error")
	}
}

func TestConnectNodes(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "root")
	AddNode(tr, model.Action, "child")

	if err := ConnectNodes(tr, "n1", "n2", "yes"); err != nil {
		t.Fatalf("ConnectNodes: %v", err)
	}
	if !tr.HasEdge("n1", "n2") {
		t.Error("edge should exist")
	}
}

func TestConnectNodesMissingSource(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Action, "child")
	if err := ConnectNodes(tr, "missing", "n1", ""); err == nil {
		t.Error("expected error for missing source")
	}
}

func TestConnectNodesMissingTarget(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "root")
	if err := ConnectNodes(tr, "n1", "missing", ""); err == nil {
		t.Error("expected error for missing target")
	}
}

func TestConnectNodesDuplicate(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "a")
	AddNode(tr, model.Action, "b")
	ConnectNodes(tr, "n1", "n2", "")
	if err := ConnectNodes(tr, "n1", "n2", ""); err == nil {
		t.Error("expected error for duplicate edge")
	}
}

func TestConnectNodesSingleParent(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "a")
	AddNode(tr, model.Decision, "b")
	AddNode(tr, model.Action, "c")
	ConnectNodes(tr, "n1", "n3", "")
	if err := ConnectNodes(tr, "n2", "n3", ""); err == nil {
		t.Error("expected error: node already has parent")
	}
}

func TestConnectNodesCycleDetection(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "a")
	AddNode(tr, model.Decision, "b")
	AddNode(tr, model.Decision, "c")
	ConnectNodes(tr, "n1", "n2", "")
	ConnectNodes(tr, "n2", "n3", "")

	// Try to create cycle: n3 -> n1
	if err := ConnectNodes(tr, "n3", "n1", ""); err == nil {
		t.Error("expected cycle detection error")
	}
}

func TestConnectNodesSelfLoop(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "a")
	if err := ConnectNodes(tr, "n1", "n1", ""); err == nil {
		t.Error("expected error for self-loop")
	}
}

func TestDisconnectNodes(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "a")
	AddNode(tr, model.Action, "b")
	ConnectNodes(tr, "n1", "n2", "")

	if err := DisconnectNodes(tr, "n1", "n2"); err != nil {
		t.Fatalf("DisconnectNodes: %v", err)
	}
	if tr.HasEdge("n1", "n2") {
		t.Error("edge should be removed")
	}
}

func TestDisconnectNodesNotFound(t *testing.T) {
	tr := model.NewTree("test")
	if err := DisconnectNodes(tr, "a", "b"); err == nil {
		t.Error("expected error")
	}
}

func TestEditNodeLabel(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "old")
	if err := EditNodeLabel(tr, "n1", "new"); err != nil {
		t.Fatalf("EditNodeLabel: %v", err)
	}
	if tr.GetNode("n1").Label != "new" {
		t.Error("label not updated")
	}
}

func TestEditNodeLabelNotFound(t *testing.T) {
	tr := model.NewTree("test")
	if err := EditNodeLabel(tr, "missing", "x"); err == nil {
		t.Error("expected error")
	}
}

func TestEditNodeType(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "x")
	if err := EditNodeType(tr, "n1", model.Action); err != nil {
		t.Fatalf("EditNodeType: %v", err)
	}
	if tr.GetNode("n1").Type != model.Action {
		t.Error("type not updated")
	}
}

func TestSetRoot(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "root")
	if err := SetRoot(tr, "n1"); err != nil {
		t.Fatalf("SetRoot: %v", err)
	}
	if tr.RootID != "n1" {
		t.Errorf("RootID = %q, want %q", tr.RootID, "n1")
	}
}

func TestSetRootNotFound(t *testing.T) {
	tr := model.NewTree("test")
	if err := SetRoot(tr, "missing"); err == nil {
		t.Error("expected error")
	}
}

func TestListNodes(t *testing.T) {
	tr := model.NewTree("test")
	AddNode(tr, model.Decision, "q1")
	AddNode(tr, model.Action, "a1")
	tr.RootID = "n1"

	lines := ListNodes(tr)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != `n1 [decision] "q1" (root)` {
		t.Errorf("line[0] = %q", lines[0])
	}
	if lines[1] != `n2 [action] "a1"` {
		t.Errorf("line[1] = %q", lines[1])
	}
}
