# Bimagic v2.0.0 - The Go-Powered Git Wizard 🔮

<p align="center">
  <img width="400" style="border-radius: 12px;" alt="Bimagic Logo" src="./Sample/logo.png" />
</p>

<p align="center">By Bimbok and adityapaul26</p>

A powerful, multi-threaded, Go-based Git automation tool that simplifies your GitHub workflow with an interactive magic menu system.

---

## 🚀 Overview (v2.0.0 Major Release)

Bimagic is an interactive command-line tool rewritten in **Go** for maximum performance, multi-threaded speed, and rock-solid reliability. It streamlines common and advanced Git operations, making version control accessible through a user-friendly terminal interface. It handles repository initialization (`main` branch default), committing, branching, remote operations with GitHub integration, reflog recovery, stash mastery, and interactive sparse checkouts.

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
- ⚡ **Command Transparency**: Displays the exact Git commands being executed.

---

## 🛠️ Installation

### Automated Installer (Recommended)

Run this one-line command to install Bimagic:

```bash
curl -fsSL https://raw.githubusercontent.com/bimagic/bimagic-go/main/install.sh | bash
```

Or run locally from a cloned repository:

```bash
./install.sh
```

*(Use `--system` flag to install into `/usr/local/bin` for all users, or default to `~/.local/bin` without root).*

### Quick Access (Keybinding)

You can set up a **Ctrl + B** keybinding for **Zsh**, **Bash**, and **Fish** shells to summon the Git Wizard from anywhere in your terminal instantly!

- **Zsh**: Add to `~/.zshrc`:
  ```zsh
  bimagic-widget() { wz; zle redisplay }
  zle -N bimagic-widget
  bindkey '^B' bimagic-widget
  ```
- **Bash**: Add to `~/.bashrc`:
  ```bash
  bind -x '"\C-b": wz'
  ```
- **Fish**: Add to `~/.config/fish/config.fish`:
  ```fish
  bind \cb 'wz; commandline -f repaint'
  ```

### Neovim Integration

You can use Bimagic directly inside Neovim! This integration wraps the CLI tool in a floating terminal window using `toggleterm.nvim` for a seamless workflow.

#### LazyVim / Toggleterm Setup

Create a new plugin file (e.g., `~/.config/nvim/lua/plugins/bimagic.lua`) with the following configuration. This sets up a `<leader>gm` keybinding to launch the wizard in a floating popup.

```lua
return {
  {
    "akinsho/toggleterm.nvim",
    opts = function(_, opts)
      opts.size = 20
      opts.open_mapping = [[<c-\>]]
    end,
    keys = {
      {
        "<leader>gm",
        function()
          local Terminal = require("toggleterm.terminal").Terminal
          local bimagic = Terminal:new({
            cmd = "wz", -- Uses 'wz' or 'bimagic' in PATH
            hidden = true,
            direction = "float",
            float_opts = {
              border = "curved", -- 'single', 'double', 'shadow', 'curved'
              width = 100,
              height = 25,
              title = "  Bimagic Git Wizard ",
            },
            close_on_exit = true,

            on_open = function(term)
              vim.cmd("startinsert!")
              vim.api.nvim_buf_set_keymap(term.bufnr, "n", "q", "<cmd>close<CR>", { noremap = true, silent = true })
            end,
          })
          bimagic:toggle()
        end,
        desc = "Bimagic (Git Wizard)",
      },
    },
  },
}
```

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
# Option 1: For user-local installation (no sudo required)
mkdir -p ~/.local/bin
mv bimagic ~/.local/bin/
ln -sf ~/.local/bin/bimagic ~/.local/bin/wz

