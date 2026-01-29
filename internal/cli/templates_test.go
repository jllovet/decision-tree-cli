package cli

import (
	"testing"
)

func TestTemplatesBuildsValidTrees(t *testing.T) {
	for _, tmpl := range templates {
		t.Run(tmpl.Name, func(t *testing.T) {
			tr := tmpl.Build()
			if len(tr.Nodes) == 0 {
				t.Error("template built empty tree")
			}
			if tr.RootID == "" {
				t.Error("template tree has no root")
			}
			if tr.GetNode(tr.RootID) == nil {
				t.Error("root node not found in tree")
			}
			if len(tr.Edges) == 0 {
				t.Error("template tree has no edges")
			}
			if err := tr.Validate(); err != nil {
				t.Errorf("template tree invalid: %v", err)
			}
		})
	}
}

func TestFindTemplate(t *testing.T) {
	tmpl := findTemplate("auth-flow")
	if tmpl == nil {
		t.Fatal("expected to find auth-flow template")
	}
	if tmpl.Name != "auth-flow" {
		t.Errorf("name = %q, want %q", tmpl.Name, "auth-flow")
	}
}

func TestFindTemplateUnknown(t *testing.T) {
	tmpl := findTemplate("nonexistent")
	if tmpl != nil {
		t.Error("expected nil for unknown template")
	}
}
