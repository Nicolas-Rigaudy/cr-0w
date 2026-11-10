// Stub service worker: makes Nest installable. No offline cache yet; every
// request goes to the network, and /api/* must never be cached (item_167).
// Web Push arrives in a later slice.
self.addEventListener('install', () => self.skipWaiting())
self.addEventListener('activate', (event) => event.waitUntil(self.clients.claim()))
self.addEventListener('fetch', () => {})
