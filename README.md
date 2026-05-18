``
   █████████                                     ██████████     █████  ███  ███████████             ███
  ███▒▒▒▒▒███                                   ▒▒███▒▒▒▒▒█    ▒▒███  ▒▒▒  ▒█▒▒▒███▒▒▒█            ▒▒▒
 ▒███    ▒███  █████ █████ ████████   ██████     ▒███  █ ▒   ███████  ████ ▒   ▒███  ▒  █████ ████ ████
 ▒███████████ ▒▒███ ▒▒███ ▒▒███▒▒███ ███▒▒███    ▒██████    ███▒▒███ ▒▒███     ▒███    ▒▒███ ▒███ ▒▒███
 ▒███▒▒▒▒▒███  ▒███  ▒███  ▒███ ▒▒▒ ▒███ ▒███    ▒███▒▒█   ▒███ ▒███  ▒███     ▒███     ▒███ ▒███  ▒███
 ▒███    ▒███  ▒▒███ ███   ▒███     ▒███ ▒███    ▒███ ▒   █▒███ ▒███  ▒███     ▒███     ▒███ ▒███  ▒███
 █████   █████  ▒▒█████    █████    ▒▒██████     ██████████▒▒████████ █████    █████    ▒▒████████ ████
▒▒▒▒▒   ▒▒▒▒▒    ▒▒▒▒▒    ▒▒▒▒▒      ▒▒▒▒▒▒     ▒▒▒▒▒▒▒▒▒▒  ▒▒▒▒▒▒▒▒ ▒▒▒▒▒    ▒▒▒▒▒      ▒▒▒▒▒▒▒▒ ▒▒▒▒▒
```

# avedit

A terminal-based Avro schema editor inspired by lazygit. Navigate, edit, and manage `.avsc` schemas entirely from your terminal with vim-like keybindings, undo/redo, and live theme switching.

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)

## Why avedit?

Working with Avro schemas in a text editor means juggling deeply nested JSON, remembering type syntax, and manually tracking named type references. avedit gives you:

- **Visual tree** of your schema structure — instantly understand nesting and types
- **Safe editing** — undo/redo everything, type pickers prevent typos
- **Named type awareness** — rename a record and all references update automatically
- **Zero setup** — single binary, no runtime dependencies

## Features

- **Tree navigation** — expand, collapse, and browse complex schemas with depth indicators
- **Inline editing** — edit field names, types, docs, defaults, and custom attributes
- **Structural mutations** — add, delete, copy, move fields; replace types
- **Named type propagation** — renaming a type updates all references across the schema
- **Undo/Redo** — full history stack (500 levels deep) with descriptions
- **Search** — prefix filters (`n:`, `t:`, `p:`, `ns:`, `a:`) and fuzzy text matching
- **File Explorer** — browse, filter, and open `.avsc` files with recursive search
- **Commands** — vim-style `:w`, `:q`, `:wq`, `:export`, `:theme`, `:notifications`
- **Themes** — JSON-based themes with runtime switching (ships with 10 themes)
- **Categorized type picker** — types grouped by Primitives / Complex / Named / Aliases
- **Notifications** — colored flash messages with scrollable history log
- **Persistent history** — command and search history saved across sessions
- **Config** — user preferences loaded from `~/.config/avedit/config.json`

## Quick Start

```bash
# Open a schema file directly
avedit schema.avsc

# Browse a directory of schemas
avedit ./schemas/

