import { useParams, useNavigate } from "react-router-dom";
import { usePayment } from "./usePayment";

export default function PaymentPage() {
  const { salesOrderCode } = useParams();
  const navigate = useNavigate();
  const { order, details, loading, error, payLoading, isPending, handlePay } =
    usePayment(salesOrderCode);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="h-8 w-8 border-4 border-gray-200 border-t-blue-600 rounded-full animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4">
        <div className="bg-white rounded-2xl border border-gray-100 p-6 max-w-sm w-full text-center">
          <p className="text-sm text-red-500 mb-4">{error}</p>
          <button
            onClick={() => navigate("/pesanan")}
            className="text-sm text-blue-600"
          >
            Lihat Pesanan Saya
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-lg mx-auto space-y-4">

        {/* Header */}
        <div>
          <p className="text-xs text-gray-400 mb-0.5">Kode Pesanan</p>
          <h1 className="text-xl font-semibold text-gray-900">{salesOrderCode}</h1>
        </div>

        {/* Item list */}
        <div className="bg-white rounded-2xl border border-gray-100 divide-y divide-gray-50">
          {details.map((item) => (
            <div key={item.detail_id} className="flex items-start gap-3 p-4">
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-900 truncate">
                  {item.product_name}
                </p>
                <p className="text-xs text-gray-400">{item.product_code}</p>
              </div>
              <div className="text-right shrink-0">
                <p className="text-sm text-gray-700">
                  {item.quantity} × {item.final_unit_price_format}
                </p>
                {parseFloat(item.discount_amount) > 0 && (
                  <p className="text-xs text-red-400 line-through">
                    {item.price_before_discount_format}
                  </p>
                )}
                <p className="text-sm font-medium text-gray-900">
                  {item.final_total_format}
                </p>
              </div>
            </div>
          ))}
        </div>

        {/* Ringkasan biaya */}
        <div className="bg-white rounded-2xl border border-gray-100 p-4 space-y-2">
          <Row label="Subtotal" value={order.total_before_discount_format} />
          {parseFloat(order.discount_amount) > 0 && (
            <Row
              label="Diskon"
              value={`- ${order.discount_amount_format}`}
              valueClass="text-green-600"
            />
          )}
          <Row
            label={`Ongkir (${order.shipping_display_name})`}
            value={order.shipping_fee_format}
          />
          <div className="border-t border-gray-100 pt-2 mt-2">
            <Row label="Total Pembayaran" value={order.final_total_format} bold />
          </div>
        </div>

        {/* Info pengiriman */}
        <div className="bg-white rounded-2xl border border-gray-100 p-4 space-y-1">
          <p className="text-xs font-medium text-gray-400 uppercase tracking-wide mb-2">
            Pengiriman
          </p>
          <p className="text-sm text-gray-800">
            {order.shipping_display_name} · {order.shipping_etd}
          </p>
          <p className="text-sm text-gray-500">{order.customer_address_snapshot}</p>
        </div>

        {/* Tombol bayar / status */}
        {isPending ? (
          <button
            onClick={handlePay}
            disabled={payLoading}
            className="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-60 text-white text-sm font-medium py-3 rounded-2xl transition-colors"
          >
            {payLoading ? "Memuat pembayaran..." : `Bayar ${order.final_total_format}`}
          </button>
        ) : (
          <div className="bg-white rounded-2xl border border-gray-100 p-4 text-center">
            <p className="text-sm text-gray-500">
              Status pesanan:{" "}
              <span className="font-medium text-gray-800">{order.status}</span>
            </p>
            <button
              onClick={() => navigate(`/pesanan/${salesOrderCode}`)}
              className="text-sm text-blue-600 mt-2"
            >
              Lihat Detail Pesanan
            </button>
          </div>
        )}

      </div>
    </div>
  );
}

function Row({ label, value, bold, valueClass }) {
  return (
    <div className="flex justify-between items-center">
      <span className={`text-sm ${bold ? "font-semibold text-gray-900" : "text-gray-500"}`}>
        {label}
      </span>
      <span
        className={`text-sm ${bold ? "font-semibold text-gray-900" : "text-gray-700"} ${valueClass ?? ""}`}
      >
        {value}
      </span>
    </div>
  );
}