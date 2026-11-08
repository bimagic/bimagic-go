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
		if err := git.RunGitCmdInteractive("submodule", "update", "--init", "--recursive"); err == nil {
			ui.PrintStatus("Submodules initialized and updated successfully!")
		} else {
			ui.PrintError("Failed to initialize or update submodules.")
		}

	case "󰓗 Sync Submodules (git submodule sync --recursive)":
		ui.PrintCommand("git submodule sync --recursive")
		if err := git.RunGitCmdInteractive("submodule", "sync", "--recursive"); err == nil {
			ui.PrintStatus("Submodules synced successfully!")
		} else {
			ui.PrintError("Failed to sync submodules.")
		}

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
		err := git.RunGitCmdInteractive("submodule", "add", subURL, subPath)
		if err == nil {
			ui.PrintStatus(fmt.Sprintf("Added submodule '%s' at '%s'!", strings.TrimSuffix(subURL, ".git"), subPath))
		} else {
			ui.PrintError("Failed to add submodule.")
		}

	case "󱖫 Show Submodule Status":
		git.RunCmdOutToScreen("git", "submodule", "status")
	}
}
