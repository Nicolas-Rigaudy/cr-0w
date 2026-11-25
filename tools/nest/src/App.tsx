import { useMemo } from 'react'
import { scopeSnapshot } from './data/scope'
import { useFeed } from './data/useData'
import { useScope } from './data/useScope'
import { match, usePath } from './router'
import { pages as brickPages } from 'virtual:nest-modules'
import { BrickModuleSections } from './modules/Sections'
import { Home, INNER } from './screens/Home'
import { Inbox } from './screens/Inbox'
import { Mrs } from './screens/Mrs'
import { Chat } from './screens/Chat'
import { Desk } from './screens/Desk'
import { Ops } from './screens/Ops'
import { Placeholder } from './screens/Placeholder'
import { Shell } from './ui/Shell'
import { ButlerRail } from './ui/ButlerRail'
import { ButlerBadge, ButlerSheet } from './ui/ButlerSheet'
import { ButlerProvider } from './data/useButler'
import { LAPTOP, useMedia } from './ui/useMedia'
import { screenTab, shownRoute } from './ui/tabs'
import { BrickSections } from './ui/BrickSections'
import { listSections } from './data/modules'

const LATER: Record<string, [string, string]> = {
  agenda: ['Agenda', 'The calendar arrives with the live API.'],
  reports: ['Reports', 'Reports arrive with the live API.'],
  missing: ['Not found', 'This page does not exist.'],
}

export function App() {
  // A page left out of crow.yaml modules.nest.pages falls back to Home (item_267).
  const route = shownRoute(match(usePath()))
  const [raw, refresh] = useFeed()
  const [scope, setScope] = useScope()
  // Every screen reads the scoped snapshot; Chat reads the whole one so a link to any session opens.
  const snap = useMemo(() => raw && scopeSnapshot(raw, scope), [raw, scope])
  const active = screenTab(route)
  const laptop = useMedia(LAPTOP)
  const inner = useMedia(INNER)
  let body = null
  if (snap && raw) {
    if (route.tab === 'home') body = <Home snap={snap} />
    else if (route.tab === 'ops') body = <Ops snap={snap} refresh={refresh} />
    else if (route.tab === 'inbox') body = <Inbox snap={snap} refresh={refresh} />
    else if (route.tab === 'mrs') body = <Mrs snap={snap} refresh={refresh} />
    else if (route.tab === 'desk') body = <Desk snap={snap} refresh={refresh} />
    else if (route.tab === 'chat') body = <Chat snap={raw} refresh={refresh} target={route.target} />
    else if (route.tab === 'brick') {
      const Page = brickPages.find((p) => p.key === route.key)!.mod.default
      body = <Page snap={snap} />
    } else body = <Placeholder name={LATER[route.tab][0]} note={LATER[route.tab][1]} />
    // Brick TSX sections (item_267) go after the page's own panels, before the generic ones.
    if (route.tab === 'home' || route.tab === 'desk') body = <BrickModuleSections snap={snap} page={route.tab}>{body}</BrickModuleSections>
    // Generic brick sections go after the page's own panels (item_237).
    if (route.tab !== 'chat' && route.tab !== 'missing' && route.tab !== 'brick') body = listSections(snap.modules ?? null, route.tab).length ? <div className="with-bricks">{body}<BrickSections snap={snap} page={route.tab} /></div> : body
  }
  // Butler rail on every laptop screen (item_176).
  // Home has its own rail (item_259): keyed apart so its fold and width never leak to other pages.
  const home = route.tab === 'home'
  const rail = laptop && snap ? <ButlerRail key={home ? 'home' : 'page'} snap={snap} home={home} /> : null
  // Phone and Fold: a header badge and a chat sheet instead of the rail.
  // Fold inner Home shows the chat in its side (item_259): no badge there.
  const badge = !laptop && snap && !(home && inner) ? <ButlerBadge snap={snap} /> : null
  return (
    <ButlerProvider>
      <Shell active={active} snap={snap} scope={scope} setScope={setScope} rail={rail} badge={badge}>{body}</Shell>
      {!laptop && snap && <ButlerSheet snap={snap} />}
    </ButlerProvider>
  )
}
