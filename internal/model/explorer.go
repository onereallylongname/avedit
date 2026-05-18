package model

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/onereallylongname/avedit/internal/theme"
)

// FileEntry represents a file or directory in the explorer.
type FileEntry struct {
	Name     string
	Path     string
	IsDir    bool
	Depth    int
	Expanded bool
}

// ExplorerModel manages the file explorer panel.
type ExplorerModel struct {
	theme  *theme.Theme
	width  int
	height int

	rootDir string
	entries []FileEntry
	cursor  int
	offset  int
	visible bool

	// Search/filter state
	filterQuery   string
	matchIndices  []int // indices into entries that match filter
	matchIdx      int   // current match position within matchIndices
}

// NewExplorerModel creates a file explorer rooted at the given directory.
func NewExplorerModel(th *theme.Theme, rootDir string) ExplorerModel {
	// Always use absolute path for reliable file operations
	abs, err := filepath.Abs(rootDir)
	if err == nil {
		rootDir = abs
	}
	e := ExplorerModel{
		theme:   th,
		rootDir: rootDir,
		visible: false,
	}
	e.rebuild()
	return e
}

// Visible returns whether the explorer is shown.
func (e ExplorerModel) Visible() bool { return e.visible }

// RootDir returns the explorer's current root directory.
func (e ExplorerModel) RootDir() string { return e.rootDir }

// Toggle shows/hides the explorer. When showing, refreshes the file listing.
func (e *ExplorerModel) Toggle() {
	e.visible = !e.visible
	if e.visible {
		e.rebuild()
	}
}

// Refresh rebuilds the file listing from disk.
func (e *ExplorerModel) Refresh() { e.rebuild() }

// SetSize updates panel dimensions.
func (e *ExplorerModel) SetSize(w, h int) {
	e.width = w
	e.height = h
}

// SelectedEntry returns the currently selected file entry, or nil.
func (e *ExplorerModel) SelectedEntry() *FileEntry {
	if len(e.entries) == 0 || e.cursor < 0 || e.cursor >= len(e.entries) {
		return nil
	}
	return &e.entries[e.cursor]
}

// HandleKey processes navigation keys for the explorer.
func (e ExplorerModel) HandleKey(msg tea.KeyPressMsg) ExplorerModel {
	switch msg.String() {
	case "j", "down":
		if e.cursor < len(e.entries)-1 {
			e.cursor++
			e.adjustScroll()
		}
	case "k", "up":
		if e.cursor > 0 {
			e.cursor--
			e.adjustScroll()
		}
	case "g", "home":
		e.cursor = 0
		e.offset = 0
	case "G", "end":
		e.cursor = len(e.entries) - 1
		e.adjustScroll()
	case "l", "right", "space":
		e.expandSelected()
	case "h", "left":
		e.collapseSelected()
	case "backspace", "-":
		e.goUp()
	case ".":
		e.setRootToSelected()
	}
	return e
}

// setRootToSelected sets the selected directory as the new explorer root.
func (e *ExplorerModel) setRootToSelected() {
	entry := e.SelectedEntry()
	if entry == nil {
		return
	}
	var newRoot string
	if entry.IsDir && entry.Name != ".." {
		newRoot = entry.Path
	} else if entry.IsDir && entry.Name == ".." {
		newRoot = entry.Path
	} else {
		// File selected — use its parent dir
		newRoot = filepath.Dir(entry.Path)
	}
	if newRoot == e.rootDir {
		return
	}
	e.rootDir = newRoot
	e.cursor = 0
	e.offset = 0
	e.rebuild()
}

// EnterSelected opens the selected directory or returns the file path.
// Returns (filePath, true) if a file was selected, ("", false) otherwise.
func (e *ExplorerModel) EnterSelected() (string, bool) {
	entry := e.SelectedEntry()
	if entry == nil {
		return "", false
	}
	// ".." entry navigates to parent
	if entry.Name == ".." {
		e.goUp()
		return "", false
	}
	if entry.IsDir {
		e.toggleExpand(e.cursor)
		return "", false
	}
	// Only return .avsc files
	if strings.HasSuffix(strings.ToLower(entry.Name), ".avsc") {
		return entry.Path, true
	}
	return "", false
}

// goUp navigates to the parent directory.
func (e *ExplorerModel) goUp() {
	parent := filepath.Dir(e.rootDir)
	if parent == e.rootDir {
		return // already at root
	}
	e.rootDir = parent
	e.cursor = 0
	e.offset = 0
	e.rebuild()
}

func (e *ExplorerModel) expandSelected() {
	entry := e.SelectedEntry()
	if entry == nil || !entry.IsDir {
		return
	}
	if !entry.Expanded {
		e.toggleExpand(e.cursor)
	}
}

