import React, { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button, Space, Typography, Divider, message, Modal } from 'antd'
import {
  ArrowLeftOutlined,
  SendOutlined,
  DeleteOutlined,
  StopOutlined,
  PlayCircleOutlined,
  QuestionCircleOutlined,
  CheckCircleOutlined,
} from '@ant-design/icons'
import AppLayout from '@components/Layout/Layout'
import AddClaimItemModal from '@/components/AddClaimItemModal/AddClaimItemModal'
import AddAttachmentModal from '@/components/AddAttachmentModal/AddAttachmentModal'
import CustomerInfo from '@/components/ClaimDetail/CustomerInfo'
import VehicleInfo from '@/components/ClaimDetail/VehicleInfo'
import ClaimInfo from '@/components/ClaimDetail/ClaimInfo'
import ClaimItemsTable from '@/components/ClaimDetail/ClaimItemsTable'
import ClaimAttachments from '@/components/ClaimDetail/ClaimAttachments'
import WarrantyPolicyCard from '@/components/ClaimDetail/WarrantyPolicyCard'
import PolicyCoverageModal from '@/components/ClaimDetail/PolicyCoverageModal'
import useClaimData from '@/hooks/useClaimData'
import useClaimPermissions from '@/hooks/useClaimPermissions'
import { getClaimsBasePath } from '@/utils/navigationHelpers'
import { claimsApi, claimItemsApi } from '@/services/claimsApi'

const { Title } = Typography

const ClaimDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()

  // Modal state
  const [addItemModalVisible, setAddItemModalVisible] = useState(false)
  const [addAttachmentModalVisible, setAddAttachmentModalVisible] = useState(false)
  const [policyCoverageModalVisible, setPolicyCoverageModalVisible] = useState(false)
  const [selectedCategoryId, setSelectedCategoryId] = useState<number | null>(null)
  const [selectedCategoryName, setSelectedCategoryName] = useState<string>('')
  const [submitLoading, setSubmitLoading] = useState(false)
  const [deleteLoading, setDeleteLoading] = useState(false)
  const [cancelLoading, setCancelLoading] = useState(false)
  const [startReviewLoading, setStartReviewLoading] = useState(false)
  const [requestInfoLoading, setRequestInfoLoading] = useState(false)
  const [completeLoading, setCompleteLoading] = useState(false)

  // Custom hooks
  const {
    claim,
    customer,
    vehicle,
    claimItems,
    attachments,
    partCategories,
    parts,
    warrantyPolicy,
    claimLoading,
    customerLoading,
    vehicleLoading,
    itemsLoading,
    attachmentsLoading,
    warrantyPolicyLoading,
    refetchClaim,
    refetchClaimItems,
    refetchAttachments,
  } = useClaimData(id)

  const {
    canAddItems,
    canAddAttachments,
    canSubmitClaim,
    canDeleteClaim,
    canCancelClaim,
    canStartReview,
    canRequestMoreInfo,
    canApproveClaimItems,
    canRejectClaimItems,
    canCompleteClaim,
    canViewWarrantyPolicy,
    canViewPolicyCoverage,
  } = useClaimPermissions(claim, claimItems)

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

  const handleOpenAddAttachmentModal = () => {
    setAddAttachmentModalVisible(true)
  }

  const handleCloseAddAttachmentModal = () => {
    setAddAttachmentModalVisible(false)
  }

  const handleAddAttachmentSuccess = () => {
    // Refresh attachments
    refetchAttachments()
    message.success('Attachments uploaded successfully')
    setAddAttachmentModalVisible(false)
  }

  const handleSubmitClaim = async () => {
    if (!id) return

    try {
      setSubmitLoading(true)
      await claimsApi.submit(id)
      message.success('Claim submitted successfully')
      refetchClaim() // Refresh claim to update status
    } catch (error) {
      console.error('Error submitting claim:', error)
      message.error('Failed to submit claim')
    } finally {
      setSubmitLoading(false)
    }
  }

  const handleDeleteClaim = () => {
    if (!id) return

    Modal.confirm({
      title: 'Delete Claim',
      content: 'Are you sure you want to delete this claim? This action cannot be undone.',
      okText: 'Delete',
      cancelText: 'Cancel',
      okType: 'danger',
      onOk: async () => {
        try {
          setDeleteLoading(true)
          await claimsApi.delete(id)
          message.success('Claim deleted successfully')
          navigate(getClaimsBasePath(window.location.pathname))
        } catch (error) {
          console.error('Error deleting claim:', error)
          message.error('Failed to delete claim')
        } finally {
          setDeleteLoading(false)
        }
      },
    })
  }

  const handleCancelClaim = () => {
    if (!id) return

    Modal.confirm({
      title: 'Cancel Claim',
      content: 'Are you sure you want to cancel this claim?',
      okText: 'Cancel Claim',
      cancelText: 'Keep Claim',
      okType: 'danger',
      onOk: async () => {
        try {
          setCancelLoading(true)
          await claimsApi.cancel(id)
          message.success('Claim cancelled successfully')
          refetchClaim() // Refresh claim to update status
        } catch (error) {
          console.error('Error cancelling claim:', error)
          message.error('Failed to cancel claim')
        } finally {
          setCancelLoading(false)
        }
      },
    })
  }

  const handleStartReview = async () => {
    if (!id) return

    try {
      setStartReviewLoading(true)
      await claimsApi.review(id)
      message.success('Review started successfully')
      refetchClaim() // Refresh claim to update status
    } catch (error) {
      console.error('Error starting review:', error)
      message.error('Failed to start review')
    } finally {
      setStartReviewLoading(false)
    }
  }

  const handleRequestMoreInfo = async () => {
    if (!id) return

    try {
      setRequestInfoLoading(true)
      await claimsApi.requestInfo(id)
      message.success('More information requested successfully')
      refetchClaim() // Refresh claim to update status
    } catch (error) {
      console.error('Error requesting more info:', error)
      message.error('Failed to request more information')
    } finally {
      setRequestInfoLoading(false)
    }
  }

  const handleApproveClaimItem = async (itemId: string) => {
    if (!id) return

    try {
      await claimItemsApi.approve(id, itemId)
      message.success('Claim item approved successfully')
      refetchClaimItems() // Refresh claim items
      refetchClaim() // Refresh claim to update status if needed
    } catch (error) {
      console.error('Error approving claim item:', error)
      message.error('Failed to approve claim item')
    }
  }

  const handleRejectClaimItem = async (itemId: string) => {
    if (!id) return

    try {
      await claimItemsApi.reject(id, itemId)
      message.success('Claim item rejected successfully')
      refetchClaimItems() // Refresh claim items
      refetchClaim() // Refresh claim to update status if needed
    } catch (error) {
      console.error('Error rejecting claim item:', error)
      message.error('Failed to reject claim item')
    }
  }

  const handleCompleteClaim = async () => {
    if (!id) return

    try {
      setCompleteLoading(true)
      await claimsApi.complete(id)
      message.success('Claim completed successfully')
      refetchClaim() // Refresh claim to update status
    } catch (error) {
      console.error('Error completing claim:', error)
      message.error('Failed to complete claim')
    } finally {
      setCompleteLoading(false)
    }
  }

  // Policy coverage modal handlers
  const handleViewPolicyCoverage = (categoryId: string, categoryName: string) => {
    setSelectedCategoryId(Number(categoryId))
    setSelectedCategoryName(categoryName)
    setPolicyCoverageModalVisible(true)
  }

  const handleClosePolicyCoverageModal = () => {
    setPolicyCoverageModalVisible(false)
    setSelectedCategoryId(null)
    setSelectedCategoryName('')
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
          <Button icon={<ArrowLeftOutlined />} onClick={handleBack} type="default">
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

          {/* Warranty Policy Information - Only for EVM Staff during review */}
          {canViewWarrantyPolicy && (
            <WarrantyPolicyCard warrantyPolicy={warrantyPolicy} loading={warrantyPolicyLoading} />
          )}

          <Divider />

          {/* Claim Items */}
          <ClaimItemsTable
            claimItems={claimItems}
            parts={parts}
            partCategories={partCategories}
            loading={itemsLoading}
            canAddItems={canAddItems}
            canApproveClaimItems={canApproveClaimItems}
            canRejectClaimItems={canRejectClaimItems}
            canViewPolicyCoverage={canViewPolicyCoverage}
            onAddItem={handleOpenAddItemModal}
            onApproveItem={handleApproveClaimItem}
            onRejectItem={handleRejectClaimItem}
            onViewCoverage={handleViewPolicyCoverage}
          />

          {/* Attachments */}
          <ClaimAttachments
            attachments={attachments}
            loading={attachmentsLoading}
            canAddAttachments={canAddAttachments}
            onAddAttachment={handleOpenAddAttachmentModal}
          />

          {/* Action Buttons */}
          {(canSubmitClaim ||
            canDeleteClaim ||
            canCancelClaim ||
            canStartReview ||
            canRequestMoreInfo ||
            canCompleteClaim) && (
            <div style={{ textAlign: 'center', marginTop: '24px' }}>
              <Space size="middle">
                {canDeleteClaim && (
                  <Button
                    type="default"
                    danger
                    size="large"
                    icon={<DeleteOutlined />}
                    loading={deleteLoading}
                    onClick={handleDeleteClaim}
                  >
                    Delete Claim
                  </Button>
                )}
                {canCancelClaim && (
                  <Button
                    type="default"
                    danger
                    size="large"
                    icon={<StopOutlined />}
                    loading={cancelLoading}
                    onClick={handleCancelClaim}
                  >
                    Cancel Claim
                  </Button>
                )}
                {canRequestMoreInfo && (
                  <Button
                    type="default"
                    size="large"
                    icon={<QuestionCircleOutlined />}
                    loading={requestInfoLoading}
                    onClick={handleRequestMoreInfo}
                  >
                    Request More Info
                  </Button>
                )}
                {canStartReview && (
                  <Button
                    type="primary"
                    size="large"
                    icon={<PlayCircleOutlined />}
                    loading={startReviewLoading}
                    onClick={handleStartReview}
                  >
                    Start Review
                  </Button>
                )}
                {canCompleteClaim && (
                  <Button
                    type="primary"
                    size="large"
                    icon={<CheckCircleOutlined />}
                    loading={completeLoading}
                    onClick={handleCompleteClaim}
                  >
                    Complete Claim
                  </Button>
                )}
                {canSubmitClaim && (
                  <Button
                    type="primary"
                    size="large"
                    icon={<SendOutlined />}
                    loading={submitLoading}
                    onClick={handleSubmitClaim}
                  >
                    Submit Claim
                  </Button>
                )}
              </Space>
            </div>
          )}
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

        {/* Add Attachment Modal */}
        {addAttachmentModalVisible && id && (
          <AddAttachmentModal
            visible={addAttachmentModalVisible}
            onCancel={handleCloseAddAttachmentModal}
            onSuccess={handleAddAttachmentSuccess}
            claimId={id}
          />
        )}

        {/* Policy Coverage Modal */}
        {policyCoverageModalVisible && selectedCategoryId && warrantyPolicy && (
          <PolicyCoverageModal
            visible={policyCoverageModalVisible}
            onCancel={handleClosePolicyCoverageModal}
            policyId={warrantyPolicy.id}
            categoryId={selectedCategoryId.toString()}
            categoryName={selectedCategoryName}
            partCategories={partCategories}
          />
        )}
      </div>
    </AppLayout>
  )
}

export default ClaimDetail
