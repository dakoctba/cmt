package main

import (
	"fmt"
	"os"

	"github.com/dakoctba/cmt/internal/commit"
	"github.com/dakoctba/cmt/internal/config"
	"github.com/spf13/cobra"
)

var (
	// Build-time variables (injected via ldflags)
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"

	// Command flags
	cfgFile string
	model   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cmt",
		Short: "Generate conventional commit messages using AI",
		Long: `cmt is a tool that generates conventional commit messages from staged Git changes using AI models.

It analyzes your staged changes and generates a commit message following the Conventional Commits specification.`,
		Version:       version,
		RunE:          commit.RunCommit,
		SilenceUsage:  true, // Don't show usage on error
		SilenceErrors: true, // Don't show error messages automatically
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmt)")
	rootCmd.PersistentFlags().StringVar(&model, "model", "", "specify the model to use")

	// Initialize config
	config.InitConfig(cfgFile, model)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
