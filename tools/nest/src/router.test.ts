import { afterEach, describe, expect, it, vi } from 'vitest'
import { navigate, setReloadOnRoute } from './router'

afterEach(() => {
  setReloadOnRoute(false)
  vi.unstubAllGlobals()
})

const stub = () => {
  const loc = { pathname: '/', search: '', assign: vi.fn(), reload: vi.fn() }
  const hist = { pushState: vi.fn() }
  vi.stubGlobal('location', loc)
  vi.stubGlobal('history', hist)
  vi.stubGlobal('dispatchEvent', vi.fn())
  vi.stubGlobal('PopStateEvent', class { constructor(public type: string) {} })
  return { loc, hist }
}

describe('route change after a new Nest build (item_283)', () => {
  it('routes in place while the page runs the served build', () => {
    const { loc, hist } = stub()
    navigate('/ops')
    expect(hist.pushState).toHaveBeenCalledWith(null, '', '/ops')
    expect(loc.assign).not.toHaveBeenCalled()
  })

  it('loads the next route from core once a new build is served', () => {
    const { loc, hist } = stub()
    setReloadOnRoute(true)
    navigate('/ops')
    expect(loc.assign).toHaveBeenCalledWith('/ops')
    expect(hist.pushState).not.toHaveBeenCalled()
  })

  it('never reloads for a link to the page already open', () => {
    const { loc } = stub()
    setReloadOnRoute(true)
    navigate('/')
    expect(loc.assign).not.toHaveBeenCalled()
  })
})
