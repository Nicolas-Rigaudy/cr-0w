import { useEffect, useMemo, useState } from 'react'
import { changedKeys, draftOf, saveBody, type Draft, type NestBuild, type NestSettings, type SettingKey, type SettingsApi } from './settings'

// What the Settings page shows (item_284): crow.yaml as loaded (base), the owner's draft,
// the save in flight or its outcome, and the Nest build that follows a save.
export interface SettingsView {
  status: 'off' | 'loading' | 'ready' | 'error'
  error?: string
  settings?: NestSettings
  base?: Draft
  draft?: Draft
  saving: boolean
  saveError?: string
  // The keys the last save changed; cleared by the next edit.
  saved?: SettingKey[]
  build?: NestBuild
}

export interface SettingsControl {
  load(): Promise<void>
  edit(change: (d: Draft) => Draft): void
  discard(): void
  save(): Promise<void>
  // Reads the build state once; the hook calls it while a build runs.
  tick(): Promise<void>
}

const msg = (e: unknown) => (e instanceof Error ? e.message : String(e))

// The page state without React, so tests drive it like the feed.
export function settingsControl(api: SettingsApi | undefined, emit: (v: SettingsView) => void): SettingsControl {
  let v: SettingsView = { status: api ? 'loading' : 'off', saving: false }
  const set = (p: Partial<SettingsView>) => emit((v = { ...v, ...p }))
  return {
    async load() {
      if (!api) return set({})
      try {
        const [s, build] = await Promise.all([api.get(), api.build().catch(() => undefined)])
        const base = draftOf(s)
        set({ status: 'ready', error: undefined, settings: s, base, draft: base, build })
      } catch (e) {
        set({ status: 'error', error: msg(e) })
      }
    },
    edit(change) {
      if (!v.draft || v.saving) return
      set({ draft: change(v.draft), saveError: undefined, saved: undefined })
    },
    discard() {
      if (v.base && !v.saving) set({ draft: v.base, saveError: undefined })
    },
    async save() {
      if (!api || !v.base || !v.draft || !v.settings || v.saving) return
      const body = saveBody(v.base, v.draft, v.settings.known)
      if (!Object.keys(body).length) return
      const keys = changedKeys(v.base, v.draft)
      set({ saving: true, saveError: undefined, saved: undefined })
      try {
        const r = await api.save(body)
        const base = draftOf(r)
        set({ saving: false, settings: r, base, draft: base, saved: keys, build: r.build })
      } catch (e) {
        set({ saving: false, saveError: msg(e) })
      }
    },
    async tick() {
      if (!api) return
      try {
        set({ build: await api.build() })
      } catch { /* core away: the next tick tries again */ }
    },
  }
}

const POLL = 1500

export function useSettings(api: SettingsApi | undefined): [SettingsView, SettingsControl] {
  const [view, setView] = useState<SettingsView>({ status: api ? 'loading' : 'off', saving: false })
  const ctl = useMemo(() => settingsControl(api, setView), [api])
  useEffect(() => { void ctl.load() }, [ctl])
  const building = view.build?.state === 'building'
  useEffect(() => {
    if (!building) return
    const t = setInterval(() => void ctl.tick(), POLL)
    return () => clearInterval(t)
  }, [building, ctl])
  return [view, ctl]
}
