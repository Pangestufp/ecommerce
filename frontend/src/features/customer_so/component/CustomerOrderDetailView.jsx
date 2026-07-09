
// Komponen view detail order — dipakai oleh admin maupun user detail page.

import StatusBadge from "../../../shared/ui/Statusbadge";

// Tidak ada logic di sini, pure presentational.
export default function CustomerOrderDetailView({ detail }) {
  if (!detail) return null;

  const { sales_order: order, details, histories, actions } = detail;

  return (
    <div className="space-y-6">
      {/* Header order */}
      <div className="bg-white border border-gray-200 rounded-xl p-5 space-y-3">
        <div className="flex items-start justify-between gap-4 flex-wrap">
          <div>
            <p className="text-xs text-gray-400 font-medium uppercase tracking-wide">Kode Order</p>
            <p className="text-lg font-semibold text-gray-800 mt-0.5">{order.sales_order_code}</p>
          </div>
          <StatusBadge status={order.status} size="base" />
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2 border-t border-gray-100 text-sm">
          <InfoRow label="Pelanggan" value={order.customer_name} />
          <InfoRow label="Email" value={order.customer_email} />
          <InfoRow label="Telepon" value={order.customer_phone} />
          {order.tracking_number && (
            <InfoRow label="No. Resi" value={order.tracking_number} />
          )}
          <InfoRow label="Kurir" value={`${order.shipping_display_name} – ${order.shipping_service}`} />
          <InfoRow label="ETD" value={order.shipping_etd} />
          <InfoRow label="Total Berat" value={`${order.total_weight} gram`} />
          <InfoRow label="Ongkir" value={order.shipping_fee_format} />
        </div>

        {order.notes && (
          <div className="pt-2 border-t border-gray-100">
            <p className="text-xs text-gray-400 font-medium mb-1">Catatan</p>
            <p className="text-sm text-gray-600">{order.notes}</p>
          </div>
        )}

        {/* Alamat */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2 border-t border-gray-100">
          <div>
            <p className="text-xs text-gray-400 font-medium mb-1">Alamat Pengiriman</p>
            <p className="text-sm text-gray-600 leading-relaxed">{order.customer_address_snapshot}</p>
          </div>
          <div>
            <p className="text-xs text-gray-400 font-medium mb-1">Alamat Asal</p>
            <p className="text-sm text-gray-600 leading-relaxed">{order.origin_address_snapshot}</p>
          </div>
        </div>
      </div>

      {/* Ringkasan harga */}
      <div className="bg-white border border-gray-200 rounded-xl p-5">
        <p className="text-sm font-semibold text-gray-700 mb-3">Ringkasan Pembayaran</p>
        <div className="space-y-1.5 text-sm">
          <PriceRow label="Subtotal" value={order.total_before_discount_format} />
          <PriceRow label="Diskon" value={`- ${order.discount_amount_format}`} className="text-green-600" />
          <PriceRow label="Ongkos Kirim" value={order.shipping_fee_format} />
          <div className="pt-2 border-t border-gray-100">
            <PriceRow label="Total" value={order.final_total_format} bold />
          </div>
        </div>
      </div>

      {/* Detail produk */}
      <div className="bg-white border border-gray-200 rounded-xl overflow-hidden">
        <p className="text-sm font-semibold text-gray-700 px-5 pt-4 pb-3 border-b border-gray-100">
          Item Pesanan
        </p>
        <div className="divide-y divide-gray-100">
          {details.map((item) => (
            <div key={item.detail_id} className="flex gap-3 px-5 py-3">
              {item.image_url && (
                <img
                  src={item.image_url}
                  alt={item.product_name}
                  className="w-14 h-14 rounded-lg object-cover border border-gray-100 shrink-0"
                />
              )}
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-800 truncate">{item.product_name}</p>
                <p className="text-xs text-gray-400 mt-0.5">{item.product_code}</p>
                <div className="flex flex-wrap gap-x-4 gap-y-0.5 mt-1 text-xs text-gray-500">
                  <span>{item.quantity} pcs × {item.final_unit_price_format}</span>
                  <span>{item.total_weight} gram</span>
                </div>
                {item.discount_amount && item.discount_amount !== "0" && (
                  <p className="text-xs text-green-600 mt-0.5">Diskon: {item.discount_amount_format}</p>
                )}
              </div>
              <p className="text-sm font-semibold text-gray-800 shrink-0">{item.final_total_format}</p>
            </div>
          ))}
        </div>
      </div>

      {/* Riwayat status */}
      <div className="bg-white border border-gray-200 rounded-xl p-5">
        <p className="text-sm font-semibold text-gray-700 mb-3">Riwayat Status</p>
        <div className="relative pl-4">
          <div className="absolute left-1.5 top-2 bottom-2 w-px bg-gray-200" />
          <div className="space-y-4">
            {histories.map((h) => (
              <div key={h.history_id} className="relative">
                <div className="absolute -left-[11px] top-1 w-2.5 h-2.5 rounded-full bg-blue-500 border-2 border-white" />
                <div className="ml-3">
                  <div className="flex items-center gap-2 flex-wrap">
                    <StatusBadge status={h.status} />
                    {h.previous_status && (
                      <span className="text-xs text-gray-400">dari <StatusBadge status={h.previous_status} /></span>
                    )}
                  </div>
                  <p className="text-xs text-gray-500 mt-1">{h.note}</p>
                  <p className="text-xs text-gray-400 mt-0.5">
                    {h.created_name} · {new Date(h.created_at).toLocaleString("id-ID")}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Aksi tersedia */}
      {actions && actions.length > 0 && (
        <div className="bg-amber-50 border border-amber-200 rounded-xl p-4">
          <p className="text-xs font-semibold text-amber-700 mb-2">Aksi Tersedia</p>
          <div className="flex flex-wrap gap-2">
            {actions.map((action) => (
              <span
                key={action}
                className="px-3 py-1 bg-amber-100 text-amber-800 border border-amber-300 rounded-full text-xs font-medium"
              >
                {action}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function InfoRow({ label, value }) {
  return (
    <div>
      <p className="text-xs text-gray-400 font-medium">{label}</p>
      <p className="text-sm text-gray-700 mt-0.5">{value ?? "-"}</p>
    </div>
  );
}

function PriceRow({ label, value, bold, className = "" }) {
  return (
    <div className="flex justify-between">
      <span className={`text-gray-500 ${bold ? "font-semibold text-gray-700" : ""}`}>{label}</span>
      <span className={`${bold ? "font-bold text-gray-800" : "text-gray-700"} ${className}`}>{value}</span>
    </div>
  );
}