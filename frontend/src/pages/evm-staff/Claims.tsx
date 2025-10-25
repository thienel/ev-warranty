import React from "react";
import AppLayout from "@components/Layout/Layout.tsx";
import ClaimManagement from "@components/ClaimManagement/ClaimManagement.tsx";

const EVMStaffClaims: React.FC = () => {
  return (
    <AppLayout title="Claim Management">
      <ClaimManagement />
    </AppLayout>
  );
};

export default EVMStaffClaims;
