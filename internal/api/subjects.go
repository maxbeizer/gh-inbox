package api

import (
	"fmt"
)

// SubjectDetail holds the rich details of a notification's subject
// fetched via GraphQL for the preview panel.
type SubjectDetail struct {
	Title  string
	Body   string
	State  string
	Author string
	Labels []string
	URL    string
}

// prQuery is the GraphQL query for pull request details.
type prQuery struct {
	Resource struct {
		PullRequest struct {
			Title  string
			Body   string
			State  string
			URL    string
			Author struct {
				Login string
			}
			Labels struct {
				Nodes []struct {
					Name string
				}
			} `graphql:"labels(first: 20)"`
		} `graphql:"... on PullRequest"`
	} `graphql:"resource(url: $url)"`
}

// issueQuery is the GraphQL query for issue details.
type issueQuery struct {
	Resource struct {
		Issue struct {
			Title  string
			Body   string
			State  string
			URL    string
			Author struct {
				Login string
			}
			Labels struct {
				Nodes []struct {
					Name string
				}
			} `graphql:"labels(first: 20)"`
		} `graphql:"... on Issue"`
	} `graphql:"resource(url: $url)"`
}

// FetchSubjectDetail retrieves the full detail of a PR or issue for preview.
func (c *Client) FetchSubjectDetail(htmlURL string, isPR bool) (*SubjectDetail, error) {
	if isPR {
		return c.fetchPRDetail(htmlURL)
	}
	return c.fetchIssueDetail(htmlURL)
}

func (c *Client) fetchPRDetail(url string) (*SubjectDetail, error) {
	var q prQuery
	vars := map[string]interface{}{
		"url": url,
	}
	if err := c.gql.Query("PRDetail", &q, vars); err != nil {
		return nil, fmt.Errorf("fetching PR detail: %w", err)
	}

	pr := q.Resource.PullRequest
	labels := make([]string, 0, len(pr.Labels.Nodes))
	for _, l := range pr.Labels.Nodes {
		labels = append(labels, l.Name)
	}

	return &SubjectDetail{
		Title:  pr.Title,
		Body:   pr.Body,
		State:  pr.State,
		Author: pr.Author.Login,
		Labels: labels,
		URL:    pr.URL,
	}, nil
}

func (c *Client) fetchIssueDetail(url string) (*SubjectDetail, error) {
	var q issueQuery
	vars := map[string]interface{}{
		"url": url,
	}
	if err := c.gql.Query("IssueDetail", &q, vars); err != nil {
		return nil, fmt.Errorf("fetching issue detail: %w", err)
	}

	issue := q.Resource.Issue
	labels := make([]string, 0, len(issue.Labels.Nodes))
	for _, l := range issue.Labels.Nodes {
		labels = append(labels, l.Name)
	}

	return &SubjectDetail{
		Title:  issue.Title,
		Body:   issue.Body,
		State:  issue.State,
		Author: issue.Author.Login,
		Labels: labels,
		URL:    issue.URL,
	}, nil
}
