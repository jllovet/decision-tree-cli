package main

import (
	"os"

	"github.com/jllovet/decision-tree-cli/internal/cli"
)

func main() {
	cli.Run(os.Stdin, os.Stdout)
}
