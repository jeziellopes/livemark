package github

import (
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
