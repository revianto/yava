'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { IconArrowLeft } from '@/components/icons'

export default function NewGroupPage() {
  const router = useRouter()
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [submitted, setSubmitted] = useState(false)

  const canSubmit = name.trim().length >= 3

  const handleSubmit = () => {
    if (!canSubmit) return
    setSubmitted(true)
    // mock: redirect to first group (in real app would create and redirect to new group)
    setTimeout(() => router.push('/groups/g-senayan'), 800)
  }

  return (
    <main className="container" style={{ paddingBottom: 64 }}>
      <div style={{ maxWidth: 560, margin: '0 auto', paddingTop: 32 }}>
        <button
          onClick={() => router.back()}
          style={{ background: 'transparent', border: 0, cursor: 'pointer', display: 'inline-flex', alignItems: 'center', gap: 6, color: 'var(--muted)', fontWeight: 500, fontSize: 13, padding: 0, marginBottom: 32 }}
        >
          <IconArrowLeft size={14} /> Batal
        </button>

        <div style={{ marginBottom: 36 }}>
          <div className="t-label muted" style={{ marginBottom: 8 }}>Grup baru</div>
          <div className="t-h1">Buat grup kopi</div>
          <div className="t-body muted" style={{ marginTop: 8 }}>
            Undang teman untuk berbagi, mendiskusikan, dan menyetujui resep bersama.
          </div>
        </div>

        <div className="col gap-3">
          <div>
            <div className="t-label muted" style={{ marginBottom: 8 }}>Nama grup *</div>
            <input
              className="input"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. Komunitas Senayan"
              style={{ width: '100%' }}
              maxLength={60}
            />
            <div className="t-caption" style={{ marginTop: 6, textAlign: 'right' }}>{name.length}/60</div>
          </div>

          <div>
            <div className="t-label muted" style={{ marginBottom: 8 }}>Deskripsi</div>
            <textarea
              className="input"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Ceritakan fokus atau tema grup ini..."
              rows={3}
              style={{ width: '100%', height: 'auto', resize: 'vertical', paddingTop: 10, paddingBottom: 10 }}
              maxLength={200}
            />
            <div className="t-caption" style={{ marginTop: 6, textAlign: 'right' }}>{description.length}/200</div>
          </div>

          {/* Info box */}
          <div className="card" style={{ padding: 20, background: 'rgba(45,82,74,.05)', border: '1px solid var(--hairline)' }}>
            <div className="col gap-2">
              {[
                'Kamu otomatis jadi Founder dan bisa promote anggota ke Admin.',
                'Invite code unik akan dibuat otomatis — bagikan ke teman.',
                'Resep yang disubmit anggota perlu persetujuan Admin sebelum aktif.',
              ].map((line, i) => (
                <div key={i} className="row gap-2" style={{ alignItems: 'flex-start' }}>
                  <span style={{ width: 18, height: 18, borderRadius: '50%', background: 'var(--deep-ink)', color: 'var(--lilac)', display: 'inline-flex', alignItems: 'center', justifyContent: 'center', fontWeight: 700, fontSize: 10, flex: 'none', marginTop: 1 }}>{i + 1}</span>
                  <span className="t-small muted">{line}</span>
                </div>
              ))}
            </div>
          </div>

          <div className="row between" style={{ marginTop: 8 }}>
            <button className="btn btn--secondary" onClick={() => router.back()}>Batal</button>
            <button
              className="btn btn--primary"
              onClick={handleSubmit}
              disabled={!canSubmit || submitted}
              style={{ opacity: canSubmit && !submitted ? 1 : .4 }}
            >
              {submitted ? 'Membuat grup…' : 'Buat grup'}
            </button>
          </div>
        </div>
      </div>
    </main>
  )
}
