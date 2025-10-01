import React, { useEffect, useState } from 'react'
import { message } from 'antd'
import api from '@services/api'
import { API_ENDPOINTS } from '@constants'
import OfficeModal from '@components/OfficeManagement/OfficeModal/OfficeModal.jsx'
import OfficeTable from '@components/OfficeManagement/OfficeTable/OfficeTable.jsx'
import OfficeActionBar from '@components/OfficeManagement/OfficeActionBar/OfficeActionBar.jsx'
import { useDelay } from '@/hooks/index.js'

const OfficeManagement = () => {
  const [offices, setOffices] = useState([])

  const [loading, setLoading] = useState(false)
  const [searchText, setSearchText] = useState('')

  const [updateOffice, setUpdateOffice] = useState(null)
  const [isUpdate, setIsUpdate] = useState(false)
  const [isOpenModal, setIsOpenModal] = useState(false)

  const handleOpenModal = (office = null, isUpdate = false) => {
    setUpdateOffice(office)
    setIsUpdate(isUpdate)
    setIsOpenModal(true)
  }

  const fetchOffices = async () => {
    try {
      const response = await api.get(API_ENDPOINTS.OFFICE)

      if (response.data.success) {
        const officeData = response.data.data || []
        setOffices(officeData)
      }
    } catch (error) {
      message.error(error.response?.data?.message || 'Failed to load offices')
      console.error('Error fetching offices:', error)
    }
  }

  useEffect(() => {
    fetchOffices()
  }, [])

  const delay = useDelay(300)

  const handleReset = async () => {
    setLoading(true)
    delay(async () => {
      setSearchText('')
      setIsOpenModal(false)
      setUpdateOffice(null)
      await fetchOffices()
      setLoading(false)
    })
  }

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
