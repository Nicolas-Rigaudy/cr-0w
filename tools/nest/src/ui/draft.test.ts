import { describe, expect, it } from 'vitest'
import { readDraft, writeDraft } from './draft'

describe('drafts', () => {
  it('keeps typed text per key and drops it once sent', () => {
    writeDraft('desk.ask', 'half a question')
    expect(readDraft('desk.ask')).toBe('half a question')
    expect(readDraft('desk.task')).toBe('')
    writeDraft('desk.ask', '')
    expect(readDraft('desk.ask')).toBe('')
  })
})
