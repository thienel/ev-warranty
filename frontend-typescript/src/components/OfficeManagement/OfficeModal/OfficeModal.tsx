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
import { BankOutlined } from "@ant-design/icons";
import { API_ENDPOINTS, OFFICE_TYPES } from "@constants/common-constants.js";
import { type OfficeModalProps, type OfficeFormData } from "@/types/index.js";
import api from "@services/api.js";
import useHandleApiError from "@/hooks/useHandleApiError.js";

const { Option } = Select;

const OfficeModal: React.FC<OfficeModalProps> = ({
  loading,
  setLoading,
  onClose,
  office = null,
  opened = false,
  isUpdate,
}) => {
  const [form] = Form.useForm<OfficeFormData>();
  const handleError = useHandleApiError();

  // Populate form when office prop changes or modal opens
  useEffect(() => {
    if (opened) {
      if (office && isUpdate) {
        // When editing, populate form with office data
        const formData: OfficeFormData = {
          office_name: office.office_name,
          office_type: office.office_type as "EVM" | "SC",
          address: office.address,
          is_active: office.is_active,
        };
        form.setFieldsValue(formData);
      } else {
        // When creating new, reset to default values
        form.setFieldsValue({ is_active: true });
      }
    }
  }, [form, office, isUpdate, opened]);

  // Clear form when modal closes
  useEffect(() => {
    if (!opened) {
      form.resetFields();
    }
  }, [form, opened]);

  const handleSubmit = async (values: OfficeFormData): Promise<void> => {
    setLoading(true);
    console.log("Submitting office data:", values);
    try {
      const payload = { ...values };

      if (isUpdate) {
        await api.put(`${API_ENDPOINTS.OFFICES}/${office?.id}`, payload);
        message.success("Office updated successfully");
      } else {
        await api.post(API_ENDPOINTS.OFFICES, payload);
        message.success("Office created successfully");
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
          {isUpdate ? "Edit Office" : "Add New Office"}
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
        key={office?.id || "new"}
      >
        <Form.Item
          label="Office Name"
          name="office_name"
          validateFirst
          rules={[
            { required: true, message: "Please enter office name" },
            { min: 2, message: "Office name must be at least 2 characters" },
            { max: 100, message: "Office name cannot exceed 100 characters" },
          ]}
        >
          <Input
            placeholder="Enter office name"
            prefix={<BankOutlined />}
            size="large"
          />
        </Form.Item>

        <Form.Item
          label="Office Type"
          name="office_type"
          rules={[{ required: true, message: "Please select office type" }]}
        >
          <Select placeholder="Select office type" size="large">
            <Option value={OFFICE_TYPES.EVM}>
              <Space>
                <BankOutlined />
                EVM
              </Space>
            </Option>
            <Option value={OFFICE_TYPES.SC}>
              <Space>
                <BankOutlined />
                Service Center
              </Space>
            </Option>
          </Select>
        </Form.Item>

        <Form.Item
          label="Address"
          name="address"
          validateFirst
          rules={[
            { required: true, message: "Please enter address" },
            { min: 5, message: "Address must be at least 5 characters" },
            { max: 200, message: "Address cannot exceed 200 characters" },
          ]}
        >
          <Input.TextArea
            placeholder="Enter office address"
            rows={3}
            size="large"
          />
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
              {isUpdate ? "Update Office" : "Create Office"}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default OfficeModal;
