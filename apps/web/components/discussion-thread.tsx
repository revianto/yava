'use client'

import { useState } from 'react'
import type { Discussion } from '@/types'

function PinIcon({ pinned }: { pinned: boolean }) {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill={pinned ? 'currentColor' : 'none'} stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round">
      <path d="M12 2L8 8H3l4 4-2 6 7-4 7 4-2-6 4-4h-5z" />
    </svg>
  )
}

interface Props {
  recipeId: string
  initialDiscussions: Discussion[]
  isAdmin?: boolean
}

export function DiscussionThread({ recipeId: _recipeId, initialDiscussions, isAdmin = false }: Props) {
  const [discussions, setDiscussions] = useState(initialDiscussions)
  const [newComment, setNewComment] = useState('')
  const [replyingTo, setReplyingTo] = useState<string | null>(null)
  const [replyText, setReplyText] = useState('')
  const [expandedReplies, setExpandedReplies] = useState<Set<string>>(new Set())

  const pinnedCount = discussions.filter((d) => d.pinned).length

  const handlePost = () => {
    if (!newComment.trim()) return
    const d: Discussion = {
      id: `d-${Date.now()}`,
      recipeId: _recipeId,
      authorName: 'Kamu',
      authorInitials: 'ND',
      content: newComment.trim(),
      createdAt: 'baru saja',
      pinned: false,
      replies: [],
    }
    setDiscussions((prev) => [...prev, d])
    setNewComment('')
  }

  const handleReply = (discussionId: string) => {
    if (!replyText.trim()) return
    setDiscussions((prev) => prev.map((d) => {
      if (d.id !== discussionId) return d
      return {
        ...d,
        replies: [...d.replies, {
          id: `r-${Date.now()}`,
          authorName: 'Kamu',
          authorInitials: 'ND',
          content: replyText.trim(),
          createdAt: 'baru saja',
        }],
      }
    }))
    setReplyText('')
    setReplyingTo(null)
    setExpandedReplies((prev) => new Set([...prev, discussionId]))
  }

  const togglePin = (discussionId: string) => {
    setDiscussions((prev) => prev.map((d) =>
      d.id === discussionId ? { ...d, pinned: !d.pinned } : d
    ))
  }

  const sorted = [...discussions].sort((a, b) => {
    if (a.pinned && !b.pinned) return -1
    if (!a.pinned && b.pinned) return 1
    return 0
  })

  return (
    <div className="card" style={{ padding: 24, marginTop: 16 }}>
      <div className="row between" style={{ marginBottom: 16 }}>
        <div className="t-h3">Diskusi grup</div>
        <span className="t-small muted">{discussions.length} komentar{pinnedCount > 0 && ` · ${pinnedCount} disematkan`}</span>
      </div>

      {/* Comment list */}
      <div>
        {sorted.map((d) => {
          const showReplies = expandedReplies.has(d.id)
          return (
            <div key={d.id} style={{ borderTop: '1px solid var(--hairline)', paddingTop: 14, paddingBottom: 4, marginBottom: 4 }}>
              {/* Main comment */}
              <div className="row gap-2" style={{ alignItems: 'flex-start' }}>
                <span className="avatar avatar--sm" style={{ flexShrink: 0 }}>{d.authorInitials}</span>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div className="row gap-2" style={{ alignItems: 'center', flexWrap: 'wrap', marginBottom: 4 }}>
                    <span className="t-small" style={{ fontWeight: 700 }}>{d.authorName}</span>
                    <span className="t-small muted">{d.createdAt}</span>
                    {d.pinned && (
                      <span className="tag tag--default" style={{ fontSize: 10 }}>DISEMATKAN</span>
                    )}
                  </div>
                  <div className="t-small" style={{ color: 'var(--deep-ink)', lineHeight: 1.55 }}>{d.content}</div>

                  {/* Actions */}
                  <div className="row gap-3" style={{ marginTop: 8, alignItems: 'center' }}>
                    <button
                      onClick={() => { setReplyingTo(replyingTo === d.id ? null : d.id); setReplyText('') }}
                      style={{ background: 'none', border: 0, cursor: 'pointer', fontSize: 12, fontWeight: 600, color: 'var(--muted)', padding: 0 }}
                    >
                      Balas
                    </button>
                    {d.replies.length > 0 && (
                      <button
                        onClick={() => setExpandedReplies((prev) => {
                          const s = new Set(prev)
                          s.has(d.id) ? s.delete(d.id) : s.add(d.id)
                          return s
                        })}
                        style={{ background: 'none', border: 0, cursor: 'pointer', fontSize: 12, fontWeight: 600, color: 'var(--muted)', padding: 0 }}
                      >
                        {showReplies ? 'Sembunyikan' : `${d.replies.length} balasan`}
                      </button>
                    )}
                    {isAdmin && (
                      <button
                        onClick={() => togglePin(d.id)}
                        style={{
                          background: 'none', border: 0, cursor: 'pointer', fontSize: 12, fontWeight: 600, padding: 0,
                          color: d.pinned ? 'var(--coral-red)' : 'var(--muted)',
                          display: 'inline-flex', alignItems: 'center', gap: 4,
                        }}
                        title={d.pinned ? 'Lepas sematan' : 'Sematkan'}
                      >
                        <PinIcon pinned={d.pinned} /> {d.pinned ? 'Lepas sematan' : 'Sematkan'}
                      </button>
                    )}
                  </div>
                </div>
              </div>

              {/* Replies */}
              {showReplies && d.replies.length > 0 && (
                <div style={{ marginLeft: 40, marginTop: 10, borderLeft: '2px solid var(--hairline)', paddingLeft: 16 }}>
                  {d.replies.map((r) => (
                    <div key={r.id} className="row gap-2" style={{ alignItems: 'flex-start', marginBottom: 10 }}>
                      <span className="avatar avatar--sm" style={{ flexShrink: 0, width: 26, height: 26, fontSize: 10 }}>{r.authorInitials}</span>
                      <div>
                        <div className="row gap-2" style={{ alignItems: 'center', marginBottom: 3 }}>
                          <span className="t-small" style={{ fontWeight: 700 }}>{r.authorName}</span>
                          <span className="t-small muted">{r.createdAt}</span>
                        </div>
                        <div className="t-small" style={{ color: 'var(--deep-ink)', lineHeight: 1.55 }}>{r.content}</div>
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {/* Reply input */}
              {replyingTo === d.id && (
                <div className="row gap-2" style={{ marginLeft: 40, marginTop: 10 }}>
                  <span className="avatar avatar--sm" style={{ flexShrink: 0, width: 26, height: 26, fontSize: 10 }}>ND</span>
                  <div style={{ flex: 1, display: 'flex', gap: 8 }}>
                    <input
                      className="input"
                      value={replyText}
                      onChange={(e) => setReplyText(e.target.value)}
                      placeholder={`Balas ${d.authorName}…`}
                      style={{ flex: 1, height: 36, fontSize: 13 }}
                      onKeyDown={(e) => e.key === 'Enter' && handleReply(d.id)}
                      autoFocus
                    />
                    <button
                      className="btn btn--primary"
                      onClick={() => handleReply(d.id)}
                      disabled={!replyText.trim()}
                      style={{ height: 36, padding: '0 16px', fontSize: 13, opacity: replyText.trim() ? 1 : 0.4 }}
                    >
                      Kirim
                    </button>
                  </div>
                </div>
              )}
            </div>
          )
        })}
      </div>

      {/* New comment input */}
      <div className="row gap-2" style={{ marginTop: 16, paddingTop: 14, borderTop: '1px solid var(--hairline)' }}>
        <span className="avatar avatar--sm" style={{ flexShrink: 0 }}>ND</span>
        <div style={{ flex: 1, display: 'flex', gap: 8 }}>
          <input
            className="input"
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            placeholder="Tulis komentar…"
            style={{ flex: 1, fontSize: 13 }}
            onKeyDown={(e) => e.key === 'Enter' && handlePost()}
          />
          <button
            className="btn btn--primary"
            onClick={handlePost}
            disabled={!newComment.trim()}
            style={{ opacity: newComment.trim() ? 1 : 0.4 }}
          >
            Kirim
          </button>
        </div>
      </div>
    </div>
  )
}
