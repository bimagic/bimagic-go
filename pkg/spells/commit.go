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

func CommitWizard() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}

	nc := "\033[0m"
	fmt.Printf("%s=== The Alchemist's Commit ===%s\n", config.GetAnsiEsc(config.Theme["BIMAGIC_PRIMARY"]), nc)

	typeStr := ui.GumChoose("Select change type:", "", "",
		"feat: A new feature",
		"fix: A bug fix",
		"docs: Documentation only changes",
		"style: Changes that do not affect the meaning of the code",
		"refactor: A code change that neither fixes a bug nor adds a feature",
		"perf: A code change that improves performance",
		"test: Adding missing tests or correcting existing tests",
		"chore: Changes to the build process or auxiliary tools",
	)
	if typeStr == "" {
		return
	}
	typeStr = strings.SplitN(typeStr, ":", 2)[0]

	scope := ui.GumInput("Scope (optional, e.g., 'login', 'ui'). Press Enter to skip.", "")
	summary := ui.GumInput("Short description (imperative mood, e.g., 'add generic login')", "")
	if summary == "" {
		ui.PrintWarning("Summary is required!")
		return
	}

	body := ""
	if ui.GumConfirm("Add a longer description (body)?") {
		body = ui.GumWrite("Enter detailed description...")
	}

	breaking := ""
	if ui.GumConfirm("Is this a BREAKING CHANGE?") {
		breaking = "!"
	}

	commitMsg := typeStr
	if scope != "" {
		commitMsg += "(" + scope + ")"
	}
	commitMsg += breaking + ": " + summary

	if body != "" {
		commitMsg += "\n\n" + body
	}

	fmt.Println()
	cmd := exec.Command("gum", "style", "--border", "rounded", "--border-foreground", config.Theme["BIMAGIC_PRIMARY"], "--padding", "1 2", "PREVIEW:", commitMsg)
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println()

	if ui.GumConfirm("Commit with this message?") {
		ui.PrintCommand(`git commit -m "` + commitMsg + `"`)
		git.RunGitCmd("commit", "-m", commitMsg)
		ui.PrintStatus("󱝁 Mischief managed! (Commit successful)")
	} else {
		ui.PrintWarning("Commit cancelled.")
	}
}
