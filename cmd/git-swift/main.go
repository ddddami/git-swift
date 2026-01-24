package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ddddami/git-swift/internal/git"
	"github.com/ddddami/git-swift/internal/ui"
)

var (
	version = "dev"
	commit  = "dirty"
)

func main() {
	branches, currentBranch, err := git.GetBranches()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("git-swift %s (%s)\n", version, commit)
			os.Exit(0)
		case "--help", "-h":
			fmt.Println("git-swift - Fast branch switching with fuzzy search")
			fmt.Println("\nUsage:")
			fmt.Println("  git swift              Interactive mode")
			fmt.Println("  git swift <query>      Fuzzy search for <query>")
			fmt.Println("  git swift <flags>      Pass flags to git switch (e.g. -c new-branch)")
			os.Exit(0)
		}
	}

	// Check if any argument looks like a flag (pass everything to git switch directly)
	for _, arg := range os.Args[1:] {
		if len(arg) > 0 && arg[0] == '-' {
			cmd := exec.Command("git", append([]string{"switch"}, os.Args[1:]...)...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					os.Exit(exitErr.ExitCode())
				}
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	var initialQuery string
	if len(os.Args) > 1 {
		branchName := os.Args[1]
		// Try direct switch only if it's not a new branch creation attempt
		if git.TryDirectSwitch(branches, branchName, currentBranch) {
			os.Exit(0)
		}
		initialQuery = branchName
	}

	if err := ui.Run(branches, currentBranch, initialQuery); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
