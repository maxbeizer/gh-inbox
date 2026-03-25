package tui

import (
	"github.com/maxbeizer/gh-inbox/internal/api"
	"github.com/maxbeizer/gh-inbox/internal/model"
)

// NotificationsFetchedMsg is sent when notifications are loaded from the API.
type NotificationsFetchedMsg struct {
	Notifications []model.Notification
	Err           error
}

// SubjectDetailFetchedMsg is sent when a subject detail is loaded for preview.
type SubjectDetailFetchedMsg struct {
	Detail *api.SubjectDetail
	ID     string
	Err    error
}

// ActionCompleteMsg is sent when a notification action (read, done, etc.) completes.
type ActionCompleteMsg struct {
	Action string
	ID     string
	Err    error
}

// AllReadMsg is sent when "mark all read" completes.
type AllReadMsg struct {
	Err error
}

// StatusMsg is a temporary status bar message.
type StatusMsg struct {
	Text    string
	IsError bool
}

// WindowSizeMsg wraps tea.WindowSizeMsg for component distribution.
type WindowSizeMsg struct {
	Width  int
	Height int
}
