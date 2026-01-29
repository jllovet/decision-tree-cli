package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/jllovet/decision-tree-cli/internal/terminal"
)

// Run starts the REPL loop with the given reader and writer.
func Run(r io.Reader, w io.Writer) {
	session := NewSession(w)
	fmt.Fprintln(w, "Decision Tree CLI (type 'help' for commands)")

	lr := terminal.NewLineReader(r, w)
	defer lr.Close()
	for {
		line, err := lr.ReadLine("> ")
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		cmd := Parse(line)
		if !session.Execute(cmd) {
			fmt.Fprintln(w, "Goodbye!")
			return
		}
	}
}
