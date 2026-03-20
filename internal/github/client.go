package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const APIBase = "https://api.github.com"

// Repo represents a GitHub repository.
type Repo struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	HTMLURL         string    `json:"html_url"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	Fork            bool      `json:"fork"`
	Private         bool      `json:"private"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Size            int       `json:"size"`
}

// PinnedRepo represents a pinned repository from the GitHub GraphQL API.
type PinnedRepo struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	StargazerCount  int       `json:"stargazerCount"`
	PushedAt        time.Time `json:"pushedAt"`
	Size            int       `json:"diskUsage"`
}

// Event represents a GitHub public event.
type Event struct {
	Type    string    `json:"type"`
	Repo    EventRepo `json:"repo"`
	Payload PRPayload `json:"payload"`
}

// EventRepo is the repo reference inside an event.
type EventRepo struct {
	Name string `json:"name"`
}

// PRPayload is the payload for PullRequestEvent.
type PRPayload struct {
	PullRequest PRRef `json:"pull_request"`
}

// PRRef is the minimal PR reference returned in event payloads.
type PRRef struct {
	URL    string `json:"url"`
	Number int    `json:"number"`
}

// PullRequest is the full pull request object from the PRs API.
type PullRequest struct {
	Title    string  `json:"title"`
	HTMLURL  string  `json:"html_url"`
	State    string  `json:"state"`
	MergedAt *string `json:"merged_at"`
	Base     PRBase  `json:"base"`
}

// PRBase holds the base branch info of a PR.
type PRBase struct {
	Repo BaseRepo `json:"repo"`
}

// BaseRepo holds the repo info from a PR base.
type BaseRepo struct {
	Private bool `json:"private"`
}

// Client is an authenticated GitHub API client.
type Client struct {
	token      string
	httpClient *http.Client
}

// New creates a new GitHub API client.
func New(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// Get fetches a GitHub API path and decodes the JSON response into out.
func (c *Client) Get(path string, out any) error {
	req, err := http.NewRequest("GET", APIBase+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("GitHub API %s returned %d: %s", path, resp.StatusCode, body)
	}
	return json.Unmarshal(body, out)
}

// Post executes a GitHub GraphQL query and decodes the response data into out.
func (c *Client) Post(query string, variables map[string]any, out any) error {
	payload, err := json.Marshal(map[string]any{"query": query, "variables": variables})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", APIBase+"/graphql", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("GitHub GraphQL returned %d: %s", resp.StatusCode, body)
	}

	var wrapper struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return err
	}
	if len(wrapper.Errors) > 0 {
		return fmt.Errorf("GitHub GraphQL error: %s", wrapper.Errors[0].Message)
	}
	return json.Unmarshal(wrapper.Data, out)
}

// PullRequestNode represents a PR returned by the GraphQL pullRequests query.
type PullRequestNode struct {
	Title    string  `json:"title"`
	URL      string  `json:"url"`
	State    string  `json:"state"`
	Merged   bool    `json:"merged"`
	MergedAt *string `json:"mergedAt"`
	Repository struct {
		NameWithOwner string `json:"nameWithOwner"`
		URL           string `json:"url"`
		IsPrivate     bool   `json:"isPrivate"`
		Owner         struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
}

// FetchAuthoredPRs returns pull requests authored by the user, excluding
// their own repos and private repos. It requests up to 100 PRs (the GraphQL
// max) to provide headroom when many results are private or self-owned.
func (c *Client) FetchAuthoredPRs(username string, limit int) ([]PullRequestNode, error) {
	const query = `
	query($login: String!, $count: Int!) {
		user(login: $login) {
			pullRequests(first: $count, states: [OPEN, MERGED, CLOSED],
			             orderBy: {field: CREATED_AT, direction: DESC}) {
				nodes {
					title
					url
					state
					merged
					mergedAt
					repository {
						nameWithOwner
						url
						isPrivate
						owner { login }
					}
				}
			}
		}
	}`

	// Always request 100 (the GraphQL max) to have plenty of headroom after
	// filtering out private repos and self-owned repos.
	const fetchCount = 100

	var data struct {
		User struct {
			PullRequests struct {
				Nodes []PullRequestNode `json:"nodes"`
			} `json:"pullRequests"`
		} `json:"user"`
	}
	if err := c.Post(query, map[string]any{"login": username, "count": fetchCount}, &data); err != nil {
		return nil, err
	}

	var result []PullRequestNode
	for _, node := range data.User.PullRequests.Nodes {
		if node.Repository.IsPrivate || node.Repository.Owner.Login == username {
			continue
		}
		result = append(result, node)
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

// FetchPinnedRepos returns the user's pinned repositories via the GraphQL API.
func (c *Client) FetchPinnedRepos(username string) ([]PinnedRepo, error) {
	const query = `
	query($login: String!) {
		user(login: $login) {
			pinnedItems(first: 6, types: [REPOSITORY]) {
				nodes {
					... on Repository {
						name
						description
						url
						stargazerCount
						pushedAt
						diskUsage
					}
				}
			}
		}
	}`

	var data struct {
		User struct {
			PinnedItems struct {
				Nodes []PinnedRepo `json:"nodes"`
			} `json:"pinnedItems"`
		} `json:"user"`
	}
	if err := c.Post(query, map[string]any{"login": username}, &data); err != nil {
		return nil, err
	}
	return data.User.PinnedItems.Nodes, nil
}
