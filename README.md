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

Run `livemark` and it rewrites each zone with live data from GitHub — featured projects ranked by stars and recency, and recent OSS contributions. Everything outside the markers stays untouched.

## Quickstart

### 1. Add zone markers to your README

```markdown
<!-- PROJECTS_START -->
<!-- PROJECTS_END -->

<!-- OSS_START -->
<!-- OSS_END -->
```

### 2. Run livemark

```bash
GH_TOKEN=your_pat GH_GITHUB_USERNAME=yourusername go run github.com/jeziellopes/livemark@latest
```

Or install it:

```bash
go install github.com/jeziellopes/livemark@latest
GH_TOKEN=your_pat livemark --username yourusername
```

### 3. Automate with GitHub Actions

Add this workflow to your profile repo (`.github/workflows/update-readme.yml`):

```yaml
name: Update README

on:
  schedule:
    - cron: "0 6 * * *"
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

## Flags

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `--username` | `GH_GITHUB_USERNAME` | _(required)_ | GitHub username |
| `--readme` | — | `README.md` | Path to README file |
| `--projects` | — | `4` | Number of featured projects |
| `--oss` | — | `5` | Number of OSS contributions |

## Zone reference

| Zone | Content |
|------|---------|
| `PROJECTS` | Top public repos ranked by stars + recency, rendered as prose bullets |
| `OSS` | Recent merged/open PRs on external public repos |

Zones are idempotent — if data hasn't changed since the last run, the file is not modified.

## Privacy

- Private repos are **never exposed** — the OSS zone only shows PRs on public repos (checked via the full PR object, not just the event payload)
- The `GH_TOKEN` PAT is used only for API authentication and rate limit headroom

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
