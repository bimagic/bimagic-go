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

func TestDrainStdin(t *testing.T) {
	// Ensure DrainStdin runs without panicking
	DrainStdin()
}

func TestGumSpin(t *testing.T) {
	// Test success command
	if !GumSpin("Testing success", "echo", "ok") {
		t.Errorf("Expected GumSpin with echo to return true")
	}

	// Test failing command
	if GumSpin("Testing failure", "sh", "-c", "exit 1") {
		t.Errorf("Expected GumSpin with exit 1 to return false")
	}
}