# Open in current directory
avedit
```

Once inside:
1. Use `j`/`k` to navigate the tree, `h`/`l` to collapse/expand
2. Press `Enter` to open the details panel for the selected node
3. Press `e` or `Enter` on an attribute to edit it
4. Press `a` to add a new field, `d` to delete
5. Use `u` to undo, `Ctrl+R` to redo
6. Type `:w` to save, `:q` to quit

## Installation

### From source

```bash
go install github.com/onereallylongname/avedit/cmd/avedit@latest
```

_**Note**: Add the go installation path to you $PATH_

## Usage

```bash
avedit <schema.avsc>
avedit <directory>
avedit                  # opens explorer in CWD
avedit --help
avedit --version
```

## Keybindings

### Global

| Key       | Action              |
|-----------|---------------------|
| `q`       | Quit                |
| `?`       | Show help overlay   |
| `Tab`     | Switch panel focus  |
| `1`/`2`/`3` | Focus panel directly |
| `Ctrl+E`  | Toggle file explorer |
| `Esc`     | Cancel / exit mode  |
| `Ctrl+S`  | Save                |

### Tree (Normal Mode)

| Key       | Action                |
|-----------|-----------------------|
| `j`/`k`   | Navigate up/down      |
| `h`/`l`   | Collapse/expand node  |
| `g`/`G`   | Jump to top/bottom    |
| `Space`   | Toggle expand         |
| `Enter`   | Focus details pane    |
| `a`       | Add field             |
| `d`       | Delete node           |
| `c`       | Copy node             |
| `m`       | Move to record        |
| `u`       | Undo                  |
| `Ctrl+R`  | Redo                  |
| `/`       | Search                |

### Explorer (Normal Mode)

| Key            | Action                    |
|----------------|---------------------------|
| `j`/`k`        | Navigate up/down          |
| `h`/`l`        | Collapse/expand dir       |
| `Enter`/`Space`| Open file / toggle dir    |
| `g`/`G`        | Jump to top/bottom        |
| `Backspace`/`-`| Navigate to parent dir    |
| `.`            | Set selected dir as root  |
| `/`            | Search files (recursive)  |

### Details (Normal Mode)

| Key       | Action                    |
|-----------|---------------------------|
| `j`/`k`   | Navigate attributes       |
| `Enter`/`e` | Edit attribute          |
| `a`       | Add custom attribute      |
| `d`       | Delete attribute          |
| `R`       | Rename custom attr key    |

### Edit Mode

| Key            | Action                        |
|----------------|-------------------------------|
| `Enter`        | Confirm edit                  |
| `Esc`          | Cancel edit                   |
| `Shift+Enter`  | Newline (multiline fields)    |

### Search Mode

| Key    | Action              |
|--------|---------------------|
| `/`    | Open search         |
| `Enter`| Jump to match       |
| `n`/`N`| Next/prev match     |
| `Esc`  | Close search        |

### Search Prefixes

| Prefix | Filters by           | Example          |
|--------|----------------------|------------------|
| `n:`   | Field/record name    | `n:order`        |
| `t:`   | Type                 | `t:record`       |
| `p:`   | Parent name          | `p:Customer`     |
| `ns:`  | Namespace            | `ns:com.payment` |
| `a:`   | Alias                | `a:OrderEvent`   |
| (none) | Free text (all fields) | `timestamp`    |

### Command Mode

| Command           | Action                      |
|-------------------|-----------------------------|
| `:w [file]`       | Save (optionally to file)   |
| `:q`              | Quit                        |
| `:q!`             | Force quit (no save prompt) |
| `:wq`             | Save and quit               |
| `:export [file]`  | Export .avsc                 |
| `:theme [name]`   | Switch theme (picker if no name) |
| `:open <file>`    | Open schema file            |
| `:notifications`  | Show notification history   |

## How It Works

avedit uses a **projection model** architecture:

```
.avsc → Parse → Projection (flat node map) → Render UI
                       ↑                          |
                   Commands ←───────── User input (keys)
                       |
                       v
               Projection → Emit → .avsc
