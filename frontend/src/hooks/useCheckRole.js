import { useSelector } from 'react-redux'

const useCheckRole = (role) => {
  const { user } = useSelector((state) => state.auth)

  return user && role === user.role
}

export default useCheckRole
