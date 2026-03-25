package statusbar

import (
	"fmt"

	"charm.land/lipgloss/v2"

	"github.com/maxbeizer/gh-inbox/internal/tui/theme"
)

// Model is the bottom status bar component.
type Model struct {
	width      int
	statusText string
	isError    bool
}

// New creates a new statusbar model.
func New() Model {
	return Model{}
}

// SetWidth sets the status bar width.
func (m *Model) SetWidth(w int) {
	m.width = w
}

// SetStatus sets a temporary status message.
func (m *Model) SetStatus(text string, isError bool) {
	m.statusText = text
	m.isError = isError
}

// ClearStatus removes the status message.
func (m *Model) ClearStatus() {
	m.statusText = ""
	m.isError = false
}

// View renders the status bar.
func (m Model) View() string {
	var left string
	if m.statusText != "" {
		if m.isError {
			left = theme.ErrorStyle.Render("✗ " + m.statusText)
		} else {
			left = theme.SuccessStyle.Render("✓ " + m.statusText)
		}
	}

	hints := []struct{ key, desc string }{
		{"j/k", "navigate"},
		{"m", "read"},
		{"d", "done"},
		{"u", "unsub"},
		{"o", "open"},
		{"p", "preview"},
		{"/", "search"},
		{"?", "help"},
	}

	var right string
	for i, h := range hints {
		if i > 0 {
			right += "  "
		}
		right += fmt.Sprintf("%s %s",
			theme.StatusBarKeyStyle.Render(h.key),
			theme.StatusBarValueStyle.Render(h.desc),
		)
	}

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	spacer := lipgloss.NewStyle().Width(gap).Render("")

	row := lipgloss.JoinHorizontal(lipgloss.Top, left, spacer, right)

	return theme.StatusBarStyle.Width(m.width).Render(row)
}
