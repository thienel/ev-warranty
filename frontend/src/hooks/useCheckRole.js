import { useEffect } from 'react'
import { useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'

const useCheckRole = (rolesAllowed = []) => {
  const { user } = useSelector((state) => state.auth)
  const navigate = useNavigate()

  useEffect(() => {
    if (!user || !rolesAllowed.includes(user?.role)) {
      navigate('/unauthorized')
    }
  }, [rolesAllowed, navigate, user])
}

export default useCheckRole
