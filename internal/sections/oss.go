package sections

import (
	"fmt"
	"sort"
	"strings"

	gh "github.com/jeziellopes/livemark/internal/github"
)

type contribution struct {
	title, prURL, repoName, repoURL string
	merged                          bool
	createdAt                       string
}

// BuildOSS returns the markdown content for the OSS zone.
// It uses the GraphQL pullRequests API to find merged PRs authored by the user
// on external public repos, sorted by creation date (newest first).
func BuildOSS(client *gh.Client, username string, count int) (string, error) {
	nodes, err := client.FetchAuthoredPRs(username, count)
	if err != nil {
		return "", err
	}

	var contribs []contribution
	for _, node := range nodes {
		if !node.Merged {
			continue
		}
		contribs = append(contribs, contribution{
			title:     node.Title,
			prURL:     node.URL,
			repoName:  node.Repository.NameWithOwner,
			repoURL:   node.Repository.URL,
			merged:    true,
			createdAt: node.CreatedAt,
		})
	}

	sortContribs(contribs)

	if len(contribs) == 0 {
		return "### Recent Open Source contributions\n\n_No recent merged contributions found._", nil
	}

	var sb strings.Builder
	sb.WriteString("### Recent Open Source contributions\n\n")
	for _, c := range contribs {
		fmt.Fprintf(&sb, "- ✅ Merged **[%s](%s)** into [%s](%s)\n", c.title, c.prURL, c.repoName, c.repoURL)
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

// sortContribs sorts contributions in-place by creation date descending (newest first).
func sortContribs(cs []contribution) {
	sort.SliceStable(cs, func(i, j int) bool {
		return cs[i].createdAt > cs[j].createdAt
	})
}
