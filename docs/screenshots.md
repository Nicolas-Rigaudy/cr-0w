# Regenerating the README screenshots

The shots use the sample fixtures, so they contain no real data.

```sh
S=$(mktemp -d)
(cd tools/crow-deck && SHOTS=$S go test -run TestShots .)
for v in Agents Usage; do
  bin/ansi-shot $S/lazytf-$v-wide.ans $S/$v.png
  magick $S/$v.png -colors 128 PNG8:docs/img/watch-${v,,}.png
done
```

Colours come from your kitty theme when it exists, else the bundled `share/ansi-shot/default.conf` (`--theme FILE` for any kitty-format theme). See `runbooks/tui-render.md`.

## Deck live shots

`TestLiveShots` renders the deck Home (all and work modes) and the Agents page from the real sources (registry, open threads, accounts, runs, usage) at 235x62 and 118x60. Output contains private data: keep it under `local/`.

```sh
D=local/state/shots/deck/live
(cd tools/crow-deck && LIVESHOTS=$PWD/../../$D go test -run TestLiveShots -count=1 .)
for f in $D/*.ans; do bin/ansi-shot $f ${f%.ans}.png; done
```
