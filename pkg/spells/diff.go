package spells

import (
	"fmt"

	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func DiffWizard() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	for {
		choice := ui.GumChoose("Diff & Inspection Wizard", "", "",
			"🔍 View Unstaged Changes (git diff)",
			"📦 View Staged Changes (git diff --staged)",
			"📄 View Specific File Diff",
			"🔀 Compare Branches",
			"󰌍 Back",
		)

		switch choice {
		case "🔍 View Unstaged Changes (git diff)":
			ui.PrintCommand("git diff")
			git.RunCmdOutToScreen("git", "diff")
		case "📦 View Staged Changes (git diff --staged)":
			ui.PrintCommand("git diff --staged")
			git.RunCmdOutToScreen("git", "diff", "--staged")
		case "📄 View Specific File Diff":
			files := git.GetGitOutput("ls-files", "-m", "-d", "-o", "--exclude-standard")
			if files == "" {
				ui.PrintWarning("No modified files to inspect.")
				continue
			}
			file := ui.GumFilterStdin(files, "Select file to view diff", false)
			if file != "" {
				ui.PrintCommand("git diff " + file)
				git.RunCmdOutToScreen("git", "diff", file)
			}
		case "🔀 Compare Branches":
			branches := git.GetGitOutput("branch", "--format=%(refname:short)")
			if branches == "" {
				ui.PrintWarning("No branches found.")
				continue
			}
			b1 := ui.GumFilterStdin(branches, "Select base branch (branch 1)", false)
			if b1 == "" {
				continue
			}
			b2 := ui.GumFilterStdin(branches, "Select compare branch (branch 2)", false)
			if b2 == "" {
				continue
			}
			ui.PrintCommand(fmt.Sprintf("git diff %s..%s", b1, b2))
			git.RunCmdOutToScreen("git", "diff", fmt.Sprintf("%s..%s", b1, b2))
		case "󰌍 Back":
			return
		}
		fmt.Println()
	}
}
