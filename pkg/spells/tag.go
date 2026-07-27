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
			"🏷️  List tags",
			"➕ Create tag",
			"🚀 Push tags to remote",
			"❌ Delete local tag",
			"🔥 Delete remote tag",
			"󰌍 Back",
		)

		switch choice {
		case "🏷️  List tags":
			tags := git.GetGitOutput("tag", "-l", "-n1")
			if tags == "" {
				ui.PrintWarning("No tags found.")
			} else {
				ui.PrintStatus("Repository Tags:")
				fmt.Println(tags)
			}
		case "➕ Create tag":
			tagName := ui.GumInput("Enter tag name (e.g. v1.0.0)", "")
			if tagName == "" {
				ui.PrintWarning("Tag name is required.")
				continue
			}
			tagMsg := ui.GumInput("Enter tag message (optional for annotated tag)", "")
			if tagMsg != "" {
				ui.PrintCommand(fmt.Sprintf(`git tag -a "%s" -m "%s"`, tagName, tagMsg))
				if git.RunGitCmd("tag", "-a", tagName, "-m", tagMsg) == nil {
					ui.PrintStatus("Annotated tag created: " + tagName)
				}
			} else {
				ui.PrintCommand(fmt.Sprintf(`git tag "%s"`, tagName))
				if git.RunGitCmd("tag", tagName) == nil {
					ui.PrintStatus("Lightweight tag created: " + tagName)
				}
			}
		case "🚀 Push tags to remote":
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
				if ui.GumSpin("Pushing tags...", "git", "push", remote, "--tags") {
					ui.PrintStatus("Tags pushed successfully!")
				}
			}
		case "❌ Delete local tag":
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
		case "🔥 Delete remote tag":
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
					git.RunGitCmd("push", remote, "--delete", selectedTag)
					ui.PrintStatus("Remote tag deleted: " + selectedTag)
				}
			}
		case "󰌍 Back":
			return
		}
		fmt.Println()
	}
}
