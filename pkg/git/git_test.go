package git

import (
	"testing"
)

func TestStripAnsi(t *testing.T) {
	colored := "\033[38;5;212mHello World\033[0m"
	clean := stripAnsi(colored)
	if clean != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", clean)
	}
}

func TestDisplayCols(t *testing.T) {
	// Simple ASCII string
	w1 := displayCols("Hello")
	if w1 != 5 {
		t.Errorf("Expected width 5, got %d", w1)
	}

	// Colored string with ANSI escapes
	colored := "\033[38;2;191;194;255mHello\033[0m"
	w2 := displayCols(colored)
	if w2 != 5 {
		t.Errorf("Expected width 5 (ANSI stripped), got %d", w2)
	}

	// String with emoji (2 terminal columns)
	withEmoji := "🪄 Hello"
	w3 := displayCols(withEmoji)
	// '🪄' (2) + ' ' (1) + 'Hello' (5) = 8
	if w3 != 8 {
		t.Errorf("Expected width 8 for emoji string, got %d", w3)
	}
}
