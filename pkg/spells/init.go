package spells

import (
	"os"

	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func InitRepo() {
	dirname := ui.GumInput("Enter repo directory name (or '.' for current dir)", "")
	if dirname == "" {
		ui.PrintWarning("Operation cancelled.")
		return
	}

	targetDir := dirname
	if dirname == "." {
		cwd, _ := os.Getwd()
		targetDir = cwd
		
		ui.PrintCommand("git init -b main")
		err := git.RunGitCmd("init", "-b", "main")
		if err != nil {
			ui.PrintCommand("git init")
			git.RunGitCmd("init")
		}

		branch := git.GetGitOutput("symbolic-ref", "--short", "HEAD")
		if branch == "master" || branch == "" {
			ui.PrintCommand("git branch -M main")
			git.RunGitCmd("branch", "-M", "main")
		}
		ui.PrintStatus("Repo initialized with 'main' branch in current directory: " + targetDir)
	} else {
		os.MkdirAll(dirname, 0o755)
		originalDir, _ := os.Getwd()
		os.Chdir(dirname)

		ui.PrintCommand("git init -b main")
		err := git.RunGitCmd("init", "-b", "main")
		if err != nil {
			ui.PrintCommand("git init")
			git.RunGitCmd("init")
		}

		branch := git.GetGitOutput("symbolic-ref", "--short", "HEAD")
		if branch == "master" || branch == "" {
			ui.PrintCommand("git branch -M main")
			git.RunGitCmd("branch", "-M", "main")
		}

		os.Chdir(originalDir)
		ui.PrintStatus("Repo initialized with 'main' branch in new directory: " + dirname)
	}
}
