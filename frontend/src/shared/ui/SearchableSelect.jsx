import { useMemo, useRef, useState } from "react";

export default function SearchableSelect({
  label,
  name,
  value,
  onChange,
  options = [],
  placeholder,
  getValue,
  getLabel,
  required = false,
}) {
  const inputRef = useRef(null);
  const [query, setQuery] = useState("");
  const [open, setOpen] = useState(false);
  const [dropdownStyle, setDropdownStyle] = useState({});

  const selected = options.find((o) => getValue(o) === value);

  const filtered = useMemo(() => {
    if (!query) return options;
    return options.filter((o) =>
      getLabel(o).toLowerCase().includes(query.toLowerCase())
    );
  }, [query, options, getLabel]);

  const handleSelect = (opt) => {
    onChange({ target: { name, value: getValue(opt) } });
    setQuery(getLabel(opt));
    setOpen(false);
  };

  const handleFocus = () => {
    if (inputRef.current) {
      const rect = inputRef.current.getBoundingClientRect();
      setDropdownStyle({
        position: "fixed",
        top: rect.bottom + 4,
        left: rect.left,
        width: rect.width,
        zIndex: 9999,
      });
    }
    setOpen(true);
  };

  return (
    <div className="flex flex-col sm:flex-row sm:items-start gap-2 w-full">
      {label && (
        <div className="w-full sm:w-28 text-center sm:text-left">
          <label
            htmlFor={name}
            className="w-28 pt-2 text-sm font-medium text-gray-700"
          >
            {label}
          </label>
        </div>
      )}

      <div className="flex flex-col flex-1">
        <input
          autoComplete="off"
          ref={inputRef}
          id={name}
          name={name}
          value={query || (selected ? getLabel(selected) : "")}
          placeholder={placeholder}
          required={required}
          onFocus={handleFocus}
          onBlur={() => setOpen(false)}
          onChange={(e) => {
            setQuery(e.target.value);
            setOpen(true);
          }}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm outline-none transition focus:ring-2 focus:ring-blue-200"
        />

        {open && (
          <div
            style={dropdownStyle}
            className="bg-white border rounded-lg max-h-60 overflow-auto shadow-lg"
          >
            {filtered.length === 0 ? (
              <div className="px-3 py-2 text-sm text-gray-500">
                Tidak ditemukan
              </div>
            ) : (
              filtered.map((opt) => (
                <div
                  key={getValue(opt)}
                  onMouseDown={(e) => {
                    e.preventDefault();
                    handleSelect(opt);
                  }}
                  className="px-3 py-2 text-sm hover:bg-gray-100 cursor-pointer text-left"
                >
                  {getLabel(opt)}
                </div>
              ))
            )}
          </div>
        )}
      </div>
    </div>
  );
}