import { useNavigate } from "react-router-dom";
import { useAdminSalesOrder } from "./useAdminSalesOrder";
import Table from "../../shared/table/Table";
import Button from "../../shared/ui/Button";
import StatusBadge from "../../shared/ui/StatusBadge";
import StatusFilter from "../../shared/ui/StatusFilter";

const columns = [
  { key: "sales_order_code", label: "Kode Order", align: "left" },
  { key: "customer_name",    label: "Pelanggan",   align: "left" },
  {
    key: "status",
    label: "Status",
    align: "left",
    render: (val) => <StatusBadge status={val} />,
  },
  { key: "final_total_format",  label: "Total",    align: "right" },
  { key: "shipping_display_name", label: "Kurir",  align: "left" },
  {
    key: "created_at",
    label: "Tanggal",
    align: "left",
    render: (val) => new Date(val).toLocaleDateString("id-ID", { day: "2-digit", month: "short", year: "numeric" }),
  },
];

export default function AdminSalesOrderPage() {
  const navigate = useNavigate();
  const {
    loading,
    orders,
    page,
    hasNext,
    hasPrev,
    selectedStatuses,
    toggleStatus,
    clearStatuses,
    next,
    prev,
  } = useAdminSalesOrder();

  return (
    <div className="p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-lg font-semibold text-gray-800">Sales Order</h1>
      </div>

      {/* Filter status */}
      <StatusFilter
        selected={selectedStatuses}
        onToggle={toggleStatus}
        onClear={clearStatuses}
      />

      {/* Table — klik row navigasi ke detail */}
      <div className="bg-white border border-gray-200 rounded-xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                {columns.map((col) => (
                  <th
                    key={col.key}
                    className={`px-4 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wide
                      ${col.align === "right" ? "text-right" : "text-left"}`}
                  >
                    {col.label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {loading && orders.length === 0 ? (
                <tr>
                  <td colSpan={columns.length} className="px-4 py-8 text-center text-sm text-gray-400">
                    Memuat...
                  </td>
                </tr>
              ) : orders.length === 0 ? (
                <tr>
                  <td colSpan={columns.length} className="px-4 py-8 text-center text-sm text-gray-400">
                    Tidak ada data
                  </td>
                </tr>
              ) : (
                orders.map((row) => (
                  <tr
                    key={row.sales_order_id}
                    onClick={() => navigate(`/admin/sales-order/${row.sales_order_code}`)}
                    className="hover:bg-gray-50 cursor-pointer transition-colors"
                  >
                    {columns.map((col) => (
                      <td
                        key={col.key}
                        className={`px-4 py-3 text-gray-700 ${col.align === "right" ? "text-right" : ""}`}
                      >
                        {col.render ? col.render(row[col.key], row) : row[col.key]}
                      </td>
                    ))}
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Pagination */}
      <div className="flex items-center justify-end gap-2">
        <Button variant="secondary" onClick={prev} disabled={loading || !hasPrev}>
          Prev
        </Button>
        <span className="text-sm text-gray-500">Page {page}</span>
        <Button variant="secondary" onClick={next} disabled={loading || !hasNext}>
          Next
        </Button>
      </div>
    </div>
  );
}