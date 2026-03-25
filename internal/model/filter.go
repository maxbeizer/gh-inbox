package model

// FilterMode controls which notifications to display.
type FilterMode int

const (
	FilterUnread FilterMode = iota
	FilterAll
	FilterParticipating
)

// Label returns a display label for the filter mode.
func (f FilterMode) Label() string {
	switch f {
	case FilterAll:
		return "All"
	case FilterParticipating:
		return "Participating"
	default:
		return "Unread"
	}
}

// Next cycles to the next filter mode.
func (f FilterMode) Next() FilterMode {
	return (f + 1) % 3
}

// SortField controls notification sort order.
type SortField int

const (
	SortUpdated SortField = iota
	SortRepo
	SortReason
)

// Label returns a display label for the sort field.
func (s SortField) Label() string {
	switch s {
	case SortRepo:
		return "Repository"
	case SortReason:
		return "Reason"
	default:
		return "Updated"
	}
}

// Next cycles to the next sort field.
func (s SortField) Next() SortField {
	return (s + 1) % 3
}

// Filters holds the current filter/sort/search state.
type Filters struct {
	Mode       FilterMode
	Sort       SortField
	SearchText string
	RepoFilter string
	TypeFilter SubjectType
}
