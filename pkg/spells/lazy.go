package spells

import (
	"fmt"
	"os"

	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func LazyWizard(cliMsg string) {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		os.Exit(1)
	}
	if cliMsg == "" {
		ui.PrintError("Error: Commit message required for Lazy Wizard (-z)")
		fmt.Println("Usage: bimagic -z \"commit message\"")
		os.Exit(1)
	}

	ui.PrintStatus("  Lazy Wizard invoked!")

	ui.PrintCommand("git add .")
	if ui.GumSpin("Adding files...", "git", "add", ".") {
		ui.PrintStatus("Files added.")
	} else {
		ui.PrintError("Failed to add files.")
		os.Exit(1)
	}

	ui.PrintCommand(`git commit -m "` + cliMsg + `"`)
	if git.RunGitCmd("commit", "-m", cliMsg) == nil {
		ui.PrintStatus("Committed: " + cliMsg)
	} else {
		if git.RunGitCmd("rev-parse", "HEAD") == nil {
			ui.PrintCommand(`git commit --amend -m "` + cliMsg + `"`)
			if git.RunGitCmd("commit", "--amend", "-m", cliMsg) == nil {
				ui.PrintStatus("Updated last commit message to: " + cliMsg)
			} else {
				ui.PrintWarning("No new changes to commit. Proceeding to push...")
			}
		} else {
			ui.PrintWarning("No commits or changes to commit. Proceeding to push...")
		}
	}

	branch := git.GetCurrentBranch()
	ui.PrintStatus("Pushing to " + branch + "...")

	ui.PrintCommand("git push")
	if ui.GumSpin("Pushing...", "git", "push") {
		ui.PrintStatus("󱝂 Magic complete!")
	} else {
		ui.PrintWarning("Standard push failed. Trying to set upstream...")
		ui.PrintCommand(`git push -u origin "` + branch + `"`)
		if ui.GumSpin("Pushing (upstream)...", "git", "push", "-u", "origin", branch) {
			ui.PrintStatus("󱝂 Magic complete (upstream set)!")
		} else {
			ui.PrintError("Push failed.")
			os.Exit(1)
		}
	}
}
