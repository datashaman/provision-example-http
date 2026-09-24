#!/usr/bin/env bash
set -euo pipefail

repo_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
target_arch="${1:-amd64}"
case "$target_arch" in
  amd64|arm64) ;;
  *) echo "usage: $0 [amd64|arm64]" >&2; exit 2 ;;
esac
output_dir="$repo_dir/dist"
mkdir -p "$output_dir"
build_dir="$(mktemp -d)"
trap 'rm -rf "$build_dir"' EXIT
(cd "$repo_dir" && GOOS=linux GOARCH="$target_arch" CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags='-buildid=' -o "$build_dir/provision-example-http" .)
chmod 0755 "$build_dir/provision-example-http"
archive="$output_dir/provision-example-http-linux-$target_arch.tar.gz"
(cd "$repo_dir" && go run ./cmd/package --input "$build_dir/provision-example-http" --output "$archive")
printf 'Artifact: %s\n' "$archive"
if command -v sha256sum >/dev/null 2>&1; then
  sha256sum "$archive"
else
  shasum -a 256 "$archive"
fi
