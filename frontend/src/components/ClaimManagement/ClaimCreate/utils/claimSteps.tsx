import React from 'react'
import { UserOutlined, CarOutlined, FileTextOutlined } from '@ant-design/icons'
import type { Customer, Vehicle } from '@/types'

export interface ClaimStep {
  title: string
  icon: React.ReactNode
  description: string
}

export const getClaimSteps = (): ClaimStep[] => [
  {
    title: 'Select Customer',
    icon: <UserOutlined />,
    description: 'Search and select the customer',
  },
  {
    title: 'Select Vehicle',
    icon: <CarOutlined />,
    description: 'Choose the vehicle for warranty claim',
  },
  {
    title: 'Claim Details',
    icon: <FileTextOutlined />,
    description: 'Provide claim description',
  },
]

export const isStepValid = (
  step: number,
  selectedCustomer: Customer | null,
  selectedVehicle: Vehicle | null,
): boolean => {
  switch (step) {
    case 0:
      return !!selectedCustomer
    case 1:
      return !!selectedCustomer && !!selectedVehicle
    case 2:
      return !!selectedCustomer && !!selectedVehicle
    default:
      return false
  }
}

export const getStepStatus = (
  currentStep: number,
  stepIndex: number,
): 'finish' | 'process' | 'wait' => {
  if (currentStep > stepIndex) return 'finish'
  if (currentStep === stepIndex) return 'process'
  return 'wait'
}

export const canNavigateToStep = (
  targetStep: number,
  selectedCustomer: Customer | null,
  selectedVehicle: Vehicle | null,
): boolean => {
  if (targetStep === 0) return true
  if (targetStep === 1) return !!selectedCustomer
  if (targetStep === 2) return !!selectedCustomer && !!selectedVehicle
  return false
}
