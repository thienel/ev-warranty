import localStorage from 'redux-persist/lib/storage'
import sessionStorage from 'redux-persist/lib/storage/session'

const storageSwitcher = {
  getItem: (key) => {
    return localStorage.getItem(key).then((res) => {
      if (res) return res
      return sessionStorage.getItem(key)
    })
  },
  setItem: (key, value) => {
    const state = JSON.parse(value)
    const remember = state?.remember ?? false

    console.log("state", state)
    console.log("remember", remember)

    if (remember) {
      return localStorage.setItem(key, value)
    } else {
      return sessionStorage.setItem(key, value)
    }
  },
  removeItem: (key) => {
    return Promise.all([
      localStorage.removeItem(key),
      sessionStorage.removeItem(key),
    ])
  },
}

export default storageSwitcher
