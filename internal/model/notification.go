package model

import (
	"fmt"
	"strings"
	"time"
)

// Notification represents a parsed GitHub notification thread.
type Notification struct {
	ID             string
	Unread         bool
	Reason         Reason
	UpdatedAt      time.Time
	LastReadAt     *time.Time
	Subject        Subject
	Repository     Repository
	URL            string
	SubscriptionURL string
}

// Subject is the item the notification is about.
type Subject struct {
	Title           string
	URL             string
	LatestCommentURL string
	Type            SubjectType
	Number          int
}

// Repository holds repo metadata from the notification.
type Repository struct {
	FullName string
	Owner    string
	Name     string
	Private  bool
	HTMLURL  string
}

// SubjectType is the kind of subject (PR, Issue, etc).
type SubjectType string

const (
	SubjectPullRequest SubjectType = "PullRequest"
	SubjectIssue       SubjectType = "Issue"
	SubjectCommit      SubjectType = "Commit"
	SubjectRelease     SubjectType = "Release"
	SubjectDiscussion  SubjectType = "Discussion"
	SubjectCheckSuite  SubjectType = "CheckSuite"
	SubjectUnknown     SubjectType = "Unknown"
)

// Icon returns a display icon for the subject type.
func (t SubjectType) Icon() string {
	switch t {
	case SubjectPullRequest:
		return "🔀"
	case SubjectIssue:
		return "🔴"
	case SubjectCommit:
		return "📝"
	case SubjectRelease:
		return "🏷️"
	case SubjectDiscussion:
		return "💬"
	case SubjectCheckSuite:
		return "✅"
	default:
		return "📌"
	}
}

// Reason is why the user received the notification.
type Reason string

const (
	ReasonApprovalRequested      Reason = "approval_requested"
	ReasonAssign                 Reason = "assign"
	ReasonAuthor                 Reason = "author"
	ReasonCIActivity             Reason = "ci_activity"
	ReasonComment                Reason = "comment"
	ReasonInvitation             Reason = "invitation"
	ReasonManual                 Reason = "manual"
	ReasonMemberFeatureRequested Reason = "member_feature_requested"
	ReasonMention                Reason = "mention"
	ReasonReviewRequested        Reason = "review_requested"
	ReasonSecurityAdvisoryCredit Reason = "security_advisory_credit"
	ReasonSecurityAlert          Reason = "security_alert"
	ReasonStateChange            Reason = "state_change"
	ReasonSubscribed             Reason = "subscribed"
	ReasonTeamMention            Reason = "team_mention"
)

// Icon returns a display icon for the reason.
func (r Reason) Icon() string {
	switch r {
	case ReasonMention:
		return "💬"
	case ReasonTeamMention:
		return "👥"
	case ReasonReviewRequested:
		return "👀"
	case ReasonAssign:
		return "👤"
	case ReasonAuthor:
		return "✍️"
	case ReasonComment:
		return "💭"
	case ReasonStateChange:
		return "🔄"
	case ReasonSubscribed:
		return "🔔"
	case ReasonManual:
		return "📌"
	case ReasonCIActivity:
		return "⚙️"
	case ReasonSecurityAlert:
		return "🚨"
	case ReasonApprovalRequested:
		return "✋"
	default:
		return "📫"
	}
}

// Label returns a human-readable label for the reason.
func (r Reason) Label() string {
	return strings.ReplaceAll(string(r), "_", " ")
}

// HTMLURL returns the web URL for viewing the notification subject.
func (n *Notification) HTMLURL() string {
	if n.Subject.Number > 0 {
		switch n.Subject.Type {
		case SubjectPullRequest:
			return fmt.Sprintf("%s/pull/%d", n.Repository.HTMLURL, n.Subject.Number)
		case SubjectIssue:
			return fmt.Sprintf("%s/issues/%d", n.Repository.HTMLURL, n.Subject.Number)
		case SubjectDiscussion:
			return fmt.Sprintf("%s/discussions/%d", n.Repository.HTMLURL, n.Subject.Number)
		}
	}
	return n.Repository.HTMLURL
}

// RelativeTime returns a human-readable relative time string.
func RelativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1m ago"
		}
		return fmt.Sprintf("%dm ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1h ago"
		}
		return fmt.Sprintf("%dh ago", h)
	case d < 30*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1d ago"
		}
		return fmt.Sprintf("%dd ago", days)
	default:
		months := int(d.Hours() / 24 / 30)
		if months == 1 {
			return "1mo ago"
		}
		return fmt.Sprintf("%dmo ago", months)
	}
}
