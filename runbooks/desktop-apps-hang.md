# Desktop apps will not open (GTK app freezes, no window)

Seen 2026-11-24: Thunar and LibreOffice never drew a window; new launches just handed off to the frozen instance.

1. Check the user daemons first: `ps -o pid,etime,time,rss -C gvfsd-metadata`. High RSS (GB) or tens of minutes of CPU means it is wedged.
2. `kill <pid>` (it restarts on demand), then `pkill -x thunar` and relaunch with `setsid thunar >/dev/null 2>&1 &`.
3. Only then look at image loaders (glycin, gdk-pixbuf). A hung `glycin-svg` child is a symptom: the loader decodes in ms and the app waits on gvfs metadata. Downgrading glycin did not help.
4. Suspect a recent bulk cleanup (Trash, `gio trash` of big folders) or a first `stow -R` that deployed old configs; check both before chasing packages.
5. Test with a PNG load under a 10 s timeout. Open a visible app only for the one relaunch the owner asked for.
