export default function Button({
  children,
  type = "button",
  onClick,
  loading = false,
  disabled = false,
  variant = "primary",
  className = "",
}) {
  const base =
    "px-4 py-2 rounded font-medium transition disabled:opacity-50 disabled:cursor-not-allowed"

  const variants = {
    primary: "bg-blue-600 text-white hover:bg-blue-700",
    secondary: "bg-gray-200 text-black hover:bg-gray-300",
    danger: "bg-red-600 text-white hover:bg-red-700",
  }

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled || loading}
      className={`${base} ${variants[variant]} ${className}`}
    >
      {loading ? "Loading..." : children}
    </button>
  )
}