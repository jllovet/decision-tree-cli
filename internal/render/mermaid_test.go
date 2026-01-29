package render

import (
	"strings"
	"testing"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

func TestMermaidRenderer(t *testing.T) {
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

	r := &MermaidRenderer{}
	out, err := r.Render(tr)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if !strings.Contains(out, "flowchart TB") {
		t.Error("missing flowchart header")
	}
	if !strings.Contains(out, "n1([Start])") {
		t.Errorf("missing n1 startend node in:\n%s", out)
	}
	if !strings.Contains(out, "n2{Authenticated?}") {
		t.Errorf("missing n2 decision node in:\n%s", out)
	}
	if !strings.Contains(out, "n3[Grant access]") {
		t.Errorf("missing n3 action node in:\n%s", out)
	}
	if !strings.Contains(out, "n4[/Show login/]") {
		t.Errorf("missing n4 IO node in:\n%s", out)
	}
	if !strings.Contains(out, "n1 --> n2") {
		t.Error("missing edge n1->n2")
	}
	if !strings.Contains(out, "n2 -- yes --> n3") {
		t.Error("missing labeled edge n2->n3")
	}
	if !strings.Contains(out, "n2 -- no --> n4") {
		t.Error("missing labeled edge n2->n4")
	}
}

func TestMermaidEmptyTree(t *testing.T) {
	tr := model.NewTree("empty")
	r := &MermaidRenderer{}
	out, err := r.Render(tr)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.HasPrefix(out, "flowchart TB") {
		t.Error("should start with flowchart TB")
	}
}
