package spells

import (
	"fmt"
	"os"
	"strings"

	"bimagic-go/pkg/config"
	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func AddFilesLogic() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	for {
		addChoice := ui.GumChoose("", "", "", " Stage Files", "󰷊 Preview Files", "󰌍 Back")
		if addChoice == "󰌍 Back" {
			break
		}
		if addChoice == "󰷊 Preview Files" {
			ScryingGlass()
			continue
		}

		out := git.GetGitOutput("ls-files", "--others", "--modified", "--exclude-standard")
		filesList := "[ALL]\n" + out
		files := ui.GumFilterStdin(filesList, "Select files to add", true)

		if files == "" {
			ui.PrintWarning("No files selected.")
		} else {
			if strings.Contains(files, "[ALL]") {
				ui.PrintCommand("git add .")
				git.RunGitCmd("add", ".")
				ui.PrintStatus("All files staged.")
			} else {
				for _, f := range strings.Split(files, "\n") {
					if f != "" {
						ui.PrintCommand(`git add "` + f + `"`)
						git.RunGitCmd("add", f)
					}
				}
				ui.PrintStatus("Selected files staged.")
				fmt.Println(files)
			}
		}
		break
	}
}

func RemoveFilesLogic() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	for {
		removeChoice := ui.GumChoose("", "", "", "Remove Files", "Preview Files", "Back")
		if removeChoice == "Back" {
			break
		}
		if removeChoice == "Preview Files" {
			ScryingGlass()
			continue
		}

		out := git.GetGitOutput("ls-files", "--cached", "--others", "--exclude-standard")
		files := ui.GumFilterStdin(out, "Select files/folders to remove", true)

		if files == "" {
			ui.PrintWarning("No files selected.")
			break
		}

		fmt.Println("Files selected for removal:")
		fmt.Printf("\033[33m%s\033[0m\n\n", files)

		if ui.GumConfirm("Confirm removal? This cannot be undone.") {
			for _, f := range strings.Split(files, "\n") {
				if f == "" {
					continue
				}
				if git.RunGitCmd("ls-files", "--error-unmatch", f) == nil {
					ui.PrintCommand(`git rm -rf "` + f + `"`)
					git.RunGitCmd("rm", "-rf", f)
				} else {
					ui.PrintCommand(`rm -rf "` + f + `"`)
					os.RemoveAll(f)
				}
			}
			ui.PrintStatus("Selected files/folders have been removed.")
		} else {
			ui.PrintStatus("Operation cancelled.")
		}
		break
	}
}

func UninitializeRepo() {
	ui.PrintWarning("This will completely uninitialize the Git repository in this folder.")
	fmt.Println("This action will delete the .git directory and cannot be undone!")
	fmt.Println()

	if ui.GumConfirm("Are you sure you want to continue?") {
		if _, err := os.Stat(".git"); !os.IsNotExist(err) {
			ui.PrintCommand("rm -rf .git")
			os.RemoveAll(".git")
			ui.PrintStatus("Git repository has been uninitialized.")
		} else {
			ui.PrintError("No .git directory found here. Nothing to do.")
		}
	} else {
		ui.PrintStatus("Operation cancelled.")
	}
}

func DrawGitGraphBox() {
	yellow := config.GetAnsiEsc(config.Theme["BIMAGIC_WARNING"])
	nc := "\033[0m"
	line1 := "Git graph"
	line2 := "[INFO] press 'q' to exit"

	fmt.Printf("%s╭%s╮%s\n", yellow, strings.Repeat("─", 30), nc)
	fmt.Printf("%s│%s %-28s %s│%s\n", yellow, nc, line1, yellow, nc)
	fmt.Printf("%s│%s %-28s %s│%s\n", yellow, nc, line2, yellow, nc)
	fmt.Printf("%s╰%s╯%s\n", yellow, strings.Repeat("─", 30), nc)
}

func PrettyGitLog() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	ui.PrintCommand("git log --graph --oneline --decorate --all")
	git.RunCmdOutToScreen("git", "log", "--graph", "--abbrev-commit", "--decorate", "--date=short", "--format=%C(auto)%h%Creset %C(blue)%ad%Creset %C(green)%an%Creset %C(yellow)%d%Creset %Creset%s", "--all")
}
