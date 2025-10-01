import React from 'react'
import { API_ENDPOINTS } from '@constants'
import OfficeModal from '@components/OfficeManagement/OfficeModal/OfficeModal.jsx'
import OfficeTable from '@components/OfficeManagement/OfficeTable/OfficeTable.jsx'
import OfficeActionBar from '@components/OfficeManagement/OfficeActionBar/OfficeActionBar.jsx'
import useManagement from '@/hooks/useManagement.js'

const OfficeManagement = () => {
  const {
    items: offices,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updateOffice,
    isUpdate,
    isOpenModal,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.OFFICE, 'office')

  return (
    <>
      <OfficeActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleOpenModal}
        loading={loading}
      />

      <OfficeTable
        loading={loading}
        setLoading={setLoading}
        offices={offices}
        searchText={searchText}
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
      />

      <OfficeModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        office={updateOffice}
        opened={isOpenModal}
        isUpdate={isUpdate}
      />
    </>
  )
}

export default OfficeManagement
