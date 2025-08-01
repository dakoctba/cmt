package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cmt" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cmt")
	}

	// Set defaults
	viper.SetDefault("model", "llama3.1")

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create it with defaults
			createDefaultConfig()
		} else {
			// Config file was found but another error was produced
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			os.Exit(1)
		}
	}
}

func createDefaultConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(home, ".cmt.yaml")

	// Create default config
	viper.Set("model", "llama3.1")

	// Write config file
	if err := viper.WriteConfigAs(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created default config file: %s\n", configPath)
}
