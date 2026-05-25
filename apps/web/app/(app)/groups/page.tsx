'use client'

import Link from 'next/link'
import { GROUPS } from '@/lib/mock-data'
import { IconPlus, IconArrow, IconGroup } from '@/components/icons'

const ROLE_LABEL: Record<string, string> = {
  founder: 'Founder',
  admin: 'Admin',
  member: 'Anggota',
}

export default function GroupsPage() {
  const activeCount = GROUPS.reduce((a, g) => a + g.recipes.filter((r) => r.status === 'active').length, 0)
  const pendingCount = GROUPS.reduce((a, g) => a + g.recipes.filter((r) => r.status === 'pending').length, 0)

  return (
    <main className="container col gap-4" style={{ paddingTop: 16, gap: 48, paddingBottom: 64 }}>
      {/* Header */}
      <section>
        <div className="row between" style={{ alignItems: 'flex-end', marginBottom: 32 }}>
          <div>
            <div className="t-label" style={{ color: 'var(--muted)', marginBottom: 8, display: 'flex', alignItems: 'center', gap: 6 }}>
              <IconGroup size={12} /> Komunitas
            </div>
            <div className="t-display" style={{ fontSize: 48, letterSpacing: '-.025em' }}>
              Grup saya<span style={{ color: 'var(--coral-red)' }}>.</span>
            </div>
            <div className="t-body muted" style={{ marginTop: 6 }}>
              {GROUPS.length} grup · {activeCount} resep aktif{pendingCount > 0 && ` · ${pendingCount} menunggu approval`}
            </div>
          </div>
          <div className="row gap-2">
            <Link href="/groups/join" className="btn btn--light-secondary">
              Gabung grup
            </Link>
            <Link href="/groups/new" className="btn btn--light-primary">
              <IconPlus size={16} /> Buat grup
            </Link>
          </div>
        </div>

        {/* Group cards */}
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 16 }}>
          {GROUPS.map((g) => {
            const active = g.recipes.filter((r) => r.status === 'active').length
            const pending = g.recipes.filter((r) => r.status === 'pending').length
            const isAdmin = g.myRole === 'admin' || g.myRole === 'founder'
            return (
              <Link
                key={g.id}
                href={`/groups/${g.id}`}
                style={{ textDecoration: 'none', color: 'inherit' }}
              >
                <div
                  className="card card--dark card--hero"
                  style={{ padding: 28, cursor: 'pointer', display: 'flex', flexDirection: 'column', gap: 0, minHeight: 240 }}
                >
                  {/* Top row */}
                  <div className="row gap-2" style={{ marginBottom: 16 }}>
                    <span className="tag" style={{ background: 'rgba(225,191,145,.12)', color: 'var(--lilac)', border: 0 }}>
                      {ROLE_LABEL[g.myRole]}
                    </span>
                    {isAdmin && pending > 0 && (
                      <span className="tag" style={{ background: 'var(--coral-red)', color: '#FFF8E7', border: 0 }}>
                        {pending} PENDING
                      </span>
                    )}
                  </div>

                  <div className="t-h1" style={{ color: '#FFF8E7', lineHeight: 1.05, marginBottom: 8 }}>{g.name}</div>
                  <div className="t-small muted-dark" style={{ marginBottom: 20, flex: 1, lineHeight: 1.5 }}>{g.description}</div>

                  {/* Stats */}
                  <div className="row gap-4" style={{ marginBottom: 20, borderTop: '1px solid rgba(225,191,145,.12)', paddingTop: 16 }}>
                    <div>
                      <div className="t-h2" style={{ color: '#FFF8E7' }}>{g.members.length}</div>
                      <div className="t-label muted-dark">Anggota</div>
                    </div>
                    <div>
                      <div className="t-h2" style={{ color: '#FFF8E7' }}>{active}</div>
                      <div className="t-label muted-dark">Resep aktif</div>
                    </div>
                    {isAdmin && pending > 0 && (
                      <div>
                        <div className="t-h2" style={{ color: 'var(--coral-red)' }}>{pending}</div>
                        <div className="t-label muted-dark">Pending</div>
                      </div>
                    )}
                  </div>

                  {/* Avatar stack + CTA */}
                  <div className="row between" style={{ alignItems: 'center' }}>
                    <div className="row">
                      {g.members.slice(0, 4).map((m, i) => (
                        <span
                          key={m.id}
                          className="avatar avatar--sm"
                          style={{
                            marginLeft: i === 0 ? 0 : -8,
                            border: '2px solid var(--deep-ink)',
                            background: 'var(--lilac)',
                            color: 'var(--deep-ink)',
                            fontWeight: 700,
                          }}
                        >{m.initials.slice(0, 2)}</span>
                      ))}
                      {g.members.length > 4 && (
                        <span
                          className="avatar avatar--sm"
                          style={{ marginLeft: -8, border: '2px solid var(--deep-ink)', background: 'rgba(225,191,145,.18)', color: 'var(--lilac)' }}
                        >+{g.members.length - 4}</span>
                      )}
                    </div>
                    <span className="row gap-1 t-small" style={{ color: 'var(--lilac)', fontWeight: 700 }}>
                      Buka <IconArrow size={14} />
                    </span>
                  </div>
                </div>
              </Link>
            )
          })}

          {/* Create new CTA */}
          <Link
            href="/groups/new"
            style={{
              textDecoration: 'none',
              display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
              minHeight: 240, borderRadius: 22, gap: 8,
              border: '1.5px dashed var(--deep-ink)', color: 'var(--deep-ink)',
              fontWeight: 700,
            }}
          >
            <IconPlus size={28} />
            <span>Buat grup baru</span>
            <span className="t-caption" style={{ color: 'var(--muted)', fontWeight: 400 }}>Undang teman untuk berbagi resep</span>
          </Link>
        </div>
      </section>
    </main>
  )
}
