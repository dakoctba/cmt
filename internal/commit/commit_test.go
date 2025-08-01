package commit

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/dakoctba/cmt/internal/git"
	"github.com/dakoctba/cmt/internal/ollama"
	"github.com/dakoctba/cmt/internal/spinner"
)

func TestCheckOllama(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should check if ollama is installed",
			wantErr: false, // Assuming ollama is installed in test environment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ollama.CheckInstallation()
			if (err != nil) != tt.wantErr {
				t.Errorf("ollama.CheckInstallation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckGitRepo(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should check if we're in a git repository",
			wantErr: false, // Assuming we're in a git repo during tests
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := git.CheckRepo()
			if (err != nil) != tt.wantErr {
				t.Errorf("git.CheckRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetStagedDiff(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should get staged diff",
			wantErr: false, // May fail if no staged changes, but that's expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, err := git.GetStagedDiff()
			if (err != nil) != tt.wantErr {
				t.Errorf("git.GetStagedDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// diff can be empty if no staged changes, which is valid
			_ = diff
		})
	}
}

func TestGenerateCommitMessage(t *testing.T) {
	tests := []struct {
		name    string
		diff    string
		model   string
		wantErr bool
	}{
		{
			name:    "should generate commit message with valid diff",
			diff:    "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..1234567\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1,1 @@\n+test content\n",
			model:   "llama3.1",
			wantErr: false, // May fail if ollama is not available, but that's expected
		},
		{
			name:    "should handle empty diff",
			diff:    "",
			model:   "llama3.1",
			wantErr: false, // May fail if ollama is not available, but that's expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := ollama.GenerateCommitMessage(tt.diff, tt.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("ollama.GenerateCommitMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && message == "" {
				t.Error("ollama.GenerateCommitMessage() returned empty message")
			}
		})
	}
}

func TestShowLoadingSpinner(t *testing.T) {
	tests := []struct {
		name  string
		model string
	}{
		{
			name:  "should show loading spinner",
			model: "llama3.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spinner := spinner.New()

			// Start spinner
			spinner.Start(tt.model)

			// Let it run for a short time
			time.Sleep(200 * time.Millisecond)

			// Stop spinner
			spinner.Stop()

			// Give it time to clean up
			time.Sleep(50 * time.Millisecond)
		})
	}
}

// Mock functions for testing
func mockCheckOllama() error {
	return nil
}

func mockCheckGitRepo() error {
	return nil
}

func mockGetStagedDiff() (string, error) {
	return "mock diff content", nil
}

func mockGenerateCommitMessage(diff, model string) (string, error) {
	return "git commit -m \"test: add mock commit message\" -m \"This is a test commit message\"", nil
}

// Test helper function to check if command exists
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Test helper function to check if we're in a git repository
func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	return cmd.Run() == nil
}

func TestCommandExists(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want bool
	}{
		{
			name: "should find git command",
			cmd:  "git",
			want: true,
		},
		{
			name: "should not find non-existent command",
			cmd:  "nonexistentcommand12345",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := commandExists(tt.cmd)
			if got != tt.want {
				t.Errorf("commandExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsGitRepo(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "should detect git repository",
			want: true, // Assuming we're in a git repo during tests
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isGitRepo()
			if got != tt.want {
				t.Errorf("isGitRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test the prompt generation logic
func TestPromptGeneration(t *testing.T) {
	diff := "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..1234567\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1,1 @@\n+test content\n"
	model := "llama3.1"

	message, err := ollama.GenerateCommitMessage(diff, model)

	// This test may fail if ollama is not available, but that's expected
	if err != nil {
		t.Logf("generateCommitMessage() failed (expected if ollama not available): %v", err)
		return
	}

	// Check if the response contains expected elements
	if !strings.Contains(message, "git commit") {
		t.Error("generateCommitMessage() response should contain 'git commit'")
	}

	if !strings.Contains(message, "-m") {
		t.Error("generateCommitMessage() response should contain '-m' flag")
	}
}
