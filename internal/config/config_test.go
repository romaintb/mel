package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if cfg.Email.Maildir == "" {
		t.Error("Expected Maildir to be set")
	}

	if cfg.UI.Keybindings.Leader != " " {
		t.Errorf("Expected Leader to be space, got '%s'", cfg.UI.Keybindings.Leader)
	}
}

func TestConfigLoad(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}
}
