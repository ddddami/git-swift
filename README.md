# gitch

> Lightning-fast branch switching with fuzzy search

A minimal CLI tool that lets you quickly find and switch between Git branches using a clean, interactive interface.

![gitch demo](https://github.com/ddddami/gitch/raw/main/demo.gif)

## Features

- ⚡ **Fast branch switching** - Find branches faster than you can type their full name
- ⌨️ **Keyboard-driven** - Navigate without reaching for the mouse
- 🔢 **Numbered selections** - Easy visual reference
- 🔖 **Uses native `git switch`** under the hood, supports all its parameters

## Installation

### Using npm

```bash
npm install -g gitch
```

### Using Homebrew

```bash
brew tap ddddami/gitch
brew install gitch
```

### From source

```bash
go install github.com/ddddami/swift-git@latest
```

## Usage

### Interactive mode
Simply run:

```bash
gitch
```

This opens an interactive UI where you can:
- Type to filter branches
- Use ↑/↓ arrows to navigate
- Press Enter to switch to the selected branch
- Press Esc to quit

### Direct mode
If you know part of the branch name:

```bash
gitch branch-name
```

This will switch directly to the branch if an exact match is found.

## Why gitch?

- **Minimal UI** - Just the information you need, nothing more
- **Lightweight** - Fast startup time, small memory footprint
- **Zero configuration** - Works out of the box

## License

MIT
