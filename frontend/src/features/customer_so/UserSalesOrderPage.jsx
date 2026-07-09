import { useNavigate } from "react-router-dom";
import { useUserSalesOrder } from "./useUserSalesOrder";
import Button from "../../shared/ui/Button";
import StatusFilter from "../../shared/ui/Statusfilter";
import StatusBadge from "../../shared/ui/Statusbadge";

export default function UserSalesOrderPage() {
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
  } = useUserSalesOrder();

  return (
    <div className="p-4 sm:p-6 space-y-4 max-w-2xl mx-auto">
      <h1 className="text-lg font-semibold text-gray-800">Pesanan Saya</h1>

      {/* Filter status — scrollable horizontal di mobile */}
      <div className="overflow-x-auto pb-1">
        <StatusFilter
          selected={selectedStatuses}
          onToggle={toggleStatus}
          onClear={clearStatuses}
        />
      </div>

      {/* Card list */}
      {loading && orders.length === 0 ? (
        <div className="py-12 text-center text-sm text-gray-400">Memuat...</div>
      ) : orders.length === 0 ? (
        <div className="py-12 text-center text-sm text-gray-400">Belum ada pesanan</div>
      ) : (
        <div className="space-y-3">
          {orders.map((order) => (
            <button
              key={order.sales_order_id}
              onClick={() => navigate(`/pesanan/${order.sales_order_code}`)}
              className="w-full text-left bg-white border border-gray-200 rounded-xl p-4 hover:border-gray-300 hover:shadow-sm transition-all active:scale-[0.99]"
            >
              {/* Baris atas: kode + status */}
              <div className="flex items-start justify-between gap-2">
                <p className="text-sm font-semibold text-gray-800">{order.sales_order_code}</p>
                <StatusBadge status={order.status} />
              </div>

              {/* Nama pelanggan & kurir */}
              <p className="text-xs text-gray-500 mt-1.5">
                {order.shipping_display_name} · {order.shipping_service}
              </p>

              {/* Tanggal */}
              <p className="text-xs text-gray-400 mt-0.5">
                {new Date(order.created_at).toLocaleDateString("id-ID", {
                  day: "2-digit",
                  month: "long",
                  year: "numeric",
                })}
              </p>

              {/* Divider */}
              <div className="border-t border-gray-100 mt-3 pt-3 flex items-center justify-between">
                <span className="text-xs text-gray-400">Total Pembayaran</span>
                <span className="text-sm font-bold text-gray-800">{order.final_total_format}</span>
              </div>
            </button>
          ))}
        </div>
      )}

      {/* Pagination */}
      <div className="flex items-center justify-center gap-3 pt-2">
        <Button variant="secondary" onClick={prev} disabled={loading || !hasPrev} className="flex-1 sm:flex-none">
          ← Sebelumnya
        </Button>
        <span className="text-sm text-gray-400 shrink-0">Hal. {page}</span>
        <Button variant="secondary" onClick={next} disabled={loading || !hasNext} className="flex-1 sm:flex-none">
          Berikutnya →
        </Button>
      </div>
    </div>
  );
}