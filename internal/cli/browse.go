package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jllovet/decision-tree-cli/internal/model"
	"github.com/jllovet/decision-tree-cli/internal/terminal"
	"github.com/jllovet/decision-tree-cli/internal/tree"
)

// flatRow represents one line in the flattened tree view.
type flatRow struct {
	nodeID string
	text   string
}

// flattenTree produces a flat list of rows by DFS-walking the tree,
// mirroring the ASCII preview rendering.
func flattenTree(t *model.Tree) []flatRow {
	if t.RootID == "" {
		return nil
	}
	if t.GetNode(t.RootID) == nil {
		return nil
	}
	var rows []flatRow
	flattenNode(&rows, t, t.RootID, "", "", true, true)
	return rows
}

func flattenNode(rows *[]flatRow, t *model.Tree, nodeID, edgeLabel, prefix string, isLast, isRoot bool) {
	n := t.GetNode(nodeID)
	if n == nil {
		return
	}

	edgePart := ""
	if edgeLabel != "" {
		edgePart = "[" + edgeLabel + "] "
	}

	var text string
	if isRoot {
		text = edgePart + nodeDecorator(n)
	} else {
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		text = prefix + connector + edgePart + nodeDecorator(n)
	}
	*rows = append(*rows, flatRow{nodeID: nodeID, text: text})

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
		flattenNode(rows, t, e.ToID, e.Label, childPrefix, last, false)
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

// browser implements the interactive tree navigator.
type browser struct {
	session *Session
	in      io.Reader
	out     io.Writer
	fd      uintptr

	rows        []flatRow
	cursor      int
	offset      int // first visible row index
	height      int // visible rows (terminal rows minus status/message)
	width       int // terminal columns
	message     string
	connectFrom string // when non-empty, browser is in connect mode
}

func newBrowser(s *Session, in io.Reader, out io.Writer) *browser {
	b := &browser{
		session: s,
		in:      in,
		out:     out,
	}
	if f, ok := in.(*os.File); ok {
		b.fd = f.Fd()
	}
	return b
}

func (b *browser) run() error {
	orig, err := terminal.EnableRawMode(b.fd)
	if err != nil {
		return fmt.Errorf("browse requires a terminal: %w", err)
	}
	defer terminal.DisableRawMode(b.fd, &orig)

	// Switch to alternate screen buffer (like vim/less)
	fmt.Fprint(b.out, "\x1b[?1049h")
	defer fmt.Fprint(b.out, "\x1b[?1049l")

	b.refresh()
	b.render()

	for {
		key := b.readKey()

		// Connect mode: only navigation, Enter, and Esc are active.
		if b.connectFrom != "" {
			switch key {
			case keyUp:
				if b.cursor > 0 {
					b.cursor--
					b.scrollToCursor()
				}
			case keyDown:
				if b.cursor < len(b.rows)-1 {
					b.cursor++
					b.scrollToCursor()
				}
			case keyEnter:
				b.finishConnect()
			case keyEsc, keyQuit:
				b.connectFrom = ""
				b.message = "Connect cancelled"
			default:
				continue
			}
			b.render()
			continue
		}

		switch key {
		case keyQuit, keyEsc:
			return nil
		case keyUp:
			if b.cursor > 0 {
				b.cursor--
				b.scrollToCursor()
			}
		case keyDown:
			if b.cursor < len(b.rows)-1 {
				b.cursor++
				b.scrollToCursor()
			}
		case keyEdit:
			b.opEditLabel()
		case keyCycleType:
			b.opCycleType()
		case keySetRoot:
			b.opSetRoot()
		case keyDelete:
			b.opDelete()
		case keyAddChild:
			b.opAddChild()
		case keyCopy:
			b.opCopy()
		case keyPaste:
			b.opPaste()
		case keyConnect:
			b.opConnect()
		case keyDisconnect:
			b.opDisconnect()
		case keyInit:
			b.opInit()
		case keyUndo:
			b.opUndo()
		case keyRedo:
			b.opRedo()
		default:
			continue
		}
		b.render()
	}
}

func (b *browser) refresh() {
	b.rows = flattenTree(b.session.Tree)
	if b.cursor >= len(b.rows) {
		b.cursor = len(b.rows) - 1
	}
	if b.cursor < 0 {
		b.cursor = 0
	}
	b.updateSize()
	b.scrollToCursor()
}

func (b *browser) updateSize() {
	rows, cols := terminal.TermSize(b.fd)
	b.width = cols
	// Reserve 2 lines: status bar + message line
	b.height = rows - 2
	if b.height < 1 {
		b.height = 1
	}
}

func (b *browser) scrollToCursor() {
	if b.cursor < b.offset {
		b.offset = b.cursor
	}
	if b.cursor >= b.offset+b.height {
		b.offset = b.cursor - b.height + 1
	}
}

func (b *browser) render() {
	b.updateSize()
	// Move cursor to top-left; each row clears to end of line
	fmt.Fprint(b.out, "\x1b[H")

	if len(b.rows) == 0 {
		fmt.Fprint(b.out, "(empty tree — press 'a' to add root, 'i' to init from a template)\x1b[K\r\n")
	} else {
		end := b.offset + b.height
		if end > len(b.rows) {
			end = len(b.rows)
		}
		for i := b.offset; i < end; i++ {
			marker := "  "
			if i == b.cursor {
				marker = "> "
			}
			if b.connectFrom != "" && b.rows[i].nodeID == b.connectFrom {
				marker = "+ "
			}
			if i == b.cursor {
				fmt.Fprintf(b.out, "\x1b[7m%s%s\x1b[0m\x1b[K\r\n", marker, b.rows[i].text)
			} else if b.connectFrom != "" && b.rows[i].nodeID == b.connectFrom {
				fmt.Fprintf(b.out, "\x1b[33m%s%s\x1b[0m\x1b[K\r\n", marker, b.rows[i].text)
			} else {
				fmt.Fprintf(b.out, "%s%s\x1b[K\r\n", marker, b.rows[i].text)
			}
		}
		// Fill remaining lines if tree is shorter than viewport
		for i := end - b.offset; i < b.height; i++ {
			fmt.Fprint(b.out, "~\x1b[K\r\n")
		}
	}

	// Message line
	if b.message != "" {
		fmt.Fprintf(b.out, "\x1b[33m%s\x1b[0m\x1b[K\r\n", b.message)
		b.message = ""
	} else {
		fmt.Fprint(b.out, "\x1b[K\r\n")
	}

	// Status bar (reverse video), padded to full width
	var status string
	if b.connectFrom != "" {
		status = fmt.Sprintf(" Connect %s \u2192 ? | \u2191\u2193 Navigate  Enter Confirm  Esc Cancel", b.connectFrom)
	} else {
		status = " \u2191\u2193/jk Navigate  e Edit  t Type  r Root  d Delete  a Add  y Copy  p Paste  c Connect  D Detach  u Undo  ^R Redo  q Quit"
	}
	if runeLen := len([]rune(status)); runeLen > b.width {
		status = string([]rune(status)[:b.width])
	} else if runeLen < b.width {
		status += strings.Repeat(" ", b.width-runeLen)
	}
	fmt.Fprintf(b.out, "\x1b[7m%s\x1b[0m", status)
}

// Key constants
const (
	keyUp         = iota + 256
	keyDown
	keyQuit
	keyEsc
	keyEnter
	keyEdit
	keyCycleType
	keySetRoot
	keyDelete
	keyAddChild
	keyCopy
	keyPaste
	keyConnect
	keyDisconnect
	keyUndo
	keyRedo
	keyInit
)

func (b *browser) readKey() int {
	buf := make([]byte, 1)
	if _, err := b.in.Read(buf); err != nil {
		return keyQuit
	}
	switch buf[0] {
	case '\r', '\n':
		return keyEnter
	case 'q':
		return keyQuit
	case 0x1b: // Escape sequence
		seq := make([]byte, 2)
		b.in.Read(seq)
		if seq[0] == '[' {
			switch seq[1] {
			case 'A':
				return keyUp
			case 'B':
				return keyDown
			}
		}
		// Bare Esc (seq[0] was not '[' or unrecognized)
		if seq[0] != '[' {
			return keyEsc
		}
		return -1
	case 'k':
		return keyUp
	case 'j':
		return keyDown
	case 'e':
		return keyEdit
	case 't':
		return keyCycleType
	case 'r':
		return keySetRoot
	case 'd':
		return keyDelete
	case 'a':
		return keyAddChild
	case 'y':
		return keyCopy
	case 'p':
		return keyPaste
	case 'c':
		return keyConnect
	case 'D':
		return keyDisconnect
	case 'i':
		return keyInit
	case 'u':
		return keyUndo
	case 0x12: // Ctrl+R
		return keyRedo
	}
	return -1
}

// prompt displays a mini-prompt on the message line and reads text input.
// Returns the entered text and true, or empty string and false if cancelled (Esc).
func (b *browser) prompt(label string) (string, bool) {
	buf := make([]byte, 0, 128)

	redraw := func() {
		// Move to the message line (height + 1 from top)
		fmt.Fprintf(b.out, "\x1b[%d;1H\x1b[K%s%s", b.height+1, label, string(buf))
	}
	redraw()

	raw := make([]byte, 1)
	for {
		if _, err := b.in.Read(raw); err != nil {
			return "", false
		}
		switch {
		case raw[0] == '\r' || raw[0] == '\n':
			return string(buf), true
		case raw[0] == 0x1b:
			return "", false
		case raw[0] == 0x7f || raw[0] == 0x08:
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				redraw()
			}
		case raw[0] == 0x15: // Ctrl+U: clear input
			buf = buf[:0]
			redraw()
		case raw[0] >= 0x20:
			buf = append(buf, raw[0])
			redraw()
		}
	}
}

