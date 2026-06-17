import { formatRupiah, getUnitPrice } from "../checkoutHelpers";
import DiscountPicker from "./DiscountPicker";

/**
 * ProductItem
 * Satu baris produk di checkout: gambar, nama, harga, qty (read-only), discount picker.
 *
 * Props:
 *  - item              : object produk dari API (product_price[n])
 *  - qty               : qty dari cart (read-only)
 *  - selectedDiscountId: discount_id aktif (null = tanpa diskon)
 *  - onDiscountChange(id): callback ganti diskon
 *  - isLast            : boolean, hilangkan border-bottom jika true
 */
export default function ProductItem({
  item,
  qty,
  selectedDiscountId,
  onDiscountChange,
  isLast,
}) {
  const selectedDisc = item.discount?.find((d) => d.discount_id === selectedDiscountId) ?? null;
  const unitPrice = getUnitPrice(item, selectedDisc);
  const lineTotal = unitPrice * qty;

  return (
    <div className={`py-4 ${!isLast ? "border-b border-gray-100" : ""}`}>
      <div className="flex items-start gap-3">
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
          <p className="text-sm font-medium text-gray-900 truncate pr-1">
            {item.product_name}
          </p>

          {selectedDisc ? (
            <div className="flex items-center gap-1.5 flex-wrap mt-0.5">
              <span className="text-xs line-through text-gray-400">
                {item.product_price_format}
              </span>
              <span className="text-xs font-semibold text-blue-600">
                {selectedDisc.final_value}
              </span>
            </div>
          ) : (
            <p className="text-xs text-gray-500 mt-0.5">{item.product_price_format}</p>
          )}

          <p className="text-[11px] text-gray-400 mt-1">x{qty}</p>
        </div>

        <div className="shrink-0 text-right">
          <p className="text-sm font-semibold text-gray-900">{formatRupiah(lineTotal)}</p>
        </div>
      </div>

      <div className="mt-3">
        <DiscountPicker
          discounts={item.discount}
          selectedId={selectedDiscountId}
          onSelect={onDiscountChange}
          productPrice={item.product_price}
        />
      </div>
    </div>
  );
}