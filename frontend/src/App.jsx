import React from 'react'
import { Navigate, Outlet, useRoutes } from 'react-router-dom'
import Login from '@pages/auth/Login/Login.jsx'
import { useSelector } from 'react-redux'
import AuthCallBack from '@pages/auth/AuthCallBack.jsx'
import Users from '@pages/Users.jsx'
import AppLayout from '@components/Layout/Layout.jsx'
import Offices from '@pages/Offices.jsx'

export const ProtectedRoute = () => {
  const { isAuthenticated } = useSelector((state) => state.auth)

  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />
}

export const PublicRoute = () => {
  const { isAuthenticated } = useSelector((state) => state.auth)

  return !isAuthenticated ? <Outlet /> : <Navigate to="/" replace />
}

const App = () => {
  const routes = [
    {
      element: <ProtectedRoute />,
      children: [
        { path: '/', element: <AppLayout /> },
        {
          path: '/users',
          element: <Users />,
        },
        {
          path: '/offices',
          element: <Offices />,
        },
      ],
    },
    {
      element: <PublicRoute />,
      children: [
        { path: '/login', element: <Login /> },
        { path: '/auth/callback', element: <AuthCallBack /> },
      ],
    },
  ]

  return useRoutes(routes)
}

export default App
