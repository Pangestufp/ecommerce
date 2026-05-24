import { formatRupiah, getUnitPrice } from "../checkoutHelpers";

/**
 * OrderSummary
 * Ringkasan harga: breakdown per produk dan grand total.
 *
 * Props:
 *  - products       : array product_price dari API
 *  - productStates  : { [product_id]: { qty, selectedDiscountId } }
 *  - grandTotal     : number, total keseluruhan
 */
export default function OrderSummary({ products, productStates, grandTotal }) {
  return (
    <div className="bg-white rounded-2xl border border-gray-100 px-4 py-4 mb-3">
      <p className="text-sm font-semibold text-gray-800 mb-3">Ringkasan Harga</p>

      <div className="flex flex-col gap-2">
        {products.map((item) => {
          const state = productStates[item.product_id];
          if (!state) return null;
          const disc = item.discount?.find((d) => d.discount_id === state.selectedDiscountId) ?? null;
          const lineTotal = getUnitPrice(item, disc) * state.qty;

          return (
            <div key={item.product_id} className="flex justify-between items-center">
              <span className="text-xs text-gray-500 truncate flex-1 mr-2">
                {item.product_name}{" "}
                <span className="text-gray-400">x{state.qty}</span>
              </span>
              <span className="text-xs font-medium text-gray-800 shrink-0">
                {formatRupiah(lineTotal)}
              </span>
            </div>
          );
        })}
      </div>

      <div className="border-t border-gray-100 mt-3 pt-3 flex justify-between items-center">
        <span className="text-sm font-semibold text-gray-700">Total</span>
        <span className="text-base font-bold text-gray-900">{formatRupiah(grandTotal)}</span>
      </div>
    </div>
  );
}