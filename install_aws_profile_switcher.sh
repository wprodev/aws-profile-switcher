#!/usr/bin/env bash

set -e

error() {
    echo "Error: $1" >&2
    exit 1
}

if ! command -v go >/dev/null 2>&1; then
    error "Go is not installed. Please install Go from https://go.dev/dl/ and try again."
fi

if ! command -v git >/dev/null 2>&1; then
    error "Git is not installed. Please install Git from https://git-scm.com/downloads and try again."
fi

REPO_URL="https://github.com/wprodev/aws-profile-switcher.git"
BINARY_NAME="aws-profile-switcher"
SYMLINK_NAME="aps"
INSTALL_DIR="/usr/local/bin"
LOCAL_INSTALL_DIR="$HOME/.local/bin"

# Determine the installation directory
# Prefer /usr/local/bin, fallback to $HOME/.local/bin
if [ -w "$INSTALL_DIR" ]; then
    TARGET_DIR="$INSTALL_DIR"
else
    mkdir -p "$LOCAL_INSTALL_DIR"
    TARGET_DIR="$LOCAL_INSTALL_DIR"
    echo "Using local installation directory: $TARGET_DIR"
fi

BUILD_DIR="$(mktemp -d)"
cleanup() {
    rm -rf "$BUILD_DIR"
}
trap cleanup EXIT

echo "Cloning repository from $REPO_URL..."
git clone "$REPO_URL" "$BUILD_DIR" >/dev/null 2>&1 || error "Failed to clone repository."
echo "Repository cloned."

echo "Building the binary..."
cd "$BUILD_DIR"
go build -o "$BINARY_NAME" >/dev/null 2>&1 || error "Failed to build the binary."
echo "Binary built successfully."

prompt_overwrite() {
    local file="$1"
    if [ -e "$file" ]; then
        while true; do
            read -r -p "The file '$file' already exists. Overwrite? [y/n]: " response
            case "$response" in
                [Yy]*) 
                    echo "Overwriting '$file'..."
                    return 0
                    ;;
                [Nn]*) 
                    echo "Skipping '$file'."
                    return 1
                    ;;
                *) 
                    echo "Please answer y or n."
                    ;;
            esac
        done
    else
        return 0
    fi
}

TARGET_BINARY="$TARGET_DIR/$BINARY_NAME"
if prompt_overwrite "$TARGET_BINARY"; then
    if [ "$TARGET_DIR" = "$INSTALL_DIR" ] && [ ! -w "$INSTALL_DIR" ]; then
        echo "Installing to $INSTALL_DIR requires sudo permissions."
        sudo install -m 755 "$BINARY_NAME" "$TARGET_BINARY" || error "Failed to install the binary."
    else
        install -m 755 "$BINARY_NAME" "$TARGET_BINARY" || error "Failed to install the binary."
    fi
    echo "Installed '$BINARY_NAME' to '$TARGET_DIR'."
fi

TARGET_SYMLINK="$TARGET_DIR/$SYMLINK_NAME"
if prompt_overwrite "$TARGET_SYMLINK"; then
    ln -sf "$BINARY_NAME" "$TARGET_SYMLINK" || error "Failed to create symlink '$SYMLINK_NAME'."
    echo "Created symlink '$SYMLINK_NAME' pointing to '$BINARY_NAME'."
fi

if ! echo "$PATH" | grep -q "$TARGET_DIR"; then
    echo "Warning: '$TARGET_DIR' is not in your PATH."
    echo "You might need to add it to your PATH to run '$BINARY_NAME' or '$SYMLINK_NAME' directly."
    echo "For example, add the following line to your ~/.bashrc or ~/.zshrc:"
    echo "  export PATH=\"\$PATH:$TARGET_DIR\""
fi

echo "Installation complete. You can run '$BINARY_NAME' or '$SYMLINK_NAME' now."
