import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import { PURGE } from 'redux-persist'
import type { AuthState, LoginPayload } from '../types'

const initialState: AuthState = {
  user: null,
  token: null,
  isAuthenticated: false,
  remember: false,
  isLoading: false,
}

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    loginSuccess: (state, action: PayloadAction<LoginPayload>) => {
      state.user = action.payload.user
      state.token = action.payload.token
      state.isAuthenticated = true
      state.remember = action.payload.remember
      state.isLoading = false
    },
    logout: (state) => {
      state.user = null
      state.token = null
      state.isAuthenticated = false
      state.isLoading = false
      state.remember = false
    },
    setToken: (state, action: PayloadAction<string>) => {
      state.token = action.payload
      if (action.payload) {
        state.isAuthenticated = true
      }
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload
    },
    clearAuth: (state) => {
      state.user = null
      state.token = null
      state.isAuthenticated = false
      state.isLoading = false
      state.remember = false
    },
  },
  extraReducers: (builder) => {
    builder.addCase(PURGE, () => initialState)
  },
})

export const { loginSuccess, logout, setToken, setLoading, clearAuth } = authSlice.actions
export default authSlice.reducer
