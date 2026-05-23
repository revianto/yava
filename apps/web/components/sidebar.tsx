'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { cn } from '@/lib/utils'
import { BookOpen, Compass, Users, Archive, Coffee } from 'lucide-react'

const navItems = [
  { href: '/recipes', label: 'Resep Saya', icon: BookOpen },
  { href: '/explore', label: 'Explore', icon: Compass },
  { href: '/groups', label: 'Grup', icon: Users },
  { href: '/archive', label: 'Arsip', icon: Archive },
]

export function Sidebar() {
  const pathname = usePathname()

  return (
    <aside className="hidden md:flex w-60 flex-col border-r border-border bg-sidebar shrink-0">
      <div className="flex items-center gap-2 px-6 py-5 border-b border-border">
        <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-primary text-primary-foreground">
          <Coffee className="w-4 h-4" />
        </div>
        <span className="font-semibold text-lg tracking-tight text-sidebar-foreground">YAVA</span>
      </div>

      <nav className="flex-1 px-3 py-4 space-y-1">
        {navItems.map(({ href, label, icon: Icon }) => (
          <Link
            key={href}
            href={href}
            className={cn(
              'flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors',
              pathname.startsWith(href)
                ? 'bg-sidebar-accent text-sidebar-foreground'
                : 'text-sidebar-foreground/60 hover:bg-sidebar-accent hover:text-sidebar-foreground'
            )}
          >
            <Icon className="w-4 h-4" />
            {label}
          </Link>
        ))}
      </nav>

      <div className="px-6 py-4 border-t border-border">
        <p className="text-xs text-muted-foreground">YAVA v0.1.0</p>
      </div>
    </aside>
  )
}
