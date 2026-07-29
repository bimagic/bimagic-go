package spells

import (
	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func BlameWizard() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}

	filesRaw := git.GetGitOutput("ls-files")
	if filesRaw == "" {
		ui.PrintError("No tracked files found in repository.")
		return
	}

	selectedFile := ui.GumFilterStdin(filesRaw, "Select file to inspect line blame & author history", false)
	if selectedFile == "" {
		return
	}

	ui.PrintCommand("git blame -w " + selectedFile)
	if ui.HasCmd("bat") {
		git.RunCmdOutToScreen("sh", "-c", "git blame -w "+selectedFile+" | bat --language=log --paging=always")
	} else {
		git.RunCmdOutToScreen("sh", "-c", "git blame -w "+selectedFile+" | less -R")
	}
}
