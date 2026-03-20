#!/usr/bin/env bash
set -euo pipefail

BINARY="livemark"
MODULE="github.com/jeziellopes/livemark"

# Check Go is installed
if ! command -v go &>/dev/null; then
  echo "error: Go is not installed. Get it at https://go.dev/dl/" >&2
  exit 1
fi

echo "Installing $BINARY..."
go install "${MODULE}@latest"

# Resolve where go install puts binaries
GOBIN="$(go env GOPATH)/bin"
if [ -n "$(go env GOBIN)" ]; then
  GOBIN="$(go env GOBIN)"
fi

echo "Installed to $GOBIN/$BINARY"

# Warn if the bin dir is not in PATH
if ! echo "$PATH" | tr ':' '\n' | grep -qx "$GOBIN"; then
  echo ""
  echo "  ⚠️  $GOBIN is not in your PATH."
  echo "  Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
  echo ""
  echo "    export PATH=\"\$PATH:$GOBIN\""
  echo ""
fi

echo "Done. Run: $BINARY --help"
