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
import { USER_ROLES } from "@constants/common-constants.js";

import EVMStaffClaims from "@pages/evm-staff/Claims";
import SCStaffClaims from "@pages/sc-staff/Claims";
import SCTechnicianClaims from "@pages/sc-technician/Claims";
import ClaimDetail from "@pages/claims/ClaimDetail";
import type { RootState } from "@redux/store";

export const ProtectedRoute: React.FC = () => {
  const { isAuthenticated } = useSelector((state: RootState) => state.auth);

  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
};

export const PublicRoute: React.FC = () => {
  const { isAuthenticated } = useSelector((state: RootState) => state.auth);

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
                { path: ":id", element: <ClaimDetail /> },
              ],
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
