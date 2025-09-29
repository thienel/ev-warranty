import React from 'react'
import { Navigate, Outlet, useRoutes } from 'react-router-dom'
import Home from '@pages/Home.jsx'
import Login from '@pages/auth/Login/Login.jsx'
import { useSelector} from 'react-redux'
import AuthCallBack from '@pages/auth/AuthCallBack.jsx'

export const ProtectedRoute = () => {
  const { isAuthenticated} = useSelector((state) => state.auth)

  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />
}

export const PublicRoute = () => {
  const { isAuthenticated} = useSelector((state) => state.auth)


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
      children: [
        { path: '/login', element: <Login /> },
        { path: '/auth/callback', element: <AuthCallBack /> },
      ],
    },
  ]

  return useRoutes(routes)
}

export default App
