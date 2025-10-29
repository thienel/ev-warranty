import { useCallback, useEffect, useState } from 'react'

const useDelay = (ms = 1000) => {
  const [timer, setTimer] = useState<NodeJS.Timeout | null>(null)

  const run = useCallback(
    (callback: () => void) => {
      if (timer) clearTimeout(timer)
      const newTimer = setTimeout(() => {
        callback()
      }, ms)
      setTimer(newTimer)
    },
    [ms, timer],
  )

  useEffect(() => {
    return () => {
      if (timer) clearTimeout(timer)
    }
  }, [timer])

  return run
}

export default useDelay
