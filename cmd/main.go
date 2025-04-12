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
	branches := strings.Split(string(out), "\n")

	idx, err := fuzzyfinder.Find(
		branches, func(i int) string {
			return branches[i]
		},
	)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Printf("selected: %v\n", idx)
}
