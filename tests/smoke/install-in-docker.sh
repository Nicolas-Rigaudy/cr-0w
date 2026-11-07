#!/usr/bin/env bash
# Install round trip in a clean container: install.sh, crow version, crow update, --restore.
#   tests/smoke/install-in-docker.sh [image]   (default python:3.13-alpine)
set -euo pipefail
root="$(cd "$(dirname "$(readlink -f "$0")")/../.." && pwd)"
work="$root/.smoke"
rm -rf "$work"; mkdir -p "$work"
trap 'rm -rf "$work"' EXIT
python3 "$root/tests/smoke/build_release.py" "$work"
cp "$root/install.sh" "$root/tests/smoke/in_container.sh" "$work/"
docker run --rm -v "$work:/smoke:ro" "${1:-python:3.13-alpine}" sh /smoke/in_container.sh
