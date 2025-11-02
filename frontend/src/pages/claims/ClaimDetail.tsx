import React, { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button, Space, Typography, Divider, message } from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import AppLayout from '@components/Layout/Layout'
import AddClaimItemModal from '@/components/AddClaimItemModal/AddClaimItemModal'
import CustomerInfo from '@/components/ClaimDetail/CustomerInfo'
import VehicleInfo from '@/components/ClaimDetail/VehicleInfo'
import ClaimInfo from '@/components/ClaimDetail/ClaimInfo'
import ClaimItemsTable from '@/components/ClaimDetail/ClaimItemsTable'
import ClaimAttachments from '@/components/ClaimDetail/ClaimAttachments'
import useClaimData from '@/hooks/useClaimData'
import useClaimPermissions from '@/hooks/useClaimPermissions'
import { getClaimsBasePath } from '@/utils/navigationHelpers'

const { Title } = Typography

const ClaimDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()

  // Modal state
  const [addItemModalVisible, setAddItemModalVisible] = useState(false)

  // Custom hooks
  const {
    claim,
    customer,
    vehicle,
    claimItems,
    attachments,
    partCategories,
    parts,
    claimLoading,
    customerLoading,
    vehicleLoading,
    itemsLoading,
    attachmentsLoading,
    refetchClaim,
    refetchClaimItems,
  } = useClaimData(id)

  const { canAddItems } = useClaimPermissions(claim)

  // Navigation handler
  const handleBack = () => {
    const location = window.location.pathname
    navigate(getClaimsBasePath(location))
  }

  // Modal handlers
  const handleOpenAddItemModal = () => {
    setAddItemModalVisible(true)
  }

  const handleCloseAddItemModal = () => {
    setAddItemModalVisible(false)
  }

  const handleAddItemSuccess = () => {
    // Refresh claim items and claim data
    refetchClaimItems()
    refetchClaim()
    message.success('Claim item added successfully')
    setAddItemModalVisible(false)
  }

  if (!id) {
    return (
      <AppLayout>
        <div style={{ textAlign: 'center', padding: '50px' }}>
          <Title level={3}>Invalid Claim ID</Title>
          <Button type="primary" onClick={handleBack}>
            Go Back
          </Button>
        </div>
      </AppLayout>
    )
  }

  return (
    <AppLayout>
      <div style={{ padding: '24px' }}>
        {/* Header */}
        <Space style={{ marginBottom: '24px' }}>
          <Button
            icon={<ArrowLeftOutlined />}
            onClick={handleBack}
            type="default"
          >
            Back to Claims
          </Button>
          <Title level={2} style={{ margin: 0 }}>
            Claim Details
          </Title>
        </Space>

        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          {/* Claim Information */}
          <ClaimInfo claim={claim} loading={claimLoading} />

          {/* Customer Information */}
          <CustomerInfo customer={customer} loading={customerLoading} />

          {/* Vehicle Information */}
          <VehicleInfo vehicle={vehicle} loading={vehicleLoading} />

          <Divider />

          {/* Claim Items */}
          <ClaimItemsTable
            claimItems={claimItems}
            parts={parts}
            partCategories={partCategories}
            loading={itemsLoading}
            canAddItems={canAddItems}
            onAddItem={handleOpenAddItemModal}
          />

          {/* Attachments */}
          <ClaimAttachments
            attachments={attachments}
            loading={attachmentsLoading}
          />
        </Space>

        {/* Add Item Modal */}
        {addItemModalVisible && id && (
          <AddClaimItemModal
            visible={addItemModalVisible}
            onCancel={handleCloseAddItemModal}
            onSuccess={handleAddItemSuccess}
            claimId={id}
            partCategories={partCategories}
          />
        )}
      </div>
    </AppLayout>
  )
}

export default ClaimDetail