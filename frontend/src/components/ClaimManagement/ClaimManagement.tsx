import React from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { API_ENDPOINTS } from "@constants/common-constants";
import { type Claim } from "@/types/index";
import useClaimsManagement from "@/components/ClaimManagement/useClaimsManagement";
import GenericActionBar from "@components/common/GenericActionBar/GenericActionBar";
import GenericTable from "@components/common/GenericTable/GenericTable";
import GenerateColumns from "./claimTableColumns";

const ClaimManagement: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const {
    loading,
    setLoading,
    searchText,
    setSearchText,
    handleOpenModal,
    handleReset,
    allowCreate,
  } = useClaimsManagement();

  // Mock data for demonstration since APIs are not available yet
  const mockClaims: Claim[] = [
    {
      id: "claim-001-abc-def",
      status: "SUBMITTED",
      customer_id: "cust-123-xyz",
      customer_name: "John Smith",
      vehicle_id: "veh-456-abc",
      vehicle_info: "Tesla Model 3 2023",
      description: "Battery replacement issue",
      total_cost: 2500.0,
      created_at: "2024-10-20T10:30:00Z",
      updated_at: "2024-10-20T10:30:00Z",
    },
    {
      id: "claim-002-def-ghi",
      status: "APPROVED",
      customer_id: "cust-124-xyz",
      customer_name: "Sarah Johnson",
      vehicle_id: "veh-457-def",
      vehicle_info: "Nissan Leaf 2022",
      description: "Motor controller malfunction",
      total_cost: 1800.0,
      created_at: "2024-10-19T14:15:00Z",
      updated_at: "2024-10-21T09:20:00Z",
    },
    {
      id: "claim-003-ghi-jkl",
      status: "PROCESSING",
      customer_id: "cust-125-xyz",
      customer_name: "Michael Davis",
      vehicle_id: "veh-458-ghi",
      vehicle_info: "BMW iX 2023",
      description: "Charging port repair",
      total_cost: 950.0,
      created_at: "2024-10-18T11:45:00Z",
      updated_at: "2024-10-22T16:30:00Z",
    },
    {
      id: "claim-004-jkl-mno",
      status: "APPROVED",
      customer_id: "cust-126-xyz",
      customer_name: "Emily Wilson",
      vehicle_id: "veh-459-jkl",
      vehicle_info: "Audi e-tron 2022",
      description: "Software update and diagnostics",
      total_cost: 350.0,
      created_at: "2024-10-15T08:20:00Z",
      updated_at: "2024-10-20T12:00:00Z",
    },
    {
      id: "claim-005-mno-pqr",
      status: "REJECTED",
      customer_id: "cust-127-xyz",
      customer_name: "Robert Brown",
      vehicle_id: "veh-460-mno",
      vehicle_info: "Ford Mustang Mach-E 2023",
      description: "Cosmetic damage claim",
      total_cost: 0.0,
      created_at: "2024-10-14T13:10:00Z",
      updated_at: "2024-10-16T10:45:00Z",
    },
    {
      id: "claim-006-pqr-stu",
      status: "SUBMITTED",
      customer_id: "cust-128-xyz",
      customer_name: "Lisa Anderson",
      vehicle_id: "veh-461-pqr",
      vehicle_info: "Hyundai Ioniq 5 2023",
      description: "Brake system inspection needed",
      total_cost: 0.0,
      created_at: "2024-10-23T09:15:00Z",
      updated_at: "2024-10-23T09:15:00Z",
    },
    {
      id: "claim-007-stu-vwx",
      status: "PROCESSING",
      customer_id: "cust-129-xyz",
      customer_name: "David Kim",
      vehicle_id: "veh-462-stu",
      vehicle_info: "Kia EV6 2022",
      description: "Infotainment system malfunction",
      total_cost: 1200.0,
      created_at: "2024-10-17T16:30:00Z",
      updated_at: "2024-10-19T11:20:00Z",
    },
    {
      id: "claim-008-vwx-yza",
      status: "COMPLETED",
      customer_id: "cust-130-xyz",
      customer_name: "Maria Garcia",
      vehicle_id: "veh-463-vwx",
      vehicle_info: "Volkswagen ID.4 2023",
      description: "Cancelled due to customer request",
      total_cost: 0.0,
      created_at: "2024-10-12T14:45:00Z",
      updated_at: "2024-10-13T10:30:00Z",
    },
  ];

  const handleViewDetails = (claim: Claim): void => {
    // Determine the base path from current location
    const currentPath = location.pathname;
    let basePath = "/claims";

    if (currentPath.includes("/evm-staff/")) {
      basePath = "/evm-staff/claims";
    } else if (currentPath.includes("/sc-staff/")) {
      basePath = "/sc-staff/claims";
    } else if (currentPath.includes("/sc-technician/")) {
      basePath = "/sc-technician/claims";
    }

    navigate(`${basePath}/${claim.id}`);
  };

  const handleCreateClaim = (): void => {
    // Determine the base path from current location
    const currentPath = location.pathname;
    let basePath = "/claims";

    if (currentPath.includes("/sc-staff/")) {
      basePath = "/sc-staff/claims";
    }

    navigate(`${basePath}/create`);
  };

  const searchFields = [
    "customer_name",
    "vehicle_info",
    "description",
    "status",
  ];

  return (
    <>
      <GenericActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleCreateClaim}
        loading={loading}
        searchPlaceholder="Search by customer name, vehicle, description, or status..."
        addButtonText="Create Claim"
        allowCreate={allowCreate}
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={
          mockClaims as unknown as (Record<string, unknown> & {
            id: string | number;
          })[]
        }
        onOpenModal={(record) =>
          record && handleOpenModal(record as Claim, false)
        }
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFields}
        deleteEndpoint={API_ENDPOINTS.CLAIMS}
        deleteSuccessMessage="Claim deleted successfully"
        additionalProps={{ onViewDetails: handleViewDetails }}
      />
    </>
  );
};

export default ClaimManagement;
