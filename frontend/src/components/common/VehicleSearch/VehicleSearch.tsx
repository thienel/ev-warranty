import React, { useState, useCallback, useEffect, useRef } from "react";
import { AutoComplete, Input, Space, Typography, Spin, Empty } from "antd";
import { CarOutlined, SearchOutlined } from "@ant-design/icons";
import { vehiclesApi, vehicleModelsApi } from "@services/index";
import type { Vehicle, Customer, VehicleModel } from "@/types";
import useHandleApiError from "@/hooks/useHandleApiError";

const { Text } = Typography;

interface VehicleOption {
  value: string;
  label: React.ReactNode;
  vehicle: Vehicle;
}

interface VehicleSearchProps {
  onSelect: (vehicle: Vehicle | null) => void;
  selectedVehicle?: Vehicle | null;
  selectedCustomer?: Customer | null;
  placeholder?: string;
  allowClear?: boolean;
  disabled?: boolean;
  className?: string;
}

const VehicleSearch: React.FC<VehicleSearchProps> = ({
  onSelect,
  selectedVehicle,
  selectedCustomer,
  placeholder = "Search vehicles by VIN or license plate...",
  allowClear = true,
  disabled = false,
  className,
}) => {
  const [options, setOptions] = useState<VehicleOption[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchValue, setSearchValue] = useState("");
  const [vehicleModels, setVehicleModels] = useState<VehicleModel[]>([]);
  const handleError = useHandleApiError();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Initialize with selected vehicle value
  useEffect(() => {
    if (selectedVehicle) {
      setSearchValue(
        selectedVehicle.vin ||
          selectedVehicle.license_plate ||
          selectedVehicle.id
      );
    } else {
      setSearchValue("");
    }
  }, [selectedVehicle]);

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

  const searchVehicles = useCallback(
    async (searchText: string) => {
      if (!searchText || searchText.length < 2) {
        setOptions([]);
        return;
      }

      try {
        setLoading(true);

        // Build search parameters based on search text and selected customer
        const searchPromises = [];

        // If customer is selected, search only their vehicles
        if (selectedCustomer) {
          searchPromises.push(
            vehiclesApi.getAll({ customerId: selectedCustomer.id })
          );
        } else {
          // Search by VIN (most common for exact searches)
          searchPromises.push(vehiclesApi.getAll({ vin: searchText }));

          // Search by license plate
          searchPromises.push(vehiclesApi.getAll({ licensePlate: searchText }));
        }

        const results = await Promise.allSettled(searchPromises);

        // Combine and deduplicate results
        const allVehicles: Vehicle[] = [];
        const seenIds = new Set<string>();

        results.forEach((result) => {
          if (result.status === "fulfilled" && result.value?.data) {
            let vehiclesData = result.value.data;

            // Handle nested data structure
            if (
              vehiclesData &&
              typeof vehiclesData === "object" &&
              "data" in vehiclesData
            ) {
              vehiclesData = (vehiclesData as { data: unknown })
                .data as Vehicle[];
            }

            if (Array.isArray(vehiclesData)) {
              vehiclesData.forEach((vehicle) => {
                if (!seenIds.has(vehicle.id)) {
                  seenIds.add(vehicle.id);
                  allVehicles.push(vehicle);
                }
              });
            }
          }
        });

        // Filter vehicles if we have a search text and we're searching all vehicles
        let filteredVehicles = allVehicles;
        if (!selectedCustomer && searchText) {
          filteredVehicles = allVehicles.filter(
            (vehicle) =>
              vehicle.vin?.toLowerCase().includes(searchText.toLowerCase()) ||
              vehicle.license_plate
                ?.toLowerCase()
                .includes(searchText.toLowerCase())
          );
        }

        // Convert to AutoComplete options
        const vehicleOptions: VehicleOption[] = filteredVehicles.map(
          (vehicle) => {
            const modelInfo = getVehicleModelInfo(vehicle.model_id);
            const vehicleInfo = [
              vehicle.vin && `VIN: ${vehicle.vin}`,
              vehicle.license_plate && `License: ${vehicle.license_plate}`,
            ]
              .filter(Boolean)
              .join(" • ");

            return {
              value: vehicle.id,
              vehicle,
              label: (
                <div style={{ padding: "8px 0" }}>
                  <Space
                    direction="vertical"
                    size={2}
                    style={{ width: "100%" }}
                  >
                    <Space>
                      <CarOutlined style={{ color: "#697565" }} />
                      <Text strong>{modelInfo}</Text>
                    </Space>
                    {vehicleInfo && (
                      <Text
                        type="secondary"
                        style={{ fontSize: "12px", marginLeft: "20px" }}
                      >
                        {vehicleInfo}
                      </Text>
                    )}
                    {vehicle.purchase_date && (
                      <Text
                        type="secondary"
                        style={{ fontSize: "11px", marginLeft: "20px" }}
                      >
                        Purchased:{" "}
                        {new Date(vehicle.purchase_date).toLocaleDateString()}
                      </Text>
                    )}
                  </Space>
                </div>
              ),
            };
          }
        );

        setOptions(vehicleOptions);
      } catch (error) {
        console.error("Failed to search vehicles:", error);
        handleError(error as Error);
        setOptions([]);
      } finally {
        setLoading(false);
      }
    },
    [handleError, selectedCustomer, getVehicleModelInfo]
  );

  const handleSearch = (value: string) => {
    setSearchValue(value);
    if (
      value !==
      (selectedVehicle?.vin ||
        selectedVehicle?.license_plate ||
        selectedVehicle?.id)
    ) {
      // Clear selection if user is typing a different value
      onSelect(null);
    }

    // Clear previous timeout
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    // Set new timeout for debounced search
    timeoutRef.current = setTimeout(() => {
      searchVehicles(value);
    }, 300);
  };

  const handleSelect = (_: string, option: VehicleOption) => {
    const vehicle = option.vehicle;
    setSearchValue(vehicle.vin || vehicle.license_plate || vehicle.id);
    onSelect(vehicle);
    setOptions([]); // Clear options after selection
  };

  // Load customer vehicles when customer is selected
  useEffect(() => {
    if (selectedCustomer && !searchValue) {
      // Automatically load customer's vehicles when customer is selected
      searchVehicles("");
    } else if (!selectedCustomer) {
      // Clear options when no customer is selected
      setOptions([]);
    }
  }, [selectedCustomer, searchVehicles, searchValue]);

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  const getPlaceholder = () => {
    if (!selectedCustomer) {
      return "Please select a customer first";
    }
    return selectedCustomer
      ? `Search vehicles for ${selectedCustomer.full_name || `${selectedCustomer.first_name} ${selectedCustomer.last_name}`}...`
      : placeholder;
  };

  const getNotFoundContent = () => {
    if (loading) return <Spin size="small" />;
    if (!selectedCustomer)
      return (
        <Empty
          description="Please select a customer first"
          image={Empty.PRESENTED_IMAGE_SIMPLE}
        />
      );
    return "No vehicles found";
  };

  return (
    <AutoComplete
      className={className}
      options={options}
      onSearch={handleSearch}
      onSelect={handleSelect}
      value={searchValue}
      placeholder={getPlaceholder()}
      allowClear={allowClear}
      disabled={disabled || !selectedCustomer}
      notFoundContent={getNotFoundContent()}
      popupMatchSelectWidth={false}
      style={{ width: "100%" }}
    >
      <Input
        prefix={<SearchOutlined />}
        suffix={loading ? <Spin size="small" /> : null}
      />
    </AutoComplete>
  );
};

export default VehicleSearch;
