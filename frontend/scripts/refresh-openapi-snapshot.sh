#!/usr/bin/env bash
set -euo pipefail

# Regenerate the OpenAPI snapshot from Go types (no server or database needed).
# Must be run from the frontend/ directory.

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
FRONTEND_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKEND_DIR="$FRONTEND_DIR/../backend"
SNAPSHOT="$FRONTEND_DIR/openapi-snapshot.json"

if ! command -v go &>/dev/null; then
	echo "error: 'go' is required on PATH to regenerate the OpenAPI snapshot" >&2
	exit 1
fi

cd "$BACKEND_DIR"
go run ./cmd/openapi-gen >"$SNAPSHOT"
echo "OpenAPI snapshot written to frontend/openapi-snapshot.json"
