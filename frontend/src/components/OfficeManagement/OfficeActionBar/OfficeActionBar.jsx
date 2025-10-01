import React from 'react'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar.jsx'

const OfficeActionBar = ({ searchText, setSearchText, onReset, onOpenModal, loading }) => {
  return (
    <GenericActionBar
      searchText={searchText}
      setSearchText={setSearchText}
      onReset={onReset}
      onOpenModal={onOpenModal}
      loading={loading}
      searchPlaceholder="Search by office name, type or address..."
      addButtonText="Add Office"
    />
  )
}

export default OfficeActionBar
