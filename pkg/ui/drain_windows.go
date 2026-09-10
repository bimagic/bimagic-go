//go:build windows

package ui

// DrainStdin is a no-op on Windows
func DrainStdin() {
}