```

- The schema is parsed into a flat map of nodes linked by parent/child IDs
- Every edit creates a reversible Command (with Do/Undo closures)
- Commands are pushed to a history stack enabling unlimited undo/redo
- Export reconstructs valid Avro JSON by walking the tree

This means:
- Edits never corrupt the schema structure
- Undo works perfectly for any operation (add, delete, move, rename, type change)
- Named type renames propagate automatically to all references

For full architecture details, see [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).

## Configuration

avedit reads configuration from `~/.config/avedit/config.json`:

```json
{
  "theme": "catppuccin-mocha",
  "themes_dir": ""
}
```

| Field        | Description                                          | Default  |
|--------------|------------------------------------------------------|----------|
| `theme`      | Default theme to load on startup                     | `"dark"` |
| `themes_dir` | Custom directory for additional theme files (optional) | `""`     |

### Theme discovery order

1. `./themes/` (current working directory)
2. Next to the `avedit` executable
3. `~/.config/avedit/themes/`

## Themes

Themes are JSON files. avedit ships with: `dark`, `light`, `monokai`, `catppuccin-mocha`, `tokyo-night`, `rose-pine`, `dracula`, `gruvbox`, `one-dark-pro`, `vscode-light`.

Switch at runtime with `:theme` (opens picker) or `:theme <name>`.

### Theme file format

Create a `.json` file in your themes directory:

```json
{
  "name": "my-theme",
  "fg": "#c0caf5",
  "bg": "#1a1b26",
  "dim": "#787c99",
  "muted": "#414868",
  "primary": "#7aa2f7",
  "secondary": "#bb9af7",
  "success": "#9ece6a",
  "warning": "#e0af68",
  "error": "#f7768e",
  "badge_record": "#7aa2f7",
  "badge_field": "#9ece6a",
  "badge_enum": "#bb9af7",
  "badge_array": "#e0af68",
  "badge_map": "#73daca",
  "badge_union": "#ff9e64",
  "badge_fixed": "#2ac3de",
  "badge_named": "#bb9af7",
  "status_bg": "#16161e"
}
```

| Color           | Purpose                             |
|-----------------|-------------------------------------|
| `fg`            | Primary text foreground             |
| `bg`            | Main background                     |
| `dim`           | Subdued text (guides, primitives)   |
| `muted`         | Very subdued (borders, separators)  |
| `primary`       | Accent (focus, active items)        |
| `secondary`     | Secondary accent (headings, enums)  |
| `success`       | Positive indicators (fields, edit mode) |
| `warning`       | Warning indicators (arrays, search mode) |
| `error`         | Error indicators (flash errors)     |
| `badge_record`  | Record node type color              |
| `badge_field`   | Field node type color               |
| `badge_enum`    | Enum node type color                |
| `badge_array`   | Array node type color               |
| `badge_map`     | Map node type color                 |
| `badge_union`   | Union node type color               |
| `badge_fixed`   | Fixed node type color               |
| `badge_named`   | Named reference type color          |
| `status_bg`     | Status bar background               |

> **Tip:** All colors are optional. Omitted values fall back to sensible defaults derived from the primary/secondary accent colors.

## Setup example

```bash
# Create config directory
mkdir -p ~/.config/avedit/themes

# Write config
cat > ~/.config/avedit/config.json << 'EOF'
{
  "theme": "rose-pine"
}
EOF

# Copy or create a custom theme
cat > ~/.config/avedit/themes/custom.json << 'EOF'
{
  "name": "custom",
  "fg": "#d4d4d4",
  "bg": "#1e1e1e",
  "dim": "#808080",
  "muted": "#404040",
  "primary": "#569cd6",
  "secondary": "#c586c0",
  "success": "#6a9955",
  "warning": "#d7ba7d",
  "error": "#f44747",
  "badge_record": "#569cd6",
  "badge_field": "#6a9955",
  "badge_enum": "#c586c0",
  "badge_array": "#d7ba7d",
  "badge_map": "#4ec9b0",
  "badge_union": "#ce9178",
  "badge_fixed": "#9cdcfe",
  "status_bg": "#252526"
}
EOF
```

## Architecture

```
cmd/avedit/main.go          Entry point + CLI flags
internal/
  config/                   User configuration + persistent history
  model/                    Bubbletea TUI models (app, tree, details, explorer, picker)
  projection/               Core schema tree (nodes, build, emit, clone, stats)
  command/                  Undo/redo mutation commands (create, remove, move, copy, replace, rename)
  search/                   Query engine (prefix filters, fuzzy scoring)
  schema/                   Avro type system constants (kinds, slots, templates)
  io/                       File load/export + clipboard
  theme/                    Theme system (JSON loader, registry, defaults)
themes/                     Bundled theme JSON files
docs/                       Architecture docs + plan
```

See [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) for detailed design documentation.

## Tips & Tricks

- **Named type references** — When you change a record/enum/fixed name, all fields referencing that type update automatically
- **Type picker** — Press `/` inside the picker to filter by typing (case-insensitive)
- **Multiline editing** — For `doc` fields, press `Ctrl+Enter` to add newlines
- **Explorer search** — `/` in the explorer searches recursively through all subdirectories
- **History** — Use `↑`/`↓` in search (`/`) and command (`:`) modes to browse previous entries
- **Unsaved changes** — avedit prompts before quitting if you have unsaved edits
- **Persistent config** — Theme choice is saved to `~/.config/avedit/config.json` automatically
- **Multiple theme dirs** — Set `themes_dir` in config to load additional themes alongside built-ins

## TODO

- [ ] Packaging 
   - [ ] Homebrew (macOS/Linux) `brew install onereallylongname/tap/avedit`
   - [ ] Scoop (Windows)
```powershell
scoop bucket add avedit https://github.com/onereallylongname/scoop-bucket
scoop install avedit
```

### Pre-built binaries

Download from [Releases](https://github.com/onereallylongname/avro-schema-viz/releases).

## License

MIT
