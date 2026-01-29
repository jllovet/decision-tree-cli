package model

// Edge represents a directed connection between two nodes.
type Edge struct {
	FromID string `json:"from"`
	ToID   string `json:"to"`
	Label  string `json:"label,omitempty"`
}
