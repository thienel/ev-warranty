import React from 'react'
import { Button, Col, Input, Row } from 'antd'
import { PlusOutlined, ReloadOutlined, SearchOutlined } from '@ant-design/icons'
import './GenericActionBar.less'

const GenericActionBar = ({
  searchText,
  setSearchText,
  onReset,
  onOpenModal,
  loading,
  searchPlaceholder = 'Search...',
  addButtonText = 'Add Item',
  className = 'generic-action-bar',
}) => {
  return (
    <Row gutter={[16, 16]} className={className}>
      <Col lg={18}>
        <Input
          placeholder={searchPlaceholder}
          prefix={<SearchOutlined />}
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          allowClear
        />
      </Col>
      <Col lg={6} className={'action-group'}>
        <Row>
          <Col lg={12}>
            <Button icon={<ReloadOutlined />} onClick={onReset} loading={loading}>
              Refresh
            </Button>
          </Col>
          <Col lg={12}>
            <Button type="primary" icon={<PlusOutlined />} onClick={onOpenModal}>
              {addButtonText}
            </Button>
          </Col>
        </Row>
      </Col>
    </Row>
  )
}

export default GenericActionBar
