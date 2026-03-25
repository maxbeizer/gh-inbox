package api

import (
	"testing"
)

func TestParseNumber(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "pulls URL",
			url:  "https://api.github.com/repos/owner/repo/pulls/42",
			want: 42,
		},
		{
			name: "issues URL",
			url:  "https://api.github.com/repos/owner/repo/issues/123",
			want: 123,
		},
		{
			name: "commits URL (no number)",
			url:  "https://api.github.com/repos/owner/repo/commits/abc123",
			want: 0,
		},
		{
			name: "empty URL",
			url:  "",
			want: 0,
		},
		{
			name: "large number",
			url:  "https://api.github.com/repos/o/r/pulls/99999",
			want: 99999,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNumber(tt.url); got != tt.want {
				t.Errorf("parseNumber(%q) = %d, want %d", tt.url, got, tt.want)
			}
		})
	}
}

func TestToSubjectType(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"PullRequest", "PullRequest"},
		{"Issue", "Issue"},
		{"Commit", "Commit"},
		{"Release", "Release"},
		{"Discussion", "Discussion"},
		{"CheckSuite", "CheckSuite"},
		{"AgentSessionThread", "Unknown"},
		{"", "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toSubjectType(tt.input)
			if string(got) != tt.want {
				t.Errorf("toSubjectType(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
