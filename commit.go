package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runCommit(cmd *cobra.Command, args []string) error {
	// Check if ollama is installed
	if err := checkOllama(); err != nil {
		return err
	}

	// Check if we're in a git repository
	if err := checkGitRepo(); err != nil {
		return err
	}

	// Get staged changes
	diff, err := getStagedDiff()
	if err != nil {
		return err
	}

	if diff == "" {
		return fmt.Errorf("no staged changes found. Please stage your changes using 'git add' first")
	}

	// Get model from config or flag
	model := viper.GetString("model")
	if model == "" {
		model = "llama3.1"
	}

	// Show loading message with spinner
	done := make(chan bool)
	go showLoadingSpinner(model, done)

	// Generate commit message
	commitMessage, err := generateCommitMessage(diff, model)

	// Stop spinner
	done <- true

	if err != nil {
		return err
	}

	fmt.Println("\nGenerated commit message:\n")
	fmt.Println(commitMessage)

	return nil
}

func checkOllama() error {
	_, err := exec.LookPath("ollama")
	if err != nil {
		return fmt.Errorf("ollama is not installed. Please install Ollama")
	}
	return nil
}

func checkGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("this is not a Git repository. Please run this command inside a Git repository")
	}
	return nil
}

func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged diff: %v", err)
	}
	return string(output), nil
}

func generateCommitMessage(diff, model string) (string, error) {
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

func showLoadingSpinner(model string, done chan bool) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0

	for {
		select {
		case <-done:
			// Clear the entire line and move to next line
			fmt.Print("\r\033[K")
			return
		default:
			fmt.Printf("\r🤔 %s Thinking with %s model...", spinner[i], model)
			time.Sleep(100 * time.Millisecond)
			i = (i + 1) % len(spinner)
		}
	}
}
