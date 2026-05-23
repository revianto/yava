'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { Plus, Coffee } from 'lucide-react'
import { Button, buttonVariants } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { mockUser } from '@/lib/mock-data'
import { cn } from '@/lib/utils'

export function Navbar() {
  const router = useRouter()

  return (
    <header className="flex items-center justify-between px-4 md:px-6 py-3 border-b border-border bg-background">
      <div className="flex md:hidden items-center gap-2">
        <div className="flex items-center justify-center w-7 h-7 rounded-lg bg-primary text-primary-foreground">
          <Coffee className="w-3.5 h-3.5" />
        </div>
        <span className="font-semibold tracking-tight">YAVA</span>
      </div>

      <div className="hidden md:block" />

      <div className="flex items-center gap-3">
        <Link href="/recipes/new" className={cn(buttonVariants({ size: 'sm' }), 'gap-1.5')}>
          <Plus className="w-4 h-4" />
          <span className="hidden sm:inline">Resep Baru</span>
        </Link>

        <DropdownMenu>
          <DropdownMenuTrigger className="rounded-full outline-none focus-visible:ring-2 focus-visible:ring-ring">
            <Avatar className="w-8 h-8">
              <AvatarImage src={mockUser.avatarUrl} alt={mockUser.name} />
              <AvatarFallback className="bg-accent text-accent-foreground text-xs font-medium">
                {mockUser.name.slice(0, 2).toUpperCase()}
              </AvatarFallback>
            </Avatar>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-48">
            <div className="px-2 py-1.5">
              <p className="text-sm font-medium leading-none">{mockUser.name}</p>
              <p className="text-xs text-muted-foreground mt-0.5">{mockUser.email}</p>
            </div>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={() => router.push('/profile')}>Profil</DropdownMenuItem>
            <DropdownMenuItem onClick={() => router.push('/settings')}>Pengaturan</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="text-destructive focus:text-destructive">
              Keluar
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  )
}
