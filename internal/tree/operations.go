package tree

import (
	"fmt"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

// AddNode creates a new node and adds it to the tree. Returns the new node's ID.
func AddNode(t *model.Tree, nodeType model.NodeType, label string) string {
	id := t.NextID()
	t.Nodes[id] = &model.Node{ID: id, Type: nodeType, Label: label}
	return id
}

// RemoveNode removes a node and all its connected edges from the tree.
func RemoveNode(t *model.Tree, id string) error {
	if _, ok := t.Nodes[id]; !ok {
		return fmt.Errorf("node %q not found", id)
	}

	// Remove all edges involving this node
	filtered := t.Edges[:0]
	for _, e := range t.Edges {
		if e.FromID != id && e.ToID != id {
			filtered = append(filtered, e)
		}
	}
	t.Edges = filtered

	delete(t.Nodes, id)

	if t.RootID == id {
		t.RootID = ""
	}
	return nil
}

// ConnectNodes creates a directed edge between two nodes.
func ConnectNodes(t *model.Tree, fromID, toID, label string) error {
	if _, ok := t.Nodes[fromID]; !ok {
		return fmt.Errorf("source node %q not found", fromID)
	}
	if _, ok := t.Nodes[toID]; !ok {
		return fmt.Errorf("target node %q not found", toID)
	}
	if t.HasEdge(fromID, toID) {
		return fmt.Errorf("edge %s -> %s already exists", fromID, toID)
	}
	// Enforce single parent: check if target already has a parent
	if p := t.Parent(toID); p != nil {
		return fmt.Errorf("node %q already has parent %q", toID, p.FromID)
	}
	// Cycle detection: fromID must not be a descendant of toID
	if wouldCreateCycle(t, fromID, toID) {
		return fmt.Errorf("connecting %s -> %s would create a cycle", fromID, toID)
	}

	t.Edges = append(t.Edges, model.Edge{FromID: fromID, ToID: toID, Label: label})
	return nil
}

// wouldCreateCycle checks if adding fromID -> toID would create a cycle.
// This happens if toID is an ancestor of fromID.
func wouldCreateCycle(t *model.Tree, fromID, toID string) bool {
	ancestors := t.Ancestors(fromID)
	return ancestors[toID] || fromID == toID
}

// DisconnectNodes removes the edge between two nodes.
func DisconnectNodes(t *model.Tree, fromID, toID string) error {
	for i, e := range t.Edges {
		if e.FromID == fromID && e.ToID == toID {
			t.Edges = append(t.Edges[:i], t.Edges[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("no edge from %s to %s", fromID, toID)
}

// EditNodeLabel changes the label of a node.
func EditNodeLabel(t *model.Tree, id, label string) error {
	n := t.GetNode(id)
	if n == nil {
		return fmt.Errorf("node %q not found", id)
	}
	n.Label = label
	return nil
}

// EditNodeType changes the type of a node.
func EditNodeType(t *model.Tree, id string, nodeType model.NodeType) error {
	n := t.GetNode(id)
	if n == nil {
		return fmt.Errorf("node %q not found", id)
	}
	n.Type = nodeType
	return nil
}

// SetRoot sets the root node of the tree.
func SetRoot(t *model.Tree, id string) error {
	if _, ok := t.Nodes[id]; !ok {
		return fmt.Errorf("node %q not found", id)
	}
	t.RootID = id
	return nil
}

// ListNodes returns a formatted list of all nodes in the tree.
func ListNodes(t *model.Tree) []string {
	ids := t.NodeIDs()
	lines := make([]string, len(ids))
	for i, id := range ids {
		n := t.Nodes[id]
		root := ""
		if id == t.RootID {
			root = " (root)"
		}
		lines[i] = fmt.Sprintf("%s [%s] %q%s", n.ID, n.Type, n.Label, root)
	}
	return lines
}
