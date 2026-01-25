# git swift

> Lightning-fast branch switching with fuzzy search

A minimal CLI tool that lets you quickly find and switch between Git branches using a clean, interactive interface.

<p align="center">
  <img src="https://github.com/user-attachments/assets/14c5b6ee-d9d9-4b11-97fe-dfac7d1a1956" alt="git swift demo" width="800px" />
</p>

## Features

- ⚡ **Fast branch switching** - Find branches faster than you can type their full name
- ⌨️ **Keyboard-driven** - Navigate and delete branches without reaching for the mouse
- 🔢 **Smart shortcuts** - Quick-switch to recent branches (`Alt+1..9`)
- 🛠️ **Git-compatible** - Supports standard `git switch` flags (create, detach, etc.)

## Installation

### Using npm (Recommended)

```bash
npm install -g git-swift
```

<!-- ### Using Homebrew -->
<!---->
<!-- ```bash -->
<!-- brew tap ddddami/git-swift -->
<!-- brew install git-swift -->
<!-- ``` -->

### From source

```bash
go install github.com/ddddami/git-swift@latest
```

## Usage

### Interactive mode

Simply run:

```bash
git swift
```

This opens an interactive UI where you can:

| Key | Action |
| --- | --- |
| `Type` | Filter branches |
| `↑` / `↓` | Navigate list |
| `Enter` | Switch to branch |
| `Alt` / `Option` + `1-9` | Quick switch to branch # |
| `Delete` / `Ctrl`+`X` | Delete selected branch |
| `Esc` | Quit |

### Direct mode

If you know part of the branch name:

```bash
git swift branch-name
```

This will switch directly to the branch if an exact match is found.

### Pass-through mode

Pass standard `git switch` flags directly:

```bash
git swift -c new-feature    # Create & switch
git swift --detach commit   # Detach HEAD
```

## Why git swift?

- **Minimal UI** - Just the information you need, nothing more
- **Lightweight** - Fast startup time, small memory footprint
- **Zero configuration** - Works out of the box

## License

MIT
