const STATUS_CONFIG = {
  PENDING_PAYMENT: { label: "Menunggu Pembayaran", color: "bg-yellow-100 text-yellow-700 border-yellow-200" },
  PAYED:           { label: "Sudah Dibayar",        color: "bg-blue-100 text-blue-700 border-blue-200" },
  EXPIRED:         { label: "Kadaluarsa",            color: "bg-gray-100 text-gray-500 border-gray-200" },
  CANCELLED_UNPAID:{ label: "Dibatalkan (Belum Bayar)", color: "bg-gray-100 text-gray-500 border-gray-200" },
  ACCEPTED:        { label: "Diterima",              color: "bg-indigo-100 text-indigo-700 border-indigo-200" },
  PACKED:          { label: "Dikemas",               color: "bg-purple-100 text-purple-700 border-purple-200" },
  SHIPPED:         { label: "Dikirim",               color: "bg-cyan-100 text-cyan-700 border-cyan-200" },
  CANCELLED:       { label: "Dibatalkan",            color: "bg-red-100 text-red-600 border-red-200" },
  FINISHED:        { label: "Selesai",               color: "bg-green-100 text-green-700 border-green-200" },
};

export const ALL_STATUSES = Object.keys(STATUS_CONFIG);

export default function StatusBadge({ status, size = "sm" }) {
  const cfg = STATUS_CONFIG[status] ?? { label: status, color: "bg-gray-100 text-gray-500 border-gray-200" };
  const text = size === "sm" ? "text-xs" : "text-sm";
  return (
    <span className={`inline-flex items-center px-2 py-0.5 rounded-full border font-medium ${text} ${cfg.color}`}>
      {cfg.label}
    </span>
  );
}

export { STATUS_CONFIG };