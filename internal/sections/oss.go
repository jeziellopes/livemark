package sections

import (
	"fmt"
	"sort"
	"strings"

	gh "github.com/jeziellopes/livemark/internal/github"
)

type contribution struct {
	title, prURL, repoName, repoURL, status string
	merged                                  bool
	open                                    bool
}

// BuildOSS returns the markdown content for the OSS zone.
// It scans the user's recent public events for PullRequestEvents on external
// repos, paginating up to maxEventPages pages to surface contributions that
// may be buried in a busy activity stream. Results are sorted merged first
// then open then closed, and up to count entries are returned.
func BuildOSS(client *gh.Client, username string, count int) (string, error) {
	const maxEventPages = 3

	seen := map[string]bool{}
	var contribs []contribution

	for page := 1; page <= maxEventPages; page++ {
		var events []gh.Event
		url := fmt.Sprintf("/users/%s/events?per_page=100&page=%d", username, page)
		if err := client.Get(url, &events); err != nil {
			return "", err
		}

		for _, e := range events {
			if e.Type != "PullRequestEvent" {
				continue
			}
			parts := strings.SplitN(e.Repo.Name, "/", 2)
			if len(parts) < 2 || parts[0] == username {
				continue
			}

			apiURL := e.Payload.PullRequest.URL
			if apiURL == "" || seen[apiURL] {
				continue
			}
			seen[apiURL] = true

			// Fetch full PR details — the event payload is minimal
			var pr gh.PullRequest
			apiPath := strings.TrimPrefix(apiURL, gh.APIBase)
			if err := client.Get(apiPath, &pr); err != nil {
				continue
			}

			// Never show PRs from private repos
			if pr.Base.Repo.Private {
				continue
			}

			var status string
			switch {
			case pr.MergedAt != nil:
				status = "✅ Merged"
			case pr.State == "open":
				status = "🔄 Open"
			default:
				status = "❌ Closed"
			}

			prURL := pr.HTMLURL
			if prURL == "" {
				prURL = fmt.Sprintf("https://github.com/%s/pull/%d", e.Repo.Name, e.Payload.PullRequest.Number)
			}
			title := pr.Title
			if title == "" {
				title = e.Repo.Name
			}

			contribs = append(contribs, contribution{
				title:    title,
				prURL:    prURL,
				repoName: e.Repo.Name,
				repoURL:  "https://github.com/" + e.Repo.Name,
				status:   status,
				merged:   pr.MergedAt != nil,
				open:     pr.State == "open",
			})
		}

		// Stop early if this was the last page or we have enough contributions.
		if len(events) < 100 || len(contribs) >= count {
			break
		}
	}

	// Sort: merged first, then open, then closed.
	sortContribs(contribs)

	if len(contribs) > count {
		contribs = contribs[:count]
	}

	if len(contribs) == 0 {
		return "### Recent OSS\n\n_No recent external contributions found._", nil
	}

	var sb strings.Builder
	sb.WriteString("### Recent OSS\n\n")
	for _, c := range contribs {
		fmt.Fprintf(&sb, "- %s **[%s](%s)** into [%s](%s)\n", c.status, c.title, c.prURL, c.repoName, c.repoURL)
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

// rank returns a sort key: 0 = merged, 1 = open, 2 = closed.
func rank(c contribution) int {
	if c.merged {
		return 0
	}
	if c.open {
		return 1
	}
	return 2
}

// sortContribs sorts contributions in-place: merged first, then open, then closed.
func sortContribs(cs []contribution) {
	sort.SliceStable(cs, func(i, j int) bool {
		return rank(cs[i]) < rank(cs[j])
	})
}
