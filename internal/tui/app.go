package tui

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"

	"github.com/maxbeizer/gh-inbox/internal/api"
	"github.com/maxbeizer/gh-inbox/internal/model"
	"github.com/maxbeizer/gh-inbox/internal/tui/components/filter"
	"github.com/maxbeizer/gh-inbox/internal/tui/components/header"
	"github.com/maxbeizer/gh-inbox/internal/tui/components/help"
	"github.com/maxbeizer/gh-inbox/internal/tui/components/preview"
	"github.com/maxbeizer/gh-inbox/internal/tui/components/statusbar"
	"github.com/maxbeizer/gh-inbox/internal/tui/components/table"
	"github.com/maxbeizer/gh-inbox/internal/tui/theme"
)

// App is the root Bubble Tea model.
type App struct {
	client *api.Client
	keys   KeyMap

	// Data
	allNotifications []model.Notification
	filtered         []model.Notification
	filters          model.Filters

	// Components
	header    header.Model
	table     table.Model
	preview   preview.Model
	statusbar statusbar.Model
	filter    filter.Model
	help      help.Model

	// State
	width    int
	height   int
	loading  bool
	ready    bool
	quitting bool
}

// NewApp creates the root TUI application model.
func NewApp(client *api.Client) App {
	return App{
		client:    client,
		keys:      DefaultKeyMap(),
		header:    header.New(),
		table:     table.New(),
		preview:   preview.New(),
		statusbar: statusbar.New(),
		filter:    filter.New(),
		help:      help.New(),
	}
}

// Init starts the application by fetching notifications.
func (a App) Init() tea.Cmd {
	return a.fetchNotifications()
}

// Update handles messages.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true
		a.layoutComponents()
		return a, nil

	case NotificationsFetchedMsg:
		a.loading = false
		a.header.SetLoading(false)
		if msg.Err != nil {
			a.statusbar.SetStatus(fmt.Sprintf("Error: %v", msg.Err), true)
			return a, nil
		}
		a.allNotifications = msg.Notifications
		a.applyFilters()
		a.statusbar.SetStatus(fmt.Sprintf("Loaded %d notifications", len(msg.Notifications)), false)
		return a, nil

	case SubjectDetailFetchedMsg:
		if msg.Err != nil {
			a.preview.SetDetail(msg.ID, nil)
			return a, nil
		}
		a.preview.SetDetail(msg.ID, msg.Detail)
		return a, nil

	case ActionCompleteMsg:
		if msg.Err != nil {
			a.statusbar.SetStatus(fmt.Sprintf("Error: %v", msg.Err), true)
			return a, nil
		}
		switch msg.Action {
		case "read":
			a.markNotifRead(msg.ID)
			a.statusbar.SetStatus("Marked as read", false)
		case "done":
			a.removeNotif(msg.ID)
			a.statusbar.SetStatus("Marked as done", false)
		case "unsubscribe":
			a.removeNotif(msg.ID)
			a.statusbar.SetStatus("Unsubscribed", false)
		}
		return a, nil

	case tea.KeyMsg:
		return a.handleKey(msg)

	case tea.MouseMsg:
		return a.handleMouse(msg)
	}

	return a, tea.Batch(cmds...)
}

