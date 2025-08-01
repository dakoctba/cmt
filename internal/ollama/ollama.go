package ollama

import (
	"fmt"
	"os/exec"
	"strings"
)

// CheckInstallation verifies if Ollama is installed
func CheckInstallation() error {
	_, err := exec.LookPath("ollama")
	if err != nil {
		return fmt.Errorf("ollama is not installed. Please install Ollama")
	}
	return nil
}

// GenerateCommitMessage generates a commit message using the specified model
func GenerateCommitMessage(diff, model string) (string, error) {
	prompt := fmt.Sprintf(`You are given a Git diff. Your task is to generate a clear and concise commit message that follows the Conventional Commits specification.

Conventional Commits summary:

A Conventional Commit consists of a structured message with a type, an optional scope, and a short description. The format is:

<type>(<optional scope>): <short description>

Common types:
	•	feat: A new feature
	•	fix: A bug fix
	•	docs: Documentation-only changes
	•	style: Code style changes (formatting, missing semicolons, etc.)
	•	refactor: Code change that neither fixes a bug nor adds a feature
	•	perf: Performance improvements
	•	test: Adding or updating tests
	•	chore: Routine tasks (build process, dependencies, etc.)

⸻

Your task:
	1.	Analyze the diff.
	2.	Create a short, meaningful commit title that clearly summarizes the change using the Conventional Commits format.
	3.	Optionally, write a description explaining what was changed and why.

Return the result as a Git commit command in the following format:

git commit -m "<title>" -m "<description>"

❗ Do not include any additional text or explanations in your response. Only return the git commit instruction.
%s`, diff)

	cmd := exec.Command("ollama", "run", model, prompt)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}
