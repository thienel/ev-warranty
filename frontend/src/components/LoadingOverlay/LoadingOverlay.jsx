import React from 'react'
import { Spin } from 'antd'
import './LoadingOverlay.less'
import { LoadingOutlined } from '@ant-design/icons'

const LoadingOverlay = ({ children, loading = false }) => {
  return (
    <div className="loading-overlay-container">
      {children}
      {loading && (
        <div className="loading-overlay">
          <div className="loading-content">
            <Spin indicator={<LoadingOutlined style={{ fontSize: 50 }} spin />} />
          </div>
        </div>
      )}
    </div>
  )
}

export default LoadingOverlay
