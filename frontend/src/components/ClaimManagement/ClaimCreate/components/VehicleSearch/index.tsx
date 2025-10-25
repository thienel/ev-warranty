import React, { useState, useCallback, useEffect } from "react";
import { Typography, Spin, Empty, Space } from "antd";
import { CarOutlined } from "@ant-design/icons";
import { vehiclesApi, vehicleModelsApi } from "@services/index";
import type { Vehicle, Customer, VehicleModel } from "@/types";
import useHandleApiError from "@/hooks/useHandleApiError";
import "./VehicleSearch.less";

const { Text } = Typography;

interface VehicleListProps {
  onSelect: (vehicle: Vehicle | null) => void;
  selectedVehicle?: Vehicle | null;
  selectedCustomer?: Customer | null;
  disabled?: boolean;
  className?: string;
}

const VehicleSearch: React.FC<VehicleListProps> = ({
  onSelect,
  selectedVehicle,
  selectedCustomer,
  disabled = false,
  className,
}) => {
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [loading, setLoading] = useState(false);
  const [vehicleModels, setVehicleModels] = useState<VehicleModel[]>([]);
  const handleError = useHandleApiError();

  // Load vehicle models for display purposes
  useEffect(() => {
    const loadVehicleModels = async () => {
      try {
        const response = await vehicleModelsApi.getAll();
        let modelsData = response.data;

        if (
          modelsData &&
          typeof modelsData === "object" &&
          "data" in modelsData
        ) {
          modelsData = (modelsData as { data: unknown }).data as VehicleModel[];
        }

        if (Array.isArray(modelsData)) {
          setVehicleModels(modelsData);
        }
      } catch (error) {
        console.error("Failed to load vehicle models:", error);
        setVehicleModels([]);
      }
    };

    loadVehicleModels();
  }, []);

  // Get vehicle model info by ID
  const getVehicleModelInfo = useCallback(
    (modelId: string) => {
      const model = vehicleModels.find((m) => m.id === modelId);
      return model
        ? `${model.brand} ${model.model_name} ${model.year}`
        : "Unknown Model";
    },
    [vehicleModels]
  );

  // Load customer vehicles when customer is selected
  useEffect(() => {
    const loadCustomerVehicles = async () => {
      if (!selectedCustomer) {
        setVehicles([]);
        return;
      }

      try {
        setLoading(true);
        const response = await vehiclesApi.getAll({
          customerId: selectedCustomer.id,
        });

        let vehiclesData = response.data;

        // Handle nested data structure
        if (
          vehiclesData &&
          typeof vehiclesData === "object" &&
          "data" in vehiclesData
        ) {
          vehiclesData = (vehiclesData as { data: unknown }).data as Vehicle[];
        }

        if (Array.isArray(vehiclesData)) {
          setVehicles(vehiclesData);
        } else {
          setVehicles([]);
        }
      } catch (error) {
        console.error("Failed to load customer vehicles:", error);
        handleError(error as Error);
        setVehicles([]);
      } finally {
        setLoading(false);
      }
    };

    loadCustomerVehicles();
  }, [selectedCustomer, handleError]);

  const handleVehicleSelect = (vehicle: Vehicle) => {
    onSelect(vehicle);
  };

  if (!selectedCustomer) {
    return (
      <div className={className}>
        <Empty
          description="Please select a customer first to view their vehicles"
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        />
      </div>
    );
  }

  if (loading) {
    return (
      <div
        className={className}
        style={{ textAlign: "center", padding: "20px" }}
      >
        <Spin size="large" />
        <div style={{ marginTop: "16px" }}>
          <Text type="secondary">Loading vehicles...</Text>
        </div>
      </div>
    );
  }

  if (vehicles.length === 0) {
    return (
      <div className={className}>
        <Empty
          description={`No vehicles found for ${selectedCustomer.full_name || `${selectedCustomer.first_name} ${selectedCustomer.last_name}`}`}
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        />
      </div>
    );
  }

  return (
    <div className={`vehicle-list ${className || ""}`}>
      <div style={{ marginBottom: "16px" }}>
        <Text type="secondary">
          Select a vehicle for{" "}
          {selectedCustomer.full_name ||
            `${selectedCustomer.first_name} ${selectedCustomer.last_name}`}
          :
        </Text>
      </div>

      <div className="vehicle-options">
        {vehicles.map((vehicle) => {
          const modelInfo = getVehicleModelInfo(vehicle.model_id);
          const isSelected = selectedVehicle?.id === vehicle.id;

          return (
            <div
              key={vehicle.id}
              className={`vehicle-option ${isSelected ? "selected" : ""} ${disabled ? "disabled" : ""}`}
              onClick={() => !disabled && handleVehicleSelect(vehicle)}
              style={{
                border: "2px solid",
                borderColor: isSelected ? "#697565" : "#e0e6dd",
                borderRadius: "8px",
                padding: "16px",
                marginBottom: "12px",
                cursor: disabled ? "not-allowed" : "pointer",
                backgroundColor: isSelected ? "#f5f7f3" : "#ffffff",
                transition: "all 0.3s ease",
                opacity: disabled ? 0.6 : 1,
              }}
            >
              <Space direction="vertical" size={4} style={{ width: "100%" }}>
                <Space>
                  <CarOutlined
                    style={{ color: isSelected ? "#697565" : "#8b9788" }}
                  />
                  <Text
                    strong
                    style={{ color: isSelected ? "#697565" : "#2c2d2a" }}
                  >
                    {modelInfo}
                  </Text>
                </Space>

                <div style={{ marginLeft: "20px" }}>
                  {vehicle.vin && (
                    <div>
                      <Text type="secondary" style={{ fontSize: "13px" }}>
                        VIN: {vehicle.vin}
                      </Text>
                    </div>
                  )}

                  {vehicle.license_plate && (
                    <div>
                      <Text type="secondary" style={{ fontSize: "13px" }}>
                        License Plate: {vehicle.license_plate}
                      </Text>
                    </div>
                  )}

                  {vehicle.purchase_date && (
                    <div>
                      <Text type="secondary" style={{ fontSize: "12px" }}>
                        Purchased:{" "}
                        {new Date(vehicle.purchase_date).toLocaleDateString()}
                      </Text>
                    </div>
                  )}
                </div>
              </Space>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default VehicleSearch;