func (e *ExplorerModel) collapseSelected() {
	entry := e.SelectedEntry()
	if entry == nil {
		return
	}
	if entry.IsDir && entry.Expanded {
		e.toggleExpand(e.cursor)
		return
	}
	// If on a file or collapsed dir, jump to parent directory
	if entry.Depth > 0 {
		for i := e.cursor - 1; i >= 0; i-- {
			if e.entries[i].IsDir && e.entries[i].Depth < entry.Depth {
				e.cursor = i
				e.adjustScroll()
				return
			}
		}
	}
}

func (e *ExplorerModel) toggleExpand(idx int) {
	if idx < 0 || idx >= len(e.entries) {
		return
	}
	entry := &e.entries[idx]
	if !entry.IsDir {
		return
	}
	entry.Expanded = !entry.Expanded
	e.rebuild()
	// Restore cursor position (entry might have shifted)
	for i, en := range e.entries {
		if en.Path == entry.Path {
			e.cursor = i
			break
		}
	}
	e.adjustScroll()
}

func (e *ExplorerModel) adjustScroll() {
	viewH := e.height - 2 // account for borders
	if viewH < 1 {
		viewH = 1
	}
	if e.cursor < e.offset {
		e.offset = e.cursor
	}
	if e.cursor >= e.offset+viewH {
		e.offset = e.cursor - viewH + 1
	}
}

// rebuild regenerates the flat entry list from disk.
func (e *ExplorerModel) rebuild() {
	// Keep track of which directories were expanded
	expanded := make(map[string]bool)
	for _, en := range e.entries {
		if en.IsDir && en.Expanded {
			expanded[en.Path] = true
		}
	}

	e.entries = nil

	// Add ".." entry if not at filesystem root
	parent := filepath.Dir(e.rootDir)
	if parent != e.rootDir {
		e.entries = append(e.entries, FileEntry{
			Name:  "..",
			Path:  parent,
			IsDir: true,
			Depth: 0,
		})
	}

	e.buildDir(e.rootDir, 0, expanded)
	e.reapplyFilter()
}

func (e *ExplorerModel) buildDir(dir string, depth int, expanded map[string]bool) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	// Sort: directories first, then files, alphabetical within each group
	var dirs, files []os.DirEntry
	for _, de := range dirEntries {
		name := de.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if de.IsDir() {
			dirs = append(dirs, de)
		} else {
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".avsc" || ext == ".json" {
				files = append(files, de)
			}
		}
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	for _, d := range dirs {
		path := filepath.Join(dir, d.Name())
		if !dirHasRelevantFiles(path) {
			continue
		}
		isExpanded := expanded[path]
		e.entries = append(e.entries, FileEntry{
			Name:     d.Name(),
			Path:     path,
			IsDir:    true,
			Depth:    depth,
			Expanded: isExpanded,
		})
		if isExpanded {
			e.buildDir(path, depth+1, expanded)
		}
	}
	for _, f := range files {
		e.entries = append(e.entries, FileEntry{
			Name:  f.Name(),
			Path:  filepath.Join(dir, f.Name()),
			IsDir: false,
			Depth: depth,
		})
	}
}

// dirHasRelevantFiles checks recursively if a directory contains any .avsc or .json files.
func dirHasRelevantFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, de := range entries {
		name := de.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if de.IsDir() {
			if dirHasRelevantFiles(filepath.Join(dir, name)) {
				return true
			}
		} else {
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".avsc" || ext == ".json" {
				return true
			}
		}
	}
	return false
}

