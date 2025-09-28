import React, { useEffect } from 'react'
import { Navigate, Outlet, useRoutes } from 'react-router-dom'
import Home from '@pages/Home.jsx'
import Login from '@pages/Login'
import { useSelector, useDispatch } from 'react-redux'
import { setInitialized } from '@redux/authSlice.js'
import LoadingOverlay from '@components/LoadingOverlay/LoadingOverlay.jsx'

export const ProtectedRoute = () => {
  const { isAuthenticated, isLoading, isInitialized } = useSelector((state) => state.auth)

  if (!isInitialized || isLoading) {
    return <LoadingOverlay loading={true} children={null} />
  }

  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />
}

export const PublicRoute = () => {
  const { isAuthenticated, isLoading, isInitialized } = useSelector((state) => state.auth)

  if (!isInitialized || isLoading) {
    return <LoadingOverlay loading={true} children={null} />
  }

  return !isAuthenticated ? <Outlet /> : <Navigate to="/" replace />
}

const App = () => {
  const dispatch = useDispatch()
  const { isInitialized } = useSelector((state) => state.auth)

  useEffect(() => {
    if (!isInitialized) {
      const timer = setTimeout(() => {
        dispatch(setInitialized())
      }, 100)
      return () => clearTimeout(timer)
    }
  }, [dispatch, isInitialized])

  const routes = [
    {
      element: <ProtectedRoute />,
      children: [{ path: '/', element: <Home /> }],
    },
    {
      element: <PublicRoute />,
      children: [{ path: '/login', element: <Login /> }],
    },
  ]

  return useRoutes(routes)
}

export default App
