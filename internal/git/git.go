package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ddddami/git-swift/internal/utils"
)

func GetBranches() ([]string, string, error) {
	cmd := exec.Command("git", "branch")
	output, err := cmd.Output()
	if err != nil {
		return nil, "", fmt.Errorf("error getting branches: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	branches := make([]string, 0, len(lines))
	currentBranch := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		isCurrent := strings.HasPrefix(line, "*")
		name := strings.TrimSpace(strings.TrimPrefix(line, "*"))

		if isCurrent {
			currentBranch = name
		}

		branches = append(branches, name)
	}

	return branches, currentBranch, nil
}

func SwitchBranch(branchName string) error {
	cmd := exec.Command("git", "switch", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error switching branch: %s", string(output))
	}
	return nil
}

func IsAlreadyOnBranch(output []byte, branchName string) bool {
	return strings.Contains(strings.ToLower(string(output)), "already on")
}

func TryDirectSwitch(branches []string, branchName string, currentBranch string) bool {
	if branchName == currentBranch {
		fmt.Printf("\n  ▶ Already on branch '%s'\n", branchName)
		return true
	}

	cmd := exec.Command("git", "switch", branchName)
	output, err := cmd.CombinedOutput()

	if IsAlreadyOnBranch(output, branchName) {
		fmt.Printf("\n  ▶ Already on branch '%s'\n", branchName)
		return true
	}

	if err == nil {
		fmt.Printf("\n  ▶ Switched to branch '%s'\n", branchName)
		return true
	}

	matches := []string{}
	for _, branch := range branches {
		if utils.FuzzyMatch(branch, branchName) {
			matches = append(matches, branch)
		}
	}

	if len(matches) == 1 {
		if matches[0] == currentBranch {
			fmt.Printf("\n  ▶ Already on branch '%s'\n", matches[0])
			return true
		}

		cmd := exec.Command("git", "switch", matches[0])
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n  Error switching branch: %s\n", string(output))
			return false
		}

		fmt.Printf("\n  ▶ Fuzzy match found; Switched to branch '%s'\n", matches[0])
		return true
	}

	return false
}

func DeleteBranch(branchName string) error {
	cmd := exec.Command("git", "branch", "--delete", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := strings.TrimSpace(string(output))
		if strings.Contains(outputStr, "not fully merged") ||
			strings.Contains(outputStr, "not merged") {
			return &UnmergedBranchError{
				BranchName: branchName,
				Message:    outputStr,
			}
		}
		return fmt.Errorf("%s", outputStr)
	}
	return nil
}

func ForceDeleteBranch(branchName string) error {
	cmd := exec.Command("git", "branch", "--delete", "--force", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(output)))
	}
	return nil
}

type UnmergedBranchError struct {
	BranchName string
	Message    string
}

func (e *UnmergedBranchError) Error() string {
	return e.Message
}
