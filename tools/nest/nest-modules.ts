// Build-time Nest modules (item_267): themes, hero, pages and sections from the user's bricks.
// `virtual:nest-modules` is generated from `crow-modules json` (on bricks and the nest brick
// config in crow.yaml) and each on brick's nest/ folder. Nothing loads at runtime.
import { execFileSync } from 'node:child_process'
import { readdirSync, readFileSync } from 'node:fs'
import { join, relative } from 'node:path'
import type { Plugin } from 'vite'

export const ENGINE_THEMES = ['modern', 'cyber']
export const ENGINE_HEROES = ['holotable', 'none']
export const ENGINE_PAGES = ['home', 'desk', 'ops', 'inbox', 'mrs', 'agenda', 'reports']
const NAME = /^[a-z][a-z0-9-]*$/

export interface BrickRow { name: string; on: boolean; dir: string; config: Record<string, unknown> }
export interface NestFile { brick: string; name: string; file: string }
export interface NestModules {
  config: { theme: string; hero: string; pages: string[] }
  themes: NestFile[]
  heroes: NestFile[]
  pages: NestFile[]
  sections: NestFile[]
  errors: string[]
}

function files(dir: string, ext: RegExp): string[] {
  try {
    return readdirSync(dir).filter((n) => ext.test(n)).sort()
  } catch {
    return []
  }
}

// Reads the bricks as crow-modules reports them. Every problem is an error: the build fails on it.
export function collect(rows: BrickRow[], read: (f: string) => string = (f) => readFileSync(f, 'utf8')): NestModules {
  const errors: string[] = []
  const nest = rows.find((r) => r.name === 'nest')?.config ?? {}
  const out: NestModules = { config: { theme: 'modern', hero: 'holotable', pages: [] }, themes: [], heroes: [], pages: [], sections: [], errors }
  for (const b of rows) {
    if (!b.on || !b.dir) continue
    const root = join(b.dir, 'nest')
    for (const n of files(join(root, 'themes'), /\.css$/)) {
      const name = n.slice(0, -4)
      const file = join(root, 'themes', n)
      if (!NAME.test(name)) errors.push(`brick ${b.name}: theme file ${n}: name must match ${NAME}`)
      else if (ENGINE_THEMES.includes(name) || out.themes.some((t) => t.name === name)) errors.push(`brick ${b.name}: theme ${name} already exists`)
      else if (!read(file).includes(`[data-theme='${name}']`)) errors.push(`brick ${b.name}: theme ${n} has no [data-theme='${name}'] block`)
      else out.themes.push({ brick: b.name, name, file })
    }
    const hero = files(join(root, 'hero'), /^index\.tsx?$/)[0]
    if (hero) out.heroes.push({ brick: b.name, name: b.name, file: join(root, 'hero', hero) })
    for (const n of files(join(root, 'pages'), /\.tsx$/)) {
      const name = n.slice(0, -4)
      if (!NAME.test(name)) errors.push(`brick ${b.name}: page file ${n}: name must match ${NAME}`)
      else if (ENGINE_PAGES.includes(name) || out.pages.some((p) => p.name === name)) errors.push(`brick ${b.name}: page ${name} already exists`)
      else out.pages.push({ brick: b.name, name, file: join(root, 'pages', n) })
    }
    for (const n of files(join(root, 'sections'), /\.tsx$/)) {
      const name = n.slice(0, -4)
      if (!NAME.test(name)) errors.push(`brick ${b.name}: section file ${n}: name must match ${NAME}`)
      else out.sections.push({ brick: b.name, name, file: join(root, 'sections', n) })
    }
  }
  const themes = [...ENGINE_THEMES, ...out.themes.map((t) => t.name)]
  const heroes = [...ENGINE_HEROES, ...out.heroes.map((h) => h.name)]
  const pages = [...ENGINE_PAGES, ...out.pages.map((p) => p.name)]
  if (nest.theme !== undefined && nest.theme !== null) {
    if (typeof nest.theme === 'string' && themes.includes(nest.theme)) out.config.theme = nest.theme
    else errors.push(`crow.yaml modules.nest.theme: unknown theme ${JSON.stringify(nest.theme)} (known: ${themes.join(', ')})`)
  }
  if (nest.hero !== undefined && nest.hero !== null) {
    if (typeof nest.hero === 'string' && heroes.includes(nest.hero)) out.config.hero = nest.hero
    else errors.push(`crow.yaml modules.nest.hero: unknown hero ${JSON.stringify(nest.hero)} (known: ${heroes.join(', ')})`)
  }
  if (nest.pages !== undefined && nest.pages !== null) {
    const list = Array.isArray(nest.pages) ? nest.pages : null
    if (!list) errors.push('crow.yaml modules.nest.pages: expects a list of page names')
    for (const p of list ?? []) {
      if (typeof p !== 'string' || !pages.includes(p)) errors.push(`crow.yaml modules.nest.pages: unknown page ${JSON.stringify(p)} (known: ${pages.join(', ')})`)
      else if (!out.config.pages.includes(p)) out.config.pages.push(p)
    }
  }
  return out
}

