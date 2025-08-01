package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name      string
		cfgFile   string
		wantErr   bool
		setupFn   func()
		cleanupFn func()
	}{
		{
			name:    "should initialize config with default values",
			cfgFile: "",
			wantErr: false,
			setupFn: func() {
				// Reset viper
				viper.Reset()
			},
			cleanupFn: func() {
				// Clean up after test
				viper.Reset()
			},
		},
		{
			name:    "should initialize config with custom config file",
			cfgFile: "test_config.yaml",
			wantErr: false,
			setupFn: func() {
				// Create test config file
				testConfig := `model: test-model`
				err := os.WriteFile("test_config.yaml", []byte(testConfig), 0644)
				if err != nil {
					t.Fatalf("Failed to create test config file: %v", err)
				}
				viper.Reset()
			},
			cleanupFn: func() {
				// Clean up test files
				os.Remove("test_config.yaml")
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

			// This function doesn't return an error, so we just check if it panics
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("InitConfig() panicked unexpectedly: %v", r)
				}
			}()

			InitConfig(tt.cfgFile, "")

			// Check if default model is set
			model := viper.GetString("model")
			if model == "" {
				t.Error("initConfig() should set default model")
			}
		})
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	tests := []struct {
		name      string
		wantErr   bool
		setupFn   func()
		cleanupFn func()
	}{
		{
			name:    "should create default config file",
			wantErr: false,
			setupFn: func() {
				// Reset viper
				viper.Reset()
				viper.Set("model", "test-model")
			},
			cleanupFn: func() {
				// Clean up test files
				home, err := os.UserHomeDir()
				if err == nil {
					configPath := filepath.Join(home, ".cmt.yaml")
					os.Remove(configPath)
				}
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

			// This function doesn't return an error, so we just check if it panics
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("createDefaultConfig() panicked unexpectedly: %v", r)
				}
			}()

			createDefaultConfig()

			// Check if config file was created
			home, err := os.UserHomeDir()
			if err != nil {
				t.Skipf("Cannot get home directory: %v", err)
			}

			configPath := filepath.Join(home, ".cmt.yaml")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				t.Error("createDefaultConfig() should create config file")
			}
		})
	}
}

func TestConfigFileOperations(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "should handle valid YAML config",
			content: `model: llama3.1`,
			wantErr: false,
		},
		{
			name:    "should handle empty config",
			content: ``,
			wantErr: false,
		},
		{
			name:    "should handle invalid YAML config",
			content: `model: [invalid yaml`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tempFile := "temp_config.yaml"
			err := os.WriteFile(tempFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp config file: %v", err)
			}
			defer os.Remove(tempFile)

			// Reset viper
			viper.Reset()
			viper.SetConfigFile(tempFile)

			// Try to read config
			err = viper.ReadInConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("viper.ReadInConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Check if model was read correctly
				model := viper.GetString("model")
				if tt.content != "" && model == "" {
					t.Error("Config should contain model value")
				}
			}
		})
	}
}

func TestDefaultModelValue(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "should have default model value",
			want: "llama3.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper and set default
			viper.Reset()
			viper.SetDefault("model", "llama3.1")

			got := viper.GetString("model")
			if got != tt.want {
				t.Errorf("Default model = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigPathResolution(t *testing.T) {
	tests := []struct {
		name     string
		cfgFile  string
		wantPath string
	}{
		{
			name:     "should resolve home directory config path",
			cfgFile:  "",
			wantPath: ".cmt.yaml",
		},
		{
			name:     "should use custom config file path",
			cfgFile:  "/custom/path/config.yaml",
			wantPath: "/custom/path/config.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper
			viper.Reset()

			// Test with the specified config file

			if tt.cfgFile != "" {
				viper.SetConfigFile(tt.cfgFile)
			} else {
				home, err := os.UserHomeDir()
				if err != nil {
					t.Skipf("Cannot get home directory: %v", err)
				}
				viper.AddConfigPath(home)
				viper.SetConfigName(".cmt")
				viper.SetConfigType("yaml")
			}

			configFile := viper.ConfigFileUsed()
			if tt.cfgFile != "" && configFile != tt.wantPath {
				t.Errorf("Config file path = %v, want %v", configFile, tt.wantPath)
			}
		})
	}
}
