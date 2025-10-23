import React from "react";
import { API_ENDPOINTS } from "@constants/common-constants.js";
import { type Office } from "@/types/index.js";
import OfficeModal from "@components/OfficeManagement/OfficeModal/OfficeModal.tsx";
import useManagement from "@/hooks/useManagement.js";
import GenericActionBar from "@components/common/GenericActionBar/GenericActionBar.tsx";
import GenericTable from "@components/common/GenericTable/GenericTable.tsx";
import GenerateColumns from "./officeTableColumns.tsx";

const OfficeManagement: React.FC = () => {
  const {
    items: offices,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updateOffice,
    isUpdate,
    isOpenModal,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.OFFICE);

  const searchFields = ["office_name", "office_type", "address"];

  return (
    <>
      <GenericActionBar
        searchText={searchText}
        setSearchText={setSearchText}
        onReset={handleReset}
        onOpenModal={handleOpenModal}
        loading={loading}
        searchPlaceholder="Search by office name, type or address..."
        addButtonText="Add Office"
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={offices as (Record<string, unknown> & { id: string | number })[]}
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFields}
        deleteEndpoint={API_ENDPOINTS.OFFICE}
        deleteSuccessMessage="Office deleted successfully"
      />

      <OfficeModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        office={updateOffice ? (updateOffice as unknown as Office) : null}
        opened={isOpenModal}
        isUpdate={isUpdate}
      />
    </>
  );
};

export default OfficeManagement;
