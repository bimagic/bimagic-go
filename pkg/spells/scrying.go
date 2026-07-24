package spells

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"bimagic-go/pkg/config"
	"bimagic-go/pkg/git"
	"bimagic-go/pkg/ui"
)

func ScryingGlass() {
	ui.PrintStatus("󰈈 Summoning the Scrying Glass...")

	if ui.HasCmd("fzf") {
		previewCmd := "cat {}"
		if ui.HasCmd("bat") {
			previewCmd = "bat --color=always --style=numbers {}"
		}

		gitOut := git.GetGitOutput("ls-files", "--cached", "--others", "--exclude-standard")

		fzfCmd := exec.Command("fzf",
			"--preview", previewCmd,
			"--preview-window=right:60%",
			"--height=80%",
			"--layout=reverse",
			"--border",
			"--cycle",
			"--prompt=󰈈 Peer into: ",
			fmt.Sprintf("--color=bg+:-1,fg+:%s,hl:%s,hl+:%s,prompt:%s,pointer:%s,marker:%s,header:%s,spinner:%s,info:%s",
				config.Theme["BIMAGIC_PRIMARY"], config.Theme["BIMAGIC_SECONDARY"], config.Theme["BIMAGIC_SECONDARY"],
				config.Theme["BIMAGIC_INFO"], config.Theme["BIMAGIC_PRIMARY"], config.Theme["BIMAGIC_SUCCESS"],
				config.Theme["BIMAGIC_PRIMARY"], config.Theme["BIMAGIC_PRIMARY"], config.Theme["BIMAGIC_MUTED"]),
		)
		fzfCmd.Stdin = strings.NewReader(gitOut)
		fzfCmd.Stderr = os.Stderr
		out, _ := fzfCmd.Output()
		file := strings.TrimSpace(string(out))

		if file == "" {
			ui.PrintStatus("The glass goes dark (Cancelled).")
			return
		}

		if ui.GumConfirm("Open in full pager?") {
			if ui.HasCmd("bat") {
				batOut, _ := exec.Command("bat", "--color=always", file).Output()
				pager := exec.Command("gum", "pager")
				pager.Stdin = bytes.NewReader(batOut)
				pager.Stdout = os.Stdout
				pager.Run()
			} else {
				pager := exec.Command("gum", "pager")
				fileCont, _ := os.ReadFile(file)
				pager.Stdin = bytes.NewReader(fileCont)
				pager.Stdout = os.Stdout
				pager.Run()
			}
		}
	} else {
		gitOut := git.GetGitOutput("ls-files", "--cached", "--others", "--exclude-standard")
		file := ui.GumFilterStdin(gitOut, "Select a file to peer into...", false)

		if file == "" {
			ui.PrintStatus("The glass goes dark (Cancelled).")
			return
		}

		ui.PrintStatus("Peering into: " + file)
		if ui.HasCmd("bat") {
			batOut, _ := exec.Command("bat", "--color=always", file).Output()
			pager := exec.Command("gum", "pager")
			pager.Stdin = bytes.NewReader(batOut)
			pager.Stdout = os.Stdout
			pager.Run()
		} else {
			pager := exec.Command("gum", "pager")
			fileCont, _ := os.ReadFile(file)
			pager.Stdin = bytes.NewReader(fileCont)
			pager.Stdout = os.Stdout
			pager.Run()
		}
	}
}
