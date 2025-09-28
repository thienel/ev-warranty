import { createSlice } from '@reduxjs/toolkit'
import { PURGE } from 'redux-persist'

const initialState = {
  user: null,
  token: null,
  isAuthenticated: false,
  isInitialized: false,
}

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setInitialized: (state) => {
      state.isInitialized = true
    },
    loginSuccess: (state, action) => {
      state.user = action.payload.user
      state.token = action.payload.token
      state.isAuthenticated = true
      state.isInitialized = true
    },
    loginFailure: (state) => {
      state.isInitialized = true
    },
    logout: (state) => {
      state.user = null
      state.token = null
      state.isAuthenticated = false
      state.isLoading = false
    },
    setToken: (state, action) => {
      state.token = action.payload
    },
  },
  extraReducers: (builder) => {
    builder.addCase(PURGE, () => initialState)
  },
})

export const { setInitialized, loginSuccess, loginFailure, logout, setToken } = authSlice.actions
export default authSlice.reducer
