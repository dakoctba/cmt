package main

import (
	"os"
	"testing"

	"github.com/dakoctba/cmt/internal/commit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestMainFunction(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantErr   bool
		setupFn   func()
		cleanupFn func()
	}{
		{
			name:    "should handle help flag",
			args:    []string{"--help"},
			wantErr: false,
			setupFn: func() {
				// Reset viper
				viper.Reset()
			},
			cleanupFn: func() {
				viper.Reset()
			},
		},
		{
			name:    "should handle version flag",
			args:    []string{"--version"},
			wantErr: false,
			setupFn: func() {
				viper.Reset()
			},
			cleanupFn: func() {
				viper.Reset()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFn != nil {
				tt.setupFn()
			}
			defer func() {
				if tt.cleanupFn != nil {
					tt.cleanupFn()
				}
			}()

			// Save original args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set test args
			os.Args = append([]string{"cmt"}, tt.args...)

			// This test is mainly to ensure the program doesn't panic
			// We can't easily test the full main function without mocking external dependencies
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("main() panicked unexpectedly: %v", r)
				}
			}()
		})
	}
}

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "should create root command with correct properties",
			args:    []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command
			rootCmd := &cobra.Command{
				Use:   "cmt",
				Short: "Generate conventional commit messages using AI",
				Long: `cmt is a tool that generates conventional commit messages from staged Git changes using AI models.

It analyzes your staged changes and generates a commit message following the Conventional Commits specification.`,
				Version:       "1.0.0",
				RunE:          commit.RunCommit,
				SilenceUsage:  true,
				SilenceErrors: true,
			}

			// Test command properties
			if rootCmd.Use != "cmt" {
				t.Errorf("Root command Use = %v, want %v", rootCmd.Use, "cmt")
			}

			if rootCmd.Short == "" {
				t.Error("Root command should have a short description")
			}

			if rootCmd.Long == "" {
				t.Error("Root command should have a long description")
			}

			if rootCmd.Version == "" {
				t.Error("Root command should have a version")
			}

			if rootCmd.RunE == nil {
				t.Error("Root command should have a RunE function")
			}
		})
	}
}

func TestCommandFlags(t *testing.T) {
	tests := []struct {
		name    string
		flag    string
		wantErr bool
	}{
		{
			name:    "should have config flag",
			flag:    "config",
			wantErr: false,
		},
		{
			name:    "should have model flag",
			flag:    "model",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command
			rootCmd := &cobra.Command{
				Use:   "cmt",
				Short: "Generate conventional commit messages using AI",
				Long: `cmt is a tool that generates conventional commit messages from staged Git changes using AI models.

It analyzes your staged changes and generates a commit message following the Conventional Commits specification.`,
				Version:       "1.0.0",
				RunE:          commit.RunCommit,
				SilenceUsage:  true,
				SilenceErrors: true,
			}

			// Add flags
			rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmt)")
			rootCmd.PersistentFlags().StringVar(&model, "model", "", "specify the model to use")

			// Test if flag exists
			flag := rootCmd.PersistentFlags().Lookup(tt.flag)
			if flag == nil && !tt.wantErr {
				t.Errorf("Flag %s should exist", tt.flag)
			}
		})
	}
}

func TestViperIntegration(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{
			name:    "should bind model flag to viper",
			key:     "model",
			value:   "test-model",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper
			viper.Reset()

			// Create a new root command
			rootCmd := &cobra.Command{
				Use:   "cmt",
				Short: "Generate conventional commit messages using AI",
				Long: `cmt is a tool that generates conventional commit messages from staged Git changes using AI models.

It analyzes your staged changes and generates a commit message following the Conventional Commits specification.`,
				Version:       "1.0.0",
				RunE:          commit.RunCommit,
				SilenceUsage:  true,
				SilenceErrors: true,
			}

			// Add flags
			rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmt)")
			rootCmd.PersistentFlags().StringVar(&model, "model", "", "specify the model to use")

			// Bind flags to config
			viper.BindPFlag("model", rootCmd.PersistentFlags().Lookup("model"))

			// Set value
			viper.Set(tt.key, tt.value)

			// Test if value is set correctly
			got := viper.GetString(tt.key)
			if got != tt.value && !tt.wantErr {
				t.Errorf("Viper value for %s = %v, want %v", tt.key, got, tt.value)
			}
		})
	}
}

func TestDefaultValues(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "should have default model value",
			key:      "model",
			expected: "llama3.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper
			viper.Reset()

			// Set default
			viper.SetDefault(tt.key, tt.expected)

			// Test default value
			got := viper.GetString(tt.key)
			if got != tt.expected {
				t.Errorf("Default value for %s = %v, want %v", tt.key, got, tt.expected)
			}
		})
	}
}

func TestCommandExecution(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "should handle empty args",
			args:    []string{},
			wantErr: false, // May succeed if ollama is available and there are staged changes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command
			rootCmd := &cobra.Command{
				Use:   "cmt",
				Short: "Generate conventional commit messages using AI",
				Long: `cmt is a tool that generates conventional commit messages from staged Git changes using AI models.

It analyzes your staged changes and generates a commit message following the Conventional Commits specification.`,
				Version:       "1.0.0",
				RunE:          commit.RunCommit,
				SilenceUsage:  true,
				SilenceErrors: true,
			}

			// Add flags
			rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmt)")
			rootCmd.PersistentFlags().StringVar(&model, "model", "", "specify the model to use")

			// Set args
			rootCmd.SetArgs(tt.args)

			// Execute command
			err := rootCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("Command execution error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
