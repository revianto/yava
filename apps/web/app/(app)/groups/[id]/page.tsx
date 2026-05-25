'use client'

import { use, useState, useRef, useEffect } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { notFound } from 'next/navigation'
import { GROUPS } from '@/lib/mock-data'
import type { GroupMember, GroupRecipeItem } from '@/types'
import {
  IconArrowLeft, IconArrow, IconCheck, IconClose, IconCopy,
  IconEdit, IconGroup, IconPlay, IconPlus,
} from '@/components/icons'

const ROLE_LABEL: Record<string, string> = { founder: 'Founder', admin: 'Admin', member: 'Anggota' }
const ROLE_ORDER: Record<string, number> = { founder: 0, admin: 1, member: 2 }

// ── small components ─────────────────────────────────────────────────────────

function RoleBadge({ role }: { role: string }) {
  if (role === 'founder') return <span className="tag" style={{ background: 'var(--coral-red)', color: '#FFF8E7', border: 0 }}>FOUNDER</span>
  if (role === 'admin') return <span className="tag" style={{ background: 'var(--deep-ink)', color: 'var(--lilac)', border: 0 }}>ADMIN</span>
  return <span className="tag tag--default">ANGGOTA</span>
}

function StatusBadge({ status }: { status: string }) {
  if (status === 'active') return <span className="tag tag--v60">AKTIF</span>
  if (status === 'pending') return <span className="tag" style={{ background: 'rgba(173,130,87,.15)', color: 'var(--coral-red)', border: '1px solid var(--coral-red)' }}>PENDING</span>
  return <span className="tag tag--default" style={{ color: 'var(--muted)' }}>DITOLAK</span>
}

// ── Tab: Resep ────────────────────────────────────────────────────────────────

