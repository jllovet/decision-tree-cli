# Architecture

## Package Layout

```
internal/
  model/     Data structures (Node, Edge, Tree)
  tree/      Business logic (operations, clipboard, undo/redo)
  render/    Output renderers (DOT, Mermaid)
  preview/   ASCII tree visualization
  storage/   JSON persistence
  terminal/  Terminal raw mode, line editing, input history
  cli/       User interface (parser, commands, REPL)
```

## Data Model

### Node
Each node has an auto-generated ID (`n1`, `n2`, ...), a `NodeType` (Decision, Action, StartEnd, IO), and a `Label`.

### Edge
Directed edges connect nodes by ID, with an optional label (e.g., "yes"/"no" for decision branches).

### Tree
The tree holds a name, a root node ID, a map of nodes, a slice of edges, and an ID counter. It enforces:
- **Single parent**: each node has at most one incoming edge
- **No cycles**: connecting nodes checks the ancestor chain
- **Referential integrity**: edges only reference existing nodes

## Design Decisions

### Command Pattern for Undo/Redo
Every mutating operation is wrapped in a `Command` interface with `Execute` and `Undo` methods. A `History` manager maintains undo/redo stacks. Executing a new command clears the redo stack.

### Clipboard with ID Remapping
Copy performs a DFS deep-copy of a subtree. Paste generates new IDs via `NextID()` and creates a mapping from old to new IDs, preserving structure without collisions.

### Renderer Interface
Both DOT and Mermaid renderers implement `Renderer.Render(*model.Tree) (string, error)`, making it easy to add new output formats.

### Testability
The REPL accepts `io.Reader` and `io.Writer` parameters, allowing full integration testing via piped input/output without needing a real terminal.

### Minimal External Dependencies
The project uses `golang.org/x/sys` for portable terminal raw mode (ioctl access) and otherwise relies only on Go's standard library.
