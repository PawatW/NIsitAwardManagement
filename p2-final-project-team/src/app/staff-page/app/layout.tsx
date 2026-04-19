import type { Metadata } from 'next'
import '../styles/dashboard.css'

export const metadata: Metadata = {
  title: 'Nisit D-Den ',
  description: 'ระบบพิจารณารางวัลนิสิตดีเด่น',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="staff-page-wrapper">
      {children}
    </div>
  )
}

