import React, { useCallback, useEffect, useState } from 'react'
import { API_ENDPOINTS } from '@constants/common-constants'
import { type Part, type PartCategory, type Office } from '@/types/index'
import { partCategoriesApi, officesApi } from '@services/index'
import PartModal from '@components/PartManagement/PartModal/PartModal'
import useManagement from '@/hooks/useManagement'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar'
import GenericTable from '@components/common/GenericTable/GenericTable'
import GenerateColumns from './partTableColumns'
import useHandleApiError from '@/hooks/useHandleApiError'

const PartManagement: React.FC = () => {
  const {
    items: parts,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updatePart,
    isUpdate,
    isOpenModal,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.PARTS)

  const [partCategories, setPartCategories] = useState<PartCategory[]>([])
  const [partCategoriesLoading, setPartCategoriesLoading] = useState(false)
  const [offices, setOffices] = useState<Office[]>([])
  const [officesLoading, setOfficesLoading] = useState(false)
  const handleError = useHandleApiError()

  const fetchPartCategories = useCallback(async (): Promise<void> => {
    try {
      setPartCategoriesLoading(true)
      const response = await partCategoriesApi.getAll()
      let categoriesData = response.data
      if (categoriesData && typeof categoriesData === 'object' && 'data' in categoriesData) {
        categoriesData = (categoriesData as { data: unknown }).data as PartCategory[]
      }
      if (Array.isArray(categoriesData)) {
        setPartCategories(categoriesData)
      } else {
        setPartCategories([])
      }
    } catch (error) {
      handleError(error as Error)
      setPartCategories([])
    } finally {
      setPartCategoriesLoading(false)
    }
  }, [handleError])

  const fetchOffices = useCallback(async (): Promise<void> => {
    try {
      setOfficesLoading(true)
      const response = await officesApi.getAll()
      let officesData = response.data
      if (officesData && typeof officesData === 'object' && 'data' in officesData) {
        officesData = (officesData as { data: unknown }).data as Office[]
      }
      if (Array.isArray(officesData)) {
        setOffices(officesData)
      } else {
        setOffices([])
      }
    } catch (error) {
      handleError(error as Error)
      setOffices([])
    } finally {
      setOfficesLoading(false)
    }
  }, [handleError])

  useEffect(() => {
    fetchPartCategories()
    fetchOffices()
  }, [fetchPartCategories, fetchOffices])

  // Refetch when modal opens if lists are empty
  useEffect(() => {
    if (isOpenModal) {
      if (partCategories.length === 0 && !partCategoriesLoading) {
        fetchPartCategories()
      }
      if (offices.length === 0 && !officesLoading) {
        fetchOffices()
      }
    }
  }, [
    isOpenModal,
    partCategories.length,
    partCategoriesLoading,
    offices.length,
    officesLoading,
    fetchPartCategories,
    fetchOffices,
  ])

  const getOfficeName = (officeId: string): string => {
    if (!Array.isArray(offices)) {
      return 'N/A'
    }

    const office = offices.find((o) => o.id === officeId)
    return office ? office.office_name : 'N/A'
  }

  const searchFields = ['serial_number', 'part_name', 'category_name', 'status']

  return (
    <>
      <GenericActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleOpenModal}
        loading={loading}
        searchPlaceholder="Search by serial number, part name, category or status..."
        addButtonText="Add Part"
        allowCreate={true}
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={parts as (Record<string, unknown> & { id: string | number })[]}
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFields}
        deleteEndpoint={API_ENDPOINTS.PARTS}
        deleteSuccessMessage="Part deleted successfully"
        additionalProps={{ getOfficeName }}
      />

      <PartModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        part={updatePart ? (updatePart as unknown as Part) : null}
        opened={isOpenModal}
        partCategories={partCategories}
        offices={offices}
        partCategoriesLoading={partCategoriesLoading}
        officesLoading={officesLoading}
        isUpdate={isUpdate}
      />
    </>
  )
}

export default PartManagement
