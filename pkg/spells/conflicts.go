package spells

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func ConflictAssistant() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}

	conflictedRaw := git.GetGitOutput("diff", "--name-only", "--diff-filter=U")
	if strings.TrimSpace(conflictedRaw) == "" {
		ui.PrintStatus("󰄬 No merge conflicts detected in working tree!")
		return
	}

	files := strings.Split(strings.TrimSpace(conflictedRaw), "\n")
	ui.PrintWarning(fmt.Sprintf("󰅖 Detected %d conflicted file(s):", len(files)))
	for _, f := range files {
		fmt.Printf("  • %s\n", f)
	}
	fmt.Println()

	selectedFile := ui.GumChoose("Select file to resolve conflict:", "", "", append(files, "󰿅 Done / Return")...)
	if selectedFile == "" || selectedFile == "󰿅 Done / Return" {
		return
	}

	action := ui.GumChoose(fmt.Sprintf("Resolve conflict for '%s':", selectedFile), "", "",
		"󰄬 Accept Ours / Current Branch (git checkout --ours)",
		"󰄬 Accept Theirs / Incoming Branch (git checkout --theirs)",
		"󰏫 Open in $EDITOR (Manual edit)",
		"󰈈 View conflict markers (git diff)",
	)

	switch action {
	case "󰄬 Accept Ours / Current Branch (git checkout --ours)":
		ui.PrintCommand(fmt.Sprintf("git checkout --ours %s", selectedFile))
		if err := git.RunGitCmd("checkout", "--ours", selectedFile); err == nil {
			git.RunGitCmd("add", selectedFile)
			ui.PrintStatus(fmt.Sprintf("Resolved '%s' with Ours (staged)!", selectedFile))
		} else {
			ui.PrintError("Failed to checkout ours for " + selectedFile)
		}

	case "󰄬 Accept Theirs / Incoming Branch (git checkout --theirs)":
		ui.PrintCommand(fmt.Sprintf("git checkout --theirs %s", selectedFile))
		if err := git.RunGitCmd("checkout", "--theirs", selectedFile); err == nil {
			git.RunGitCmd("add", selectedFile)
			ui.PrintStatus(fmt.Sprintf("Resolved '%s' with Theirs (staged)!", selectedFile))
		} else {
			ui.PrintError("Failed to checkout theirs for " + selectedFile)
		}

	case "󰏫 Open in $EDITOR (Manual edit)":
		editor := os.Getenv("EDITOR")
		if editor == "" {
			if ui.HasCmd("nvim") {
				editor = "nvim"
			} else if ui.HasCmd("vim") {
				editor = "vim"
			} else {
				editor = "nano"
			}
		}
		ui.PrintCommand(fmt.Sprintf("%s %s", editor, selectedFile))
		cmd := exec.Command(editor, selectedFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		if ui.GumConfirm(fmt.Sprintf("Mark '%s' as resolved and stage?", selectedFile)) {
			git.RunGitCmd("add", selectedFile)
			ui.PrintStatus(fmt.Sprintf("Staged '%s' as resolved!", selectedFile))
		}

	case "󰈈 View conflict markers (git diff)":
		git.RunCmdOutToScreen("git", "diff", selectedFile)
	}

	remaining := git.GetGitOutput("diff", "--name-only", "--diff-filter=U")
	if strings.TrimSpace(remaining) == "" {
		ui.PrintStatus("󰄬 All conflicts resolved! You can now commit or continue rebase/merge.")
		if ui.GumConfirm("Commit resolved merge changes now?") {
			spellsCommit()
		}
	}
}

func spellsCommit() {
	msg := ui.GumInput("Enter merge commit message", "Merge resolved conflicts")
	if msg != "" {
		git.RunGitCmd("commit", "-m", msg)
		ui.PrintStatus("Merge commit completed!")
	}
}
