package github

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v58/github"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

// Client wraps GitHub REST and GraphQL clients
type Client struct {
	Rest    *github.Client
	GraphQL *graphql.Client
	Token   string
}

// NewClient creates a new GitHub client with authentication
func NewClient(token string) *Client {
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	restClient := github.NewClient(tc)
	graphqlClient := graphql.NewClient("https://api.github.com/graphql", tc)

	return &Client{
		Rest:    restClient,
		GraphQL: graphqlClient,
		Token:   token,
	}
}

// Repository represents a GitHub repository
type Repository struct {
	Owner struct {
		Login string
	}
	Name        string
	FullName    string `graphql:"nameWithOwner"`
	Description string
	IsPrivate   bool
	IsFork      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Languages   struct {
		Nodes []struct {
			Name string
		}
	} `graphql:"languages(first: 1)"`
	StargazerCount int
	ForkCount      int
}

// GetRepositories fetches repositories for an owner (user or organization) using GraphQL
func (c *Client) GetRepositories(ctx context.Context, owner string, limit int) ([]Repository, error) {
	// Use repositoryOwner which works for both users and organizations
	var query struct {
		RepositoryOwner struct {
			Repositories struct {
				Nodes []Repository
			} `graphql:"repositories(first: $limit, orderBy: {field: UPDATED_AT, direction: DESC})"`
		} `graphql:"repositoryOwner(login: $owner)"`
	}

	variables := map[string]interface{}{
		"owner": graphql.String(owner),
		"limit": graphql.Int(limit),
	}

	err := c.GraphQL.Query(ctx, &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories for owner '%s': %w", owner, err)
	}

	return query.RepositoryOwner.Repositories.Nodes, nil
}

// Commit represents a commit with author information
type Commit struct {
	OID    string
	Author struct {
		Name  string
		Email string
		Date  time.Time
	}
	Additions int
	Deletions int
	Message   string
}

// GetCommits fetches commits for a repository within a date range using REST API
func (c *Client) GetCommits(ctx context.Context, owner, repo string, since, until time.Time) ([]Commit, error) {
	var allCommits []Commit
	page := 1
	perPage := 100

	for {
		// Use REST API for better date filtering
		opts := &github.CommitsListOptions{
			Since: since,
			Until: until,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		}

		commits, resp, err := c.Rest.Repositories.ListCommits(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch commits: %w", err)
		}

		// Convert REST API commits to our internal format
		for _, commit := range commits {
			if commit.Commit == nil || commit.Commit.Author == nil {
				continue
			}

			c := Commit{
				OID:     commit.GetSHA(),
				Message: commit.Commit.GetMessage(),
			}

			if commit.Commit.Author != nil {
				c.Author.Name = commit.Commit.Author.GetName()
				c.Author.Email = commit.Commit.Author.GetEmail()
				if commit.Commit.Author.Date != nil {
					c.Author.Date = commit.Commit.Author.Date.Time
				}
			}

			if commit.Stats != nil {
				c.Additions = commit.Stats.GetAdditions()
				c.Deletions = commit.Stats.GetDeletions()
			}

			allCommits = append(allCommits, c)
		}

		// Check if there are more pages
		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage

		// Safety limit to avoid infinite loops
		if len(allCommits) > 10000 {
			break
		}
	}

	return allCommits, nil
}

// PullRequest represents a pull request
type PullRequest struct {
	Number    int
	Title     string
	State     string
	IsDraft   bool
	CreatedAt time.Time
	MergedAt  *time.Time
	ClosedAt  *time.Time
	Author    struct {
		Login string
	}
	Reviews struct {
		Nodes []struct {
			State  string
			Author struct {
				Login string
			}
			CreatedAt time.Time
		}
	} `graphql:"reviews(first: 100)"`
	Comments struct {
		TotalCount int
		Nodes      []struct {
			Author struct {
				Login string
			}
			CreatedAt time.Time
		}
	} `graphql:"comments(first: 100)"`
}

// GetPullRequests fetches pull requests for a repository within a date range
func (c *Client) GetPullRequests(ctx context.Context, owner, repo string, since, until time.Time) ([]PullRequest, error) {
	var allPRs []PullRequest
	hasNextPage := true
	cursor := (*string)(nil)

	for hasNextPage {
		var query struct {
			Repository struct {
				PullRequests struct {
					PageInfo struct {
						HasNextPage bool
						EndCursor   string
					}
					Nodes []PullRequest
				} `graphql:"pullRequests(first: 100, orderBy: {field: CREATED_AT, direction: DESC}, after: $cursor)"`
			} `graphql:"repository(owner: $owner, name: $repo)"`
		}

		variables := map[string]interface{}{
			"owner":  graphql.String(owner),
			"repo":   graphql.String(repo),
			"cursor": (*graphql.String)(cursor),
		}

		err := c.GraphQL.Query(ctx, &query, variables)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch pull requests: %w", err)
		}

		prs := query.Repository.PullRequests.Nodes

		// Filter by date range
		var filteredPRs []PullRequest
		for _, pr := range prs {
			if pr.CreatedAt.After(since) && pr.CreatedAt.Before(until) {
				filteredPRs = append(filteredPRs, pr)
			}
		}

		allPRs = append(allPRs, filteredPRs...)

		hasNextPage = query.Repository.PullRequests.PageInfo.HasNextPage
		if hasNextPage {
			cursor = &query.Repository.PullRequests.PageInfo.EndCursor
		}

		// If we've gone past our date range, stop fetching
		if len(prs) > 0 && prs[len(prs)-1].CreatedAt.Before(since) {
			break
		}
	}

	return allPRs, nil
}

// Issue represents a GitHub issue
type Issue struct {
	Number    int
	Title     string
	State     string
	CreatedAt time.Time
	Author    struct {
		Login string
	}
	Comments struct {
		TotalCount int
		Nodes      []struct {
			Author struct {
				Login string
			}
			CreatedAt time.Time
		}
	} `graphql:"comments(first: 100)"`
}

// GetIssues fetches issues for a repository within a date range
func (c *Client) GetIssues(ctx context.Context, owner, repo string, since, until time.Time) ([]Issue, error) {
	var allIssues []Issue
	hasNextPage := true
	cursor := (*string)(nil)

	for hasNextPage {
		var query struct {
			Repository struct {
				Issues struct {
					PageInfo struct {
						HasNextPage bool
						EndCursor   string
					}
					Nodes []Issue
				} `graphql:"issues(first: 100, orderBy: {field: CREATED_AT, direction: DESC}, after: $cursor)"`
			} `graphql:"repository(owner: $owner, name: $repo)"`
		}

		variables := map[string]interface{}{
			"owner":  graphql.String(owner),
			"repo":   graphql.String(repo),
			"cursor": (*graphql.String)(cursor),
		}

		err := c.GraphQL.Query(ctx, &query, variables)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch issues: %w", err)
		}

		issues := query.Repository.Issues.Nodes

		// Filter by date range
		var filteredIssues []Issue
		for _, issue := range issues {
			if issue.CreatedAt.After(since) && issue.CreatedAt.Before(until) {
				filteredIssues = append(filteredIssues, issue)
			}
		}

		allIssues = append(allIssues, filteredIssues...)

		hasNextPage = query.Repository.Issues.PageInfo.HasNextPage
		if hasNextPage {
			cursor = &query.Repository.Issues.PageInfo.EndCursor
		}

		// If we've gone past our date range, stop fetching
		if len(issues) > 0 && issues[len(issues)-1].CreatedAt.Before(since) {
			break
		}
	}

	return allIssues, nil
}
