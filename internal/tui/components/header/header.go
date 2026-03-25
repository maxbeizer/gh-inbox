package header

import (
	"fmt"

	"charm.land/lipgloss/v2"

	"github.com/maxbeizer/gh-inbox/internal/model"
	"github.com/maxbeizer/gh-inbox/internal/tui/theme"
)

// Model is the header bar component.
type Model struct {
	width   int
	filters model.Filters
	count   int
	total   int
	loading bool
}

// New creates a new header model.
func New() Model {
	return Model{}
}

// SetWidth sets the header width.
func (m *Model) SetWidth(w int) {
	m.width = w
}

// SetFilters updates the displayed filter state.
func (m *Model) SetFilters(f model.Filters) {
	m.filters = f
}

// SetCount sets the displayed/total notification counts.
func (m *Model) SetCount(visible, total int) {
	m.count = visible
	m.total = total
}

// SetLoading sets the loading indicator state.
func (m *Model) SetLoading(loading bool) {
	m.loading = loading
}

// View renders the header.
func (m Model) View() string {
	title := theme.HeaderTitleStyle.Render("📬 gh-inbox")

	filterLabel := fmt.Sprintf(" [%s]", m.filters.Mode.Label())
	sortLabel := fmt.Sprintf(" sort:%s", m.filters.Sort.Label())
	filter := theme.HeaderFilterStyle.Render(filterLabel + sortLabel)

	if m.filters.SearchText != "" {
		filter += theme.SearchPromptStyle.Render(fmt.Sprintf(" 🔍 %q", m.filters.SearchText))
	}

	var status string
	if m.loading {
		status = theme.SpinnerStyle.Render(" ⟳ loading...")
	} else {
		status = theme.HeaderFilterStyle.Render(fmt.Sprintf(" %d/%d", m.count, m.total))
	}

	left := title + filter
	right := status

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 1
	}
	spacer := lipgloss.NewStyle().Width(gap).Render("")

	row := lipgloss.JoinHorizontal(lipgloss.Top, left, spacer, right)

	return theme.HeaderStyle.Width(m.width).Render(row)
}
