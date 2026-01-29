package model

import "fmt"

// NodeType represents the visual shape/type of a node in a decision tree.
type NodeType int

const (
	Decision NodeType = iota // Diamond shape
	Action                   // Rectangle shape
	StartEnd                 // Oval/ellipse shape
	IO                       // Parallelogram shape
)

func (t NodeType) String() string {
	switch t {
	case Decision:
		return "decision"
	case Action:
		return "action"
	case StartEnd:
		return "startend"
	case IO:
		return "io"
	default:
		return "unknown"
	}
}

// ParseNodeType converts a string to a NodeType.
func ParseNodeType(s string) (NodeType, error) {
	switch s {
	case "decision":
		return Decision, nil
	case "action":
		return Action, nil
	case "startend":
		return StartEnd, nil
	case "io":
		return IO, nil
	default:
		return 0, fmt.Errorf("unknown node type: %q", s)
	}
}

// Node represents a single node in a decision tree.
type Node struct {
	ID    string   `json:"id"`
	Type  NodeType `json:"type"`
	Label string   `json:"label"`
}
