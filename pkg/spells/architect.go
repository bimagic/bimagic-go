package spells

import (
	"os"
	"strings"

	"bimagic-go/pkg/ui"
)

func SummonGitignore() {
	if !ui.HasCmd("curl") {
		ui.PrintError("Error: curl is not installed. Required to summon the Architect.")
		return
	}

	ui.PrintStatus("📜 Summoning the Architect...")

	if _, err := os.Stat(".gitignore"); !os.IsNotExist(err) {
		ui.PrintWarning("A .gitignore file already exists in this directory.")
		if !ui.GumConfirm("Do you want to overwrite it?") {
			ui.PrintStatus("Operation cancelled.")
			return
		}
	}

	templates := []string{
		"Actionscript", "Ada", "Android", "Angular", "AppEngine", "ArchLinuxPackages", "Autotools",
		"C++", "C", "CMake", "CUDA", "CakePHP", "ChefCookbook", "Clojure", "CodeIgniter", "Composer",
		"Dart", "Delphi", "Dotnet", "Drupal", "Elixir", "Elm", "Erlang", "Flutter", "Fortran",
		"Go", "Godot", "Gradle", "Grails", "Haskell", "Haxe", "Java", "Jekyll", "Joomla", "Julia",
		"Kotlin", "Laravel", "Lua", "Magento", "Maven", "Nextjs", "Nim", "Nix", "Node", "Objective-C",
		"Opa", "Perl", "Phalcon", "PlayFramework", "Prestashop", "Processing", "Python", "Qt",
		"R", "ROS", "Rails", "Ruby", "Rust", "Scala", "Scheme", "Smalltalk", "Swift", "Symfony",
		"Terraform", "TeX", "Unity", "UnrealEngine", "VisualStudio", "WordPress", "Zig",
	}

	template := ui.GumFilterStdin(strings.Join(templates, "\n"), "Search for a blueprint (e.g., Python, Node, Rust)...", false)

	if template == "" {
		ui.PrintStatus("Cancelled.")
		return
	}

	ui.PrintStatus("Drawing the magic circle for " + template + "...")
	url := "https://raw.githubusercontent.com/github/gitignore/main/" + template + ".gitignore"

	ui.PrintCommand(`curl -sL "` + url + `" -o .gitignore`)
	if ui.GumSpin("Fetching template...", "curl", "-sL", url, "-o", ".gitignore") {
		b, _ := os.ReadFile(".gitignore")
		if strings.Contains(string(b), "404: Not Found") {
			ui.PrintError("Failed to summon template: 404 Not Found at " + url)
			os.Remove(".gitignore")
			return
		}
		ui.PrintStatus("✨ .gitignore for " + template + " created successfully!")
	} else {
		ui.PrintError("Failed to summon template. Check your internet connection.")
	}
}
