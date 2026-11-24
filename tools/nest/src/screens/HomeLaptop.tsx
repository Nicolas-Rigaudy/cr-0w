import type { Snapshot } from '../data/useData'
import { HeroSlot } from '../hero/Slot'

// Laptop home (item_176, item_259): the hero fills the page; needs you, next and running
// sit in the Home rail header next to the chat (ui/HomeHud).
export function HomeLaptop({ snap }: { snap: Snapshot }) {
  const { missions, inbox, now } = snap
  return (
    <div className="lp home-lp">
      <section className="slot hero-slot" aria-label="CR-0W overview">
        <HeroSlot snap={snap} missions={missions} inbox={inbox} now={now} sections={snap.sections} />
      </section>
    </div>
  )
}
