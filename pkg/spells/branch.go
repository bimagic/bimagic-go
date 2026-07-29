package spells

import (
	"fmt"
	"strings"
	"sync"

	"bimagic-go/pkg/config"
	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func CreateSwitchBranch() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	var (
		wg             sync.WaitGroup
		currentBranch  string
		branchesOutput string
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		currentBranch = git.GetCurrentBranch()
	}()
	go func() {
		defer wg.Done()
		branchesOutput = git.GetGitOutput("branch", "-a", "--format=%(refname:short)")
	}()
	wg.Wait()

	ui.PrintStatus("Current branch: " + currentBranch)
	fmt.Println()
	ui.PrintStatus("Available branches:")

	branches := strings.Split(branchesOutput, "\n")
	uniqueBranches := make(map[string]bool)

	green := config.GetAnsiEsc(config.Theme["BIMAGIC_SUCCESS"])
	nc := "\033[0m"

	for _, b := range branches {
		b = strings.TrimSpace(b)
		if b == "" || uniqueBranches[b] {
			continue
		}
		uniqueBranches[b] = true
		if b == currentBranch {
			fmt.Printf("%s➤ %s%s (current)\n", green, b, nc)
		} else {
			fmt.Printf("  %s\n", b)
		}
	}
	fmt.Println()

	branchOpt := ui.GumChoose("Branch Operations", "", "",
		"Switch to existing branch",
		"Create new branch",
		"Rename current branch",
		"Delete local branch",
		"Cancel",
	)

	switch branchOpt {
	case "Switch to existing branch":
		branchesOnly := git.GetGitOutput("branch", "--format=%(refname:short)")
		existingBranch := ui.GumFilterStdin(branchesOnly, "Select branch to switch to", false)
		if existingBranch != "" {
			ui.PrintCommand(`git checkout "` + existingBranch + `"`)
			git.RunGitCmd("checkout", existingBranch)
			ui.PrintStatus("Switched to branch: " + existingBranch)
		} else {
			ui.PrintWarning("No branch selected.")
		}
	case "Create new branch":
		newBranch := ui.GumInput("Enter new branch name", "")
		if newBranch != "" {
			ui.PrintCommand(`git checkout -b "` + newBranch + `"`)
			git.RunGitCmd("checkout", "-b", newBranch)
			ui.PrintStatus("Created and switched to new branch: " + newBranch)
		} else {
			ui.PrintError("No branch name provided.")
		}
	case "Rename current branch":
		RenameBranch()
	case "Delete local branch":
		DeleteBranch()
	default:
		ui.PrintWarning("Operation cancelled.")
	}
}

func DeleteBranch() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	currentBranch := git.GetCurrentBranch()
	branchesOut := git.GetGitOutput("branch", "--format=%(refname:short)")

	var candidateBranches []string
	for _, b := range strings.Split(branchesOut, "\n") {
		b = strings.TrimSpace(b)
		if b != "" && b != currentBranch {
			candidateBranches = append(candidateBranches, b)
		}
	}

	if len(candidateBranches) == 0 {
		ui.PrintWarning("No other local branches available to delete.")
		return
	}

	targetBranch := ui.GumFilterStdin(strings.Join(candidateBranches, "\n"), "Select branch to delete", false)
	if targetBranch == "" {
		ui.PrintWarning("No branch selected.")
		return
	}

	force := "-d"
	if ui.GumConfirm(fmt.Sprintf("Force delete branch '%s' (-D)? (Choose No for safe delete -d)", targetBranch)) {
		force = "-D"
	}

	if ui.GumConfirm(fmt.Sprintf("Are you sure you want to delete branch '%s' (%s)?", targetBranch, force)) {
		ui.PrintCommand(fmt.Sprintf(`git branch %s "%s"`, force, targetBranch))
		if git.RunGitCmd("branch", force, targetBranch) == nil {
			ui.PrintStatus(fmt.Sprintf("Branch '%s' deleted successfully.", targetBranch))
		} else {
			ui.PrintError(fmt.Sprintf("Failed to delete branch '%s'.", targetBranch))
		}
	}
}

func RenameBranch() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	currentBranch := git.GetCurrentBranch()
	newName := ui.GumInput(fmt.Sprintf("Enter new name for current branch '%s'", currentBranch), "")
	if newName == "" {
		ui.PrintWarning("No name provided. Rename cancelled.")
		return
	}

	ui.PrintCommand(fmt.Sprintf(`git branch -m "%s"`, newName))
	if git.RunGitCmd("branch", "-m", newName) == nil {
		ui.PrintStatus(fmt.Sprintf("Renamed branch '%s' -> '%s'", currentBranch, newName))
	} else {
		ui.PrintError("Failed to rename branch.")
	}
}

