package render

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jllovet/decision-tree-cli/internal/storage"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testdata")
}

func TestDOTGolden(t *testing.T) {
	treePath := filepath.Join(testdataDir(), "sample-tree.json")
	goldenPath := filepath.Join(testdataDir(), "expected-dot.txt")

	tr, err := storage.Load(treePath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	r := &DOTRenderer{}
	got, err := r.Render(tr)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if got != string(expected) {
		t.Errorf("DOT output mismatch.\nGot:\n%s\nExpected:\n%s", got, string(expected))
	}
}

func TestMermaidGolden(t *testing.T) {
	treePath := filepath.Join(testdataDir(), "sample-tree.json")
	goldenPath := filepath.Join(testdataDir(), "expected-mermaid.txt")

	tr, err := storage.Load(treePath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	r := &MermaidRenderer{}
	got, err := r.Render(tr)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if got != string(expected) {
		t.Errorf("Mermaid output mismatch.\nGot:\n%s\nExpected:\n%s", got, string(expected))
	}
}
