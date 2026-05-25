'use client'

import { useState, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { GROUPS } from '@/lib/mock-data'
import { IconArrowLeft, IconCheck } from '@/components/icons'

export default function JoinGroupPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [code, setCode] = useState(searchParams.get('code') ?? '')
  const [status, setStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle')
  const [foundGroup, setFoundGroup] = useState<(typeof GROUPS)[0] | null>(null)

  useEffect(() => {
    const c = searchParams.get('code')
    if (c) setCode(c.toUpperCase())
  }, [searchParams])

  const handleJoin = () => {
    if (!code.trim()) return
    setStatus('loading')
    setTimeout(() => {
      const match = GROUPS.find((g) => g.inviteCode === code.trim().toUpperCase())
      if (match) {
        setFoundGroup(match)
        setStatus('success')
        setTimeout(() => router.push(`/groups/${match.id}`), 1200)
      } else {
        setStatus('error')
      }
    }, 600)
  }

  return (
    <main className="container" style={{ paddingBottom: 64 }}>
      <div style={{ maxWidth: 480, margin: '0 auto', paddingTop: 32 }}>
        <button
          onClick={() => router.back()}
          style={{ background: 'transparent', border: 0, cursor: 'pointer', display: 'inline-flex', alignItems: 'center', gap: 6, color: 'var(--muted)', fontWeight: 500, fontSize: 13, padding: 0, marginBottom: 32 }}
        >
          <IconArrowLeft size={14} /> Kembali
        </button>

        <div style={{ marginBottom: 36 }}>
          <div className="t-label muted" style={{ marginBottom: 8 }}>Bergabung</div>
          <div className="t-h1">Gabung grup</div>
          <div className="t-body muted" style={{ marginTop: 8 }}>
            Masukkan invite code yang kamu terima dari admin grup.
          </div>
        </div>

        {status === 'success' && foundGroup ? (
          <div className="card card--dark" style={{ padding: 28, textAlign: 'center' }}>
            <div style={{
              width: 56, height: 56, borderRadius: '50%',
              background: 'var(--coral-red)', color: '#FFF8E7',
              display: 'inline-flex', alignItems: 'center', justifyContent: 'center',
              marginBottom: 16,
            }}>
              <IconCheck size={24} />
            </div>
            <div className="t-h2" style={{ color: '#FFF8E7', marginBottom: 6 }}>Berhasil bergabung!</div>
            <div className="t-small muted-dark">Mengalihkan ke <strong style={{ color: '#FFF8E7' }}>{foundGroup.name}</strong>…</div>
          </div>
        ) : (
          <div className="col gap-3">
            <div>
              <div className="t-label muted" style={{ marginBottom: 8 }}>Invite code</div>
              <input
                className="input"
                value={code}
                onChange={(e) => { setCode(e.target.value.toUpperCase()); setStatus('idle') }}
                placeholder="e.g. SNYN2026"
                style={{ width: '100%', fontFamily: 'monospace', letterSpacing: '.08em', fontWeight: 700, fontSize: 16 }}
                onKeyDown={(e) => e.key === 'Enter' && handleJoin()}
              />
              {status === 'error' && (
                <div className="t-small" style={{ color: 'var(--coral-red)', marginTop: 6 }}>
                  Kode tidak ditemukan. Pastikan kode benar dan coba lagi.
                </div>
              )}
            </div>

            <button
              className="btn btn--primary btn--block btn--lg"
              onClick={handleJoin}
              disabled={!code.trim() || status === 'loading'}
              style={{ opacity: code.trim() && status !== 'loading' ? 1 : .4 }}
            >
              {status === 'loading' ? 'Mencari grup…' : 'Gabung sekarang'}
            </button>

            <div className="divider" style={{ margin: '8px 0' }} />

            <div className="t-caption muted" style={{ textAlign: 'center' }}>
              Belum punya kode? Minta admin grup untuk membagikan invite link.
            </div>
          </div>
        )}
      </div>
    </main>
  )
}
