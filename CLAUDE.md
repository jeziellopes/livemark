# {project-name}

{One sentence: what this does and who it's for.}

## Stack

- Runtime: {e.g., Go 1.23 / Node.js 22 / Python 3.12}
- Framework: {e.g., Cobra / NestJS / FastAPI — or "none"}
- Database: {e.g., PostgreSQL} _(if applicable)_
- Cache: {e.g., Redis} _(if applicable)_
- Queue: {e.g., BullMQ} _(if applicable)_

## Commands

- Install: `{install-command}`
- Dev: `{dev-command}`
- Test: `{test-command}`
- Test single: `{test-single-command}`
- Lint: `{lint-command}`
- Build: `{build-command}`

## Architecture

{Brief description of the high-level architecture and data flow.}

- `{src-dir}/` — {description}
- `{src-dir}/` — {description}

## Key Patterns

- {Pattern 1 — e.g., "Repository pattern for all database access"}
- {Pattern 2 — e.g., "Errors wrapped with context at every boundary"}

## Observability

- Logging: {library + format — e.g., "slog, structured JSON" or "pino, structured JSON with correlation IDs"}
- Metrics: {library or service — e.g., "Prometheus"}
- Tracing: {library or service — e.g., "OpenTelemetry"}

## Deployment

- Environment: {devops-platform}
- CI/CD: {ci-cd-platform}
- Secrets: {how secrets are managed — e.g., "env vars, never committed"}

## Conventions

- {Only things that differ from framework defaults}
- {Add project-specific patterns here}

## Secrets

- Manager: {secrets-manager}
- Load locally: `{secrets-run-cmd}`
- See `docs/secrets.md` for setup instructions

## Specs

- Feature specs live in `specs/` with acceptance criteria
- Use `/spec` to generate, `/red` → `/green` → `/refactor` for TDD workflow
