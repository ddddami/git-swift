package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

type Branch struct {
	Name      string
	IsCurrent bool
}

func main() {
	listBranchCmd := exec.Command("git", "branch")
	out, err := listBranchCmd.CombinedOutput()
	branches := strings.Split(strings.TrimSpace(string(out)), "\n")

	if err != nil {
		fmt.Printf("git branch list failed: %s - %s", err, string(out))
	}

	idx, err := fuzzyfinder.Find(
		branches, func(i int) string {
			return branches[i]
		},
	)
	if err != nil {
		fmt.Printf("Error: %s", err) // fix -  out is still holding last reference even if cmd failed
	}

	branch := strings.TrimSpace(strings.TrimPrefix(branches[idx], "*"))
	checkoutCmd := exec.Command("git", "switch", branch)
	out, err = checkoutCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("git checkout failed: %s - %s", err, string(out))
		return
	}
	fmt.Printf("Switched to branch '%s'\n", branch)
}
