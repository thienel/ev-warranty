import React from 'react'
import { Button, Result } from 'antd'
import { useNavigate } from 'react-router-dom'
import './NotFound.less'

const NotFound = () => {
  const navigate = useNavigate()

  const handleGoHome = () => {
    navigate('/')
  }

  const handleGoBack = () => {
    navigate(-1)
  }

  return (
    <div className="not-found-container">
      <Result
        status="404"
        title="404"
        subTitle="Sorry, the page you visited does not exist."
        extra={
          <div className="not-found-actions">
            <Button type="primary" onClick={handleGoHome}>
              Back Home
            </Button>
            <Button onClick={handleGoBack}>Go Back</Button>
          </div>
        }
      />
    </div>
  )
}

export default NotFound
