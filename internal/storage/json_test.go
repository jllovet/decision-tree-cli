package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")

	tree := model.NewTree("test-tree")
	tree.Counter = 2
	tree.RootID = "n1"
	tree.Nodes["n1"] = &model.Node{ID: "n1", Type: model.Decision, Label: "Is it raining?"}
	tree.Nodes["n2"] = &model.Node{ID: "n2", Type: model.Action, Label: "Take umbrella"}
	tree.Edges = []model.Edge{{FromID: "n1", ToID: "n2", Label: "yes"}}

	if err := Save(tree, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Name != "test-tree" {
		t.Errorf("Name = %q, want %q", loaded.Name, "test-tree")
	}
	if loaded.RootID != "n1" {
		t.Errorf("RootID = %q, want %q", loaded.RootID, "n1")
	}
	if loaded.Counter != 2 {
		t.Errorf("Counter = %d, want %d", loaded.Counter, 2)
	}
	if len(loaded.Nodes) != 2 {
		t.Errorf("Nodes count = %d, want 2", len(loaded.Nodes))
	}
	if len(loaded.Edges) != 1 {
		t.Errorf("Edges count = %d, want 1", len(loaded.Edges))
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("{invalid"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadInvalidTree(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.json")
	// Tree with edge referencing non-existent node
	data := `{"name":"bad","root_id":"","nodes":{"n1":{"id":"n1","type":0,"label":"x"}},"edges":[{"from":"n1","to":"missing"}],"counter":1}`
	os.WriteFile(path, []byte(data), 0644)

	_, err := Load(path)
	if err == nil {
		t.Error("expected validation error")
	}
}
