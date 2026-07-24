#!/usr/bin/env bash
# ==============================================================================
#  Bimagic Go - Automated Installer Script
#  Repository: https://github.com/bimagic/bimagic-go.git
# ==============================================================================

set -e

REPO_URL="https://github.com/bimagic/bimagic-go.git"
BINARY_NAME="bimagic"

# Styling & Colors
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    PURPLE='\033[0;35m'
    CYAN='\033[0;36m'
    BOLD='\033[1m'
    RESET='\033[0m'
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    PURPLE=''
    CYAN=''
    BOLD=''
    RESET=''
fi

info() {
    printf "${BLUE}[INFO]${RESET} %s\n" "$1"
}

success() {
    printf "${GREEN}[SUCCESS]${RESET} %s\n" "$1"
}

warn() {
    printf "${YELLOW}[WARNING]${RESET} %s\n" "$1"
}

error() {
    printf "${RED}[ERROR]${RESET} %s\n" "$1" >&2
}

banner() {
    printf "${CYAN}${BOLD}"
    cat << 'EOF'
 ▗▖   ▄ ▄▄▄▄   ▗▄▖  ▗▄▄▖▄  ▗▄▄▖
 ▐▌   ▄ █ █ █ ▐▌ ▐▌▐▌   ▄ ▐▌   
 ▐▛▀▚▖█ █   █ ▐▛▀▜▌▐▌▝▜▌█ ▐▌   
 ▐▙▄▞▘█       ▐▌ ▐▌▝▚▄▞▘█ ▝▚▄▄▖
EOF
    printf "${RESET}\n"
    printf "${PURPLE}${BOLD}✨ Bimagic Go Installer ✨${RESET}\n\n"
}

cleanup() {
    if [ -n "${TMP_DIR:-}" ] && [ -d "$TMP_DIR" ]; then
        rm -rf "$TMP_DIR"
    fi
}
trap cleanup EXIT INT TERM

detect_os_arch() {
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"

    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        armv7l|armv6l) ARCH="arm" ;;
        i386|i686) ARCH="386" ;;
        *) warn "Unrecognized architecture: $ARCH. Proceeding anyway..." ;;
    esac

    info "Detected OS: ${BOLD}$OS${RESET}, Arch: ${BOLD}$ARCH${RESET}"
}

check_dependencies() {
    info "Checking required dependencies..."

    # 1. Check Git
    if ! command -v git >/dev/null 2>&1; then
        error "Git is not installed. Please install git and re-run this script."
        exit 1
    fi

    # 2. Check Go
    if ! command -v go >/dev/null 2>&1; then
        error "Go language compiler (go) is not found."
        error "Please install Go (v1.20+ recommended) from https://go.dev/doc/install and re-run this installer."
        exit 1
    fi

    # 3. Check Gum
    if ! command -v gum >/dev/null 2>&1; then
        warn "Required dependency 'gum' is not installed."
        info "Attempting to auto-install 'gum' using 'go install'..."
        if go install github.com/charmbracelet/gum@latest; then
            GOPATH_BIN="$(go env GOPATH)/bin"
            export PATH="$GOPATH_BIN:$PATH"
            success "'gum' installed successfully via go install!"
        else
            warn "Could not auto-install 'gum'. You can install it manually:"
            warn "  macOS:   brew install gum"
            warn "  Arch:    sudo pacman -S gum"
            warn "  Debian/Ubuntu/Fedora/Alpine: see https://github.com/charmbracelet/gum#installation"
        fi
    else
        success "Dependency 'gum' found!"
    fi

    # 4. Check optional dependencies
    if ! command -v bat >/dev/null 2>&1; then
        info "Optional dependency 'bat' not found (provides enhanced syntax highlighting)."
    fi
    if ! command -v fzf >/dev/null 2>&1; then
        info "Optional dependency 'fzf' not found (provides fuzzy interactive file previews)."
    fi
}

