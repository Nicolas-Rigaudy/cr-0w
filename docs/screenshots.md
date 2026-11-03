# Regenerating the README screenshots

The shots use the sample fixtures, so they contain no real data.

```sh
S=$(mktemp -d)
(cd tools/crow-watch && SHOTS=$S go test -run TestShots .)
for v in agents usage; do
  bin/ansi-shot $S/$v-laptop.ans $S/$v.png
  magick $S/$v.png -colors 128 PNG8:docs/img/watch-$v.png
done
```

Colours and font come from your kitty theme. See `runbooks/tui-render.md`.

## Deck live shots

`TestLiveShots` renders the deck Home from the real sources (registry, open threads, accounts, runs) at 235x62 and 118x60, in all and work modes. Output contains private data: keep it under `local/`.

```sh
D=local/state/shots/deck/live
(cd tools/crow-deck && LIVESHOTS=$PWD/../../$D go test -run TestLiveShots -count=1 .)
for f in $D/*.ans; do bin/ansi-shot $f ${f%.ans}.png; done
```
