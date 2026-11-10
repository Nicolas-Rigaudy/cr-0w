/// <reference types="vitest/config" />
import { defineConfig, type Plugin } from 'vite'
import react from '@vitejs/plugin-react'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'

// Shared core-contract records live outside this package. Read them from disk
// instead of importing across the package boundary.
const RECORDS = fileURLToPath(new URL('../../logics/specs/core_contract/fixtures/records/', import.meta.url))
const ID = 'virtual:core-records'

function coreRecords(): Plugin {
  return {
    name: 'core-records',
    resolveId: (id) => (id === ID ? '\0' + ID : undefined),
    load(id) {
      if (id !== '\0' + ID) return
      const read = (n: string) => readFileSync(RECORDS + n, 'utf8')
      return `export const missions = ${read('missions.json')};\nexport const inbox = ${read('inbox.json')};`
    },
  }
}

export default defineConfig({
  plugins: [react(), coreRecords()],
  server: { host: true },
  build: { target: 'es2022' },
  test: { environment: 'node' },
})
