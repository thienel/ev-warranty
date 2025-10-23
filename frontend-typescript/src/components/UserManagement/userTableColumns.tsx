import { ROLE_LABELS, USER_ROLES } from "@constants/common-constants.js";
import { Button, Popconfirm, Space, Tag } from "antd";
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  DeleteOutlined,
  EditOutlined,
  HomeOutlined,
  MailOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { type User } from "@/types/index.js";

type OnOpenModal = (item?: (Record<string, unknown> & { id: string | number }) | null, isUpdate?: boolean) => void;
type OnDelete = (itemId: string | number) => Promise<void>;

const GenerateColumns = (
  sortedInfo: Record<string, unknown>,
  filteredInfo: Record<string, unknown>,
  onOpenModal: OnOpenModal,
  onDelete: OnDelete,
  additionalProps: Record<string, unknown>
) => {
  const { getOfficeName } = additionalProps as {
    getOfficeName?: (officeId: string) => string;
  };

  return [
    {
      title: (
        <span style={{ padding: "0 14px", display: "inline-block" }}>Name</span>
      ),
      dataIndex: "name",
      key: "name",
      width: "20%",
      sorter: (a: User, b: User) => (a.name || "").localeCompare(b.name || ""),
      sortOrder:
        sortedInfo.columnKey === "name"
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
          <UserOutlined style={{ color: "#697565" }} />
          <span>{text || "N/A"}</span>
        </Space>
      ),
    },
    {
      title: (
        <span style={{ padding: "0 14px", display: "inline-block" }}>
          Email
        </span>
      ),
      dataIndex: "email",
      key: "email",
      width: "22%",
      sorter: (a: User, b: User) =>
        (a.email || "").localeCompare(b.email || ""),
      sortOrder:
        sortedInfo.columnKey === "email"
          ? (sortedInfo.order as "ascend" | "descend" | null)
          : null,
      ellipsis: true,
      render: (text: string) => (
        <Space style={{ padding: "0 14px" }}>
          <MailOutlined style={{ color: "#697565" }} />
          <span>{text || "N/A"}</span>
        </Space>
      ),
    },
    {
      title: "Role",
      dataIndex: "role",
      key: "role",
      align: "center" as const,
      width: "15%",
      filters: Object.values(USER_ROLES).map((role) => ({
        text: ROLE_LABELS[role],
        value: role,
      })),
      filteredValue: (filteredInfo.role as React.Key[] | null) || null,
      onFilter: (value: string | number | boolean, record: User) =>
        record.role === value,
      render: (role: string) => {
        return (
          <Space>{ROLE_LABELS[role as keyof typeof ROLE_LABELS] || role}</Space>
        );
      },
    },
    {
      title: "Office",
      dataIndex: "office_id",
      key: "office_id",
      align: "center" as const,
      width: "20%",
      render: (officeId: string) => (
        <Space>
          <HomeOutlined style={{ color: "#697565" }} />
          <span>{getOfficeName?.(officeId) || "N/A"}</span>
        </Space>
      ),
    },
    {
      title: "Status",
      dataIndex: "is_active",
      key: "is_active",
      align: "center" as const,
      width: "13%",
      filters: [
        { text: "Active", value: true },
        { text: "Inactive", value: false },
      ],
      filteredValue: (filteredInfo.is_active as React.Key[] | null) || null,
      onFilter: (value: string | number | boolean, record: User) =>
        record.is_active === value,
      render: (isActive: boolean) => (
        <Tag
          icon={isActive ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
          color={isActive ? "green" : "red"}
        >
          {isActive ? "Active" : "Inactive"}
        </Tag>
      ),
    },
    {
      title: "Actions",
      key: "action",
      fixed: "right" as const,
      align: "center" as const,
      width: "10%",
      render: (_: unknown, record: User) => (
        <Space size="small">
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => onOpenModal(record as unknown as Record<string, unknown> & { id: string | number }, true)}
          />
          <Popconfirm
            title="Delete user"
            description="Are you sure you want to delete this user?"
            onConfirm={() => onDelete(record.id)}
            okText="Delete"
            cancelText="Cancel"
            okButtonProps={{ danger: true }}
          >
            <Button type="text" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];
};

export default GenerateColumns;
