# livemark

Keep your GitHub profile README alive — self-updating zones powered by the GitHub API. Zero deps, written in Go.

## How it works

Add HTML comment markers to your `README.md` to define dynamic zones:

```markdown
<!-- PROJECTS_START -->
<!-- PROJECTS_END -->

<!-- OSS_START -->
<!-- OSS_END -->
```

Run `livemark` and it rewrites each zone with live data from GitHub:

- **PROJECTS** — your pinned repos appear first (in your defined order), followed by public repos ranked by stars, recency, and size. Repos not pushed to in 2+ years are excluded.
- **OSS** — recent pull requests on external repos, sorted merged → open → closed.

Everything outside the markers stays untouched. Zones are idempotent — if data hasn't changed, the file is not modified.

## Quickstart

### 1. Add zone markers to your README

```markdown
<!-- PROJECTS_START -->
<!-- PROJECTS_END -->

<!-- OSS_START -->
<!-- OSS_END -->
```

### 2. Install livemark

```bash
curl -fsSL https://raw.githubusercontent.com/jeziellopes/livemark/main/install.sh | bash
```

Or with Go directly:

```bash
go install github.com/jeziellopes/livemark@latest
```

### 3. Run livemark

```bash
GH_TOKEN=your_pat livemark --username yourusername
```

### 4. Automate with GitHub Actions

Add this workflow to your profile repo (`.github/workflows/update-readme.yml`).
It runs on a **daily schedule** and also on a **`livemark-release` dispatch** event so your
profile updates immediately whenever a new livemark version is released.

```yaml
name: Update README

on:
  repository_dispatch:
    types: [livemark-release]
  schedule:
    - cron: "0 2 * * *"  # 23:00 UTC-3 (end of work day)
  workflow_dispatch:

jobs:
  update:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: stable
      - name: Run livemark
        env:
          GH_TOKEN: ${{ secrets.GH_TOKEN }}
        run: go run github.com/jeziellopes/livemark@latest --username yourusername
      - name: Commit and push if changed
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git add README.md
          git diff --cached --quiet || git commit -m "chore: update profile README [skip ci]"
          git push
```

> **Required secret:** `GH_TOKEN` — Personal Access Token with `repo` and `read:user` scopes.

To enable instant updates on release, also add `PROFILE_DISPATCH_TOKEN` to the **livemark
repo** secrets — a PAT with `repo` scope on your profile repo. livemark's release workflow
fires the dispatch automatically after each tagged release.

## Local development

If you have the [GitHub CLI](https://cli.github.com/) installed and authenticated, no token setup is needed:

```bash
gh auth login
livemark --username yourname --dry-run
# Preview written to livemark.preview.md
```

Open `livemark.preview.md` in any markdown viewer to see exactly what livemark would inject into your README — without modifying it.

Token resolution order (first found wins):

| Priority | Source |
|----------|--------|
| 1 | `GH_TOKEN` env var |
| 2 | `GITHUB_TOKEN` env var |
| 3 | `gh auth token` (gh ≥ 2.37) |
| 4 | `gh auth status --show-token` (gh < 2.37) |

`livemark.preview.md` is gitignored by default and safe to generate locally at any time.

## Flags

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `--username` | `GH_GITHUB_USERNAME` | _(required)_ | GitHub username |
| `--readme` | — | `README.md` | Path to README file |
| `--projects` | — | `4` | Number of featured projects |
| `--oss` | — | `5` | Number of OSS contributions |
| `--dry-run` | — | `false` | Write preview to `livemark.preview.md` instead of modifying README |

## Zone reference

| Zone | Content |
|------|---------|
| `PROJECTS` | Pinned repos first, then top public repos ranked by stars × recency × size. Repos inactive 2+ years are excluded. |
| `OSS` | Recent PRs on external public repos, sorted merged → open → closed |

## Privacy

- Private repos are **never exposed** — the OSS zone only shows PRs on public repos (checked via the full PR object, not just the event payload)
- The `GH_TOKEN` PAT is used only for API authentication and rate limit headroom

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
