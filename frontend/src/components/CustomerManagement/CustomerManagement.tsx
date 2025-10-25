import React from "react";
import { API_ENDPOINTS } from "@constants/common-constants";
import { type Customer } from "@/types/index";
import CustomerModal from "@components/CustomerManagement/CustomerModal/CustomerModal";
import useManagement from "@/hooks/useManagement";
import GenericActionBar from "@components/common/GenericActionBar/GenericActionBar";
import GenericTable from "@components/common/GenericTable/GenericTable";
import GenerateColumns from "./customerTableColumns";

const CustomerManagement: React.FC = () => {
  const {
    items: customers,
    loading,
    setLoading,
    searchText,
    setSearchText,
    updateItem: updateCustomer,
    isUpdate,
    isOpenModal,
    handleOpenModal,
    handleReset,
  } = useManagement(API_ENDPOINTS.CUSTOMERS);

  const searchFields = ["first_name", "last_name", "email", "phone_number"];
  const searchFieldsWithFullName = [
    ...searchFields,
    (customer: Record<string, unknown> & { id: string | number }) => {
      const customerRecord = customer as unknown as Customer;
      return `${customerRecord.first_name || ""} ${customerRecord.last_name || ""}`.trim();
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
        searchPlaceholder="Search by name, email or phone..."
        addButtonText="Add Customer"
      />

      <GenericTable
        loading={loading}
        setLoading={setLoading}
        searchText={searchText}
        data={
          customers as (Record<string, unknown> & { id: string | number })[]
        }
        onOpenModal={handleOpenModal}
        onRefresh={handleReset}
        generateColumns={GenerateColumns}
        searchFields={searchFieldsWithFullName}
        deleteEndpoint={API_ENDPOINTS.CUSTOMERS}
        deleteSuccessMessage="Customer deleted successfully"
      />

      <CustomerModal
        loading={loading}
        setLoading={setLoading}
        onClose={handleReset}
        customer={
          updateCustomer ? (updateCustomer as unknown as Customer) : null
        }
        opened={isOpenModal}
        isUpdate={isUpdate}
      />
    </>
  );
};

export default CustomerManagement;
