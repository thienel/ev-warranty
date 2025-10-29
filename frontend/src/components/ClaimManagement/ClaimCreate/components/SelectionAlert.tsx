import React from 'react'
import { Alert, Typography } from 'antd'
import type { Customer, Vehicle } from '@/types'

const { Text } = Typography

interface SelectionAlertProps {
  selectedCustomer?: Customer | null
  selectedVehicle?: Vehicle | null
}

const SelectionAlert: React.FC<SelectionAlertProps> = ({ selectedCustomer, selectedVehicle }) => {
  if (selectedCustomer) {
    return (
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
                  <Text type="secondary">Email: {selectedCustomer.email}</Text>
                </div>
              )}
              {selectedCustomer.phone_number && (
                <div className="detail-item">
                  <Text type="secondary">Phone: {selectedCustomer.phone_number}</Text>
                </div>
              )}
              {selectedCustomer.address && (
                <div className="detail-item">
                  <Text type="secondary">Address: {selectedCustomer.address}</Text>
                </div>
              )}
            </div>
          </div>
        }
        type="success"
        showIcon
        className="selection-alert"
      />
    )
  }

  if (selectedVehicle) {
    return (
      <Alert
        description={
          <div className="vehicle-info">
            <Text strong className="vehicle-vin">
              VIN: {selectedVehicle.vin}
            </Text>
            <div className="vehicle-details">
              {selectedVehicle.license_plate && (
                <div className="detail-item">
                  <Text type="secondary">License Plate: {selectedVehicle.license_plate}</Text>
                </div>
              )}
              {selectedVehicle.purchase_date && (
                <div className="detail-item">
                  <Text type="secondary">
                    Purchase Date: {new Date(selectedVehicle.purchase_date).toLocaleDateString()}
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
    )
  }

  return null
}

export default SelectionAlert
