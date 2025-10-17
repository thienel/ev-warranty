import React from 'react'
import { Navigate, Outlet, useRoutes } from 'react-router-dom'
import Login from '@pages/Login/Login.jsx'
import { useSelector } from 'react-redux'
import CallBack from '@pages/CallBack.jsx'
import Users from '@pages/admin/Users.jsx'
import AppLayout from '@components/Layout/Layout.jsx'
import Offices from '@pages/admin/Offices.jsx'
import Error from '@pages/Error/Error.jsx'
import useCheckRole from '@/hooks/useCheckRole.js'
import { USER_ROLES } from '@constants'

export const ProtectedRoute = () => {
  const { isAuthenticated } = useSelector((state) => state.auth)

  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />
}

export const PublicRoute = () => {
  const { isAuthenticated } = useSelector((state) => state.auth)

  return !isAuthenticated ? <Outlet /> : <Navigate to="/" replace />
}

export const AdminRoute = () => {
  const isRightRole = useCheckRole(USER_ROLES.ADMIN)

  return isRightRole ? <Outlet /> : <Navigate to="/unauthorized" replace />
}

const App = () => {
  const routes = [
    {
      element: <ProtectedRoute />,
      children: [
        { path: '/', element: <AppLayout /> },
        {
          path: '/admin',
          element: <AdminRoute />,
          children: [
            {
              path: 'users',
              element: <Users />,
            },
            {
              path: 'offices',
              element: <Offices />,
            },
          ],
        },
      ],
    },
    {
      element: <PublicRoute />,
      children: [
        { path: '/login', element: <Login /> },
        { path: '/callback', element: <CallBack /> },
      ],
    },
    {
      path: '/unauthorized',
      element: <Error code={403} />,
    },
    {
      path: '/servererror',
      element: <Error code={500} />,
    },
    {
      path: '*',
      element: <Error code={404} />,
    },
  ]

  return useRoutes(routes)
}

export default App