# Option 2: For system-wide installation (requires sudo)
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
- **[fzf](https://github.com/junegunn/fzf)** (Optional): Fuzzy file selection with side-by-side preview windows.

---

## ⚙️ Configuration

### Setting Up GitHub Credentials

Bimagic requires your GitHub username and a personal access token for HTTPS remotes. Add these to your shell configuration file:

1. Open your shell config (`~/.bashrc` or `~/.zshrc`):
```bash
nano ~/.zshrc
```
2. Add your GitHub credentials:
```bash
export GITHUB_USER="your_github_username"
export GITHUB_TOKEN="your_github_personal_access_token"
```
3. Reload configuration:
```bash
source ~/.zshrc
```

### Creating a GitHub Personal Access Token

1. Go to GitHub → Settings → Developer settings → Personal access tokens
2. Click "Generate new token (classic)"
3. Give your token a descriptive name (e.g., "bimagic-cli")
4. Select the **repo** scope (full control of repositories)
5. Click "Generate token" and copy the token immediately.

### Theme Customization 🎨

Bimagic allows you to fully customize UI colors through `~/.config/bimagic/theme.wz`. Formats accept **ANSI numbers (0-255)** or **Hex codes (#RRGGBB)**.

#### Example `theme.wz`:

```bash
# Bimagic Theme - Arctic Neon
BIMAGIC_PRIMARY="#00FFFF"
BIMAGIC_SECONDARY="#00AFFF"
BIMAGIC_SUCCESS="#00FF87"
BIMAGIC_ERROR="#FF005F"
BIMAGIC_WARNING="#FFD700"
BIMAGIC_INFO="#00FFAF"
BIMAGIC_MUTED="243"

# Banner Gradients
BANNER_COLOR_1="51"
BANNER_COLOR_2="45"
BANNER_COLOR_3="39"
BANNER_COLOR_4="99"
BANNER_COLOR_5="135"
```

### Matugen Integration

#### Step 1: Create Matugen Template

Create `~/.config/matugen/templates/bimagic-theme.wz`:

```bash
# Bimagic Theme - Generated by Matugen
BIMAGIC_PRIMARY="{{colors.primary.default.hex}}"
BIMAGIC_SECONDARY="{{colors.secondary.default.hex}}"
BIMAGIC_SUCCESS="{{colors.tertiary.default.hex}}"
BIMAGIC_ERROR="{{colors.error.default.hex}}"
BIMAGIC_WARNING="{{colors.tertiary_container.default.hex}}"
BIMAGIC_INFO="{{colors.primary_container.default.hex}}"
BIMAGIC_MUTED="{{colors.outline.default.hex}}"

BANNER_COLOR_1="{{colors.primary.default.hex}}"
BANNER_COLOR_2="{{colors.primary_fixed.default.hex}}"
BANNER_COLOR_3="{{colors.secondary.default.hex}}"
BANNER_COLOR_4="{{colors.secondary_fixed.default.hex}}"
BANNER_COLOR_5="{{colors.tertiary.default.hex}}"
```

#### Step 2: Update Matugen Config

Add this block to `~/.config/matugen/config.toml`:

```toml
[templates.bimagic]
input_path = "~/.config/matugen/templates/bimagic-theme.wz"
output_path = "~/.config/bimagic/theme.wz"
```

---

## 💻 Usage

Simply run:
```bash
bimagic
# or use the short shortcut alias
wz
```

**Pro Tip:** Press **Ctrl + B** in your terminal to quickly summon the wizard from anywhere!

### Command Line Flags

Perform quick actions directly:

- **Clone Repository**: `wz -d "repo-url"`
- **Shallow Clone**: `wz -d "repo-url" --depth 1`
- **Interactive Clone**: `wz -d -i "repo-url"`
- **The Lazy Wizard** (Add + Commit/Amend + Push): `wz -z "commit message"`
- **Status Dashboard**: `wz -s`
- **Git Graph**: `wz -g`
- **The Time Turner** (Undo last commit): `wz -u`
- **The Architect** (Summon .gitignore): `wz -a`
- **Pull**: `wz -p`

---

## 🔮 Menu Spells Walkthrough

1. **Clone repository**: Clone repositories with standard or interactive sparse checkout modes featuring real-time progress bars.
2. **Init new repo**: Initialize a new Git repository automatically configured with `main` branch (`git init -b main`).
3. **Add files**: Multi-select file staging interface with `[ALL]` support.
4. **Commit changes**: Magic Commit (Conventional Commit builder) or Quick Commit (one-line).
5. **Push to remote**: Automated push with multi-remote handling and upstream tracking.
6. **Pull latest changes**: Fetch and merge from remote repositories.
7. **Create/switch branch**: Interactive branch creation and switching with current branch `➤` marker.
8. **Set remote**: Configure HTTPS (token auth) or SSH remotes.
9. **Show status**: Multi-threaded repository health dashboard (<60ms).
10. **Contributor Statistics**: Per-author contribution activity (lines, commits, highlights) over custom timeframes.
11. **Git graph**: Pretty colorized tree log visualization.
12. **Summon the Architect (.gitignore)**: Interactive `.gitignore` generator with 70+ templates directly from GitHub.
13. **Remove files/folders (rm)**: Intelligent removal (tracked via `git rm -rf`, untracked via `rm -rf`).
14. **Merge branches**: Merge target branches into current branch with conflict reporting.
15. **Uninitialize repo**: Remove `.git` tracking cleanly.
16. **Summon the Resurrection Stone**: Search reflog history to recover lost commits or spawn recovery branches.
17. **Revert commit(s)**: Multi-select revert with conflict detection.
18. **Stash operations**: Full stash management (Push with untracked options, Pop, List, Apply, Drop, Clear).
19. **The Scrying Glass**: Instant file browser with `fzf` side-by-side preview and `bat` syntax highlighting.
20. **Exit**: Clean exit with confirmation.

---

### Detailed Spell Deep-Dive

#### Clone Repository (Option 1)
- **Standard Clone**: Full or shallow (`--depth`) clone.
- **Interactive Clone (Sparse Checkout)**: Downloads repository structure (`--filter=blob:none --no-checkout`), presents interactive file list, and downloads only selected paths.

#### Commit Changes (Option 4)
- **Magic Commit (Builder)**: Follows Conventional Commits spec. Prompts for Type (`feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`), Scope, Summary, Body, and Breaking Change (`!`).
- **Quick Commit**: One-line prompt.

#### Contributor Statistics (Option 10)
Analyzes `git log --numstat` over 7d, 30d, 90d, 1y, or All time. Outputs themed progress bar percentages, lines changed, commit counts, and highlights Most Active / Most Productive contributors.

#### Resurrection Stone (Option 16)
Inspects hidden `git reflog` history. Search lost commits by message or hash and choose to:
1. Create a new branch at that commit (Safest).
2. Hard reset current branch to that commit (Emergency recovery).

#### Time Turner (Option 17 / `-u`)
Undo button with 3 levels:
- **Soft**: Cancels commit, leaves files **staged**.
- **Mixed**: Cancels commit, **unstages** files.
- **Hard**: **Destroys** last commit and changes.

---

## ❓ Troubleshooting & Security

1. **"Command not found" after installation**:
   - Ensure `$HOME/.local/bin` is in your PATH. Add `export PATH="$HOME/.local/bin:$PATH"` to your `~/.zshrc` or `~/.bashrc`.
2. **Permission denied errors**:
   - Make script/binary executable: `chmod +x ~/.local/bin/bimagic`
3. **GitHub authentication errors**:
   - Verify `GITHUB_USER` and `GITHUB_TOKEN` environment variables are set correctly in your shell configuration.
4. **Token Security**:
   - Protect your shell configuration file permissions (`chmod 600 ~/.zshrc`). Never commit tokens to version control.

---

## 📄 License & Disclaimer

This project is open-source under the **MIT License**.

---
**Cast your Git spells with confidence!** ✨
