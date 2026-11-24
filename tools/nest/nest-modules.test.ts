import { mkdirSync, mkdtempSync, writeFileSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'
import { collect, source, type BrickRow } from './nest-modules'

const GARDEN = fileURLToPath(new URL('./fixtures/bricks/garden', import.meta.url))
const row = (name: string, dir: string, config: Record<string, unknown> = {}, on = true): BrickRow => ({ name, on, dir, config })

function brick(files: Record<string, string>): string {
  const dir = mkdtempSync(join(tmpdir(), 'nest-brick-'))
  for (const [f, body] of Object.entries(files)) {
    mkdirSync(join(dir, 'nest', f, '..'), { recursive: true })
    writeFileSync(join(dir, 'nest', f), body)
  }
  return dir
}

describe('nest modules (item_267)', () => {
  it('keeps the engine defaults with no brick and no config', () => {
    const m = collect([row('nest', '')])
    expect(m.errors).toEqual([])
    expect(m.config).toEqual({ theme: 'modern', hero: 'holotable', pages: [] })
    expect(source(m)).toContain('export const themes = ["modern","cyber"]')
  })

  it('finds the theme, hero, page and section of an on brick', () => {
    const m = collect([row('nest', '', { theme: 'moss', hero: 'garden', pages: ['home', 'garden', 'inbox'] }), row('garden', GARDEN)])
    expect(m.errors).toEqual([])
    expect(m.config).toEqual({ theme: 'moss', hero: 'garden', pages: ['home', 'garden', 'inbox'] })
    expect([m.themes, m.heroes, m.pages, m.sections].map((l) => l.map((f) => `${f.brick}:${f.name}`))).toEqual([['garden:moss'], ['garden:garden'], ['garden:garden'], ['garden:beds']])
  })

  it('skips off bricks', () => {
    const m = collect([row('garden', GARDEN, {}, false)])
    expect(m.themes.length + m.heroes.length + m.pages.length + m.sections.length).toBe(0)
  })

  it('names the brick for a clash, a theme without its block or a bad file name', () => {
    const dir = brick({ 'themes/cyber.css': "[data-theme='cyber'] {}", 'themes/fern.css': ':root {}', 'pages/inbox.tsx': '', 'pages/Bad.tsx': '' })
    expect(collect([row('weeds', dir)]).errors).toEqual([
      'brick weeds: theme cyber already exists',
      "brick weeds: theme fern.css has no [data-theme='fern'] block",
      'brick weeds: page file Bad.tsx: name must match /^[a-z][a-z0-9-]*$/',
      'brick weeds: page inbox already exists',
    ])
  })

  it('rejects unknown config values', () => {
    expect(collect([row('nest', '', { theme: 'nope', hero: 'nope', pages: ['home', 'nope'] })]).errors).toEqual([
      'crow.yaml modules.nest.theme: unknown theme "nope" (known: modern, cyber)',
      'crow.yaml modules.nest.hero: unknown hero "nope" (known: holotable, none)',
      'crow.yaml modules.nest.pages: unknown page "nope" (known: home, desk, ops, inbox, mrs, agenda, reports)',
    ])
    expect(collect([row('nest', '', { pages: 'home' })]).errors).toEqual(['crow.yaml modules.nest.pages: expects a list of page names'])
  })
})
