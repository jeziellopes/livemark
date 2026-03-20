package sections

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	gh "github.com/jeziellopes/livemark/internal/github"
)

// BuildProjects returns the markdown content for the PROJECTS zone.
// Pinned repos (user-curated) are listed first; remaining slots are filled
// by public non-fork repos ranked by a composite score of stars, recency,
// and repo size.
func BuildProjects(client *gh.Client, username string, count int) (string, error) {
	// Fetch pinned repos via GraphQL (best signal — user-curated).
	pinned, err := client.FetchPinnedRepos(username)
	if err != nil {
		// Non-fatal: fall back to scored-only if GraphQL fails.
		pinned = nil
	}

	// Fetch all public repos via REST.
	var repos []gh.Repo
	if err := client.Get(fmt.Sprintf("/users/%s/repos?sort=updated&per_page=100&type=owner", username), &repos); err != nil {
		return "", err
	}

	// Build a set of pinned repo names for deduplication.
	pinnedNames := make(map[string]bool, len(pinned))
	for _, p := range pinned {
		pinnedNames[p.Name] = true
	}

	// Filter REST repos: public, non-fork, non-self, not already pinned.
	var public []gh.Repo
	for _, r := range repos {
		if !r.Private && !r.Fork && r.Name != username && !pinnedNames[r.Name] {
			public = append(public, r)
		}
	}
	sort.Slice(public, func(i, j int) bool {
		return repoScore(public[i]) > repoScore(public[j])
	})

	var sb strings.Builder
	sb.WriteString("### What I've shipped lately\n\n")

	written := 0

	// Pinned repos first (in user-defined order).
	for _, p := range pinned {
		if written >= count {
			break
		}
		desc := strings.TrimRight(p.Description, ".")
		stars := ""
		if p.StargazerCount > 0 {
			stars = fmt.Sprintf(" ⭐ %d", p.StargazerCount)
		}
		fmt.Fprintf(&sb, "- **[%s](%s)**%s — %s.\n", p.Name, p.URL, stars, desc)
		written++
	}

	// Fill remaining slots with top-scored repos.
	for _, r := range public {
		if written >= count {
			break
		}
		desc := strings.TrimRight(r.Description, ".")
		stars := ""
		if r.StargazersCount > 0 {
			stars = fmt.Sprintf(" ⭐ %d", r.StargazersCount)
		}
		fmt.Fprintf(&sb, "- **[%s](%s)**%s — %s.\n", r.Name, r.HTMLURL, stars, desc)
		written++
	}

	return strings.TrimRight(sb.String(), "\n"), nil
}

// repoScore ranks repos by a composite of log-scaled stars, recency of last
// push, and log-scaled repo size. Log-scaling prevents a handful of
// high-star repos from drowning out newer, actively developed projects.
func repoScore(r gh.Repo) float64 {
	pushedAt := r.PushedAt
	if pushedAt.IsZero() {
		pushedAt = r.UpdatedAt
	}
	daysAgo := time.Since(pushedAt).Hours() / 24
	recency := math.Max(0, 365-daysAgo) / 365
	return math.Log1p(float64(r.StargazersCount))*10 +
		recency +
		math.Log1p(float64(r.Size))*2
}
