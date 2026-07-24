package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

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

func ShowRepoStatus() {
	if !IsGitRepo() {
		ui.PrintWarning("Not inside a git repository!")
		return
	}

	var (
		wg         sync.WaitGroup
		branch     string
		upstream   string
		ahead      string = "0"
		behind     string = "0"
		cleanIndex bool
		cleanCache bool
		conflicts  string
	)

	wg.Add(4)

	// Goroutine 1: Branch and upstream status
	go func() {
		defer wg.Done()
		branch = GetCurrentBranch()
		upstream = GetGitOutput("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")

		if upstream != "" {
			var abWg sync.WaitGroup
			abWg.Add(2)
			go func() {
				defer abWg.Done()
				if a := GetGitOutput("rev-list", "--count", upstream+"..HEAD"); a != "" {
					ahead = a
				}
			}()
			go func() {
				defer abWg.Done()
				if b := GetGitOutput("rev-list", "--count", "HEAD.."+upstream); b != "" {
					behind = b
				}
			}()
			abWg.Wait()
		}
	}()

	// Goroutine 2: Working tree diff check
	go func() {
		defer wg.Done()
		cleanIndex = RunGitCmd("diff", "--quiet") == nil
	}()

	// Goroutine 3: Cached index diff check
	go func() {
		defer wg.Done()
		cleanCache = RunGitCmd("diff", "--cached", "--quiet") == nil
	}()

	// Goroutine 4: Unmerged conflicts check
	go func() {
		defer wg.Done()
		conflicts = GetGitOutput("ls-files", "-u")
	}()

	wg.Wait()

	status := "🟡 uncommitted"
	color := config.Theme["BIMAGIC_WARNING"]

	if cleanIndex && cleanCache {
		if conflicts != "" {
			status = "🔴 conflicts"
			color = config.Theme["BIMAGIC_ERROR"]
		} else {
			status = "🟢 clean"
			color = config.Theme["BIMAGIC_SUCCESS"]
		}
	}

	displayUser := os.Getenv("GITHUB_USER")
	if displayUser == "" {
		displayUser = "SSH/Local"
	}

	content := fmt.Sprintf("GITHUB USER: %s\nBRANCH: %s\nAHEAD: %s | BEHIND: %s\nSTATUS: %s", displayUser, branch, ahead, behind, status)

	fmt.Println()
	cmd := exec.Command("gum", "style", "--border", "rounded", "--margin", "1 0", "--padding", "1 2", "--border-foreground", color, content)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
