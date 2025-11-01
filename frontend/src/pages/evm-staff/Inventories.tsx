import React from 'react'
import AppLayout from '@components/Layout/Layout.tsx'
import PartManagement from '@components/PartManagement/PartManagement.tsx'

const Inventories: React.FC = () => {
  return (
    <AppLayout title="Inventory Management">
      <PartManagement />
    </AppLayout>
  )
}

export default Inventories
