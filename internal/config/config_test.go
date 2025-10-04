package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg == nil {
		t.Error("Config should not be nil")
	}
}

func TestDefaultConfig(t *testing.T) {
	if DefaultConfig.Editor.TabSize != 4 {
		t.Errorf("Expected tab size 4, got %d", DefaultConfig.Editor.TabSize)
	}

	if DefaultConfig.Theme.Current != "dark" {
		t.Errorf("Expected theme 'dark', got '%s'", DefaultConfig.Theme.Current)
	}

	if DefaultConfig.AI.Provider != "ollama" {
		t.Errorf("Expected AI provider 'ollama', got '%s'", DefaultConfig.AI.Provider)
	}
}

func TestSave(t *testing.T) {
	// Note: Testing Save is limited without being able to override getConfigPath
	// In production, config is saved to ~/.finpup.yaml
	testCfg := DefaultConfig
	testCfg.Theme.Current = "light"

	// Just verify the function doesn't panic
	// In real usage, this saves to ~/.finpup.yaml
	_ = Save(&testCfg)
}
