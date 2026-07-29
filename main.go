package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"bimagic-go/pkg/config"
	"bimagic-go/pkg/git"
	"bimagic-go/pkg/spells"
	"bimagic-go/pkg/ui"
)

func showHelp() {
	c := config.GetAnsiEsc(config.Theme["BIMAGIC_PRIMARY"])
	y := config.GetAnsiEsc(config.Theme["BIMAGIC_WARNING"])
	g := config.GetAnsiEsc(config.Theme["BIMAGIC_MUTED"])
	nc := "\033[0m"

	fmt.Printf("%sBimagic Git Wizard %s - Direct Keymaps & CLI Reference%s\n\n", c, config.Version, nc)
	fmt.Printf("%sUsage:%s wz [flag] [options]\n\n", y, nc)

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

	for _, k := range keymaps {
		fmt.Printf("  %s%-28s%s %s%s%s\n", c, k.Flag, nc, g, k.Desc, nc)
	}
	fmt.Println()
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

	// Interactive Main Loop with Keybind Badges
	for {
		ui.ClearScreen()
		git.ShowRepoStatus()

		options := []string{
			"[d]  Clone repository",
			"[I]  Init new repo",
			"[A]  Add / Stage files",
			"[U] 󰁯 Unstage files",
			"[X] 󰮘 Discard local modifications",
			"[c]  Commit changes",
			"[P]  Push to remote",
			"[p]  Pull latest changes",
			"[b]  Branch operations (Switch, Create, Rename, Delete)",
			"[t] 󰓹 Tag operations (Create, List, Push, Delete)",
			"[D] 󰈈 Diff & inspection wizard",
			"[C]  Cherry-pick commits",
			"[r]  Set remote",
			"[s] 󱖫 Show status",
			"[S]  Contributor Statistics",
			"[g]  Git graph",
			"[a] 󰓗 Summon the Architect (.gitignore)",
			"[R] 󰮘 Remove files/folders (rm)",
			"[m]  Merge branches",
			"[u]  Uninitialize repo",
			"[k] 󰔪 Summon the Resurrection Stone (Recover lost code)",
			"[v] 󰁯 Revert commit(s)",
			"[w] 󰓗 Stash operations",
			"[q] 󰈈 The Scrying Glass (Quick View)",
			"[x] 󰿅 Exit",
		}

		choice := ui.GumChoose(
			" Choose your spell: (type keybind or j/k to navigate)",
			" ",
			config.Theme["BIMAGIC_PRIMARY"],
			options...,
		)
		fmt.Println()

		switch choice {
		case "[d]  Clone repository":
			repoURL := ui.GumInput("Enter repository URL", "")
			if repoURL == "" {
				continue
			}
			repoDepth := ui.GumInput("Enter depth (empty for full clone)", "")
			cloneMode := ui.GumChoose("", "", "", "Standard Clone", "Interactive (Select files)")
			spells.CloneRepo(repoURL, cloneMode == "Interactive (Select files)", repoDepth)

		case "[w] 󰓗 Stash operations":
			spells.StashOperations()

		case "[I]  Init new repo":
			spells.InitRepo()

		case "[A]  Add / Stage files":
			spells.AddFilesLogic()

		case "[U] 󰁯 Unstage files":
			spells.UnstageFilesLogic()

		case "[X] 󰮘 Discard local modifications":
			spells.DiscardChangesLogic()

		case "[t] 󰓹 Tag operations (Create, List, Push, Delete)":
			spells.TagOperations()

		case "[D] 󰈈 Diff & inspection wizard":
			spells.DiffWizard()

		case "[C]  Cherry-pick commits":
			spells.CherryPickWizard()

		case "[c]  Commit changes":
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

		case "[P]  Push to remote":
			spells.PushToRemote()

		case "[p]  Pull latest changes":
			spells.PullChangesInteractive()

		case "[b]  Branch operations (Switch, Create, Rename, Delete)":
			spells.CreateSwitchBranch()

		case "[r]  Set remote":
			git.SetupRemote("origin")

		case "[s] 󱖫 Show status":
			git.RunCmdOutToScreen("git", "status")

		case "[R] 󰮘 Remove files/folders (rm)":
			spells.RemoveFilesLogic()

		case "[u]  Uninitialize repo":
			spells.UninitializeRepo()

		case "[q] 󰈈 The Scrying Glass (Quick View)":
			spells.ScryingGlass()

		case "[x] 󰿅 Exit":
			if ui.GumConfirm("Are you sure you want to exit?") {
				fmt.Println("Git Wizard vanishes in a puff of smoke...")
				os.Exit(0)
			} else {
				continue
			}

		case "[m]  Merge branches":
			spells.MergeBranches()

		case "[S]  Contributor Statistics":
			spells.ShowContributorStats()

		case "[a] 󰓗 Summon the Architect (.gitignore)":
			spells.SummonGitignore()

		case "[k] 󰔪 Summon the Resurrection Stone (Recover lost code)":
			spells.ResurrectCommit()

		case "[v] 󰁯 Revert commit(s)":
			spells.RevertCommits()

		case "[g]  Git graph":
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