func (b *browser) selectedNodeID() string {
	if b.cursor < 0 || b.cursor >= len(b.rows) {
		return ""
	}
	return b.rows[b.cursor].nodeID
}

func (b *browser) opEditLabel() {
	id := b.selectedNodeID()
	if id == "" {
		return
	}
	n := b.session.Tree.GetNode(id)
	if n == nil {
		return
	}
	text, ok := b.prompt(fmt.Sprintf("New label for %s [%s]: ", id, n.Label))
	if !ok || text == "" {
		b.message = "Edit cancelled"
		return
	}
	cmd := tree.NewEditLabelCmd(id, text)
	if err := b.session.History.Execute(b.session.Tree, cmd); err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	b.message = fmt.Sprintf("Updated %s label", id)
	b.refresh()
}

func (b *browser) opCycleType() {
	id := b.selectedNodeID()
	if id == "" {
		return
	}
	n := b.session.Tree.GetNode(id)
	if n == nil {
		return
	}
	// Cycle: Decision→Action→StartEnd→IO→Decision
	cycle := []model.NodeType{model.Decision, model.Action, model.StartEnd, model.IO}
	next := cycle[0]
	for i, nt := range cycle {
		if nt == n.Type {
			next = cycle[(i+1)%len(cycle)]
			break
		}
	}
	cmd := tree.NewEditTypeCmd(id, next)
	if err := b.session.History.Execute(b.session.Tree, cmd); err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	b.message = fmt.Sprintf("%s type → %s", id, next)
	b.refresh()
}

