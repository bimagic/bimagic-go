package spells

import (
	"fmt"
	"strconv"
	"strings"

	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func RebaseWizard() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}

	mode := ui.GumChoose("󰊢 Choose Rebase & Squash Operation:", "", "",
		"󰮘 Quick Squash (Combine last N commits into one)",
		"󰏫 Interactive Rebase (git rebase -i HEAD~N)",
		"󰁯 Rebase onto target branch (git rebase <branch>)",
		"󰅖 Abort ongoing rebase (git rebase --abort)",
	)

	switch mode {
	case "󰮘 Quick Squash (Combine last N commits into one)":
		logOut := git.GetGitOutput("log", "-n", "10", "--oneline")
		if logOut == "" {
			ui.PrintError("No commits found to squash.")
			return
		}
		fmt.Println("Recent commits:")
		fmt.Println(logOut)
		fmt.Println()

		numStr := ui.GumInput("How many commits back to squash into 1? (e.g., 3)", "2")
		n, err := strconv.Atoi(numStr)
		if err != nil || n < 2 {
			ui.PrintWarning("Squash requires combining at least 2 commits.")
			return
		}

		newMsg := ui.GumInput("Enter commit message for squashed commit", "")
		if newMsg == "" {
			ui.PrintWarning("No commit message provided. Cancelled.")
			return
		}

		if ui.GumConfirm(fmt.Sprintf("Squash last %d commits into one with message '%s'?", n, newMsg)) {
			ui.PrintCommand(fmt.Sprintf("git reset --soft HEAD~%d", n))
			if err := git.RunGitCmd("reset", "--soft", fmt.Sprintf("HEAD~%d", n)); err != nil {
				ui.PrintError("Failed to soft reset commits.")
				return
			}
			ui.PrintCommand(fmt.Sprintf(`git commit -m "%s"`, newMsg))
			if err := git.RunGitCmd("commit", "-m", newMsg); err == nil {
				ui.PrintStatus(fmt.Sprintf("Successfully squashed last %d commits into 1!", n))
			} else {
				ui.PrintError("Failed to create squashed commit.")
			}
		}

	case "󰏫 Interactive Rebase (git rebase -i HEAD~N)":
		numStr := ui.GumInput("How many commits back to open in interactive rebase? (e.g., 5)", "5")
		n, err := strconv.Atoi(numStr)
		if err != nil || n < 1 {
			return
		}
		ui.PrintCommand(fmt.Sprintf("git rebase -i HEAD~%d", n))
		git.RunCmdOutToScreen("git", "rebase", "-i", fmt.Sprintf("HEAD~%d", n))

	case "󰁯 Rebase onto target branch (git rebase <branch>)":
		branchesRaw := git.GetGitOutput("branch", "-a")
		if branchesRaw == "" {
			ui.PrintError("No branches available.")
			return
		}
		selectedBranch := ui.GumFilterStdin(branchesRaw, "Select target branch to rebase onto", false)
		selectedBranch = strings.TrimSpace(strings.TrimPrefix(selectedBranch, "*"))
		if selectedBranch == "" {
			return
		}
		if ui.GumConfirm(fmt.Sprintf("Rebase current branch onto '%s'?", selectedBranch)) {
			ui.PrintCommand(fmt.Sprintf("git rebase %s", selectedBranch))
			err := git.RunGitCmd("rebase", selectedBranch)
			if err == nil {
				ui.PrintStatus("Successfully rebased onto " + selectedBranch + "!")
			} else {
				ui.PrintError("Rebase encountered conflicts! Use 'wz --conflicts' to resolve.")
			}
		}

	case "󰅖 Abort ongoing rebase (git rebase --abort)":
		ui.PrintCommand("git rebase --abort")
		err := git.RunGitCmd("rebase", "--abort")
		if err == nil {
			ui.PrintStatus("Rebase aborted.")
		} else {
			ui.PrintError("No ongoing rebase to abort.")
		}
	}
}
