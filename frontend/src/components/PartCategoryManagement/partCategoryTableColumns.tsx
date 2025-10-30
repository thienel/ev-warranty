import { Button, Popconfirm, Space, Tag } from 'antd'
import { DeleteOutlined, EditOutlined, AppstoreOutlined, FileTextOutlined } from '@ant-design/icons'
import { type PartCategory } from '@/types/index.js'

type OnOpenModal = (
  item?: (Record<string, unknown> & { id: string | number }) | null,
  isUpdate?: boolean,
) => void
type OnDelete = (itemId: string | number) => Promise<void>

const GenerateColumns = (
  sortedInfo: Record<string, unknown>,
  _filteredInfo: Record<string, unknown>,
  onOpenModal: OnOpenModal,
  onDelete: OnDelete,
) => {
  return [
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Category Name</span>,
      dataIndex: 'category_name',
      key: 'category_name',
      width: '25%',
      sorter: (a: Record<string, unknown>, b: Record<string, unknown>) => {
        const aCategory = a as unknown as PartCategory
        const bCategory = b as unknown as PartCategory
        return (aCategory.category_name || '').localeCompare(bCategory.category_name || '')
      },
      sortOrder:
        sortedInfo.columnKey === 'category_name'
          ? (sortedInfo.order as 'ascend' | 'descend' | null)
          : null,
      render: (text: string) => (
        <Space
          style={{
            padding: '0 14px',
            whiteSpace: 'normal',
            wordBreak: 'break-word',
          }}
        >
          <AppstoreOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: <span style={{ padding: '0 14px', display: 'inline-block' }}>Description</span>,
      dataIndex: 'description',
      key: 'description',
      width: '30%',
      ellipsis: true,
      render: (text: string) => (
        <Space style={{ padding: '0 14px', whiteSpace: 'normal', wordBreak: 'break-word' }}>
          <FileTextOutlined style={{ color: '#697565' }} />
          <span>{text || 'N/A'}</span>
        </Space>
      ),
    },
    {
      title: 'Parent Category',
      dataIndex: 'parent_category_name',
      key: 'parent_category_name',
      align: 'center' as const,
      width: '25%',
      render: (text: string) => <span>{text || <Tag color="default">Root Category</Tag>}</span>,
    },
    {
      title: 'Actions',
      key: 'actions',
      align: 'center' as const,
      width: '20%',
      render: (_: unknown, record: Record<string, unknown>) => {
        const category = record as unknown as PartCategory
        return (
          <Space size="small">
            <Button
              type="text"
              size="small"
              icon={<EditOutlined />}
              onClick={() =>
                onOpenModal(record as Record<string, unknown> & { id: string | number }, true)
              }
              style={{ color: '#1890ff' }}
            >
              Edit
            </Button>
            <Popconfirm
              title="Delete Part Category"
              description="Are you sure you want to delete this part category?"
              onConfirm={() => onDelete(category.id)}
              okText="Yes"
              cancelText="No"
              placement="topRight"
            >
              <Button type="text" size="small" icon={<DeleteOutlined />} danger>
                Delete
              </Button>
            </Popconfirm>
          </Space>
        )
      },
    },
  ]
}

export default GenerateColumns
