#!/usr/bin/env sh
set -eu

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

if [ -x "$repo_dir/specops" ]; then
  exec "$repo_dir/specops" install-skill "$@"
fi

if command -v specops >/dev/null 2>&1; then
  exec specops install-skill "$@"
fi

exec go run "$repo_dir/cmd/specops" install-skill "$@"