func (a App) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If filter input is active, route keys there
	if a.filter.Active() {
		return a.handleFilterKey(msg)
	}

	// If help is visible, any key closes it
	if a.help.Visible() {
		if key.Matches(msg, a.keys.Help) || key.Matches(msg, a.keys.Escape) || key.Matches(msg, a.keys.Quit) {
			a.help.Toggle()
		}
		return a, nil
	}

	switch {
	case key.Matches(msg, a.keys.Quit):
		a.quitting = true
		return a, tea.Quit

	case key.Matches(msg, a.keys.Up):
		a.table.MoveUp()
		return a, a.onSelectionChange()

	case key.Matches(msg, a.keys.Down):
		a.table.MoveDown()
		return a, a.onSelectionChange()

	case key.Matches(msg, a.keys.Top):
		a.table.GoToTop()
		return a, a.onSelectionChange()

	case key.Matches(msg, a.keys.Bottom):
		a.table.GoToBottom()
		return a, a.onSelectionChange()

	case key.Matches(msg, a.keys.PageDown):
		if a.preview.Visible() {
			a.preview.ScrollDown()
		} else {
			a.table.PageDown()
		}
		return a, a.onSelectionChange()

	case key.Matches(msg, a.keys.PageUp):
		if a.preview.Visible() {
			a.preview.ScrollUp()
		} else {
			a.table.PageUp()
		}
		return a, a.onSelectionChange()

	case key.Matches(msg, a.keys.Preview), key.Matches(msg, a.keys.Enter):
		a.preview.Toggle()
		a.layoutComponents()
		if a.preview.Visible() {
			return a, a.fetchSelectedDetail()
		}
		return a, nil

	case key.Matches(msg, a.keys.Escape):
		if a.preview.Visible() {
			a.preview.Toggle()
			a.layoutComponents()
		}
		return a, nil

	case key.Matches(msg, a.keys.Refresh):
		return a, a.refreshNotifications(false)

	case key.Matches(msg, a.keys.RefreshAll):
		return a, a.refreshNotifications(true)

	case key.Matches(msg, a.keys.MarkRead):
		return a, a.markSelectedRead()

	case key.Matches(msg, a.keys.MarkUnread):
		return a, a.markSelectedUnread()

	case key.Matches(msg, a.keys.MarkDone):
		return a, a.markSelectedDone()

	case key.Matches(msg, a.keys.Unsubscribe):
		return a, a.unsubscribeSelected()

	case key.Matches(msg, a.keys.Open):
		return a, a.openSelected()

	case key.Matches(msg, a.keys.CopyURL):
		a.copySelectedURL()
		return a, nil

	case key.Matches(msg, a.keys.Search):
		a.filter.Activate()
		return a, nil

	case key.Matches(msg, a.keys.Filter):
		a.filters.Mode = a.filters.Mode.Next()
		a.applyFilters()
		return a, nil

	case key.Matches(msg, a.keys.Sort):
		a.filters.Sort = a.filters.Sort.Next()
		a.applyFilters()
		return a, nil

	case key.Matches(msg, a.keys.Help):
		a.help.Toggle()
		return a, nil
	}

	return a, nil
}

func (a App) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, a.keys.Escape):
		a.filter.Clear()
		a.filters.SearchText = ""
		a.applyFilters()
		return a, nil

	case key.Matches(msg, a.keys.Enter):
		a.filters.SearchText = a.filter.Value()
		a.filter.Deactivate()
		a.applyFilters()
		return a, nil

	default:
		ti := a.filter.TextInput()
		newTI, cmd := ti.Update(msg)
		*ti = newTI
		return a, cmd
	}
}

func (a App) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Basic mouse support: click on row to select
	if click, ok := msg.(tea.MouseClickMsg); ok {
		m := click.Mouse()
		// Approximate: header=1 line, each row after that
		row := m.Y - 2 // account for header bar + table header
		if row >= 0 && row < len(a.filtered) {
			for row > a.table.Cursor() {
				a.table.MoveDown()
			}
			for row < a.table.Cursor() {
				a.table.MoveUp()
			}
			return a, a.onSelectionChange()
		}
	}
	return a, nil
}

