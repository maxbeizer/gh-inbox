package help

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/maxbeizer/gh-inbox/internal/tui/theme"
)

type binding struct {
	key  string
	desc string
}

var sections = []struct {
	title    string
	bindings []binding
}{
	{
		title: "Navigation",
		bindings: []binding{
			{"j / ↓", "Move down"},
			{"k / ↑", "Move up"},
			{"g", "Go to top"},
			{"G", "Go to bottom"},
			{"ctrl+d", "Page down"},
			{"ctrl+u", "Page up"},
		},
	},
	{
		title: "Actions (matches GitHub web)",
		bindings: []binding{
			{"e", "Mark as done"},
			{"⇧I", "Mark as read"},
			{"⇧U", "Mark as unread"},
			{"⇧M", "Unsubscribe"},
			{"o", "Open in browser"},
			{"y", "Copy URL to clipboard"},
		},
	},
	{
		title: "View",
		bindings: []binding{
			{"p / Enter", "Toggle preview panel"},
			{"r", "Refresh notifications"},
			{"R", "Refresh all (include read)"},
			{"/", "Search / filter by text"},
			{"f", "Cycle filter mode"},
			{"s", "Cycle sort order"},
			{"Esc", "Close panel / clear search"},
		},
	},
	{
		title: "General",
		bindings: []binding{
			{"?", "Toggle this help"},
			{"q / ctrl+c", "Quit"},
		},
	},
}

// Model is the full-screen help overlay.
type Model struct {
	visible bool
	width   int
	height  int
}

// New creates a new help model.
func New() Model {
	return Model{}
}

// SetSize updates dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Visible returns whether help is showing.
func (m Model) Visible() bool {
	return m.visible
}

// Toggle flips help visibility.
func (m *Model) Toggle() {
	m.visible = !m.visible
}

// View renders the help overlay.
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	var b strings.Builder

	b.WriteString(theme.HelpTitleStyle.Render("⌨️  Keyboard Shortcuts"))
	b.WriteString("\n\n")

	for _, section := range sections {
		b.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.ColorPeach).
			MarginBottom(0).
			Render(section.title))
		b.WriteString("\n")

		for _, bind := range section.bindings {
			key := theme.HelpKeyStyle.Render(bind.key)
			desc := theme.HelpDescStyle.Render(bind.desc)
			b.WriteString("  " + key + desc + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(theme.PreviewMetaStyle.Render("Press ? or Esc to close"))

	content := b.String()

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(theme.ColorBase).
		Padding(2, 4).
		Render(content)
}
