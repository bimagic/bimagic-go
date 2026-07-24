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

func StashOperations() {
	for {
		stashChoice := ui.GumChoose("Stash Operations", "", "",
			"󱉛 Push (Save) changes",
			"󱉙 Pop latest stash",
			" List stashes",
			" Apply specific stash",
			" Drop specific stash",
			"󰎟 Clear all stashes",
			"󰌍 Back",
		)

		switch stashChoice {
		case "󱉛 Push (Save) changes":
			msg := ui.GumInput("Optional stash message", "")
			includeUntracked := ""
			if ui.GumConfirm("Include untracked files?") {
				includeUntracked = "-u"
			}
			cmdStr := "git stash push"
			if includeUntracked != "" {
				cmdStr += " " + includeUntracked
			}
			if msg != "" {
				cmdStr += ` -m "` + msg + `"`
			}
			ui.PrintCommand(cmdStr)
			args := []string{"stash", "push"}
			if includeUntracked != "" {
				args = append(args, includeUntracked)
			}
			if msg != "" {
				args = append(args, "-m", msg)
			}
			if git.RunGitCmd(args...) == nil {
				ui.PrintStatus("Changes stashed successfully!")
			} else {
				ui.PrintError("Failed to stash changes.")
			}
		case "󱉙 Pop latest stash":
			ui.PrintCommand("git stash pop")
			if git.RunGitCmd("stash", "pop") == nil {
				ui.PrintStatus("Stash popped successfully!")
			} else {
				ui.PrintError("Failed to pop stash (possible conflicts).")
			}
		case " List stashes":
			stashes := git.GetGitOutput("stash", "list")
			if stashes == "" {
				ui.PrintWarning("No stashes found.")
			} else {
				cmd := exec.Command("gum", "style", "--border", "normal", "--border-foreground", config.Theme["BIMAGIC_PRIMARY"], "--padding", "0 1", stashes)
				cmd.Stdout = os.Stdout
				cmd.Run()
			}
		case " Apply specific stash":
			stashes := git.GetGitOutput("stash", "list")
			if stashes == "" {
				ui.PrintWarning("No stashes found.")
				continue
			}
			stashEntry := ui.GumFilterStdin(stashes, "Select stash to apply", false)
			if stashEntry != "" {
				stashID := strings.SplitN(stashEntry, ":", 2)[0]
				ui.PrintCommand("git stash apply " + stashID)
				if git.RunGitCmd("stash", "apply", stashID) == nil {
					ui.PrintStatus("Applied " + stashID)
				} else {
					ui.PrintError("Failed to apply " + stashID)
				}
			}
		case " Drop specific stash":
			stashes := git.GetGitOutput("stash", "list")
			if stashes == "" {
				ui.PrintWarning("No stashes found.")
				continue
			}
			stashEntry := ui.GumFilterStdin(stashes, "Select stash to drop", false)
			if stashEntry != "" {
				stashID := strings.SplitN(stashEntry, ":", 2)[0]
				if ui.GumConfirm("Are you sure you want to drop " + stashID + "?") {
					ui.PrintCommand("git stash drop " + stashID)
					if git.RunGitCmd("stash", "drop", stashID) == nil {
						ui.PrintStatus("Dropped " + stashID)
					} else {
						ui.PrintError("Failed to drop " + stashID)
					}
				}
			}
		case "󰎟 Clear all stashes":
			if git.GetGitOutput("stash", "list") == "" {
				ui.PrintWarning("No stashes found.")
				continue
			}
			if ui.GumConfirm("DANGER: This will delete ALL stashes. Continue?") {
				ui.PrintCommand("git stash clear")
				if git.RunGitCmd("stash", "clear") == nil {
					ui.PrintStatus("All stashes cleared.")
				} else {
					ui.PrintError("Failed to clear stashes.")
				}
			} else {
				ui.PrintStatus("Operation cancelled.")
			}
		case "󰌍 Back":
			return
		}
		fmt.Println()
		ui.GumStyleWithArgs(fmt.Sprintf("--foreground=%s", config.Theme["BIMAGIC_MUTED"]), "Press Enter to continue...")
		ui.WaitForEnter()
	}
}
