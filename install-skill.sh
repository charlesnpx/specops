#!/usr/bin/env sh
set -eu

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

if [ -x "$repo_dir/specops" ]; then
  exec "$repo_dir/specops" install-skill "$@"
fi

version=${SPECOPS_VERSION:-$(git -C "$repo_dir" describe --tags --exact-match --match 'v*' 2>/dev/null || true)}
version=${version#v}
version=${version:-0.1.3-dev}

exec go run -ldflags "-X main.version=$version" "$repo_dir/cmd/specops" install-skill "$@"
