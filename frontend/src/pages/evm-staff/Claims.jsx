import React from 'react'
import AppLayout from '@components/Layout/Layout.jsx'
import ClaimManagement from '@components/ClaimManagement/ClaimManagement.jsx'

const Claims = () => {
  return (
    <AppLayout title={'Claim Management'}>
      <ClaimManagement />
    </AppLayout>
  )
}

export default Claims
