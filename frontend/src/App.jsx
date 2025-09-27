import React from 'react'
import { Navigate, Outlet, useRoutes } from 'react-router-dom'
import Home from '@pages/Home.jsx'
import Login from '@components/Login'
import { useSelector } from 'react-redux'

export const ProtectedRoute = () => {
  const isAuthenticated = useSelector((state) => state.auth.isAuthenticated)
  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />
}

export const PublicRoute = () => {
  const isAuthenticated = useSelector((state) => state.auth.isAuthenticated)
  return !isAuthenticated ? <Outlet /> : <Navigate to="/" replace />
}

const App = () => {
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
