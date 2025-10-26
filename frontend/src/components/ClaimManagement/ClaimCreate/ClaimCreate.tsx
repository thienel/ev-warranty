import React from "react";
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
  Steps,
} from "antd";
import {
  SaveOutlined,
  CheckCircleOutlined,
  UserOutlined,
  CarOutlined,
  FileTextOutlined,
} from "@ant-design/icons";
import { useNavigate, useLocation } from "react-router-dom";
import CustomerSearch from "./components/CustomerSearch";
import VehicleSearch from "./components/VehicleSearch";
import ClaimFormHeader from "./components/ClaimFormHeader";
import StepCard from "./components/StepCard";
import SelectionAlert from "./components/SelectionAlert";
import { useClaimFormState } from "./hooks/useClaimFormState";
import { useClaimSubmission } from "./hooks/useClaimSubmission";
import { getClaimsBasePath } from "@/utils/navigationHelpers";
import { getClaimSteps, isStepValid, getStepStatus } from "./utils/claimSteps";
import "./ClaimCreate.less";

const { Text } = Typography;
const { TextArea } = Input;

interface ClaimFormValues {
  description: string;
}

const ClaimCreate: React.FC = () => {
  const [form] = Form.useForm<ClaimFormValues>();
  const navigate = useNavigate();
  const location = useLocation();

  // Custom hooks for state management
  const {
    selectedCustomer,
    selectedVehicle,
    currentStep,
    handleCustomerSelect,
    handleVehicleSelect,
    handleStepChange,
  } = useClaimFormState();

  const { loading, submitClaim } = useClaimSubmission();

  // Get steps configuration
  const steps = getClaimSteps();

  // Handle form submission
  const handleSubmit = async (values: ClaimFormValues) => {
    const result = await submitClaim(values, selectedCustomer, selectedVehicle);
    if (!result.success && result.shouldNavigateToStep !== undefined) {
      handleStepChange(result.shouldNavigateToStep);
    }
  };

  // Handle going back
  const handleBack = () => {
    const basePath = getClaimsBasePath(location.pathname);
    navigate(basePath);
  };

  return (
    <div className="claim-create">
      {/* Header Section */}
      <ClaimFormHeader
        onBack={handleBack}
        currentStep={currentStep}
        totalSteps={steps.length}
      />

      {/* Steps Navigation */}
      <Card className="steps-card">
        <Steps
          current={currentStep}
          items={steps.map((step, index) => ({
            ...step,
            status: getStepStatus(currentStep, index),
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
              <StepCard
                title="Customer Selection"
                icon={<UserOutlined className="card-icon" />}
                isActive={currentStep === 0}
                isCompleted={!!selectedCustomer}
              >
                <CustomerSearch
                  onSelect={handleCustomerSelect}
                  selectedCustomer={selectedCustomer}
                  placeholder="Search by customer name or email..."
                />
                <SelectionAlert selectedCustomer={selectedCustomer} />
              </StepCard>
            </Col>

            {/* Vehicle Selection */}
            <Col xs={24} lg={12}>
              <StepCard
                title="Vehicle Selection"
                icon={<CarOutlined className="card-icon" />}
                isActive={currentStep === 1}
                isCompleted={!!selectedVehicle}
                isDisabled={!selectedCustomer}
                hoverable={!!selectedCustomer}
              >
                <VehicleSearch
                  onSelect={handleVehicleSelect}
                  selectedVehicle={selectedVehicle}
                  selectedCustomer={selectedCustomer}
                  disabled={!selectedCustomer}
                />
                <SelectionAlert selectedVehicle={selectedVehicle} />
              </StepCard>
            </Col>

            {/* Claim Description */}
            <Col xs={24}>
              <StepCard
                title="Claim Description"
                icon={<FileTextOutlined className="card-icon" />}
                isActive={currentStep === 2}
                isCompleted={false}
                isDisabled={!selectedCustomer || !selectedVehicle}
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
                    autoSize={{ minRows: 6, maxRows: 12 }}
                    style={{ marginTop: 12 }}
                    rows={8}
                    placeholder="Describe the issue, symptoms, and any relevant details about the warranty claim..."
                    disabled={!selectedCustomer || !selectedVehicle}
                    maxLength={1000}
                  />
                </Form.Item>
              </StepCard>
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
                disabled={!isStepValid(2, selectedCustomer, selectedVehicle)}
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
