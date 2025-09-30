import React from 'react'
import { Button, Col, Input, Row } from 'antd'
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
              Add User
            </Button>
          </Col>
        </Row>
      </Col>
    </Row>
  )
}

export default UserActionBar
