# Secrets Management Guide

Secrets (API keys, tokens, passwords, connection strings) must never be committed to git.
This project uses **gitleaks** to block accidental commits + a secrets manager for local dev and CI.

## Quick Start

### 1. Activate the pre-commit hook

```bash
pip install pre-commit   # or: brew install pre-commit
pre-commit install
```

After this, every `git commit` automatically scans for secrets. First-time setup is one minute.

---

## Secrets Managers (pick one)

### Infisical (self-hosted) — recommended if you have your own infra

Open source (MIT). Run your own instance. CLI syncs secrets to local dev and CI.

```bash
# Install CLI
brew install infisical/get-cli/infisical

# Point to your instance (add to shell profile)
export INFISICAL_API_URL=https://your-infisical.example.com

# Login and link project
infisical login
infisical init

# Run commands with secrets injected
infisical run -- npm start
infisical run -- go run .

# GitHub Actions: use the Infisical action
# https://infisical.com/docs/integrations/cicd/githubactions
```

`.infisical.json` (in this repo) stores the project workspace ID — no secrets, safe to commit.

---

### Doppler — best SaaS option (no infra required)

```bash
brew install dopplerhq/cli/doppler
doppler login
doppler setup           # link this project

doppler run -- npm start
doppler run -- go run .
```

GitHub Actions: add `DOPPLER_TOKEN` as a repo secret, use `dopplerhq/cli-action`.

---

### SOPS + age — best offline/air-gap option

Encrypts your `.env` file so it can be committed to git safely.

```bash
brew install sops age

# Generate a key
age-keygen -o ~/.config/sops/age/keys.txt

# Encrypt
sops --age $(cat ~/.config/sops/age/keys.txt | grep public | awk '{print $4}') \
     --encrypt .env > .env.enc

# Use
sops exec-env .env.enc -- npm start
```

---

### direnv — minimal option (no secrets management, just scoped loading)

```bash
brew install direnv
echo 'eval "$(direnv hook bash)"' >> ~/.bashrc  # or zsh
direnv allow .
```

`.envrc` (in this repo) runs `dotenv` to load `.env` automatically when you `cd` into the project.

---

## What NOT to do

- ❌ Never commit `.env` (it's in `.gitignore`)
- ❌ Never paste tokens in chat, PRs, or issue comments
- ❌ Never hardcode secrets in source files
- ❌ Never log secrets (even in debug mode)
- ✅ Rotate tokens immediately if exposed
