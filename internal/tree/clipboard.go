package tree

import "github.com/jllovet/decision-tree-cli/internal/model"

// Clipboard holds a copied subtree for paste operations.
type Clipboard struct {
	Nodes []model.Node
	Edges []model.Edge
	Root  string // root of the copied subtree
}

// CopySubtree performs a DFS deep-copy of a subtree rooted at nodeID.
func CopySubtree(t *model.Tree, nodeID string) (*Clipboard, error) {
	node := t.GetNode(nodeID)
	if node == nil {
		return nil, errNodeNotFound(nodeID)
	}

	cb := &Clipboard{Root: nodeID}
	visited := make(map[string]bool)
	var dfs func(id string)
	dfs = func(id string) {
		if visited[id] {
			return
		}
		visited[id] = true
		n := t.GetNode(id)
		if n == nil {
			return
		}
		cb.Nodes = append(cb.Nodes, *n)
		for _, e := range t.Children(id) {
			cb.Edges = append(cb.Edges, e)
			dfs(e.ToID)
		}
	}
	dfs(nodeID)
	return cb, nil
}

// PasteSubtree pastes the clipboard contents into the tree, remapping IDs.
// Returns a map from old IDs to new IDs.
func PasteSubtree(t *model.Tree, cb *Clipboard) map[string]string {
	idMap := make(map[string]string)

	// Create new nodes with remapped IDs
	for _, n := range cb.Nodes {
		newID := t.NextID()
		idMap[n.ID] = newID
		t.Nodes[newID] = &model.Node{
			ID:    newID,
			Type:  n.Type,
			Label: n.Label,
		}
	}

	// Create new edges with remapped IDs
	for _, e := range cb.Edges {
		t.Edges = append(t.Edges, model.Edge{
			FromID: idMap[e.FromID],
			ToID:   idMap[e.ToID],
			Label:  e.Label,
		})
	}

	return idMap
}

func errNodeNotFound(id string) error {
	return &nodeNotFoundError{id}
}

type nodeNotFoundError struct {
	id string
}

func (e *nodeNotFoundError) Error() string {
	return "node " + e.id + " not found"
}
