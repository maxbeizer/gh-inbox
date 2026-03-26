package demo

import (
	"time"

	"github.com/maxbeizer/gh-inbox/internal/model"
)

// Notifications returns a realistic set of fake notifications for demo/recording.
func Notifications() []model.Notification {
	now := time.Now()
	return []model.Notification{
		{
			ID: "1", Unread: true,
			Reason: model.ReasonReviewRequested, UpdatedAt: now.Add(-2 * time.Minute),
			Subject: model.Subject{Title: "feat: add streaming support for large file uploads", Type: model.SubjectPullRequest, Number: 3421, URL: ""},
			Repository: model.Repository{FullName: "acme/api-gateway", HTMLURL: "https://github.com/acme/api-gateway"},
		},
		{
			ID: "2", Unread: true,
			Reason: model.ReasonMention, UpdatedAt: now.Add(-8 * time.Minute),
			Subject: model.Subject{Title: "Bug: dashboard widgets don't render on Safari 18", Type: model.SubjectIssue, Number: 892, URL: ""},
			Repository: model.Repository{FullName: "acme/web-console", HTMLURL: "https://github.com/acme/web-console"},
		},
		{
			ID: "3", Unread: true,
			Reason: model.ReasonAssign, UpdatedAt: now.Add(-15 * time.Minute),
			Subject: model.Subject{Title: "Migrate user auth from JWT to Paseto tokens", Type: model.SubjectIssue, Number: 1205, URL: ""},
			Repository: model.Repository{FullName: "acme/identity-service", HTMLURL: "https://github.com/acme/identity-service"},
		},
		{
			ID: "4", Unread: true,
			Reason: model.ReasonComment, UpdatedAt: now.Add(-32 * time.Minute),
			Subject: model.Subject{Title: "refactor: extract notification handler into separate package", Type: model.SubjectPullRequest, Number: 567, URL: ""},
			Repository: model.Repository{FullName: "acme/event-bus", HTMLURL: "https://github.com/acme/event-bus"},
		},
		{
			ID: "5", Unread: true,
			Reason: model.ReasonCIActivity, UpdatedAt: now.Add(-1 * time.Hour),
			Subject: model.Subject{Title: "CI: nightly integration test suite", Type: model.SubjectCheckSuite, Number: 0, URL: ""},
			Repository: model.Repository{FullName: "acme/platform", HTMLURL: "https://github.com/acme/platform"},
		},
		{
			ID: "6", Unread: true,
			Reason: model.ReasonTeamMention, UpdatedAt: now.Add(-2 * time.Hour),
			Subject: model.Subject{Title: "RFC: adopt OpenTelemetry for distributed tracing", Type: model.SubjectDiscussion, Number: 42, URL: ""},
			Repository: model.Repository{FullName: "acme/engineering-rfcs", HTMLURL: "https://github.com/acme/engineering-rfcs"},
		},
		{
			ID: "7", Unread: false,
			Reason: model.ReasonStateChange, UpdatedAt: now.Add(-3 * time.Hour),
			Subject: model.Subject{Title: "fix: race condition in connection pool recycling", Type: model.SubjectPullRequest, Number: 1102, URL: ""},
			Repository: model.Repository{FullName: "acme/api-gateway", HTMLURL: "https://github.com/acme/api-gateway"},
		},
		{
			ID: "8", Unread: false,
			Reason: model.ReasonSubscribed, UpdatedAt: now.Add(-5 * time.Hour),
			Subject: model.Subject{Title: "v2.4.0 — Improved caching and rate limit handling", Type: model.SubjectRelease, Number: 0, URL: ""},
			Repository: model.Repository{FullName: "acme/go-sdk", HTMLURL: "https://github.com/acme/go-sdk"},
		},
		{
			ID: "9", Unread: false,
			Reason: model.ReasonAuthor, UpdatedAt: now.Add(-8 * time.Hour),
			Subject: model.Subject{Title: "docs: add runbook for on-call database failover", Type: model.SubjectPullRequest, Number: 89, URL: ""},
			Repository: model.Repository{FullName: "acme/ops-handbook", HTMLURL: "https://github.com/acme/ops-handbook"},
		},
		{
			ID: "10", Unread: false,
			Reason: model.ReasonSecurityAlert, UpdatedAt: now.Add(-12 * time.Hour),
			Subject: model.Subject{Title: "Dependabot: bump golang.org/x/crypto from 0.21 to 0.23", Type: model.SubjectPullRequest, Number: 4410, URL: ""},
			Repository: model.Repository{FullName: "acme/api-gateway", HTMLURL: "https://github.com/acme/api-gateway"},
		},
		{
			ID: "11", Unread: false,
			Reason: model.ReasonManual, UpdatedAt: now.Add(-1 * 24 * time.Hour),
			Subject: model.Subject{Title: "Tracking: Q2 performance optimization initiative", Type: model.SubjectIssue, Number: 300, URL: ""},
			Repository: model.Repository{FullName: "acme/platform", HTMLURL: "https://github.com/acme/platform"},
		},
		{
			ID: "12", Unread: false,
			Reason: model.ReasonApprovalRequested, UpdatedAt: now.Add(-2 * 24 * time.Hour),
			Subject: model.Subject{Title: "Deploy: canary rollout for search indexer v3", Type: model.SubjectPullRequest, Number: 78, URL: ""},
			Repository: model.Repository{FullName: "acme/deploy-infra", HTMLURL: "https://github.com/acme/deploy-infra"},
		},
	}
}
