package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/jllovet/decision-tree-cli/internal/model"
	"github.com/jllovet/decision-tree-cli/internal/preview"
	"github.com/jllovet/decision-tree-cli/internal/render"
	"github.com/jllovet/decision-tree-cli/internal/storage"
	"github.com/jllovet/decision-tree-cli/internal/tree"
)

// Session holds the state for a CLI session.
type Session struct {
	Tree      *model.Tree
	History   *tree.History
	Clipboard *tree.Clipboard
	In        io.Reader
	Out       io.Writer
}

// NewSession creates a new CLI session with an empty tree.
func NewSession(w io.Writer) *Session {
	return &Session{
		Tree:    model.NewTree("untitled"),
		History: tree.NewHistory(),
		Out:     w,
	}
}

// Execute dispatches a parsed command to the appropriate handler.
// Returns true if the session should continue, false to quit.
func (s *Session) Execute(cmd ParsedCommand) bool {
	switch cmd.Name {
	case "":
		return true
	case "add":
		s.cmdAdd(cmd.Args)
	case "connect":
		s.cmdConnect(cmd.Args)
	case "disconnect":
		s.cmdDisconnect(cmd.Args)
	case "remove":
		s.cmdRemove(cmd.Args)
	case "edit":
		s.cmdEdit(cmd.Args)
	case "set-root":
		s.cmdSetRoot(cmd.Args)
	case "list":
		s.cmdList()
	case "preview":
		s.cmdPreview()
	case "render":
		s.cmdRender(cmd.Args)
	case "copy":
		s.cmdCopy(cmd.Args)
	case "paste":
		s.cmdPaste()
	case "save":
		s.cmdSave(cmd.Args)
	case "load":
		s.cmdLoad(cmd.Args)
	case "browse":
		s.cmdBrowse()
	case "undo":
		s.cmdUndo()
	case "redo":
		s.cmdRedo()
	case "help":
		s.cmdHelp()
	case "quit", "exit":
		return false
	default:
		fmt.Fprintf(s.Out, "Unknown command: %s (type 'help' for commands)\n", cmd.Name)
	}
	return true
}

func (s *Session) cmdAdd(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(s.Out, "Usage: add <type> <label>")
		fmt.Fprintln(s.Out, "Types: decision, action, startend, io")
		return
	}
	nodeType, err := model.ParseNodeType(args[0])
	if err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	label := strings.Join(args[1:], " ")
	cmd := tree.NewAddNodeCmd(nodeType, label)
	if err := s.History.Execute(s.Tree, cmd); err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	// Get the ID from the command
	type idGetter interface{ ID() string }
	if ig, ok := cmd.(idGetter); ok {
		fmt.Fprintf(s.Out, "Added node %s\n", ig.ID())
	}
}

func (s *Session) cmdConnect(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(s.Out, "Usage: connect <from> <to> [label]")
		return
	}
	label := ""
	if len(args) >= 3 {
		label = strings.Join(args[2:], " ")
	}
	cmd := tree.NewConnectCmd(args[0], args[1], label)
	if err := s.History.Execute(s.Tree, cmd); err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	fmt.Fprintf(s.Out, "Connected %s -> %s\n", args[0], args[1])
}

func (s *Session) cmdDisconnect(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(s.Out, "Usage: disconnect <from> <to>")
		return
	}
	cmd := tree.NewDisconnectCmd(args[0], args[1])
	if err := s.History.Execute(s.Tree, cmd); err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	fmt.Fprintf(s.Out, "Disconnected %s -> %s\n", args[0], args[1])
}

func (s *Session) cmdRemove(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(s.Out, "Usage: remove <node-id>")
		return
	}
	cmd := tree.NewRemoveNodeCmd(args[0])
	if err := s.History.Execute(s.Tree, cmd); err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	fmt.Fprintf(s.Out, "Removed node %s\n", args[0])
}

func (s *Session) cmdEdit(args []string) {
	if len(args) < 3 {
		fmt.Fprintln(s.Out, "Usage: edit <node-id> label <new-label>")
		fmt.Fprintln(s.Out, "       edit <node-id> type <new-type>")
		return
	}
	id := args[0]
	field := strings.ToLower(args[1])
	value := strings.Join(args[2:], " ")

	switch field {
	case "label":
		cmd := tree.NewEditLabelCmd(id, value)
		if err := s.History.Execute(s.Tree, cmd); err != nil {
			fmt.Fprintf(s.Out, "Error: %v\n", err)
			return
		}
		fmt.Fprintf(s.Out, "Updated %s label\n", id)
	case "type":
		nt, err := model.ParseNodeType(value)
		if err != nil {
			fmt.Fprintf(s.Out, "Error: %v\n", err)
			return
		}
		cmd := tree.NewEditTypeCmd(id, nt)
		if err := s.History.Execute(s.Tree, cmd); err != nil {
			fmt.Fprintf(s.Out, "Error: %v\n", err)
			return
		}
		fmt.Fprintf(s.Out, "Updated %s type\n", id)
	default:
		fmt.Fprintf(s.Out, "Unknown field %q (use 'label' or 'type')\n", field)
	}
}

