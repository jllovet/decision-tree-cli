package model

import "testing"

func TestNodeTypeString(t *testing.T) {
	tests := []struct {
		nt   NodeType
		want string
	}{
		{Decision, "decision"},
		{Action, "action"},
		{StartEnd, "startend"},
		{IO, "io"},
		{NodeType(99), "unknown"},
	}
	for _, tc := range tests {
		if got := tc.nt.String(); got != tc.want {
			t.Errorf("NodeType(%d).String() = %q, want %q", tc.nt, got, tc.want)
		}
	}
}

func TestParseNodeType(t *testing.T) {
	tests := []struct {
		input   string
		want    NodeType
		wantErr bool
	}{
		{"decision", Decision, false},
		{"action", Action, false},
		{"startend", StartEnd, false},
		{"io", IO, false},
		{"bogus", 0, true},
	}
	for _, tc := range tests {
		got, err := ParseNodeType(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseNodeType(%q) expected error", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseNodeType(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseNodeType(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestNewTree(t *testing.T) {
	tr := NewTree("test")
	if tr.Name != "test" {
		t.Errorf("Name = %q, want %q", tr.Name, "test")
	}
	if tr.Nodes == nil {
		t.Error("Nodes map should be initialized")
	}
	if len(tr.Nodes) != 0 {
		t.Error("Nodes should be empty")
	}
}

func TestNextID(t *testing.T) {
	tr := NewTree("test")
	id1 := tr.NextID()
	id2 := tr.NextID()
	if id1 != "n1" {
		t.Errorf("first ID = %q, want %q", id1, "n1")
	}
	if id2 != "n2" {
		t.Errorf("second ID = %q, want %q", id2, "n2")
	}
}

func TestTreeChildren(t *testing.T) {
	tr := NewTree("test")
	tr.Nodes["n1"] = &Node{ID: "n1", Type: Decision, Label: "root"}
	tr.Nodes["n2"] = &Node{ID: "n2", Type: Action, Label: "a"}
	tr.Nodes["n3"] = &Node{ID: "n3", Type: Action, Label: "b"}
	tr.Edges = []Edge{
		{FromID: "n1", ToID: "n2", Label: "yes"},
		{FromID: "n1", ToID: "n3", Label: "no"},
	}

	children := tr.Children("n1")
	if len(children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(children))
	}
}

func TestTreeParent(t *testing.T) {
	tr := NewTree("test")
	tr.Nodes["n1"] = &Node{ID: "n1", Type: Decision, Label: "root"}
	tr.Nodes["n2"] = &Node{ID: "n2", Type: Action, Label: "child"}
	tr.Edges = []Edge{{FromID: "n1", ToID: "n2", Label: "yes"}}

	p := tr.Parent("n2")
	if p == nil {
		t.Fatal("expected parent")
	}
	if p.FromID != "n1" {
		t.Errorf("parent FromID = %q, want %q", p.FromID, "n1")
	}

	if tr.Parent("n1") != nil {
		t.Error("root should have no parent")
	}
}

func TestHasEdge(t *testing.T) {
	tr := NewTree("test")
	tr.Edges = []Edge{{FromID: "a", ToID: "b"}}
	if !tr.HasEdge("a", "b") {
		t.Error("expected edge a->b")
	}
	if tr.HasEdge("b", "a") {
		t.Error("should not find edge b->a")
	}
}

func TestValidate(t *testing.T) {
	tr := NewTree("test")
	tr.RootID = "missing"
	if err := tr.Validate(); err == nil {
		t.Error("expected error for missing root")
	}

	tr.RootID = ""
	tr.Nodes["n1"] = &Node{ID: "n1"}
	tr.Edges = []Edge{{FromID: "n1", ToID: "missing"}}
	if err := tr.Validate(); err == nil {
		t.Error("expected error for dangling edge")
	}

	tr.Edges = []Edge{{FromID: "missing", ToID: "n1"}}
	if err := tr.Validate(); err == nil {
		t.Error("expected error for dangling edge source")
	}

	tr.Edges = nil
	if err := tr.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAncestors(t *testing.T) {
	tr := NewTree("test")
	tr.Nodes["n1"] = &Node{ID: "n1"}
	tr.Nodes["n2"] = &Node{ID: "n2"}
	tr.Nodes["n3"] = &Node{ID: "n3"}
	tr.Edges = []Edge{
		{FromID: "n1", ToID: "n2"},
		{FromID: "n2", ToID: "n3"},
	}
	anc := tr.Ancestors("n3")
	if !anc["n2"] || !anc["n1"] {
		t.Errorf("expected n1 and n2 as ancestors, got %v", anc)
	}
	if len(anc) != 2 {
		t.Errorf("expected 2 ancestors, got %d", len(anc))
	}
}

func TestNodeIDs(t *testing.T) {
	tr := NewTree("test")
	tr.Nodes["b"] = &Node{ID: "b"}
	tr.Nodes["a"] = &Node{ID: "a"}
	ids := tr.NodeIDs()
	if len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Errorf("NodeIDs() = %v, want [a b]", ids)
	}
}
