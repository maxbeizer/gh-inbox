package model

import (
	"testing"
	"time"
)

func TestSubjectTypeIcon(t *testing.T) {
	tests := []struct {
		name string
		st   SubjectType
		want string
	}{
		{"PullRequest", SubjectPullRequest, "🔀"},
		{"Issue", SubjectIssue, "🔴"},
		{"Commit", SubjectCommit, "📝"},
		{"Release", SubjectRelease, "🏷️"},
		{"Discussion", SubjectDiscussion, "💬"},
		{"CheckSuite", SubjectCheckSuite, "✅"},
		{"Unknown", SubjectUnknown, "📌"},
		{"empty", SubjectType(""), "📌"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.st.Icon(); got != tt.want {
				t.Errorf("SubjectType(%q).Icon() = %q, want %q", tt.st, got, tt.want)
			}
		})
	}
}

func TestReasonIcon(t *testing.T) {
	tests := []struct {
		name string
		r    Reason
		want string
	}{
		{"mention", ReasonMention, "💬"},
		{"review_requested", ReasonReviewRequested, "👀"},
		{"assign", ReasonAssign, "👤"},
		{"subscribed", ReasonSubscribed, "🔔"},
		{"unknown", Reason("unknown_reason"), "📫"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Icon(); got != tt.want {
				t.Errorf("Reason(%q).Icon() = %q, want %q", tt.r, got, tt.want)
			}
		})
	}
}

func TestReasonLabel(t *testing.T) {
	tests := []struct {
		r    Reason
		want string
	}{
		{ReasonReviewRequested, "review requested"},
		{ReasonStateChange, "state change"},
		{ReasonMention, "mention"},
	}
	for _, tt := range tests {
		t.Run(string(tt.r), func(t *testing.T) {
			if got := tt.r.Label(); got != tt.want {
				t.Errorf("Reason(%q).Label() = %q, want %q", tt.r, got, tt.want)
			}
		})
	}
}

func TestNotificationHTMLURL(t *testing.T) {
	tests := []struct {
		name string
		n    Notification
		want string
	}{
		{
			name: "PR",
			n: Notification{
				Subject: Subject{
					Type:   SubjectPullRequest,
					Number: 42,
				},
				Repository: Repository{
					HTMLURL: "https://github.com/owner/repo",
				},
			},
			want: "https://github.com/owner/repo/pull/42",
		},
		{
			name: "Issue",
			n: Notification{
				Subject: Subject{
					Type:   SubjectIssue,
					Number: 123,
				},
				Repository: Repository{
					HTMLURL: "https://github.com/owner/repo",
				},
			},
			want: "https://github.com/owner/repo/issues/123",
		},
		{
			name: "no number fallback",
			n: Notification{
				Subject: Subject{
					Type: SubjectCommit,
				},
				Repository: Repository{
					HTMLURL: "https://github.com/owner/repo",
				},
			},
			want: "https://github.com/owner/repo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.HTMLURL(); got != tt.want {
				t.Errorf("HTMLURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRelativeTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		{"just now", now.Add(-10 * time.Second), "just now"},
		{"minutes", now.Add(-5 * time.Minute), "5m ago"},
		{"1 minute", now.Add(-1 * time.Minute), "1m ago"},
		{"hours", now.Add(-3 * time.Hour), "3h ago"},
		{"1 hour", now.Add(-1 * time.Hour), "1h ago"},
		{"days", now.Add(-5 * 24 * time.Hour), "5d ago"},
		{"1 day", now.Add(-1 * 24 * time.Hour), "1d ago"},
		{"months", now.Add(-60 * 24 * time.Hour), "2mo ago"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RelativeTime(tt.t); got != tt.want {
				t.Errorf("RelativeTime() = %q, want %q", got, tt.want)
			}
		})
	}
}
