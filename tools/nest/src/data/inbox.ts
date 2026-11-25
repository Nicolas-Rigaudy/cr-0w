import type { Snapshot } from './feed'
import { mailOf, type Mail } from './live'
import type { InboxItem } from './types'
import { workOf, type WorkItem } from './work'

// The Inbox "To answer" list (item_204, owner 2026-11-16 second review): only what the
// owner reads or answers. A question is a crow inbox question (`ask`, from a mission) or a
// live session or run waiting on a prompt (`prompt`); mail is flagged or important mail.
// Failed runs and red-pipeline MRs stay on Desk.
export type AnswerRow =
  | { key: string; kind: 'ask'; at: string; item: InboxItem }
  | { key: string; kind: 'prompt'; at?: string; item: WorkItem }
  | { key: string; kind: 'mail'; at?: undefined; item: Mail }

export const isAsk = (r: InboxItem) => r.state === 'open' && r.type !== 'fyi'
export const isFyi = (r: InboxItem) => r.state === 'open' && r.type === 'fyi'
// A session or run waiting on a prompt or choice. MRs waiting on a review stay on Desk
// and MRs. A mission joins only through its inbox question (an ask); a mission waiting
// on the owner without one (a review, say) stays on Desk.
export const isPrompt = (w: WorkItem) => w.state === 'needs' && (w.kind === 'session' || w.kind === 'run')

const valid = (iso: string | undefined) => !!iso && !isNaN(Date.parse(iso))
// Oldest first; rows with no time (mail) go last in source order.
const byWait = (a: { at?: string }, b: { at?: string }) => {
  const x = valid(a.at) ? Date.parse(a.at!) : Infinity
  const y = valid(b.at) ? Date.parse(b.at!) : Infinity
  return x === y ? 0 : x < y ? -1 : 1
}

// When a work item started waiting: now minus wait_s, else its since.
function waitFrom(w: WorkItem, now?: Date): string | undefined {
  if (now && w.wait_s != null) return new Date(now.getTime() - w.wait_s * 1000).toISOString()
  return w.since
}

// Questions and mail by how long they wait.
export function toAnswer(inbox: InboxItem[], work: WorkItem[] | null, mail: Mail[], now?: Date): AnswerRow[] {
  const asks: AnswerRow[] = inbox.filter(isAsk).map((r) => ({ key: `q:${r.id}`, kind: 'ask' as const, at: r.created, item: r }))
  const prompts: AnswerRow[] = (work ?? []).filter(isPrompt)
    .map((w) => ({ key: `w:${w.id}`, kind: 'prompt' as const, at: waitFrom(w, now), item: w }))
  const mails: AnswerRow[] = mail.map((m) => ({ key: `m:${m.id}`, kind: 'mail' as const, item: m }))
  return [...asks, ...prompts, ...mails].sort(byWait)
}

export const answerOf = (s: Snapshot) => toAnswer(s.inbox, workOf(s.sections.work), mailOf(s.sections.today), s.now)

// Inbox kinds (item_223, board 223-B): the kind rail on the laptop and the chips on the phone.
// Follow-ups wait on someone else and never count in the badge.
export type InboxKind = 'decision' | 'review' | 'approval' | 'ask' | 'mail' | 'followup'
export const KINDS: Array<[InboxKind, string]> = [
  ['decision', 'Decisions'], ['review', 'Reviews'], ['approval', 'Approvals'], ['ask', 'Asks'], ['mail', 'Mail'], ['followup', 'Follow-ups'],
]
export const KIND_WORD: Record<InboxKind, string> = {
  decision: 'decision', review: 'review', approval: 'approval', ask: 'ask', mail: 'mail', followup: 'follow-up',
}

const IMAGE = /\.(png|jpe?g|webp|gif|svg)$/i
export const images = (r: InboxItem) => r.artifacts.filter((a) => IMAGE.test(a))
// A review with images is a design review (item_211); without, a doc review.
export const isDesign = (r: InboxItem) => r.source === 'mockup' || ((r.type === 'review' || r.type === 'choice') && images(r).length > 1)

// The options and scenarios of a design review (item_211, board 211-B). A one-letter token
// in the file name is the option, the words after it the scenario (`boot-A-dark-room.png`);
// files without a letter become options A, B, C in order.
export type Design = { letters: string[]; scenarios: string[]; file: (letter: string, scenario: string) => string | undefined }
export function designOf(r: InboxItem): Design {
  const files = images(r).map((path, n) => {
    const words = (path.split('/').pop() ?? path).replace(IMAGE, '').split(/[-_ .]+/).filter(Boolean)
    const at = words.findIndex((w) => /^[A-Z]$/.test(w))
    if (at < 0) return { path, letter: String.fromCharCode(65 + n), scenario: '' }
    const after = words.slice(at + 1)
    return { path, letter: words[at], scenario: (after.length ? after : words.slice(0, at)).join(' ') }
  })
  const uniq = (xs: string[]) => [...new Set(xs)]
  return {
    letters: uniq(files.map((f) => f.letter)),
    scenarios: uniq(files.map((f) => f.scenario)),
    file: (letter, scenario) => (files.find((f) => f.letter === letter && f.scenario === scenario) ?? files.find((f) => f.letter === letter))?.path,
  }
}

export function kindOfItem(r: InboxItem): InboxKind {
  if (isDesign(r) || r.type === 'review') return 'review'
  if (r.type === 'approval') return 'approval'
  if (r.type === 'choice') return 'decision'
  return 'ask'
}

export function kindOf(r: AnswerRow): InboxKind {
  if (r.kind === 'mail') return 'mail'
  if (r.kind === 'prompt') return 'ask'
  return kindOfItem(r.item)
}

// What the row is, in words: "design review", "doc review", or the kind.
export function kindWord(r: AnswerRow): string {
  const k = kindOf(r)
  if (k !== 'review' || r.kind !== 'ask') return KIND_WORD[k]
  return isDesign(r.item) ? 'design review' : 'doc review'
}

const sameDay = (a: Date, b: Date) => a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate()
// Done today: rows answered or hidden since midnight, newest first.
export function doneToday(inbox: InboxItem[], now: Date): InboxItem[] {
  return inbox
    .filter((r) => r.state !== 'open' && r.type !== 'fyi' && valid(r.answered ?? undefined) && sameDay(new Date(r.answered!), now))
    .sort((a, b) => b.answered!.localeCompare(a.answered!))
}

// Inbox tab badge (owner 2026-11-16): the number of rows in "To answer", so the badge
// and the list always agree.
export function inboxBadge(s: Snapshot): { n: number } {
  return { n: answerOf(s).length }
}
