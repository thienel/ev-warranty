import React from 'react'
import AppLayout from '@components/Layout/Layout'
import CustomerManagement from '@components/CustomerManagement/CustomerManagement'

const Customers: React.FC = () => {
  return (
    <AppLayout title="Customer Management">
      <CustomerManagement />
    </AppLayout>
  )
}

export default Customers
