import React from "react";
import AppLayout from "@components/Layout/Layout";
import VehicleModelManagement from "@components/VehicleModelManagement/VehicleModelManagement";

const VehicleModels: React.FC = () => {
  return (
    <AppLayout title="Vehicle Model Management">
      <VehicleModelManagement />
    </AppLayout>
  );
};

export default VehicleModels;
