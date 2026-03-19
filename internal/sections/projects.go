package sections

import (
	"fmt"
	"sort"
	"strings"
	"time"

	gh "github.com/jeziellopes/livemark/internal/github"
)

// BuildProjects returns the markdown content for the PROJECTS zone.
// It fetches the user's public, non-fork repos and ranks them by stars + recency,
// returning up to count entries as prose bullet points.
func BuildProjects(client *gh.Client, username string, count int) (string, error) {
	var repos []gh.Repo
	if err := client.Get(fmt.Sprintf("/users/%s/repos?sort=updated&per_page=100&type=owner", username), &repos); err != nil {
		return "", err
	}

	var public []gh.Repo
	for _, r := range repos {
		if !r.Private && !r.Fork && r.Name != username {
			public = append(public, r)
		}
	}
	sort.Slice(public, func(i, j int) bool {
		return repoScore(public[i]) > repoScore(public[j])
	})
	if len(public) > count {
		public = public[:count]
	}

	var sb strings.Builder
	sb.WriteString("### What I've shipped lately\n\n")
	for _, r := range public {
		desc := strings.TrimRight(r.Description, ".")
		stars := ""
		if r.StargazersCount > 0 {
			stars = fmt.Sprintf(" ⭐ %d", r.StargazersCount)
		}
		fmt.Fprintf(&sb, "- **[%s](%s)**%s — %s.\n", r.Name, r.HTMLURL, stars, desc)
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

func repoScore(r gh.Repo) float64 {
	daysAgo := time.Since(r.UpdatedAt).Hours() / 24
	recency := max(0, 365-daysAgo) / 365
	return float64(r.StargazersCount)*10 + recency
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
