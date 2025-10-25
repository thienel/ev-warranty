import React, { useState, useCallback, useEffect, useRef } from "react";
import { AutoComplete, Input, Space, Typography, Spin } from "antd";
import { UserOutlined, SearchOutlined } from "@ant-design/icons";
import { customersApi } from "@services/index";
import type { Customer } from "@/types";
import useHandleApiError from "@/hooks/useHandleApiError";

const { Text } = Typography;

interface CustomerOption {
  value: string;
  label: React.ReactNode;
  customer: Customer;
}

interface CustomerSearchProps {
  onSelect: (customer: Customer | null) => void;
  selectedCustomer?: Customer | null;
  placeholder?: string;
  allowClear?: boolean;
  disabled?: boolean;
  className?: string;
}

const CustomerSearch: React.FC<CustomerSearchProps> = ({
  onSelect,
  selectedCustomer,
  placeholder = "Search customers by name, email, or phone...",
  allowClear = true,
  disabled = false,
  className,
}) => {
  const [options, setOptions] = useState<CustomerOption[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchValue, setSearchValue] = useState("");
  const handleError = useHandleApiError();
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Initialize with selected customer value
  useEffect(() => {
    if (selectedCustomer) {
      setSearchValue(
        selectedCustomer.full_name ||
          `${selectedCustomer.first_name} ${selectedCustomer.last_name}`
      );
    } else {
      setSearchValue("");
    }
  }, [selectedCustomer]);

  const searchCustomers = useCallback(
    async (searchText: string) => {
      if (!searchText || searchText.length < 2) {
        setOptions([]);
        return;
      }

      try {
        setLoading(true);

        // Try different search approaches
        const searchPromises = [];

        // Search by name (most common)
        searchPromises.push(customersApi.getAll({ name: searchText }));

        // If it looks like an email, search by email
        if (searchText.includes("@")) {
          searchPromises.push(customersApi.getAll({ email: searchText }));
        }

        // If it looks like a phone number, search by phone
        if (/^\d+/.test(searchText)) {
          searchPromises.push(customersApi.getAll({ phone: searchText }));
        }

        const results = await Promise.allSettled(searchPromises);

        // Combine and deduplicate results
        const allCustomers: Customer[] = [];
        const seenIds = new Set<string>();

        results.forEach((result) => {
          if (result.status === "fulfilled" && result.value?.data) {
            let customersData = result.value.data;

            // Handle nested data structure
            if (
              customersData &&
              typeof customersData === "object" &&
              "data" in customersData
            ) {
              customersData = (customersData as { data: unknown })
                .data as Customer[];
            }

            if (Array.isArray(customersData)) {
              customersData.forEach((customer) => {
                if (!seenIds.has(customer.id)) {
                  seenIds.add(customer.id);
                  allCustomers.push(customer);
                }
              });
            }
          }
        });

        // Convert to AutoComplete options
        const customerOptions: CustomerOption[] = allCustomers.map(
          (customer) => {
            const displayName =
              customer.full_name ||
              `${customer.first_name} ${customer.last_name}`;
            const displayInfo = [
              customer.email,
              customer.phone_number,
              customer.address,
            ]
              .filter(Boolean)
              .join(" • ");

            return {
              value: customer.id,
              customer,
              label: (
                <div style={{ padding: "8px 0" }}>
                  <Space
                    direction="vertical"
                    size={2}
                    style={{ width: "100%" }}
                  >
                    <Space>
                      <UserOutlined style={{ color: "#697565" }} />
                      <Text strong>{displayName}</Text>
                    </Space>
                    {displayInfo && (
                      <Text
                        type="secondary"
                        style={{ fontSize: "12px", marginLeft: "20px" }}
                      >
                        {displayInfo}
                      </Text>
                    )}
                  </Space>
                </div>
              ),
            };
          }
        );

        setOptions(customerOptions);
      } catch (error) {
        console.error("Failed to search customers:", error);
        handleError(error as Error);
        setOptions([]);
      } finally {
        setLoading(false);
      }
    },
    [handleError]
  );

  const handleSearch = (value: string) => {
    setSearchValue(value);
    if (
      value !==
      (selectedCustomer?.full_name ||
        `${selectedCustomer?.first_name} ${selectedCustomer?.last_name}`)
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
      searchCustomers(value);
    }, 300);
  };

  const handleSelect = (_: string, option: CustomerOption) => {
    const customer = option.customer;
    setSearchValue(
      customer.full_name || `${customer.first_name} ${customer.last_name}`
    );
    onSelect(customer);
    setOptions([]); // Clear options after selection
  };

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <AutoComplete
      className={className}
      options={options}
      onSearch={handleSearch}
      onSelect={handleSelect}
      value={searchValue}
      placeholder={placeholder}
      allowClear={allowClear}
      disabled={disabled}
      notFoundContent={loading ? <Spin size="small" /> : "No customers found"}
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

export default CustomerSearch;
