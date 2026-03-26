package api

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/maxbeizer/gh-inbox/internal/model"
)

// apiNotification matches the GitHub REST API notification response shape.
type apiNotification struct {
	ID              string            `json:"id"`
	Unread          bool              `json:"unread"`
	Reason          string            `json:"reason"`
	UpdatedAt       time.Time         `json:"updated_at"`
	LastReadAt      *time.Time        `json:"last_read_at"`
	URL             string            `json:"url"`
	SubscriptionURL string            `json:"subscription_url"`
	Subject         apiSubject        `json:"subject"`
	Repository      apiRepository     `json:"repository"`
}

type apiSubject struct {
	Title            string `json:"title"`
	URL              string `json:"url"`
	LatestCommentURL string `json:"latest_comment_url"`
	Type             string `json:"type"`
}

type apiRepository struct {
	FullName string   `json:"full_name"`
	Private  bool     `json:"private"`
	HTMLURL  string   `json:"html_url"`
	Owner    apiOwner `json:"owner"`
	Name     string   `json:"name"`
}

type apiOwner struct {
	Login string `json:"login"`
}

var numberFromURL = regexp.MustCompile(`/(\d+)$`)

// parseNumber extracts the issue/PR number from an API URL like
// https://api.github.com/repos/owner/repo/pulls/123
func parseNumber(url string) int {
	m := numberFromURL.FindStringSubmatch(url)
	if len(m) < 2 {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return n
}

func toSubjectType(s string) model.SubjectType {
	switch s {
	case "PullRequest":
		return model.SubjectPullRequest
	case "Issue":
		return model.SubjectIssue
	case "Commit":
		return model.SubjectCommit
	case "Release":
		return model.SubjectRelease
	case "Discussion":
		return model.SubjectDiscussion
	case "CheckSuite":
		return model.SubjectCheckSuite
	default:
		return model.SubjectUnknown
	}
}

func toNotification(a apiNotification) model.Notification {
	return model.Notification{
		ID:              a.ID,
		Unread:          a.Unread,
		Reason:          model.Reason(a.Reason),
		UpdatedAt:       a.UpdatedAt,
		LastReadAt:      a.LastReadAt,
		URL:             a.URL,
		SubscriptionURL: a.SubscriptionURL,
		Subject: model.Subject{
			Title:            a.Subject.Title,
			URL:              a.Subject.URL,
			LatestCommentURL: a.Subject.LatestCommentURL,
			Type:             toSubjectType(a.Subject.Type),
			Number:           parseNumber(a.Subject.URL),
		},
		Repository: model.Repository{
			FullName: a.Repository.FullName,
			Owner:    a.Repository.Owner.Login,
			Name:     a.Repository.Name,
			Private:  a.Repository.Private,
			HTMLURL:  a.Repository.HTMLURL,
		},
	}
}

const (
	perPage  = 50
	maxPages = 10 // Cap at 500 notifications to avoid rate limiting
)

// ListNotifications fetches notifications from the GitHub API.
// Fetches up to maxPages pages (500 notifications) to avoid excessive API calls.
func (c *Client) ListNotifications(all, participating bool) ([]model.Notification, error) {
	var allNotifications []model.Notification

	for page := 1; page <= maxPages; page++ {
		var raw []apiNotification
		path := fmt.Sprintf("notifications?all=%t&participating=%t&per_page=%d&page=%d",
			all, participating, perPage, page)

		err := c.rest.Get(path, &raw)
		if err != nil {
			return nil, fmt.Errorf("fetching notifications page %d: %w", page, err)
		}

		for _, a := range raw {
			allNotifications = append(allNotifications, toNotification(a))
		}

		if len(raw) < perPage {
			break
		}
	}

	return allNotifications, nil
}

// MarkThreadRead marks a single notification thread as read.
func (c *Client) MarkThreadRead(threadID string) error {
	path := fmt.Sprintf("notifications/threads/%s", threadID)
	return c.rest.Patch(path, nil, nil)
}

// MarkThreadDone marks a notification thread as done (dismisses it).
func (c *Client) MarkThreadDone(threadID string) error {
	path := fmt.Sprintf("notifications/threads/%s", threadID)
	return c.rest.Delete(path, nil)
}

// Unsubscribe removes the subscription for a notification thread.
func (c *Client) Unsubscribe(threadID string) error {
	path := fmt.Sprintf("notifications/threads/%s/subscription", threadID)
	return c.rest.Delete(path, nil)
}
