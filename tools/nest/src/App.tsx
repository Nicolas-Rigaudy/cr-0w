import { useFeed } from './data/useData'
import { match, usePath } from './router'
import { Home } from './screens/Home'
import { Inbox } from './screens/Inbox'
import { Desk } from './screens/Desk'
import { MissionDetail } from './screens/MissionDetail'
import { Missions } from './screens/Missions'
import { Placeholder } from './screens/Placeholder'
import { Shell } from './ui/Shell'

const LATER: Record<string, [string, string]> = {
  mrs: ['MRs', 'Merge requests arrive with the live API.'],
  agenda: ['Agenda', 'The calendar arrives with the live API.'],
  reports: ['Reports', 'Reports arrive with the live API.'],
  missing: ['Not found', 'This page does not exist.'],
}

export function App() {
  const route = match(usePath())
  const [snap, refresh] = useFeed()
  const active = route.tab === 'mission' ? 'missions' : route.tab === 'missing' ? 'home' : route.tab
  let body = null
  if (snap) {
    if (route.tab === 'home') body = <Home snap={snap} />
    else if (route.tab === 'missions') body = <Missions snap={snap} />
    else if (route.tab === 'mission') body = <MissionDetail snap={snap} id={route.id} />
    else if (route.tab === 'inbox') body = <Inbox snap={snap} />
    else if (route.tab === 'desk') body = <Desk snap={snap} refresh={refresh} />
    else body = <Placeholder name={LATER[route.tab][0]} note={LATER[route.tab][1]} />
  }
  return <Shell active={active} snap={snap}>{body}</Shell>
}
