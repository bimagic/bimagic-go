# Bimagic v2.0.0 - The Go-Powered Git Wizard 🔮

<p align="center">
  <img width="400" style="border-radius: 12px;" alt="Bimagic Logo" src="./Sample/logo.png" />
</p>

<p align="center">By Bimbok and adityapaul26</p>

A lightning-fast, multi-threaded, Go-based Git automation tool that simplifies your GitHub workflow with an interactive magic menu.

---

## 🚀 Overview (v2.0.0 Major Release)

Bimagic v2.0.0 is a major breaking release rebuilt from the ground up with a modular Go package structure (`pkg/`), multi-threaded CPU parallel execution, and pervasive dynamic theme integration. It streamlines all common and advanced Git operations through a fluid terminal menu interface.

---

## Sample

<img width="1920" height="1080" alt="Interactive Menu" src="Sample/1.png" />
<p align="center">Interactive magic spell menu with custom theme styling</p>

<img width="1920" height="1080" alt="Confirmation Modal" src="Sample/2.png" />
<p align="center">Interactive confirmation dialog with dynamic button themes</p>

---

## ✨ Features & What's New in v2.0.0

- ⚡ **Multi-Threaded Speed Engine**: Multi-core CPU parallelism (`runtime.GOMAXPROCS` & `sync.WaitGroup`) rendering status dashboards in under **60ms**.
- 📦 **Modular Package Architecture**: Clean domain separation (`pkg/config`, `pkg/ui`, `pkg/git`, `pkg/spells`) for maximum stability and maintainability.
- 🪄 **Shortcut Symlink (`wz`)**: Automatically provisions `wz -> bimagic` shortcut symlinks during installation. Run with `bimagic` or `wz`!
- 🎨 **Full Theme Adaptation**: Every button, prompt, selector, input box, and modal dynamically inherits colors from `~/.config/bimagic/theme.wz`.
- 🔮 **Interactive Magic Menu**: Driven by Charm's `gum` with Nerd Font iconography.
- 🧙‍♂️ **Resilient Lazy Wizard (`-z`)**: Intelligently handles already-committed repositories by updating commit messages via `--amend` and pushing directly.
- 🔐 **Secure GitHub Auth**: Integration via personal access tokens or SSH.
- 📦 **Instant Initialization**: Setup new repositories guaranteed with `main` branch default (`git init -b main`).
- 📥 **Smart & Sparse Cloning**: Standard clone or interactive file selection (`--filter=blob:none`).
- 📊 **Real-time Progress Bar**: Visual feedback for git object downloads and delta resolution.
- 🗜️ **Shallow Clone Support**: Use `--depth` for lightweight clones.
- 🔄 **Simplified Remote Ops**: Automated push/pull and upstream tracking.
- 🌿 **Branch Mastery**: Concurrent creation, switching, and merging between branches.
- 📊 **Status Dashboard**: Live view of ahead/behind counts, branch state, and merge conflict warnings.
- 🛡️ **Safe File Management**: Intelligent file/folder removal (`git rm` vs `rm -rf`).
- 📈 **Contributor Statistics**: Contribution reports across custom time ranges (7d, 30d, 90d, 1y, all time).
- 🌐 **Git Graph Viewer**: Colorized tree log visualization of repository history.
- 📜 **The Architect**: Interactive `.gitignore` generator with 70+ industry templates.
- ⏪ **Revert Spells**: Multi-select commit revert with conflict safety.
- 🪨 **Resurrection Stone**: Recover deleted commits directly from the Git reflog.
- ⏳ **Time Turner**: Undo recent commits with Soft, Mixed, or Hard level resets.
- 🗃️ **Stash Mastery**: Complete control over Git stashes (Push, Pop, List, Apply, Drop, Clear).
- 🔍 **The Scrying Glass**: Instant file preview with `fzf` & `bat` syntax highlighting.

---

## 🛠️ Installation

### Automated Installer (Recommended)

Run the automated cross-platform installer script:

```bash
curl -fsSL https://raw.githubusercontent.com/bimagic/bimagic-go/main/install.sh | bash
```

Or run locally from a cloned repository:

```bash
./install.sh
```

*(Use `--system` flag to install into `/usr/local/bin` for all users, or default to `~/.local/bin` without root).*

### From Source (Requires Go)

```bash
go install github.com/bimagic/bimagic-go@latest
```

### Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/bimagic/bimagic-go.git
```
2. Build the binary:
```bash
go build -ldflags="-s -w" -o bimagic main.go
```
3. Move to your PATH and create `wz` symlink:
```bash
sudo mv bimagic /usr/local/bin/
sudo ln -sf /usr/local/bin/bimagic /usr/local/bin/wz
```

---

## 🗑️ Uninstallation

To remove Bimagic and its shortcut symlink:

```bash
curl -fsSL https://raw.githubusercontent.com/bimagic/bimagic-go/main/uninstall.sh | bash
```

Or run locally:

```bash
./uninstall.sh [--purge]
```

*(Use `--purge` or `-p` flag to also remove configuration files in `~/.config/bimagic`).*

---

## ⚡ Dependencies

- **[gum](https://github.com/charmbracelet/gum)** (Required): Powers interactive TUI elements.
- **Git** (Required): Version control engine.
- **[bat](https://github.com/sharkdp/bat)** (Optional): Syntax highlighting in The Scrying Glass.
- **[fzf](https://github.com/junegunn/fzf)** (Optional): Fuzzy file selection with preview windows.

---

## ⚙️ Configuration & Themes

### GitHub Credentials

Configure environment variables in your shell profile (`.zshrc`, `.bashrc`, etc.):

```bash
export GITHUB_USER="your_username"
export GITHUB_TOKEN="your_personal_access_token"
```

### Theme Customization 🎨

Customize your wizard's appearance at `~/.config/bimagic/theme.wz`:

```bash
BIMAGIC_PRIMARY="#00FFFF"
BIMAGIC_SECONDARY="#00AFFF"
BIMAGIC_SUCCESS="#00FF87"
BIMAGIC_ERROR="#FF005F"
BIMAGIC_WARNING="#FFD700"
BIMAGIC_INFO="#00FFAF"
BIMAGIC_MUTED="243"
```

---

## 💻 Usage

Cast your spells using either command:
```bash
bimagic
# or
wz
```

### Command Line Flags

Perform quick actions directly:

- **Clone**: `wz -d "url" [--depth 1] [-i]`
- **Lazy Wizard** (Add + Commit/Amend + Push): `wz -z "commit message"`
- **Status**: `wz -s`
- **Graph**: `wz -g`
- **Undo**: `wz -u`
- **Architect**: `wz -a`
- **Pull**: `wz -p`

---

## 📄 License

This project is open-source under the **MIT License**.

---
**Cast your Git spells with confidence!** ✨
