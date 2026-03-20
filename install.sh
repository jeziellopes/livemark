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

# If GOBIN is not in PATH, add it to the appropriate shell profile
if ! echo "$PATH" | tr ':' '\n' | grep -qx "$GOBIN"; then
  # Detect shell profile
  case "${SHELL:-}" in
    */zsh)  PROFILE="$HOME/.zshrc" ;;
    */bash) PROFILE="$HOME/.bashrc" ;;
    *)      PROFILE="$HOME/.profile" ;;
  esac

  EXPORT_LINE="export PATH=\"\$PATH:$GOBIN\""

  if ! grep -qF "$GOBIN" "$PROFILE" 2>/dev/null; then
    echo "" >> "$PROFILE"
    echo "# Added by livemark installer" >> "$PROFILE"
    echo "$EXPORT_LINE" >> "$PROFILE"
    echo ""
    echo "  ✅ Added PATH export to $PROFILE"
    echo "  Run this to apply now:"
    echo ""
    echo "    source $PROFILE"
    echo ""
  else
    echo ""
    echo "  ℹ️  $GOBIN already referenced in $PROFILE (restart your shell to apply)"
    echo ""
  fi
fi

echo "Done. Run: $BINARY --help"
