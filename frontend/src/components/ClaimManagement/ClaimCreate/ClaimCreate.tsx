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
  Tag,
  Progress,
} from "antd";
import {
  UserOutlined,
  CarOutlined,
  FileTextOutlined,
  SaveOutlined,
  ArrowLeftOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
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
    <div className="claim-create">
      {/* Header Section */}
      <div className="claim-create-header">
        <div className="header-nav">
          <Button
            icon={<ArrowLeftOutlined />}
            onClick={handleBack}
            type="text"
            size="large"
            className="back-button"
          >
            Back to Claims
          </Button>
        </div>

        <div className="header-content">
          <div className="header-text">
            <Title level={1} className="page-title">
              Create New Warranty Claim
            </Title>
            <Text className="page-description">
              Create a new warranty claim for a customer's vehicle. Follow the
              steps below to complete the process.
            </Text>
          </div>

          <div className="header-progress">
            <Progress
              type="circle"
              percent={Math.round(((currentStep + 1) / 3) * 100)}
              size={80}
              strokeColor={{
                "0%": "#697565",
                "100%": "#5a6358",
              }}
              format={() => (
                <div className="progress-content">
                  <div className="step-number">{currentStep + 1}</div>
                  <div className="step-total">of 3</div>
                </div>
              )}
            />
          </div>
        </div>
      </div>

      {/* Steps Navigation */}
      <Card className="steps-card">
        <Steps
          current={currentStep}
          items={steps.map((step, index) => ({
            ...step,
            status:
              currentStep > index
                ? "finish"
                : currentStep === index
                  ? "process"
                  : "wait",
            icon: currentStep > index ? <CheckCircleOutlined /> : step.icon,
          }))}
          onChange={handleStepChange}
          className="claim-steps"
        />
      </Card>

      {/* Form Content */}
      <div className="form-container">
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
                  <div className="card-header">
                    <div className="card-title">
                      <UserOutlined className="card-icon" />
                      <span>Customer Selection</span>
                    </div>
                    <div className="card-status">
                      {selectedCustomer ? (
                        <Tag color="success" icon={<CheckCircleOutlined />}>
                          Selected
                        </Tag>
                      ) : (
                        <Tag
                          color="warning"
                          icon={<ExclamationCircleOutlined />}
                        >
                          Required
                        </Tag>
                      )}
                    </div>
                  </div>
                }
                className={`step-card ${currentStep === 0 ? "active" : ""} ${selectedCustomer ? "completed" : ""}`}
                hoverable
              >
                <div className="card-content">
                  <CustomerSearch
                    onSelect={handleCustomerSelect}
                    selectedCustomer={selectedCustomer}
                    placeholder="Search by customer name or email..."
                  />

                  {selectedCustomer && (
                    <Alert
                      style={{ marginTop: 18 }}
                      description={
                        <div className="customer-info">
                          <Text strong className="customer-name">
                            {selectedCustomer.full_name ||
                              `${selectedCustomer.first_name} ${selectedCustomer.last_name}`}
                          </Text>
                          <div className="customer-details">
                            {selectedCustomer.email && (
                              <div className="detail-item">
                                <Text type="secondary">
                                  Email: {selectedCustomer.email}
                                </Text>
                              </div>
                            )}
                            {selectedCustomer.phone_number && (
                              <div className="detail-item">
                                <Text type="secondary">
                                  Phone: {selectedCustomer.phone_number}
                                </Text>
                              </div>
                            )}
                            {selectedCustomer.address && (
                              <div className="detail-item">
                                <Text type="secondary">
                                  Address: {selectedCustomer.address}
                                </Text>
                              </div>
                            )}
                          </div>
                        </div>
                      }
                      type="success"
                      showIcon
                      className="selection-alert"
                    />
                  )}
                </div>
              </Card>
            </Col>

            {/* Vehicle Selection */}
            <Col xs={24} lg={12}>
              <Card
                title={
                  <div className="card-header">
                    <div className="card-title">
                      <CarOutlined className="card-icon" />
                      <span>Vehicle Selection</span>
                    </div>
                    <div className="card-status">
                      {selectedVehicle ? (
                        <Tag color="success" icon={<CheckCircleOutlined />}>
                          Selected
                        </Tag>
                      ) : selectedCustomer ? (
                        <Tag
                          color="warning"
                          icon={<ExclamationCircleOutlined />}
                        >
                          Required
                        </Tag>
                      ) : (
                        <Tag color="default">Waiting</Tag>
                      )}
                    </div>
                  </div>
                }
                className={`step-card ${currentStep === 1 ? "active" : ""} ${!selectedCustomer ? "disabled" : ""} ${selectedVehicle ? "completed" : ""}`}
                hoverable={!!selectedCustomer}
              >
                <div className="card-content">
                  <VehicleSearch
                    onSelect={handleVehicleSelect}
                    selectedVehicle={selectedVehicle}
                    selectedCustomer={selectedCustomer}
                    disabled={!selectedCustomer}
                  />

                  {selectedVehicle && (
                    <Alert
                      description={
                        <div className="vehicle-info">
                          <Text strong className="vehicle-vin">
                            VIN: {selectedVehicle.vin}
                          </Text>
                          <div className="vehicle-details">
                            {selectedVehicle.license_plate && (
                              <div className="detail-item">
                                <Text type="secondary">
                                  License Plate: {selectedVehicle.license_plate}
                                </Text>
                              </div>
                            )}
                            {selectedVehicle.purchase_date && (
                              <div className="detail-item">
                                <Text type="secondary">
                                  Purchase Date:{" "}
                                  {new Date(
                                    selectedVehicle.purchase_date
                                  ).toLocaleDateString()}
                                </Text>
                              </div>
                            )}
                          </div>
                        </div>
                      }
                      type="success"
                      showIcon
                      className="selection-alert"
                    />
                  )}
                </div>
              </Card>
            </Col>

            {/* Claim Description */}
            <Col xs={24}>
              <Card
                title={
                  <div className="card-header">
                    <div className="card-title">
                      <FileTextOutlined className="card-icon" />
                      <span>Claim Description</span>
                    </div>
                    <div className="card-status">
                      {!selectedCustomer || !selectedVehicle ? (
                        <Tag color="default">Waiting</Tag>
                      ) : (
                        <Tag
                          color="warning"
                          icon={<ExclamationCircleOutlined />}
                        >
                          Required
                        </Tag>
                      )}
                    </div>
                  </div>
                }
                className={`step-card ${currentStep === 2 ? "active" : ""} ${!selectedCustomer || !selectedVehicle ? "disabled" : ""}`}
                hoverable={!!(selectedCustomer && selectedVehicle)}
              >
                <Form.Item
                  name="description"
                  label={
                    <div className="form-label">
                      <span>Claim Description</span>
                      <Text type="secondary" className="label-helper">
                        Provide detailed information about the issue
                      </Text>
                    </div>
                  }
                  rules={[
                    {
                      required: true,
                      message: "Please provide a description of the issue",
                    },
                    {
                      min: 10,
                      message:
                        "Description must be at least 10 characters long",
                    },
                    {
                      max: 1000,
                      message: "Description cannot exceed 1000 characters",
                    },
                  ]}
                >
                  <TextArea
                    style={{ marginTop: 12 }}
                    rows={8}
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
          <div style={{ textAlign: "center" }}>
            <Space size="large">
              <Button onClick={handleBack} disabled={loading} size="large">
                Cancel
              </Button>
              <Button
                type="primary"
                htmlType="submit"
                icon={<SaveOutlined />}
                loading={loading}
                disabled={!isStepValid(2)}
                size="large"
                className="submit-button"
              >
                Create Claim
              </Button>
            </Space>
          </div>
        </Form>
      </div>
    </div>
  );
};

export default ClaimCreate;
