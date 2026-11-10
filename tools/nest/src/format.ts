import type { Mission } from './data/types'

const hm = new Intl.DateTimeFormat('en-GB', { hour: '2-digit', minute: '2-digit' })
const day = new Intl.DateTimeFormat('en-GB', { weekday: 'short', day: '2-digit', month: 'short' })
const dayShort = new Intl.DateTimeFormat('en-GB', { weekday: 'short' })

export const time = (iso: string | Date) => hm.format(new Date(iso))
export const date = (d: Date) => day.format(d)
export const weekday = (d: Date) => dayShort.format(d)

export function duration(fromIso: string, to: Date | string) {
  const min = Math.max(0, Math.round((new Date(to).getTime() - new Date(fromIso).getTime()) / 60000))
  const h = Math.floor(min / 60)
  return h ? `${h}h ${String(min % 60).padStart(2, '0')}m` : `${min}m`
}

export function sameDay(a: Date, b: Date) {
  return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate()
}

// "20261109-build-demo-login" -> "Demo login"
const ACRONYMS = new Set(['api', 'mr', 'pr', 'sse', 'ui', 'ux', 'pwa', 'ci', 'url', 'id', 'tui', 'cli', 'mcp'])

export function title(m: Pick<Mission, 'id' | 'kind'>) {
  const slug = m.id.replace(/^\d{8}-/, '').replace(new RegExp(`^${m.kind}-`), '')
  const words = slug.split('-').map((w) => (ACRONYMS.has(w) ? w.toUpperCase() : w)).join(' ')
  return words.charAt(0).toUpperCase() + words.slice(1)
}

export function greeting(now: Date) {
  const h = now.getHours()
  return h < 12 ? 'Good morning, Sir' : h < 18 ? 'Good afternoon, Sir' : 'Good evening, Sir'
}

// Age of an event: "12m", "3h", "4d".
export function ago(iso: string, now: Date) {
  const min = Math.max(0, Math.round((now.getTime() - new Date(iso).getTime()) / 60000))
  if (min < 60) return `${min}m`
  if (min < 1440) return `${Math.floor(min / 60)}h`
  return `${Math.floor(min / 1440)}d`
}

// Elapsed run time as a clock: "01:05".
export function elapsed(fromIso: string, to: Date) {
  const min = Math.max(0, Math.floor((to.getTime() - new Date(fromIso).getTime()) / 60000))
  return `${String(Math.floor(min / 60)).padStart(2, '0')}:${String(min % 60).padStart(2, '0')}`
}

const wd = new Intl.DateTimeFormat('en-GB', { weekday: 'short' })
const mo = new Intl.DateTimeFormat('en-GB', { month: 'short' })
const dd = (d: Date) => String(d.getDate()).padStart(2, '0')

// Laptop header clock: "MON 09 NOV · 14:32" (modern), "MON 09.11 · 14:32" (cyber).
export function headerClock(d: Date, cyber: boolean) {
  const day = wd.format(d).toUpperCase()
  const md = cyber ? `${dd(d)}.${String(d.getMonth() + 1).padStart(2, '0')}` : `${dd(d)} ${mo.format(d).toUpperCase()}`
  return `${day} ${md} · ${time(d)}`
}
