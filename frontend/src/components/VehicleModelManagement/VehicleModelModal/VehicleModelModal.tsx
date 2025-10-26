import React, { useEffect } from "react";
import { Modal, Button, Form, Input, message, Space, InputNumber } from "antd";
import { CarOutlined, TagOutlined, CalendarOutlined } from "@ant-design/icons";
import {
  type VehicleModelModalProps,
  type VehicleModelFormData,
} from "@/types/index";
import { vehicleModelsApi } from "@services/index";
import useHandleApiError from "@/hooks/useHandleApiError";

const VehicleModelModal: React.FC<VehicleModelModalProps> = ({
  loading,
  setLoading,
  onClose,
  vehicleModel = null,
  opened = false,
  isUpdate,
}) => {
  const [form] = Form.useForm<VehicleModelFormData>();
  const handleError = useHandleApiError();

  // Populate form when vehicleModel prop changes or modal opens
  useEffect(() => {
    if (opened) {
      if (vehicleModel && isUpdate) {
        // When editing, populate form with vehicle model data
        form.setFieldsValue({
          brand: vehicleModel.brand,
          model_name: vehicleModel.model_name,
          year: vehicleModel.year,
        });
      } else {
        // When creating new, reset form
        form.resetFields();
      }
    }
  }, [form, vehicleModel, isUpdate, opened]);

  // Clear form when modal closes
  useEffect(() => {
    if (!opened) {
      form.resetFields();
    }
  }, [form, opened]);

  const handleSubmit = async (values: VehicleModelFormData): Promise<void> => {
    setLoading(true);
    try {
      const payload = { ...values };

      if (isUpdate) {
        await vehicleModelsApi.update(vehicleModel?.id || "", payload);
        message.success("Vehicle model updated successfully");
      } else {
        await vehicleModelsApi.create(payload);
        message.success("Vehicle model created successfully");
      }

      onClose();
    } catch (error) {
      handleError(error as Error);
    } finally {
      setLoading(false);
    }
  };

  const currentYear = new Date().getFullYear();

  return (
    <Modal
      title={
        <Space style={{ margin: "14px 0" }}>
          {isUpdate ? "Edit Vehicle Model" : "Add New Vehicle Model"}
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
        key={vehicleModel?.id || "new"}
      >
        <Form.Item
          label="Brand"
          name="brand"
          validateFirst
          rules={[
            { required: true, message: "Please enter brand name" },
            { min: 2, message: "Brand name must be at least 2 characters" },
            { max: 50, message: "Brand name cannot exceed 50 characters" },
            {
              pattern: /^[\p{L}\d\s&.-]+$/u,
              message:
                "Brand name can only contain letters, numbers, spaces, and common symbols (&.-)",
            },
          ]}
        >
          <Input
            placeholder="Enter brand name (e.g., Toyota, BMW)"
            prefix={<TagOutlined />}
            size="large"
          />
        </Form.Item>

        <Form.Item
          label="Model Name"
          name="model_name"
          validateFirst
          rules={[
            { required: true, message: "Please enter model name" },
            { min: 1, message: "Model name must be at least 1 character" },
            { max: 50, message: "Model name cannot exceed 50 characters" },
            {
              pattern: /^[\p{L}\d\s&.-]+$/u,
              message:
                "Model name can only contain letters, numbers, spaces, and common symbols (&.-)",
            },
          ]}
        >
          <Input
            placeholder="Enter model name (e.g., Camry, X5)"
            prefix={<CarOutlined />}
            size="large"
          />
        </Form.Item>

        <Form.Item
          label="Year"
          name="year"
          validateFirst
          rules={[
            { required: true, message: "Please enter manufacturing year" },
            {
              type: "number",
              min: 2000,
              max: currentYear,
              message: `Year must be between 2000 and ${currentYear}`,
            },
          ]}
        >
          <InputNumber
            placeholder={`Enter year (2000-${currentYear})`}
            prefix={<CalendarOutlined />}
            size="large"
            style={{ width: "100%" }}
            min={2000}
            max={currentYear}
            precision={0}
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
              {isUpdate ? "Update Vehicle Model" : "Create Vehicle Model"}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default VehicleModelModal;
