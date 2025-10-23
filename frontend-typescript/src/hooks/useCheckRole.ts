import { useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import type { RootState } from '@redux/store'
import type { UserRole } from '@constants/common-constants'

const useCheckRole = (role: UserRole): boolean => {
  const { user } = useSelector((state: RootState) => state.auth)

  const navigate = useNavigate()
  if (!user) {
    navigate('/login')
    return false
  }

  return role === user.role
}

export default useCheckRole