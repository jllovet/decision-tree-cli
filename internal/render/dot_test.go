package render

import (
	"strings"
	"testing"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

func TestDOTRenderer(t *testing.T) {
	tr := model.NewTree("auth-flow")
	tr.Nodes["n1"] = &model.Node{ID: "n1", Type: model.StartEnd, Label: "Start"}
	tr.Nodes["n2"] = &model.Node{ID: "n2", Type: model.Decision, Label: "Authenticated?"}
	tr.Nodes["n3"] = &model.Node{ID: "n3", Type: model.Action, Label: "Grant access"}
	tr.Nodes["n4"] = &model.Node{ID: "n4", Type: model.IO, Label: "Show login"}
	tr.Edges = []model.Edge{
		{FromID: "n1", ToID: "n2"},
		{FromID: "n2", ToID: "n3", Label: "yes"},
		{FromID: "n2", ToID: "n4", Label: "no"},
	}

	r := &DOTRenderer{}
	out, err := r.Render(tr)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	// Check structure
	if !strings.Contains(out, "digraph auth_flow") {
		t.Error("missing digraph header")
	}
	if !strings.Contains(out, "rankdir=TB") {
		t.Error("missing rankdir")
	}
	if !strings.Contains(out, `n1 [label="Start", shape=ellipse]`) {
		t.Errorf("missing/wrong n1 definition in:\n%s", out)
	}
	if !strings.Contains(out, `n2 [label="Authenticated?", shape=diamond]`) {
		t.Errorf("missing/wrong n2 definition in:\n%s", out)
	}
	if !strings.Contains(out, `n3 [label="Grant access", shape=box]`) {
		t.Errorf("missing/wrong n3 definition in:\n%s", out)
	}
	if !strings.Contains(out, `n4 [label="Show login", shape=parallelogram]`) {
		t.Errorf("missing/wrong n4 definition in:\n%s", out)
	}
	if !strings.Contains(out, `n2 -> n3 [label="yes"]`) {
		t.Error("missing edge n2->n3")
	}
	if !strings.Contains(out, `n2 -> n4 [label="no"]`) {
		t.Error("missing edge n2->n4")
	}
	if !strings.Contains(out, "n1 -> n2;") {
		t.Error("missing unlabeled edge n1->n2")
	}
}

func TestDOTEscaping(t *testing.T) {
	tr := model.NewTree("test")
	tr.Nodes["n1"] = &model.Node{ID: "n1", Type: model.Action, Label: `Say "hello"`}
	tr.Edges = nil

	r := &DOTRenderer{}
	out, err := r.Render(tr)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, `label="Say \"hello\""`) {
		t.Errorf("escaping issue in:\n%s", out)
	}
}

func TestDotIDEmpty(t *testing.T) {
	if got := dotID(""); got != "tree" {
		t.Errorf("dotID('') = %q, want %q", got, "tree")
	}
}
