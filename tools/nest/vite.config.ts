/// <reference types="vitest/config" />
import { defineConfig, type Plugin } from 'vite'
import react from '@vitejs/plugin-react'
import { readdirSync, readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'

// Shared core-contract records live outside this package. Read them from disk
// instead of importing across the package boundary.
const RECORDS = fileURLToPath(new URL('../../logics/specs/core_contract/fixtures/records/', import.meta.url))
const API = fileURLToPath(new URL('../../logics/specs/core_contract/fixtures/api/', import.meta.url))
const ID = 'virtual:core-records'
const API_ID = 'virtual:core-api'
// npm run dev proxies /api to the core over its unix socket, which the core
// trusts as the owner. The TCP port would need the tailnet identity header.
const CORE_SOCK = fileURLToPath(new URL('../../local/state/core.sock', import.meta.url))

function coreRecords(): Plugin {
  return {
    name: 'core-records',
    resolveId: (id) => (id === ID || id === API_ID ? '\0' + id : undefined),
    load(id) {
      if (id === '\0' + API_ID) {
        // API response fixtures as raw text, keyed by file name (tests only).
        const files = readdirSync(API).filter((n) => n.endsWith('.json'))
        return `export default ${JSON.stringify(Object.fromEntries(files.map((n) => [n, readFileSync(API + n, 'utf8')])))};`
      }
      if (id !== '\0' + ID) return
      const read = (n: string) => readFileSync(RECORDS + n, 'utf8')
      return `export const missions = ${read('missions.json')};\nexport const inbox = ${read('inbox.json')};`
    },
  }
}

export default defineConfig({
  plugins: [react(), coreRecords()],
  server: {
    host: true,
    proxy: { '/api': { target: { socketPath: CORE_SOCK } as never } },
  },
  // tailscale serve proxies the phone in with the tailnet hostname
  preview: { allowedHosts: ['.ts.net'] },
  build: { target: 'es2022' },
  test: { environment: 'node' },
})
