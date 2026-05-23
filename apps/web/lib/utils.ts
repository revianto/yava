import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  if (m === 0) return `${s}d`
  if (s === 0) return `${m}m`
  return `${m}m ${s}d`
}

export function formatDate(isoString: string): string {
  return new Intl.DateTimeFormat('id-ID', { dateStyle: 'medium' }).format(new Date(isoString))
}
