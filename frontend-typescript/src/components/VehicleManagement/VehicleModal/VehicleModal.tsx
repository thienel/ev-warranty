import React, { useEffect } from "react";
import {
  Modal,
  Button,
  Form,
  Input,
  message,
  Select,
  Space,
  DatePicker,
} from "antd";
import { CarOutlined, UserOutlined, FileTextOutlined } from "@ant-design/icons";
import { API_ENDPOINTS } from "@constants/common-constants";
import { type VehicleModalProps, type VehicleFormData } from "@/types/index";
import api from "@services/api";
import useHandleApiError from "@/hooks/useHandleApiError";
import dayjs from "dayjs";

const { Option } = Select;

const VehicleModal: React.FC<VehicleModalProps> = ({
  loading,
  setLoading,
  onClose,
  vehicle = null,
  opened = false,
  customers,
  vehicleModels,
  customersLoading = false,
  vehicleModelsLoading = false,
  isUpdate,
}) => {
  const [form] = Form.useForm<VehicleFormData>();
  const handleError = useHandleApiError();

  // Populate form when vehicle prop changes or modal opens
  useEffect(() => {
    if (opened) {
      if (vehicle && isUpdate) {
        // When editing, populate form with vehicle data
        form.setFieldsValue({
          vin: vehicle.vin,
          license_plate: vehicle.license_plate || "",
          customer_id: vehicle.customer_id,
          model_id: vehicle.model_id,
          purchase_date: vehicle.purchase_date
            ? dayjs(vehicle.purchase_date)
            : undefined,
        });
      } else {
        // When creating new, reset form
        form.resetFields();
      }
    }
  }, [form, vehicle, isUpdate, opened]);

  // Clear form when modal closes
  useEffect(() => {
    if (!opened) {
      form.resetFields();
    }
  }, [form, opened]);

  const handleSubmit = async (values: VehicleFormData): Promise<void> => {
    setLoading(true);
    try {
      const payload = {
        ...values,
        purchase_date: values.purchase_date
          ? (values.purchase_date as { format: (f: string) => string }).format(
              "YYYY-MM-DD"
            )
          : undefined,
      };

      if (isUpdate) {
        await api.put(`${API_ENDPOINTS.VEHICLES}/${vehicle?.id}`, payload);
        message.success("Vehicle updated successfully");
      } else {
        await api.post(API_ENDPOINTS.VEHICLES, payload);
        message.success("Vehicle created successfully");
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
          {isUpdate ? "Edit Vehicle" : "Add New Vehicle"}
        </Space>
      }
      open={opened}
      onCancel={onClose}
      style={{ margin: "auto" }}
      footer={null}
      width={600}
      destroyOnHidden
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        autoComplete="off"
        key={vehicle?.id || "new"}
      >
        <Form.Item
          label="VIN"
          name="vin"
          validateFirst
          rules={[
            { required: true, message: "Please enter VIN" },
            { min: 17, max: 17, message: "VIN must be exactly 17 characters" },
            {
              pattern: /^[A-HJ-NPR-Z0-9]{17}$/,
              message: "Invalid VIN format (no I, O, Q allowed)",
            },
          ]}
        >
          <Input
            placeholder="Enter VIN (17 characters)"
            prefix={<FileTextOutlined />}
            size="large"
            maxLength={17}
            style={{ textTransform: "uppercase" }}
            onChange={(e) => {
              const value = e.target.value.toUpperCase();
              form.setFieldValue("vin", value);
            }}
          />
        </Form.Item>

        <Form.Item
          label="License Plate"
          name="license_plate"
          rules={[
            { max: 20, message: "License plate cannot exceed 20 characters" },
          ]}
        >
          <Input
            placeholder="Enter license plate (optional)"
            prefix={<CarOutlined />}
            size="large"
          />
        </Form.Item>

        <Form.Item
          label="Customer"
          name="customer_id"
          rules={[{ required: true, message: "Please select a customer" }]}
        >
          <Select
            placeholder={
              customersLoading
                ? "Loading customers..."
                : customers.length === 0
                  ? "No customers available"
                  : "Select customer"
            }
            size="large"
            showSearch
            optionFilterProp="children"
            loading={customersLoading}
            disabled={customersLoading}
            filterOption={(input, option) =>
              option?.children
                ?.toString()
                .toLowerCase()
                .includes(input.toLowerCase()) ?? false
            }
          >
            {customers.map((customer) => (
              <Option key={customer.id} value={customer.id}>
                <Space>
                  <UserOutlined />
                  {`${customer.first_name || ""} ${customer.last_name || ""}`.trim()}
                </Space>
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Vehicle Model"
          name="model_id"
          rules={[{ required: true, message: "Please select a vehicle model" }]}
        >
          <Select
            placeholder={
              vehicleModelsLoading
                ? "Loading vehicle models..."
                : vehicleModels.length === 0
                  ? "No vehicle models available"
                  : "Select vehicle model"
            }
            size="large"
            showSearch
            optionFilterProp="children"
            loading={vehicleModelsLoading}
            disabled={vehicleModelsLoading}
            filterOption={(input, option) =>
              option?.children
                ?.toString()
                .toLowerCase()
                .includes(input.toLowerCase()) ?? false
            }
          >
            {vehicleModels.map((model) => (
              <Option key={model.id} value={model.id}>
                <Space>
                  <CarOutlined />
                  {`${model.brand} ${model.model_name} (${model.year})`}
                </Space>
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Purchase Date"
          name="purchase_date"
          rules={[
            {
              validator: (_, value) => {
                if (value && value.isAfter(dayjs())) {
                  return Promise.reject(
                    new Error("Purchase date cannot be in the future")
                  );
                }
                return Promise.resolve();
              },
            },
          ]}
        >
          <DatePicker
            placeholder="Select purchase date (optional)"
            size="large"
            style={{ width: "100%" }}
            format="YYYY-MM-DD"
            disabledDate={(current) =>
              current && current > dayjs().endOf("day")
            }
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
              {isUpdate ? "Update Vehicle" : "Create Vehicle"}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default VehicleModal;
