import { missions as coreMissions, inbox as coreInbox } from 'virtual:core-records'
import stateMissions from '../../fixtures/states/missions.json'
import stateInbox from '../../fixtures/states/inbox.json'
import type { DataSource } from './source'
import type { InboxItem, Mission } from './types'

// Shared core-contract records plus Nest-local records that cover every UI state.
const MISSIONS = [...coreMissions, ...stateMissions] as Mission[]
const INBOX = [...coreInbox, ...stateInbox] as InboxItem[]
const NOW = new Date('2026-11-09T16:10:00+01:00')

export function fixtureSource(missions: Mission[] = MISSIONS, inbox: InboxItem[] = INBOX, now: Date = NOW): DataSource {
  return {
    missions: async () => missions,
    mission: async (id) => missions.find((m) => m.id === id),
    inbox: async () => inbox,
    now: () => now,
  }
}
