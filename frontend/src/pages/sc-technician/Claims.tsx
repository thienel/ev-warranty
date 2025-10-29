import React from 'react'
import AppLayout from '@components/Layout/Layout.tsx'
import ClaimManagement from '@components/ClaimManagement/ClaimManagement.tsx'

const SCTechnicianClaims: React.FC = () => {
  return (
    <AppLayout title="Claim Management">
      <ClaimManagement />
    </AppLayout>
  )
}

export default SCTechnicianClaims
