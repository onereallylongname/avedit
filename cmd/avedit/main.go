// avedit — Avro schema editor TUI
package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/onereallylongname/avedit/internal/config"
	avio "github.com/onereallylongname/avedit/internal/io"
	"github.com/onereallylongname/avedit/internal/logging"
	"github.com/onereallylongname/avedit/internal/model"
)

// Version is set at build time via ldflags.
var Version = "0.1.0"

func main() {
	// Parse flags manually (consistent with existing style)
	var logLevel string
	var target string
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case arg == "--version" || arg == "-v":
			fmt.Printf("avedit v%s\n", Version)
			os.Exit(0)
		case arg == "--help" || arg == "-h":
			printUsage()
			os.Exit(0)
		case arg == "--log-level" && i+1 < len(os.Args):
			i++
			logLevel = os.Args[i]
		case strings.HasPrefix(arg, "--log-level="):
			logLevel = strings.TrimPrefix(arg, "--log-level=")
		default:
			if target == "" {
				target = arg
			}
		}
	}

	cfg := config.Load()

	// Resolve log level: CLI flag > config > default (error)
	levelStr := cfg.LogLevel
	if logLevel != "" {
		levelStr = logLevel
	}
	if levelStr == "" {
		levelStr = "error"
	}
	level := logging.ParseLevel(levelStr)

	// Initialize logging
	cleanup, err := logging.Init(logging.LogDir(), level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: logging init failed: %v\n", err)
	} else {
		defer cleanup()
	}

	// Determine target: file, directory, or cwd
	if target == "" {
		target = "."
	}
	absTarget, _ := filepath.Abs(target)

	info, statErr := os.Stat(absTarget)
	if statErr != nil && target != "." {
		slog.Error("target not found", "path", absTarget, "error", statErr)
		fmt.Fprintf(os.Stderr, "Error: %v\n", statErr)
		os.Exit(1)
	}

	var app tea.Model
	if statErr == nil && info.IsDir() {
		app = model.NewAppExplorer(absTarget, cfg)
	} else {
		proj, _, loadErr := avio.LoadAvroFromFile(absTarget)
		if loadErr != nil {
			slog.Error("schema load failed", "path", absTarget, "error", loadErr)
			fmt.Fprintf(os.Stderr, "Error loading schema: %v\n", loadErr)
			os.Exit(1)
		}
		app = model.NewApp(proj, absTarget, cfg)
	}

	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		slog.Error("TUI runtime error", "error", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("avedit - Avro schema editor TUI")
	fmt.Println()
	fmt.Printf("Usage: avedit [flags] [<schema.avsc> | <directory>]\n\n")
	fmt.Println("  No arguments    Open in current directory with file explorer")
	fmt.Println("  <schema.avsc>   Open a specific schema file")
	fmt.Println("  <directory>     Open in directory with file explorer")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -h, --help              Show this help message")
	fmt.Println("  -v, --version           Show version information")
	fmt.Println("  --log-level <level>     Set log level (debug, info, warn, error)")
	fmt.Println()
	fmt.Print(model.KeybindingsText())
}