determine_install_dir() {
    # If user specifies --system or is root, try /usr/local/bin
    if [ "${1:-}" = "--system" ] || [ "$(id -u 2>/dev/null || echo 1)" -eq 0 ]; then
        INSTALL_DIR="/usr/local/bin"
        if [ ! -w "$INSTALL_DIR" ] && command -v sudo >/dev/null 2>&1; then
            USE_SUDO=true
        fi
    elif [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    else
        INSTALL_DIR="$HOME/.local/bin"
    fi

    mkdir -p "$INSTALL_DIR" 2>/dev/null || true
    info "Installation target directory: ${BOLD}$INSTALL_DIR${RESET}"
}

build_and_install() {
    TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t 'bimagic')"
    BUILD_DIR=""

    # Check if running inside local source directory containing main.go
    if [ -f "main.go" ] && [ -f "go.mod" ]; then
        info "Building Bimagic from local source..."
        BUILD_DIR="."
    else
        info "Cloning Bimagic repository from $REPO_URL..."
        git clone --depth 1 "$REPO_URL" "$TMP_DIR/bimagic-go"
        BUILD_DIR="$TMP_DIR/bimagic-go"
    fi

    info "Compiling binary with Go..."
    (
        cd "$BUILD_DIR"
        go build -ldflags="-s -w" -o "$TMP_DIR/$BINARY_NAME" main.go
    )

    info "Installing $BINARY_NAME to $INSTALL_DIR..."
    if [ "${USE_SUDO:-false}" = "true" ] && [ ! -w "$INSTALL_DIR" ]; then
        sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
        sudo chmod 755 "$INSTALL_DIR/$BINARY_NAME"
    else
        mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
        chmod 755 "$INSTALL_DIR/$BINARY_NAME"
    fi

    success "Successfully installed $BINARY_NAME to $INSTALL_DIR/$BINARY_NAME!"

    # Create 'wz' shortcut symlink
    info "Creating 'wz' shortcut symlink..."
    if [ "${USE_SUDO:-false}" = "true" ] && [ ! -w "$INSTALL_DIR" ]; then
        sudo ln -sf "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/wz"
    else
        ln -sf "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/wz"
    fi
    success "Shortcut symlink created: ${BOLD}$INSTALL_DIR/wz -> $INSTALL_DIR/$BINARY_NAME${RESET}"
}

setup_config_and_theme() {
    CONFIG_DIR="$HOME/.config/bimagic"
    THEME_FILE="$CONFIG_DIR/theme.wz"

    if [ ! -d "$CONFIG_DIR" ]; then
        mkdir -p "$CONFIG_DIR"
    fi

    if [ ! -f "$THEME_FILE" ]; then
        info "Creating default theme configuration at $THEME_FILE..."
        cat << 'EOF' > "$THEME_FILE"
# Bimagic Custom Theme File
# Colors accept ANSI (0-255) or HEX (#RRGGBB) values.

BIMAGIC_PRIMARY="#00FFFF"
BIMAGIC_SECONDARY="#00AFFF"
BIMAGIC_SUCCESS="#00FF87"
BIMAGIC_ERROR="#FF005F"
BIMAGIC_WARNING="#FFD700"
BIMAGIC_INFO="#00FFAF"
BIMAGIC_MUTED="243"
EOF
    fi
}

check_path_variable() {
    case ":$PATH:" in
        *":$INSTALL_DIR:"*) ;;
        *)
            warn "$INSTALL_DIR is not currently in your \$PATH!"
            info "To run 'bimagic' from anywhere, add the following line to your shell profile (~/.bashrc, ~/.zshrc, or ~/.config/fish/config.fish):"
            printf "  ${BOLD}export PATH=\"\$PATH:%s\"${RESET}\n\n" "$INSTALL_DIR"
            ;;
    esac
}

main() {
    banner
    detect_os_arch
    check_dependencies
    determine_install_dir "${1:-}"
    build_and_install
    setup_config_and_theme
    check_path_variable

    success "Installation complete! Cast your Git spells by running: ${BOLD}bimagic${RESET}"
}

main "$@"
