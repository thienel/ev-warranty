import React from 'react'
import AppLayout from '@components/Layout/Layout.tsx'
import PartCategoryManagement from '@components/PartCategoryManagement/PartCategoryManagement.tsx'

const PartCategories: React.FC = () => {
  return (
    <AppLayout title="Part Category Management">
      <PartCategoryManagement />
    </AppLayout>
  )
}

export default PartCategories