func (s *Session) cmdSetRoot(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(s.Out, "Usage: set-root <node-id>")
		return
	}
	cmd := tree.NewSetRootCmd(args[0])
	if err := s.History.Execute(s.Tree, cmd); err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	fmt.Fprintf(s.Out, "Root set to %s\n", args[0])
}

func (s *Session) cmdList() {
	lines := tree.ListNodes(s.Tree)
	if len(lines) == 0 {
		fmt.Fprintln(s.Out, "(no nodes)")
		return
	}
	for _, line := range lines {
		fmt.Fprintln(s.Out, line)
	}
}

func (s *Session) cmdPreview() {
	fmt.Fprintln(s.Out, preview.Render(s.Tree))
}

func (s *Session) cmdRender(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(s.Out, "Usage: render <dot|mermaid>")
		return
	}
	var r render.Renderer
	switch strings.ToLower(args[0]) {
	case "dot":
		r = &render.DOTRenderer{}
	case "mermaid":
		r = &render.MermaidRenderer{}
	default:
		fmt.Fprintf(s.Out, "Unknown format: %s (use 'dot' or 'mermaid')\n", args[0])
		return
	}
	out, err := r.Render(s.Tree)
	if err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	fmt.Fprint(s.Out, out)
}

func (s *Session) cmdCopy(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(s.Out, "Usage: copy <node-id>")
		return
	}
	cb, err := tree.CopySubtree(s.Tree, args[0])
	if err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	s.Clipboard = cb
	fmt.Fprintf(s.Out, "Copied subtree from %s (%d nodes)\n", args[0], len(cb.Nodes))
}

func (s *Session) cmdPaste() {
	if s.Clipboard == nil {
		fmt.Fprintln(s.Out, "Clipboard is empty")
		return
	}
	cmd := tree.NewPasteSubtreeCmd(s.Clipboard)
	if err := s.History.Execute(s.Tree, cmd); err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	type pastedIDsGetter interface{ PastedIDs() map[string]string }
	if pg, ok := cmd.(pastedIDsGetter); ok {
		idMap := pg.PastedIDs()
		fmt.Fprintf(s.Out, "Pasted %d nodes (root: %s -> %s)\n", len(idMap), s.Clipboard.Root, idMap[s.Clipboard.Root])
	}
}

func (s *Session) cmdSave(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(s.Out, "Usage: save <filename>")
		return
	}
	if err := storage.Save(s.Tree, args[0]); err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	fmt.Fprintf(s.Out, "Saved to %s\n", args[0])
}

func (s *Session) cmdLoad(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(s.Out, "Usage: load <filename>")
		return
	}
	loaded, err := storage.Load(args[0])
	if err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	s.Tree = loaded
	s.History = tree.NewHistory()
	s.Clipboard = nil
	fmt.Fprintf(s.Out, "Loaded %q (%d nodes)\n", loaded.Name, len(loaded.Nodes))
}

func (s *Session) cmdUndo() {
	if err := s.History.Undo(s.Tree); err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	fmt.Fprintln(s.Out, "Undone")
}

func (s *Session) cmdRedo() {
	if err := s.History.Redo(s.Tree); err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
		return
	}
	fmt.Fprintln(s.Out, "Redone")
}

func (s *Session) cmdBrowse() {
	if s.In == nil {
		fmt.Fprintln(s.Out, "Error: browse requires an interactive terminal")
		return
	}
	b := newBrowser(s, s.In, s.Out)
	if err := b.run(); err != nil {
		fmt.Fprintf(s.Out, "Error: %v\n", err)
	}
}

func (s *Session) cmdHelp() {
	help := `Commands:
  add <type> <label>         Add a node (types: decision, action, startend, io)
  connect <from> <to> [label] Connect two nodes with an optional edge label
  disconnect <from> <to>     Remove edge between two nodes
  remove <node-id>           Remove a node and its edges
  edit <id> label <text>     Edit a node's label
  edit <id> type <type>      Edit a node's type
  set-root <node-id>         Set the root node
  list                       List all nodes
  preview                    Show ASCII tree preview
  browse                     Interactive tree browser
  render <dot|mermaid>       Render as DOT or Mermaid diagram
  copy <node-id>             Copy a subtree to clipboard
  paste                      Paste clipboard contents
  save <filename>            Save tree to JSON file
  load <filename>            Load tree from JSON file
  undo                       Undo last action
  redo                       Redo last undone action
  help                       Show this help
  quit                       Exit the program
`
	fmt.Fprint(s.Out, help)
}
