package preview

import (
	"strings"
	"testing"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

func TestAsciiPreview(t *testing.T) {
	tr := model.NewTree("test")
	tr.RootID = "n1"
	tr.Nodes["n1"] = &model.Node{ID: "n1", Type: model.Decision, Label: "Is it raining?"}
	tr.Nodes["n2"] = &model.Node{ID: "n2", Type: model.Action, Label: "Take umbrella"}
	tr.Nodes["n3"] = &model.Node{ID: "n3", Type: model.Action, Label: "Go outside"}
	tr.Edges = []model.Edge{
		{FromID: "n1", ToID: "n2", Label: "yes"},
		{FromID: "n1", ToID: "n3", Label: "no"},
	}

	out := Render(tr)

	if !strings.Contains(out, "<Is it raining?>") {
		t.Errorf("missing root node in:\n%s", out)
	}
	if !strings.Contains(out, "[yes]") {
		t.Errorf("missing yes label in:\n%s", out)
	}
	if !strings.Contains(out, "[no]") {
		t.Errorf("missing no label in:\n%s", out)
	}
	if !strings.Contains(out, "[Take umbrella]") {
		t.Errorf("missing action node in:\n%s", out)
	}
	if !strings.Contains(out, "├") || !strings.Contains(out, "└") {
		t.Errorf("missing box-drawing chars in:\n%s", out)
	}
}

func TestAsciiNoRoot(t *testing.T) {
	tr := model.NewTree("test")
	out := Render(tr)
	if out != "(no root set)" {
		t.Errorf("got %q, want %q", out, "(no root set)")
	}
}

func TestAsciiRootNotFound(t *testing.T) {
	tr := model.NewTree("test")
	tr.RootID = "missing"
	out := Render(tr)
	if out != "(root node not found)" {
		t.Errorf("got %q", out)
	}
}

func TestAsciiDeepTree(t *testing.T) {
	tr := model.NewTree("test")
	tr.RootID = "n1"
	tr.Nodes["n1"] = &model.Node{ID: "n1", Type: model.StartEnd, Label: "Start"}
	tr.Nodes["n2"] = &model.Node{ID: "n2", Type: model.Decision, Label: "Check"}
	tr.Nodes["n3"] = &model.Node{ID: "n3", Type: model.IO, Label: "Input"}
	tr.Edges = []model.Edge{
		{FromID: "n1", ToID: "n2"},
		{FromID: "n2", ToID: "n3", Label: "next"},
	}

	out := Render(tr)
	if !strings.Contains(out, "([Start])") {
		t.Errorf("missing startend decorator in:\n%s", out)
	}
	if !strings.Contains(out, "//Input//") {
		t.Errorf("missing IO decorator in:\n%s", out)
	}
}
