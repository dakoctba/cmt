package commit

import (
	"fmt"

	"github.com/dakoctba/cmt/internal/config"
	"github.com/dakoctba/cmt/internal/git"
	"github.com/dakoctba/cmt/internal/ollama"
	"github.com/dakoctba/cmt/internal/spinner"
	"github.com/spf13/cobra"
)

// RunCommit is the main function for generating commit messages
func RunCommit(cmd *cobra.Command, args []string) error {
	// Check if ollama is installed
	if err := ollama.CheckInstallation(); err != nil {
		return err
	}

	// Check if we're in a git repository
	if err := git.CheckRepo(); err != nil {
		return err
	}

	// Get staged changes
	diff, err := git.GetStagedDiff()
	if err != nil {
		return err
	}

	if diff == "" {
		return fmt.Errorf("no staged changes found. Please stage your changes using 'git add' first")
	}

	// Get model from config
	model := config.GetModel()
	if model == "" {
		model = "llama3.1"
	}

	// Show loading message with spinner
	spinner := spinner.New()
	spinner.Start(model)

	// Generate commit message
	commitMessage, err := ollama.GenerateCommitMessage(diff, model)

	// Stop spinner
	spinner.Stop()

	if err != nil {
		return err
	}

	fmt.Println("\nGenerated commit message:")
	fmt.Println(commitMessage)

	return nil
}
