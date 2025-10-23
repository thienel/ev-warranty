import React, { useEffect, useState } from "react";
import { API_ENDPOINTS, ROLE_LABELS } from "@constants/common-constants.js";
import { type User, type Office } from "@/types/index.js";
import api from "@services/api.js";
import UserModal from "@components/UserManagement/UserModal/UserModal.tsx";
import useManagement from "@/hooks/useManagement.js";
import GenericActionBar from "@components/common/GenericActionBar/GenericActionBar.tsx";
import GenericTable from "@components/common/GenericTable/GenericTable.tsx";
import GenerateColumns from "./userTableColumns.tsx";
import useHandleApiError from "@/hooks/useHandleApiError.js";

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
  } = useManagement(API_ENDPOINTS.USER);

  const [offices, setOffices] = useState<Office[]>([]);
  const handleError = useHandleApiError();

  const fetchOffices = async (): Promise<void> => {
    try {
      const response = await api.get(API_ENDPOINTS.OFFICE);
      setOffices(response.data.data || []);
    } catch (error) {
      handleError(error as Error);
    }
  };

  useEffect(() => {
    fetchOffices();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const getOfficeName = (officeId: string): string => {
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
        deleteEndpoint={API_ENDPOINTS.USER}
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
