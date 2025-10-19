import React, { useEffect } from 'react'
import LoadingOverlay from '@components/LoadingOverlay/LoadingOverlay.jsx'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { loginSuccess, setToken } from '@redux/authSlice.js'
import { message } from 'antd'
import { useDispatch } from 'react-redux'
import api from '@services/api.js'
import { API_ENDPOINTS } from '@constants/common-constants.js'

const CallBack = () => {
  const [searchParams] = useSearchParams()
  const dispatch = useDispatch()
  const navigate = useNavigate()

  useEffect(() => {
    const handleLogin = async () => {
      const token = searchParams.get('token')
      try {
        dispatch(setToken(token))
        const res = (await api.get(API_ENDPOINTS.AUTH.TOKEN, { withCredentials: true })).data
        if (!res.data.valid) {
          console.error(res.message)
          message.error('Login failed. Please check your credentials.')
        } else {
          navigate('/')
          dispatch(loginSuccess({ user: res.data.user, token }))
        }
      } catch (error) {
        console.error(error)
        message.error('Login failed. Please check your credentials.')
        navigate('/login')
      }
    }

    handleLogin()
  }, [dispatch, searchParams, navigate])

  return <LoadingOverlay loading={true} />
}

export default CallBack
