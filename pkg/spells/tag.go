package spells

import (
	"fmt"
	"strings"

	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func TagOperations() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	for {
		choice := ui.GumChoose("Tag Operations", "", "",
			"󰓹 List tags",
			"󰐕 Create tag",
			" Push tags to remote",
			"󰅖 Delete local tag",
			"󰮘 Delete remote tag",
			"󰌍 Back",
		)

		switch choice {
		case "󰓹 List tags":
			tags := git.GetGitOutput("tag", "-l", "-n1")
			if tags == "" {
				ui.PrintWarning("No tags found.")
			} else {
				ui.PrintStatus("Repository Tags:")
				fmt.Println(tags)
			}
		case "󰐕 Create tag":
			tagName := ui.GumInput("Enter tag name (e.g. v1.0.0)", "")
			if tagName == "" {
				ui.PrintWarning("Tag name is required.")
				continue
			}
			tagMsg := ui.GumInput("Enter tag message (optional for annotated tag)", "")
			if tagMsg != "" {
				ui.PrintCommand(fmt.Sprintf(`git tag -a "%s" -m "%s"`, tagName, tagMsg))
				if git.RunGitCmdInteractive("tag", "-a", tagName, "-m", tagMsg) == nil {
					ui.PrintStatus("Annotated tag created: " + tagName)
				} else {
					ui.PrintError("Failed to create annotated tag.")
				}
			} else {
				ui.PrintCommand(fmt.Sprintf(`git tag "%s"`, tagName))
				if git.RunGitCmdInteractive("tag", tagName) == nil {
					ui.PrintStatus("Lightweight tag created: " + tagName)
				} else {
					ui.PrintError("Failed to create lightweight tag.")
				}
			}
		case " Push tags to remote":
			remotesStr := git.GetGitOutput("remote")
			if remotesStr == "" {
				ui.PrintError("No remote configured.")
				continue
			}
			remote := "origin"
			remotes := strings.Split(remotesStr, "\n")
			if len(remotes) > 1 {
				remote = ui.GumFilterStdin(remotesStr, "Select remote to push tags", false)
			}
			if remote != "" {
				ui.PrintCommand(fmt.Sprintf(`git push "%s" --tags`, remote))
				if git.RunGitCmdInteractive("push", remote, "--tags") == nil {
					ui.PrintStatus("Tags pushed successfully!")
				} else {
					ui.PrintError("Failed to push tags.")
				}
			}
		case "󰅖 Delete local tag":
			tags := git.GetGitOutput("tag", "-l")
			if tags == "" {
				ui.PrintWarning("No tags found to delete.")
				continue
			}
			selectedTag := ui.GumFilterStdin(tags, "Select local tag to delete", false)
			if selectedTag != "" {
				if ui.GumConfirm("Delete local tag '" + selectedTag + "'?") {
					ui.PrintCommand(`git tag -d "` + selectedTag + `"`)
					git.RunGitCmd("tag", "-d", selectedTag)
					ui.PrintStatus("Local tag deleted: " + selectedTag)
				}
			}
		case "󰮘 Delete remote tag":
			tags := git.GetGitOutput("tag", "-l")
			if tags == "" {
				ui.PrintWarning("No tags found.")
				continue
			}
			selectedTag := ui.GumFilterStdin(tags, "Select remote tag to delete", false)
			if selectedTag != "" {
				remotesStr := git.GetGitOutput("remote")
				if remotesStr == "" {
					ui.PrintError("No remote configured.")
					continue
				}
				remote := "origin"
				remotes := strings.Split(remotesStr, "\n")
				if len(remotes) > 1 {
					remote = ui.GumFilterStdin(remotesStr, "Select remote", false)
				}
				if remote != "" && ui.GumConfirm(fmt.Sprintf("Delete remote tag '%s' from '%s'?", selectedTag, remote)) {
					ui.PrintCommand(fmt.Sprintf(`git push "%s" --delete "%s"`, remote, selectedTag))
					if git.RunGitCmdInteractive("push", remote, "--delete", selectedTag) == nil {
						ui.PrintStatus("Remote tag deleted: " + selectedTag)
					} else {
						ui.PrintError("Failed to delete remote tag.")
					}
				}
			}
		case "󰌍 Back":
			return
		}
		fmt.Println()
	}
}
