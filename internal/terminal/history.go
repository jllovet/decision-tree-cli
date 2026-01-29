package terminal

// InputHistory maintains a list of previously entered lines with cursor-based
// navigation. This is distinct from tree.History (Command-pattern undo/redo
// for tree mutations); input history is a flat string list for REPL recall.
type InputHistory struct {
	entries []string
	maxSize int
	cursor  int
}

// NewInputHistory creates a history buffer that retains up to maxSize entries.
func NewInputHistory(maxSize int) *InputHistory {
	return &InputHistory{
		entries: nil,
		maxSize: maxSize,
	}
}

// Add appends a line to history, skipping consecutive duplicates.
func (h *InputHistory) Add(line string) {
	if line == "" {
		return
	}
	if len(h.entries) > 0 && h.entries[len(h.entries)-1] == line {
		return
	}
	h.entries = append(h.entries, line)
	if len(h.entries) > h.maxSize {
		h.entries = h.entries[len(h.entries)-h.maxSize:]
	}
	h.Reset()
}

// Prev moves the cursor back and returns the previous entry (up arrow).
func (h *InputHistory) Prev() (string, bool) {
	if len(h.entries) == 0 {
		return "", false
	}
	if h.cursor > 0 {
		h.cursor--
	}
	return h.entries[h.cursor], true
}

// Next moves the cursor forward and returns the next entry (down arrow).
// Returns ("", true) when moving past the newest entry to indicate the
// user's fresh input line.
func (h *InputHistory) Next() (string, bool) {
	if len(h.entries) == 0 {
		return "", false
	}
	if h.cursor < len(h.entries) {
		h.cursor++
	}
	if h.cursor >= len(h.entries) {
		return "", true
	}
	return h.entries[h.cursor], true
}

// Reset sets the cursor past the end so the next Prev call returns the most
// recent entry. Call this at the start of each new ReadLine.
func (h *InputHistory) Reset() {
	h.cursor = len(h.entries)
}