func PushToRemote() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	branch := git.GetCurrentBranch()
	remotesStr := git.GetGitOutput("remote")

	var remote string
	if remotesStr == "" {
		ui.PrintError("No remote set!")
		if git.SetupRemote("origin") {
			remote = "origin"
		} else {
			return
		}
	} else {
		remotes := strings.Split(remotesStr, "\n")
		if len(remotes) == 1 {
			remote = remotes[0]
		} else {
			remote = ui.GumFilterStdin(remotesStr, "Select remote to push to", false)
		}
	}

	if remote == "" {
		return
	}

	if ui.GumConfirm(fmt.Sprintf("Push branch '%s' to '%s'?", branch, remote)) {
		fmt.Printf("Pushing branch '%s' to '%s'...\n", branch, remote)
		ui.PrintCommand(fmt.Sprintf(`git push -u "%s" "%s"`, remote, branch))
		ui.GumSpin("Pushing...", "git", "push", "-u", remote, branch)
	} else {
		ui.PrintStatus("Push cancelled.")
	}
}

func PullLatestChanges() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	ui.PrintCommand("git fetch --all")
	if ui.GumSpin("Fetching updates...", "git", "fetch", "--all") {
		ui.PrintStatus("Fetch complete.")
	} else {
		ui.PrintWarning("Fetch encountered issues during fetch.")
	}

	ui.PrintCommand("git pull --all")
	if ui.GumSpin("Pulling all...", "git", "pull", "--all") {
		ui.PrintStatus("Pull all complete.")
	} else {
		ui.PrintError("Pull failed. There might be conflicts or no upstream set.")
	}
}

func PullChangesInteractive() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	if ui.GumSpin("Fetching updates...", "git", "fetch", "--all") {
		ui.PrintStatus("Fetch complete.")
	} else {
		ui.PrintWarning("Fetch encountered issues.")
	}

	pullChoice := ui.GumChoose("Select pull mode:", "", "",
		" Pull with Rebase (git pull --rebase)",
		" Standard Pull (git pull)",
		" Pull specific branch",
		" Pull all remotes (git pull --all)",
	)

	switch pullChoice {
	case " Pull with Rebase (git pull --rebase)":
		ui.PrintCommand("git pull --rebase")
		err := git.RunGitCmd("pull", "--rebase")
		if err == nil {
			ui.PrintStatus("Successfully pulled and rebased local commits!")
		} else {
			ui.PrintError("Pull with rebase encountered conflicts! Use 'wz --conflicts' to resolve.")
		}

	case " Standard Pull (git pull)":
		ui.PrintCommand("git pull")
		err := git.RunGitCmd("pull")
		if err == nil {
			ui.PrintStatus("Pull complete.")
		} else {
			ui.PrintError("Pull failed.")
		}

	case " Pull all remotes (git pull --all)":
		if ui.GumConfirm("Run 'git pull --all'?") {
			ui.PrintCommand("git pull --all")
			ui.GumSpin("Pulling all...", "git", "pull", "--all")
			ui.PrintStatus("Pull all complete.")
		} else {
			ui.PrintStatus("Pull cancelled.")
		}

	case " Pull specific branch":
		branch := ui.GumInput("Enter branch to pull", "main")
		if branch == "" {
			branch = "main"
		}

		remotesStr := git.GetGitOutput("remote")
		if remotesStr == "" {
			ui.PrintError("No remote set! Cannot pull.")
			return
		}
		var remote string
		remotes := strings.Split(remotesStr, "\n")
		if len(remotes) == 1 {
			remote = remotes[0]
		} else {
			remote = ui.GumFilterStdin(remotesStr, "Select remote to pull from", false)
		}

		if remote != "" {
			if ui.GumConfirm(fmt.Sprintf("Pull branch '%s' from '%s'?", branch, remote)) {
				ui.PrintCommand(fmt.Sprintf(`git pull "%s" "%s"`, remote, branch))
				ui.GumSpin("Pulling...", "git", "pull", remote, branch)
			} else {
				ui.PrintStatus("Pull cancelled.")
			}
		} else {
			ui.PrintWarning("No remote selected.")
		}
	}
}

func MergeBranches() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}
	currentBranch := git.GetCurrentBranch()
	ui.PrintStatus("You are on branch: " + currentBranch)
	fmt.Println()

	branchesOut := git.GetGitOutput("branch", "--format=%(refname:short)")
	filtered := ""
	for _, b := range strings.Split(branchesOut, "\n") {
		if b != currentBranch {
			filtered += b + "\n"
		}
	}

	mergeBranch := ui.GumFilterStdin(filtered, "Select branch to merge into "+currentBranch, false)

	if mergeBranch == "" {
		ui.PrintWarning("No branch selected. Merge cancelled.")
	} else {
		if ui.GumConfirm(fmt.Sprintf("Merge branch '%s' into '%s'?", mergeBranch, currentBranch)) {
			fmt.Printf("Merging branch '%s' into '%s'...\n", mergeBranch, currentBranch)
			ui.PrintCommand(`git merge "` + mergeBranch + `"`)
			if ui.GumSpin("Merging...", "git", "merge", mergeBranch) {
				ui.PrintStatus("Merge successful!")
			} else {
				ui.PrintError("Merge had conflicts! Resolve them manually.")
			}
		} else {
			ui.PrintStatus("Merge cancelled.")
		}
	}
}
