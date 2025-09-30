import React from 'react'
import { Button, Col, Input, Row, Space } from 'antd'
import { PlusOutlined, ReloadOutlined, SearchOutlined } from '@ant-design/icons'
import './UserActionBar.less'

const UserActionBar = ({ searchText, setSearchText, onReset, onOpenModal, loading }) => {
  return (
    <Row gutter={[16, 16]} className={'user-action-bar'}>
      <Col lg={18}>
        <Input
          placeholder="Search by name, email or role..."
          prefix={<SearchOutlined />}
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          allowClear
          size="large"
        />
      </Col>
      <Col lg={6} style={{ textAlign: 'right' }}>
        <Space wrap>
          <Button icon={<ReloadOutlined />} onClick={onReset} loading={loading}>
            Refresh
          </Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={onOpenModal}>
            Add User
          </Button>
        </Space>
      </Col>
    </Row>
  )
}

export default UserActionBar
