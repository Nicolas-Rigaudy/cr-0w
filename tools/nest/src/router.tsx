import { createContext, useContext, useEffect, useState, type AnchorHTMLAttributes, type MouseEvent } from 'react'

const PathContext = createContext('/')
export const usePath = () => useContext(PathContext)

export function navigate(to: string) {
  if (to === location.pathname) return
  history.pushState(null, '', to)
  dispatchEvent(new PopStateEvent('popstate'))
}

export function Router({ children }: { children: React.ReactNode }) {
  const [path, setPath] = useState(location.pathname)
  useEffect(() => {
    const on = () => {
      setPath(location.pathname)
      document.querySelector('.main')?.scrollTo(0, 0)
    }
    addEventListener('popstate', on)
    return () => removeEventListener('popstate', on)
  }, [])
  return <PathContext.Provider value={path}>{children}</PathContext.Provider>
}

// A real <a href> that routes in place on a plain click.
export function Link({ href, onClick, ...rest }: AnchorHTMLAttributes<HTMLAnchorElement> & { href: string }) {
  const click = (e: MouseEvent<HTMLAnchorElement>) => {
    onClick?.(e)
    if (e.defaultPrevented || e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) return
    e.preventDefault()
    navigate(href)
  }
  return <a href={href} onClick={click} {...rest} />
}

export type Route =
  | { tab: 'home' | 'desk' | 'missions' | 'inbox' | 'mrs' | 'agenda' | 'reports' }
  | { tab: 'missing' }

export function match(path: string): Route {
  const p = path.replace(/\/+$/, '') || '/'
  if (p === '/') return { tab: 'home' }
  // Old mission links (no detail screen since item_168) land on Missions.
  if (/^\/missions\/[^/]+$/.test(p)) return { tab: 'missions' }
  const key = p.slice(1)
  if (['desk', 'missions', 'inbox', 'mrs', 'agenda', 'reports'].includes(key)) return { tab: key as 'desk' }
  return { tab: 'missing' }
}