// Filter applies a search query recursively and expands directories to show matches.
func (e *ExplorerModel) Filter(query string) {
	e.filterQuery = strings.ToLower(query)
	e.matchIndices = nil
	e.matchIdx = 0
	if e.filterQuery == "" {
		return
	}

	// Walk the filesystem to find all matching files/dirs.
	// Collect paths to expand (ancestor dirs of matches).
	dirsToExpand := make(map[string]bool)
	_ = filepath.WalkDir(e.rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(name))
			if ext != ".avsc" && ext != ".json" {
				return nil
			}
		}
		if strings.Contains(strings.ToLower(name), e.filterQuery) {
			// Mark all ancestor directories for expansion
			p := filepath.Dir(path)
			for {
				if p == e.rootDir || p == "." || p == string(filepath.Separator) {
					break
				}
				// Safety: stop if we've gone above rootDir
				rel, err := filepath.Rel(e.rootDir, p)
				if err != nil || strings.HasPrefix(rel, "..") {
					break
				}
				dirsToExpand[p] = true
				next := filepath.Dir(p)
				if next == p {
					break
				}
				p = next
			}
		}
		return nil
	})

	if len(dirsToExpand) == 0 && e.filterQuery != "" {
		// No matches found — still rebuild match indices from current entries
		for i, entry := range e.entries {
			if strings.Contains(strings.ToLower(entry.Name), e.filterQuery) {
				e.matchIndices = append(e.matchIndices, i)
			}
		}
		return
	}

	// Merge with currently expanded directories
	expanded := make(map[string]bool)
	for _, en := range e.entries {
		if en.IsDir && en.Expanded {
			expanded[en.Path] = true
		}
	}
	for p := range dirsToExpand {
		expanded[p] = true
	}

	// Rebuild entries with expanded dirs
	e.entries = nil
	parent := filepath.Dir(e.rootDir)
	if parent != e.rootDir {
		e.entries = append(e.entries, FileEntry{
			Name: "..", Path: parent, IsDir: true, Depth: 0,
		})
	}
	e.buildDir(e.rootDir, 0, expanded)

	// Build match indices for entries whose name matches the query
	for i, entry := range e.entries {
		if strings.Contains(strings.ToLower(entry.Name), e.filterQuery) {
			e.matchIndices = append(e.matchIndices, i)
		}
	}
	if len(e.matchIndices) > 0 {
		e.cursor = e.matchIndices[0]
		e.adjustScroll()
	}
}

// reapplyFilter re-runs the current filter after a structural change.
func (e *ExplorerModel) reapplyFilter() {
	if e.filterQuery == "" {
		return
	}
	e.matchIndices = nil
	e.matchIdx = 0
	for i, entry := range e.entries {
		if strings.Contains(strings.ToLower(entry.Name), e.filterQuery) {
			e.matchIndices = append(e.matchIndices, i)
		}
	}
}

// ClearFilter removes the current filter.
func (e *ExplorerModel) ClearFilter() {
	e.filterQuery = ""
	e.matchIndices = nil
	e.matchIdx = 0
}

// FilteredCount returns the number of matches.
func (e *ExplorerModel) FilteredCount() int {
	return len(e.matchIndices)
}

// FilterIdx returns the current match index (0-based).
func (e *ExplorerModel) FilterIdx() int {
	return e.matchIdx
}

// NextMatch moves to the next filter match.
func (e *ExplorerModel) NextMatch() {
	if len(e.matchIndices) == 0 {
		return
	}
	e.matchIdx = (e.matchIdx + 1) % len(e.matchIndices)
	e.cursor = e.matchIndices[e.matchIdx]
	e.adjustScroll()
}

// PrevMatch moves to the previous filter match.
func (e *ExplorerModel) PrevMatch() {
	if len(e.matchIndices) == 0 {
		return
	}
	e.matchIdx--
	if e.matchIdx < 0 {
		e.matchIdx = len(e.matchIndices) - 1
	}
	e.cursor = e.matchIndices[e.matchIdx]
	e.adjustScroll()
}

// View renders the file explorer panel.
func (e ExplorerModel) View() string {
	if len(e.entries) == 0 {
		return "  (empty)"
	}

	viewH := e.height - 2
	if viewH < 1 {
		viewH = 1
	}

	var b strings.Builder
	end := e.offset + viewH
	if end > len(e.entries) {
		end = len(e.entries)
	}

	for i := e.offset; i < end; i++ {
		if i > e.offset {
			b.WriteByte('\n')
		}
		entry := e.entries[i]
		indent := strings.Repeat("  ", entry.Depth)

		var icon string
		if entry.IsDir {
			if entry.Expanded {
				icon = " "
			} else {
				icon = " "
			}
		} else {
			ext := strings.ToLower(filepath.Ext(entry.Name))
			if ext == ".avsc" {
				icon = "󰈙 "
			} else {
				icon = " "
			}
		}

		line := indent + icon + entry.Name

		// Truncate to width
		maxW := e.width - 2
		if maxW < 4 {
			maxW = 4
		}
		if len(line) > maxW {
			line = line[:maxW-1] + "…"
		}
		// Pad to width
		if len(line) < maxW {
			line += strings.Repeat(" ", maxW-len(line))
		}

		if i == e.cursor {
			b.WriteString(e.theme.TreeNodeActive.Render(line))
		} else if e.isMatch(i) {
			b.WriteString(lipgloss.NewStyle().Foreground(e.theme.Warning).Render(line))
		} else if entry.IsDir {
			b.WriteString(lipgloss.NewStyle().Foreground(e.theme.Primary).Render(line))
		} else {
			b.WriteString(line)
		}
	}

	return b.String()
}

// isMatch returns true if the entry at index i is a filter match.
func (e ExplorerModel) isMatch(i int) bool {
	for _, idx := range e.matchIndices {
		if idx == i {
			return true
		}
	}
	return false
}
