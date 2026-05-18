// Package command implements the undo/redo command pattern for schema mutations.
//
// Design: Each mutation is a Command struct containing Do/Undo closures and a description.
// Commands capture old state at creation time (via closure), enabling zero-copy undo.
// The History struct manages undo/redo stacks with a configurable depth limit.
//
// Command factories (CreateNode, RemoveNode, etc.) validate preconditions and return
// a Command or error. The caller then passes the Command to History.Execute().
//
// See docs/ARCHITECTURE.md "Command Layer" for full details.
package command

// Command represents a reversible mutation on the projection.
type Command struct {
	DoFn   func() error
	UndoFn func() error
	Desc   string
}

// Do executes the command's forward action.
func (c *Command) Do() error {
	return c.DoFn()
}

// Undo executes the command's reverse action.
func (c *Command) Undo() error {
	return c.UndoFn()
}

// History manages undo/redo stacks of executed commands.
type History struct {
	undoStack []*Command
	redoStack []*Command
	limit     int
}

// NewHistory creates a command history with the given max depth.
func NewHistory(limit int) *History {
	if limit <= 0 {
		limit = 500
	}
	return &History{
		undoStack: make([]*Command, 0, 64),
		redoStack: make([]*Command, 0, 64),
		limit:     limit,
	}
}

// Execute runs a command and pushes it onto the undo stack.
func (h *History) Execute(cmd *Command) error {
	if err := cmd.Do(); err != nil {
		return err
	}

	h.undoStack = append(h.undoStack, cmd)
	h.redoStack = h.redoStack[:0] // clear redo on new action

	// enforce max size
	if len(h.undoStack) > h.limit {
		h.undoStack = h.undoStack[1:]
	}

	return nil
}

// Undo reverses the last command.
func (h *History) Undo() error {
	if len(h.undoStack) == 0 {
		return nil
	}

	cmd := h.undoStack[len(h.undoStack)-1]
	h.undoStack = h.undoStack[:len(h.undoStack)-1]

	if err := cmd.Undo(); err != nil {
		return err
	}

	h.redoStack = append(h.redoStack, cmd)
	return nil
}

// Redo re-applies the last undone command.
func (h *History) Redo() error {
	if len(h.redoStack) == 0 {
		return nil
	}

	cmd := h.redoStack[len(h.redoStack)-1]
	h.redoStack = h.redoStack[:len(h.redoStack)-1]

	if err := cmd.Do(); err != nil {
		return err
	}

	h.undoStack = append(h.undoStack, cmd)
	return nil
}

// CanUndo returns true if there are commands to undo.
func (h *History) CanUndo() bool {
	return len(h.undoStack) > 0
}

// CanRedo returns true if there are commands to redo.
func (h *History) CanRedo() bool {
	return len(h.redoStack) > 0
}

// UndoCount returns the number of commands in the undo stack.
func (h *History) UndoCount() int {
	return len(h.undoStack)
}

// RedoCount returns the number of commands in the redo stack.
func (h *History) RedoCount() int {
	return len(h.redoStack)
}

// LastDesc returns the description of the last undoable command.
func (h *History) LastDesc() string {
	if len(h.undoStack) == 0 {
		return ""
	}
	return h.undoStack[len(h.undoStack)-1].Desc
}

// Reset clears all history.
func (h *History) Reset() {
	h.undoStack = h.undoStack[:0]
	h.redoStack = h.redoStack[:0]
}
