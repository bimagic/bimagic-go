package spells

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func CloneRepo(url string, interactive bool, depth string) {
	repoName := strings.TrimSuffix(filepath.Base(url), ".git")
	if _, err := os.Stat(repoName); !os.IsNotExist(err) {
		ui.PrintError(fmt.Sprintf("Directory '%s' already exists.", repoName))
		return
	}

	depthArgs := []string{}
	if depth != "" {
		depthArgs = []string{"--depth", depth}
	}

	if interactive {
		ui.PrintStatus(fmt.Sprintf("Initializing interactive clone for %s...", repoName))

		args := []string{"clone", "--filter=blob:none", "--no-checkout"}
		args = append(args, depthArgs...)
		args = append(args, url, repoName)

		cmdStr := fmt.Sprintf("git %s", strings.Join(args, " "))
		ui.PrintCommand(cmdStr)
		ui.PrintStatus("Cloning structure for " + repoName + "...")

		args = append([]string{"clone", "--progress", "--filter=blob:none", "--no-checkout"}, depthArgs...)
		args = append(args, url, repoName)
		RunGitCloneWithProgress(args)

		if _, err := os.Stat(repoName); os.IsNotExist(err) {
			ui.PrintError("Clone failed.")
			return
		}

		originalDir, _ := os.Getwd()
		os.Chdir(repoName)

		ui.PrintStatus("Fetching file list...")
		allFiles := git.GetGitOutput("ls-tree", "-r", "--name-only", "HEAD")
		selectedPaths := ui.GumFilterStdin(allFiles, "Select files/folders to download (Space to select)", true)

		if selectedPaths == "" {
			ui.PrintWarning("No files selected. Aborting checkout.")
			os.Chdir(originalDir)
			os.RemoveAll(repoName)
			return
		}

		ui.GumSpin("Configuring sparse checkout...", "git", "sparse-checkout", "init", "--no-cone")

		setCmd := exec.Command("git", "sparse-checkout", "set", "--stdin")
		setCmd.Stdin = strings.NewReader(selectedPaths)
		setCmd.Run()

		ui.PrintStatus("Downloading selected files...")
		RunGitCloneWithProgress([]string{"checkout", "--progress", "HEAD"})

		os.Chdir(originalDir)
		ui.PrintStatus("Successfully cloned selected files into '" + repoName + "'!")

	} else {
		args := []string{"clone"}
		args = append(args, depthArgs...)
		args = append(args, url, repoName)

		cmdStr := fmt.Sprintf("git %s", strings.Join(args, " "))
		ui.PrintCommand(cmdStr)
		ui.PrintStatus(fmt.Sprintf("Cloning %s into %s...", url, repoName))

		args = append([]string{"clone", "--progress"}, depthArgs...)
		args = append(args, url, repoName)
		success := RunGitCloneWithProgress(args)

		if success {
			ui.PrintStatus(fmt.Sprintf("Successfully cloned '%s' into '%s'!", url, repoName))
		} else {
			ui.PrintError("Clone failed.")
		}
	}
}

func RunGitCloneWithProgress(args []string) bool {
	cmd := exec.Command("git", args...)
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
			return i + 1, data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	})

	objRegex := regexp.MustCompile(`Receiving objects:\s+(\d+)%`)
	speedRegex := regexp.MustCompile(`\|\s*([0-9\.]+\s*[KMGT]?i?B/s)`)
	deltaRegex := regexp.MustCompile(`Resolving deltas:\s+(\d+)%`)
	updateRegex := regexp.MustCompile(`Updating files:\s+(\d+)%`)

	phaseStartTime := time.Now()
	lastPhase := ""

	for scanner.Scan() {
		line := scanner.Text()
		var label string
		var p int
		var speedStr string
		var etaStr string

		if m := objRegex.FindStringSubmatch(line); len(m) > 1 {
			label = "Receiving Objects"
			p, _ = strconv.Atoi(m[1])
			if sm := speedRegex.FindStringSubmatch(line); len(sm) > 1 {
				speedStr = "⚡ " + strings.TrimSpace(sm[1])
			}
		} else if m := deltaRegex.FindStringSubmatch(line); len(m) > 1 {
			label = "Resolving Deltas"
			p, _ = strconv.Atoi(m[1])
		} else if m := updateRegex.FindStringSubmatch(line); len(m) > 1 {
			label = "Updating Files"
			p, _ = strconv.Atoi(m[1])
		} else {
			continue
		}

		if label != lastPhase {
			phaseStartTime = time.Now()
			lastPhase = label
		}

		if p > 0 && p < 100 {
			elapsed := time.Since(phaseStartTime)
			totalEstimated := time.Duration(float64(elapsed) / (float64(p) / 100.0))
			remaining := totalEstimated - elapsed
			if remaining < 0 {
				remaining = 0
			}
			secs := int(remaining.Seconds())
			mins := secs / 60
			secs = secs % 60
			if mins > 0 {
				etaStr = fmt.Sprintf("󱎫 ETA: %02dm%02ds", mins, secs)
			} else {
				etaStr = fmt.Sprintf("󱎫 ETA: %02ds", secs)
			}
		} else if p >= 100 {
			etaStr = "󰄬 Done"
		}

		ui.DrawProgressBarWithStats(label, p, speedStr, etaStr)
	}

	err := cmd.Wait()
	fmt.Println()
	return err == nil
}
