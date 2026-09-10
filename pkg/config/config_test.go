package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAnsiEsc(t *testing.T) {
	// Test ANSI 256 number
	esc := GetAnsiEsc("212")
	if esc != "\033[38;5;212m" {
		t.Errorf("Expected \\033[38;5;212m, got %q", esc)
	}

	// Test Hex TrueColor
	escHex := GetAnsiEsc("#bfc2ff")
	// #bfc2ff -> R: 191, G: 194, B: 255
	expectedHex := "\033[38;2;191;194;255m"
	if escHex != expectedHex {
		t.Errorf("Expected %q, got %q", expectedHex, escHex)
	}
}

func TestLoadTheme(t *testing.T) {
	tmpDir := t.TempDir()
	themeFile := filepath.Join(tmpDir, "theme.wz")

	themeContent := `# Test Theme
BIMAGIC_PRIMARY = "#ff0055"
BIMAGIC_SUCCESS = "46"
# Comment line
CUSTOM_KEY = "123"
`
	if err := os.WriteFile(themeFile, []byte(themeContent), 0o644); err != nil {
		t.Fatalf("Failed to write test theme file: %v", err)
	}

	LoadTheme(themeFile)

	if Theme["BIMAGIC_PRIMARY"] != "#ff0055" {
		t.Errorf("Expected #ff0055 for BIMAGIC_PRIMARY, got %q", Theme["BIMAGIC_PRIMARY"])
	}
	if Theme["BIMAGIC_SUCCESS"] != "46" {
		t.Errorf("Expected 46 for BIMAGIC_SUCCESS, got %q", Theme["BIMAGIC_SUCCESS"])
	}
	if Theme["CUSTOM_KEY"] != "123" {
		t.Errorf("Expected 123 for CUSTOM_KEY, got %q", Theme["CUSTOM_KEY"])
	}
}
