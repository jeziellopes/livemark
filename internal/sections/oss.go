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
	createdAt                               string
}

// BuildOSS returns the markdown content for the OSS zone.
// It uses the GraphQL pullRequests API to find PRs authored by the user on
// external public repos, sorted merged first then open then closed.
func BuildOSS(client *gh.Client, username string, count int) (string, error) {
	nodes, err := client.FetchAuthoredPRs(username, count)
	if err != nil {
		return "", err
	}

	var contribs []contribution
	for _, node := range nodes {
		var status string
		switch {
		case node.Merged:
			status = "✅ Merged"
		case node.State == "OPEN":
			status = "🔄 Open"
		default:
			status = "❌ Closed"
		}
		contribs = append(contribs, contribution{
			title:     node.Title,
			prURL:     node.URL,
			repoName:  node.Repository.NameWithOwner,
			repoURL:   node.Repository.URL,
			status:    status,
			merged:    node.Merged,
			open:      node.State == "OPEN",
			createdAt: node.CreatedAt,
		})
	}

	sortContribs(contribs)

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
// Within each status group, newer contributions appear first (descending by createdAt).
func sortContribs(cs []contribution) {
	sort.SliceStable(cs, func(i, j int) bool {
		ri, rj := rank(cs[i]), rank(cs[j])
		if ri != rj {
			return ri < rj
		}
		// Same rank: sort by createdAt descending (newer first)
		return cs[i].createdAt > cs[j].createdAt
	})
}
