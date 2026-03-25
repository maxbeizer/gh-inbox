package preview

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"

	"github.com/maxbeizer/gh-inbox/internal/api"
	"github.com/maxbeizer/gh-inbox/internal/model"
	"github.com/maxbeizer/gh-inbox/internal/tui/theme"
)

// Model is the preview side panel component.
type Model struct {
	viewport viewport.Model
	width    int
	height   int
	visible  bool
	detail   *api.SubjectDetail
	notif    *model.Notification
	loading  bool
	cache    map[string]*api.SubjectDetail
}

// New creates a new preview model.
func New() Model {
	vp := viewport.New()
	return Model{
		viewport: vp,
		cache:    make(map[string]*api.SubjectDetail),
	}
}

// SetSize updates the preview panel dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.SetWidth(w - 4) // padding
	m.viewport.SetHeight(h - 2)
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
		m.cache[id] = detail
	}
	m.loading = false
	m.updateViewport()
}

// SetNotification sets which notification is being previewed.
func (m *Model) SetNotification(n *model.Notification) {
	m.notif = n
	if n != nil {
		if cached, ok := m.cache[n.ID]; ok {
			m.detail = cached
			m.loading = false
			m.updateViewport()
			return
		}
	}
	m.detail = nil
	m.updateViewport()
}

// GetCached returns cached detail for a notification ID, if available.
func (m *Model) GetCached(id string) *api.SubjectDetail {
	return m.cache[id]
}

// ScrollDown scrolls the preview viewport down.
func (m *Model) ScrollDown() {
	m.viewport.ScrollDown(5)
}

// ScrollUp scrolls the preview viewport up.
func (m *Model) ScrollUp() {
	m.viewport.ScrollUp(5)
}

func (m *Model) updateViewport() {
	content := m.renderContent()
	m.viewport.SetContent(content)
	m.viewport.GotoTop()
}

func (m Model) renderContent() string {
	if m.notif == nil {
		return lipgloss.NewStyle().
			Foreground(theme.ColorOverlay1).
			Render("Select a notification to preview")
	}

	var b strings.Builder

	// Notification metadata
	b.WriteString(theme.PreviewTitleStyle.Render(m.notif.Subject.Title))
	b.WriteString("\n\n")

	meta := []string{
		fmt.Sprintf("%s  %s  %s  %s #%d",
			m.notif.Subject.Type.Icon(),
			m.notif.Repository.FullName,
			m.notif.Reason.Icon(),
			m.notif.Reason.Label(),
			m.notif.Subject.Number,
		),
	}
	b.WriteString(theme.PreviewMetaStyle.Render(strings.Join(meta, "\n")))
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
		b.WriteString(lipgloss.NewStyle().
			Width(m.width - 6).
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

	content := m.viewport.View()

	return theme.PreviewBorderStyle.
		Width(m.width).
		Height(m.height).
		Render(content)
}
