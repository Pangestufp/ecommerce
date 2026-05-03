import { Trash2 } from "lucide-react";

export default function CartLineItem({ item, checked, onCheck, onQtyChange, onRemove }) {
  return (
    <div className="flex items-center gap-3 py-4 border-b border-gray-100 last:border-0">

      {/* Checkbox */}
      <input
        type="checkbox"
        checked={checked}
        onChange={() => onCheck(item.product_id)}
        className="w-4 h-4 rounded accent-gray-900 shrink-0 cursor-pointer"
      />

      {/* Gambar */}
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

      {/* Info */}
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-gray-900 truncate">
          {item.product_name}
        </p>
        <p className="text-sm font-medium text-gray-700 mt-0.5">
          {item.Price_format}
        </p>
      </div>

      {/* Qty control */}
      <div className="flex items-center gap-1.5 shrink-0">
        <button
          onClick={() => onQtyChange(item.product_id, item.qty - 1)}
          className="w-7 h-7 rounded-lg border border-gray-200 text-gray-600 hover:border-gray-400 transition-colors flex items-center justify-center text-base leading-none"
        >
          −
        </button>
        <span className="text-sm font-medium text-gray-900 w-5 text-center">
          {item.qty}
        </span>
        <button
          onClick={() => onQtyChange(item.product_id, item.qty + 1)}
          className="w-7 h-7 rounded-lg border border-gray-200 text-gray-600 hover:border-gray-400 transition-colors flex items-center justify-center text-base leading-none"
        >
          +
        </button>
      </div>

      {/* Hapus */}
      <button
        onClick={() => onRemove(item.product_id)}
        className="shrink-0 text-gray-300 hover:text-red-400 transition-colors"
      >
        <Trash2 size={16} />
      </button>

    </div>
  );
}