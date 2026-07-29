package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"bimagic-go/pkg/config"
	"bimagic-go/pkg/git"
	"bimagic-go/pkg/spells"
	"bimagic-go/pkg/ui"
)

func getTerminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err == nil {
		parts := strings.Fields(string(out))
		if len(parts) >= 2 {
			if w, err := strconv.Atoi(parts[1]); err == nil && w > 0 {
				return w
			}
		}
	}
	return 80
}

func showHelp() {
	c := config.GetAnsiEsc(config.Theme["BIMAGIC_PRIMARY"])
	y := config.GetAnsiEsc(config.Theme["BIMAGIC_WARNING"])
	g := config.GetAnsiEsc(config.Theme["BIMAGIC_MUTED"])
	nc := "\033[0m"

	fmt.Printf("%s󱝁 Bimagic Git Wizard %s - Direct Keymaps Reference%s\n", c, config.Version, nc)
	fmt.Printf("%sUsage:%s wz [flag] [options]\n", y, nc)

	keymaps := []struct {
		Flag string
		Desc string
	}{
		{"-d, --clone [url] [-i]", "Clone repository (optional -i interactive sparse)"},
		{"-I, --init", "Initialize a new Git repository (main default)"},
		{"-A, --add", "Stage files interactively"},
		{"-U, --unstage", "Unstage files (git restore --staged)"},
		{"-X, --discard", "Discard local uncommitted edits (git checkout --)"},
		{"-c, --commit", "Magic Commit (Conventional Commits builder)"},
		{"-P, --push", "Push local commits to remote repository"},
		{"-p, --pull", "Pull latest changes from remote"},
		{"-b, --branch", "Branch operations (switch, create, rename, delete)"},
		{"-t, --tag", "Tag operations (create, list, push, delete)"},
		{"-D, --diff", "Diff & inspection wizard (unstaged, staged, file, branch)"},
		{"-C, --cherry", "Cherry-pick commits onto current branch"},
		{"-r, --remote", "Configure HTTPS token or SSH remotes"},
		{"-s, --status", "Show status dashboard (<5ms single-pass)"},
		{"-S, --stats", "Contributor statistics and activity analysis"},
		{"-g, --graph", "Display pretty git log tree graph"},
		{"-a, --architect", "Summon the Architect (.gitignore generator)"},
		{"-R, --remove", "Safely remove files/folders with git integration"},
		{"-m, --merge", "Merge branches with conflict detection"},
		{"--uninit", "Uninitialize Git repository (remove .git)"},
		{"-k, --resurrect", "Resurrection Stone (recover lost reflog commits)"},
		{"-v, --revert", "Revert one or more commits (multi-select)"},
		{"-w, --stash", "Stash operations (push, pop, list, apply, drop, clear)"},
		{"-q, --quickview", "The Scrying Glass (instant file browser)"},
		{"-z, --lazy [msg]", "The Lazy Wizard (Add + Commit + Push)"},
		{"-h, --help", "Show this power user direct keymap guide"},
	}

	termWidth := getTerminalWidth()

	if termWidth >= 120 {
		// Desktop Wide Screen: 2-Column Side-by-Side Table
		flagW := 24
		descW := 31

		half := (len(keymaps) + 1) / 2
		fmt.Println()
		fmt.Printf("%s╭%s┬%s┬%s┬%s╮%s\n", c, strings.Repeat("─", flagW+2), strings.Repeat("─", descW+2), strings.Repeat("─", flagW+2), strings.Repeat("─", descW+2), nc)
		fmt.Printf("%s│%s %-*s %s│%s %-*s %s│%s %-*s %s│%s %-*s %s│%s\n",
			c, y, flagW, "DIRECT KEYMAP", c, y, descW, "DESCRIPTION", c, y, flagW, "DIRECT KEYMAP", c, y, descW, "DESCRIPTION", c, nc)
		fmt.Printf("%s├%s┼%s┼%s┼%s┤%s\n", c, strings.Repeat("─", flagW+2), strings.Repeat("─", descW+2), strings.Repeat("─", flagW+2), strings.Repeat("─", descW+2), nc)

		for i := 0; i < half; i++ {
			left := keymaps[i]
			leftFlag := left.Flag
			if len(leftFlag) > flagW {
				leftFlag = leftFlag[:flagW]
			}
			leftDesc := left.Desc
			if len(leftDesc) > descW {
				leftDesc = leftDesc[:descW-3] + "..."
			}

			rightFlag := ""
			rightDesc := ""
			if i+half < len(keymaps) {
				right := keymaps[i+half]
				rightFlag = right.Flag
				if len(rightFlag) > flagW {
					rightFlag = rightFlag[:flagW]
				}
				rightDesc = right.Desc
				if len(rightDesc) > descW {
					rightDesc = rightDesc[:descW-3] + "..."
				}
			}

			fmt.Printf("%s│%s %-*s %s│%s %-*s %s│%s %-*s %s│%s %-*s %s│%s\n",
				c, c, flagW, leftFlag, c, g, descW, leftDesc, c, c, flagW, rightFlag, c, g, descW, rightDesc, c, nc)
		}
		fmt.Printf("%s╰%s┴%s┴%s┴%s╯%s\n", c, strings.Repeat("─", flagW+2), strings.Repeat("─", descW+2), strings.Repeat("─", flagW+2), strings.Repeat("─", descW+2), nc)

	} else if termWidth >= 80 {
		// Medium Screen: Single Column Table
		flagW := 28
		descW := 55

		fmt.Println()
		fmt.Printf("%s╭%s┬%s╮%s\n", c, strings.Repeat("─", flagW+2), strings.Repeat("─", descW+2), nc)
		fmt.Printf("%s│%s %-*s %s│%s %-*s %s│%s\n", c, y, flagW, "DIRECT KEYMAP FLAG", c, y, descW, "COMMAND DESCRIPTION", c, nc)
		fmt.Printf("%s├%s┼%s┤%s\n", c, strings.Repeat("─", flagW+2), strings.Repeat("─", descW+2), nc)

		for _, k := range keymaps {
			flagStr := k.Flag
			if len(flagStr) > flagW {
				flagStr = flagStr[:flagW]
			}
			descStr := k.Desc
			if len(descStr) > descW {
				descStr = descStr[:descW-3] + "..."
			}
			fmt.Printf("%s│%s %-*s %s│%s %-*s %s│%s\n", c, c, flagW, flagStr, c, g, descW, descStr, c, nc)
		}
		fmt.Printf("%s╰%s┴%s╯%s\n", c, strings.Repeat("─", flagW+2), strings.Repeat("─", descW+2), nc)
	} else {
		// Small Screen: Compact List Format
		fmt.Println()
		for _, k := range keymaps {
			fmt.Printf("  %s%-28s%s %s%s%s\n", c, k.Flag, nc, g, k.Desc, nc)
		}
		fmt.Println()
	}
}

