# AWS Profile Switcher (`aws-profile-switcher` / `aps`)

![AWS Profile Switcher](https://img.shields.io/badge/AWS_Profile_Switcher-%F0%9F%93%A6-blue)

A sleek, interactive terminal-based tool to manage your AWS profiles effortlessly. Switch between your AWS profiles with ease, view detailed configurations, and ensure your default profile is always up-to-date—all from your terminal.

## Quickstart

By defauft when executing any AWS CLI commands the `[default]` profile (`~/.aws/config`) with it's corresponding access keys (`~/.aws/credentials`) are taken into account.

AWS CLI allows to use specific profile in 2 ways:
- Execute any `aws` cli command with `--profile` flag, e.g `aws s3 ls --profile my-profile`
- Use env var e.g `export AWS_PROFILE=my-profile`

In both cases you need to remember your profile name and type more text in your terminal, manually set default profile values or create some shell aliases.

To avoid all that just install `aps` (aws-profile-switcher):
```sh
curl -sSL https://raw.githubusercontent.com/wprodev/aws-profile-switcher/main/install_aws_profile_switcher.sh | sh
```
Run with:
```sh
aps
# or
aws-profile-switcher
```

## 📋 Features

- **Profile Listing:** Displays all your AWS profiles defined in `~/.aws/config`.
- **Interactive Selection:** Navigate profiles using arrow keys and select with Enter.
- **Real-time Filtering:** Start typing to filter profiles dynamically.
- **Profile Details:** View detailed key-value pairs of the selected profile side-by-side.
- **Automatic Reloading:** Detects changes to your `~/.aws/config` and reloads profiles automatically.
- **Warnings & Infos:** Notifies you of any misconfigurations or important information.
- **Multi-Platform Support:** Compatible with macOS and Linux.
- **Installation Convenience:** Easy installation script to get you started quickly.
- **Flexible Execution:** Run the tool using `aws-profile-switcher` or its alias `aps`.

## 🚀 Installation

### Prerequisites

- **Go:** Ensure you have Go installed. [Download Go](https://go.dev/dl/)
- **Git:** Ensure Git is installed. [Download Git](https://git-scm.com/downloads)

### Quick Install via Shell Script

Use the provided installation script to download, build, and install `aws-profile-switcher` along with its alias `aps`.

1. **Download and Run the Install Script:**

    ```sh
    curl -sSL https://raw.githubusercontent.com/wprodev/aws-profile-switcher/main/install_aws_profile_switcher.sh | sh
    ```
    Or using `wget`:
    ```sh
    wget -qO- https://raw.githubusercontent.com/wprodev/aws-profile-switcher/main/install_aws_profile_switcher.sh | sh
    ```

    *NOTE: If `sudo` is not available and the user lacks write permissions to system directories (like `/usr/local/bin`), the script will fallback to install the binary in a user-specific directory `$HOME/.local/bin`*
2. **Script Breakdown:**

    - Checks for Dependencies: Verifies that both go and git are installed.
    - Clones the Repository: Downloads the latest version from GitHub.
    - Builds the Binary: Compiles the Go program.
    - Installs the Binary: Places aws-profile-switcher in /usr/local/bin or falls back to $HOME/.local/bin if necessary.
    - Creates Symlink: Sets up aps as an alias for quick access.
    - Handles Overwrites: Prompts before overwriting existing binaries or symlinks.
    - PATH Verification: Alerts if the installation directory is not in your PATH.

3. **Post-Installation:**

    Verify Installation:

    ```sh
    aws-profile-switcher --version
    # or
    aps --version
    ```
    Add to PATH (*if necessary*):

    If the script warns that the installation directory is not in your `PATH`, add it by appending the following line to your shell configuration file (`~/.bashrc`, `~/.zshrc`, etc.):

    `export PATH="$PATH:$HOME/.local/bin"`

    Then, reload your shell:
      ```sh
      source ~/.bashrc
      # or
      source ~/.zshrc
      ```

### 🛠️ Build from Source

If you prefer to build `aws-profile-switcher` from source, follow the steps below. This method is ideal for developers or users who want to customize the tool or contribute to its development.

#### 📋 Prerequisites

Before you begin, ensure that you have the following installed on your system:

- **Go** (version 1.16 or later): [Download Go](https://go.dev/dl/)
- **Git**: [Download Git](https://git-scm.com/downloads)

#### 📦 Clone the Repository

Start by cloning the `aws-profile-switcher` repository from GitHub to your local machine:

```sh
git clone https://github.com/yourusername/yourrepo.git
```
Replace yourusername and yourrepo with your actual GitHub username and repository name.

Navigate to the cloned directory:
```sh
cd yourrepo
```

#### 🔧 Build the Binary

Use the go build command to compile the aws-profile-switcher binary.

```sh
go build -o aws-profile-switcher
```

#### 📂 Install the Binary

To make aws-profile-switcher accessible system-wide, you can move the compiled binary to a directory that's included in your PATH, such as `/usr/local/bin` or `$HOME/.local/bin`.

**Option 1: Install to /usr/local/bin (Requires sudo)**

```sh
sudo install -m 755 aws-profile-switcher /usr/local/bin/
```

**Option 2: Install to $HOME/.local/bin**

First, ensure that the directory exists:

```sh
mkdir -p "$HOME/.local/bin"
```

Then, move the binary:

```sh
install -m 755 aws-profile-switcher "$HOME/.local/bin/"
```

*Note: If you choose to install to `$HOME/.local/bin`, make sure this directory is included in your PATH. You can add the following line to your shell configuration file (e.g., `~/.bashrc`, `~/.zshrc`):*

```sh
export PATH="$HOME/.local/bin:$PATH"
```

After adding, reload your shell configuration:
```sh
source ~/.bashrc
# or
source ~/.zshrc
```

#### 🔗 Create a Symlink (Optional)

For convenience, you can create a symbolic link (aps) that points to aws-profile-switcher, allowing you to execute the tool using either name.

Example:

```sh
ln -s "$(which aws-profile-switcher)" "$(dirname "$(which aws-profile-switcher)")/aps"
```

Ensure that the installation directory (e.g., `/usr/local/bin` or `$HOME/.local/bin`) is writable or use sudo if necessary.
#### 🌐 Verify the Installation

Confirm that aws-profile-switcher is installed correctly by checking its version:
```sh
aws-profile-switcher --version
# or using the symlink
aps --version
```

#### 🚀 Run the Tool

Launch the aws-profile-switcher interactive interface:

```sh
aws-profile-switcher
# or using the symlink
aps
```

#### 🔄 Updating the Tool

To update aws-profile-switcher to the latest version:

```sh
cd yourrepo
```

Pull the Latest Changes:
```sh
git pull origin main
```

Rebuild the Binary:
```sh
go build -o aws-profile-switcher
```

## 🎯 Usage

  Launch the AWS Profile Switcher:
    
  ```sh
  aws-profile-switcher
  # or
  aps
  ```

## 🔍 Troubleshooting

  * **Emoji Rendering Issues:**
    
    If warning or info icons don't render correctly, ensure your terminal supports UTF-8 and is using a font that includes emoji glyphs.

  * **Installation Permissions:**
    
    If installing to `/usr/local/bin `fails due to permission issues, ensure you have the necessary rights or use the local installation path.

  * **PATH Issues:**
    
    If you can't run `aws-profile-switcher` or `aps` after installation, verify that the installation directory is included in your `PATH`.

  * **Config File Not Detected:**
    
    Ensure your AWS config file is located at `~/.aws/config` and follows the correct format with `[profile ...]` sections.

## 🤝 Contributing

Contributions are welcome! Please open an issue or submit a pull request for any bugs, features, or improvements.

  1. Fork the Repository
  2. Create a Feature Branch
  3. Commit Your Changes
  4. Push to Your Fork
  5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 📚 References

  - [AWS CLI Configuration Files](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)
  - [Bubble Tea](https://github.com/charmbracelet/bubbletea)
  - [Lip Gloss](https://github.com/charmbracelet/lipgloss)
  - [fsnotify](https://github.com/fsnotify/fsnotify)