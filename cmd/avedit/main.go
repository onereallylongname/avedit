// avedit — Avro schema editor TUI
package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/onereallylongname/avedit/internal/config"
	avio "github.com/onereallylongname/avedit/internal/io"
	"github.com/onereallylongname/avedit/internal/model"
)

// Version is set at build time via ldflags.
var Version = "0.1.0"

func main() {
	arg := ""
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	switch arg {
	case "--version", "-v":
		fmt.Printf("avedit v%s\n", Version)
		os.Exit(0)
	case "--help", "-h":
		printUsage()
		os.Exit(0)
	}

	cfg := config.Load()

	// Determine target: file, directory, or cwd
	target := "."
	if arg != "" {
		target = arg
	}
	absTarget, _ := filepath.Abs(target)

	info, err := os.Stat(absTarget)
	if err != nil && arg != "" {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var app tea.Model
	if err == nil && info.IsDir() {
		// Open in directory mode with explorer visible, no file loaded
		app = model.NewAppExplorer(absTarget, cfg)
	} else {
		// Open a specific file
		proj, _, loadErr := avio.LoadAvroFromFile(absTarget)
		if loadErr != nil {
			fmt.Fprintf(os.Stderr, "Error loading schema: %v\n", loadErr)
			os.Exit(1)
		}
		app = model.NewApp(proj, absTarget, cfg)
	}

	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
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
	fmt.Println("  -h, --help      Show this help message")
	fmt.Println("  -v, --version   Show version information")
	fmt.Println()
	fmt.Print(model.KeybindingsText())
}
