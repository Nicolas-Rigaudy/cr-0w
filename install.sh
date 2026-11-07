#!/bin/sh
# install.sh: install a tagged crow release into a user prefix (Linux, macOS, WSL).
#
#   sh install.sh [--version X.Y.Z] [--prefix DIR] [--bin-dir DIR] [--restore GIT_URL] [--engine-only]
#
# Downloads crow-<version>.tar.gz and SHA256SUMS from the GitHub release, verifies the checksum,
# unpacks into <prefix>/versions/<version>, points <prefix>/current at it and links <bin-dir>/crow.
# Creates the data dir (CROW_HOME, default $XDG_CONFIG_HOME/crow) with local/ from the example
# files; existing data is never overwritten. --restore clones the owner's private data repo into
# <data>/local and lists the steps git cannot carry (keyring tokens, cdx accounts, the vault).
# --engine-only skips the data dir (crow update uses it).
#
# Environment: CROW_REPO (owner/name on GitHub), CROW_RELEASE_URL (base holding download/v<ver>/),
# CROW_API_URL (base answering /releases/latest), CROW_PREFIX, CROW_BIN_DIR, CROW_HOME.
set -eu

REPO="${CROW_REPO:-Nicolas-Rigaudy/cr-0w}"
RELEASE_URL="${CROW_RELEASE_URL:-https://github.com/$REPO/releases}"
API_URL="${CROW_API_URL:-https://api.github.com/repos/$REPO}"
PREFIX="${CROW_PREFIX:-${XDG_DATA_HOME:-$HOME/.local/share}/crow}"
BIN_DIR="${CROW_BIN_DIR:-$HOME/.local/bin}"
DATA="${CROW_HOME:-${XDG_CONFIG_HOME:-$HOME/.config}/crow}"
VERSION=""
RESTORE=""
ENGINE_ONLY=""

die() { echo "install.sh: $*" >&2; exit 1; }
say() { echo "crow: $*"; }

while [ $# -gt 0 ]; do
  case "$1" in
    --version) VERSION="${2:?--version needs a value}"; shift 2 ;;
    --prefix)  PREFIX="${2:?--prefix needs a value}"; shift 2 ;;
    --bin-dir) BIN_DIR="${2:?--bin-dir needs a value}"; shift 2 ;;
    --restore) RESTORE="${2:?--restore needs a git url}"; shift 2 ;;
    --engine-only) ENGINE_ONLY=1; shift ;;
    -h|--help) sed -n '2,14p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *) die "unknown option '$1' (--help)" ;;
  esac
done

fetch() {  # fetch <url> <file>
  if command -v curl >/dev/null 2>&1; then curl -fsSL "$1" -o "$2"
  elif command -v wget >/dev/null 2>&1; then wget -q "$1" -O "$2"
  else die "needs curl or wget"; fi
}

sha256() {
  if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | cut -d' ' -f1
  elif command -v shasum >/dev/null 2>&1; then shasum -a 256 "$1" | cut -d' ' -f1
  else die "needs sha256sum or shasum"; fi
}

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

if [ -z "$VERSION" ]; then
  fetch "$API_URL/releases/latest" "$tmp/latest.json" || die "cannot read the latest release from $API_URL"
  VERSION=$(sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"v\{0,1\}\([^"]*\)".*/\1/p' "$tmp/latest.json" | head -n 1)
  [ -n "$VERSION" ] || die "no tag_name in the latest release"
fi
VERSION="${VERSION#v}"
asset="crow-$VERSION.tar.gz"
base="$RELEASE_URL/download/v$VERSION"

say "installing $VERSION into $PREFIX"
fetch "$base/$asset" "$tmp/$asset" || die "cannot download $base/$asset"
fetch "$base/SHA256SUMS" "$tmp/SHA256SUMS" || die "cannot download $base/SHA256SUMS"
want=$(awk -v a="$asset" '{ f = $2; sub(/^\*/, "", f); if (f == a) print $1 }' "$tmp/SHA256SUMS")
[ -n "$want" ] || die "SHA256SUMS has no line for $asset"
got=$(sha256 "$tmp/$asset")
[ "$want" = "$got" ] || die "checksum mismatch for $asset (expected $want, got $got); nothing installed"

mkdir -p "$tmp/x"
tar -xzf "$tmp/$asset" -C "$tmp/x"
src="$tmp/x/crow-$VERSION"
[ -x "$src/bin/crow" ] || die "$asset has no crow-$VERSION/bin/crow"
[ -e "$src/local" ] && die "$asset contains local/; refusing a release that carries owner data"
printf 'method=script\nrepo=%s\nprefix=%s\nbin_dir=%s\n' "$REPO" "$PREFIX" "$BIN_DIR" > "$src/.crow-install"

mkdir -p "$PREFIX/versions" "$BIN_DIR"
dest="$PREFIX/versions/$VERSION"
rm -rf "$dest.new"
mv "$src" "$dest.new"
rm -rf "$dest"
mv "$dest.new" "$dest"
ln -sfn "versions/$VERSION" "$PREFIX/current"
ln -sf "$PREFIX/current/bin/crow" "$BIN_DIR/crow"
say "linked $BIN_DIR/crow"

# Data dir: created from the examples, never overwritten.
if [ -n "$ENGINE_ONLY" ]; then
  :
elif [ -n "$RESTORE" ]; then
  if [ -d "$DATA/local" ] && [ -n "$(ls -A "$DATA/local" 2>/dev/null)" ]; then
    die "$DATA/local already has content; move it away before --restore"
  fi
  command -v git >/dev/null 2>&1 || die "--restore needs git"
  mkdir -p "$DATA"
  rm -rf "$DATA/local"
  git clone -q "$RESTORE" "$DATA/local" || die "cannot clone $RESTORE"
  say "restored $DATA/local from $RESTORE"
else
  mkdir -p "$DATA/local/state"
  for ex in "$dest"/docs/*.example.yaml; do
    [ -e "$ex" ] || continue
    name=$(basename "$ex" .example.yaml).yaml
    [ -e "$DATA/local/$name" ] || cp "$ex" "$DATA/local/$name"
  done
  say "data dir $DATA (existing files kept)"
fi

case ":$PATH:" in *":$BIN_DIR:"*) ;; *) say "add $BIN_DIR to PATH" ;; esac
[ -n "$ENGINE_ONLY" ] || [ "$DATA" = "${XDG_CONFIG_HOME:-$HOME/.config}/crow" ] || say "set CROW_HOME=$DATA in your shell profile"

if [ -n "$RESTORE" ]; then
  echo
  echo "Left to do by hand (git does not carry these, see runbooks/install-and-update.md):"
  if [ -f "$DATA/local/tokens.yaml" ]; then
    sed -n 's/^[[:space:]]*-[[:space:]]*service:[[:space:]]*\([^[:space:]#]*\).*/\1/p' "$DATA/local/tokens.yaml" \
      | while read -r svc; do echo "  crow secret set $svc"; done
  fi
  if [ -f "$DATA/local/accounts.yaml" ]; then
    sed -n 's/.*name:[[:space:]]*\([A-Za-z0-9_.-]*\).*/\1/p' "$DATA/local/accounts.yaml" \
      | while read -r acc; do echo "  cdx add $acc && cdx login $acc"; done
  fi
  echo "  clone or sync the vault, then check vault.path in $DATA/local/projects.yaml"
  echo "  crow doctor"
fi
say "done: crow version"
