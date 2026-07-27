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

func main() {
	// Enable multi-core CPU parallelism
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 1. Handle Version Flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("Bimagic Git Wizard %s\n", config.Version)
		os.Exit(0)
	}

	// 2. Load Config & Theme
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "bimagic")
	themeFile := filepath.Join(configDir, "theme.wz")
	config.LoadTheme(themeFile)

	// 3. Ensure gum is installed
	if !ui.HasCmd("gum") {
		fmt.Println("Error: gum is not installed.")
		fmt.Println("Please install it: https://github.com/charmbracelet/gum")
		os.Exit(1)
	}

	// 4. Parse CLI arguments
	var cliMode, cliURL, cliMsg, cliDepth string
	cliInteractive := false

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-d":
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
		case "-i":
			cliInteractive = true
		case "-z":
			cliMode = "lazy"
			if i+1 < len(args) {
				cliMsg = args[i+1]
				i++
			}
		case "-s":
			cliMode = "status"
		case "-u":
			cliMode = "undo"
		case "-g":
			cliMode = "graph"
		case "-p":
			cliMode = "pull"
		case "-a", "--architect":
			cliMode = "architect"
		case "-t", "--tag":
			cliMode = "tag"
		case "--diff":
			cliMode = "diff"
		default:
			if cliMode == "clone" && cliURL == "" {
				cliURL = args[i]
			} else if cliMode == "lazy" && cliMsg == "" {
				cliMsg = args[i]
			}
		}
	}

	// 5. Handle direct CLI Modes
	switch cliMode {
	case "clone":
		if cliURL == "" {
			ui.PrintError("Error: Repository URL required with -d")
			os.Exit(1)
		}
		spells.CloneRepo(cliURL, cliInteractive, cliDepth)
		os.Exit(0)
	case "status":
		git.ShowRepoStatus()
		os.Exit(0)
	case "pull":
		spells.PullLatestChanges()
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
	case "undo":
		spells.TimeTurner()
		os.Exit(0)
	case "lazy":
		spells.LazyWizard(cliMsg)
		os.Exit(0)
	case "tag":
		spells.TagOperations()
		os.Exit(0)
	case "diff":
		spells.DiffWizard()
		os.Exit(0)
	}

	// Warn if credentials are not set
	if os.Getenv("GITHUB_USER") == "" || os.Getenv("GITHUB_TOKEN") == "" {
		ui.PrintWarning("GITHUB_USER or GITHUB_TOKEN not set. Defaulting to SSH/Local mode.")
	}

	// 6. Welcome Banner Logic
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

	// 7. Interactive Main Loop
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
