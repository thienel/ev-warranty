import React from 'react'
import { message } from 'antd'
import { API_ENDPOINTS } from '@constants/common-constants.js'
import useManagement from '@/hooks/useManagement.js'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar.jsx'
import GenericTable from '@components/common/GenericTable/GenericTable.jsx'
import GenerateColumns from './claimTableColumns.jsx'

const ClaimManagement = () => {
  const {
    loading,
    setLoading,
    searchText,
    setSearchText,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.CLAIM, 'claim')

  // Mock data for demonstration since APIs are not available yet
  const mockClaims = [
    {
      id: 'claim-001-abc-def',
      status: 'SUBMITTED',
      customer_id: 'cust-123-xyz',
      customer_name: 'John Smith',
      vehicle_id: 'veh-456-abc',
      vehicle_info: 'Tesla Model 3 2023',
      description: 'Battery replacement issue',
      total_cost: 2500.00,
      created_at: '2024-10-20T10:30:00Z',
      updated_at: '2024-10-20T10:30:00Z',
    },
    {
      id: 'claim-002-def-ghi',
      status: 'APPROVED',
      customer_id: 'cust-124-xyz',
      customer_name: 'Sarah Johnson',
      vehicle_id: 'veh-457-def',
      vehicle_info: 'Nissan Leaf 2022',
      description: 'Motor controller malfunction',
      total_cost: 1800.00,
      created_at: '2024-10-19T14:15:00Z',
      updated_at: '2024-10-21T09:20:00Z',
    },
    {
      id: 'claim-003-ghi-jkl',
      status: 'REVIEWING',
      customer_id: 'cust-125-xyz',
      customer_name: 'Michael Davis',
      vehicle_id: 'veh-458-ghi',
      vehicle_info: 'BMW iX 2023',
      description: 'Charging port repair',
      total_cost: 950.00,
      created_at: '2024-10-18T11:45:00Z',
      updated_at: '2024-10-22T16:30:00Z',
    },
    {
      id: 'claim-004-jkl-mno',
      status: 'PARTIALLY_APPROVED',
      customer_id: 'cust-126-xyz',
      customer_name: 'Emily Wilson',
      vehicle_id: 'veh-459-jkl',
      vehicle_info: 'Audi e-tron 2022',
      description: 'Software update and diagnostics',
      total_cost: 350.00,
      created_at: '2024-10-15T08:20:00Z',
      updated_at: '2024-10-20T12:00:00Z',
    },
    {
      id: 'claim-005-mno-pqr',
      status: 'REJECTED',
      customer_id: 'cust-127-xyz',
      customer_name: 'Robert Brown',
      vehicle_id: 'veh-460-mno',
      vehicle_info: 'Ford Mustang Mach-E 2023',
      description: 'Cosmetic damage claim',
      total_cost: 0.00,
      created_at: '2024-10-14T13:10:00Z',
      updated_at: '2024-10-16T10:45:00Z',
    },
    {
      id: 'claim-006-pqr-stu',
      status: 'DRAFT',
      customer_id: 'cust-128-xyz',
      customer_name: 'Lisa Anderson',
      vehicle_id: 'veh-461-pqr',
      vehicle_info: 'Hyundai Ioniq 5 2023',
      description: 'Brake system inspection needed',
      total_cost: 0.00,
      created_at: '2024-10-23T09:15:00Z',
      updated_at: '2024-10-23T09:15:00Z',
    },
    {
      id: 'claim-007-stu-vwx',
      status: 'REQUEST_INFO',
      customer_id: 'cust-129-xyz',
      customer_name: 'David Kim',
      vehicle_id: 'veh-462-stu',
      vehicle_info: 'Kia EV6 2022',
      description: 'Infotainment system malfunction',
      total_cost: 1200.00,
      created_at: '2024-10-17T16:30:00Z',
      updated_at: '2024-10-19T11:20:00Z',
    },
    {
      id: 'claim-008-vwx-yza',
      status: 'CANCELLED',
      customer_id: 'cust-130-xyz',
      customer_name: 'Maria Garcia',
      vehicle_id: 'veh-463-vwx',
      vehicle_info: 'Volkswagen ID.4 2023',
      description: 'Cancelled due to customer request',
      total_cost: 0.00,
      created_at: '2024-10-12T14:45:00Z',
      updated_at: '2024-10-13T10:30:00Z',
    },
  ]

  const handleViewDetails = (claim) => {
    message.info(`Viewing details for claim ${claim.id.slice(0, 8)}. Detail page will be implemented later.`)
  }

  const handleCreateClaim = () => {
    message.info('Create claim page will be implemented later.')
  }

  const searchFields = ['customer_name', 'vehicle_info', 'description', 'status']

  return (
    <>
      <GenericActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleCreateClaim}
        loading={loading}
        searchPlaceholder="Search by customer name, vehicle, description, or status..."
        addButtonText="Create Claim"
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={mockClaims} // Using mock data instead of claims until APIs are available
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFields}
        deleteEndpoint={API_ENDPOINTS.CLAIM}
        deleteSuccessMessage="Claim deleted successfully"
        deleteErrorMessage="Failed to delete claim"
        additionalProps={{ onViewDetails: handleViewDetails }}
        showDeleteButton={false} // Hide delete button since you mentioned no modals
      />
    </>
  )
}

export default ClaimManagement
