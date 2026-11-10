import { missions as coreMissions, inbox as coreInbox } from 'virtual:core-records'
import stateMissions from '../../fixtures/states/missions.json'
import stateInbox from '../../fixtures/states/inbox.json'
import stateAgenda from '../../fixtures/states/agenda.json'
import type { DataSource } from './source'
import type { AgendaEvent, InboxItem, Mission, SectionBody } from './types'

// Shared core-contract records plus Nest-local records that cover every UI state.
const MISSIONS = [...coreMissions, ...stateMissions] as Mission[]
const INBOX = [...coreInbox, ...stateInbox] as InboxItem[]
const AGENDA = stateAgenda as AgendaEvent[]
const NOW = new Date('2026-11-09T16:10:00+01:00')
// Script sections stay empty until their screens bring state fixtures (item_171/172).
const NO_SECTION: SectionBody = { data: null, stale: false, updated_at: null }

export function fixtureSource(missions: Mission[] = MISSIONS, inbox: InboxItem[] = INBOX, now: Date = NOW, agenda: AgendaEvent[] = AGENDA): DataSource {
  return {
    missions: async () => missions,
    mission: async (id) => missions.find((m) => m.id === id),
    inbox: async () => inbox,
    agenda: async () => agenda,
    section: async () => NO_SECTION,
    now: () => now,
  }
}
