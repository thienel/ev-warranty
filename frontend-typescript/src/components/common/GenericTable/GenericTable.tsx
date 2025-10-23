import { useEffect, useMemo, useState } from "react";
import { message, Table } from "antd";
import type { TableProps } from "antd";
import api from "@services/api";
import useHandleApiError from "@/hooks/useHandleApiError";
import type { ErrorResponse } from "@/constants/error-messages";

type ColumnType<T> = NonNullable<TableProps<T>["columns"]>[number];

interface GenericTableProps<T = Record<string, unknown>> {
  loading: boolean;
  setLoading: (loading: boolean) => void;
  searchText: string;
  data: T[];
  onOpenModal: (item?: T | null, isUpdate?: boolean) => void;
  onRefresh: () => Promise<void>;
  generateColumns: (
    sortedInfo: Record<string, unknown>,
    filteredInfo: Record<string, unknown>,
    onOpenModal: (item?: T | null, isUpdate?: boolean) => void,
    handleDelete: (itemId: string | number) => Promise<void>,
    additionalProps: Record<string, unknown>
  ) => ColumnType<T>[];
  searchFields?: Array<string | ((item: T) => string)>;
  deleteEndpoint: string;
  deleteSuccessMessage?: string;
  additionalProps?: Record<string, unknown>;
}

interface PaginationState {
  current: number;
  pageSize: number;
  total: number;
}

const GenericTable = <
  T extends Record<string, unknown> & { id: string | number },
>({
  loading,
  setLoading,
  searchText,
  data,
  onOpenModal,
  onRefresh,
  generateColumns,
  searchFields = [],
  deleteEndpoint,
  deleteSuccessMessage = "Item deleted successfully",
  additionalProps = {},
}: GenericTableProps<T>) => {
  const [filteredInfo, setFilteredInfo] = useState<Record<string, unknown>>({});
  const [sortedInfo, setSortedInfo] = useState<Record<string, unknown>>({});
  const [pagination, setPagination] = useState<PaginationState>({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const handleError = useHandleApiError();

  useEffect(() => {
    setFilteredInfo({});
    setSortedInfo({});
    setPagination((prev) => ({
      ...prev,
      total: data.length,
    }));
  }, [data]);

  const filteredData = useMemo(() => {
    if (!searchText) return data;

    const searchLower = searchText.trim().toLowerCase();
    return data.filter((item) =>
      searchFields.some((field) => {
        let value: string;
        if (typeof field === "function") {
          value = field(item);
        } else if (field.includes(".")) {
          const nestedValue = field
            .split(".")
            .reduce(
              (obj: Record<string, unknown>, key) =>
                (obj?.[key] as Record<string, unknown>) || {},
              item
            );
          value = String(nestedValue);
        } else {
          value = String(item[field]);
        }
        return value.toLowerCase().includes(searchLower);
      })
    );
  }, [data, searchText, searchFields]);

  const handleTableChange = (
    newPagination: unknown,
    filters: unknown,
    sorter: unknown
  ) => {
    if (newPagination && typeof newPagination === "object") {
      setPagination(newPagination as PaginationState);
    }
    setFilteredInfo((filters as Record<string, unknown>) || {});
    setSortedInfo((sorter as Record<string, unknown>) || {});
  };

  const handleDelete = async (itemId: string | number) => {
    setLoading(true);
    try {
      await api.delete(`${deleteEndpoint}/${itemId}`);
      message.success(deleteSuccessMessage);
      await onRefresh();
    } catch (error) {
      handleError(error as ErrorResponse);
    } finally {
      setLoading(false);
    }
  };

  const columns = generateColumns(
    sortedInfo,
    filteredInfo,
    onOpenModal,
    handleDelete,
    additionalProps
  );

  return (
    <Table<T>
      size="middle"
      columns={columns}
      rowKey="id"
      loading={loading}
      dataSource={filteredData}
      showSorterTooltip={false}
      pagination={{
        ...pagination,
        total: filteredData.length,
        showSizeChanger: true,
        showQuickJumper: true,
        pageSizeOptions: ["10", "20", "50", "100"],
        size: "default",
      }}
      onChange={handleTableChange}
      scroll={{ x: 1000 }}
      bordered
    />
  );
};

export default GenericTable;
