package preview

import (
	"fmt"
	"strings"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

// Render produces an ASCII tree preview using box-drawing characters.
func Render(t *model.Tree) string {
	if t.RootID == "" {
		return "(no root set)"
	}
	if t.GetNode(t.RootID) == nil {
		return "(root node not found)"
	}
	var b strings.Builder
	renderNode(&b, t, t.RootID, "", "", true, true)
	return b.String()
}

func renderNode(b *strings.Builder, t *model.Tree, nodeID, edgeLabel, prefix string, isLast, isRoot bool) {
	n := t.GetNode(nodeID)
	if n == nil {
		return
	}

	// Edge label prefix
	edgePart := ""
	if edgeLabel != "" {
		edgePart = "[" + edgeLabel + "] "
	}

	if isRoot {
		b.WriteString(edgePart + nodeDecorator(n) + "\n")
	} else {
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		b.WriteString(prefix + connector + edgePart + nodeDecorator(n) + "\n")
	}

	// Child prefix
	var childPrefix string
	if isRoot {
		childPrefix = ""
	} else if isLast {
		childPrefix = prefix + "    "
	} else {
		childPrefix = prefix + "│   "
	}

	children := t.Children(nodeID)
	for i, e := range children {
		last := i == len(children)-1
		renderNode(b, t, e.ToID, e.Label, childPrefix, last, false)
	}
}

func nodeDecorator(n *model.Node) string {
	switch n.Type {
	case model.Decision:
		return fmt.Sprintf("<%s>", n.Label)
	case model.Action:
		return fmt.Sprintf("[%s]", n.Label)
	case model.StartEnd:
		return fmt.Sprintf("([%s])", n.Label)
	case model.IO:
		return fmt.Sprintf("//%s//", n.Label)
	default:
		return n.Label
	}
}
