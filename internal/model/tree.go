package model

import (
	"fmt"
	"sort"
)

// Tree represents a decision tree with nodes and edges.
type Tree struct {
	Name    string           `json:"name"`
	RootID  string           `json:"root_id"`
	Nodes   map[string]*Node `json:"nodes"`
	Edges   []Edge           `json:"edges"`
	Counter int              `json:"counter"`
}

// NewTree creates a new empty tree with the given name.
func NewTree(name string) *Tree {
	return &Tree{
		Name:  name,
		Nodes: make(map[string]*Node),
	}
}

// NextID generates the next unique node ID.
func (t *Tree) NextID() string {
	t.Counter++
	return fmt.Sprintf("n%d", t.Counter)
}

// GetNode returns the node with the given ID, or nil if not found.
func (t *Tree) GetNode(id string) *Node {
	return t.Nodes[id]
}

// Children returns the IDs of all children of the given node, along with their edge labels.
func (t *Tree) Children(nodeID string) []Edge {
	var children []Edge
	for _, e := range t.Edges {
		if e.FromID == nodeID {
			children = append(children, e)
		}
	}
	return children
}

// Parent returns the parent edge for a node, or nil if none exists.
func (t *Tree) Parent(nodeID string) *Edge {
	for i := range t.Edges {
		if t.Edges[i].ToID == nodeID {
			return &t.Edges[i]
		}
	}
	return nil
}

// HasEdge checks if an edge exists between two nodes.
func (t *Tree) HasEdge(fromID, toID string) bool {
	for _, e := range t.Edges {
		if e.FromID == fromID && e.ToID == toID {
			return true
		}
	}
	return false
}

// Validate checks the tree for structural issues.
func (t *Tree) Validate() error {
	// Check root exists if set
	if t.RootID != "" {
		if _, ok := t.Nodes[t.RootID]; !ok {
			return fmt.Errorf("root node %q not found", t.RootID)
		}
	}

	// Check all edges reference existing nodes
	for _, e := range t.Edges {
		if _, ok := t.Nodes[e.FromID]; !ok {
			return fmt.Errorf("edge references non-existent source node %q", e.FromID)
		}
		if _, ok := t.Nodes[e.ToID]; !ok {
			return fmt.Errorf("edge references non-existent target node %q", e.ToID)
		}
	}

	return nil
}

// NodeIDs returns a sorted list of all node IDs.
func (t *Tree) NodeIDs() []string {
	ids := make([]string, 0, len(t.Nodes))
	for id := range t.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Ancestors returns the set of ancestor node IDs for the given node by walking parent edges.
func (t *Tree) Ancestors(nodeID string) map[string]bool {
	ancestors := make(map[string]bool)
	current := nodeID
	for {
		p := t.Parent(current)
		if p == nil {
			break
		}
		if ancestors[p.FromID] {
			break // cycle detected, stop
		}
		ancestors[p.FromID] = true
		current = p.FromID
	}
	return ancestors
}
