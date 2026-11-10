# Nest on the phone (tailnet PWA install)

Install the Nest PWA on an Android phone while Nest is only reachable on the tailnet.

## Serve
1. `cd tools/nest && npm run build`
2. `npx vite preview --host 127.0.0.1 --port 4173 --strictPort` (background). `vite.config.ts` allows `.ts.net` hosts in preview.
3. `tailscale serve --bg 4173` and keep Funnel off (no login in Nest; Funnel requests carry no tailnet identity).
   - First time: "Serve is not enabled" prints an admin link; enable Serve + HTTPS there.
   - "Access denied": the owner runs `sudo tailscale set --operator=<user>` once in a real terminal (the `!` prompt cannot read a sudo password).
   - The very first HTTPS request can stall while the cert is issued.
4. `tailscale serve status` prints the `https://<host>.<tailnet>.ts.net/` URL.

## Install
1. Phone: Tailscale connected, open the URL in Chrome, ⋮ → Install app.
2. Chrome builds the WebAPK on Google's side fine (private host is not a problem). The Play Store then parks the install at `SCHEDULED` while the Tailscale VPN is up.
3. Fix: switch Tailscale off on the phone until the app appears, then back on. Every earlier tap that was parked installs too: remove the extra copies.

## Debug an install
- Host: `android-tools` (adb). Phone: Developer options → USB debugging, accept the RSA prompt.
- `adb forward tcp:9222 localabstract:chrome_devtools_remote`, filter `curl -s 127.0.0.1:9222/json` to the Nest URL only (the list holds every open tab).
- CDP over a Node `WebSocket`: `Page.getInstallabilityErrors`, `Page.getAppManifest`.
- `adb logcat -v time | grep -iE "webapk|Finsky"` while tapping Install shows minting and Play queue states.
- `adb forward --remove-all` when done.
