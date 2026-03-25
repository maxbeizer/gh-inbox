package filter

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"

	"github.com/maxbeizer/gh-inbox/internal/tui/theme"
)

// Model is the search/filter input component.
type Model struct {
	input  textinput.Model
	active bool
	width  int
}

// New creates a new filter model.
func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Search notifications..."
	ti.CharLimit = 100
	return Model{input: ti}
}

// SetWidth sets the filter bar width.
func (m *Model) SetWidth(w int) {
	m.width = w
}

// Active returns whether the filter input is focused.
func (m Model) Active() bool {
	return m.active
}

// Activate focuses the filter input.
func (m *Model) Activate() {
	m.active = true
	m.input.Focus()
}

// Deactivate unfocuses the filter input.
func (m *Model) Deactivate() {
	m.active = false
	m.input.Blur()
}

// Value returns the current search text.
func (m Model) Value() string {
	return m.input.Value()
}

// SetValue sets the search text.
func (m *Model) SetValue(v string) {
	m.input.SetValue(v)
}

// Clear resets the search text and deactivates.
func (m *Model) Clear() {
	m.input.SetValue("")
	m.Deactivate()
}

// TextInput returns a reference to the underlying text input for Update routing.
func (m *Model) TextInput() *textinput.Model {
	return &m.input
}

// View renders the filter input bar.
func (m Model) View() string {
	if !m.active {
		return ""
	}

	prompt := theme.SearchPromptStyle.Render("🔍 ")
	inputView := m.input.View()
	row := lipgloss.JoinHorizontal(lipgloss.Center, prompt, inputView)

	return lipgloss.NewStyle().
		Width(m.width).
		Background(theme.ColorSurface0).
		Padding(0, 1).
		Render(row)
}
