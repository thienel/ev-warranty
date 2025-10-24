import React, { useEffect, useState } from "react";
import { API_ENDPOINTS, ROLE_LABELS } from "@constants/common-constants";
import { type User, type Office } from "@/types/index";
import { officesApi } from "@services/index";
import UserModal from "@components/UserManagement/UserModal/UserModal";
import useManagement from "@/hooks/useManagement";
import GenericActionBar from "@components/common/GenericActionBar/GenericActionBar";
import GenericTable from "@components/common/GenericTable/GenericTable";
import GenerateColumns from "./userTableColumns";
import useHandleApiError from "@/hooks/useHandleApiError";

const UserManagement: React.FC = () => {
  const {
    items: users,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updateUser,
    isUpdate,
    isOpenModal,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.USERS);

  const [offices, setOffices] = useState<Office[]>([]);
  const handleError = useHandleApiError();

  const fetchOffices = async (): Promise<void> => {
    try {
      const response = await officesApi.getAll();
      // Ensure we always have an array
      const officesData = response.data;
      if (Array.isArray(officesData)) {
        setOffices(officesData);
      } else {
        console.warn("API returned non-array data for offices:", officesData);
        setOffices([]);
      }
    } catch (error) {
      console.error("Failed to fetch offices:", error);
      handleError(error as Error);
      setOffices([]); // Set empty array on error
    }
  };

  useEffect(() => {
    console.log("UserManagement: Fetching offices...");
    fetchOffices();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Debug log to track offices state changes
  useEffect(() => {
    console.log("UserManagement: offices state changed:", {
      isArray: Array.isArray(offices),
      length: Array.isArray(offices) ? offices.length : "N/A",
      data: offices,
    });
  }, [offices]);

  const getOfficeName = (officeId: string): string => {
    // Safety check: ensure offices is an array before calling find
    if (!Array.isArray(offices)) {
      console.warn("offices is not an array:", offices);
      return "N/A";
    }

    const office = offices.find((o) => o.id === officeId);
    return office ? office.office_name : "N/A";
  };

  const searchFields = ["name", "email"];
  const searchFieldsWithRole = [
    ...searchFields,
    (user: Record<string, unknown> & { id: string | number }) => {
      const userRecord = user as unknown as User;
      return (
        ROLE_LABELS[userRecord.role as keyof typeof ROLE_LABELS] ||
        userRecord.role
      );
    },
  ];

  return (
    <>
      <GenericActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleOpenModal}
        loading={loading}
        searchPlaceholder="Search by name, email or role..."
        addButtonText="Add User"
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={users as (Record<string, unknown> & { id: string | number })[]}
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFieldsWithRole}
        deleteEndpoint={API_ENDPOINTS.USERS}
        deleteSuccessMessage="User deleted successfully"
        additionalProps={{ getOfficeName }}
      />

      <UserModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        user={updateUser ? (updateUser as unknown as User) : null}
        opened={isOpenModal}
        offices={offices}
        isUpdate={isUpdate}
      />
    </>
  );
};

export default UserManagement;
