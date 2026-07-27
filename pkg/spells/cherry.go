package spells

import (
	"fmt"
	"strings"

	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func CherryPickWizard() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}

	currentBranch := git.GetCurrentBranch()
	ui.PrintStatus("Current branch: " + currentBranch)
	fmt.Println()

	logOut := git.GetGitOutput("log", "--all", "--oneline", "--decorate", "-n", "100")
	if logOut == "" {
		ui.PrintError("No commits found in history.")
		return
	}

	selected := ui.GumFilterStdin(logOut, "Select commit to cherry-pick into "+currentBranch, false)
	if selected == "" {
		ui.PrintWarning("Cherry-pick cancelled.")
		return
	}

	commitHash := strings.Fields(selected)[0]

	if ui.GumConfirm(fmt.Sprintf("Cherry-pick commit %s onto '%s'?", commitHash, currentBranch)) {
		ui.PrintCommand("git cherry-pick " + commitHash)
		if git.RunGitCmd("cherry-pick", commitHash) == nil {
			ui.PrintStatus(fmt.Sprintf("Commit %s cherry-picked successfully!", commitHash))
		} else {
			ui.PrintError("Cherry-pick encountered conflicts!")
			fmt.Println("Please resolve conflicts, then run:")
			fmt.Println("  git cherry-pick --continue")
		}
	} else {
		ui.PrintStatus("Cherry-pick cancelled.")
	}
}
