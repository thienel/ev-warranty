import React, { useState, useEffect } from 'react'
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
  InputNumber,
  Select,
  Spin,
} from 'antd'
import {
  SaveOutlined,
  CheckCircleOutlined,
  UserOutlined,
  CarOutlined,
  FileTextOutlined,
  DashboardOutlined,
  TeamOutlined,
} from '@ant-design/icons'
import { useNavigate, useLocation } from 'react-router-dom'
import CustomerSearch from './components/CustomerSearch'
import VehicleSearch from './components/VehicleSearch'
import ClaimFormHeader from './components/ClaimFormHeader'
import StepCard from './components/StepCard'
import SelectionAlert from './components/SelectionAlert'
import { useClaimFormState } from './hooks/useClaimFormState'
import { useClaimSubmission } from './hooks/useClaimSubmission'
import { getClaimsBasePath } from '@/utils/navigationHelpers'
import { getClaimSteps, isStepValid, getStepStatus } from './utils/claimSteps'
import { techniciansApi } from '@/services'
import type { User } from '@/types'
import './ClaimCreate.less'

const { Text } = Typography
const { TextArea } = Input
const { Option } = Select

interface ClaimFormValues {
  description: string
  kilometers: number
  technician_id: string
}

const ClaimCreate: React.FC = () => {
  const [form] = Form.useForm<ClaimFormValues>()
  const navigate = useNavigate()
  const location = useLocation()

  // State for available technicians
  const [technicians, setTechnicians] = useState<User[]>([])
  const [loadingTechnicians, setLoadingTechnicians] = useState(false)

  // Custom hooks for state management
  const {
    selectedCustomer,
    selectedVehicle,
    currentStep,
    handleCustomerSelect,
    handleVehicleSelect,
    handleStepChange,
  } = useClaimFormState()

  const { loading, submitClaim } = useClaimSubmission()

  // Get steps configuration
  const steps = getClaimSteps()

  // Fetch available technicians
  useEffect(() => {
    const fetchTechnicians = async () => {
      try {
        setLoadingTechnicians(true)
        const response = await techniciansApi.getAvailable()
        let techData = response.data

        // Handle nested data structure
        if (techData && typeof techData === 'object' && 'data' in techData) {
          techData = (techData as { data: unknown }).data as User[]
        }

        setTechnicians(Array.isArray(techData) ? techData : [])
      } catch (error) {
        console.error('Failed to fetch technicians:', error)
        setTechnicians([])
      } finally {
        setLoadingTechnicians(false)
      }
    }

    fetchTechnicians()
  }, [])

  // Handle form submission
  const handleSubmit = async (values: ClaimFormValues) => {
    const result = await submitClaim(values, selectedCustomer, selectedVehicle)
    if (!result.success && result.shouldNavigateToStep !== undefined) {
      handleStepChange(result.shouldNavigateToStep)
    }
  }

  // Handle going back
  const handleBack = () => {
    const basePath = getClaimsBasePath(location.pathname)
    navigate(basePath)
  }

  return (
    <div className="claim-create">
      {/* Header Section */}
      <ClaimFormHeader onBack={handleBack} currentStep={currentStep} totalSteps={steps.length} />

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
        <Form form={form} layout="vertical" onFinish={handleSubmit} disabled={loading}>
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
                      message: 'Please provide a description of the issue',
                    },
                    {
                      min: 10,
                      message: 'Description must be at least 10 characters long',
                    },
                    {
                      max: 1000,
                      message: 'Description cannot exceed 1000 characters',
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

                <Row gutter={[16, 16]}>
                  {/* Kilometers Field */}
                  <Col xs={24} md={12}>
                    <Form.Item
                      name="kilometers"
                      label={
                        <div className="form-label">
                          <DashboardOutlined style={{ marginRight: 4 }} />
                          <span>Vehicle Kilometers</span>
                        </div>
                      }
                      rules={[
                        {
                          required: true,
                          message: 'Please enter vehicle kilometers',
                        },
                        {
                          type: 'number',
                          min: 0,
                          message: 'Kilometers must be a positive number',
                        },
                      ]}
                    >
                      <InputNumber
                        style={{ width: '100%' }}
                        placeholder="Enter current kilometers"
                        disabled={!selectedCustomer || !selectedVehicle}
                        min={0}
                        formatter={(value) => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
                      />
                    </Form.Item>
                  </Col>

                  {/* Technician Selection */}
                  <Col xs={24} md={12}>
                    <Form.Item
                      name="technician_id"
                      label={
                        <div className="form-label">
                          <TeamOutlined style={{ marginRight: 4 }} />
                          <span>Assign Technician</span>
                        </div>
                      }
                      rules={[
                        {
                          required: true,
                          message: 'Please select a technician',
                        },
                      ]}
                    >
                      <Select
                        placeholder="Select a technician"
                        disabled={!selectedCustomer || !selectedVehicle}
                        loading={loadingTechnicians}
                        showSearch
                        optionFilterProp="children"
                        filterOption={(input, option) => {
                          const label = (option?.label as string) || ''
                          return label.toLowerCase().includes(input.toLowerCase())
                        }}
                        notFoundContent={
                          loadingTechnicians ? <Spin size="small" /> : 'No technicians available'
                        }
                      >
                        {technicians.map((tech) => (
                          <Option
                            key={tech.id}
                            value={tech.id}
                            label={`${tech.name} (${tech.email})`}
                          >
                            {tech.name} ({tech.email})
                          </Option>
                        ))}
                      </Select>
                    </Form.Item>
                  </Col>
                </Row>
              </StepCard>
            </Col>
          </Row>

          <Divider />

          {/* Action Buttons */}
          <div style={{ textAlign: 'center' }}>
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
  )
}

export default ClaimCreate
