import React from "react";
import { Navigate, Outlet, useRoutes } from "react-router-dom";
import Login from "@pages/Login/Login";
import { useSelector } from "react-redux";
import CallBack from "@pages/CallBack";
import Users from "@pages/admin/Users";
import AppLayout from "@components/Layout/Layout";
import Offices from "@pages/admin/Offices";
import Error from "@pages/Error/Error";
import useCheckRole from "@/hooks/useCheckRole";
import { USER_ROLES } from "@constants/common-constants";
import { useTokenRefresh } from "@/hooks/useTokenRefresh";

import EVMStaffClaims from "@pages/evm-staff/Claims";
import SCStaffClaims from "@pages/sc-staff/Claims";
import SCTechnicianClaims from "@pages/sc-technician/Claims";
import ClaimDetail from "@pages/claims/ClaimDetail";
import ClaimCreate from "@pages/sc-staff/ClaimCreate";

// New management pages
import Customers from "@pages/sc-staff/Customers";
import Vehicles from "@pages/evm-staff/Vehicles";
import VehicleModels from "@pages/evm-staff/VehicleModels";

import type { RootState } from "@redux/store";

export const ProtectedRoute: React.FC = () => {
  const authState = useSelector((state: RootState) => state.auth);
  const isAuthenticated = authState?.isAuthenticated || false;

  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
};

export const PublicRoute: React.FC = () => {
  const authState = useSelector((state: RootState) => state.auth);
  const isAuthenticated = authState?.isAuthenticated || false;

  return !isAuthenticated ? <Outlet /> : <Navigate to="/" replace />;
};

export const AdminRoute: React.FC = () => {
  const isRightRole = useCheckRole(USER_ROLES.ADMIN);

  return isRightRole ? <Outlet /> : <Navigate to="/unauthorized" replace />;
};

export const EVMStaffRoute: React.FC = () => {
  const isRightRole = useCheckRole(USER_ROLES.EVM_STAFF);

  return isRightRole ? <Outlet /> : <Navigate to="/unauthorized" replace />;
};

export const SCStaffRoute: React.FC = () => {
  const isRightRole = useCheckRole(USER_ROLES.SC_STAFF);

  return isRightRole ? <Outlet /> : <Navigate to="/unauthorized" replace />;
};

export const SCTechnicianRoute: React.FC = () => {
  const isRightRole = useCheckRole(USER_ROLES.SC_TECHNICIAN);

  return isRightRole ? <Outlet /> : <Navigate to="/unauthorized" replace />;
};

const App: React.FC = () => {
  // Initialize token refresh management
  useTokenRefresh();

  const routes = [
    {
      element: <ProtectedRoute />,
      children: [
        { path: "/", element: <AppLayout /> },
        {
          path: "/admin",
          element: <AdminRoute />,
          children: [
            {
              path: "users",
              element: <Users />,
            },
            {
              path: "offices",
              element: <Offices />,
            },
          ],
        },
        {
          path: "/evm-staff",
          element: <EVMStaffRoute />,
          children: [
            {
              path: "claims",
              children: [
                { path: "", element: <EVMStaffClaims /> },
                { path: ":id", element: <ClaimDetail /> },
              ],
            },
            {
              path: "vehicles",
              element: <Vehicles />,
            },
            {
              path: "vehicle-models",
              element: <VehicleModels />,
            },
          ],
        },
        {
          path: "/sc-staff",
          element: <SCStaffRoute />,
          children: [
            {
              path: "claims",
              children: [
                { path: "", element: <SCStaffClaims /> },
                { path: "create", element: <ClaimCreate /> },
                { path: ":id", element: <ClaimDetail /> },
              ],
            },
            {
              path: "customers",
              element: <Customers />,
            },
          ],
        },
        {
          path: "/sc-technician",
          element: <SCTechnicianRoute />,
          children: [
            {
              path: "claims",
              children: [
                { path: "", element: <SCTechnicianClaims /> },
                { path: ":id", element: <ClaimDetail /> },
              ],
            },
          ],
        },
      ],
    },
    {
      element: <PublicRoute />,
      children: [
        { path: "/login", element: <Login /> },
        { path: "/callback", element: <CallBack /> },
      ],
    },
    {
      path: "/unauthorized",
      element: <Error code={403} />,
    },
    {
      path: "/servererror",
      element: <Error code={500} />,
    },
    {
      path: "*",
      element: <Error code={404} />,
    },
  ];

  return useRoutes(routes);
};

export default App;
