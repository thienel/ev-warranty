import { CLAIM_STATUS_LABELS } from '@constants/common-constants.js'
import { Button, Space } from 'antd'
import { EyeOutlined, UserOutlined, CarOutlined, CalendarOutlined } from '@ant-design/icons'

const GenerateColumns = (sortedInfo, filteredInfo, onOpenModal, onDelete, additionalProps) => {
  const { onViewDetails } = additionalProps

  return [
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Customer Name</span>,
      dataIndex: 'customer_name',
      key: 'customer_name',
      width: '25%',
      sorter: (a, b) => (a.customer_name || '').localeCompare(b.customer_name || ''),
      sortOrder: sortedInfo.columnKey === 'customer_name' ? sortedInfo.order : null,
      render: (text, record) => (
        <Space style={{ padding: '0 14px', whiteSpace: 'normal', wordBreak: 'break-word' }}>
          <UserOutlined style={{ color: '#697565' }} />
          <span>{text || `Customer ${record.customer_id?.slice(0, 8)}` || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Vehicle</span>,
      dataIndex: 'vehicle_info',
      key: 'vehicle_info',
      width: '25%',
      render: (text, record) => (
        <Space style={{ padding: '0 14px', whiteSpace: 'normal', wordBreak: 'break-word' }}>
          <CarOutlined style={{ color: '#697565' }} />
          <span>{text || `Vehicle ${record.vehicle_id?.slice(0, 8)}` || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      align: 'center',
      width: '15%',
      filters: [
        { text: 'Draft', value: 'DRAFT' },
        { text: 'Submitted', value: 'SUBMITTED' },
        { text: 'Reviewing', value: 'REVIEWING' },
        { text: 'Request Info', value: 'REQUEST_INFO' },
        { text: 'Approved', value: 'APPROVED' },
        { text: 'Partially Approved', value: 'PARTIALLY_APPROVED' },
        { text: 'Rejected', value: 'REJECTED' },
        { text: 'Cancelled', value: 'CANCELLED' },
      ],
      filteredValue: filteredInfo.status || null,
      onFilter: (value, record) => record.status === value,
      render: (status) => {
        return <Space>{CLAIM_STATUS_LABELS[status] || status || 'Unknown'}</Space>
      },
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Total Cost</span>,
      dataIndex: 'total_cost',
      key: 'total_cost',
      width: '15%',
      align: 'right',
      sorter: (a, b) => (a.total_cost || 0) - (b.total_cost || 0),
      sortOrder: sortedInfo.columnKey === 'total_cost' ? sortedInfo.order : null,
      render: (cost) => (
        <Space style={{ padding: '0 14px' }}>
          <span>{`${cost.toLocaleString('vi-VN', { style: 'currency', currency: 'VND' })} `}</span>
        </Space>
      ),
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Created Date</span>,
      dataIndex: 'created_at',
      key: 'created_at',
      width: '15%',
      align: 'center',
      sorter: (a, b) => new Date(a.created_at || 0) - new Date(b.created_at || 0),
      sortOrder: sortedInfo.columnKey === 'created_at' ? sortedInfo.order : null,
      render: (date) => (
        <Space style={{ padding: '0 14px' }}>
          <CalendarOutlined style={{ color: '#697565' }} />
          <span>{date ? new Date(date).toLocaleDateString('vi-VN') : 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Details',
      key: 'actions',
      width: '5%',
      align: 'center',
      render: (_, record) => (
        <Space size="middle">
          <Button
            type="text"
            icon={<EyeOutlined />}
            onClick={() => onViewDetails(record)}
            title="View Details"
            style={{ color: '#1890ff' }}
          />
        </Space>
      ),
    },
  ]
}

export default GenerateColumns
