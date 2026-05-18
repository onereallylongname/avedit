// Package config handles user configuration for avedit.
// Configuration is loaded from ~/.config/avedit/config.json.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds user preferences.
type Config struct {
	Theme     string `json:"theme"`      // Default theme name
	ThemesDir string `json:"themes_dir"` // Custom themes directory (optional)
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Theme: "dark",
	}
}

// Load reads the config from the standard path (~/.config/avedit/config.json).
// Returns default config if the file doesn't exist.
func Load() Config {
	cfg := DefaultConfig()
	path := FilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, &cfg)
	return cfg
}

// Dir returns the config directory path.
func Dir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "avedit")
}

// FilePath returns the full path to the config file.
func FilePath() string {
	return filepath.Join(Dir(), "config.json")
}

// ThemesDirFromConfig returns the themes directory, preferring config override.
func (c Config) ThemesDirResolved() string {
	if c.ThemesDir != "" {
		return c.ThemesDir
	}
	return ""
}

// Save writes the config to the standard path (~/.config/avedit/config.json).
func (c Config) Save() error {
	dir := Dir()
	if dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePath(), data, 0644)
}

// SetTheme updates the theme name in the config and persists it.
func SetTheme(name string) {
	cfg := Load()
	cfg.Theme = name
	_ = cfg.Save()
}
