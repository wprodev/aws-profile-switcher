#!/bin/sh

set -e

error() {
    echo "Error: $1" >&2
    exit 1
}

info() {
    echo "Info: $1"
}

prompt_yes_no() {
    PROMPT_MESSAGE="$1"
    DEFAULT_CHOICE="$2"

    if [ -t 0 ]; then
        while true; do
            printf "%s [y/n]: " "$PROMPT_MESSAGE"
            read RESPONSE
            case "$RESPONSE" in
                [Yy]* )
                    return 0
                    ;;
                [Nn]* )
                    return 1
                    ;;
                * )
                    echo "Please answer y or n."
                    ;;
            esac
        done
    else
        if [ "$DEFAULT_CHOICE" = "y" ]; then
            return 0
        else
            return 1
        fi
    fi
}

if ! command -v go >/dev/null 2>&1; then
    error "Go is not installed. Please install Go from https://go.dev/dl/ and try again."
fi

if ! command -v git >/dev/null 2>&1; then
    error "Git is not installed. Please install Git from https://git-scm.com/downloads and try again."
fi

if command -v sudo >/dev/null 2>&1; then
    SUDO_AVAILABLE=true
else
    SUDO_AVAILABLE=false
    info "'sudo' is not installed. Installation to system directories will require appropriate permissions."
fi

REPO_URL="https://github.com/wprodev/aws-profile-switcher.git"
BINARY_NAME="aws-profile-switcher"
SYMLINK_NAME="aps"
INSTALL_DIR="/usr/local/bin"
LOCAL_INSTALL_DIR="$HOME/.local/bin"

if [ -w "$INSTALL_DIR" ]; then
    TARGET_DIR="$INSTALL_DIR"
else
    if [ "$SUDO_AVAILABLE" = true ]; then
        if prompt_yes_no "The installation directory '$INSTALL_DIR' is not writable. Do you want to install '$BINARY_NAME' to '$INSTALL_DIR' using sudo?" "n"; then
            USE_SUDO=true
            TARGET_DIR="$INSTALL_DIR"
        else
            TARGET_DIR="$LOCAL_INSTALL_DIR"
            mkdir -p "$TARGET_DIR"
            info "Using local installation directory: '$TARGET_DIR'"
        fi
    else
        info "Cannot install to '$INSTALL_DIR' because 'sudo' is not available and the directory is not writable."
        info "Falling back to local installation directory: '$LOCAL_INSTALL_DIR'."
        TARGET_DIR="$LOCAL_INSTALL_DIR"
        mkdir -p "$TARGET_DIR"
    fi
fi

BUILD_DIR="$(mktemp -d)"
cleanup() {
    rm -rf "$BUILD_DIR"
}
trap cleanup EXIT

info "Cloning repository from '$REPO_URL'..."
git clone "$REPO_URL" "$BUILD_DIR" >/dev/null 2>&1 || error "Failed to clone repository."
info "Repository cloned."

info "Building the binary..."
cd "$BUILD_DIR"
go build -o "$BINARY_NAME" >/dev/null 2>&1 || error "Failed to build the binary."
info "Binary built successfully."

prompt_overwrite() {
    FILE_PATH="$1"
    DESCRIPTION="$2"

    if [ -e "$FILE_PATH" ]; then
        if [ "$FORCE_INSTALL" = true ]; then
            info "Overwriting '$FILE_PATH'..."
            return 0
        fi

        if [ -t 0 ]; then
            if prompt_yes_no "The file '$FILE_PATH' already exists. Do you want to overwrite it?" "n"; then
                return 0
            else
                info "Skipping '$DESCRIPTION'."
                return 1
            fi
        else
            info "Non-interactive mode: Skipping '$DESCRIPTION'."
            return 1
        fi
    else
        return 0
    fi
}

TARGET_BINARY="$TARGET_DIR/$BINARY_NAME"
if prompt_overwrite "$TARGET_BINARY" "the binary '$BINARY_NAME'"; then
    if [ "$TARGET_DIR" = "$INSTALL_DIR" ] && [ "$SUDO_AVAILABLE" = true ] && [ "$USE_SUDO" = true ]; then
        info "Installing '$BINARY_NAME' to '$INSTALL_DIR' with sudo..."
        sudo install -m 755 "$BINARY_NAME" "$TARGET_BINARY" || error "Failed to install the binary with sudo."
    else
        info "Installing '$BINARY_NAME' to '$TARGET_DIR'..."
        install -m 755 "$BINARY_NAME" "$TARGET_BINARY" || error "Failed to install the binary."
    fi
    info "Installed '$BINARY_NAME' to '$TARGET_DIR'."
fi

TARGET_SYMLINK="$TARGET_DIR/$SYMLINK_NAME"
if prompt_overwrite "$TARGET_SYMLINK" "the symlink '$SYMLINK_NAME'"; then
    if [ "$TARGET_DIR" = "$INSTALL_DIR" ] && [ "$SUDO_AVAILABLE" = true ] && [ "$USE_SUDO" = true ]; then
        info "Creating symlink '$SYMLINK_NAME' pointing to '$BINARY_NAME' with sudo..."
        sudo ln -sf "$BINARY_NAME" "$TARGET_SYMLINK" || error "Failed to create symlink with sudo."
    else
        info "Creating symlink '$SYMLINK_NAME' pointing to '$BINARY_NAME'..."
        ln -sf "$BINARY_NAME" "$TARGET_SYMLINK" || error "Failed to create symlink."
    fi
    info "Created symlink '$SYMLINK_NAME' pointing to '$BINARY_NAME'."
fi

PATH_INCLUDES_TARGET=false
OLD_IFS="$IFS"
IFS=:
for dir in $PATH; do
    if [ "$dir" = "$TARGET_DIR" ]; then
        PATH_INCLUDES_TARGET=true
        break
    fi
done
IFS="$OLD_IFS"

if [ "$PATH_INCLUDES_TARGET" = false ]; then
    echo "Warning: '$TARGET_DIR' is not in your PATH."
    echo "You might need to add it to your PATH to run '$BINARY_NAME' or '$SYMLINK_NAME' directly."
    echo "For example, add the following line to your ~/.bashrc or ~/.zshrc:"
    echo "  export PATH=\"\$PATH:$TARGET_DIR\""
fi

echo "Installation complete. You can run '$BINARY_NAME' or '$SYMLINK_NAME' now."
