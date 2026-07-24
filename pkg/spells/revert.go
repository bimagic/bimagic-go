package spells

import (
	"fmt"
	"strings"

	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func RevertCommits() {
	ui.PrintStatus("Fetching commit history...")
	fmt.Println()

	history := git.GetGitOutput("log", "--oneline", "--decorate")
	commitsSelection := ui.GumFilterStdin(history, "Select commit(s) to revert", true)

	if commitsSelection == "" {
		ui.PrintWarning("No commit selected. Revert cancelled.")
		return
	}

	var hashes []string
	for _, line := range strings.Split(commitsSelection, "\n") {
		if line != "" {
			hashes = append(hashes, strings.Fields(line)[0])
		}
	}

	fmt.Println("You selected:")
	fmt.Println(strings.Join(hashes, "\n"))
	fmt.Println()

	if ui.GumConfirm("Confirm revert?") {
		for _, c := range hashes {
			fmt.Printf("Reverting commit %s...\n", c)
			ui.PrintCommand(fmt.Sprintf("git revert --no-edit %s", c))
			if git.RunGitCmd("revert", "--no-edit", c) == nil {
				ui.PrintStatus(fmt.Sprintf("Commit %s reverted.", c))
			} else {
				ui.PrintError(fmt.Sprintf("Conflict occurred while reverting %s!", c))
				fmt.Println("Please resolve conflicts, then run:")
				fmt.Println("  git revert --continue")
				break
			}
		}
	} else {
		ui.PrintStatus("Revert cancelled.")
	}
}
