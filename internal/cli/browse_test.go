package cli

import (
	"testing"

	"github.com/jllovet/decision-tree-cli/internal/model"
	"github.com/jllovet/decision-tree-cli/internal/tree"
)

func buildSampleTree() *model.Tree {
	t := model.NewTree("test")
	tree.AddNode(t, model.StartEnd, "Start")   // n1
	tree.AddNode(t, model.Decision, "Auth?")    // n2
	tree.AddNode(t, model.Action, "Grant")      // n3
	tree.AddNode(t, model.IO, "Show login")     // n4
	tree.SetRoot(t, "n1")
	tree.ConnectNodes(t, "n1", "n2", "")
	tree.ConnectNodes(t, "n2", "n3", "yes")
	tree.ConnectNodes(t, "n2", "n4", "no")
	return t
}

func TestFlattenTree(t *testing.T) {
	tr := buildSampleTree()
	rows := flattenTree(tr)

	if len(rows) != 4 {
		t.Fatalf("got %d rows, want 4", len(rows))
	}

	// Check node IDs in DFS order
	wantIDs := []string{"n1", "n2", "n3", "n4"}
	for i, want := range wantIDs {
		if rows[i].nodeID != want {
			t.Errorf("row %d nodeID = %q, want %q", i, rows[i].nodeID, want)
		}
	}

	// Root row should have StartEnd decorator
	if rows[0].text != "([Start])" {
		t.Errorf("row 0 text = %q, want %q", rows[0].text, "([Start])")
	}

	// Second row should have tree connector
	if rows[1].text != "└── <Auth?>" {
		t.Errorf("row 1 text = %q, want %q", rows[1].text, "└── <Auth?>")
	}
}

func TestFlattenTreeEmpty(t *testing.T) {
	tr := model.NewTree("empty")
	rows := flattenTree(tr)
	if len(rows) != 0 {
		t.Errorf("got %d rows for empty tree, want 0", len(rows))
	}
}

func TestFlattenTreeNoRoot(t *testing.T) {
	tr := model.NewTree("noroot")
	tree.AddNode(tr, model.Action, "Orphan")
	rows := flattenTree(tr)
	if len(rows) != 0 {
		t.Errorf("got %d rows for tree with no root, want 0", len(rows))
	}
}

func TestCursorClampAfterDelete(t *testing.T) {
	tr := buildSampleTree()
	b := &browser{
		session: &Session{
			Tree:    tr,
			History: tree.NewHistory(),
		},
	}
	b.refresh()
	// Set cursor to last row
	b.cursor = len(b.rows) - 1

	// Delete the last node
	tree.RemoveNode(tr, b.rows[b.cursor].nodeID)
	b.refresh()

	if b.cursor >= len(b.rows) {
		t.Errorf("cursor %d should be < len(rows) %d after deletion", b.cursor, len(b.rows))
	}
}

func TestViewportOffset(t *testing.T) {
	b := &browser{
		height: 3,
		rows:   make([]flatRow, 10),
		cursor: 0,
		offset: 0,
	}

	// Move cursor to row 4 — should scroll
	b.cursor = 4
	b.scrollToCursor()
	if b.offset != 2 {
		t.Errorf("offset = %d, want 2 (cursor at 4 with height 3)", b.offset)
	}

	// Move cursor back to row 0 — should scroll up
	b.cursor = 0
	b.scrollToCursor()
	if b.offset != 0 {
		t.Errorf("offset = %d, want 0", b.offset)
	}
}
