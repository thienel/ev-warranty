import React from 'react'
import { Card, Tag } from 'antd'
import { CheckCircleOutlined, ExclamationCircleOutlined } from '@ant-design/icons'

interface StepCardProps {
  title: string
  icon: React.ReactNode
  children: React.ReactNode
  isActive: boolean
  isCompleted: boolean
  isDisabled?: boolean
  hoverable?: boolean
}

const StepCard: React.FC<StepCardProps> = ({
  title,
  icon,
  children,
  isActive,
  isCompleted,
  isDisabled = false,
  hoverable = true,
}) => {
  const getStatusTag = () => {
    if (isCompleted) {
      return (
        <Tag color="success" icon={<CheckCircleOutlined />}>
          Selected
        </Tag>
      )
    }

    if (isDisabled) {
      return <Tag color="default">Waiting</Tag>
    }

    return (
      <Tag color="warning" icon={<ExclamationCircleOutlined />}>
        Required
      </Tag>
    )
  }

  const getCardClassName = () => {
    const baseClass = 'step-card'
    const classes = [baseClass]

    if (isActive) classes.push('active')
    if (isCompleted) classes.push('completed')
    if (isDisabled) classes.push('disabled')

    return classes.join(' ')
  }

  return (
    <Card
      title={
        <div className="card-header">
          <div className="card-title">
            {icon}
            <span>{title}</span>
          </div>
          <div className="card-status">{getStatusTag()}</div>
        </div>
      }
      className={getCardClassName()}
      hoverable={hoverable && !isDisabled}
    >
      <div className="card-content">{children}</div>
    </Card>
  )
}

export default StepCard
