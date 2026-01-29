package terminal

import "testing"

func TestHistoryPrevNext(t *testing.T) {
	h := NewInputHistory(100)
	h.Add("first")
	h.Add("second")
	h.Add("third")

	// Walk backwards
	if s, ok := h.Prev(); !ok || s != "third" {
		t.Errorf("Prev() = %q, %v; want third, true", s, ok)
	}
	if s, ok := h.Prev(); !ok || s != "second" {
		t.Errorf("Prev() = %q, %v; want second, true", s, ok)
	}
	if s, ok := h.Prev(); !ok || s != "first" {
		t.Errorf("Prev() = %q, %v; want first, true", s, ok)
	}
	// At beginning, stays on first
	if s, ok := h.Prev(); !ok || s != "first" {
		t.Errorf("Prev() at start = %q, %v; want first, true", s, ok)
	}

	// Walk forward
	if s, ok := h.Next(); !ok || s != "second" {
		t.Errorf("Next() = %q, %v; want second, true", s, ok)
	}
	if s, ok := h.Next(); !ok || s != "third" {
		t.Errorf("Next() = %q, %v; want third, true", s, ok)
	}
	// Past end returns empty (fresh input line)
	if s, ok := h.Next(); !ok || s != "" {
		t.Errorf("Next() past end = %q, %v; want \"\", true", s, ok)
	}
}

func TestHistoryEmpty(t *testing.T) {
	h := NewInputHistory(100)
	if _, ok := h.Prev(); ok {
		t.Error("Prev() on empty should return false")
	}
	if _, ok := h.Next(); ok {
		t.Error("Next() on empty should return false")
	}
}

func TestHistorySkipDuplicates(t *testing.T) {
	h := NewInputHistory(100)
	h.Add("same")
	h.Add("same")
	h.Add("same")

	if s, ok := h.Prev(); !ok || s != "same" {
		t.Errorf("Prev() = %q, %v", s, ok)
	}
	// Should only have one entry
	if s, _ := h.Prev(); s != "same" {
		t.Errorf("should stay on same, got %q", s)
	}
}

func TestHistorySkipEmpty(t *testing.T) {
	h := NewInputHistory(100)
	h.Add("")
	if _, ok := h.Prev(); ok {
		t.Error("empty strings should not be added")
	}
}

func TestHistoryMaxSize(t *testing.T) {
	h := NewInputHistory(3)
	h.Add("a")
	h.Add("b")
	h.Add("c")
	h.Add("d")

	// "a" should have been evicted; walk back 3 times
	h.Reset()
	s1, _ := h.Prev()
	s2, _ := h.Prev()
	s3, _ := h.Prev()
	if s1 != "d" || s2 != "c" || s3 != "b" {
		t.Errorf("entries = [%s, %s, %s], want [d, c, b]", s1, s2, s3)
	}
	// One more Prev stays at the oldest
	s4, _ := h.Prev()
	if s4 != "b" {
		t.Errorf("clamped Prev = %q, want b", s4)
	}
}

func TestHistoryReset(t *testing.T) {
	h := NewInputHistory(100)
	h.Add("one")
	h.Add("two")

	h.Prev()
	h.Prev()

	h.Reset()
	if s, ok := h.Prev(); !ok || s != "two" {
		t.Errorf("after Reset, Prev() = %q, %v; want two, true", s, ok)
	}
}
