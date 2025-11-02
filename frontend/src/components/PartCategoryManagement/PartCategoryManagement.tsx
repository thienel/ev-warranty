import React, { useCallback, useEffect, useState } from 'react'
import { API_ENDPOINTS } from '@constants/common-constants'
import { type PartCategory } from '@/types/index'
import { partCategoriesApi } from '@services/index'
import PartCategoryModal from '@components/PartCategoryManagement/PartCategoryModal/PartCategoryModal'
import useManagement from '@/hooks/useManagement'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar'
import GenericTable from '@components/common/GenericTable/GenericTable'
import GenerateColumns from './partCategoryTableColumns'
import useHandleApiError from '@/hooks/useHandleApiError'

const PartCategoryManagement: React.FC = () => {
  const {
    items: partCategories,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updatePartCategory,
    isUpdate,
    isOpenModal,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.PART_CATEGORIES)

  const [allPartCategories, setAllPartCategories] = useState<PartCategory[]>([])
  const [partCategoriesLoading, setPartCategoriesLoading] = useState(false)
  const handleError = useHandleApiError()

  const fetchPartCategories = useCallback(async (): Promise<void> => {
    try {
      setPartCategoriesLoading(true)
      const response = await partCategoriesApi.getAll()
      // Handle different response structures - same as useManagement hook
      let categoriesData = response.data
      // If response has nested data property, use that
      if (categoriesData && typeof categoriesData === 'object' && 'data' in categoriesData) {
        categoriesData = (categoriesData as { data: unknown }).data as PartCategory[]
      }
      // Ensure we always have an array
      if (Array.isArray(categoriesData)) {
        setAllPartCategories(categoriesData)
      } else {
        setAllPartCategories([])
      }
    } catch (error) {
      handleError(error as Error)
      setAllPartCategories([]) // Set empty array on error
    } finally {
      setPartCategoriesLoading(false)
    }
  }, [handleError])

  useEffect(() => {
    fetchPartCategories()
  }, [fetchPartCategories])

  // Refetch part categories when modal opens if list is empty
  useEffect(() => {
    if (isOpenModal && allPartCategories.length === 0 && !partCategoriesLoading) {
      fetchPartCategories()
    }
  }, [isOpenModal, allPartCategories.length, partCategoriesLoading, fetchPartCategories])

  const searchFields = ['category_name', 'description', 'parent_category_name']

  return (
    <>
      <GenericActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleOpenModal}
        loading={loading}
        searchPlaceholder="Search by category name, description or parent..."
        addButtonText="Add Part Category"
        allowCreate={true}
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={partCategories as (Record<string, unknown> & { id: string | number })[]}
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFields}
        deleteEndpoint={API_ENDPOINTS.PART_CATEGORIES}
        deleteSuccessMessage="Part category deleted successfully"
      />

      <PartCategoryModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        partCategory={updatePartCategory ? (updatePartCategory as unknown as PartCategory) : null}
        opened={isOpenModal}
        partCategories={allPartCategories}
        partCategoriesLoading={partCategoriesLoading}
        isUpdate={isUpdate}
      />
    </>
  )
}

export default PartCategoryManagement
