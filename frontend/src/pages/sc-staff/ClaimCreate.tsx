import React from 'react'
import AppLayout from '@components/Layout/Layout'
import ClaimCreateComponent from '@components/ClaimManagement/ClaimCreate/ClaimCreate'

const ClaimCreate: React.FC = () => {
  return (
    <AppLayout title="Create New Claim">
      <ClaimCreateComponent />
    </AppLayout>
  )
}

export default ClaimCreate
