import { ALL_STATUSES, STATUS_CONFIG } from "./StatusBadge";

// Pill-based multi-select filter — dipakai di kedua halaman list.
export default function StatusFilter({ selected, onToggle, onClear }) {
  return (
    <div className="flex flex-wrap gap-2 items-center">
      <button
        onClick={onClear}
        className={`px-3 py-1 rounded-full border text-xs font-medium transition
          ${selected.length === 0
            ? "bg-gray-800 text-white border-gray-800"
            : "bg-white text-gray-500 border-gray-300 hover:border-gray-400"
          }`}
      >
        Semua
      </button>

      {ALL_STATUSES.map((status) => {
        const cfg = STATUS_CONFIG[status];
        const active = selected.includes(status);
        return (
          <button
            key={status}
            onClick={() => onToggle(status)}
            className={`px-3 py-1 rounded-full border text-xs font-medium transition
              ${active ? cfg.color + " ring-1 ring-offset-0" : "bg-white text-gray-500 border-gray-300 hover:border-gray-400"}`}
          >
            {cfg.label}
          </button>
        );
      })}
    </div>
  );
}