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
	cmd := exec.Command("git", "branch")
	out, _ := cmd.Output()
	fmt.Println(string(out))
	branches := strings.Split(strings.TrimSpace(string(out)), "\n")

	idx, err := fuzzyfinder.Find(
		branches, func(i int) string {
			return branches[i]
		},
	)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	branch := strings.TrimSpace(branches[idx])
	fmt.Println(branch)
	cmd = exec.Command("git", "switch", branch)
	out, err = cmd.Output()
	_ = out
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("switched to %s branch\n", branch)
	}
}
