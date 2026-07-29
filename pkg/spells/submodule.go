package spells

import (
	"fmt"
	"strings"

	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func SubmoduleManager() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}

	subStatus := git.GetGitOutput("submodule", "status")
	if subStatus != "" {
		fmt.Println("Active Submodules:")
		fmt.Println(subStatus)
		fmt.Println()
	} else {
		ui.PrintWarning("No submodules configured in this repository.")
		fmt.Println()
	}

	choice := ui.GumChoose("󰓗 Submodule Operations:", "", "",
		"󰁯 Init & Update Submodules (git submodule update --init --recursive)",
		"󰓗 Sync Submodules (git submodule sync --recursive)",
		" Add New Submodule (git submodule add <url> <path>)",
		"󱖫 Show Submodule Status",
	)

	switch choice {
	case "󰁯 Init & Update Submodules (git submodule update --init --recursive)":
		ui.PrintCommand("git submodule update --init --recursive")
		ui.GumSpin("Updating submodules...", "git", "submodule", "update", "--init", "--recursive")
		ui.PrintStatus("Submodules initialized and updated successfully!")

	case "󰓗 Sync Submodules (git submodule sync --recursive)":
		ui.PrintCommand("git submodule sync --recursive")
		ui.GumSpin("Syncing submodules...", "git", "submodule", "sync", "--recursive")
		ui.PrintStatus("Submodules synced successfully!")

	case " Add New Submodule (git submodule add <url> <path>)":
		subURL := ui.GumInput("Enter submodule repository URL", "")
		if subURL == "" {
			return
		}
		subPath := ui.GumInput("Enter local path for submodule (e.g., vendor/lib)", "")
		if subPath == "" {
			return
		}
		ui.PrintCommand(fmt.Sprintf("git submodule add %s %s", subURL, subPath))
		err := git.RunGitCmd("submodule", "add", subURL, subPath)
		if err == nil {
			ui.PrintStatus(fmt.Sprintf("Added submodule '%s' at '%s'!", strings.TrimSuffix(subURL, ".git"), subPath))
		} else {
			ui.PrintError("Failed to add submodule.")
		}

	case "󱖫 Show Submodule Status":
		git.RunCmdOutToScreen("git", "submodule", "status")
	}
}
