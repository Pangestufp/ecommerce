import { useState, useEffect } from "react";

export default function SearchBar({
  placeholder = "Search...",
  onChange,
  delay = 500,
}) {
  const [value, setValue] = useState("");
  const [isFirstRender, setIsFirstRender] = useState(true);

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
    <div className="w-full max-w-sm">
      <input
        type="text"
        value={value}
        placeholder={placeholder}
        onChange={(e) => setValue(e.target.value)}
        className="
          w-full
          px-3 py-2
          border border-gray-300
          rounded-lg
          text-sm
          outline-none
          transition
          focus:ring-2
          focus:ring-blue-200
        "
      />
    </div>
  );
}