function TabResep({
  recipes, isAdmin, onApprove, onReject,
}: {
  recipes: GroupRecipeItem[]
  isAdmin: boolean
  onApprove: (id: string) => void
  onReject: (id: string, reason: string) => void
}) {
  const active = recipes.filter((r) => r.status === 'active')
  const pending = recipes.filter((r) => r.status === 'pending')
  const rejected = recipes.filter((r) => r.status === 'rejected')
  const [rejectTarget, setRejectTarget] = useState<string | null>(null)
  const [rejectReason, setRejectReason] = useState('')

  return (
    <div className="col gap-4">
      {/* Active recipes */}
      <div>
        <div className="row between" style={{ marginBottom: 14, alignItems: 'flex-end' }}>
          <div>
            <div className="t-label muted" style={{ marginBottom: 4 }}>Koleksi grup</div>
            <div className="t-h2">{active.length} resep aktif</div>
          </div>
        </div>
        {active.length === 0 ? (
          <div className="card" style={{ padding: 32, textAlign: 'center' }}>
            <div className="t-body muted">Belum ada resep aktif. Submit resep kamu!</div>
          </div>
        ) : (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 12 }}>
            {active.map((r) => (
              <Link key={r.id} href={`/recipes/${r.recipeId}`} style={{ textDecoration: 'none', color: 'inherit' }}>
                <div className="recipe-card">
                  <div className="row gap-2 wrap" style={{ marginBottom: 4 }}>
                    <span className="tag tag--group">{r.recipeType.toUpperCase()}</span>
                  </div>
                  <div className="recipe-card__title">{r.recipeName}</div>
                  <div className="recipe-card__meta">
                    <span>{r.recipeSubtype}</span>
                    <span className="dot" />
                    <span>oleh {r.submittedBy}</span>
                  </div>
                  <div className="row between" style={{ marginTop: 'auto', paddingTop: 10 }}>
                    <span className="t-caption">{r.submittedAt}</span>
                    <span className="row gap-1" style={{ color: 'var(--deep-ink)', fontWeight: 700, fontSize: 13 }}>
                      Lihat <IconArrow size={14} />
                    </span>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>

      {/* Pending — admin only */}
      {isAdmin && pending.length > 0 && (
        <div>
          <div style={{ marginBottom: 14 }}>
            <div className="t-label" style={{ color: 'var(--coral-red)', marginBottom: 4 }}>Menunggu approval</div>
            <div className="t-h2">{pending.length} resep pending</div>
          </div>
          <div className="col gap-3">
            {pending.map((r) => (
              <div key={r.id} className="card" style={{ padding: 20 }}>
                <div className="row between" style={{ marginBottom: 12, alignItems: 'flex-start' }}>
                  <div>
                    <div style={{ fontWeight: 700, fontSize: 16, marginBottom: 2 }}>{r.recipeName}</div>
                    <div className="t-small muted">{r.recipeType} · {r.recipeSubtype} · oleh {r.submittedBy} · {r.submittedAt}</div>
                  </div>
                  <StatusBadge status={r.status} />
                </div>

                {rejectTarget === r.id ? (
                  <div className="col gap-2">
                    <input
                      className="input"
                      value={rejectReason}
                      onChange={(e) => setRejectReason(e.target.value)}
                      placeholder="Alasan penolakan (opsional)…"
                      style={{ width: '100%' }}
                      autoFocus
                    />
                    <div className="row gap-2">
                      <button
                        className="btn btn--destructive"
                        onClick={() => { onReject(r.id, rejectReason); setRejectTarget(null); setRejectReason('') }}
                      >
                        Tolak resep
                      </button>
                      <button className="btn btn--secondary" onClick={() => { setRejectTarget(null); setRejectReason('') }}>
                        Batal
                      </button>
                    </div>
                  </div>
                ) : (
                  <div className="row gap-2">
                    <button
                      className="btn btn--primary"
                      onClick={() => onApprove(r.id)}
                      style={{ flex: 1 }}
                    >
                      <IconCheck size={16} /> Setujui
                    </button>
                    <button
                      className="btn btn--secondary"
                      onClick={() => setRejectTarget(r.id)}
                      style={{ flex: 1 }}
                    >
                      <IconClose size={16} /> Tolak
                    </button>
                    <Link href={`/recipes/${r.recipeId}`} className="btn btn--secondary" style={{ textDecoration: 'none' }}>
                      <IconArrow size={14} /> Lihat
                    </Link>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Rejected — admin only */}
      {isAdmin && rejected.length > 0 && (
        <div>
          <div className="t-label muted" style={{ marginBottom: 12 }}>Riwayat ditolak ({rejected.length})</div>
          <div className="col gap-2">
            {rejected.map((r) => (
              <div key={r.id} className="card" style={{ padding: 16, opacity: 0.7 }}>
                <div className="row between" style={{ marginBottom: r.rejectionReason ? 8 : 0 }}>
                  <div>
                    <div style={{ fontWeight: 700, fontSize: 14 }}>{r.recipeName}</div>
                    <div className="t-caption">oleh {r.submittedBy} · {r.submittedAt}</div>
                  </div>
                  <StatusBadge status={r.status} />
                </div>
                {r.rejectionReason && (
                  <div className="t-small muted" style={{ borderTop: '1px solid var(--hairline)', paddingTop: 8, marginTop: 4 }}>
                    Alasan: {r.rejectionReason}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

// ── Tab: Anggota ──────────────────────────────────────────────────────────────

function TabAnggota({
  members, isAdmin, isFounder, groupId,
  onRemove, onPromote, onDemote,
}: {
  members: GroupMember[]
  isAdmin: boolean
  isFounder: boolean
  groupId: string
  onRemove: (id: string) => void
  onPromote: (id: string) => void
  onDemote: (id: string) => void
}) {
  const sorted = [...members].sort((a, b) => ROLE_ORDER[a.role] - ROLE_ORDER[b.role])
  const inviteUrl = `${typeof window !== 'undefined' ? window.location.origin : ''}/groups/join?code=${
    GROUPS.find((g) => g.id === groupId)?.inviteCode ?? ''
  }`
  const [copied, setCopied] = useState(false)

  const copyInvite = () => {
    navigator.clipboard.writeText(inviteUrl).catch(() => {})
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="col gap-4">
      {/* Invite section — admin only */}
      {isAdmin && (
        <div className="card card--dark" style={{ padding: 20 }}>
          <div className="t-label muted-dark" style={{ marginBottom: 8 }}>Invite link</div>
          <div className="row gap-2">
            <div style={{
              flex: 1, padding: '10px 14px', borderRadius: 10,
              background: 'rgba(225,191,145,.08)', border: '1px solid rgba(225,191,145,.16)',
              fontSize: 13, fontFamily: 'monospace', color: 'var(--lilac)', letterSpacing: '.04em',
              overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap',
            }}>
              {inviteUrl}
            </div>
            <button className="btn btn--light-primary" onClick={copyInvite} style={{ flex: 'none' }}>
              {copied ? <><IconCheck size={14} /> Disalin</> : <><IconCopy size={14} /> Salin</>}
            </button>
          </div>
        </div>
      )}

      {/* Member list */}
      <div>
        <div className="t-label muted" style={{ marginBottom: 12 }}>{members.length} Anggota</div>
        <div className="col">
          {sorted.map((m, i) => {
            const isMe = m.id.includes('me')
            const canManage = isAdmin && !isMe && m.role !== 'founder'
            return (
              <div
                key={m.id}
                className="row gap-3"
                style={{ padding: '14px 0', borderTop: i === 0 ? 0 : '1px solid var(--hairline)', alignItems: 'center' }}
              >
                <span className="avatar" style={{ background: isMe ? 'var(--coral-red)' : 'var(--deep-ink)', color: isMe ? '#FFF8E7' : 'var(--lilac)' }}>
                  {m.initials.slice(0, 2)}
                </span>
                <div className="grow">
                  <div style={{ fontWeight: 700, fontSize: 15 }}>{m.name}{isMe && <span className="t-caption muted"> (kamu)</span>}</div>
                  <div className="t-caption">Bergabung {m.joinedAt}</div>
                </div>
                <RoleBadge role={m.role} />
                {canManage && (
                  <div className="row gap-1">
                    {m.role === 'member' && isFounder && (
                      <button
                        className="btn btn--secondary"
                        style={{ height: 32, padding: '0 12px', fontSize: 12 }}
                        onClick={() => onPromote(m.id)}
                        title="Jadikan Admin"
                      >
                        Promote
                      </button>
                    )}
                    {m.role === 'admin' && isFounder && (
                      <button
                        className="btn btn--secondary"
                        style={{ height: 32, padding: '0 12px', fontSize: 12 }}
                        onClick={() => onDemote(m.id)}
                        title="Turunkan ke Member"
                      >
                        Demote
                      </button>
                    )}
                    <button
                      className="btn btn--destructive"
                      style={{ height: 32, padding: '0 12px', fontSize: 12 }}
                      onClick={() => onRemove(m.id)}
                    >
                      Keluarkan
                    </button>
                  </div>
                )}
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}

// ── Tab: Pengaturan ───────────────────────────────────────────────────────────

function TabPengaturan({
  group, isAdmin, isFounder, onDisband,
}: {
  group: (typeof GROUPS)[0]
  isAdmin: boolean
  isFounder: boolean
  onDisband: () => void
}) {
  const [name, setName] = useState(group.name)
  const [desc, setDesc] = useState(group.description)
  const [saved, setSaved] = useState(false)
  const [confirmDisband, setConfirmDisband] = useState(false)

  const save = () => { setSaved(true); setTimeout(() => setSaved(false), 2000) }

  return (
    <div className="col gap-4" style={{ maxWidth: 560 }}>
      {isAdmin ? (
        <div className="card" style={{ padding: 24 }}>
          <div className="t-h3" style={{ marginBottom: 20 }}>Info grup</div>
          <div className="col gap-3">
            <div>
              <div className="t-label muted" style={{ marginBottom: 8 }}>Nama grup</div>
              <input className="input" value={name} onChange={(e) => setName(e.target.value)} style={{ width: '100%' }} />
            </div>
            <div>
              <div className="t-label muted" style={{ marginBottom: 8 }}>Deskripsi</div>
              <textarea
                className="input"
                value={desc}
                onChange={(e) => setDesc(e.target.value)}
                rows={3}
                style={{ width: '100%', height: 'auto', resize: 'vertical', paddingTop: 10, paddingBottom: 10 }}
              />
            </div>
            <div>
              <button className="btn btn--primary" onClick={save}>
                {saved ? <><IconCheck size={14} /> Tersimpan</> : <><IconEdit size={14} /> Simpan perubahan</>}
              </button>
            </div>
          </div>
        </div>
      ) : (
        <div className="card" style={{ padding: 24 }}>
          <div className="t-h3" style={{ marginBottom: 8 }}>{group.name}</div>
          <div className="t-body muted">{group.description}</div>
          <div className="t-caption muted" style={{ marginTop: 12 }}>Dibuat {group.createdAt} · Hanya Admin yang bisa mengubah info grup.</div>
        </div>
      )}

      <div className="card" style={{ padding: 24 }}>
        <div className="t-h3" style={{ marginBottom: 6 }}>Invite code</div>
        <div className="t-small muted" style={{ marginBottom: 14 }}>Bagikan kode ini untuk mengundang anggota baru.</div>
        <div style={{
          display: 'inline-flex', alignItems: 'center', gap: 12, padding: '12px 20px',
          borderRadius: 12, background: 'rgba(45,82,74,.06)', border: '1px solid var(--hairline)',
        }}>
          <span style={{ fontFamily: 'monospace', fontWeight: 700, fontSize: 22, letterSpacing: '.08em', color: 'var(--deep-ink)' }}>
            {group.inviteCode}
          </span>
        </div>
      </div>

      {/* Danger zone — founder only */}
      {isFounder && (
        <div className="card" style={{ padding: 24, border: '1px solid rgba(173,130,87,.35)' }}>
          <div className="t-h3" style={{ marginBottom: 6, color: 'var(--coral-red)' }}>Danger zone</div>
          <div className="t-small muted" style={{ marginBottom: 16 }}>
            Membubarkan grup akan menghapus semua data grup secara permanen. Aksi ini tidak bisa dibatalkan.
          </div>
          {confirmDisband ? (
            <div className="col gap-2">
              <div className="t-small" style={{ fontWeight: 700 }}>Yakin ingin membubarkan grup ini?</div>
              <div className="row gap-2">
                <button className="btn btn--destructive" onClick={onDisband}>Ya, bubarkan</button>
                <button className="btn btn--secondary" onClick={() => setConfirmDisband(false)}>Batal</button>
              </div>
            </div>
          ) : (
            <button className="btn btn--destructive" onClick={() => setConfirmDisband(true)}>
              Bubarkan grup
            </button>
          )}
        </div>
      )}
    </div>
  )
}

// ── Root ──────────────────────────────────────────────────────────────────────

export default function GroupDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params)
  const router = useRouter()
  const original = GROUPS.find((g) => g.id === id)
  if (!original) notFound()

  const [recipes, setRecipes] = useState(original.recipes)
  const [members, setMembers] = useState(original.members)
  const [activeTab, setActiveTab] = useState<'resep' | 'anggota' | 'pengaturan'>('resep')

  const isAdmin = original.myRole === 'admin' || original.myRole === 'founder'
  const isFounder = original.myRole === 'founder'

  const activeCount = recipes.filter((r) => r.status === 'active').length
  const pendingCount = recipes.filter((r) => r.status === 'pending').length

  const handleApprove = (recipeItemId: string) =>
    setRecipes((prev) => prev.map((r) => r.id === recipeItemId ? { ...r, status: 'active' as const } : r))

  const handleReject = (recipeItemId: string, reason: string) =>
    setRecipes((prev) => prev.map((r) => r.id === recipeItemId ? { ...r, status: 'rejected' as const, rejectionReason: reason || 'Tidak memenuhi kriteria grup.' } : r))

  const handleRemoveMember = (memberId: string) =>
    setMembers((prev) => prev.filter((m) => m.id !== memberId))

  const handlePromote = (memberId: string) =>
    setMembers((prev) => prev.map((m) => m.id === memberId ? { ...m, role: 'admin' as const } : m))

  const handleDemote = (memberId: string) =>
    setMembers((prev) => prev.map((m) => m.id === memberId ? { ...m, role: 'member' as const } : m))

  const handleDisband = () => router.push('/groups')

  const TABS = [
    { key: 'resep', label: `Resep (${activeCount})` },
    { key: 'anggota', label: `Anggota (${members.length})` },
    { key: 'pengaturan', label: 'Pengaturan' },
  ] as const

  return (
    <main className="container" style={{ paddingBottom: 64 }}>
      {/* Back */}
      <button
        onClick={() => router.back()}
        style={{ background: 'transparent', border: 0, cursor: 'pointer', display: 'inline-flex', alignItems: 'center', gap: 6, color: 'var(--muted)', fontWeight: 500, fontSize: 13, padding: 0, marginBottom: 24 }}
      >
        <IconArrowLeft size={14} /> Grup saya
      </button>

      {/* Header card */}
      <div className="card card--electric card--hero" style={{ padding: 32, marginBottom: 32, position: 'relative', overflow: 'hidden' }}>
        <div className="row gap-2" style={{ marginBottom: 16, flexWrap: 'wrap' }}>
          <span className="tag" style={{ background: 'rgba(225,191,145,.12)', color: 'var(--lilac)', border: 0 }}>
            <IconGroup size={11} /> {ROLE_LABEL[original.myRole]}
          </span>
          {isAdmin && pendingCount > 0 && (
            <span className="tag" style={{ background: 'var(--coral-red)', color: '#FFF8E7', border: 0 }}>
              {pendingCount} PENDING
            </span>
          )}
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr auto', gap: 32, alignItems: 'flex-end' }}>
          <div>
            <div className="t-h1" style={{ color: '#FFF8E7', marginBottom: 8 }}>{original.name}</div>
            <div className="t-body muted-dark" style={{ maxWidth: 540 }}>{original.description}</div>
          </div>
          <div className="row gap-3">
            {[
              { label: 'Anggota', val: members.length },
              { label: 'Resep aktif', val: activeCount },
              ...(isAdmin && pendingCount > 0 ? [{ label: 'Pending', val: pendingCount }] : []),
            ].map(({ label, val }) => (
              <div key={label} style={{ textAlign: 'right' }}>
                <div className="t-h1" style={{ color: '#FFF8E7' }}>{val}</div>
                <div className="t-label muted-dark">{label}</div>
              </div>
            ))}
          </div>
        </div>

        <div className="row between" style={{ marginTop: 24, paddingTop: 16, borderTop: '1px solid rgba(225,191,145,.14)', alignItems: 'center' }}>
          <div className="row">
            {members.slice(0, 5).map((m, i) => (
              <span key={m.id} className="avatar avatar--sm" style={{
                marginLeft: i === 0 ? 0 : -8, border: '2px solid var(--abyss)',
                background: 'var(--lilac)', color: 'var(--deep-ink)', fontWeight: 700,
              }}>{m.initials.slice(0, 2)}</span>
            ))}
            {members.length > 5 && (
              <span className="avatar avatar--sm" style={{ marginLeft: -8, border: '2px solid var(--abyss)', background: 'rgba(225,191,145,.18)', color: 'var(--lilac)' }}>+{members.length - 5}</span>
            )}
          </div>
          <Link href={`/recipes/new`} className="btn btn--light-primary" style={{ textDecoration: 'none' }}>
            <IconPlus size={14} /> Submit resep
          </Link>
        </div>
      </div>

      {/* Tabs */}
      <div className="row gap-2" style={{ marginBottom: 28 }}>
        {TABS.map(({ key, label }) => (
          <button
            key={key}
            className={`tab${activeTab === key ? ' tab--active' : ''}`}
            onClick={() => setActiveTab(key)}
          >{label}</button>
        ))}
      </div>

      {/* Tab content */}
      {activeTab === 'resep' && (
        <TabResep
          recipes={recipes}
          isAdmin={isAdmin}
          onApprove={handleApprove}
          onReject={handleReject}
        />
      )}
      {activeTab === 'anggota' && (
        <TabAnggota
          members={members}
          isAdmin={isAdmin}
          isFounder={isFounder}
          groupId={id}
          onRemove={handleRemoveMember}
          onPromote={handlePromote}
          onDemote={handleDemote}
        />
      )}
      {activeTab === 'pengaturan' && (
        <TabPengaturan
          group={original}
          isAdmin={isAdmin}
          isFounder={isFounder}
          onDisband={handleDisband}
        />
      )}
    </main>
  )
}
