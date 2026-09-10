//go:build !windows

package ui

import (
	"os"
	"syscall"
)

// DrainStdin flushes any unread terminal status responses from the stdin buffer
func DrainStdin() {
	syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), 0x540B, 0)
}
