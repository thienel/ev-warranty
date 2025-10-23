import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import { PURGE } from 'redux-persist'
import type { AuthState, LoginPayload } from '../types'

const initialState: AuthState = {
  user: null,
  token: null,
  isAuthenticated: false,
  remember: false,
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
    },
  },
  extraReducers: (builder) => {
    builder.addCase(PURGE, () => initialState)
  },
})

export const { loginSuccess, logout, setToken } = authSlice.actions
export default authSlice.reducer