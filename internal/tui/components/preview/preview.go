package preview

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/maxbeizer/gh-inbox/internal/api"
	"github.com/maxbeizer/gh-inbox/internal/model"
	"github.com/maxbeizer/gh-inbox/internal/tui/theme"
)

const maxCacheSize = 100

// Model is the preview side panel component.
type Model struct {
	width      int
	height     int
	visible    bool
	detail     *api.SubjectDetail
	notif      *model.Notification
	loading    bool
	cache      map[string]*api.SubjectDetail
	scrollPos  int
	contentLen int
}

// New creates a new preview model.
func New() Model {
	return Model{
		cache: make(map[string]*api.SubjectDetail),
	}
}

// SetSize updates the preview panel dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Visible returns whether the preview is showing.
func (m Model) Visible() bool {
	return m.visible
}

// Toggle flips the preview visibility.
func (m *Model) Toggle() {
	m.visible = !m.visible
}

// SetLoading sets the loading state.
func (m *Model) SetLoading(loading bool) {
	m.loading = loading
}

// SetDetail sets the subject detail and caches it.
func (m *Model) SetDetail(id string, detail *api.SubjectDetail) {
	m.detail = detail
	if detail != nil {
		// Evict oldest entries if cache is full
		if len(m.cache) >= maxCacheSize {
			m.cache = make(map[string]*api.SubjectDetail)
		}
		m.cache[id] = detail
	}
	m.loading = false
	m.scrollPos = 0
}

// SetNotification sets which notification is being previewed.
func (m *Model) SetNotification(n *model.Notification) {
	m.notif = n
	m.scrollPos = 0
	if n != nil {
		if cached, ok := m.cache[n.ID]; ok {
			m.detail = cached
			m.loading = false
			return
		}
	}
	m.detail = nil
}

// GetCached returns cached detail for a notification ID, if available.
func (m *Model) GetCached(id string) *api.SubjectDetail {
	return m.cache[id]
}

// ClearCache removes all cached subject details.
func (m *Model) ClearCache() {
	m.cache = make(map[string]*api.SubjectDetail)
}

// ScrollDown scrolls the preview content down.
func (m *Model) ScrollDown() {
	m.scrollPos += 5
	maxScroll := m.contentLen - m.innerHeight()
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scrollPos > maxScroll {
		m.scrollPos = maxScroll
	}
}

// ScrollUp scrolls the preview content up.
func (m *Model) ScrollUp() {
	m.scrollPos -= 5
	if m.scrollPos < 0 {
		m.scrollPos = 0
	}
}

func (m Model) innerHeight() int {
	// Account for border + padding
	return max(m.height-4, 1)
}

func (m *Model) renderContent() string {
	if m.notif == nil {
		return lipgloss.NewStyle().
			Foreground(theme.ColorOverlay1).
			Render("Select a notification to preview")
	}

	var b strings.Builder

	// Notification metadata
	b.WriteString(theme.PreviewTitleStyle.Render(m.notif.Subject.Title))
	b.WriteString("\n\n")

	meta := fmt.Sprintf("%s  %s  %s  %s",
		m.notif.Subject.Type.Icon(),
		m.notif.Repository.FullName,
		m.notif.Reason.Icon(),
		m.notif.Reason.Label(),
	)
	if m.notif.Subject.Number > 0 {
		meta += fmt.Sprintf("  #%d", m.notif.Subject.Number)
	}
	b.WriteString(theme.PreviewMetaStyle.Render(meta))
	b.WriteString("\n")

	if m.loading {
		b.WriteString("\n")
		b.WriteString(theme.SpinnerStyle.Render("⟳ Loading details..."))
		return b.String()
	}

	if m.detail == nil {
		b.WriteString("\n")
		b.WriteString(theme.PreviewMetaStyle.Render("No detail available for this notification type"))
		return b.String()
	}

	// State and author
	b.WriteString("\n")
	stateIcon := "⚪"
	switch strings.ToLower(m.detail.State) {
	case "open":
		stateIcon = "🟢"
	case "closed":
		stateIcon = "🔴"
	case "merged":
		stateIcon = "🟣"
	}
	b.WriteString(fmt.Sprintf("%s %s  by @%s\n", stateIcon, m.detail.State, m.detail.Author))

	// Labels
	if len(m.detail.Labels) > 0 {
		b.WriteString("\n")
		for _, label := range m.detail.Labels {
			b.WriteString(theme.PreviewLabelStyle.Render(label))
			b.WriteString(" ")
		}
		b.WriteString("\n")
	}

	// Body
	if m.detail.Body != "" {
		b.WriteString("\n")
		bodyWidth := m.width - 6
		if bodyWidth < 20 {
			bodyWidth = 20
		}
		b.WriteString(lipgloss.NewStyle().
			Width(bodyWidth).
			Foreground(theme.ColorText).
			Render(m.detail.Body))
	}

	return b.String()
}

// View renders the preview panel.
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	content := m.renderContent()
	lines := strings.Split(content, "\n")
	m.contentLen = len(lines)

	// Apply scroll
	start := m.scrollPos
	if start >= len(lines) {
		start = max(len(lines)-1, 0)
	}
	end := start + m.innerHeight()
	if end > len(lines) {
		end = len(lines)
	}
	visible := strings.Join(lines[start:end], "\n")

	return theme.PreviewBorderStyle.
		Width(m.width).
		Height(m.height).
		Render(visible)
}
