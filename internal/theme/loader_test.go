package theme

import (
	"path/filepath"
	"runtime"
	"testing"
)

func themesDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "themes")
}

func TestLoadFile_CatppuccinMocha(t *testing.T) {
	th, err := LoadFile(filepath.Join(themesDir(), "catppuccin-mocha.json"))
	if err != nil {
		t.Fatal(err)
	}
	if th.Name != "catppuccin-mocha" {
		t.Fatalf("expected name catppuccin-mocha, got %s", th.Name)
	}
}

func TestLoadFile_TokyoNight(t *testing.T) {
	th, err := LoadFile(filepath.Join(themesDir(), "tokyo-night.json"))
	if err != nil {
		t.Fatal(err)
	}
	if th.Name != "tokyo-night" {
		t.Fatalf("expected name tokyo-night, got %s", th.Name)
	}
}

func TestLoadFile_RosePine(t *testing.T) {
	th, err := LoadFile(filepath.Join(themesDir(), "rose-pine.json"))
	if err != nil {
		t.Fatal(err)
	}
	if th.Name != "rose-pine" {
		t.Fatalf("expected name rose-pine, got %s", th.Name)
	}
}

func TestRegistry(t *testing.T) {
	reg := NewRegistry(themesDir())
	// Built-in dark always present
	if reg.Get("dark") == nil {
		t.Fatal("built-in dark theme not found")
	}
	// File-loaded themes
	for _, name := range []string{"catppuccin-mocha", "tokyo-night", "rose-pine", "light", "monokai"} {
		if reg.Get(name) == nil {
			t.Fatalf("expected theme %s to be loaded", name)
		}
	}
	names := reg.Names()
	if len(names) < 6 {
		t.Fatalf("expected at least 6 themes, got %d: %v", len(names), names)
	}
}
