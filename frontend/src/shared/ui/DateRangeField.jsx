export default function DateRangeField({
  label,
  startName = "start_at",
  endName = "expired_at",
  startValue,
  endValue,
  onChange,
  required = false,
  disabled = false,
}) {
  return (
    <div className="flex flex-col gap-1 w-full">
      {label && (
        <label className="text-sm font-medium text-gray-700">
          {label}
        </label>
      )}

      <div className="flex gap-2">
        <input
          type="date"
          name={startName}
          value={startValue || ""}
          onChange={onChange}
          required={required}
          disabled={disabled}
          className="border border-gray-300 rounded px-3 py-2 text-sm w-full focus:outline-none focus:ring-2 focus:ring-blue-500"
        />

        <input
          type="date"
          name={endName}
          value={endValue || ""}
          onChange={onChange}
          required={required}
          disabled={disabled}
          min={startValue}
          className="border border-gray-300 rounded px-3 py-2 text-sm w-full focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <p className="text-xs text-gray-400">
        Berlaku dari awal hari sampai akhir hari
      </p>
    </div>
  );
}