import { Minus, Plus } from "lucide-react";
import { formatRupiah, getUnitPrice } from "../checkoutHelpers";
import DiscountPicker from "./DiscountPicker";

/**
 * ProductItem
 * Satu baris produk di checkout: gambar, nama, harga, qty stepper, discount picker.
 *
 * Props:
 *  - item              : object produk dari API (product_price[n])
 *  - qty               : qty saat ini
 *  - selectedDiscountId: discount_id aktif (null = tanpa diskon)
 *  - onQtyChange(delta): callback +1 / -1
 *  - onDiscountChange(id): callback ganti diskon
 *  - isLast            : boolean, hilangkan border-bottom jika true
 */
export default function ProductItem({
  item,
  qty,
  selectedDiscountId,
  onQtyChange,
  onDiscountChange,
  isLast,
}) {
  const selectedDisc = item.discount?.find((d) => d.discount_id === selectedDiscountId) ?? null;
  const unitPrice = getUnitPrice(item, selectedDisc);
  const lineTotal = unitPrice * qty;

  return (
    <div className={`py-4 ${!isLast ? "border-b border-gray-100" : ""}`}>
      {/* Baris atas: gambar + info + total */}
      <div className="flex items-start gap-3">
        {/* Thumbnail */}
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

        {/* Nama + harga satuan */}
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
        </div>

        {/* Total baris */}
        <div className="shrink-0 text-right">
          <p className="text-sm font-semibold text-gray-900">{formatRupiah(lineTotal)}</p>
        </div>
      </div>

      {/* Baris bawah: discount picker (kiri) + qty stepper (kanan) */}
      <div className="flex items-center justify-between mt-3 gap-2">
        <DiscountPicker
          discounts={item.discount}
          selectedId={selectedDiscountId}
          onSelect={onDiscountChange}
          productPrice={item.product_price}
        />

        {/* Qty stepper */}
        <div className="flex items-center gap-2 shrink-0">
          <button
            type="button"
            onClick={() => onQtyChange(-1)}
            disabled={qty <= 1}
            className="w-7 h-7 rounded-lg border border-gray-200 flex items-center justify-center text-gray-500 hover:bg-gray-50 disabled:opacity-30 transition"
          >
            <Minus size={13} />
          </button>
          <span className="text-sm font-medium text-gray-800 w-5 text-center tabular-nums">
            {qty}
          </span>
          <button
            type="button"
            onClick={() => onQtyChange(+1)}
            disabled={qty >= item.available_stock}
            className="w-7 h-7 rounded-lg border border-gray-200 flex items-center justify-center text-gray-500 hover:bg-gray-50 disabled:opacity-30 transition"
          >
            <Plus size={13} />
          </button>
        </div>
      </div>

      {/* Info stok */}
      <p className="text-[10px] text-gray-400 mt-1.5">
        Stok tersisa: {item.available_stock}
      </p>
    </div>
  );
}