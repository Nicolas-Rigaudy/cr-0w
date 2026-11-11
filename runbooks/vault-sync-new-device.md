# Add a device to the vault's Syncthing mesh

Use when a new machine (desktop, laptop, phone) needs the Obsidian vault.

1. Check the hub machine's config first: `grep -oE '<device id="[^"]+" name="[^"]*"|<folder id="[^"]+" label="[^"]*"' ~/.local/state/syncthing/config.xml` (use the real home path under cdx). A folder only shows up on the new device once a peer shares it with that device.
2. Install Syncthing on the new device (Windows: SyncTrayzor or the official Windows setup; Linux: package + `systemctl --user enable --now syncthing`).
3. Pair it directly with the hub machine (Add Remote Device with its ID, accept on both sides), then tick the vault in the folder's Sharing tab. Pairing only through a phone works but makes the phone a required relay.
4. Accept the folder on the new device into an empty path (any drive). Wait for the first sync before opening it in Obsidian.
5. Copy the ignore patterns from the hub machine's `.stignore` (they do not sync): workspace files, noisy plugin data, `.*.crow-tmp`.
6. Use one sync tool per vault per device: leave Obsidian Sync off where Syncthing runs.
7. After the first sync, list `*.sync-conflict-*` files, diff each against its main file, and delete stale copies once the main file is confirmed current (deleting propagates to all peers).
