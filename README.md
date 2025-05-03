# git swift

> Lightning-fast branch switching with fuzzy search

A minimal CLI tool that lets you quickly find and switch between Git branches using a clean, interactive interface.

<p align="center">
  <img src="https://github.com/ddddami/git-swift/raw/master/images/demo.gif" alt="git swift demo" width="800px" />
</p>

## Features

- ⚡ **Fast branch switching** - Find branches faster than you can type their full name
- ⌨️ **Keyboard-driven** - Navigate without reaching for the mouse
- 🔢 **Numbered selections** - Easy visual reference
- 🔖 **Uses native `git switch`** under the hood, supports all its parameters

## Installation

### Using npm

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

- Type to filter branches
- Use ↑/↓ arrows to navigate
- Press Enter to switch to the selected branch
- Press Esc to quit

### Direct mode

If you know part of the branch name:

```bash
git swift branch-name
```

This will switch directly to the branch if an exact match is found.

## Why git swift?

- **Minimal UI** - Just the information you need, nothing more
- **Lightweight** - Fast startup time, small memory footprint
- **Zero configuration** - Works out of the box

## License

MIT
