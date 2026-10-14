#!/usr/bin/env bash
# Enforce the cgo boundary: only internal/joltc may import "C".
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

violations="$(grep -rIl --include='*.go' -E '^import "C"$' . | grep -v '^\./internal/joltc/' || true)"

if [[ -n "$violations" ]]; then
  echo "error: import \"C\" outside internal/joltc:" >&2
  echo "$violations" >&2
  exit 1
fi

echo "cgo boundary OK"
