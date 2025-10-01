import React from 'react'
import AppLayout from '@components/Layout/Layout.jsx'
import UserManagement from '@components/UserManagement/UserManagement.jsx'

const Users = () => {
  return (
    <AppLayout title={'Users Management'}>
      <UserManagement />
    </AppLayout>
  )
}

export default Users
