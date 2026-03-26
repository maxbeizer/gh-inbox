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

type legend struct {
	icon string
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

var reasonLegend = []legend{
	{"💬", "Mentioned"},
	{"👥", "Team mentioned"},
	{"👀", "Review requested"},
	{"👤", "Assigned"},
	{"✍️", "You authored"},
	{"💭", "You commented"},
	{"🔄", "State changed"},
	{"🔔", "Subscribed / watching"},
	{"📌", "Manually subscribed"},
	{"⚙️", "CI activity"},
	{"🚨", "Security alert"},
	{"✋", "Approval requested"},
	{"📫", "Other"},
}

var typeLegend = []legend{
	{"🔀", "Pull Request"},
	{"🔴", "Issue"},
	{"📝", "Commit"},
	{"🏷️", "Release"},
	{"💬", "Discussion"},
	{"✅", "Check Suite"},
	{"📌", "Other"},
}

var statusLegend = []legend{
	{"●", "Unread"},
	{"○", "Read"},
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

	var left strings.Builder
	var right strings.Builder

	// ── Left column: keyboard shortcuts ──
	left.WriteString(theme.HelpTitleStyle.Render("⌨️  Keyboard Shortcuts"))
	left.WriteString("\n\n")

	for _, section := range sections {
		left.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.ColorPeach).
			Render(section.title))
		left.WriteString("\n")

		for _, bind := range section.bindings {
			key := theme.HelpKeyStyle.Render(bind.key)
			desc := theme.HelpDescStyle.Render(bind.desc)
			left.WriteString("  " + key + desc + "\n")
		}
		left.WriteString("\n")
	}

	// ── Right column: icon legends ──
	right.WriteString(theme.HelpTitleStyle.Render("📖  Icon Legend"))
	right.WriteString("\n\n")

	iconStyle := lipgloss.NewStyle().Width(4)
	descStyle := theme.HelpDescStyle

	// Reason icons (first column in each row)
	right.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.ColorPeach).
		Render("Reason (why you were notified)"))
	right.WriteString("\n")
	for _, item := range reasonLegend {
		right.WriteString("  " + iconStyle.Render(item.icon) + descStyle.Render(item.desc) + "\n")
	}
	right.WriteString("\n")

	// Type icons (second column in each row)
	right.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.ColorPeach).
		Render("Type (what the notification is about)"))
	right.WriteString("\n")
	for _, item := range typeLegend {
		right.WriteString("  " + iconStyle.Render(item.icon) + descStyle.Render(item.desc) + "\n")
	}
	right.WriteString("\n")

	// Status
	right.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.ColorPeach).
		Render("Status"))
	right.WriteString("\n")
	for _, item := range statusLegend {
		right.WriteString("  " + iconStyle.Render(item.icon) + descStyle.Render(item.desc) + "\n")
	}

	// Layout: two columns
	colW := (m.width - 10) / 2
	leftCol := lipgloss.NewStyle().Width(colW).Render(left.String())
	rightCol := lipgloss.NewStyle().Width(colW).Render(right.String())

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, "  ", rightCol)
	content += "\n\n" + theme.PreviewMetaStyle.Render("Press ? or Esc to close")

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		MaxHeight(m.height).
		Background(theme.ColorBase).
		Padding(2, 4).
		Render(content)
}
