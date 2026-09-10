package ui

import (
	"strings"
	"testing"
)

func TestGenerateBar(t *testing.T) {
	bar0 := GenerateBar(0)
	if !strings.Contains(bar0, "[") || !strings.Contains(bar0, "]") {
		t.Errorf("Expected brackets in bar, got %q", bar0)
	}

	bar100 := GenerateBar(100)
	if !strings.Contains(bar100, "█") {
		t.Errorf("Expected filled blocks in 100%% bar, got %q", bar100)
	}
}
