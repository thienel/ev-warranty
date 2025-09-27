import React, { useState } from 'react'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UploadOutlined,
  UserOutlined,
  VideoCameraOutlined,
  DashboardOutlined,
  SettingOutlined,
  LogoutOutlined,
} from '@ant-design/icons'
import { Button, Layout, Menu, theme, Typography, Avatar, Space } from 'antd'

const { Header, Sider, Content } = Layout
const { Text } = Typography

const AppLayout = () => {
  const [collapsed, setCollapsed] = useState(false)
  const {
    token: {
      colorBgContainer,
      borderRadiusLG,
      colorPrimary,
      colorText,
      colorTextSecondary,
      boxShadow,
    },
  } = theme.useToken()

  const menuItems = [
    {
      key: '1',
      icon: <DashboardOutlined />,
      label: 'Dashboard',
      style: { marginTop: '8px' },
    },
    { key: '2', icon: <UserOutlined />, label: 'Users' },
    { key: '3', icon: <VideoCameraOutlined />, label: 'Media' },
    { key: '4', icon: <UploadOutlined />, label: 'Uploads' },
    {
      key: '5',
      icon: <SettingOutlined />,
      label: 'Settings',
      style: { marginTop: 'auto' },
    },
  ]

  return (
    <Layout
      style={{
        height: '100vh',
        overflow: 'hidden',
      }}
    >
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        width={260}
        collapsedWidth={80}
        style={{
          boxShadow: '2px 0 8px rgba(0,0,0,0.15)',
          position: 'relative',
          zIndex: 10,
        }}
      >
        {/* Logo Section */}
        <div
          style={{
            height: '64px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: collapsed ? 'center' : 'flex-start',
            padding: collapsed ? '0' : '0 24px',
            borderBottom: '1px solid rgba(255,255,255,0.1)',
            transition: 'all 0.3s ease',
          }}
        >
          {collapsed ? (
            <div
              style={{
                width: '32px',
                height: '32px',
                borderRadius: '8px',
                background: 'rgba(255,255,255,0.1)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: 'white',
                fontSize: '16px',
                fontWeight: 'bold',
              }}
            >
              A
            </div>
          ) : (
            <Space>
              <div
                style={{
                  width: '32px',
                  height: '32px',
                  borderRadius: '8px',
                  background: 'rgba(255,255,255,0.1)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: 'white',
                  fontSize: '16px',
                  fontWeight: 'bold',
                }}
              >
                A
              </div>
              <Text style={{ color: 'white', fontSize: '18px', fontWeight: '600' }}>
                Admin Panel
              </Text>
            </Space>
          )}
        </div>

        {/* Menu Section */}
        <Menu
          theme="dark"
          mode="inline"
          defaultSelectedKeys={['1']}
          items={menuItems}
          style={{
            backgroundColor: 'transparent',
            border: 'none',
            padding: '16px 8px',
          }}
        />

        {/* User Profile Section */}
        {!collapsed && (
          <div
            style={{
              position: 'absolute',
              bottom: '16px',
              left: '16px',
              right: '16px',
              padding: '12px',
              borderRadius: '8px',
              background: 'rgba(255,255,255,0.1)',
              backdropFilter: 'blur(10px)',
            }}
          >
            <Space>
              <Avatar size="small" style={{ backgroundColor: '#87d068' }}>
                U
              </Avatar>
              <div>
                <Text
                  style={{ color: 'white', fontSize: '14px', display: 'block', lineHeight: 1.2 }}
                >
                  John Doe
                </Text>
                <Text style={{ color: 'rgba(255,255,255,0.7)', fontSize: '12px' }}>
                  Administrator
                </Text>
              </div>
            </Space>
          </div>
        )}
      </Sider>

      <Layout style={{ height: '100vh', overflow: 'hidden' }}>
        {/* Header */}
        <Header
          style={{
            padding: '0 24px',
            background: colorBgContainer,
            borderBottom: `1px solid #f0f0f0`,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            boxShadow: '0 1px 4px rgba(0,0,0,0.1)',
            zIndex: 9,
          }}
        >
          <Space>
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
              style={{
                fontSize: '16px',
                width: '40px',
                height: '40px',
                color: colorText,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            />

            <div style={{ marginLeft: '16px' }}>
              <Text style={{ fontSize: '18px', fontWeight: '600', color: colorText }}>
                Dashboard
              </Text>
              <Text
                style={{
                  fontSize: '14px',
                  color: colorTextSecondary,
                  display: 'block',
                  marginTop: '-2px',
                }}
              >
                Welcome back, manage your application
              </Text>
            </div>
          </Space>

          {/* Header Actions */}
          <Space>
            <Button
              type="text"
              icon={<SettingOutlined />}
              style={{
                color: colorTextSecondary,
                border: 'none',
              }}
            />
            <Button
              type="text"
              icon={<LogoutOutlined />}
              style={{
                color: colorTextSecondary,
                border: 'none',
              }}
            />
            <Avatar style={{ backgroundColor: colorPrimary, marginLeft: '8px' }}>JD</Avatar>
          </Space>
        </Header>

        {/* Content */}
        <Content
          style={{
            margin: '24px',
            padding: '24px',
            minHeight: 'calc(100vh - 112px)',
            background: colorBgContainer,
            borderRadius: borderRadiusLG,
            boxShadow: boxShadow,
            overflow: 'auto',
            position: 'relative',
          }}
        >
          {/* Content Header */}
          <div style={{ marginBottom: '24px' }}>
            <Text style={{ fontSize: '24px', fontWeight: '600', color: colorText }}>
              Main Content Area
            </Text>
            <Text
              style={{
                fontSize: '14px',
                color: colorTextSecondary,
                display: 'block',
                marginTop: '4px',
              }}
            >
              This is your main content area where you can display your application content
            </Text>
          </div>

          {/* Demo Content */}
          <div
            style={{
              height: '400px',
              background: `linear-gradient(135deg, ${colorPrimary}20, ${colorPrimary}10)`,
              borderRadius: borderRadiusLG,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              border: `1px dashed ${colorPrimary}40`,
            }}
          >
            <Text style={{ color: colorTextSecondary, fontSize: '16px' }}>
              Your content goes here...
            </Text>
          </div>
        </Content>
      </Layout>
    </Layout>
  )
}

export default AppLayout