// The on bricks from crow-modules. Tests and NEST_MODULES=none build the engine defaults.
export function readBricks(engine: string): BrickRow[] {
  if (process.env.VITEST || process.env.NEST_MODULES === 'none') return []
  const out = execFileSync('python3', [join(engine, 'bin', 'crow-modules'), 'json'], { encoding: 'utf8' })
  return (JSON.parse(out).modules as BrickRow[]).map((m) => ({ name: m.name, on: m.on, dir: m.dir, config: m.config ?? {} }))
}

export function source(m: NestModules): string {
  const q = JSON.stringify
  const lines: string[] = []
  const comp = (kind: string, list: NestFile[]) =>
    list.map((f, i) => `import * as ${kind}${i} from ${q(f.file)};`)
  lines.push(...comp('hero', m.heroes), ...comp('page', m.pages), ...comp('section', m.sections))
  lines.push(`export const config = ${q(m.config)};`)
  lines.push(`export const themes = ${q([...ENGINE_THEMES, ...m.themes.map((t) => t.name)])};`)
  lines.push(`export const heroes = {${m.heroes.map((h, i) => `${q(h.name)}: hero${i}.default`).join(', ')}};`)
  lines.push(`export const pages = [${m.pages.map((p, i) => `{ key: ${q(p.name)}, brick: ${q(p.brick)}, mod: page${i} }`).join(', ')}];`)
  lines.push(`export const sections = [${m.sections.map((s, i) => `{ id: ${q(s.name)}, brick: ${q(s.brick)}, mod: section${i} }`).join(', ')}];`)
  return lines.join('\n') + '\n'
}

const ID = 'virtual:nest-modules'
// Brick theme CSS, imported by main.tsx right after the built-in token files.
const THEMES_ID = 'virtual:nest-themes'
export const themeSource = (m: NestModules) => m.themes.map((t) => `import ${JSON.stringify(t.file)};`).join('\n') + '\n'

export function nestModules(engine: string): Plugin {
  let mods: NestModules | undefined
  const get = () => {
    if (!mods) {
      mods = collect(readBricks(engine))
      if (mods.errors.length) throw new Error(`nest modules:\n  ${mods.errors.join('\n  ')}`)
    }
    return mods
  }
  return {
    name: 'nest-modules',
    config: () => {
      // vite dev serves brick files from outside the package root.
      const dirs = get().themes.concat(get().heroes, get().pages, get().sections).map((f) => f.file)
      return dirs.length ? { server: { fs: { allow: [engine, ...dirs] } } } : undefined
    },
    resolveId: (id) => (id === ID || id === THEMES_ID ? '\0' + id : undefined),
    load: (id) => (id === '\0' + ID ? source(get()) : id === '\0' + THEMES_ID ? themeSource(get()) : undefined),
    // A brick file that fails to compile names its brick. Rolldown wraps errors in `errors`
    // and may print the path relative to the package, so match ids and messages both ways.
    buildEnd(error) {
      if (!error || !mods) return
      const all = [error, ...((error as { errors?: Error[] }).errors ?? [])] as (Error & { id?: string })[]
      const text = all.map((e) => `${e.id ?? ''}\n${e.message}`).join('\n')
      const hit = [...mods.themes, ...mods.heroes, ...mods.pages, ...mods.sections].find((f) => text.includes(f.file) || text.includes(relative(process.cwd(), f.file)))
      if (hit) throw new Error(`brick ${hit.brick}: ${hit.file} does not build:\n${error.message}`)
    },
    // The first paint uses the configured theme; a saved pick wins only when it still exists.
    transformIndexHtml: (html) => {
      const m = get()
      const known = JSON.stringify([...ENGINE_THEMES, ...m.themes.map((t) => t.name)])
      return html
        .replace('data-theme="modern"', `data-theme="${m.config.theme}"`)
        .replace('if (t) document', `if (t && ${known}.includes(t)) document`)
    },
  }
}
