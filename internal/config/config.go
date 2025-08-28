package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	// Email settings
	Email EmailConfig `yaml:"email"`

	// UI settings
	UI UIConfig `yaml:"ui"`

	// External tools configuration
	ExternalTools ExternalToolsConfig `yaml:"external_tools"`
}

// EmailConfig contains email-related configuration
type EmailConfig struct {
	// Maildir path (default: ~/Mail)
	Maildir string `yaml:"maildir"`

	// Default account to use
	DefaultAccount string `yaml:"default_account"`

	// Auto-sync interval in seconds (0 to disable)
	AutoSyncInterval int `yaml:"auto_sync_interval"`
}

// UIConfig contains UI-related configuration
type UIConfig struct {
	// Theme settings
	Theme ThemeConfig `yaml:"theme"`

	// Keybindings (neovim-style, non-remappable)
	Keybindings KeybindingsConfig `yaml:"keybindings"`
}

// ThemeConfig contains theme-related settings
type ThemeConfig struct {
	// Color scheme (light, dark, auto)
	ColorScheme string `yaml:"color_scheme"`

	// Show unread indicators
	ShowUnreadIndicators bool `yaml:"show_unread_indicators"`

	// Show sync status
	ShowSyncStatus bool `yaml:"show_sync_status"`
}

// KeybindingsConfig contains keybinding settings
type KeybindingsConfig struct {
	// Leader key (default: space)
	Leader string `yaml:"leader"`
}

// ExternalToolsConfig contains external tool paths
type ExternalToolsConfig struct {
	// Path to mbsync executable
	Mbsync string `yaml:"mbsync"`

	// Path to notmuch executable
	Notmuch string `yaml:"notmuch"`

	// Path to msmtp executable
	Msmtp string `yaml:"msmtp"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()

	return &Config{
		Email: EmailConfig{
			Maildir:          filepath.Join(homeDir, "Mail"),
			DefaultAccount:   "",
			AutoSyncInterval: 300, // 5 minutes
		},
		UI: UIConfig{
			Theme: ThemeConfig{
				ColorScheme:          "auto",
				ShowUnreadIndicators: true,
				ShowSyncStatus:       true,
			},
			Keybindings: KeybindingsConfig{
				Leader: " ",
			},
		},
		ExternalTools: ExternalToolsConfig{
			Mbsync:  "mbsync",
			Notmuch: "notmuch",
			Msmtp:   "msmtp",
		},
	}
}

// Load loads the configuration from file or returns default
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read and parse config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal and save config
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigPath returns the path to the configuration file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".config", "mel", "config.yaml"), nil
}