func main() {
	// Enable multi-core CPU parallelism
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load Config & Theme
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "bimagic")
	themeFile := filepath.Join(configDir, "theme.wz")
	config.LoadTheme(themeFile)

	// Handle Version / Help Flags
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "--version" || arg == "-version" {
			fmt.Printf("Bimagic Git Wizard %s\n", config.Version)
			os.Exit(0)
		} else if arg == "-h" || arg == "--help" {
			showHelp()
			os.Exit(0)
		}
	}

	// Ensure gum is installed
	if !ui.HasCmd("gum") {
		fmt.Println("Error: gum is not installed.")
		fmt.Println("Please install it: https://github.com/charmbracelet/gum")
		os.Exit(1)
	}

	// Parse CLI arguments
	var cliMode, cliURL, cliMsg, cliDepth string
	cliInteractive := false

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d", "--clone":
			cliMode = "clone"
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				cliURL = args[i+1]
				i++
			}
		case "--depth":
			if i+1 < len(args) {
				cliDepth = args[i+1]
				i++
			}
		case "-i", "--interactive":
			cliInteractive = true
		case "-z", "--lazy":
			cliMode = "lazy"
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				cliMsg = args[i+1]
				i++
			}
		case "-I", "--init":
			cliMode = "init"
		case "-A", "--add":
			cliMode = "add"
		case "-U", "--unstage":
			cliMode = "unstage"
		case "-X", "--discard":
			cliMode = "discard"
		case "-c", "--commit":
			cliMode = "commit"
		case "-P", "--push":
			cliMode = "push"
		case "-p", "--pull":
			cliMode = "pull"
		case "-b", "--branch":
			cliMode = "branch"
		case "-t", "--tag":
			cliMode = "tag"
		case "-D", "--diff":
			cliMode = "diff"
		case "-C", "--cherry":
			cliMode = "cherry"
		case "-r", "--remote":
			cliMode = "remote"
		case "-s", "--status":
			cliMode = "status"
		case "-S", "--stats":
			cliMode = "stats"
		case "-g", "--graph":
			cliMode = "graph"
		case "-a", "--architect":
			cliMode = "architect"
		case "-R", "--remove":
			cliMode = "remove"
		case "-m", "--merge":
			cliMode = "merge"
		case "--uninit":
			cliMode = "uninit"
		case "-k", "--resurrect":
			cliMode = "resurrect"
		case "-v", "--revert":
			cliMode = "revert"
		case "-w", "--stash":
			cliMode = "stash"
		case "-q", "--quickview", "--scrying":
			cliMode = "scrying"
		default:
			if cliMode == "clone" && cliURL == "" {
				cliURL = args[i]
			} else if cliMode == "lazy" && cliMsg == "" {
				cliMsg = args[i]
			}
		}
	}

	// Dispatch Direct Power User CLI Keymaps
	switch cliMode {
	case "clone":
		if cliURL == "" {
			ui.PrintError("Error: Repository URL required with -d / --clone")
			os.Exit(1)
		}
		spells.CloneRepo(cliURL, cliInteractive, cliDepth)
		os.Exit(0)
	case "init":
		spells.InitRepo()
		os.Exit(0)
	case "add":
		spells.AddFilesLogic()
		os.Exit(0)
	case "unstage":
		spells.UnstageFilesLogic()
		os.Exit(0)
	case "discard":
		spells.DiscardChangesLogic()
		os.Exit(0)
	case "commit":
		spells.CommitWizard()
		os.Exit(0)
	case "push":
		spells.PushToRemote()
		os.Exit(0)
	case "pull":
		spells.PullChangesInteractive()
		os.Exit(0)
	case "branch":
		spells.CreateSwitchBranch()
		os.Exit(0)
	case "tag":
		spells.TagOperations()
		os.Exit(0)
	case "diff":
		spells.DiffWizard()
		os.Exit(0)
	case "cherry":
		spells.CherryPickWizard()
		os.Exit(0)
	case "remote":
		git.SetupRemote("origin")
		os.Exit(0)
	case "status":
		git.ShowRepoStatus()
		os.Exit(0)
	case "stats":
		spells.ShowContributorStats()
		os.Exit(0)
	case "graph":
		if !git.IsGitRepo() {
			ui.PrintError("Not a git repository!")
			os.Exit(1)
		}
		spells.PrettyGitLog()
		os.Exit(0)
	case "architect":
		spells.SummonGitignore()
		os.Exit(0)
	case "remove":
		spells.RemoveFilesLogic()
		os.Exit(0)
	case "merge":
		spells.MergeBranches()
		os.Exit(0)
	case "uninit":
		spells.UninitializeRepo()
		os.Exit(0)
	case "resurrect":
		spells.ResurrectCommit()
		os.Exit(0)
	case "revert":
		spells.RevertCommits()
		os.Exit(0)
	case "stash":
		spells.StashOperations()
		os.Exit(0)
	case "scrying":
		spells.ScryingGlass()
		os.Exit(0)
	case "undo":
		spells.TimeTurner()
		os.Exit(0)
	case "lazy":
		spells.LazyWizard(cliMsg)
		os.Exit(0)
	}

	// Warn if credentials are not set
	if os.Getenv("GITHUB_USER") == "" || os.Getenv("GITHUB_TOKEN") == "" {
		ui.PrintWarning("GITHUB_USER or GITHUB_TOKEN not set. Defaulting to SSH/Local mode.")
	}

	// Welcome Banner Logic
	versionFile := filepath.Join(configDir, "version")
	storedVersion := ""
	if b, err := os.ReadFile(versionFile); err == nil {
		storedVersion = strings.TrimSpace(string(b))
	}

	if storedVersion != config.Version {
		config.ShowWelcomeBanner(versionFile, configDir, storedVersion == "")
	} else {
		fmt.Println("Welcome to the Git Wizard! Let's work some magic...")
		fmt.Println()
	}

	// Interactive Main Loop
	for {
		ui.ClearScreen()
		git.ShowRepoStatus()

		options := []string{
			" Clone repository",
			" Init new repo",
			" Add / Stage files",
			"󰁯 Unstage files",
			"󰮘 Discard local modifications",
			" Commit changes",
			" Push to remote",
			" Pull latest changes",
			" Branch operations (Switch, Create, Rename, Delete)",
			"󰓹 Tag operations (Create, List, Push, Delete)",
			"󰈈 Diff & inspection wizard",
			" Cherry-pick commits",
			" Set remote",
			"󱖫 Show status",
			" Contributor Statistics",
			" Git graph",
			"󰓗 Summon the Architect (.gitignore)",
			"󰮘 Remove files/folders (rm)",
			" Merge branches",
			" Uninitialize repo",
			"󰔪 Summon the Resurrection Stone (Recover lost code)",
			"󰁯 Revert commit(s)",
			"󰓗 Stash operations",
			"󰈈 The Scrying Glass (Quick View)",
			"󰿅 Exit",
		}

		choice := ui.GumChoose(
			" Choose your spell: (j/k to navigate)",
			" ",
			config.Theme["BIMAGIC_PRIMARY"],
			options...,
		)
		fmt.Println()

		switch choice {
		case " Clone repository":
			repoURL := ui.GumInput("Enter repository URL", "")
			if repoURL == "" {
				continue
			}
			repoDepth := ui.GumInput("Enter depth (empty for full clone)", "")
			cloneMode := ui.GumChoose("", "", "", "Standard Clone", "Interactive (Select files)")
			spells.CloneRepo(repoURL, cloneMode == "Interactive (Select files)", repoDepth)

		case "󰓗 Stash operations":
			spells.StashOperations()

		case " Init new repo":
			spells.InitRepo()

		case " Add / Stage files":
			spells.AddFilesLogic()

		case "󰁯 Unstage files":
			spells.UnstageFilesLogic()

		case "󰮘 Discard local modifications":
			spells.DiscardChangesLogic()

		case "󰓹 Tag operations (Create, List, Push, Delete)":
			spells.TagOperations()

		case "󰈈 Diff & inspection wizard":
			spells.DiffWizard()

		case " Cherry-pick commits":
			spells.CherryPickWizard()

		case " Commit changes":
			commitMode := ui.GumChoose("", "", "", "󰦥 Magic Commit (Builder)", "󱐋 Quick Commit (One-line)")
			if commitMode == "󰦥 Magic Commit (Builder)" {
				spells.CommitWizard()
			} else {
				msg := ui.GumInput("Enter commit message", "")
				if msg == "" {
					ui.PrintWarning("No commit message provided. Cancelled.")
					continue
				}
				if ui.GumConfirm("Commit changes?") {
					ui.PrintCommand(`git commit -m "` + msg + `"`)
					err := git.RunGitCmd("commit", "-m", msg)
					if err == nil {
						ui.PrintStatus("Commit done!")
					} else {
						ui.PrintStatus("Commit cancelled.")
					}
				} else {
					ui.PrintStatus("Commit cancelled.")
				}
			}

		case " Push to remote":
			spells.PushToRemote()

		case " Pull latest changes":
			spells.PullChangesInteractive()

		case " Branch operations (Switch, Create, Rename, Delete)":
			spells.CreateSwitchBranch()

		case " Set remote":
			git.SetupRemote("origin")

		case "󱖫 Show status":
			git.RunCmdOutToScreen("git", "status")

		case "󰮘 Remove files/folders (rm)":
			spells.RemoveFilesLogic()

		case " Uninitialize repo":
			spells.UninitializeRepo()

		case "󰈈 The Scrying Glass (Quick View)":
			spells.ScryingGlass()

		case "󰿅 Exit":
			if ui.GumConfirm("Are you sure you want to exit?") {
				fmt.Println("Git Wizard vanishes in a puff of smoke...")
				os.Exit(0)
			} else {
				continue
			}

		case " Merge branches":
			spells.MergeBranches()

		case " Contributor Statistics":
			spells.ShowContributorStats()

		case "󰓗 Summon the Architect (.gitignore)":
			spells.SummonGitignore()

		case "󰔪 Summon the Resurrection Stone (Recover lost code)":
			spells.ResurrectCommit()

		case "󰁯 Revert commit(s)":
			spells.RevertCommits()

		case " Git graph":
			spells.DrawGitGraphBox()
			ui.GumSpin("Drawing git graph...", "sleep", "2")
			spells.PrettyGitLog()

		default:
			ui.PrintError("Invalid choice! Try again.")
			fmt.Println("Git Wizard vanishes in a puff of smoke...")
			break
		}

		fmt.Println()
		ui.GumStyleWithArgs(fmt.Sprintf("--foreground=%s", config.Theme["BIMAGIC_MUTED"]), "Press Enter to continue...")
		ui.WaitForEnter()
	}
}
