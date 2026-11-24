// Brick TSX sections (item_267): nest/sections/<name>.tsx of on bricks, after Home or Desk panels.
import type { ReactNode } from 'react'
import { sections } from 'virtual:nest-modules'
import type { Snapshot } from '../data/useData'

export function BrickModuleSections({ snap, page, children }: { snap: Snapshot; page: 'home' | 'desk'; children: ReactNode }) {
  const list = sections.filter((s) => s.mod.page === page)
  if (!list.length) return <>{children}</>
  return (
    <div className="with-bricks">
      {children}
      {list.map(({ id, brick, mod: { default: Section } }) => (
        <section key={`${brick}/${id}`} className="brick-section" data-brick={brick}>
          <Section snap={snap} />
        </section>
      ))}
    </div>
  )
}
