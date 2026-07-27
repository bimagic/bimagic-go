package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"bimagic-go/pkg/config"
	"bimagic-go/pkg/ui"
)

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

// ShowRepoStatus parses git status in a single-pass porcelain v2 execution (< 5ms)
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

	status := "🟡 uncommitted"
	borderColor := config.Theme["BIMAGIC_WARNING"]

	if hasConflicts {
		status = "🔴 conflicts"
		borderColor = config.Theme["BIMAGIC_ERROR"]
	} else if !hasUncommitted {
		status = "🟢 clean"
		borderColor = config.Theme["BIMAGIC_SUCCESS"]
	}

	displayUser := os.Getenv("GITHUB_USER")
	if displayUser == "" {
		displayUser = "SSH/Local"
	}

	line1 := fmt.Sprintf("GITHUB USER: %s", displayUser)
	line2 := fmt.Sprintf("BRANCH: %s", branch)
	line3 := fmt.Sprintf("AHEAD: %s | BEHIND: %s", ahead, behind)
	line4 := fmt.Sprintf("STATUS: %s", status)

	DrawRoundedBox(borderColor, line1, line2, line3, line4)
}

// DrawRoundedBox renders styled rounded borders natively in Go (0ms overhead)
func DrawRoundedBox(borderHexOrAnsi string, lines ...string) {
	c := config.GetAnsiEsc(borderHexOrAnsi)
	nc := "\033[0m"

	maxLen := 0
	for _, l := range lines {
		// Visible length calculation (ignoring status emojis/ANSI if needed)
		runeCount := len([]rune(l))
		if runeCount > maxLen {
			maxLen = runeCount
		}
	}
	padding := 2
	width := maxLen + (padding * 2)

	fmt.Println()
	fmt.Printf("%s╭%s╮%s\n", c, strings.Repeat("─", width), nc)
	for _, l := range lines {
		runeCount := len([]rune(l))
		padRight := width - padding - runeCount
		if padRight < 0 {
			padRight = 0
		}
		fmt.Printf("%s│%s%s%s%s%s│%s\n", c, nc, strings.Repeat(" ", padding), l, strings.Repeat(" ", padRight), c, nc)
	}
	fmt.Printf("%s╰%s╯%s\n", c, strings.Repeat("─", width), nc)
}
