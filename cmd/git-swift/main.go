package main

import (
	"fmt"
	"os"

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

	var initialQuery string

	if len(os.Args) > 1 {
		branchName := os.Args[1]
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
