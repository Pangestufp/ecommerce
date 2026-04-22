export default function TextField({
  label,
  type = "text",
  placeholder,
  value,
  onChange,
  error,
  hint,
  disabled = false,
  autoFocus = false,
  name,
  id,
  min,
}) {
  const inputId = id || name || label?.toLowerCase().replace(/\s+/g, "-");

  return (
    <div className="flex flex-col sm:flex-row sm:items-start gap-2 w-full">
      {label && (
        <div className="w-full sm:w-28 text-center sm:text-left">
        <label
          htmlFor={inputId}
          className="w-28 pt-2 text-sm font-medium text-gray-700"
        >
          {label}
        </label>
        </div>
      )}

      <div className="flex flex-col flex-1">
        <input
          id={inputId}
          name={name}
          type={type}
          placeholder={placeholder}
          value={value}
          onChange={onChange}
          disabled={disabled}
          autoFocus={autoFocus}
          required
          min={type === "number" ? (min ?? 0) : undefined}
          className={`
          w-full
          px-3 py-2
          border rounded-lg
          text-sm
          outline-none
          transition
          ${error ? "border-red-500 focus:ring-red-200" : "border-gray-300 focus:ring-blue-200"}
          focus:ring-2
          disabled:bg-gray-100
        `}
        />
      </div>

      {(error || hint) && (
        <p className={`text-xs ${error ? "text-red-500" : "text-gray-500"}`}>
          {error || hint}
        </p>
      )}
    </div>
  );
}
