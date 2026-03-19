# Contributing to livemark

Thanks for your interest in contributing!

## Getting started

1. Fork the repo
2. Create a branch: `git checkout -b feat/your-feature`
3. Make your changes
4. Verify it compiles and runs: `GH_TOKEN=your_pat go run . --username yourusername`
5. Open a PR against `main`

## Guidelines

- Keep the zero external dependency policy — stdlib only
- New zones should follow the same pattern as `internal/sections/projects.go`
- Idempotency is a hard requirement: re-running with the same data must produce no diff

## Running locally

```bash
git clone https://github.com/jeziellopes/livemark
cd livemark
GH_TOKEN=your_pat go run . --username yourusername --readme path/to/README.md
```
