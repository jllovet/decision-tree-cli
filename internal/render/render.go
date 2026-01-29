package render

import "github.com/jllovet/decision-tree-cli/internal/model"

// Renderer converts a tree to a string representation.
type Renderer interface {
	Render(t *model.Tree) (string, error)
}
