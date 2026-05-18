package model

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Editor is a lightweight inline text editor for single-line values.
type Editor struct {
	value  []rune
	cursor int
	active bool
	width  int
	style  EditorStyle
}

// EditorStyle holds the visual styles for the editor.
type EditorStyle struct {
	Text       lipgloss.Style
	Cursor     lipgloss.Style
	Placeholder lipgloss.Style
}

// NewEditor creates an editor with the given initial value.
func NewEditor(initial string, style EditorStyle) Editor {
	runes := []rune(initial)
	return Editor{
		value:  runes,
		cursor: len(runes),
		active: true,
		width:  40,
		style:  style,
	}
}

// Value returns the current text content.
func (e *Editor) Value() string {
	return string(e.value)
}

// SetValue replaces the editor content and moves cursor to end.
func (e *Editor) SetValue(s string) {
	e.value = []rune(s)
	e.cursor = len(e.value)
}

// InsertRune inserts a single rune at the cursor position.
func (e *Editor) InsertRune(ch rune) {
	e.value = append(e.value[:e.cursor], append([]rune{ch}, e.value[e.cursor:]...)...)
	e.cursor++
}

// SetWidth sets the visible width.
func (e *Editor) SetWidth(w int) {
	e.width = w
}

// Active returns whether the editor is in editing state.
func (e *Editor) Active() bool {
	return e.active
}

// Activate enables editing.
func (e *Editor) Activate() {
	e.active = true
}

// Deactivate disables editing.
func (e *Editor) Deactivate() {
	e.active = false
}

// HandleKey processes a keypress. Returns true if the key was consumed.
func (e *Editor) HandleKey(msg tea.KeyPressMsg) bool {
	if !e.active {
		return false
	}

	switch msg.String() {
	case "left":
		if e.cursor > 0 {
			e.cursor--
		}
	case "right":
		if e.cursor < len(e.value) {
			e.cursor++
		}
	case "home", "ctrl+a":
		e.cursor = 0
	case "end", "ctrl+e":
		e.cursor = len(e.value)
	case "backspace":
		if e.cursor > 0 {
			e.value = append(e.value[:e.cursor-1], e.value[e.cursor:]...)
			e.cursor--
		}
	case "delete":
		if e.cursor < len(e.value) {
			e.value = append(e.value[:e.cursor], e.value[e.cursor+1:]...)
		}
	case "ctrl+u":
		e.value = e.value[e.cursor:]
		e.cursor = 0
	case "ctrl+k":
		e.value = e.value[:e.cursor]
	case "ctrl+w":
		// Delete word backward
		if e.cursor > 0 {
			pos := e.cursor - 1
			for pos > 0 && e.value[pos] == ' ' {
				pos--
			}
			for pos > 0 && e.value[pos-1] != ' ' {
				pos--
			}
			e.value = append(e.value[:pos], e.value[e.cursor:]...)
			e.cursor = pos
		}
	default:
		// Insert printable character
		if !isControlKey(msg) {
			var ch rune
			if msg.Text != "" {
				// Use Text field (actual typed character with correct case)
				for _, r := range msg.Text {
					e.value = append(e.value[:e.cursor], append([]rune{r}, e.value[e.cursor:]...)...)
					e.cursor++
				}
				return true
			} else if msg.Code > 0 && msg.Code < 128 {
				ch = msg.Code
			}
			if ch > 0 {
				e.value = append(e.value[:e.cursor], append([]rune{ch}, e.value[e.cursor:]...)...)
				e.cursor++
			}
		}
	}
	return true
}

// View renders the editor content with a visible cursor.
func (e Editor) View() string {
	if !e.active {
		return e.style.Text.Render(string(e.value))
	}

	var sb strings.Builder
	text := e.value
	cursorPos := e.cursor

	// Render text before cursor
	if cursorPos > 0 {
		sb.WriteString(e.style.Text.Render(string(text[:cursorPos])))
	}

	// Render cursor character
	if cursorPos < len(text) {
		ch := text[cursorPos]
		cursorDisplay := string(ch)
		if ch == '\n' {
			cursorDisplay = "↵\n"
		}
		sb.WriteString(e.style.Cursor.Render(cursorDisplay))
		// Render text after cursor
		if cursorPos+1 < len(text) {
			sb.WriteString(e.style.Text.Render(string(text[cursorPos+1:])))
		}
	} else {
		// Cursor at end — show block cursor on space
		sb.WriteString(e.style.Cursor.Render(" "))
	}

	return sb.String()
}

// isControlKey returns true if the keypress is a non-printable control sequence.
func isControlKey(msg tea.KeyPressMsg) bool {
	s := msg.String()
	if len(s) == 0 {
		return true
	}
	// Known control prefixes
	if strings.HasPrefix(s, "ctrl+") || strings.HasPrefix(s, "alt+") {
		return true
	}
	// Named keys that aren't printable
	switch s {
	case "enter", "tab", "esc", "up", "down", "left", "right",
		"home", "end", "pgup", "pgdown", "insert", "delete",
		"backspace", "f1", "f2", "f3", "f4", "f5", "f6",
		"f7", "f8", "f9", "f10", "f11", "f12":
		return true
	}
	return false
}
