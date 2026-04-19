import '../styles/globals.css'
import ConditionalLayout from './ConditionalLayout'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Nisit Daeden - Admin Dashboard',
  description: 'Award Management System for Kasetsart University',
}

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (

    <ConditionalLayout>
      {children}
    </ConditionalLayout>
  )
}