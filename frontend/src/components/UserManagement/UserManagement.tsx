import React, { useCallback, useEffect, useState } from 'react'
import { API_ENDPOINTS, ROLE_LABELS } from '@constants/common-constants'
import { type User, type Office } from '@/types/index'
import { officesApi } from '@services/index'
import UserModal from '@components/UserManagement/UserModal/UserModal'
import useManagement from '@/hooks/useManagement'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar'
import GenericTable from '@components/common/GenericTable/GenericTable'
import GenerateColumns from './userTableColumns'
import useHandleApiError from '@/hooks/useHandleApiError'

const UserManagement: React.FC = () => {
  const {
    items: users,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updateUser,
    isUpdate,
    isOpenModal,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.USERS)

  const [offices, setOffices] = useState<Office[]>([])
  const [officesLoading, setOfficesLoading] = useState(false)
  const handleError = useHandleApiError()

  const fetchOffices = useCallback(async (): Promise<void> => {
    try {
      setOfficesLoading(true)
      const response = await officesApi.getAll()
      // Handle different response structures - same as useManagement hook
      let officesData = response.data
      // If response has nested data property, use that
      if (officesData && typeof officesData === 'object' && 'data' in officesData) {
        officesData = (officesData as { data: unknown }).data as Office[]
      }
      // Ensure we always have an array
      if (Array.isArray(officesData)) {
        setOffices(officesData)
      } else {
        setOffices([])
      }
    } catch (error) {
      handleError(error as Error)
      setOffices([]) // Set empty array on error
    } finally {
      setOfficesLoading(false)
    }
  }, [handleError])

  useEffect(() => {
    fetchOffices()
  }, [fetchOffices])

  // Refetch offices when modal opens if offices list is empty
  useEffect(() => {
    if (isOpenModal && offices.length === 0 && !officesLoading) {
      fetchOffices()
    }
  }, [isOpenModal, offices.length, officesLoading, fetchOffices])

  const getOfficeName = (officeId: string): string => {
    // Safety check: ensure offices is an array before calling find
    if (!Array.isArray(offices)) {
      return 'N/A'
    }

    const office = offices.find((o) => o.id === officeId)
    return office ? office.office_name : 'N/A'
  }

  const searchFields = ['name', 'email']
  const searchFieldsWithRole = [
    ...searchFields,
    (user: Record<string, unknown> & { id: string | number }) => {
      const userRecord = user as unknown as User
      return ROLE_LABELS[userRecord.role as keyof typeof ROLE_LABELS] || userRecord.role
    },
  ]

  return (
    <>
      <GenericActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleOpenModal}
        loading={loading}
        searchPlaceholder="Search by name, email or role..."
        addButtonText="Add User"
        allowCreate={true}
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={
          users.filter((v) => v.email !== 'admin@example.com') as (Record<string, unknown> & {
            id: string | number
          })[]
        }
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFieldsWithRole}
        deleteEndpoint={API_ENDPOINTS.USERS}
        deleteSuccessMessage="User deleted successfully"
        additionalProps={{ getOfficeName }}
      />

      <UserModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        user={updateUser ? (updateUser as unknown as User) : null}
        opened={isOpenModal}
        offices={offices}
        officesLoading={officesLoading}
        isUpdate={isUpdate}
      />
    </>
  )
}

export default UserManagement
