import { Button, Popconfirm, Space, Tag } from "antd";
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  DeleteOutlined,
  EditOutlined,
  EnvironmentOutlined,
  BankOutlined,
} from "@ant-design/icons";
import { type Office } from "@/types/index.js";

type OnOpenModal = (item?: (Record<string, unknown> & { id: string | number }) | null, isUpdate?: boolean) => void;
type OnDelete = (itemId: string | number) => Promise<void>;

const OFFICE_TYPE_LABELS = {
  evm: "EVM",
  sc: "Service Center",
} as const;

const GenerateColumns = (
  sortedInfo: Record<string, unknown>,
  filteredInfo: Record<string, unknown>,
  onOpenModal: OnOpenModal,
  onDelete: OnDelete
) => {
  return [
    {
      title: (
        <span style={{ padding: "0 14px", display: "inline-block" }}>
          Office Name
        </span>
      ),
      dataIndex: "office_name",
      key: "office_name",
      width: "25%",
      sorter: (a: Office, b: Office) =>
        (a.office_name || "").localeCompare(b.office_name || ""),
      sortOrder:
        sortedInfo.columnKey === "office_name"
          ? (sortedInfo.order as "ascend" | "descend" | null)
          : null,
      render: (text: string) => (
        <Space
          style={{
            padding: "0 14px",
            whiteSpace: "normal",
            wordBreak: "break-word",
          }}
        >
          <BankOutlined style={{ color: "#697565" }} />
          <span>{text || "N/A"}</span>
        </Space>
      ),
    },
    {
      title: "Office Type",
      dataIndex: "office_type",
      key: "office_type",
      align: "center" as const,
      width: "15%",
      filters: [
        { text: "EVM", value: "evm" },
        { text: "Service Center", value: "sc" },
      ],
      filteredValue: (filteredInfo.office_type as React.Key[] | null) || null,
      onFilter: (value: string | number | boolean, record: Office) =>
        record.office_type === value,
      render: (office_type: "evm" | "sc") => {
        const label = OFFICE_TYPE_LABELS[office_type] || office_type;
        return <span>{label}</span>;
      },
    },
    {
      title: (
        <span style={{ padding: "0 14px", display: "inline-block" }}>
          Address
        </span>
      ),
      dataIndex: "address",
      key: "address",
      width: "30%",
      sorter: (a: Office, b: Office) =>
        (a.address || "").localeCompare(b.address || ""),
      sortOrder:
        sortedInfo.columnKey === "address"
          ? (sortedInfo.order as "ascend" | "descend" | null)
          : null,
      ellipsis: true,
      render: (text: string) => (
        <Space style={{ padding: "0 14px" }}>
          <EnvironmentOutlined style={{ color: "#697565" }} />
          <span>{text || "N/A"}</span>
        </Space>
      ),
    },
    {
      title: "Status",
      dataIndex: "is_active",
      key: "is_active",
      align: "center" as const,
      width: "12%",
      filters: [
        { text: "Active", value: true },
        { text: "Inactive", value: false },
      ],
      filteredValue: (filteredInfo.is_active as React.Key[] | null) || null,
      onFilter: (value: string | number | boolean, record: Office) =>
        record.is_active === value,
      render: (is_active: boolean) => {
        const color = is_active ? "green" : "red";
        const icon = is_active ? (
          <CheckCircleOutlined />
        ) : (
          <CloseCircleOutlined />
        );
        const text = is_active ? "Active" : "Inactive";
        return (
          <Tag color={color} icon={icon}>
            {text}
          </Tag>
        );
      },
    },
    {
      title: "Actions",
      key: "actions",
      align: "center" as const,
      width: "18%",
      render: (_: unknown, record: Office) => (
        <Space size="small">
          <Button
            type="text"
            size="small"
            icon={<EditOutlined />}
            onClick={() => onOpenModal(record as unknown as Record<string, unknown> & { id: string | number }, true)}
            style={{ color: "#1890ff" }}
          >
            Edit
          </Button>
          <Popconfirm
            title="Delete Office"
            description="Are you sure you want to delete this office?"
            onConfirm={() => onDelete(record.id)}
            okText="Yes"
            cancelText="No"
            placement="topRight"
          >
            <Button type="text" size="small" icon={<DeleteOutlined />} danger>
              Delete
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];
};

export default GenerateColumns;
