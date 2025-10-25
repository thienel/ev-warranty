import React from "react";
import AppLayout from "@components/Layout/Layout";
import VehicleManagement from "@components/VehicleManagement/VehicleManagement";

const Vehicles: React.FC = () => {
  return (
    <AppLayout title="Vehicle Management">
      <VehicleManagement />
    </AppLayout>
  );
};

export default Vehicles;
