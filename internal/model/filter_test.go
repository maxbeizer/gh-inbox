package model

import (
	"testing"
)

func TestFilterModeLabel(t *testing.T) {
	tests := []struct {
		mode FilterMode
		want string
	}{
		{FilterUnread, "Unread"},
		{FilterAll, "All"},
		{FilterParticipating, "Participating"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.mode.Label(); got != tt.want {
				t.Errorf("FilterMode(%d).Label() = %q, want %q", tt.mode, got, tt.want)
			}
		})
	}
}

func TestFilterModeNext(t *testing.T) {
	tests := []struct {
		mode FilterMode
		want FilterMode
	}{
		{FilterUnread, FilterAll},
		{FilterAll, FilterParticipating},
		{FilterParticipating, FilterUnread},
	}
	for _, tt := range tests {
		t.Run(tt.mode.Label(), func(t *testing.T) {
			if got := tt.mode.Next(); got != tt.want {
				t.Errorf("FilterMode(%d).Next() = %d, want %d", tt.mode, got, tt.want)
			}
		})
	}
}

func TestSortFieldLabel(t *testing.T) {
	tests := []struct {
		sort SortField
		want string
	}{
		{SortUpdated, "Updated"},
		{SortRepo, "Repository"},
		{SortReason, "Reason"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.sort.Label(); got != tt.want {
				t.Errorf("SortField(%d).Label() = %q, want %q", tt.sort, got, tt.want)
			}
		})
	}
}

func TestSortFieldNext(t *testing.T) {
	tests := []struct {
		sort SortField
		want SortField
	}{
		{SortUpdated, SortRepo},
		{SortRepo, SortReason},
		{SortReason, SortUpdated},
	}
	for _, tt := range tests {
		t.Run(tt.sort.Label(), func(t *testing.T) {
			if got := tt.sort.Next(); got != tt.want {
				t.Errorf("SortField(%d).Next() = %d, want %d", tt.sort, got, tt.want)
			}
		})
	}
}
