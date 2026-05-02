import { useState, useEffect } from "react";
import { Search, X } from "lucide-react";

export default function LongSearchBar({
  placeholder = "Cari produk, kategori, atau brand...",
  onChange,
  delay = 500,
}) {
  const [value, setValue] = useState("");
  const [isFirstRender, setIsFirstRender] = useState(true);
  const [isFocused, setIsFocused] = useState(false);

  useEffect(() => {
    if (isFirstRender) {
      setIsFirstRender(false);
      return;
    }
    const handler = setTimeout(() => {
      onChange?.(value);
    }, delay);

    return () => clearTimeout(handler);
  }, [value, delay]);

  return (
    <div className="w-full max-w-xl">
      <div
        className={`
          flex items-center gap-2.5
          bg-white px-4 py-2.5 rounded-xl
          border transition-all duration-150
          ${isFocused
            ? "border-blue-500 shadow-[0_0_0_3px_rgba(37,99,235,0.10)]"
            : "border-gray-200 shadow-sm"
          }
        `}
      >
        <Search
          size={18}
          className={`shrink-0 transition-colors duration-150 ${isFocused ? "text-blue-500" : "text-gray-400"}`}
        />
        <input
          type="text"
          value={value}
          placeholder={placeholder}
          onChange={(e) => setValue(e.target.value)}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          className="flex-1 border-none outline-none text-sm text-gray-900 bg-transparent placeholder:text-gray-400"
        />
        {value && (
          <button onClick={() => setValue("")} className="flex items-center p-0">
            <X size={15} className="text-gray-400 hover:text-gray-600 transition-colors" />
          </button>
        )}
      </div>
    </div>
  );
}