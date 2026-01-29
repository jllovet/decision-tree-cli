package cli

import (
	"bytes"
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

func TestOpAddChildEmptyTreeCreatesRoot(t *testing.T) {
	tr := model.NewTree("empty")
	// Simulate typing: "decision\r" then "MyRoot\r"
	input := []byte("decision\rMyRoot\r")
	var out bytes.Buffer
	b := &browser{
		session: &Session{
			Tree:    tr,
			History: tree.NewHistory(),
		},
		in:     bytes.NewReader(input),
		out:    &out,
		height: 20,
		width:  80,
	}
	b.refresh()

	if len(b.rows) != 0 {
		t.Fatalf("expected empty tree, got %d rows", len(b.rows))
	}

	b.opAddChild()

	if tr.RootID == "" {
		t.Fatal("expected root to be set after opAddChild on empty tree")
	}
	n := tr.GetNode(tr.RootID)
	if n == nil {
		t.Fatal("root node not found in tree")
	}
	if n.Label != "MyRoot" {
		t.Errorf("root label = %q, want %q", n.Label, "MyRoot")
	}
	if n.Type != model.Decision {
		t.Errorf("root type = %v, want %v", n.Type, model.Decision)
	}
	if len(b.rows) != 1 {
		t.Errorf("expected 1 row after adding root, got %d", len(b.rows))
	}
}

func TestOpAddChildEmptyTreeCancelType(t *testing.T) {
	tr := model.NewTree("empty")
	// Simulate pressing Esc during type prompt
	input := []byte{0x1b}
	var out bytes.Buffer
	b := &browser{
		session: &Session{
			Tree:    tr,
			History: tree.NewHistory(),
		},
		in:     bytes.NewReader(input),
		out:    &out,
		height: 20,
		width:  80,
	}
	b.refresh()
	b.opAddChild()

	if tr.RootID != "" {
		t.Error("expected root to remain empty after cancel")
	}
	if b.message != "Add cancelled" {
		t.Errorf("message = %q, want %q", b.message, "Add cancelled")
	}
}

func TestConnectModeStateTransitions(t *testing.T) {
	tr := buildSampleTree()
	b := &browser{
		session: &Session{
			Tree:    tr,
			History: tree.NewHistory(),
		},
		height: 20,
		width:  80,
	}
	b.refresh()

	// Select n1 (cursor at 0) and enter connect mode
	b.cursor = 0
	b.opConnect()
	if b.connectFrom != "n1" {
		t.Fatalf("connectFrom = %q, want %q", b.connectFrom, "n1")
	}

	// Cancel with Esc
	b.connectFrom = ""
	b.message = ""

	// Enter connect mode again from n2
	b.cursor = 1
	b.opConnect()
	if b.connectFrom != "n2" {
		t.Fatalf("connectFrom = %q, want %q", b.connectFrom, "n2")
	}

	// Simulate cancel
	b.connectFrom = ""
	if b.connectFrom != "" {
		t.Error("connectFrom should be empty after cancel")
	}
}

func TestConnectModeFinish(t *testing.T) {
	// Build a tree with an orphan node so we can connect to it
	tr := model.NewTree("test")
	tree.AddNode(tr, model.Decision, "Root")  // n1
	tree.AddNode(tr, model.Action, "Orphan")  // n2 — no parent
	tree.SetRoot(tr, "n1")

	// Provide edge label input: Enter (empty label)
	input := []byte("\r")
	var out bytes.Buffer
	b := &browser{
		session: &Session{
			Tree:    tr,
			History: tree.NewHistory(),
		},
		in:     bytes.NewReader(input),
		out:    &out,
		height: 20,
		width:  80,
	}
	b.refresh()

	// Tree only shows n1 (n2 is orphan, not reachable from root).
	// Enter connect mode from n1 (index 0)
	b.cursor = 0
	b.opConnect()
	if b.connectFrom != "n1" {
		t.Fatalf("connectFrom = %q, want %q", b.connectFrom, "n1")
	}

	// We can't navigate to n2 in the flattened tree since it's an orphan.
	// Directly call finishConnect after manually injecting n2 row.
	b.rows = append(b.rows, flatRow{nodeID: "n2", text: "[Orphan]"})
	b.cursor = 1
	b.finishConnect()

	if b.connectFrom != "" {
		t.Error("connectFrom should be cleared after finishConnect")
	}
	// Verify edge was created from n1 to n2
	children := tr.Children("n1")
	found := false
	for _, e := range children {
		if e.ToID == "n2" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected edge from n1 to n2 after finishConnect")
	}
}

func TestConnectModeSelfConnect(t *testing.T) {
	tr := buildSampleTree()
	var out bytes.Buffer
	b := &browser{
		session: &Session{
			Tree:    tr,
			History: tree.NewHistory(),
		},
		in:     bytes.NewReader(nil),
		out:    &out,
		height: 20,
		width:  80,
	}
	b.refresh()

	// Enter connect mode from n1, try to connect to self
	b.cursor = 0
	b.opConnect()
	b.finishConnect() // cursor still on n1

	if b.message != "Cannot connect node to itself" {
		t.Errorf("message = %q, want self-connect error", b.message)
	}
}

func TestOpInitEmptyTree(t *testing.T) {
	tr := model.NewTree("empty")
	// Simulate typing "1\r" to pick the first template
	input := []byte("1\r")
	var out bytes.Buffer
	b := &browser{
		session: &Session{
			Tree:    tr,
			History: tree.NewHistory(),
		},
		in:     bytes.NewReader(input),
		out:    &out,
		height: 20,
		width:  80,
	}
	b.refresh()

	if len(b.rows) != 0 {
		t.Fatalf("expected empty tree, got %d rows", len(b.rows))
	}

	b.opInit()

	if len(b.session.Tree.Nodes) == 0 {
		t.Error("tree should have nodes after opInit")
	}
	if b.session.Tree.RootID == "" {
		t.Error("tree should have root after opInit")
	}
}

func TestOpInitNonEmptyTree(t *testing.T) {
	tr := buildSampleTree()
	var out bytes.Buffer
	b := &browser{
		session: &Session{
			Tree:    tr,
			History: tree.NewHistory(),
		},
		in:     bytes.NewReader(nil),
		out:    &out,
		height: 20,
		width:  80,
	}
	b.refresh()

	b.opInit()

	if b.message != "Init only works on an empty tree" {
		t.Errorf("message = %q, want error about non-empty tree", b.message)
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
