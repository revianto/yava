'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { IconSearch, IconBell, IconPlus } from '@/components/icons'

const NAV_LINKS = [
  { label: 'Dashboard', href: '/' },
  { label: 'Resep saya', href: '/recipes' },
  { label: 'Explore', href: '/explore' },
  { label: 'Grup', href: '/groups' },
]

export function Topnav() {
  const pathname = usePathname()

  return (
    <header className="topnav">
      <div className="row gap-4">
        <Link href="/" className="logo" style={{ fontSize: 26 }}>
          <span>YAVA</span><span className="dot">.</span>
        </Link>
        <nav className="topnav__links">
          {NAV_LINKS.map(({ label, href }) => {
            const isActive = href === '/' ? pathname === '/' : pathname.startsWith(href)
            return (
              <Link
                key={href}
                href={href}
                className={`topnav__link${isActive ? ' topnav__link--active' : ''}`}
              >
                {label}
              </Link>
            )
          })}
        </nav>
      </div>
      <div className="row gap-2">
        <div style={{ position: 'relative' }}>
          <IconSearch size={18} style={{ position: 'absolute', left: 14, top: '50%', transform: 'translateY(-50%)', color: 'var(--muted)' }} />
          <input className="input input--search" placeholder="Cari resep, jenis, atau bahan…" style={{ width: 280 }} />
        </div>
        <button className="icon-btn" title="Notifikasi"><IconBell size={18} /></button>
        <Link href="/recipes/new" className="btn btn--light-primary">
          <IconPlus size={16} /> Buat resep
        </Link>
        <span className="avatar" style={{ background: 'var(--coral-red)' }}>ND</span>
      </div>
    </header>
  )
}
