import React from 'react'
import GenericActionBar from '@components/common/GenericActionBar/GenericActionBar.jsx'
import '@components/common/GenericActionBar/GenericActionBar.less'

const UserActionBar = ({ searchText, setSearchText, onReset, onOpenModal, loading }) => {
  return (
    <GenericActionBar
      searchText={searchText}
      setSearchText={setSearchText}
      onReset={onReset}
      onOpenModal={onOpenModal}
      loading={loading}
      searchPlaceholder="Search by name, email or role..."
      addButtonText="Add User"
    />
  )
}

export default UserActionBar
