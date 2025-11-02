import React from 'react'
import { Card, Descriptions, Typography } from 'antd'
import { CarOutlined } from '@ant-design/icons'
import type { VehicleDetail } from '@/types/index'

const { Title, Text } = Typography

interface VehicleInfoProps {
  vehicle: VehicleDetail | null
  loading: boolean
}

const VehicleInfo: React.FC<VehicleInfoProps> = ({ vehicle, loading }) => {
  return (
    <Card
      title={
        <Title level={4}>
          <CarOutlined /> Vehicle Information
        </Title>
      }
      loading={loading}
    >
      {vehicle ? (
        <Descriptions bordered column={1}>
          <Descriptions.Item label="VIN">
            <Text code strong>
              {vehicle.vin}
            </Text>
          </Descriptions.Item>
          {vehicle.license_plate && (
            <Descriptions.Item label="License Plate">
              <Text code>{vehicle.license_plate}</Text>
            </Descriptions.Item>
          )}
          {vehicle.model && (
            <Descriptions.Item label="Model">{vehicle.model.model_name}</Descriptions.Item>
          )}
          {vehicle.model?.brand && (
            <Descriptions.Item label="Brand">{vehicle.model.brand}</Descriptions.Item>
          )}
          {vehicle.model?.year && (
            <Descriptions.Item label="Year">{vehicle.model.year}</Descriptions.Item>
          )}
          {vehicle.model?.policy_name && (
            <Descriptions.Item label="Warranty Policy">
              {vehicle.model.policy_name}
            </Descriptions.Item>
          )}
          {vehicle.purchase_date && (
            <Descriptions.Item label="Purchase Date">
              {new Date(vehicle.purchase_date).toLocaleDateString()}
            </Descriptions.Item>
          )}
        </Descriptions>
      ) : !loading ? (
        <Text type="secondary">Vehicle information not available</Text>
      ) : null}
    </Card>
  )
}

export default VehicleInfo
