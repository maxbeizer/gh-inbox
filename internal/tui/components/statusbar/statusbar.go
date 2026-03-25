package statusbar

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/maxbeizer/gh-inbox/internal/model"
	"github.com/maxbeizer/gh-inbox/internal/tui/theme"
)

// Powerline separator characters
const (
	sepRight = "\ue0b0" // 
	sepLeft  = "\ue0b2" // 
)

// Model is the bottom powerline-style status bar component.
type Model struct {
	width      int
	statusText string
	isError    bool
	filters    model.Filters
	count      int
	total      int
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

// SetFilters updates the displayed filter/sort state.
func (m *Model) SetFilters(f model.Filters) {
	m.filters = f
}

// SetCount updates the notification counts.
func (m *Model) SetCount(visible, total int) {
	m.count = visible
	m.total = total
}

// segment holds a powerline segment's content and colors.
type segment struct {
	text string
	fg   color.Color
	bg   color.Color
}

// View renders the powerline-style status bar.
func (m Model) View() string {
	// ── Left side: mode + filter info + status message ──
	leftSegs := m.buildLeftSegments()

	// ── Right side: key hints ──
	rightSegs := m.buildRightSegments()

	left := renderPowerlineLeft(leftSegs)
	right := renderPowerlineRight(rightSegs)

	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	gap := m.width - leftW - rightW
	if gap < 0 {
		gap = 0
	}

	mid := lipgloss.NewStyle().
		Background(theme.ColorBase).
		Width(gap).
		Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, left, mid, right)
}

func (m Model) buildLeftSegments() []segment {
	segs := []segment{}

	// Mode badge
	modeIcon := "📬"
	modeBg := theme.ColorMauve
	switch m.filters.Mode {
	case model.FilterAll:
		modeIcon = "📭"
		modeBg = theme.ColorBlue
	case model.FilterParticipating:
		modeIcon = "🙋"
		modeBg = theme.ColorTeal
	}
	segs = append(segs, segment{
		text: fmt.Sprintf(" %s %s ", modeIcon, m.filters.Mode.Label()),
		fg:   theme.ColorBase,
		bg:   modeBg,
	})

	// Sort indicator
	segs = append(segs, segment{
		text: fmt.Sprintf("  %s ", m.filters.Sort.Label()),
		fg:   theme.ColorText,
		bg:   theme.ColorSurface1,
	})

	// Count
	segs = append(segs, segment{
		text: fmt.Sprintf(" %d/%d ", m.count, m.total),
		fg:   theme.ColorSubtext0,
		bg:   theme.ColorSurface0,
	})

	// Search text (if active)
	if m.filters.SearchText != "" {
		segs = append(segs, segment{
			text: fmt.Sprintf(" 🔍 %s ", m.filters.SearchText),
			fg:   theme.ColorYellow,
			bg:   theme.ColorSurface0,
		})
	}

	// Status message
	if m.statusText != "" {
		fg := theme.ColorGreen
		icon := "✓"
		if m.isError {
			fg = theme.ColorRed
			icon = "✗"
		}
		segs = append(segs, segment{
			text: fmt.Sprintf(" %s %s ", icon, m.statusText),
			fg:   fg,
			bg:   theme.ColorSurface0,
		})
	}

	return segs
}

func (m Model) buildRightSegments() []segment {
	type hint struct {
		key  string
		desc string
		bg   color.Color
	}

	hints := []hint{
		{"j/k", "nav", theme.ColorSurface1},
		{"e", "done", theme.ColorSurface1},
		{"⇧I", "read", theme.ColorSurface1},
		{"⇧U", "unread", theme.ColorSurface1},
		{"⇧M", "unsub", theme.ColorSurface2},
		{"o", "open", theme.ColorSurface2},
		{"p", "preview", theme.ColorSurface2},
		{"/", "search", theme.ColorSurface2},
		{"?", "help", theme.ColorMauve},
	}

	segs := make([]segment, 0, len(hints))
	for _, h := range hints {
		text := fmt.Sprintf(" %s %s ",
			lipgloss.NewStyle().Bold(true).Foreground(theme.ColorLavender).Render(h.key),
			h.desc,
		)
		segs = append(segs, segment{
			text: text,
			fg:   theme.ColorText,
			bg:   h.bg,
		})
	}
	return segs
}

// renderPowerlineLeft renders segments with right-pointing arrow separators.
func renderPowerlineLeft(segs []segment) string {
	if len(segs) == 0 {
		return ""
	}

	var result string
	for i, seg := range segs {
		body := lipgloss.NewStyle().
			Foreground(seg.fg).
			Background(seg.bg).
			Render(seg.text)
		result += body

		// Arrow separator: fg = current bg, bg = next bg (or base)
		nextBg := theme.ColorBase
		if i+1 < len(segs) {
			nextBg = segs[i+1].bg
		}
		arrow := lipgloss.NewStyle().
			Foreground(seg.bg).
			Background(nextBg).
			Render(sepRight)
		result += arrow
	}

	return result
}

// renderPowerlineRight renders segments with left-pointing arrow separators.
func renderPowerlineRight(segs []segment) string {
	if len(segs) == 0 {
		return ""
	}

	var result string
	for i, seg := range segs {
		// Arrow separator: fg = current bg, bg = previous bg (or base)
		prevBg := theme.ColorBase
		if i > 0 {
			prevBg = segs[i-1].bg
		}
		arrow := lipgloss.NewStyle().
			Foreground(seg.bg).
			Background(prevBg).
			Render(sepLeft)
		result += arrow

		body := lipgloss.NewStyle().
			Foreground(seg.fg).
			Background(seg.bg).
			Render(seg.text)
		result += body
	}

	return result
}
