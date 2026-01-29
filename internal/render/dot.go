package render

import (
	"fmt"
	"strings"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

// DOTRenderer renders a tree as a Graphviz DOT diagram.
type DOTRenderer struct{}

func (r *DOTRenderer) Render(t *model.Tree) (string, error) {
	var b strings.Builder
	b.WriteString("digraph ")
	b.WriteString(dotID(t.Name))
	b.WriteString(" {\n")
	b.WriteString("  rankdir=TB;\n\n")

	for _, id := range t.NodeIDs() {
		n := t.Nodes[id]
		shape := dotShape(n.Type)
		b.WriteString(fmt.Sprintf("  %s [label=%s, shape=%s];\n", id, dotLabel(n.Label), shape))
	}

	if len(t.Edges) > 0 {
		b.WriteString("\n")
	}
	for _, e := range t.Edges {
		if e.Label != "" {
			b.WriteString(fmt.Sprintf("  %s -> %s [label=%s];\n", e.FromID, e.ToID, dotLabel(e.Label)))
		} else {
			b.WriteString(fmt.Sprintf("  %s -> %s;\n", e.FromID, e.ToID))
		}
	}

	b.WriteString("}\n")
	return b.String(), nil
}

func dotShape(t model.NodeType) string {
	switch t {
	case model.Decision:
		return "diamond"
	case model.Action:
		return "box"
	case model.StartEnd:
		return "ellipse"
	case model.IO:
		return "parallelogram"
	default:
		return "box"
	}
}

func dotLabel(s string) string {
	escaped := strings.ReplaceAll(s, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	return `"` + escaped + `"`
}

func dotID(s string) string {
	// Simple identifier: replace non-alphanumeric with underscore
	var b strings.Builder
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			b.WriteRune(c)
		} else {
			b.WriteRune('_')
		}
	}
	result := b.String()
	if result == "" {
		return "tree"
	}
	return result
}
