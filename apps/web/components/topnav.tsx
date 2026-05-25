'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useState, useRef, useEffect } from 'react'
import { IconSearch, IconBell, IconPlus } from '@/components/icons'
import { NOTIFICATIONS } from '@/lib/mock-data'
import type { Notification } from '@/types'

const NAV_LINKS = [
  { label: 'Dashboard', href: '/' },
  { label: 'Resep saya', href: '/recipes' },
  { label: 'Explore', href: '/explore' },
  { label: 'Grup', href: '/groups' },
]

const NOTIF_ICON: Record<string, string> = {
  approved: '✓',
  rejected: '✕',
  reply: '↩',
  joined: '+',
}

const NOTIF_COLOR: Record<string, string> = {
  approved: 'var(--deep-ink)',
  rejected: 'var(--coral-red)',
  reply: '#6B7FD4',
  joined: 'var(--coral-red)',
}

export function Topnav() {
  const pathname = usePathname()
  const [notifications, setNotifications] = useState<Notification[]>(NOTIFICATIONS)
  const [bellOpen, setBellOpen] = useState(false)
  const bellRef = useRef<HTMLDivElement>(null)

  const unread = notifications.filter((n) => !n.read).length

  useEffect(() => {
    if (!bellOpen) return
    const handler = (e: MouseEvent) => {
      if (bellRef.current && !bellRef.current.contains(e.target as Node)) setBellOpen(false)
    }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
  }, [bellOpen])

  const markAllRead = () => setNotifications((prev) => prev.map((n) => ({ ...n, read: true })))

  const markRead = (id: string) => setNotifications((prev) => prev.map((n) => n.id === id ? { ...n, read: true } : n))

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

        {/* Notification bell */}
        <div ref={bellRef} style={{ position: 'relative' }}>
          <button
            className="icon-btn"
            title="Notifikasi"
            onClick={() => setBellOpen((v) => !v)}
            style={{ position: 'relative' }}
          >
            <IconBell size={18} />
            {unread > 0 && (
              <span style={{
                position: 'absolute', top: 6, right: 6,
                width: 8, height: 8, borderRadius: '50%',
                background: 'var(--coral-red)',
                border: '1.5px solid var(--lavender-fog)',
              }} />
            )}
          </button>

          {bellOpen && (
            <div style={{
              position: 'absolute', top: 'calc(100% + 10px)', right: 0, zIndex: 200,
              background: '#FBF8F1', border: '1px solid var(--hairline)', borderRadius: 16,
              boxShadow: '0 12px 40px -8px rgba(45,82,74,.22)',
              width: 340, overflow: 'hidden',
            }}>
              {/* Header */}
              <div className="row between" style={{ padding: '14px 16px 12px', borderBottom: '1px solid var(--hairline)' }}>
                <span className="t-h3" style={{ fontSize: 15 }}>Notifikasi</span>
                {unread > 0 && (
                  <button
                    onClick={markAllRead}
                    style={{ background: 'none', border: 0, cursor: 'pointer', fontSize: 12, fontWeight: 600, color: 'var(--muted)', padding: 0 }}
                  >
                    Tandai semua dibaca
                  </button>
                )}
              </div>

              {/* List */}
              <div style={{ maxHeight: 380, overflowY: 'auto' }}>
                {notifications.length === 0 ? (
                  <div style={{ padding: '32px 16px', textAlign: 'center', color: 'var(--muted)', fontSize: 13 }}>
                    Tidak ada notifikasi
                  </div>
                ) : (
                  notifications.map((n) => (
                    <Link
                      key={n.id}
                      href={n.link ?? '#'}
                      onClick={() => { markRead(n.id); setBellOpen(false) }}
                      style={{ textDecoration: 'none', display: 'block' }}
                    >
                      <div
                        style={{
                          padding: '12px 16px',
                          background: n.read ? 'transparent' : 'rgba(173,130,87,.06)',
                          borderBottom: '1px solid var(--hairline)',
                          display: 'flex', gap: 12, alignItems: 'flex-start',
                          cursor: 'pointer',
                          transition: 'background 150ms',
                        }}
                        onMouseEnter={(e) => (e.currentTarget.style.background = 'var(--lavender-fog)')}
                        onMouseLeave={(e) => (e.currentTarget.style.background = n.read ? 'transparent' : 'rgba(173,130,87,.06)')}
                      >
                        <span style={{
                          width: 30, height: 30, borderRadius: '50%', flexShrink: 0,
                          background: NOTIF_COLOR[n.type],
                          display: 'inline-flex', alignItems: 'center', justifyContent: 'center',
                          color: '#FBF8F1', fontWeight: 700, fontSize: 13,
                        }}>
                          {NOTIF_ICON[n.type]}
                        </span>
                        <div style={{ flex: 1, minWidth: 0 }}>
                          <div style={{ fontSize: 13, fontWeight: n.read ? 500 : 700, color: 'var(--deep-ink)', marginBottom: 2 }}>
                            {n.title}
                          </div>
                          <div style={{ fontSize: 12, color: 'var(--muted)', lineHeight: 1.4 }}>{n.body}</div>
                          <div style={{ fontSize: 11, color: 'var(--muted)', marginTop: 4 }}>{n.createdAt}</div>
                        </div>
                        {!n.read && (
                          <span style={{ width: 7, height: 7, borderRadius: '50%', background: 'var(--coral-red)', flexShrink: 0, marginTop: 4 }} />
                        )}
                      </div>
                    </Link>
                  ))
                )}
              </div>
            </div>
          )}
        </div>

        <Link href="/recipes/new" className="btn btn--light-primary">
          <IconPlus size={16} /> Buat resep
        </Link>
        <span className="avatar" style={{ background: 'var(--coral-red)' }}>ND</span>
      </div>
    </header>
  )
}
