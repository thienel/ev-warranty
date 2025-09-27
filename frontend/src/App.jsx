import React from 'react'
import { useRoutes } from 'react-router-dom'
import Home from '@pages/Home.jsx'
import Login from '@components/Login'

const App = () => {
  const routes = [
    {
      path: '/',
      element: <Home />,
    },
    {
      path: '/login',
      element: <Login />,
    },
  ]

  return useRoutes(routes)
}

export default App
