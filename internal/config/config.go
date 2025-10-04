package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AI    AIConfig    `yaml:"ai"`
	Theme ThemeConfig `yaml:"theme"`
	Editor EditorConfig `yaml:"editor"`
}

type AIConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Provider string `yaml:"provider"` // openai, ollama, openrouter
	APIKey   string `yaml:"api_key"`
	BaseURL  string `yaml:"base_url"`
	Model    string `yaml:"model"`
}

type ThemeConfig struct {
	Current string `yaml:"current"` // dark, light, monokai, solarized
}

type EditorConfig struct {
	TabSize      int  `yaml:"tab_size"`
	ShowLineNums bool `yaml:"show_line_numbers"`
	AutoIndent   bool `yaml:"auto_indent"`
}

var DefaultConfig = Config{
	AI: AIConfig{
		Enabled:  false,
		Provider: "ollama",
		BaseURL:  "http://localhost:11434",
		Model:    "llama3.2",
	},
	Theme: ThemeConfig{
		Current: "dark",
	},
	Editor: EditorConfig{
		TabSize:      4,
		ShowLineNums: true,
		AutoIndent:   true,
	},
}

func Load() (*Config, error) {
	configPath := getConfigPath()

	// Create default config if not exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := Save(&DefaultConfig); err != nil {
			return &DefaultConfig, nil // Return default if can't save
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return &DefaultConfig, nil
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return &DefaultConfig, err
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	configPath := getConfigPath()
	configDir := filepath.Dir(configPath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".finpup.yaml"
	}
	return filepath.Join(home, ".finpup.yaml")
}
