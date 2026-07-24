package spells

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"bimagic-go/pkg/config"
	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

type authorStat struct {
	lines   int
	commits int
}

func ShowContributorStats() {
	if !git.IsGitRepo() {
		ui.PrintError("Not a git repository!")
		return
	}

	timeRange := ui.GumChoose("Select time range", "", "", "Last 7 days", "Last 30 days", "Last 90 days", "Last year", "All time")
	since := ""
	switch timeRange {
	case "Last 7 days":
		since = "--since=7 days ago"
	case "Last 30 days":
		since = "--since=30 days ago"
	case "Last 90 days":
		since = "--since=3 months ago"
	case "Last year":
		since = "--since=1 year ago"
	case "All time":
		since = ""
	default:
		ui.PrintWarning("No time range selected.")
		return
	}

	ui.PrintStatus(fmt.Sprintf("Analyzing contributions (%s)...", timeRange))
	fmt.Println()

	args := []string{"log", "--pretty=format:COMMIT|%aN", "--numstat"}
	if since != "" {
		args = append(args, since)
	}

	cmd := exec.Command("git", args...)
	outBytes, err := cmd.Output()
	if err != nil || len(outBytes) == 0 {
		ui.PrintError("No contribution data found for the selected period.")
		return
	}

	stats := make(map[string]*authorStat)
	totalLines := 0
	currentAuthor := ""

	scanner := bufio.NewScanner(bytes.NewReader(outBytes))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "COMMIT|") {
			parts := strings.SplitN(line, "|", 2)
			if len(parts) == 2 {
				currentAuthor = parts[1]
				if stats[currentAuthor] == nil {
					stats[currentAuthor] = &authorStat{}
				}
				stats[currentAuthor].commits++
			}
		} else if len(line) > 0 && (line[0] >= '0' && line[0] <= '9' || line[0] == '-') {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				added, deleted := 0, 0
				if parts[0] != "-" {
					added, _ = strconv.Atoi(parts[0])
				}
				if parts[1] != "-" {
					deleted, _ = strconv.Atoi(parts[1])
				}
				lines := added + deleted
				if currentAuthor != "" {
					stats[currentAuthor].lines += lines
					totalLines += lines
				}
			}
		}
	}

	if len(stats) == 0 {
		ui.PrintError("No contribution data found for the selected period.")
		return
	}

	type authorInfo struct {
		name       string
		lines      int
		commits    int
		percentage float64
	}

	var authors []authorInfo
	for name, st := range stats {
		pct := float64(0)
		if totalLines > 0 {
			pct = (float64(st.lines) / float64(totalLines)) * 100
		}
		authors = append(authors, authorInfo{
			name:       strings.TrimSpace(name),
			lines:      st.lines,
			commits:    st.commits,
			percentage: pct,
		})
	}

	sort.Slice(authors, func(i, j int) bool {
		return authors[i].lines > authors[j].lines
	})

	purple := config.GetAnsiEsc(config.Theme["BIMAGIC_PRIMARY"])
	blue := config.GetAnsiEsc(config.Theme["BIMAGIC_SECONDARY"])
	yellow := config.GetAnsiEsc(config.Theme["BIMAGIC_WARNING"])
	cyan := config.GetAnsiEsc(config.Theme["BIMAGIC_INFO"])
	nc := "\033[0m"

	fmt.Printf("%sContribution Report (%s)%s\n", purple, timeRange, nc)
	fmt.Println(strings.Repeat("─", 45))

	mostActiveAuthor := ""
	mostCommits := 0
	mostProductiveAuthor := ""
	highestAvg := 0

	for _, a := range authors {
		bar := ui.GenerateBar(a.percentage)
		fmt.Printf("%s%-15s%s %s %s%5.1f%%%s (%s%d lines%s)\n", blue, a.name, nc, bar, yellow, a.percentage, nc, cyan, a.lines, nc)

		if a.commits > mostCommits {
			mostCommits = a.commits
			mostActiveAuthor = a.name
		}

		if a.commits > 0 {
			avgLines := a.lines / a.commits
			if avgLines > highestAvg {
				highestAvg = avgLines
				mostProductiveAuthor = a.name
			}
		}
	}

	fmt.Println()
	fmt.Printf("%sHighlights:%s\n", cyan, nc)
	if mostActiveAuthor != "" {
		fmt.Printf("%sMost Active:%s %s%s%s (%s%d commits%s)\n", blue, nc, purple, mostActiveAuthor, nc, yellow, mostCommits, nc)
	}
	if mostProductiveAuthor != "" {
		fmt.Printf("%sMost Productive:%s %s%s%s (%s%d lines/commit%s)\n", blue, nc, purple, mostProductiveAuthor, nc, yellow, highestAvg, nc)
	}
	fmt.Printf("%sTotal Contributors:%s %s%d%s\n", blue, nc, yellow, len(authors), nc)
}
