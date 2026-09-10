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
	if git.RunGitCmd("add", ".") != nil {
		ui.PrintError("Failed to add files.")
		os.Exit(1)
	}
	ui.PrintStatus("Files added.")

	// Check if there are staged changes to commit
	hasStaged := git.RunGitCmd("diff", "--cached", "--quiet") != nil

	if hasStaged {
		ui.PrintCommand(`git commit -m "` + cliMsg + `"`)
		if git.RunGitCmd("commit", "-m", cliMsg) == nil {
			ui.PrintStatus("Committed: " + cliMsg)
		} else {
			ui.PrintError("Commit failed.")
			os.Exit(1)
		}
	} else {
		// No staged changes - check if we have any unpushed commits
		hasUnpushed := false
		if upstream := git.GetGitOutput("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"); upstream != "" {
			unpushedCount := git.GetGitOutput("rev-list", "@{u}..HEAD", "--count")
			if unpushedCount != "" && unpushedCount != "0" {
				hasUnpushed = true
			}
		} else {
			if git.RunGitCmd("rev-parse", "HEAD") == nil {
				hasUnpushed = true
			}
		}

		if !hasUnpushed {
			ui.PrintWarning("Nothing to commit, working tree clean. Everything up to date!")
			return
		}
		ui.PrintWarning("No new changes to commit. Proceeding to push existing commits...")
	}

	branch := git.GetCurrentBranch()
	ui.PrintStatus("Pushing to " + branch + "...")

	ui.PrintCommand("git push")
	if ui.GumSpin("Pushing...", "git", "push") {
		ui.PrintStatus("󱝂 Magic complete!")
		return
	}

	// If push failed, check if upstream was not set
	hasUpstream := git.RunGitCmd("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}") == nil
	if !hasUpstream {
		ui.PrintWarning("Standard push failed (no upstream). Trying to set upstream...")
		ui.PrintCommand(`git push -u origin "` + branch + `"`)
		if ui.GumSpin("Pushing (upstream)...", "git", "push", "-u", "origin", branch) {
			ui.PrintStatus("󱝂 Magic complete (upstream set)!")
			return
		}
	}

	ui.PrintError("Push failed. You may need to pull remote changes first (git pull) or check your remote permissions.")
	os.Exit(1)
}
