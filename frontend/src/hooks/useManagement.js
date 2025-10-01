import { useEffect, useState } from 'react'
import { message } from 'antd'
import api from '@services/api'
import { useDelay } from '@/hooks/index.js'

const useManagement = (apiEndpoint, itemName = 'item') => {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(false)
  const [searchText, setSearchText] = useState('')
  const [updateItem, setUpdateItem] = useState(null)
  const [isUpdate, setIsUpdate] = useState(false)
  const [isOpenModal, setIsOpenModal] = useState(false)

  const handleOpenModal = (item = null, isUpdate = false) => {
    setUpdateItem(item)
    setIsUpdate(isUpdate)
    setIsOpenModal(true)
  }

  const fetchItems = async () => {
    try {
      const response = await api.get(apiEndpoint)
      if (response.data.success) {
        const itemData = response.data.data || []
        setItems(itemData)
      }
    } catch (error) {
      message.error(error.response?.data?.message || `Failed to load ${itemName}s`)
      console.error(`Error fetching ${itemName}s:`, error)
    }
  }

  const delay = useDelay(300)

  const handleReset = async () => {
    setLoading(true)
    delay(async () => {
      setSearchText('')
      setIsOpenModal(false)
      setUpdateItem(null)
      await fetchItems()
      setLoading(false)
    })
  }

  useEffect(() => {
    fetchItems()
  }, [])

  return {
    items,
    setItems,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem,
    setUpdateItem,
    isUpdate,
    setIsUpdate,
    isOpenModal,
    setIsOpenModal,
    handleOpenModal,
    fetchItems,
    handleReset,
  }
}

export default useManagement
