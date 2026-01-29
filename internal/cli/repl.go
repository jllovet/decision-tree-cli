package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Run starts the REPL loop with the given reader and writer.
func Run(r io.Reader, w io.Writer) {
	session := NewSession(w)
	fmt.Fprintln(w, "Decision Tree CLI (type 'help' for commands)")

	scanner := bufio.NewScanner(r)
	for {
		fmt.Fprint(w, "> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		cmd := Parse(line)
		if !session.Execute(cmd) {
			fmt.Fprintln(w, "Goodbye!")
			return
		}
	}
}