// View renders the entire application.
func (a App) View() tea.View {
	var content string

	if a.quitting {
		content = ""
	} else if !a.ready {
		content = lipgloss.NewStyle().
			Foreground(theme.ColorMauve).
			Padding(1, 2).
			Render("📬 Loading gh-inbox...")
	} else if a.help.Visible() {
		content = a.help.View()
	} else {
		headerView := a.header.View()
		footerView := a.statusbar.View()
		headerH := lipgloss.Height(headerView)
		footerH := lipgloss.Height(footerView)

		// Main content fills whatever is left
		mainH := a.height - headerH - footerH
		if a.filter.Active() {
			mainH--
		}
		if mainH < 1 {
			mainH = 1
		}

		var sections []string
		sections = append(sections, headerView)

		if a.filter.Active() {
			sections = append(sections, a.filter.View())
		}

		// Main content: table (+ optional preview) — fixed height to push footer down
		tableView := a.table.View()
		var mainContent string
		if a.preview.Visible() {
			previewView := a.preview.View()
			mainContent = lipgloss.JoinHorizontal(lipgloss.Top, tableView, previewView)
		} else {
			mainContent = tableView
		}
		mainContent = lipgloss.NewStyle().
			Height(mainH).
			Width(a.width).
			Render(mainContent)
		sections = append(sections, mainContent)
		sections = append(sections, footerView)

		content = lipgloss.JoinVertical(lipgloss.Left, sections...)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

// layoutComponents distributes width/height to child components.
func (a *App) layoutComponents() {
	a.header.SetWidth(a.width)
	a.statusbar.SetWidth(a.width)
	a.filter.SetWidth(a.width)
	a.help.SetSize(a.width, a.height)

	// Measure rendered heights to get accurate remaining space
	headerH := lipgloss.Height(a.header.View())
	footerH := max(lipgloss.Height(a.statusbar.View()), 1)
	ch := a.height - headerH - footerH
	if a.filter.Active() {
		ch--
	}
	if ch < 1 {
		ch = 1
	}

	if a.preview.Visible() {
		tableW := a.width * 55 / 100
		previewW := a.width - tableW
		a.table.SetSize(tableW, ch)
		a.preview.SetSize(previewW, ch)
	} else {
		a.table.SetSize(a.width, ch)
	}
}

// applyFilters filters and sorts notifications, then updates components.
func (a *App) applyFilters() {
	filtered := make([]model.Notification, 0, len(a.allNotifications))

	for _, n := range a.allNotifications {
		// Filter mode
		switch a.filters.Mode {
		case model.FilterUnread:
			if !n.Unread {
				continue
			}
		case model.FilterParticipating:
			// participating = mentioned, assigned, review requested, authored
			switch n.Reason {
			case model.ReasonMention, model.ReasonAssign, model.ReasonReviewRequested,
				model.ReasonAuthor, model.ReasonTeamMention, model.ReasonComment:
				// keep
			default:
				continue
			}
		}

		// Text search
		if a.filters.SearchText != "" {
			q := strings.ToLower(a.filters.SearchText)
			title := strings.ToLower(n.Subject.Title)
			repo := strings.ToLower(n.Repository.FullName)
			reason := strings.ToLower(string(n.Reason))
			if !strings.Contains(title, q) && !strings.Contains(repo, q) && !strings.Contains(reason, q) {
				continue
			}
		}

		filtered = append(filtered, n)
	}

	// Sort
	sort.Slice(filtered, func(i, j int) bool {
		switch a.filters.Sort {
		case model.SortRepo:
			if filtered[i].Repository.FullName != filtered[j].Repository.FullName {
				return filtered[i].Repository.FullName < filtered[j].Repository.FullName
			}
			return filtered[i].UpdatedAt.After(filtered[j].UpdatedAt)
		case model.SortReason:
			if filtered[i].Reason != filtered[j].Reason {
				return filtered[i].Reason < filtered[j].Reason
			}
			return filtered[i].UpdatedAt.After(filtered[j].UpdatedAt)
		default:
			return filtered[i].UpdatedAt.After(filtered[j].UpdatedAt)
		}
	})

	a.filtered = filtered
	a.table.SetNotifications(filtered)
	a.header.SetFilters(a.filters)
	a.header.SetCount(len(filtered), len(a.allNotifications))
	a.statusbar.SetFilters(a.filters)
	a.statusbar.SetCount(len(filtered), len(a.allNotifications))
}

// Command helpers

func (a *App) fetchNotifications() tea.Cmd {
	a.loading = true
	a.header.SetLoading(true)
	all := a.filters.Mode == model.FilterAll
	participating := a.filters.Mode == model.FilterParticipating
	return func() tea.Msg {
		notifs, err := a.client.ListNotifications(all, participating)
		return NotificationsFetchedMsg{Notifications: notifs, Err: err}
	}
}

func (a *App) refreshNotifications(all bool) tea.Cmd {
	a.loading = true
	a.header.SetLoading(true)
	a.statusbar.SetStatus("Refreshing...", false)
	return func() tea.Msg {
		notifs, err := a.client.ListNotifications(all, false)
		return NotificationsFetchedMsg{Notifications: notifs, Err: err}
	}
}

func (a *App) onSelectionChange() tea.Cmd {
	if a.preview.Visible() {
		return a.fetchSelectedDetail()
	}
	return nil
}

func (a *App) fetchSelectedDetail() tea.Cmd {
	n := a.table.Selected()
	if n == nil {
		return nil
	}
	a.preview.SetNotification(n)

	// Check cache
	if a.preview.GetCached(n.ID) != nil {
		return nil
	}

	// Only fetch detail for types with a subject API URL
	if n.Subject.URL == "" {
		return nil
	}

	// Support PR, Issue, and Release previews
	isPR := n.Subject.Type == model.SubjectPullRequest
	isIssue := n.Subject.Type == model.SubjectIssue
	isRelease := n.Subject.Type == model.SubjectRelease

	if !isPR && !isIssue && !isRelease {
		return nil
	}

	a.preview.SetLoading(true)
	subjectURL := n.Subject.URL
	id := n.ID
	return func() tea.Msg {
		var detail *api.SubjectDetail
		var err error
		if isRelease {
			detail, err = a.client.FetchReleaseDetail(subjectURL)
		} else {
			detail, err = a.client.FetchSubjectDetail(subjectURL, isPR)
		}
		return SubjectDetailFetchedMsg{Detail: detail, ID: id, Err: err}
	}
}

func (a *App) markSelectedRead() tea.Cmd {
	n := a.table.Selected()
	if n == nil {
		return nil
	}
	id := n.ID
	return func() tea.Msg {
		err := a.client.MarkThreadRead(id)
		return ActionCompleteMsg{Action: "read", ID: id, Err: err}
	}
}

func (a *App) markSelectedUnread() tea.Cmd {
	n := a.table.Selected()
	if n == nil {
		return nil
	}
	// Toggle unread state locally — the API has no "mark unread" endpoint,
	// but the notification will reappear as unread on next sync if still active.
	for i := range a.allNotifications {
		if a.allNotifications[i].ID == n.ID {
			a.allNotifications[i].Unread = true
			break
		}
	}
	a.applyFilters()
	a.statusbar.SetStatus("Marked as unread", false)
	return nil
}

func (a *App) markSelectedDone() tea.Cmd {
	n := a.table.Selected()
	if n == nil {
		return nil
	}
	id := n.ID
	return func() tea.Msg {
		err := a.client.MarkThreadDone(id)
		return ActionCompleteMsg{Action: "done", ID: id, Err: err}
	}
}

func (a *App) unsubscribeSelected() tea.Cmd {
	n := a.table.Selected()
	if n == nil {
		return nil
	}
	id := n.ID
	return func() tea.Msg {
		err := a.client.Unsubscribe(id)
		return ActionCompleteMsg{Action: "unsubscribe", ID: id, Err: err}
	}
}

func (a *App) openSelected() tea.Cmd {
	n := a.table.Selected()
	if n == nil {
		return nil
	}
	url := n.HTMLURL()
	return func() tea.Msg {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		default:
			cmd = exec.Command("open", url)
		}
		_ = cmd.Run()
		return nil
	}
}

func (a *App) copySelectedURL() {
	n := a.table.Selected()
	if n == nil {
		return
	}
	url := n.HTMLURL()

	// Try pbcopy (macOS), xclip (linux)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	default:
		cmd = exec.Command("xclip", "-selection", "clipboard")
	}
	cmd.Stdin = strings.NewReader(url)
	if err := cmd.Run(); err != nil {
		a.statusbar.SetStatus("Failed to copy URL", true)
		return
	}
	a.statusbar.SetStatus("Copied URL", false)
}

func (a *App) markNotifRead(id string) {
	for i := range a.allNotifications {
		if a.allNotifications[i].ID == id {
			a.allNotifications[i].Unread = false
			break
		}
	}
	a.applyFilters()
}

func (a *App) removeNotif(id string) {
	for i := range a.allNotifications {
		if a.allNotifications[i].ID == id {
			a.allNotifications = append(a.allNotifications[:i], a.allNotifications[i+1:]...)
			break
		}
	}
	a.applyFilters()
}
