import React from 'react'
import AppLayout from '@components/Layout/Layout.tsx'
import OfficeManagement from '@components/OfficeManagement/OfficeManagement.tsx'

const Offices: React.FC = () => {
  return (
    <AppLayout title="Office Management">
      <OfficeManagement />
    </AppLayout>
  )
}

export default Offices
