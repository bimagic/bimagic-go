package config

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const Version = "v2.0.0"

var Theme = map[string]string{
	"BIMAGIC_PRIMARY":   "212",
	"BIMAGIC_SECONDARY": "51",
	"BIMAGIC_SUCCESS":   "46",
	"BIMAGIC_ERROR":     "196",
	"BIMAGIC_WARNING":   "214",
	"BIMAGIC_INFO":      "39",
	"BIMAGIC_MUTED":     "240",
	"BANNER_COLOR_1":    "51",
	"BANNER_COLOR_2":    "45",
	"BANNER_COLOR_3":    "39",
	"BANNER_COLOR_4":    "99",
	"BANNER_COLOR_5":    "135",
}

func LoadTheme(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.Trim(strings.TrimSpace(parts[1]), `"' `)
			if key != "" && val != "" {
				Theme[key] = val
			}
		}
	}
}

func GetAnsiEsc(color string) string {
	color = strings.TrimSpace(color)
	if strings.HasPrefix(color, "#") && len(color) == 7 {
		r, _ := strconv.ParseInt(color[1:3], 16, 64)
		g, _ := strconv.ParseInt(color[3:5], 16, 64)
		b, _ := strconv.ParseInt(color[5:7], 16, 64)
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}
	return fmt.Sprintf("\033[38;5;%sm", color)
}

func PlaySound(soundType string) {
	switch soundType {
	case "success":
		fmt.Print("\a")
		for i := 0; i < 2; i++ {
			fmt.Print("\a")
			time.Sleep(100 * time.Millisecond)
		}
	case "error":
		for i := 0; i < 3; i++ {
			fmt.Print("\a")
			time.Sleep(50 * time.Millisecond)
		}
	case "warning":
		fmt.Print("\a")
	case "magic":
		for i := 0; i < 3; i++ {
			fmt.Print("\a")
			time.Sleep(200 * time.Millisecond)
		}
	case "progress":
		fmt.Print("\a")
	}
}

func ShowWelcomeBanner(versionFile, configDir string, isFirstTime bool) {
	fmt.Print("\033[H\033[2J")
	fmt.Printf("%s▗▖   ▄ ▄▄▄▄   ▗▄▖  ▗▄▄▖▄  ▗▄▄▖\033[0m\n", GetAnsiEsc(Theme["BANNER_COLOR_1"]))
	fmt.Printf("%s▐▌   ▄ █ █ █ ▐▌ ▐▌▐▌   ▄ ▐▌   \033[0m\n", GetAnsiEsc(Theme["BANNER_COLOR_2"]))
	fmt.Printf("%s▐▛▀▚▖█ █   █ ▐▛▀▜▌▐▌▝▜▌█ ▐▌   \033[0m\n", GetAnsiEsc(Theme["BANNER_COLOR_3"]))
	fmt.Printf("%s▐▙▄▞▘█       ▐▌ ▐▌▝▚▄▞▘█ ▝▚▄▄▖\033[0m\n", GetAnsiEsc(Theme["BANNER_COLOR_4"]))
	fmt.Printf("%s                              \033[0m\n", GetAnsiEsc(Theme["BANNER_COLOR_5"]))

	fmt.Println()
	execStyle(fmt.Sprintf("--foreground=%s", Theme["BIMAGIC_PRIMARY"]), "󱝁 Welcome to Bimagic Git Wizard "+Version+" 󱝁")
	fmt.Println()

	if isFirstTime {
		execStyle(fmt.Sprintf("--foreground=%s", Theme["BIMAGIC_SUCCESS"]), "It looks like this is your first time using Bimagic! Let's cast some spells.")
	} else {
		execStyle(fmt.Sprintf("--foreground=%s", Theme["BIMAGIC_SUCCESS"]), "Bimagic has been updated to "+Version+"! Enjoy the new magic.")
	}
	fmt.Println()

	os.MkdirAll(configDir, 0o755)
	os.WriteFile(versionFile, []byte(Version+"\n"), 0o644)

	execStyle(fmt.Sprintf("--foreground=%s", Theme["BIMAGIC_MUTED"]), "Press Enter to open the spellbook...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func execStyle(colorArg, text string) {
	// Simple inline wrapper to avoid circular dependency with ui package
	fmt.Printf("%s%s\033[0m\n", GetAnsiEsc(Theme["BIMAGIC_PRIMARY"]), text)
}
