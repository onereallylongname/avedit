// Package theme defines the visual theme system for avedit.
// Themes are collections of lipgloss styles used by all UI components.
package theme

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Theme holds all lipgloss styles used across the UI.
type Theme struct {
	Name string

	// Base colors
	Fg    color.Color
	Bg    color.Color
	Dim   color.Color
	Muted color.Color

	// Accent colors
	Primary   color.Color
	Secondary color.Color
	Success   color.Color
	Warning   color.Color
	Error     color.Color

	// Panel styles
	TreeBorder    lipgloss.Style
	DetailBorder  lipgloss.Style
	FocusedBorder lipgloss.Style

	// Tree elements
	TreeNode       lipgloss.Style
	TreeNodeActive lipgloss.Style
	TreeIndent     lipgloss.Style
	TreeExpander   lipgloss.Style

	// Node kind badges
	KindRecord lipgloss.Style
	KindField  lipgloss.Style
	KindEnum   lipgloss.Style
	KindArray  lipgloss.Style
	KindMap    lipgloss.Style
	KindUnion  lipgloss.Style
	KindFixed  lipgloss.Style
	KindPrim   lipgloss.Style
	KindNamed  lipgloss.Style

	// Details panel
	DetailKey   lipgloss.Style
	DetailValue lipgloss.Style
	DetailTitle lipgloss.Style

	// Status bar
	StatusBar     lipgloss.Style
	StatusMode    lipgloss.Style
	StatusFile    lipgloss.Style
	StatusStats   lipgloss.Style
	StatusHelp    lipgloss.Style
	StatusNormal  lipgloss.Style
	StatusEdit    lipgloss.Style
	StatusSearch  lipgloss.Style
	StatusCommand lipgloss.Style
}

// ExpanderChars defines the characters used for tree expand/collapse indicators.
type ExpanderChars struct {
	Expanded  string
	Collapsed string
	Leaf      string
}

// DefaultExpanders provides the standard tree expander characters.
var DefaultExpanders = ExpanderChars{
	Expanded:  "▼",
	Collapsed: "▶",
	Leaf:      " ",
}
