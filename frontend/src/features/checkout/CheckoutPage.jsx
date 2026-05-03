import { useCheckout } from "./useCheckout";

function formatRupiah(n) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency", currency: "IDR", maximumFractionDigits: 0,
  }).format(parseFloat(n));
}

export default function CheckoutPage() {
  const { checkoutData } = useCheckout();

  if (!checkoutData) return null;

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-xl mx-auto px-4 py-6">

        <h1 className="text-lg font-semibold text-gray-900 mb-6">
          Konfirmasi Pesanan
        </h1>

        {/* Note kalau ada perubahan dari backend */}
        {checkoutData.is_note === 1 && (
          <div className="bg-orange-50 border border-orange-100 rounded-xl px-4 py-3 mb-4">
            <p className="text-xs text-orange-700">{checkoutData.note}</p>
          </div>
        )}

        {/* List produk */}
        <div className="bg-white rounded-2xl border border-gray-100 px-4 mb-3">
          {checkoutData.list_product?.map(item => (
            <div
              key={item.product_id}
              className="flex items-center gap-3 py-4 border-b border-gray-100 last:border-0"
            >
              <div className="shrink-0 w-14 h-14 rounded-xl overflow-hidden bg-gray-50">
                {item.image ? (
                  <img
                    src={item.image}
                    alt={item.product_name}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-gray-300 text-xs">
                    No Img
                  </div>
                )}
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-900 truncate">
                  {item.product_name}
                </p>
                <p className="text-xs text-gray-400 mt-0.5">x{item.qty}</p>
              </div>
              <p className="text-sm font-medium text-gray-900 shrink-0">
                {item.Price_format}
              </p>
            </div>
          ))}
        </div>

        {/* Total */}
        <div className="bg-white rounded-2xl border border-gray-100 px-4 py-4">
          <div className="flex justify-between items-center">
            <span className="text-sm text-gray-500">Total</span>
            <span className="text-base font-semibold text-gray-900">
              {formatRupiah(checkoutData.total_now)}
            </span>
          </div>
        </div>

      </div>
    </div>
  );
}