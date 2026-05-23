'use client'

import { use, useState, useEffect, useRef, useCallback } from 'react'
import Link from 'next/link'
import { notFound } from 'next/navigation'
import { ArrowLeft, Pause, Play, RotateCcw, SkipForward, CheckCircle } from 'lucide-react'
import { Button, buttonVariants } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { mockRecipes } from '@/lib/mock-data'
import { cn } from '@/lib/utils'

type BrewState = 'prep' | 'brewing' | 'paused' | 'done'

export default function BrewPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params)
  const recipe = mockRecipes.find((r) => r.id === id)
  if (!recipe || recipe.sessions.length === 0) notFound()

  const sessions = recipe.sessions
  const [state, setState] = useState<BrewState>('prep')
  const [sessionIdx, setSessionIdx] = useState(0)
  const [elapsed, setElapsed] = useState(0)
  const [prepCountdown, setPrepCountdown] = useState(3)
  const startRef = useRef<number | null>(null)
  const rafRef = useRef<number | null>(null)

  const currentSession = sessions[sessionIdx]
  const progress = currentSession ? Math.min((elapsed / currentSession.durationSeconds) * 100, 100) : 0

  const tick = useCallback(() => {
    if (startRef.current === null) return
    const now = performance.now()
    const elapsedSec = (now - startRef.current) / 1000
    setElapsed(elapsedSec)

    if (elapsedSec >= currentSession.durationSeconds) {
      if (sessionIdx < sessions.length - 1) {
        setSessionIdx((i) => i + 1)
        setElapsed(0)
        startRef.current = performance.now()
      } else {
        setState('done')
        startRef.current = null
        return
      }
    }
    rafRef.current = requestAnimationFrame(tick)
  }, [currentSession, sessionIdx, sessions.length])

  useEffect(() => {
    if (state !== 'prep') return
    const interval = setInterval(() => {
      setPrepCountdown((c) => {
        if (c <= 1) {
          clearInterval(interval)
          setState('brewing')
          startRef.current = performance.now()
          return 0
        }
        return c - 1
      })
    }, 1000)
    return () => clearInterval(interval)
  }, [state])

  useEffect(() => {
    if (state === 'brewing') {
      rafRef.current = requestAnimationFrame(tick)
    }
    return () => {
      if (rafRef.current) cancelAnimationFrame(rafRef.current)
    }
  }, [state, tick])

  const pause = () => {
    setState('paused')
    if (rafRef.current) cancelAnimationFrame(rafRef.current)
  }

  const resume = () => {
    startRef.current = performance.now() - elapsed * 1000
    setState('brewing')
  }

  const skip = () => {
    if (sessionIdx < sessions.length - 1) {
      setSessionIdx((i) => i + 1)
      setElapsed(0)
      startRef.current = performance.now()
      setState('brewing')
    } else {
      setState('done')
    }
  }

  const reset = () => {
    setSessionIdx(0)
    setElapsed(0)
    setPrepCountdown(3)
    startRef.current = null
    setState('prep')
  }

  const formatTime = (sec: number) => {
    const m = Math.floor(sec / 60)
    const s = Math.floor(sec % 60)
    return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
  }

  const remaining = Math.max(0, currentSession?.durationSeconds - elapsed)

  if (state === 'done') {
    const total = recipe.sessions.reduce((a, s) => a + s.durationSeconds, 0)
    return (
      <div className="h-full flex flex-col items-center justify-center p-6 text-center space-y-6">
        <CheckCircle className="w-16 h-16 text-primary" />
        <div>
          <h2 className="text-2xl font-bold">Brewing Selesai!</h2>
          <p className="text-muted-foreground mt-1">Total waktu: {formatTime(total)}</p>
        </div>
        <p className="text-sm text-muted-foreground max-w-xs">
          Kopi kamu siap. Nikmati setiap tegukan ☕
        </p>
        <div className="flex gap-3">
          <Button variant="outline" onClick={reset} className="gap-2">
            <RotateCcw className="w-4 h-4" />
            Ulangi
          </Button>
          <Link href={`/recipes/${id}`} className={buttonVariants()}>
            Kembali ke Resep
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="h-full flex flex-col p-4 md:p-6 max-w-lg mx-auto">
      <div className="flex items-center gap-2 mb-6">
        <Link
          href={`/recipes/${id}`}
          className={cn(buttonVariants({ variant: 'ghost', size: 'sm' }), 'gap-2 text-muted-foreground')}
        >
          <ArrowLeft className="w-4 h-4" />
          Batalkan
        </Link>
      </div>

      <div className="flex-1 flex flex-col justify-center space-y-8">
        <div className="text-center">
          <p className="text-xs text-muted-foreground mb-1">{recipe.typeName}</p>
          <h1 className="text-xl font-semibold">{recipe.title}</h1>
        </div>

        {state === 'prep' && (
          <div className="text-center space-y-4">
            <div className="text-7xl font-mono font-bold text-primary">{prepCountdown}</div>
            <p className="text-muted-foreground text-sm">Siapkan peralatan Anda...</p>
          </div>
        )}

        {(state === 'brewing' || state === 'paused') && (
          <div className="space-y-6">
            <div className="text-center space-y-1">
              <div className="flex items-center justify-center gap-2 text-xs text-muted-foreground">
                <span>Session {sessionIdx + 1} / {sessions.length}</span>
              </div>
              <h2 className="text-3xl font-bold">{currentSession.label}</h2>
              {currentSession.waterMl && (
                <p className="text-muted-foreground text-sm">{currentSession.waterMl}ml air</p>
              )}
              {currentSession.notes && (
                <p className="text-xs text-muted-foreground mt-1 italic">{currentSession.notes}</p>
              )}
            </div>

            <div className="text-center">
              <div className={cn(
                'text-6xl font-mono font-bold transition-colors',
                state === 'paused' ? 'text-muted-foreground' : 'text-foreground'
              )}>
                {formatTime(remaining)}
              </div>
              <p className="text-xs text-muted-foreground mt-1">tersisa</p>
            </div>

            <Progress value={progress} className="h-2" />

            <div className="flex justify-center gap-1.5">
              {sessions.map((s, i) => (
                <div
                  key={s.id}
                  className={cn(
                    'h-1.5 rounded-full transition-all',
                    i < sessionIdx ? 'bg-primary w-6' :
                    i === sessionIdx ? 'bg-primary/70 w-10' :
                    'bg-muted w-6'
                  )}
                />
              ))}
            </div>
          </div>
        )}

        {(state === 'brewing' || state === 'paused') && (
          <div className="flex items-center justify-center gap-4">
            <Button variant="outline" size="icon" onClick={reset} className="w-10 h-10">
              <RotateCcw className="w-4 h-4" />
            </Button>

            <Button
              size="icon"
              className="w-16 h-16 rounded-full shadow-lg"
              onClick={state === 'brewing' ? pause : resume}
            >
              {state === 'brewing' ? (
                <Pause className="w-6 h-6" />
              ) : (
                <Play className="w-6 h-6" />
              )}
            </Button>

            <Button variant="outline" size="icon" onClick={skip} className="w-10 h-10">
              <SkipForward className="w-4 h-4" />
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}
