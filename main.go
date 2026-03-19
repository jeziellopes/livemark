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
package main

import (
	"flag"
	"fmt"
	"os"

	gh "github.com/jeziellopes/livemark/internal/github"
	"github.com/jeziellopes/livemark/internal/readme"
	"github.com/jeziellopes/livemark/internal/sections"
)

func main() {
	username := flag.String("username", os.Getenv("GH_GITHUB_USERNAME"), "GitHub username")
	readmePath := flag.String("readme", "README.md", "path to README.md")
	projectsCount := flag.Int("projects", 4, "number of featured projects")
	ossCount := flag.Int("oss", 5, "number of OSS contributions")
	flag.Parse()

	token := os.Getenv("GH_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "error: GH_TOKEN environment variable is required")
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
