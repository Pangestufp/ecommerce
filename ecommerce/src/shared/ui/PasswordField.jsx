import { Eye, EyeOff } from "lucide-react"
import { useState } from "react"

export default function PasswordField({
  label,
  placeholder,
  value,
  onChange,
  error,
  hint,
  disabled = false,
  autoFocus = false,
  name,
  id,
}) {
  const [showPassword, setShowPassword] = useState(false)

  const inputId = id || name || label?.toLowerCase().replace(/\s+/g, "-")
  const type = showPassword ? "text" : "password"

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

      <div className="relative flex flex-col flex-1">
        <input
          id={inputId}
          name={name}
          type={type}
          placeholder={placeholder}
          value={value}
          onChange={onChange}
          disabled={disabled}
          autoFocus={autoFocus}
          className={`
            w-full
            px-3 py-2 pr-16
            border rounded-lg
            text-sm
            outline-none
            transition
            ${error ? "border-red-500 focus:ring-red-200" : "border-gray-300 focus:ring-blue-200"}
            focus:ring-2
            disabled:bg-gray-100
          `}
        />

        <button
          type="button"
          onClick={() => setShowPassword((v) => !v)}
          disabled={disabled}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-sm text-gray-500 hover:text-gray-700"
        >
          {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
        </button>
      </div>

      {(error || hint) && (
        <p className={`text-xs ${error ? "text-red-500" : "text-gray-500"}`}>
          {error || hint}
        </p>
      )}
    </div>
  )
}