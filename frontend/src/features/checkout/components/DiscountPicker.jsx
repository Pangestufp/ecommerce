import { useState } from "react";
import { ChevronDown, ChevronUp, Tag } from "lucide-react";
import { formatRupiah } from "../checkoutHelpers";

/**
 * DiscountPicker
 * Expand/collapse list diskon per produk.
 * Auto-terpilih diskon terbaik saat mount (dihandle dari parent).
 *
 * Props:
 *  - discounts       : array diskon dari API
 *  - selectedId      : discount_id yang sedang aktif (null = tanpa diskon)
 *  - onSelect(id)    : callback saat user ganti pilihan (null = tanpa diskon)
 *  - productPrice    : harga asli produk (string/number)
 */
export default function DiscountPicker({ discounts, selectedId, onSelect, productPrice }) {
  const [open, setOpen] = useState(false);

  if (!discounts || discounts.length === 0) return null;

  const selected = discounts.find((d) => d.discount_id === selectedId);

  return (
    <div className="mt-2 flex-1 min-w-0">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className="flex items-center gap-1.5 text-xs text-blue-600 font-medium hover:text-blue-700 transition-colors"
      >
        <Tag size={12} className="shrink-0" />
        <span className="truncate">
          {selected ? (
            <>
              Diskon: <span className="font-semibold">{selected.discount_name}</span>
              &nbsp;→&nbsp;{selected.final_value}
            </>
          ) : (
            "Pilih diskon"
          )}
        </span>
        {open ? <ChevronUp size={12} className="shrink-0" /> : <ChevronDown size={12} className="shrink-0" />}
      </button>

      {open && (
        <div className="mt-2 flex flex-col gap-1.5">
          {/* Opsi tanpa diskon */}
          <label
            className={`flex items-start gap-2 rounded-lg border px-3 py-2 cursor-pointer transition-all
              ${!selectedId
                ? "border-blue-500 bg-blue-50"
                : "border-gray-200 bg-white hover:bg-gray-50"
              }`}
          >
            <input
              type="radio"
              className="mt-0.5 accent-blue-600"
              checked={!selectedId}
              onChange={() => onSelect(null)}
            />
            <div className="flex-1 min-w-0">
              <p className="text-xs font-medium text-gray-700">Tanpa diskon</p>
              <p className="text-xs text-gray-400">{formatRupiah(productPrice)}</p>
            </div>
          </label>

          {discounts.map((d) => (
            <label
              key={d.discount_id}
              className={`flex items-start gap-2 rounded-lg border px-3 py-2 cursor-pointer transition-all
                ${selectedId === d.discount_id
                  ? "border-blue-500 bg-blue-50"
                  : "border-gray-200 bg-white hover:bg-gray-50"
                }`}
            >
              <input
                type="radio"
                className="mt-0.5 accent-blue-600"
                checked={selectedId === d.discount_id}
                onChange={() => onSelect(d.discount_id)}
              />
              <div className="flex-1 min-w-0">
                <p className="text-xs font-medium text-gray-800 truncate">{d.discount_name}</p>
                <p className="text-xs text-gray-400">
                  {d.discount_type === "percentage"
                    ? `Potongan ${d.discount_value_format}`
                    : `Potongan ${d.discount_Amount_format}`}
                </p>
                <p className="text-xs font-semibold text-blue-600 mt-0.5">{d.final_value}</p>
                <p className="text-[10px] text-gray-400">s/d {d.expired_at_format}</p>
              </div>
            </label>
          ))}
        </div>
      )}
    </div>
  );
}