import React from 'react'
import { Button, Result } from 'antd'
import { useNavigate } from 'react-router-dom'
import './Error.less'
import { ERROR_MESSAGES } from '@constants/common-constants.js'

const Error = ({ code = 404 }) => {
  const navigate = useNavigate()

  const handleGoHome = () => {
    navigate('/')
  }

  return (
    <div className="error-container">
      <Result
        status={code}
        title={code}
        subTitle={ERROR_MESSAGES[code]}
        extra={
          <div className="error-actions">
            <Button type="primary" onClick={handleGoHome}>
              Back Home
            </Button>
          </div>
        }
      />
    </div>
  )
}

export default Error
