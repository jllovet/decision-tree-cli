package cli

import (
	"strings"
)

// ParsedCommand represents a parsed CLI command.
type ParsedCommand struct {
	Name string
	Args []string
}

// Parse tokenizes the input line and returns a ParsedCommand.
func Parse(line string) ParsedCommand {
	tokens := tokenize(line)
	if len(tokens) == 0 {
		return ParsedCommand{}
	}
	return ParsedCommand{
		Name: strings.ToLower(tokens[0]),
		Args: tokens[1:],
	}
}

// tokenize splits a line into tokens, respecting quoted strings.
func tokenize(line string) []string {
	var tokens []string
	var current strings.Builder
	inQuote := false
	quoteChar := byte(0)

	for i := 0; i < len(line); i++ {
		c := line[i]
		switch {
		case inQuote:
			if c == quoteChar {
				inQuote = false
			} else {
				current.WriteByte(c)
			}
		case c == '"' || c == '\'':
			inQuote = true
			quoteChar = c
		case c == ' ' || c == '\t':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(c)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}
