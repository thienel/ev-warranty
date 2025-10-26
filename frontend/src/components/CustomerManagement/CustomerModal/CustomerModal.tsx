import React, { useEffect } from "react";
import { Modal, Button, Form, Input, message, Space } from "antd";
import { UserOutlined, MailOutlined, PhoneOutlined } from "@ant-design/icons";
import { type CustomerModalProps, type CustomerFormData } from "@/types/index";
import { customersApi } from "@services/index";
import useHandleApiError from "@/hooks/useHandleApiError";

const CustomerModal: React.FC<CustomerModalProps> = ({
  loading,
  setLoading,
  onClose,
  customer = null,
  opened = false,
  isUpdate,
}) => {
  const [form] = Form.useForm<CustomerFormData>();
  const handleError = useHandleApiError();

  // Populate form when customer prop changes or modal opens
  useEffect(() => {
    if (opened) {
      if (customer && isUpdate) {
        // When editing, populate form with customer data
        form.setFieldsValue({
          first_name: customer.first_name,
          last_name: customer.last_name,
          email: customer.email || "",
          phone_number: customer.phone_number || "",
          address: customer.address || "",
        });
      } else {
        // When creating new, reset form
        form.resetFields();
      }
    }
  }, [form, customer, isUpdate, opened]);

  // Clear form when modal closes
  useEffect(() => {
    if (!opened) {
      form.resetFields();
    }
  }, [form, opened]);

  const handleSubmit = async (values: CustomerFormData): Promise<void> => {
    setLoading(true);
    try {
      const payload = { ...values };

      if (isUpdate) {
        await customersApi.update(customer?.id || '', payload);
        message.success("Customer updated successfully");
      } else {
        await customersApi.create(payload);
        message.success("Customer created successfully");
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
          {isUpdate ? "Edit Customer" : "Add New Customer"}
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
        key={customer?.id || "new"}
      >
        <Form.Item
          label="First Name"
          name="first_name"
          validateFirst
          rules={[
            { required: true, message: "Please enter first name" },
            { min: 2, message: "First name must be at least 2 characters" },
            { max: 50, message: "First name cannot exceed 50 characters" },
            {
              pattern: /^[\p{L}\s'-]+$/u,
              message:
                "First name can only contain letters, spaces, apostrophes or hyphens",
            },
          ]}
        >
          <Input
            placeholder="Enter first name"
            prefix={<UserOutlined />}
            size="large"
          />
        </Form.Item>

        <Form.Item
          label="Last Name"
          name="last_name"
          validateFirst
          rules={[
            { required: true, message: "Please enter last name" },
            { min: 2, message: "Last name must be at least 2 characters" },
            { max: 50, message: "Last name cannot exceed 50 characters" },
            {
              pattern: /^[\p{L}\s'-]+$/u,
              message:
                "Last name can only contain letters, spaces, apostrophes or hyphens",
            },
          ]}
        >
          <Input
            placeholder="Enter last name"
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
          />
        </Form.Item>

        <Form.Item
          label="Phone Number"
          name="phone_number"
          rules={[
            { required: true, message: "Please enter phone number" },
            {
              pattern:
                /^(0|\+84)(3[2-9]|5[6|8|9]|7[0|6-9]|8[1-9]|9[0-9])[0-9]{7}$/,
              message: "Invalid phone number format",
            },
          ]}
        >
          <Input
            placeholder="Enter phone number"
            prefix={<PhoneOutlined />}
            size="large"
          />
        </Form.Item>

        <Form.Item
          label="Address"
          name="address"
          rules={[
            { required: true, message: "Please enter address" },
            { max: 255, message: "Address cannot exceed 255 characters" },
          ]}
        >
          <Input.TextArea placeholder="Enter address" rows={3} size="large" />
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
              {isUpdate ? "Update Customer" : "Create Customer"}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CustomerModal;
