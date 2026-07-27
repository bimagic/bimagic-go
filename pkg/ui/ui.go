package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"bimagic-go/pkg/config"
)

var cmdCache sync.Map

func HasCmd(name string) bool {
	if val, ok := cmdCache.Load(name); ok {
		return val.(bool)
	}
	_, err := exec.LookPath(name)
	found := err == nil
	cmdCache.Store(name, found)
	return found
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func WaitForEnter() {
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func PrintCommand(cmd string) {
	gray := config.GetAnsiEsc(config.Theme["BIMAGIC_MUTED"])
	purple := config.GetAnsiEsc(config.Theme["BIMAGIC_PRIMARY"])
	nc := "\033[0m"
	fmt.Printf("%s %sCommand:%s %s%s\n", gray, purple, nc, cmd, nc)
}

func PrintStatus(msg string) {
	GumStyleWithArgs(fmt.Sprintf("--foreground=%s", config.Theme["BIMAGIC_PRIMARY"]), msg)
	config.PlaySound("success")
}

func PrintError(msg string) {
	GumStyleWithArgs(fmt.Sprintf("--foreground=%s", config.Theme["BIMAGIC_ERROR"]), msg)
	config.PlaySound("error")
}

func PrintWarning(msg string) {
	GumStyleWithArgs(fmt.Sprintf("--foreground=%s", config.Theme["BIMAGIC_WARNING"]), msg)
	config.PlaySound("warning")
}

func DrawProgressBar(label string, percent int) {
	width := 30
	filled := (percent * width) / 100
	empty := width - filled

	bar := strings.Repeat("█", filled)
	bg := strings.Repeat("░", empty)

	c := config.GetAnsiEsc(config.Theme["BIMAGIC_INFO"])
	p := config.GetAnsiEsc(config.Theme["BIMAGIC_PRIMARY"])
	g := config.GetAnsiEsc(config.Theme["BIMAGIC_MUTED"])
	y := config.GetAnsiEsc(config.Theme["BIMAGIC_WARNING"])
	nc := "\033[0m"

	fmt.Printf("\r\033[K%s%-20s%s [%s%s%s%s%s%s] %s%3d%%%s", c, label, nc, p, bar, nc, g, bg, nc, y, percent, nc)
}

func GenerateBar(percentage float64) string {
	width := 20
	intPercentage := int(percentage)
	if intPercentage > 100 {
		intPercentage = 100
	}
	filled := (intPercentage * width) / 100
	empty := width - filled

	colors := []string{
		config.GetAnsiEsc(config.Theme["BANNER_COLOR_1"]),
		config.GetAnsiEsc(config.Theme["BANNER_COLOR_2"]),
		config.GetAnsiEsc(config.Theme["BANNER_COLOR_3"]),
		config.GetAnsiEsc(config.Theme["BANNER_COLOR_4"]),
		config.GetAnsiEsc(config.Theme["BANNER_COLOR_5"]),
	}
	gray := config.GetAnsiEsc(config.Theme["BIMAGIC_MUTED"])
	nc := "\033[0m"

	bar := gray + "[" + nc
	for i := 0; i < filled; i++ {
		colorIdx := (i * 5) / width
		if colorIdx > 4 {
			colorIdx = 4
		}
		bar += colors[colorIdx] + "█"
	}
	bar += gray
	for i := 0; i < empty; i++ {
		bar += "░"
	}
	bar += "]" + nc
	return bar
}

// --- Charm Gum Wrappers ---

func GumChoose(header, cursor, cursorFg string, options ...string) string {
	args := []string{"choose"}
	if header != "" {
		args = append(args, "--header", header, "--header.foreground", config.Theme["BIMAGIC_PRIMARY"])
	} else {
		args = append(args, "--header.foreground", config.Theme["BIMAGIC_PRIMARY"])
	}
	if cursor != "" {
		args = append(args, "--cursor", cursor)
	} else {
		args = append(args, "--cursor", " ")
	}
	if cursorFg != "" {
		args = append(args, "--cursor.foreground", cursorFg)
	} else {
		args = append(args, "--cursor.foreground", config.Theme["BIMAGIC_PRIMARY"])
	}
	args = append(args, "--selected.foreground", config.Theme["BIMAGIC_SECONDARY"])
	args = append(args, options...)

	cmd := exec.Command("gum", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, _ := cmd.Output()
	return strings.TrimSpace(string(out))
}

func GumInput(placeholder, value string) string {
	args := []string{
		"input",
		"--prompt.foreground", config.Theme["BIMAGIC_PRIMARY"],
		"--cursor.foreground", config.Theme["BIMAGIC_PRIMARY"],
		"--header.foreground", config.Theme["BIMAGIC_PRIMARY"],
	}
	if placeholder != "" {
		args = append(args, "--placeholder", placeholder, "--placeholder.foreground", config.Theme["BIMAGIC_MUTED"])
	}
	if value != "" {
		args = append(args, "--value", value)
	}
	cmd := exec.Command("gum", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, _ := cmd.Output()
	return strings.TrimSpace(string(out))
}

func GumConfirm(prompt string) bool {
	args := []string{
		"confirm", prompt,
		"--prompt.foreground", config.Theme["BIMAGIC_PRIMARY"],
		"--selected.background", config.Theme["BIMAGIC_PRIMARY"],
		"--selected.foreground", "0",
		"--unselected.background", config.Theme["BIMAGIC_MUTED"],
		"--unselected.foreground", "255",
	}
	cmd := exec.Command("gum", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run() == nil
}

func GumSpin(title string, cmdArgs ...string) bool {
	args := []string{
		"spin",
		"--spinner.foreground", config.Theme["BIMAGIC_PRIMARY"],
		"--title.foreground", config.Theme["BIMAGIC_INFO"],
		"--title", title,
		"--",
	}
	args = append(args, cmdArgs...)
	cmd := exec.Command("gum", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run() == nil
}

func GumStyleWithArgs(colorArg string, text string) {
	fg := config.Theme["BIMAGIC_PRIMARY"]
	if strings.HasPrefix(colorArg, "--foreground=") {
		fg = strings.TrimPrefix(colorArg, "--foreground=")
	}
	c := config.GetAnsiEsc(fg)
	nc := "\033[0m"
	fmt.Printf("%s%s%s\n", c, text, nc)
}

func GumFilterStdin(items, placeholder string, noLimit bool) string {
	args := []string{
		"filter",
		"--indicator", " ",
		"--indicator.foreground", config.Theme["BIMAGIC_PRIMARY"],
		"--match.foreground", config.Theme["BIMAGIC_SECONDARY"],
		"--prompt.foreground", config.Theme["BIMAGIC_PRIMARY"],
		"--header.foreground", config.Theme["BIMAGIC_PRIMARY"],
	}
	if placeholder != "" {
		args = append(args, "--placeholder", placeholder, "--placeholder.foreground", config.Theme["BIMAGIC_MUTED"])
	}
	if noLimit {
		args = append(args, "--no-limit")
	}
	cmd := exec.Command("gum", args...)
	cmd.Stdin = strings.NewReader(items)
	cmd.Stderr = os.Stderr
	out, _ := cmd.Output()
	return strings.TrimSpace(string(out))
}

func GumWrite(placeholder string) string {
	args := []string{
		"write",
		"--cursor.foreground", config.Theme["BIMAGIC_PRIMARY"],
		"--header.foreground", config.Theme["BIMAGIC_PRIMARY"],
	}
	if placeholder != "" {
		args = append(args, "--placeholder", placeholder, "--placeholder.foreground", config.Theme["BIMAGIC_MUTED"])
	}
	cmd := exec.Command("gum", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, _ := cmd.Output()
	return strings.TrimSpace(string(out))
}
