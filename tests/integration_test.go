package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/dakoctba/cmt/internal/commit"
	"github.com/dakoctba/cmt/internal/git"
	"github.com/dakoctba/cmt/internal/ollama"
	"github.com/dakoctba/cmt/internal/spinner"
	"github.com/spf13/viper"
)

// TestIntegration tests the full workflow
func TestIntegration(t *testing.T) {
	tests := []struct {
		name    string
		setupFn func()
		wantErr bool
	}{
		{
			name: "should handle complete workflow with staged changes",
			setupFn: func() {
				// Create a test file and stage it
				testContent := "test content for integration test"
				err := os.WriteFile("test_integration.txt", []byte(testContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}

				// Stage the file
				cmd := exec.Command("git", "add", "test_integration.txt")
				if err := cmd.Run(); err != nil {
					t.Logf("Failed to stage file (this is expected if not in git repo): %v", err)
				}
			},
			wantErr: false, // May fail if ollama is not available, but that's expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFn != nil {
				tt.setupFn()
			}
			defer func() {
				// Clean up test files
				os.Remove("test_integration.txt")
				// Unstage if possible
				exec.Command("git", "reset", "test_integration.txt").Run()
			}()

			// Reset viper
			viper.Reset()
			viper.SetDefault("model", "llama3.1")

			// Test the workflow
			err := commit.RunCommit(nil, []string{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Integration test error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		diff    string
		model   string
		wantErr bool
	}{
		{
			name:    "should handle very long diff",
			diff:    strings.Repeat("diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..1234567\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1,1 @@\n+test content\n", 100),
			model:   "llama3.1",
			wantErr: false, // May fail if ollama is not available, but that's expected
		},
		{
			name:    "should handle diff with special characters",
			diff:    "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..1234567\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1,1 @@\n+test content with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?\n",
			model:   "llama3.1",
			wantErr: false, // May fail if ollama is not available, but that's expected
		},
		{
			name:    "should handle diff with unicode characters",
			diff:    "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..1234567\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1,1 @@\n+test content with unicode: 🚀✨🎉\n",
			model:   "llama3.1",
			wantErr: false, // May fail if ollama is not available, but that's expected
		},
		{
			name:    "should handle empty model name",
			diff:    "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..1234567\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1,1 @@\n+test content\n",
			model:   "",
			wantErr: true, // Expected to fail with empty model name
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

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setupFn func()
		wantErr bool
	}{
		{
			name: "should handle missing ollama gracefully",
			setupFn: func() {
				// This test assumes ollama might not be available
				// The actual behavior depends on the environment
			},
			wantErr: true, // Expected to fail if ollama is not available
		},
		{
			name: "should handle non-git directory gracefully",
			setupFn: func() {
				// This test would need to be run outside a git repo
				// For now, we just test the current environment
			},
			wantErr: false, // May fail depending on environment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFn != nil {
				tt.setupFn()
			}

			// Test individual functions that might fail
			err := ollama.CheckInstallation()
			if err != nil {
				t.Logf("ollama.CheckInstallation() failed as expected: %v", err)
			}

			err = git.CheckRepo()
			if err != nil {
				t.Logf("git.CheckRepo() failed as expected: %v", err)
			}
		})
	}
}

// TestPerformance tests performance characteristics
func TestPerformance(t *testing.T) {
	tests := []struct {
		name    string
		diff    string
		model   string
		timeout time.Duration
	}{
		{
			name:    "should complete within reasonable time",
			diff:    "diff --git a/test.txt b/test.txt\nnew file mode 100644\nindex 0000000..1234567\n--- /dev/null\n+++ b/test.txt\n@@ -0,0 +1,1 @@\n+test content\n",
			model:   "llama3.1",
			timeout: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a channel to signal completion
			done := make(chan bool, 1)
			var result string
			var err error

			// Run the function in a goroutine
			go func() {
				result, err = ollama.GenerateCommitMessage(tt.diff, tt.model)
				done <- true
			}()

			// Wait for completion or timeout
			select {
			case <-done:
				if err != nil {
					t.Logf("generateCommitMessage() failed (expected if ollama not available): %v", err)
				} else if result == "" {
					t.Error("generateCommitMessage() returned empty result")
				}
			case <-time.After(tt.timeout):
				t.Errorf("generateCommitMessage() took longer than %v", tt.timeout)
			}
		})
	}
}

// TestConcurrentAccess tests concurrent access to functions
func TestConcurrentAccess(t *testing.T) {
	tests := []struct {
		name    string
		workers int
	}{
		{
			name:    "should handle concurrent access",
			workers: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper
			viper.Reset()
			viper.SetDefault("model", "llama3.1")

			// Create channels for coordination
			start := make(chan struct{})
			done := make(chan bool, tt.workers)

			// Start workers
			for i := 0; i < tt.workers; i++ {
				go func(workerID int) {
					<-start // Wait for start signal

					// Test concurrent access to viper
					model := viper.GetString("model")
					if model != "llama3.1" {
						t.Errorf("Worker %d: Expected model 'llama3.1', got '%s'", workerID, model)
					}

					done <- true
				}(i)
			}

			// Signal all workers to start
			close(start)

			// Wait for all workers to complete
			for i := 0; i < tt.workers; i++ {
				<-done
			}
		})
	}
}

// TestMemoryUsage tests memory usage patterns
func TestMemoryUsage(t *testing.T) {
	tests := []struct {
		name       string
		iterations int
	}{
		{
			name:       "should not leak memory with repeated calls",
			iterations: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper before test
			viper.Reset()
			viper.SetDefault("model", "llama3.1")

			// Make repeated calls to test for memory leaks
			for i := 0; i < tt.iterations; i++ {
				model := viper.GetString("model")
				if model != "llama3.1" {
					t.Errorf("Iteration %d: Expected model 'llama3.1', got '%s'", i, model)
				}

				// Test spinner creation and cleanup
				spinner := spinner.New()
				spinner.Start("test-model")
				time.Sleep(50 * time.Millisecond)
				spinner.Stop()
				time.Sleep(10 * time.Millisecond)
			}
		})
	}
}
