package api

import (
	"fmt"
	"strings"
	"time"
)

// SubjectDetail holds the rich details of a notification's subject
// fetched via REST API for the preview panel.
type SubjectDetail struct {
	Title  string
	Body   string
	State  string
	Author string
	Labels []string
	URL    string
}

// apiPR matches the GitHub REST API pull request response.
type apiPR struct {
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	State   string    `json:"state"`
	HTMLURL string    `json:"html_url"`
	Merged  bool      `json:"merged"`
	User    apiUser   `json:"user"`
	Labels  []apiLabel `json:"labels"`
}

// apiIssue matches the GitHub REST API issue response.
type apiIssue struct {
	Title   string     `json:"title"`
	Body    string     `json:"body"`
	State   string     `json:"state"`
	HTMLURL string     `json:"html_url"`
	User    apiUser    `json:"user"`
	Labels  []apiLabel `json:"labels"`
}

type apiUser struct {
	Login string `json:"login"`
}

type apiLabel struct {
	Name string `json:"name"`
}

// FetchSubjectDetail retrieves the full detail of a PR or issue for preview
// using the notification's subject API URL (not the HTML URL).
func (c *Client) FetchSubjectDetail(subjectAPIURL string, isPR bool) (*SubjectDetail, error) {
	if subjectAPIURL == "" {
		return nil, fmt.Errorf("no subject URL available")
	}

	// The subject URL is a full API URL like https://api.github.com/repos/owner/repo/pulls/123
	// go-gh REST client expects a path relative to the API base, so strip the prefix.
	path := subjectAPIURL
	if idx := strings.Index(path, "/repos/"); idx >= 0 {
		path = path[idx+1:] // "repos/owner/repo/pulls/123"
	}

	if isPR {
		return c.fetchPRDetail(path)
	}
	return c.fetchIssueDetail(path)
}

func (c *Client) fetchPRDetail(path string) (*SubjectDetail, error) {
	var pr apiPR
	if err := c.rest.Get(path, &pr); err != nil {
		return nil, fmt.Errorf("fetching PR detail: %w", err)
	}

	labels := make([]string, 0, len(pr.Labels))
	for _, l := range pr.Labels {
		labels = append(labels, l.Name)
	}

	state := pr.State
	if pr.Merged {
		state = "merged"
	}

	return &SubjectDetail{
		Title:  pr.Title,
		Body:   pr.Body,
		State:  state,
		Author: pr.User.Login,
		Labels: labels,
		URL:    pr.HTMLURL,
	}, nil
}

func (c *Client) fetchIssueDetail(path string) (*SubjectDetail, error) {
	var issue apiIssue
	if err := c.rest.Get(path, &issue); err != nil {
		return nil, fmt.Errorf("fetching issue detail: %w", err)
	}

	labels := make([]string, 0, len(issue.Labels))
	for _, l := range issue.Labels {
		labels = append(labels, l.Name)
	}

	return &SubjectDetail{
		Title:  issue.Title,
		Body:   issue.Body,
		State:  issue.State,
		Author: issue.User.Login,
		Labels: labels,
		URL:    issue.HTMLURL,
	}, nil
}

// apiRelease matches the GitHub REST API release response (minimal fields).
type apiRelease struct {
	Name    string    `json:"name"`
	TagName string    `json:"tag_name"`
	Body    string    `json:"body"`
	HTMLURL string    `json:"html_url"`
	Author  apiUser   `json:"author"`
	Draft   bool      `json:"draft"`
	CreatedAt time.Time `json:"created_at"`
}

// FetchReleaseDetail retrieves release details from the subject API URL.
func (c *Client) FetchReleaseDetail(subjectAPIURL string) (*SubjectDetail, error) {
	path := subjectAPIURL
	if idx := strings.Index(path, "/repos/"); idx >= 0 {
		path = path[idx+1:]
	}

	var rel apiRelease
	if err := c.rest.Get(path, &rel); err != nil {
		return nil, fmt.Errorf("fetching release detail: %w", err)
	}

	title := rel.Name
	if title == "" {
		title = rel.TagName
	}

	state := "published"
	if rel.Draft {
		state = "draft"
	}

	return &SubjectDetail{
		Title:  title,
		Body:   rel.Body,
		State:  state,
		Author: rel.Author.Login,
		URL:    rel.HTMLURL,
	}, nil
}
