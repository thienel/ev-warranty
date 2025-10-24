import React, { useEffect } from "react";
import {
  Modal,
  Button,
  Form,
  Input,
  message,
  Select,
  Space,
  Switch,
} from "antd";
import {
  HomeOutlined,
  LockOutlined,
  MailOutlined,
  UserOutlined,
} from "@ant-design/icons";
import {
  API_ENDPOINTS,
  PASSWORD_RULES,
  ROLE_LABELS,
  USER_ROLES,
} from "@constants/common-constants.js";
import { type UserModalProps, type UserFormData } from "@/types/index.js";
import api from "@services/api.js";
import useHandleApiError from "@/hooks/useHandleApiError.js";

const { Option } = Select;

const UserModal: React.FC<UserModalProps> = ({
  loading,
  setLoading,
  onClose,
  user = null,
  opened = false,
  offices,
  officesLoading = false,
  isUpdate,
}) => {
  const [form] = Form.useForm<UserFormData>();
  const handleError = useHandleApiError();

  // Debug log to track offices prop changes
  useEffect(() => {
    console.log("UserModal: offices prop changed:", {
      isArray: Array.isArray(offices),
      length: Array.isArray(offices) ? offices.length : "N/A",
      data: offices,
      opened,
      officesLoading,
      isUpdate,
    });
  }, [offices, opened, officesLoading, isUpdate]);

  // Populate form when user prop changes or modal opens
  useEffect(() => {
    if (opened) {
      if (user && isUpdate) {
        // When editing, populate form with user data
        form.setFieldsValue(user);
      } else {
        // When creating new, reset to default values
        form.setFieldsValue({ is_active: true });
      }
    }
  }, [form, user, isUpdate, opened]);

  // Clear form when modal closes
  useEffect(() => {
    if (!opened) {
      form.resetFields();
    }
  }, [form, opened]);

  const handleSubmit = async (values: UserFormData): Promise<void> => {
    setLoading(true);
    try {
      const payload = { ...values };

      if (!isUpdate && !payload.password) {
        message.warning("Password is required for new user");
        setLoading(false);
        return;
      }

      if (isUpdate) {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const { password, email, ...updatePayload } = payload;
        await api.put(`${API_ENDPOINTS.USERS}/${user?.id}`, updatePayload);
        message.success("User updated successfully");
      } else {
        await api.post(API_ENDPOINTS.USERS, payload);
        message.success("User created successfully");
      }

      onClose();
    } catch (error) {
      handleError(error as Error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      title={
        <Space style={{ margin: "14px 0" }}>
          {isUpdate ? "Edit User" : "Add New User"}
        </Space>
      }
      open={opened}
      onCancel={onClose}
      style={{ margin: "auto" }}
      footer={null}
      width={500}
      destroyOnHidden
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        autoComplete="off"
        key={user?.id || "new"}
      >
        <Form.Item
          label="Full Name"
          name="name"
          validateFirst
          rules={[
            { required: true, message: "Please enter full name" },
            { min: 2, message: "Name must be at least 2 characters" },
            { max: 50, message: "Name cannot exceed 50 characters" },
            {
              pattern: /^[\p{L}\s'-]+$/u,
              message:
                "Name can only contain letters, spaces, apostrophes or hyphens",
            },
          ]}
        >
          <Input
            placeholder="Enter full name"
            prefix={<UserOutlined />}
            size="large"
          />
        </Form.Item>
        <Form.Item
          label="Email"
          name="email"
          rules={[
            { required: true, message: "Please enter email" },
            { type: "email", message: "Invalid email format" },
            { max: 100, message: "Email cannot exceed 100 characters" },
          ]}
        >
          <Input
            placeholder="Enter email"
            prefix={<MailOutlined />}
            size="large"
            disabled={isUpdate}
          />
        </Form.Item>

        {!isUpdate && (
          <Form.Item
            label="Password"
            name="password"
            validateFirst
            rules={PASSWORD_RULES}
          >
            <Input.Password
              placeholder="Enter password"
              prefix={<LockOutlined />}
              size="large"
            />
          </Form.Item>
        )}

        <Form.Item
          label="Role"
          name="role"
          rules={[{ required: true, message: "Please select a role" }]}
        >
          <Select placeholder="Select role" size="large">
            {Object.entries(USER_ROLES).map(([key, value]) => (
              <Option key={key} value={value}>
                <Space>
                  <UserOutlined />
                  {ROLE_LABELS[value]}
                </Space>
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Office"
          name="office_id"
          rules={[{ required: true, message: "Please select an office" }]}
        >
          <Select
            placeholder={
              officesLoading
                ? "Loading offices..."
                : offices.length === 0
                  ? "No offices available"
                  : "Select office"
            }
            size="large"
            showSearch
            optionFilterProp="children"
            loading={officesLoading}
            disabled={officesLoading}
            filterOption={(input, option) =>
              option?.children
                ?.toString()
                .toLowerCase()
                .includes(input.toLowerCase()) ?? false
            }
          >
            {offices.map((office) => (
              <Option key={office.id} value={office.id}>
                <Space>
                  <HomeOutlined />
                  {office.office_name}
                </Space>
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item label="Status" name="is_active" valuePropName="checked">
          <Switch
            checkedChildren="Active"
            unCheckedChildren="Inactive"
            size="default"
          />
        </Form.Item>

        <Form.Item
          style={{ marginBottom: 0, textAlign: "right", marginTop: "24px" }}
        >
          <Space>
            <Button size="large" onClick={onClose}>
              Cancel
            </Button>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              size="large"
            >
              {isUpdate ? "Update User" : "Create User"}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default UserModal;
