import { useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'

const useCheckRole = (role) => {
  const { user } = useSelector((state) => state.auth)

  const navigate = useNavigate()
  if (!user) {
    navigate('/login')
  }

  return role === user.role
}

export default useCheckRole
