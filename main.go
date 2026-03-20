// livemark — keep your GitHub profile README alive.
//
// Rewrites HTML-comment-delimited zones in your README.md with live
// data from the GitHub API: featured projects and OSS contributions.
//
// Usage:
//
//	GH_TOKEN=<pat> livemark [flags]
//
// Flags:
//
//	--username   GitHub username (default: GH_GITHUB_USERNAME env, required)
//	--readme     path to README.md (default: README.md)
//	--projects   number of featured projects to show (default: 4)
//	--oss        number of OSS contributions to show (default: 5)
//	--dry-run    preview generated content in livemark.preview.md without modifying README
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	gh "github.com/jeziellopes/livemark/internal/github"
	"github.com/jeziellopes/livemark/internal/readme"
	"github.com/jeziellopes/livemark/internal/sections"
)

func main() {
	username := flag.String("username", os.Getenv("GH_GITHUB_USERNAME"), "GitHub username")
	readmePath := flag.String("readme", "README.md", "path to README.md")
	projectsCount := flag.Int("projects", 4, "number of featured projects")
	ossCount := flag.Int("oss", 5, "number of OSS contributions")
	dryRun := flag.Bool("dry-run", false, "preview generated content in livemark.preview.md without modifying README")
	flag.Parse()

	token := resolveToken()
	if token == "" {
		fmt.Fprintln(os.Stderr, "error: no GitHub token found.\nSet GH_TOKEN or authenticate with: gh auth login")
		os.Exit(1)
	}
	if *username == "" {
		fmt.Fprintln(os.Stderr, "error: --username flag or GH_GITHUB_USERNAME env variable is required")
		flag.Usage()
		os.Exit(1)
	}

	client := gh.New(token)

	projects, err := sections.BuildProjects(client, *username, *projectsCount)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error fetching projects:", err)
		os.Exit(1)
	}

	oss, err := sections.BuildOSS(client, *username, *ossCount)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error fetching OSS contributions:", err)
		os.Exit(1)
	}

	if *dryRun {
		const previewPath = "livemark.preview.md"
		zones := []readme.PreviewZone{
			{Name: "PROJECTS", Content: projects},
			{Name: "OSS", Content: oss},
		}
		if err := readme.WritePreview(previewPath, *username, zones); err != nil {
			fmt.Fprintln(os.Stderr, "error writing preview:", err)
			os.Exit(1)
		}
		fmt.Printf("Preview written to %s\n", previewPath)
		return
	}

	changed := false
	for _, zone := range []struct {
		name    string
		content string
	}{
		{"PROJECTS", projects},
		{"OSS", oss},
	} {
		ok, err := readme.Rewrite(*readmePath, zone.name, zone.content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error rewriting zone %s: %v\n", zone.name, err)
			os.Exit(1)
		}
		if ok {
			changed = true
		}
	}

	if changed {
		fmt.Println("README.md updated.")
	} else {
		fmt.Println("README.md is already up to date.")
	}
}

// resolveToken returns a GitHub token from the first available source:
// GH_TOKEN env → GITHUB_TOKEN env → gh auth token (gh≥2.37) →
// gh auth status --show-token (gh<2.37).
func resolveToken() string {
	if t := os.Getenv("GH_TOKEN"); t != "" {
		return t
	}
	if t := os.Getenv("GITHUB_TOKEN"); t != "" {
		return t
	}
	// gh auth token — available in gh >= 2.37
	if out, err := exec.Command("gh", "auth", "token").Output(); err == nil {
		if t := strings.TrimSpace(string(out)); t != "" {
			return t
		}
	}
	// gh auth status --show-token — fallback for older gh versions
	if out, err := exec.Command("gh", "auth", "status", "--show-token").CombinedOutput(); err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			if _, after, ok := strings.Cut(line, "Token: "); ok {
				if t := strings.TrimSpace(after); t != "" {
					return t
				}
			}
		}
	}
	return ""
}
