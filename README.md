# Decision Tree CLI

An interactive terminal tool for building, editing, and visualizing decision trees. Renders to Graphviz DOT, Mermaid diagrams, and ASCII previews. Pure Go standard library, no external dependencies.

## Installation

```bash
# Install as `dt` on your PATH
make install

# Or build locally without installing
make build
```

## Quick Start

```bash
dt
```

### Initialize from a template

```
Decision Tree CLI (type 'help' for commands)
> init
Available templates:
  1. auth-flow — Authentication flow
  2. approval — Approval workflow
  3. troubleshooting — Troubleshooting guide
Usage: init <template-name>
> init auth-flow
Initialized tree from template "auth-flow" (5 nodes)
> preview
([Start])
└── <Authenticated?>
    ├── [yes] [Grant access]
    │   └── ([End])
    └── [no] [Show login form]
```

### Build a tree manually

```
> add startend "Start"
Added node n1
> add decision "Authenticated?"
Added node n2
> add action "Grant access"
Added node n3
> add io "Show login form"
Added node n4
> connect n1 n2
Connected n1 -> n2
> connect n2 n3 yes
Connected n2 -> n3
> connect n2 n4 no
Connected n2 -> n4
> set-root n1
Root set to n1
> preview
([Start])
├── <Authenticated?>
│   ├── [yes] [Grant access]
│   └── [no] //Show login form//
> render dot
digraph untitled {
  rankdir=TB;
  ...
}
> save auth-flow.json
Saved to auth-flow.json
> quit
Goodbye!
```

## Command Reference

| Command | Description |
|---------|-------------|
| `init [name]` | Initialize tree from a template (list templates with no args) |
| `add <type> <label>` | Add a node. Types: `decision`, `action`, `startend`, `io` |
| `connect <from> <to> [label]` | Connect two nodes with an optional edge label |
| `disconnect <from> <to>` | Remove edge between two nodes |
| `remove <node-id>` | Remove a node and its connected edges |
| `edit <id> label <text>` | Change a node's label |
| `edit <id> type <type>` | Change a node's type |
| `set-root <node-id>` | Set the root node for preview/rendering |
| `list` | List all nodes with their types |
| `preview` | ASCII tree preview with box-drawing characters |
| `render dot [file]` | Output Graphviz DOT diagram (optionally to file) |
| `render mermaid [file]` | Output Mermaid flowchart (optionally to file) |
| `copy <node-id>` | Copy a subtree to clipboard |
| `paste` | Paste clipboard contents (IDs are remapped) |
| `save <filename>` | Save tree to JSON file |
| `load <filename>` | Load tree from JSON file |
| `undo` | Undo last action |
| `redo` | Redo last undone action |
| `browse` | Open interactive full-screen tree browser |
| `help` | Show command help |
| `quit` / `exit` | Exit the program |

## Interactive Browser

Launch a full-screen tree browser with `browse`:

```
> browse
```

Key bindings:

| Key | Action |
|-----|--------|
| `j` / `k` | Move cursor down / up |
| `a` | Add child node (or root if tree is empty) |
| `d` | Delete selected node |
| `e` | Edit selected node |
| `c` | Connect mode (select source, move to target, confirm) |
| `i` | Init from template (empty tree only) |
| `u` / `r` | Undo / Redo |
| `q` | Quit browser |

## Templates

Available templates for `init`:

| Name | Description |
|------|-------------|
| `auth-flow` | Authentication flow: Start → Authenticated? → Grant access / Show login form |
| `approval` | Approval workflow: Start → Submit → Approved? → Process / Revise (loop) |
| `troubleshooting` | Troubleshooting guide: Start → Plugged in? → Plug it in / Check settings → Resolved? |

Use from the REPL with `init <name>` or from the browser with `i` on an empty tree.

## Node Types and Shapes

| Type | DOT Shape | Mermaid Syntax | ASCII Preview |
|------|-----------|----------------|---------------|
| `decision` | diamond | `{label}` | `<label>` |
| `action` | box | `[label]` | `[label]` |
| `startend` | ellipse | `([label])` | `([label])` |
| `io` | parallelogram | `[/label/]` | `//label//` |

## JSON File Format

Trees are saved as JSON with the following structure:

```json
{
  "name": "auth-flow",
  "root_id": "n1",
  "nodes": {
    "n1": { "id": "n1", "type": 0, "label": "Authenticated?" }
  },
  "edges": [
    { "from": "n1", "to": "n2", "label": "yes" }
  ],
  "counter": 2
}
```

Node type values: `0` = decision, `1` = action, `2` = startend, `3` = io.

## Project Structure

```
cmd/decision-tree-cli/   Main entrypoint
internal/
  model/                 Node, Edge, Tree data structures
  tree/                  Operations, clipboard, undo/redo history
  render/                DOT and Mermaid renderers
  preview/               ASCII tree preview
  storage/               JSON save/load
  cli/                   Parser, commands, REPL loop, templates, browser
  terminal/              Raw-mode terminal I/O and line reader
testdata/                Sample fixtures and golden files
docs/                    Architecture and examples
.github/workflows/       CI: build, test, security scanning
```

## Development

```bash
make install  # Install as `dt` to GOPATH/bin
make build    # Build binary to bin/
make test     # Run all tests
make vet      # Static analysis
make run      # Build and run
make clean    # Remove build artifacts
```

### CI

A GitHub Actions workflow runs on push and PR to `main`:

- **Build & Vet** — `go build`, `go vet`, `go test`
- **govulncheck** — Scans dependencies for known vulnerabilities
- **gosec** — Security-focused static analysis
- **staticcheck** — Advanced Go linter
- **Dependency Review** — Reviews new dependencies on PRs
- **Module Integrity** — Verifies checksums and tidy modules
