import React from 'react'
import AppLayout from '@components/Layout/Layout.tsx'
import UserManagement from '@components/UserManagement/UserManagement.tsx'

const Users: React.FC = () => {
  return (
    <AppLayout title="User Management">
      <UserManagement />
    </AppLayout>
  )
}

export default Users
