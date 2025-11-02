import { CLAIM_STATUS_LABELS } from '@constants/common-constants.js'
import { Button, Space } from 'antd'
import { EyeOutlined, UserOutlined, CarOutlined, CalendarOutlined } from '@ant-design/icons'
import { type Claim } from '@/types/index.js'

type OnViewDetails = (claim: Claim) => void
type OnOpenModal = (
  item?: (Record<string, unknown> & { id: string | number }) | null,
  isUpdate?: boolean,
) => void
type OnDelete = (itemId: string | number) => Promise<void>

interface AdditionalProps {
  onViewDetails: OnViewDetails
}

const GenerateColumns = (
  sortedInfo: Record<string, unknown>,
  filteredInfo: Record<string, unknown>,
  _onOpenModal: OnOpenModal,
  _onDelete: OnDelete,
  additionalProps: Record<string, unknown>,
) => {
  const { onViewDetails } = additionalProps as unknown as AdditionalProps

  return [
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Customer Name</span>,
      dataIndex: 'customer_name',
      key: 'customer_name',
      width: '25%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aClaim = a as unknown as Claim & { customer_name?: string }
        const bClaim = b as unknown as Claim & { customer_name?: string }
        return (aClaim.customer_name || aClaim.customer_id || '')
          .toString()
          .localeCompare((bClaim.customer_name || bClaim.customer_id || '').toString())
      },
      sortOrder:
        sortedInfo.columnKey === 'customer_name'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (_text: string, record: Record<string, unknown>) => {
        const claim = record as unknown as Claim & { customer_name?: string }
        const displayName =
          claim.customer_name || `Customer ${claim.customer_id?.slice(0, 8)}` || 'N/A'
        return (
          <Space
            style={{
              padding: '0 14px',
              whiteSpace: 'normal',
              wordBreak: 'break-word',
            }}
          >
            <UserOutlined style={{ color: '#697565' }} />
            <span>{displayName}</span>
          </Space>
        )
      },
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Vehicle</span>,
      dataIndex: 'vehicle_info',
      key: 'vehicle_info',
      width: '25%',
      render: (text: string, record: Record<string, unknown>) => {
        const claim = record as unknown as Claim
        return (
          <Space
            style={{
              padding: '0 14px',
              whiteSpace: 'normal',
              wordBreak: 'break-word',
            }}
          >
            <CarOutlined style={{ color: '#697565' }} />
            <span>{text || `Vehicle ${claim.vehicle_id?.slice(0, 8)}` || 'N/A'}</span>
          </Space>
        )
      },
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      align: 'center' as const,
      width: '15%',
      filters: [
        { text: 'Submitted', value: 'SUBMITTED' },
        { text: 'Processing', value: 'PROCESSING' },
        { text: 'Approved', value: 'APPROVED' },
        { text: 'Rejected', value: 'REJECTED' },
        { text: 'Completed', value: 'COMPLETED' },
      ],
      filteredValue: (filteredInfo.status as React.Key[] | null) || null,
      onFilter: (value: boolean | React.Key, record: Record<string, unknown>) => {
        const claim = record as unknown as Claim
        return claim.status === value
      },
      render: (status: string) => {
        return (
          <Space>
            {CLAIM_STATUS_LABELS[status as keyof typeof CLAIM_STATUS_LABELS] || status || 'Unknown'}
          </Space>
        )
      },
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Total Cost</span>,
      dataIndex: 'total_cost',
      key: 'total_cost',
      width: '15%',
      align: 'right' as const,
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aClaim = a as unknown as Claim
        const bClaim = b as unknown as Claim
        return (aClaim.total_cost || 0) - (bClaim.total_cost || 0)
      },
      sortOrder:
        sortedInfo.columnKey === 'total_cost'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (cost: number) => (
        <Space style={{ padding: '0 14px' }}>
          <span>
            {cost.toLocaleString('vi-VN', {
              style: 'currency',
              currency: 'VND',
            })}
          </span>
        </Space>
      ),
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Created Date</span>,
      dataIndex: 'created_at',
      key: 'created_at',
      width: '15%',
      align: 'center' as const,
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aClaim = a as unknown as Claim
        const bClaim = b as unknown as Claim
        return (
          new Date(aClaim.created_at || 0).getTime() - new Date(bClaim.created_at || 0).getTime()
        )
      },
      sortOrder:
        sortedInfo.columnKey === 'created_at'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (date: string) => (
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
      align: 'center' as const,
      render: (_: unknown, record: Record<string, unknown>) => {
        const claim = record as unknown as Claim
        return (
          <Space size="middle">
            <Button
              type="text"
              icon={<EyeOutlined />}
              onClick={() => onViewDetails(claim)}
              title="View Details"
              style={{ color: '#1890ff' }}
            />
          </Space>
        )
      },
    },
  ]
}

export default GenerateColumns
