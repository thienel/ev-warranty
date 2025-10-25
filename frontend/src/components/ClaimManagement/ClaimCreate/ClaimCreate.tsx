import React, { useState } from "react";
import {
  Card,
  Form,
  Input,
  Button,
  Space,
  Typography,
  Divider,
  Row,
  Col,
  message,
  Steps,
  Alert,
} from "antd";
import {
  UserOutlined,
  CarOutlined,
  FileTextOutlined,
  SaveOutlined,
  ArrowLeftOutlined,
} from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import CustomerSearch from "@components/common/CustomerSearch/CustomerSearch";
import VehicleSearch from "@components/common/VehicleSearch/VehicleSearch";
import { claimsApi } from "@services/claimsApi";
import type { Customer, Vehicle, CreateClaimRequest } from "@/types";
import useHandleApiError from "@/hooks/useHandleApiError";
import "./ClaimCreate.less";

const { Title, Text } = Typography;
const { TextArea } = Input;

interface ClaimFormValues {
  description: string;
}

const ClaimCreate: React.FC = () => {
  const [form] = Form.useForm<ClaimFormValues>();
  const navigate = useNavigate();
  const handleError = useHandleApiError();

  // State management
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(
    null
  );
  const [selectedVehicle, setSelectedVehicle] = useState<Vehicle | null>(null);
  const [loading, setLoading] = useState(false);
  const [currentStep, setCurrentStep] = useState(0);

  // Steps configuration
  const steps = [
    {
      title: "Select Customer",
      icon: <UserOutlined />,
      description: "Search and select the customer",
    },
    {
      title: "Select Vehicle",
      icon: <CarOutlined />,
      description: "Choose the vehicle for warranty claim",
    },
    {
      title: "Claim Details",
      icon: <FileTextOutlined />,
      description: "Provide claim description",
    },
  ];

  // Handle customer selection
  const handleCustomerSelect = (customer: Customer | null) => {
    setSelectedCustomer(customer);
    // Reset vehicle selection when customer changes
    if (
      selectedVehicle &&
      (!customer || customer.id !== selectedVehicle.customer_id)
    ) {
      setSelectedVehicle(null);
    }
    // Auto-advance to next step if customer is selected
    if (customer && currentStep === 0) {
      setCurrentStep(1);
    }
  };

  // Handle vehicle selection
  const handleVehicleSelect = (vehicle: Vehicle | null) => {
    setSelectedVehicle(vehicle);
    // Auto-advance to next step if vehicle is selected
    if (vehicle && currentStep === 1) {
      setCurrentStep(2);
    }
  };

  // Handle form submission
  const handleSubmit = async (values: ClaimFormValues) => {
    if (!selectedCustomer) {
      message.error("Please select a customer");
      setCurrentStep(0);
      return;
    }

    if (!selectedVehicle) {
      message.error("Please select a vehicle");
      setCurrentStep(1);
      return;
    }

    if (!values.description?.trim()) {
      message.error("Please provide a claim description");
      return;
    }

    try {
      setLoading(true);

      const claimData: CreateClaimRequest = {
        customer_id: selectedCustomer.id,
        vehicle_id: selectedVehicle.id,
        description: values.description.trim(),
      };

      await claimsApi.create(claimData);

      message.success("Claim created successfully!");
      navigate("/sc-staff/claims");
    } catch (error) {
      console.error("Failed to create claim:", error);
      handleError(error as Error);
    } finally {
      setLoading(false);
    }
  };

  // Handle going back
  const handleBack = () => {
    navigate("/sc-staff/claims");
  };

  // Handle step navigation
  const handleStepChange = (step: number) => {
    if (step === 0) {
      setCurrentStep(0);
    } else if (step === 1 && selectedCustomer) {
      setCurrentStep(1);
    } else if (step === 2 && selectedCustomer && selectedVehicle) {
      setCurrentStep(2);
    }
  };

  // Validation for current step
  const isStepValid = (step: number) => {
    switch (step) {
      case 0:
        return !!selectedCustomer;
      case 1:
        return !!selectedCustomer && !!selectedVehicle;
      case 2:
        return !!selectedCustomer && !!selectedVehicle;
      default:
        return false;
    }
  };

  return (
    <div
      className="claim-create"
      style={{ padding: "24px", maxWidth: "1200px", margin: "0 auto" }}
    >
      {/* Header */}
      <div style={{ marginBottom: "24px" }}>
        <Space size="middle" style={{ marginBottom: "16px" }}>
          <Button icon={<ArrowLeftOutlined />} onClick={handleBack} type="text">
            Back to Claims
          </Button>
        </Space>
        <Title level={2} style={{ margin: 0 }}>
          Create New Warranty Claim
        </Title>
        <Text type="secondary">
          Create a new warranty claim for a customer's vehicle. Only SC Staff
          can create claims.
        </Text>
      </div>

      {/* Steps */}
      <Card style={{ marginBottom: "24px" }}>
        <Steps
          current={currentStep}
          items={steps}
          onChange={handleStepChange}
        />
      </Card>

      {/* Form Content */}
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        disabled={loading}
      >
        <Row gutter={[24, 24]}>
          {/* Customer Selection */}
          <Col xs={24} lg={12}>
            <Card
              title={
                <Space>
                  <UserOutlined />
                  <span>Customer Selection</span>
                </Space>
              }
              className={`step-card ${currentStep === 0 ? "active" : ""}`}
            >
              <Space
                direction="vertical"
                size="middle"
                style={{ width: "100%" }}
              >
                <CustomerSearch
                  onSelect={handleCustomerSelect}
                  selectedCustomer={selectedCustomer}
                  placeholder="Search by customer name, email, or phone..."
                  className="customer-search"
                />

                {selectedCustomer && (
                  <Alert
                    message="Customer Selected"
                    description={
                      <div className="customer-info">
                        <Text strong>
                          {selectedCustomer.full_name ||
                            `${selectedCustomer.first_name} ${selectedCustomer.last_name}`}
                        </Text>
                        <br />
                        {selectedCustomer.email && (
                          <Text type="secondary">
                            Email: {selectedCustomer.email}
                          </Text>
                        )}
                        {selectedCustomer.email &&
                          selectedCustomer.phone_number && <br />}
                        {selectedCustomer.phone_number && (
                          <Text type="secondary">
                            Phone: {selectedCustomer.phone_number}
                          </Text>
                        )}
                        {(selectedCustomer.email ||
                          selectedCustomer.phone_number) &&
                          selectedCustomer.address && <br />}
                        {selectedCustomer.address && (
                          <Text type="secondary">
                            Address: {selectedCustomer.address}
                          </Text>
                        )}
                      </div>
                    }
                    type="success"
                    showIcon
                    style={{ marginTop: "8px" }}
                  />
                )}
              </Space>
            </Card>
          </Col>

          {/* Vehicle Selection */}
          <Col xs={24} lg={12}>
            <Card
              title={
                <Space>
                  <CarOutlined />
                  <span>Vehicle Selection</span>
                </Space>
              }
              className={`step-card ${currentStep === 1 ? "active" : ""} ${!selectedCustomer ? "disabled" : ""}`}
            >
              <Space
                direction="vertical"
                size="middle"
                style={{ width: "100%" }}
              >
                <VehicleSearch
                  onSelect={handleVehicleSelect}
                  selectedVehicle={selectedVehicle}
                  selectedCustomer={selectedCustomer}
                  disabled={!selectedCustomer}
                  className="vehicle-search"
                />

                {selectedVehicle && (
                  <Alert
                    message="Vehicle Selected"
                    description={
                      <div className="vehicle-info">
                        <Text strong>VIN: {selectedVehicle.vin}</Text>
                        <br />
                        {selectedVehicle.license_plate && (
                          <>
                            <Text type="secondary">
                              License Plate: {selectedVehicle.license_plate}
                            </Text>
                            <br />
                          </>
                        )}
                        {selectedVehicle.purchase_date && (
                          <Text type="secondary">
                            Purchase Date:{" "}
                            {new Date(
                              selectedVehicle.purchase_date
                            ).toLocaleDateString()}
                          </Text>
                        )}
                      </div>
                    }
                    type="success"
                    showIcon
                    style={{ marginTop: "8px" }}
                  />
                )}
              </Space>
            </Card>
          </Col>

          {/* Claim Description */}
          <Col xs={24}>
            <Card
              title={
                <Space>
                  <FileTextOutlined />
                  <span>Claim Description</span>
                </Space>
              }
              className={`step-card ${currentStep === 2 ? "active" : ""} ${!selectedCustomer || !selectedVehicle ? "disabled" : ""}`}
            >
              <Form.Item
                name="description"
                label="Claim Description"
                rules={[
                  {
                    required: true,
                    message: "Please provide a description of the issue",
                  },
                  {
                    min: 10,
                    message: "Description must be at least 10 characters long",
                  },
                  {
                    max: 1000,
                    message: "Description cannot exceed 1000 characters",
                  },
                ]}
                extra="Describe the issue or problem that requires warranty coverage. Be as detailed as possible."
              >
                <TextArea
                  rows={6}
                  placeholder="Describe the issue, symptoms, and any relevant details about the warranty claim..."
                  disabled={!selectedCustomer || !selectedVehicle}
                  showCount
                  maxLength={1000}
                />
              </Form.Item>
            </Card>
          </Col>
        </Row>

        <Divider />

        {/* Action Buttons */}
        <div className="form-actions">
          <Space size="middle">
            <Button onClick={handleBack} disabled={loading}>
              Cancel
            </Button>
            <Button
              type="primary"
              htmlType="submit"
              icon={<SaveOutlined />}
              loading={loading}
              disabled={!isStepValid(2)}
              size="large"
            >
              Create Claim
            </Button>
          </Space>
        </div>
      </Form>
    </div>
  );
};

export default ClaimCreate;
