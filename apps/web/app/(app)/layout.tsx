import { Topnav } from '@/components/topnav'

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="page">
      <Topnav />
      {children}
    </div>
  )
}
