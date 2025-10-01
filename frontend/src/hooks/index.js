import { useState, useEffect, useCallback } from 'react'

/**
 * useDelay - custom hook để delay thực thi một function
 * @param {number} ms - số mili giây delay
 * @returns {function} run - function để trigger delay
 */
export const useDelay = (ms = 1000) => {
  const [timer, setTimer] = useState(null)

  const run = useCallback(
    (callback) => {
      if (timer) clearTimeout(timer)
      const newTimer = setTimeout(() => {
        callback()
      }, ms)
      setTimer(newTimer)
    },
    [ms, timer]
  )

  useEffect(() => {
    return () => {
      if (timer) clearTimeout(timer)
    }
  }, [timer])

  return run
}

export { default as useManagement } from './useManagement.js'
