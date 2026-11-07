#!/bin/sh
# Runs inside a clean container: install, version, update, data check, restore.
set -e
apk add -q --no-cache bash curl git >/dev/null
adduser -D u
cp -r /smoke/data /srv/data && chown -R u /srv/data
# shellcheck disable=SC2016  # expands inside the user shell
su u -s /bin/sh -c '
set -e; cd; export CROW_RELEASE_URL=file:///smoke/rel CROW_API_URL=file:///smoke/api PATH=$HOME/.local/bin:$PATH
git config --global --add safe.directory "*"
echo "== install"; sh /smoke/install.sh
echo "== which"; command -v crow
echo "== version"; crow version
echo "local before:"; ls ~/.config/crow/local; echo "x" > ~/.config/crow/local/state/keep.txt
S1=$(find ~/.config/crow -type f -exec sha256sum {} + | sort | sha256sum)
echo "== update --check"; crow update --check --to 0.2.0; readlink ~/.local/share/crow/current
echo "== update"; crow update --to 0.2.0; crow version | head -1
S2=$(find ~/.config/crow -type f -exec sha256sum {} + | sort | sha256sum)
[ "$S1" = "$S2" ] && echo "data dir unchanged byte for byte"
echo "== restore on a fresh user dir"; export HOME=$HOME/fresh; mkdir -p $HOME; export PATH=$HOME/.local/bin:$PATH
sh /smoke/install.sh --restore /srv/data; ls ~/.config/crow/local; crow version
'
