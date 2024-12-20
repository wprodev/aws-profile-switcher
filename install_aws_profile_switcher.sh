#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

# Function to print error messages
error() {
    echo "Error: $1" >&2
    exit 1
}

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
    error "Go is not installed. Please install Go from https://go.dev/dl/ and try again."
fi

# Check if Git is installed
if ! command -v git >/dev/null 2>&1; then
    error "Git is not installed. Please install Git from https://git-scm.com/downloads and try again."
fi

# Configuration
REPO_URL="https://github.com/wprodev/aws-profile-switcher.git"  # Replace with your actual repo
BINARY_NAME="aws-profile-switcher"
SYMLINK_NAME="aps"
INSTALL_DIR="/usr/local/bin"
LOCAL_INSTALL_DIR="$HOME/.local/bin"

# Determine the installation directory
# Prefer /usr/local/bin, fallback to $HOME/.local/bin
if [ -w "$INSTALL_DIR" ]; then
    TARGET_DIR="$INSTALL_DIR"
else
    # Create the local install directory if it doesn't exist
    mkdir -p "$LOCAL_INSTALL_DIR"
    TARGET_DIR="$LOCAL_INSTALL_DIR"
    echo "Using local installation directory: $TARGET_DIR"
fi

# Temporary build directory
BUILD_DIR="$(mktemp -d)"
cleanup() {
    rm -rf "$BUILD_DIR"
}
trap cleanup EXIT

# Clone the repository
echo "Cloning repository from $REPO_URL..."
git clone "$REPO_URL" "$BUILD_DIR" >/dev/null 2>&1 || error "Failed to clone repository."
echo "Repository cloned."

# Build the binary
echo "Building the binary..."
cd "$BUILD_DIR"
go build -o "$BINARY_NAME" >/dev/null 2>&1 || error "Failed to build the binary."
echo "Binary built successfully."

# Function to prompt user for overwriting existing files
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

# Install the binary
TARGET_BINARY="$TARGET_DIR/$BINARY_NAME"
if prompt_overwrite "$TARGET_BINARY"; then
    # Use sudo if necessary
    if [ "$TARGET_DIR" = "$INSTALL_DIR" ] && [ ! -w "$INSTALL_DIR" ]; then
        echo "Installing to $INSTALL_DIR requires sudo permissions."
        sudo install -m 755 "$BINARY_NAME" "$TARGET_BINARY" || error "Failed to install the binary."
    else
        install -m 755 "$BINARY_NAME" "$TARGET_BINARY" || error "Failed to install the binary."
    fi
    echo "Installed '$BINARY_NAME' to '$TARGET_DIR'."
fi

# Create the symlink 'aps'
TARGET_SYMLINK="$TARGET_DIR/$SYMLINK_NAME"
if prompt_overwrite "$TARGET_SYMLINK"; then
    ln -sf "$BINARY_NAME" "$TARGET_SYMLINK" || error "Failed to create symlink '$SYMLINK_NAME'."
    echo "Created symlink '$SYMLINK_NAME' pointing to '$BINARY_NAME'."
fi

# Verify installation directory is in PATH
if ! echo "$PATH" | grep -q "$TARGET_DIR"; then
    echo "Warning: '$TARGET_DIR' is not in your PATH."
    echo "You might need to add it to your PATH to run '$BINARY_NAME' or '$SYMLINK_NAME' directly."
    echo "For example, add the following line to your ~/.bashrc or ~/.zshrc:"
    echo "  export PATH=\"\$PATH:$TARGET_DIR\""
fi

echo "Installation complete. You can run '$BINARY_NAME' or '$SYMLINK_NAME' now."
