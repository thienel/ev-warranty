import React, { useState } from "react";
import { Layout } from "antd";
import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar/Sidebar.js";
import LayoutHeader from "./LayoutHeader/LayoutHeader.js";
import LayoutContent from "./LayoutContent/LayoutContent.js";
import "./Layout.less";

interface AppLayoutProps {
  children?: React.ReactNode;
  title?: string;
}

const AppLayout: React.FC<AppLayoutProps> = ({ children, title }) => {
  const [collapsed, setCollapsed] = useState(false);

  const handleToggleCollapse = () => {
    setCollapsed(!collapsed);
  };

  return (
    <Layout
      className="app-layout"
      style={{ height: "100vh", overflow: "hidden" }}
    >
      <Sidebar collapsed={collapsed} />

      <Layout style={{ height: "100vh", overflow: "hidden" }}>
        <LayoutHeader
          collapsed={collapsed}
          onToggleCollapse={handleToggleCollapse}
          title={title}
        />
        <LayoutContent>{children || <Outlet />}</LayoutContent>
      </Layout>
    </Layout>
  );
};

export default AppLayout;
