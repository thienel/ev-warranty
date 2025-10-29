/* eslint-disable @typescript-eslint/no-explicit-any */
import localStorage from 'redux-persist/lib/storage'
import sessionStorage from 'redux-persist/lib/storage/session'

interface StorageState {
  remember?: boolean | string
}

const storageSwitcher = {
  getItem: (key: string): Promise<string | null> => {
    return localStorage.getItem(key).then((res: any) => {
      if (res) return res
      return sessionStorage.getItem(key)
    })
  },
  setItem: (key: string, value: string): Promise<void> => {
    const state: StorageState = JSON.parse(value)
    const remember = state?.remember === true || state?.remember === 'true'

    if (remember) {
      return localStorage.setItem(key, value)
    } else {
      return sessionStorage.setItem(key, value)
    }
  },
  removeItem: (key: string): Promise<void[]> => {
    return Promise.all([localStorage.removeItem(key), sessionStorage.removeItem(key)])
  },
}

export default storageSwitcher
