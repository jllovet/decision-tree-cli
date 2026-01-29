package render

import (
	"fmt"
	"strings"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

// MermaidRenderer renders a tree as a Mermaid flowchart.
type MermaidRenderer struct{}

func (r *MermaidRenderer) Render(t *model.Tree) (string, error) {
	var b strings.Builder
	b.WriteString("flowchart TB\n")

	for _, id := range t.NodeIDs() {
		n := t.Nodes[id]
		b.WriteString(fmt.Sprintf("  %s%s\n", id, mermaidNodeDef(n)))
	}

	if len(t.Edges) > 0 {
		b.WriteString("\n")
	}
	for _, e := range t.Edges {
		if e.Label != "" {
			b.WriteString(fmt.Sprintf("  %s -- %s --> %s\n", e.FromID, mermaidEscape(e.Label), e.ToID))
		} else {
			b.WriteString(fmt.Sprintf("  %s --> %s\n", e.FromID, e.ToID))
		}
	}

	return b.String(), nil
}

func mermaidNodeDef(n *model.Node) string {
	label := mermaidEscape(n.Label)
	switch n.Type {
	case model.Decision:
		return "{" + label + "}"
	case model.Action:
		return "[" + label + "]"
	case model.StartEnd:
		return "([" + label + "])"
	case model.IO:
		return "[/" + label + "/]"
	default:
		return "[" + label + "]"
	}
}

func mermaidEscape(s string) string {
	s = strings.ReplaceAll(s, `"`, "#quot;")
	return s
}
