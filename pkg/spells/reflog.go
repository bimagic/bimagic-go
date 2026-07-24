package spells

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"bimagic-go/pkg/config"
	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func ResurrectCommit() {
	ui.PrintStatus("󰔪  Summoning the Resurrection Stone...")

	reflog := git.GetGitOutput("reflog")
	if reflog == "" {
		ui.PrintError("The ancient logs are empty. Nothing to resurrect.")
		return
	}

	fmt.Println("Search the ancient timelines for your lost code:")
	formattedReflog := git.GetGitOutput("reflog", "--date=relative", "--format=%h %gd %C(blue)%ad%Creset %s")
	selectedLog := ui.GumFilterStdin(formattedReflog, "Search for a lost commit message or hash...", false)

	if selectedLog == "" {
		ui.PrintStatus("Resurrection cancelled.")
		return
	}

	targetHash := strings.Fields(selectedLog)[0]

	fmt.Println()
	cmd := exec.Command("gum", "style", "--border", "rounded", "--border-foreground", config.Theme["BIMAGIC_PRIMARY"], "--padding", "1 2", "TARGET TIMELINE:", selectedLog)
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println()

	action := ui.GumChoose("", "", "", "󰔱 Create a new branch here (Safest)", "  Hard Reset current branch to here (Dangerous)", "Cancel")

	switch action {
	case "󰔱 Create a new branch here (Safest)":
		newBranch := ui.GumInput("Enter new branch name (e.g., recovered-code)", "")
		if newBranch != "" {
			ui.PrintCommand(`git checkout -b "` + newBranch + `" "` + targetHash + `"`)
			git.RunGitCmd("checkout", "-b", newBranch, targetHash)
			ui.PrintStatus("󱝁 Timeline restored! You are now on branch: " + newBranch)
		}
	case "  Hard Reset current branch to here (Dangerous)":
		if ui.GumConfirm("This will overwrite your CURRENT work. Are you absolutely sure?") {
			ui.PrintCommand("git reset --hard " + targetHash)
			git.RunGitCmd("reset", "--hard", targetHash)
			ui.PrintStatus("󱝁 Timeline restored via hard reset!")
		}
	default:
		ui.PrintStatus("The stone goes dormant.")
	}
}
