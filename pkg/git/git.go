package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"bimagic-go/pkg/config"
	"bimagic-go/pkg/ui"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripAnsi(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func RunGitCmd(args ...string) error {
	cmd := exec.Command("git", args...)
	return cmd.Run()
}

func GetGitOutput(args ...string) string {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func RunCmdOutToScreen(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func IsGitRepo() bool {
	return RunGitCmd("rev-parse", "--git-dir") == nil
}

func GetCurrentBranch() string {
	b := GetGitOutput("branch", "--show-current")
	if b == "" {
		return "main"
	}
	return b
}

func SetupRemote(remoteName string) bool {
	if !IsGitRepo() {
		ui.PrintError("Not a git repository! Initialize it first.")
		return false
	}
	if remoteName == "" {
		remoteName = "origin"
	}

	protocol := ui.GumChoose(fmt.Sprintf("Select protocol for '%s':", remoteName), "", "", "HTTPS (Token)", "SSH")
	var remoteURL string

	if protocol == "HTTPS (Token)" {
		ghUser := os.Getenv("GITHUB_USER")
		ghToken := os.Getenv("GITHUB_TOKEN")
		if ghUser == "" || ghToken == "" {
			ui.PrintError("GITHUB_USER and GITHUB_TOKEN required for HTTPS!")
			return false
		}
		repoName := ui.GumInput("Enter repo name (example: my-repo)", "")
		if repoName == "" {
			return false
		}
		repoName = strings.TrimSuffix(repoName, ".git") + ".git"
		remoteURL = fmt.Sprintf("https://%s@github.com/%s/%s", ghToken, ghUser, repoName)
	} else if protocol == "SSH" {
		remoteURL = ui.GumInput("Enter SSH URL (e.g., git@github.com:user/repo.git)", "")
		if remoteURL == "" {
			return false
		}
	} else {
		ui.PrintWarning("No protocol selected.")
		return false
	}

	if ui.GumConfirm(fmt.Sprintf("Set remote '%s' to %s?", remoteName, remoteURL)) {
		ui.PrintCommand(fmt.Sprintf(`git remote remove "%s"`, remoteName))
		RunGitCmd("remote", "remove", remoteName)
		ui.PrintCommand(fmt.Sprintf(`git remote add "%s" "%s"`, remoteName, remoteURL))
		RunGitCmd("remote", "add", remoteName, remoteURL)
		ui.PrintStatus(fmt.Sprintf(" Remote '%s' set to %s", remoteName, remoteURL))
		return true
	}
	ui.PrintStatus("Operation cancelled.")
	return false
}

// ShowRepoStatus parses git status in a single-pass porcelain v2 execution (< 5ms) using pure Nerd Font icons
func ShowRepoStatus() {
	cmd := exec.Command("git", "status", "--porcelain=v2", "--branch")
	outBytes, err := cmd.Output()
	if err != nil {
		ui.PrintWarning("Not inside a git repository!")
		return
	}

	branch := "main"
	ahead := "0"
	behind := "0"
	hasUncommitted := false
	hasConflicts := false

	scanner := bufio.NewScanner(strings.NewReader(string(outBytes)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# branch.head ") {
			branch = strings.TrimPrefix(line, "# branch.head ")
		} else if strings.HasPrefix(line, "# branch.ab ") {
			abStr := strings.TrimPrefix(line, "# branch.ab ")
			parts := strings.Fields(abStr)
			if len(parts) >= 2 {
				ahead = strings.TrimPrefix(parts[0], "+")
				behind = strings.TrimPrefix(parts[1], "-")
			}
		} else if strings.HasPrefix(line, "u ") {
			hasConflicts = true
			hasUncommitted = true
		} else if len(line) > 0 && (line[0] == '1' || line[0] == '2' || line[0] == '?') {
			hasUncommitted = true
		}
	}

	cPrimary := config.GetAnsiEsc(config.Theme["BIMAGIC_PRIMARY"])
	cSuccess := config.GetAnsiEsc(config.Theme["BIMAGIC_SUCCESS"])
	cWarning := config.GetAnsiEsc(config.Theme["BIMAGIC_WARNING"])
	cError := config.GetAnsiEsc(config.Theme["BIMAGIC_ERROR"])
	cMuted := config.GetAnsiEsc(config.Theme["BIMAGIC_MUTED"])
	nc := "\033[0m"

	statusText := cWarning + "󰀦 uncommitted" + nc
	borderColor := config.Theme["BIMAGIC_PRIMARY"]

	if hasConflicts {
		statusText = cError + "󰅖 conflicts" + nc
		borderColor = config.Theme["BIMAGIC_ERROR"]
	} else if !hasUncommitted {
		statusText = cSuccess + "󰄬 clean" + nc
		borderColor = config.Theme["BIMAGIC_SUCCESS"]
	}

	displayUser := os.Getenv("GITHUB_USER")
	if displayUser == "" {
		displayUser = "SSH/Local"
	}

	header := fmt.Sprintf("%s BIMAGIC GIT WIZARD%s  %s%s%s", cPrimary, nc, cMuted, config.Version, nc)
	line1 := fmt.Sprintf("%s USER    :%s %s", cMuted, nc, displayUser)
	line2 := fmt.Sprintf("%s BRANCH  :%s %s%s%s", cMuted, nc, cPrimary, branch, nc)
	line3 := fmt.Sprintf("%s⇅ SYNC    :%s ⇡%s  ⇣%s", cMuted, nc, ahead, behind)
	line4 := fmt.Sprintf("%s󱖫 STATUS  :%s %s", cMuted, nc, statusText)

	DrawRoundedBox(borderColor, header, line1, line2, line3, line4)
}

func displayCols(s string) int {
	cleanStr := stripAnsi(s)
	w := 0
	for _, r := range cleanStr {
		if r == '🪄' || r == '🟡' || r == '🟢' || r == '🔴' || r == '✨' || r == 0x26A0 || r == 0xFE0F || (r >= 0x1F000 && r <= 0x1FAFF) || (r >= 0x2600 && r <= 0x27BF) {
			w += 2
		} else {
			w += 1
		}
	}
	return w
}

// DrawRoundedBox renders styled rounded borders natively in Go (0ms overhead) with precise column alignment
func DrawRoundedBox(borderHexOrAnsi string, lines ...string) {
	c := config.GetAnsiEsc(borderHexOrAnsi)
	nc := "\033[0m"

	maxLen := 0
	for _, l := range lines {
		w := displayCols(l)
		if w > maxLen {
			maxLen = w
		}
	}
	padding := 4
	width := maxLen + (padding * 2)
	minWidth := 40
	if width < minWidth {
		width = minWidth
	}

	fmt.Println()
	fmt.Printf("%s╭%s╮%s\n", c, strings.Repeat("─", width), nc)
	fmt.Printf("%s│%s%s│%s\n", c, strings.Repeat(" ", width), c, nc)
	for _, l := range lines {
		w := displayCols(l)
		padRight := width - padding - w
		if padRight < 0 {
			padRight = 0
		}
		fmt.Printf("%s│%s%s%s%s%s│%s\n", c, nc, strings.Repeat(" ", padding), l, strings.Repeat(" ", padRight), c, nc)
	}
	fmt.Printf("%s│%s%s│%s\n", c, strings.Repeat(" ", width), c, nc)
	fmt.Printf("%s╰%s╯%s\n", c, strings.Repeat("─", width), nc)
}
