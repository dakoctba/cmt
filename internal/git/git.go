package git

import (
	"fmt"
	"os/exec"
)

// CheckRepo verifies if the current directory is a Git repository
func CheckRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("this is not a Git repository. Please run this command inside a Git repository")
	}
	return nil
}

// GetStagedDiff returns the staged changes as a string
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %v", err)
	}
	return string(output), nil
}
