#!/usr/bin/env bash
# ==============================================================================
#  Bimagic Go - Uninstaller Script
#  Repository: https://github.com/bimagic/bimagic-go.git
# ==============================================================================

set -e

BINARY_NAME="bimagic"
CONFIG_DIR="$HOME/.config/bimagic"

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
    printf "${PURPLE}${BOLD}✨ Bimagic Go Uninstaller ✨${RESET}\n\n"
}

remove_binary() {
    info "Searching for installed $BINARY_NAME binary locations..."

    # Potential binary locations and shortcuts
    LOCATIONS=(
        "/usr/local/bin/$BINARY_NAME"
        "/usr/local/bin/wz"
        "/usr/bin/$BINARY_NAME"
        "/usr/bin/wz"
        "$HOME/.local/bin/$BINARY_NAME"
        "$HOME/.local/bin/wz"
        "$HOME/bin/$BINARY_NAME"
        "$HOME/bin/wz"
    )

    if command -v go >/dev/null 2>&1; then
        GOPATH_BIN="$(go env GOPATH 2>/dev/null)/bin"
        if [ -n "$GOPATH_BIN" ]; then
            LOCATIONS+=("$GOPATH_BIN/$BINARY_NAME" "$GOPATH_BIN/wz")
        fi
    fi

    FOUND=0

    for LOC in "${LOCATIONS[@]}"; do
        if [ -f "$LOC" ] || [ -L "$LOC" ]; then
            FOUND=$((FOUND + 1))
            if [ -w "$LOC" ] || [ -w "$(dirname "$LOC")" ]; then
                rm -f "$LOC"
                success "Removed $LOC"
            elif sudo -n true 2>/dev/null; then
                sudo rm -f "$LOC"
                success "Removed $LOC (via sudo)"
            else
                warn "Permission denied: Cannot remove $LOC without root privileges."
                info "Please remove manually using: ${BOLD}sudo rm -f $LOC${RESET}"
            fi
        fi
    done

    if [ "$FOUND" -eq 0 ]; then
        warn "No $BINARY_NAME binary executable was found in standard PATH locations."
    fi
}

clean_config() {
    PURGE="${1:-false}"

    if [ "$PURGE" = "true" ] || [ "${1:-}" = "--purge" ] || [ "${1:-}" = "-p" ]; then
        SHOULD_PURGE=true
    elif [ -d "$CONFIG_DIR" ] && [ -t 0 ]; then
        printf "${YELLOW}Do you also want to delete configuration and themes directory ($CONFIG_DIR)? [y/N]: ${RESET}"
        read -r CONFIRM
        case "$CONFIRM" in
            [yY][eE][sS]|[yY]) SHOULD_PURGE=true ;;
            *) SHOULD_PURGE=false ;;
        esac
    else
        SHOULD_PURGE=false
    fi

    if [ "${SHOULD_PURGE:-false}" = "true" ]; then
        if [ -d "$CONFIG_DIR" ]; then
            info "Removing configuration directory $CONFIG_DIR..."
            rm -rf "$CONFIG_DIR"
            success "Configuration directory removed."
        fi
    else
        if [ -d "$CONFIG_DIR" ]; then
            info "Configuration files in $CONFIG_DIR were preserved."
            info "(To remove them later, use: ./uninstall.sh --purge)"
        fi
    fi
}

main() {
    banner
    remove_binary
    clean_config "${1:-}"

    success "Uninstallation complete. Git Wizard has vanished in a puff of smoke! ✨"
}

main "$@"
