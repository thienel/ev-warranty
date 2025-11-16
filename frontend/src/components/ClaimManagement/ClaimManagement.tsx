import React from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { API_ENDPOINTS, USER_ROLES } from '@constants/common-constants'
import { type Claim } from '@/types/index'
import useClaimsManagement from '@/components/ClaimManagement/useClaimsManagement'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar'
import GenericTable from '@components/common/GenericTable/GenericTable'
import GenerateColumns from './claimTableColumns'
import { allowRoles, getClaimsBasePath } from '@/utils/navigationHelpers'

const ClaimManagement: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()

  const {
    claims,
    loading,
    setLoading,
    searchText,
    setSearchText,
    handleOpenModal,
    handleReset,
    allowCreate,
  } = useClaimsManagement()

  const handleViewDetails = (claim: Claim): void => {
    if (
      allowRoles(location.pathname, [
        USER_ROLES.ADMIN,
        USER_ROLES.EVM_STAFF,
        USER_ROLES.SC_STAFF,
        USER_ROLES.SC_TECHNICIAN,
      ])
    ) {
      const basePath = getClaimsBasePath(location.pathname)
      navigate(`${basePath}/${claim.id}`)
    }
  }

  const handleCreateClaim = (): void => {
    if (allowRoles(location.pathname, [USER_ROLES.SC_STAFF])) {
      const basePath = getClaimsBasePath(location.pathname)
      navigate(`${basePath}/create`)
    }
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
        allowCreate={allowCreate}
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={
          (Array.isArray(claims) ? claims : []).sort((a, b) => {
            const dateA = new Date(a.updated_at as string)
            const dateB = new Date(b.updated_at as string)
            return dateB.getTime() - dateA.getTime()
          }) as unknown as (Record<string, unknown> & {
            id: string | number
          })[]
        }
        onOpenModal={(record) => record && handleOpenModal(record as unknown as Claim, false)}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFields}
        deleteEndpoint={API_ENDPOINTS.CLAIMS}
        deleteSuccessMessage="Claim deleted successfully"
        additionalProps={{ onViewDetails: handleViewDetails }}
      />
    </>
  )
}

export default ClaimManagement
