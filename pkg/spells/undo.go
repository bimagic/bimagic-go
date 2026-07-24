package spells

import (
	"os"
	"strings"

	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func TimeTurner() {
	ui.PrintStatus("󱦟 Spinning the Time Turner...")

	if err := git.RunGitCmd("rev-parse", "HEAD"); err != nil {
		ui.PrintError("No commits to undo! This repo is empty.")
		os.Exit(1)
	}

	isInitialCommit := git.RunGitCmd("rev-parse", "HEAD~1") != nil

	undoType := ui.GumChoose("Select Undo Level:", "", "",
		"Soft (Undo commit, keep changes staged - Best for fixing typos)",
		"Mixed (Undo commit, keep changes unstaged - Best for splitting work)",
		"Hard (DESTROY changes - Revert to previous state)",
		"Cancel",
	)

	switch {
	case strings.HasPrefix(undoType, "Soft"):
		if isInitialCommit {
			ui.PrintCommand("git update-ref -d HEAD")
			git.RunGitCmd("update-ref", "-d", "HEAD")
		} else {
			ui.PrintCommand("git reset --soft HEAD~1")
			git.RunGitCmd("reset", "--soft", "HEAD~1")
		}
		ui.PrintStatus("✨ Success! I undid the commit, but kept your files ready to commit again.")
	case strings.HasPrefix(undoType, "Mixed"):
		if isInitialCommit {
			ui.PrintCommand("git update-ref -d HEAD")
			git.RunGitCmd("update-ref", "-d", "HEAD")
			ui.PrintCommand("git rm --cached -r -q .")
			git.RunGitCmd("rm", "--cached", "-r", "-q", ".")
		} else {
			ui.PrintCommand("git reset HEAD~1")
			git.RunGitCmd("reset", "HEAD~1")
		}
		ui.PrintStatus("󱞈 Success! I undid the commit and unstaged the files.")
	case strings.HasPrefix(undoType, "Hard"):
		if ui.GumConfirm(" DANGER: This deletes your work forever. Are you sure?") {
			if isInitialCommit {
				ui.PrintCommand("git update-ref -d HEAD")
				git.RunGitCmd("update-ref", "-d", "HEAD")
				ui.PrintCommand("git rm --cached -r -q .")
				git.RunGitCmd("rm", "--cached", "-r", "-q", ".")
				ui.PrintCommand("git clean -fd")
				git.RunGitCmd("clean", "-fd")
			} else {
				ui.PrintCommand("git reset --hard HEAD~1")
				git.RunGitCmd("reset", "--hard", "HEAD~1")
			}
			ui.PrintStatus("󱠇 Obliviate! The last commit and its changes are destroyed.")
		} else {
			ui.PrintStatus("Operation cancelled.")
		}
	default:
		ui.PrintStatus("Mischief managed (Cancelled).")
	}
}
