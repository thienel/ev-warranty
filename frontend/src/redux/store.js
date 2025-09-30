import { configureStore } from '@reduxjs/toolkit'
import { persistStore, persistReducer } from 'redux-persist'
import authReducer from './authSlice.js'
import storageSwitcher from '@redux/storageSwitcher.js'

const persistConfig = {
  key: 'auth',
  storage: storageSwitcher,
  whitelist: ['user', 'token', 'isAuthenticated', 'remember'],
}

const persistedAuthReducer = persistReducer(persistConfig, authReducer)

const store = configureStore({
  reducer: {
    auth: persistedAuthReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: [
          'persist/FLUSH',
          'persist/REHYDRATE',
          'persist/PAUSE',
          'persist/PERSIST',
          'persist/PURGE',
          'persist/REGISTER',
        ],
      },
    }),
})

export const persistor = persistStore(store)
export default store
