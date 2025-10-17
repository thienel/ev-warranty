import React from 'react'
import AppLayout from '@components/Layout/Layout.jsx'
import OfficeManagement from '@components/OfficeManagement/OfficeManagement.jsx'

const Offices = () => {
  return (
    <AppLayout title={'Office Management'}>
      <OfficeManagement />
    </AppLayout>
  )
}

export default Offices
