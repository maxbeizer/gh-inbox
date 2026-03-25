package table

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/maxbeizer/gh-inbox/internal/model"
	"github.com/maxbeizer/gh-inbox/internal/tui/theme"
)

// Model is the notification list table component.
type Model struct {
	notifications []model.Notification
	cursor        int
	offset        int
	width         int
	height        int
}

// New creates a new table model.
func New() Model {
	return Model{}
}

// SetSize sets the table dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// SetNotifications replaces the displayed notifications.
func (m *Model) SetNotifications(notifs []model.Notification) {
	m.notifications = notifs
	if m.cursor >= len(notifs) {
		m.cursor = max(0, len(notifs)-1)
	}
	m.clampOffset()
}

// Cursor returns the current cursor position.
func (m Model) Cursor() int {
	return m.cursor
}

// Selected returns the currently selected notification, if any.
func (m Model) Selected() *model.Notification {
	if len(m.notifications) == 0 {
		return nil
	}
	return &m.notifications[m.cursor]
}

// MoveUp moves the cursor up.
func (m *Model) MoveUp() {
	if m.cursor > 0 {
		m.cursor--
		m.clampOffset()
	}
}

// MoveDown moves the cursor down.
func (m *Model) MoveDown() {
	if m.cursor < len(m.notifications)-1 {
		m.cursor++
		m.clampOffset()
	}
}

// GoToTop moves cursor to the first item.
func (m *Model) GoToTop() {
	m.cursor = 0
	m.offset = 0
}

// GoToBottom moves cursor to the last item.
func (m *Model) GoToBottom() {
	if len(m.notifications) > 0 {
		m.cursor = len(m.notifications) - 1
		m.clampOffset()
	}
}

// PageDown scrolls down by visible height.
func (m *Model) PageDown() {
	visible := m.visibleRows()
	m.cursor = min(m.cursor+visible, len(m.notifications)-1)
	m.clampOffset()
}

// PageUp scrolls up by visible height.
func (m *Model) PageUp() {
	visible := m.visibleRows()
	m.cursor = max(m.cursor-visible, 0)
	m.clampOffset()
}

// RemoveAt removes the notification at position i and adjusts cursor.
func (m *Model) RemoveAt(i int) {
	if i < 0 || i >= len(m.notifications) {
		return
	}
	m.notifications = append(m.notifications[:i], m.notifications[i+1:]...)
	if m.cursor >= len(m.notifications) {
		m.cursor = max(0, len(m.notifications)-1)
	}
	m.clampOffset()
}

func (m *Model) clampOffset() {
	visible := m.visibleRows()
	if visible <= 0 {
		return
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}
}

func (m Model) visibleRows() int {
	// Subtract 1 for header row
	return max(m.height-1, 1)
}

// View renders the notification table.
func (m Model) View() string {
	if len(m.notifications) == 0 {
		empty := lipgloss.NewStyle().
			Foreground(theme.ColorOverlay1).
			Width(m.width).
			Align(lipgloss.Center).
			Padding(2, 0).
			Render("No notifications 🎉")
		return empty
	}

	var b strings.Builder

	// Column header
	header := m.renderHeader()
	b.WriteString(header)
	b.WriteString("\n")

	// Rows
	visible := m.visibleRows()
	end := min(m.offset+visible, len(m.notifications))
	for i := m.offset; i < end; i++ {
		row := m.renderRow(i, i == m.cursor)
		b.WriteString(row)
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (m Model) renderHeader() string {
	// Fixed widths: reason(3) + type(3) + read(3) + number(8) + time(10) = 27
	// Flexible: repo + title share remaining
	repoW, titleW := m.flexWidths()

	cols := []string{
		theme.TableHeaderStyle.Width(3).Render(""),
		theme.TableHeaderStyle.Width(3).Render(""),
		theme.TableHeaderStyle.Width(3).Render(""),
		theme.TableHeaderStyle.Width(repoW).Render("Repo"),
		theme.TableHeaderStyle.Width(8).Render("#"),
		theme.TableHeaderStyle.Width(titleW).Render("Title"),
		theme.TableHeaderStyle.Width(10).Render("Updated"),
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, cols...)
}

func (m Model) renderRow(idx int, selected bool) string {
	n := m.notifications[idx]
	repoW, titleW := m.flexWidths()

	// Read indicator
	readInd := theme.ReadIndicator.Render("○")
	if n.Unread {
		readInd = theme.UnreadIndicator.Render("●")
	}

	// Number display
	numStr := ""
	if n.Subject.Number > 0 {
		numStr = fmt.Sprintf("#%d", n.Subject.Number)
	}

	// Truncate strings to fit
	repo := truncate(n.Repository.FullName, repoW)
	title := truncate(n.Subject.Title, titleW)
	timeStr := model.RelativeTime(n.UpdatedAt)

	baseStyle := theme.NormalRowStyle
	if selected {
		baseStyle = theme.SelectedRowStyle
	} else if !n.Unread {
		baseStyle = theme.DimRowStyle
	}

	cols := []string{
		baseStyle.Width(3).Render(string(n.Reason.Icon())),
		baseStyle.Width(3).Render(string(n.Subject.Type.Icon())),
		baseStyle.Width(3).Render(readInd),
		baseStyle.Foreground(theme.ColorSapphire).Width(repoW).Render(repo),
		baseStyle.Foreground(theme.ColorOverlay1).Width(8).Render(numStr),
		baseStyle.Width(titleW).Render(title),
		baseStyle.Foreground(theme.ColorOverlay0).Width(10).Render(timeStr),
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cols...)
}

func (m Model) flexWidths() (repoW, titleW int) {
	fixed := 3 + 3 + 3 + 8 + 10 // reason + type + read + number + time
	remaining := m.width - fixed
	if remaining < 20 {
		remaining = 20
	}
	repoW = remaining * 35 / 100
	titleW = remaining - repoW
	return
}

func truncate(s string, maxW int) string {
	if maxW <= 0 {
		return ""
	}
	if len(s) <= maxW {
		return s
	}
	if maxW <= 3 {
		return s[:maxW]
	}
	return s[:maxW-1] + "…"
}
