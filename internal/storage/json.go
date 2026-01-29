package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jllovet/decision-tree-cli/internal/model"
)

// Save writes the tree to a JSON file.
func Save(tree *model.Tree, path string) error {
	data, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// Load reads a tree from a JSON file and validates it.
func Load(path string) (*model.Tree, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	var tree model.Tree
	if err := json.Unmarshal(data, &tree); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if tree.Nodes == nil {
		tree.Nodes = make(map[string]*model.Node)
	}
	if err := tree.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	return &tree, nil
}