func (b *browser) opSetRoot() {
	id := b.selectedNodeID()
	if id == "" {
		return
	}
	cmd := tree.NewSetRootCmd(id)
	if err := b.session.History.Execute(b.session.Tree, cmd); err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	b.message = fmt.Sprintf("Root set to %s", id)
	b.refresh()
}

func (b *browser) opDelete() {
	id := b.selectedNodeID()
	if id == "" {
		return
	}
	cmd := tree.NewRemoveNodeCmd(id)
	if err := b.session.History.Execute(b.session.Tree, cmd); err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	b.message = fmt.Sprintf("Deleted %s", id)
	b.refresh()
}

func (b *browser) opAddChild() {
	// Empty tree: add root node flow.
	if len(b.rows) == 0 {
		b.addRoot()
		return
	}

	parentID := b.selectedNodeID()
	if parentID == "" {
		return
	}
	typeStr, ok := b.prompt("Child type (decision/action/startend/io): ")
	if !ok || typeStr == "" {
		b.message = "Add cancelled"
		return
	}
	nodeType, err := model.ParseNodeType(strings.TrimSpace(typeStr))
	if err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	label, ok := b.prompt("Child label: ")
	if !ok || label == "" {
		b.message = "Add cancelled"
		return
	}
	// Add the node
	addCmd := tree.NewAddNodeCmd(nodeType, label)
	if err := b.session.History.Execute(b.session.Tree, addCmd); err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	// Get the new node's ID
	type idGetter interface{ ID() string }
	childID := ""
	if ig, ok := addCmd.(idGetter); ok {
		childID = ig.ID()
	}
	// Prompt for edge label
	edgeLabel, ok := b.prompt("Edge label (Enter for none): ")
	if !ok {
		edgeLabel = ""
	}
	// Connect parent to child
	connCmd := tree.NewConnectCmd(parentID, childID, edgeLabel)
	if err := b.session.History.Execute(b.session.Tree, connCmd); err != nil {
		b.message = "Error connecting: " + err.Error()
		return
	}
	b.message = fmt.Sprintf("Added %s as child of %s", childID, parentID)
	b.refresh()
}

func (b *browser) opCopy() {
	id := b.selectedNodeID()
	if id == "" {
		return
	}
	cb, err := tree.CopySubtree(b.session.Tree, id)
	if err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	b.session.Clipboard = cb
	b.message = fmt.Sprintf("Copied subtree from %s (%d nodes)", id, len(cb.Nodes))
}

