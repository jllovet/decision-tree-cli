package tree

import "github.com/jllovet/decision-tree-cli/internal/model"

// Command represents an undoable operation.
type Command interface {
	Execute(t *model.Tree) error
	Undo(t *model.Tree) error
}

// History manages undo/redo stacks.
type History struct {
	undoStack []Command
	redoStack []Command
}

// NewHistory creates a new History manager.
func NewHistory() *History {
	return &History{}
}

// Execute runs a command and pushes it onto the undo stack.
func (h *History) Execute(t *model.Tree, cmd Command) error {
	if err := cmd.Execute(t); err != nil {
		return err
	}
	h.undoStack = append(h.undoStack, cmd)
	h.redoStack = nil // clear redo on new action
	return nil
}

// Undo undoes the last command.
func (h *History) Undo(t *model.Tree) error {
	if len(h.undoStack) == 0 {
		return errNothingToUndo
	}
	cmd := h.undoStack[len(h.undoStack)-1]
	h.undoStack = h.undoStack[:len(h.undoStack)-1]
	if err := cmd.Undo(t); err != nil {
		return err
	}
	h.redoStack = append(h.redoStack, cmd)
	return nil
}

// Redo re-applies the last undone command.
func (h *History) Redo(t *model.Tree) error {
	if len(h.redoStack) == 0 {
		return errNothingToRedo
	}
	cmd := h.redoStack[len(h.redoStack)-1]
	h.redoStack = h.redoStack[:len(h.redoStack)-1]
	if err := cmd.Execute(t); err != nil {
		return err
	}
	h.undoStack = append(h.undoStack, cmd)
	return nil
}

// CanUndo returns true if there are actions to undo.
func (h *History) CanUndo() bool {
	return len(h.undoStack) > 0
}

// CanRedo returns true if there are actions to redo.
func (h *History) CanRedo() bool {
	return len(h.redoStack) > 0
}

// --- Concrete Commands ---

type addNodeCmd struct {
	nodeType model.NodeType
	label    string
	id       string // set after execute
}

func NewAddNodeCmd(nodeType model.NodeType, label string) Command {
	return &addNodeCmd{nodeType: nodeType, label: label}
}

func (c *addNodeCmd) Execute(t *model.Tree) error {
	c.id = AddNode(t, c.nodeType, c.label)
	return nil
}

func (c *addNodeCmd) Undo(t *model.Tree) error {
	return RemoveNode(t, c.id)
}

// ID returns the created node's ID (available after Execute).
func (c *addNodeCmd) ID() string {
	return c.id
}

type removeNodeCmd struct {
	id            string
	removedNode   model.Node
	removedEdges  []model.Edge
	wasRoot       bool
}

func NewRemoveNodeCmd(id string) Command {
	return &removeNodeCmd{id: id}
}

func (c *removeNodeCmd) Execute(t *model.Tree) error {
	n := t.GetNode(c.id)
	if n == nil {
		return errNodeNotFound(c.id)
	}
	c.removedNode = *n
	c.wasRoot = t.RootID == c.id

	// Save edges that will be removed
	c.removedEdges = nil
	for _, e := range t.Edges {
		if e.FromID == c.id || e.ToID == c.id {
			c.removedEdges = append(c.removedEdges, e)
		}
	}
	return RemoveNode(t, c.id)
}

func (c *removeNodeCmd) Undo(t *model.Tree) error {
	t.Nodes[c.id] = &model.Node{
		ID:    c.removedNode.ID,
		Type:  c.removedNode.Type,
		Label: c.removedNode.Label,
	}
	t.Edges = append(t.Edges, c.removedEdges...)
	if c.wasRoot {
		t.RootID = c.id
	}
	return nil
}

type connectCmd struct {
	fromID, toID, label string
}

func NewConnectCmd(fromID, toID, label string) Command {
	return &connectCmd{fromID: fromID, toID: toID, label: label}
}

func (c *connectCmd) Execute(t *model.Tree) error {
	return ConnectNodes(t, c.fromID, c.toID, c.label)
}

func (c *connectCmd) Undo(t *model.Tree) error {
	return DisconnectNodes(t, c.fromID, c.toID)
}

type disconnectCmd struct {
	fromID, toID string
	label        string // saved for undo
}

func NewDisconnectCmd(fromID, toID string) Command {
	return &disconnectCmd{fromID: fromID, toID: toID}
}

func (c *disconnectCmd) Execute(t *model.Tree) error {
	// Save the label before removing
	for _, e := range t.Edges {
		if e.FromID == c.fromID && e.ToID == c.toID {
			c.label = e.Label
			break
		}
	}
	return DisconnectNodes(t, c.fromID, c.toID)
}

func (c *disconnectCmd) Undo(t *model.Tree) error {
	return ConnectNodes(t, c.fromID, c.toID, c.label)
}

type editLabelCmd struct {
	id       string
	newLabel string
	oldLabel string
}

func NewEditLabelCmd(id, newLabel string) Command {
	return &editLabelCmd{id: id, newLabel: newLabel}
}

func (c *editLabelCmd) Execute(t *model.Tree) error {
	n := t.GetNode(c.id)
	if n == nil {
		return errNodeNotFound(c.id)
	}
	c.oldLabel = n.Label
	return EditNodeLabel(t, c.id, c.newLabel)
}

func (c *editLabelCmd) Undo(t *model.Tree) error {
	return EditNodeLabel(t, c.id, c.oldLabel)
}

type editTypeCmd struct {
	id      string
	newType model.NodeType
	oldType model.NodeType
}

func NewEditTypeCmd(id string, newType model.NodeType) Command {
	return &editTypeCmd{id: id, newType: newType}
}

func (c *editTypeCmd) Execute(t *model.Tree) error {
	n := t.GetNode(c.id)
	if n == nil {
		return errNodeNotFound(c.id)
	}
	c.oldType = n.Type
	return EditNodeType(t, c.id, c.newType)
}

func (c *editTypeCmd) Undo(t *model.Tree) error {
	return EditNodeType(t, c.id, c.oldType)
}

type setRootCmd struct {
	newRoot string
	oldRoot string
}

func NewSetRootCmd(id string) Command {
	return &setRootCmd{newRoot: id}
}

func (c *setRootCmd) Execute(t *model.Tree) error {
	c.oldRoot = t.RootID
	return SetRoot(t, c.newRoot)
}

func (c *setRootCmd) Undo(t *model.Tree) error {
	t.RootID = c.oldRoot
	return nil
}

type pasteSubtreeCmd struct {
	clipboard *Clipboard
	idMap     map[string]string
}

func NewPasteSubtreeCmd(cb *Clipboard) Command {
	return &pasteSubtreeCmd{clipboard: cb}
}

func (c *pasteSubtreeCmd) Execute(t *model.Tree) error {
	c.idMap = PasteSubtree(t, c.clipboard)
	return nil
}

func (c *pasteSubtreeCmd) Undo(t *model.Tree) error {
	// Remove all pasted nodes (which also cleans up edges via RemoveNode)
	for _, newID := range c.idMap {
		RemoveNode(t, newID)
	}
	return nil
}

// PastedIDs returns the ID mapping from the paste operation (available after Execute).
func (c *pasteSubtreeCmd) PastedIDs() map[string]string {
	return c.idMap
}

// sentinel errors
type sentinelError string

func (e sentinelError) Error() string { return string(e) }

const (
	errNothingToUndo sentinelError = "nothing to undo"
	errNothingToRedo sentinelError = "nothing to redo"
)
