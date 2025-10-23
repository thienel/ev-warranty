import { useCallback, useEffect, useState } from 'react'
import api from '@services/api'
import useDelay from '@/hooks/useDelay'
import useHandleApiError from '@/hooks/useHandleApiError'
import type { ErrorResponse } from '@/constants/error-messages'

interface UseManagementReturn<T> {
  items: T[]
  setItems: (items: T[]) => void
  loading: boolean
  setLoading: (loading: boolean) => void
  searchText: string
  setSearchText: (text: string) => void
  updateItem: T | null
  setUpdateItem: (item: T | null) => void
  isUpdate: boolean
  setIsUpdate: (isUpdate: boolean) => void
  isOpenModal: boolean
  setIsOpenModal: (isOpen: boolean) => void
  handleOpenModal: (item?: T | null, isUpdate?: boolean) => void
  fetchItems: () => Promise<void>
  handleReset: () => Promise<void>
}

const useManagement = <T = Record<string, unknown>>(apiEndpoint: string): UseManagementReturn<T> => {
  const [items, setItems] = useState<T[]>([])
  const [loading, setLoading] = useState(false)
  const [searchText, setSearchText] = useState('')
  const [updateItem, setUpdateItem] = useState<T | null>(null)
  const [isUpdate, setIsUpdate] = useState(false)
  const [isOpenModal, setIsOpenModal] = useState(false)
  const handleError = useHandleApiError()

  const handleOpenModal = (item: T | null = null, isUpdate = false) => {
    setUpdateItem(item)
    setIsUpdate(isUpdate)
    setIsOpenModal(true)
  }

  const fetchItems = useCallback(async () => {
    try {
      const response = await api.get(apiEndpoint)
      const itemData = response.data.data || []
      setItems(itemData)
    } catch (error) {
      handleError(error as ErrorResponse)
    }
  }, [apiEndpoint, handleError])

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
  }, [fetchItems])

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