func (b *browser) opPaste() {
	parentID := b.selectedNodeID()
	if parentID == "" {
		return
	}
	if b.session.Clipboard == nil {
		b.message = "Clipboard is empty"
		return
	}
	cmd := tree.NewPasteSubtreeCmd(b.session.Clipboard)
	if err := b.session.History.Execute(b.session.Tree, cmd); err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	// Get pasted root and connect to parent
	type pastedIDsGetter interface{ PastedIDs() map[string]string }
	if pg, ok := cmd.(pastedIDsGetter); ok {
		idMap := pg.PastedIDs()
		newRoot := idMap[b.session.Clipboard.Root]
		connCmd := tree.NewConnectCmd(parentID, newRoot, "")
		if err := b.session.History.Execute(b.session.Tree, connCmd); err != nil {
			b.message = "Pasted but could not connect: " + err.Error()
			b.refresh()
			return
		}
		b.message = fmt.Sprintf("Pasted %d nodes under %s", len(idMap), parentID)
	}
	b.refresh()
}

func (b *browser) addRoot() {
	typeStr, ok := b.prompt("Root type (decision/action/startend/io): ")
	if !ok || typeStr == "" {
		b.message = "Add cancelled"
		return
	}
	nodeType, err := model.ParseNodeType(strings.TrimSpace(typeStr))
	if err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	label, ok := b.prompt("Root label: ")
	if !ok || label == "" {
		b.message = "Add cancelled"
		return
	}
	addCmd := tree.NewAddNodeCmd(nodeType, label)
	if err := b.session.History.Execute(b.session.Tree, addCmd); err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	type idGetter interface{ ID() string }
	newID := ""
	if ig, ok := addCmd.(idGetter); ok {
		newID = ig.ID()
	}
	setRootCmd := tree.NewSetRootCmd(newID)
	if err := b.session.History.Execute(b.session.Tree, setRootCmd); err != nil {
		b.message = "Error setting root: " + err.Error()
		return
	}
	b.message = fmt.Sprintf("Created root node %s", newID)
	b.refresh()
}

func (b *browser) finishConnect() {
	toID := b.selectedNodeID()
	fromID := b.connectFrom
	b.connectFrom = ""
	if toID == "" {
		b.message = "Connect cancelled"
		return
	}
	if fromID == toID {
		b.message = "Cannot connect node to itself"
		return
	}
	edgeLabel, ok := b.prompt("Edge label (Enter for none): ")
	if !ok {
		edgeLabel = ""
	}
	cmd := tree.NewConnectCmd(fromID, toID, edgeLabel)
	if err := b.session.History.Execute(b.session.Tree, cmd); err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	b.message = fmt.Sprintf("Connected %s \u2192 %s", fromID, toID)
	b.refresh()
}

func (b *browser) opConnect() {
	fromID := b.selectedNodeID()
	if fromID == "" {
		return
	}
	b.connectFrom = fromID
}

func (b *browser) opDisconnect() {
	id := b.selectedNodeID()
	if id == "" {
		return
	}
	parentEdge := b.session.Tree.Parent(id)
	if parentEdge == nil {
		b.message = fmt.Sprintf("%s has no parent edge", id)
		return
	}
	cmd := tree.NewDisconnectCmd(parentEdge.FromID, id)
	if err := b.session.History.Execute(b.session.Tree, cmd); err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	b.message = fmt.Sprintf("Disconnected %s from %s", id, parentEdge.FromID)
	b.refresh()
}

func (b *browser) opInit() {
	if len(b.rows) != 0 {
		b.message = "Init only works on an empty tree"
		return
	}
	menu := "Templates:\n"
	for i, tmpl := range templates {
		menu += fmt.Sprintf("  %d. %s — %s\n", i+1, tmpl.Name, tmpl.Description)
	}
	b.message = menu
	b.render()

	choice, ok := b.prompt(fmt.Sprintf("Pick template (1-%d): ", len(templates)))
	if !ok || choice == "" {
		b.message = "Init cancelled"
		return
	}
	var idx int
	if _, err := fmt.Sscanf(choice, "%d", &idx); err != nil || idx < 1 || idx > len(templates) {
		b.message = "Invalid choice"
		return
	}
	tmpl := templates[idx-1]
	b.session.Tree = tmpl.Build()
	b.session.History = tree.NewHistory()
	b.session.Clipboard = nil
	b.message = fmt.Sprintf("Initialized from %q (%d nodes)", tmpl.Name, len(b.session.Tree.Nodes))
	b.refresh()
}

func (b *browser) opUndo() {
	if err := b.session.History.Undo(b.session.Tree); err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	b.message = "Undone"
	b.refresh()
}

func (b *browser) opRedo() {
	if err := b.session.History.Redo(b.session.Tree); err != nil {
		b.message = "Error: " + err.Error()
		return
	}
	b.message = "Redone"
	b.refresh()
}
