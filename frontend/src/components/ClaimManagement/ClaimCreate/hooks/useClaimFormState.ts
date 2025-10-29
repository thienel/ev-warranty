import { useState } from 'react'
import type { Customer, Vehicle } from '@/types'
import { canNavigateToStep } from '../utils/claimSteps'

export const useClaimFormState = () => {
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(null)
  const [selectedVehicle, setSelectedVehicle] = useState<Vehicle | null>(null)
  const [currentStep, setCurrentStep] = useState(0)

  const handleCustomerSelect = (customer: Customer | null) => {
    setSelectedCustomer(customer)
    // Reset vehicle selection when customer changes
    if (selectedVehicle && (!customer || customer.id !== selectedVehicle.customer_id)) {
      setSelectedVehicle(null)
    }
    // Auto-advance to next step if customer is selected
    if (customer && currentStep === 0) {
      setCurrentStep(1)
    }
  }

  const handleVehicleSelect = (vehicle: Vehicle | null) => {
    setSelectedVehicle(vehicle)
    // Auto-advance to next step if vehicle is selected
    if (vehicle && currentStep === 1) {
      setCurrentStep(2)
    }
  }

  const handleStepChange = (step: number) => {
    if (canNavigateToStep(step, selectedCustomer, selectedVehicle)) {
      setCurrentStep(step)
    }
  }

  return {
    selectedCustomer,
    selectedVehicle,
    currentStep,
    handleCustomerSelect,
    handleVehicleSelect,
    handleStepChange,
  }
}
