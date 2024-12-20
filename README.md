# AWS Profile Switcher (`aws-profile-switcher` / `aps`)

![AWS Profile Switcher](https://img.shields.io/badge/AWS_Profile_Switcher-%F0%9F%93%A6-blue)

A sleek, interactive terminal-based tool to manage your AWS profiles effortlessly. Switch between your AWS profiles with ease, view detailed configurations, and ensure your default profile is always up-to-date—all from your terminal.

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
