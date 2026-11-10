import type { InboxItem, Mission } from './types'

// The only door between screens and data. The fixture adapter implements it
// now; the API adapter (slice 3) implements the same interface.
export interface DataSource {
  missions(): Promise<Mission[]>
  mission(id: string): Promise<Mission | undefined>
  inbox(): Promise<InboxItem[]>
  // Reference time for "today" and elapsed clocks. Fixtures pin it, the API uses the wall clock.
  now(): Date
}
