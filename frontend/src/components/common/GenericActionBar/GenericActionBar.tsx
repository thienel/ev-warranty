import React from 'react'
import { Button, Col, Input, Row } from 'antd'
import { PlusOutlined, ReloadOutlined, SearchOutlined } from '@ant-design/icons'
import './GenericActionBar.less'

interface GenericActionBarProps {
  searchText: string
  setSearchText: (text: string) => void
  onReset: () => void
  onOpenModal: () => void
  loading: boolean
  searchPlaceholder?: string
  addButtonText?: string
  className?: string
  allowCreate: boolean
}

const GenericActionBar: React.FC<GenericActionBarProps> = ({
  searchText,
  setSearchText,
  onReset,
  onOpenModal,
  loading,
  searchPlaceholder = 'Search...',
  addButtonText = 'Add Item',
  className = 'generic-action-bar',
  allowCreate = true,
}) => {
  return (
    <Row gutter={[16, 16]} className={className}>
      <Col lg={allowCreate ? 18 : 20}>
        <Input
          placeholder={searchPlaceholder}
          prefix={<SearchOutlined />}
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          allowClear
        />
      </Col>
      <Col lg={allowCreate ? 6 : 4} className="action-group">
        <Row>
          <Col lg={allowCreate ? 12 : 24}>
            <Button icon={<ReloadOutlined />} onClick={onReset} loading={loading}>
              Refresh
            </Button>
          </Col>
          {allowCreate && (
            <Col lg={12}>
              <Button type="primary" icon={<PlusOutlined />} onClick={onOpenModal}>
                {addButtonText}
              </Button>
            </Col>
          )}
        </Row>
      </Col>
    </Row>
  )
}

export default GenericActionBar
