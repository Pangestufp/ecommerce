import { useState } from "react";
import { ChevronDown, ChevronUp, MapPin, Check } from "lucide-react";

/**
 * AddressCard
 * Tampilkan satu alamat dalam card.
 *
 * Props:
 *  - address    : object alamat dari API
 *  - isSelected : boolean, tampilkan checkmark biru jika true
 */
export function AddressCard({ address, isSelected }) {
  return (
    <div
      className={`rounded-xl border px-3 py-2.5 transition-all
        ${isSelected
          ? "border-blue-500 bg-blue-50"
          : "border-gray-200 bg-white hover:border-blue-300 hover:bg-blue-50/40"
        }`}
    >
      <div className="flex items-start justify-between gap-2">
        <div className="flex-1 min-w-0">
          {/* Nama + badge label + badge utama */}
          <div className="flex items-center gap-1.5 flex-wrap">
            <span className="text-xs font-semibold text-gray-800">
              {address.recipient_name}
            </span>
            <span className="text-[10px] bg-gray-100 text-gray-500 rounded px-1.5 py-0.5">
              {address.label}
            </span>
            {address.is_primary === 1 && (
              <span className="text-[10px] bg-blue-100 text-blue-600 rounded px-1.5 py-0.5 font-medium">
                Utama
              </span>
            )}
          </div>

          <p className="text-xs text-gray-500 mt-1">{address.phone}</p>
          <p className="text-xs text-gray-600 mt-0.5 leading-relaxed">
            {address.additional_address}, {address.sub_district_name},{" "}
            {address.district_name}, {address.city_name},{" "}
            {address.province_name} {address.zip_code}
          </p>
        </div>

        {isSelected && (
          <div className="shrink-0 w-5 h-5 rounded-full bg-blue-500 flex items-center justify-center mt-0.5">
            <Check size={11} className="text-white" />
          </div>
        )}
      </div>
    </div>
  );
}

/**
 * AddressSection
 * Tampilkan alamat terpilih + expand untuk ganti alamat.
 * Auto-pilih primary atau index pertama (dihandle dari parent).
 *
 * Props:
 *  - addresses          : array alamat dari API
 *  - selectedId         : address_id yang sedang aktif
 *  - onSelect(id)       : callback saat user pilih alamat lain
 */
export function AddressSection({ addresses, selectedId, onSelect }) {
  const [expanded, setExpanded] = useState(false);

  const selected = addresses.find((a) => a.address_id === selectedId);
  const others = addresses.filter((a) => a.address_id !== selectedId);

  if (!selected) return (
  <div className="bg-white rounded-2xl border border-gray-100 px-4 py-4 mb-3">
    <div className="flex items-center gap-2 mb-3">
      <MapPin size={14} className="text-blue-500 shrink-0" />
      <span className="text-sm font-semibold text-gray-800">Alamat Pengiriman</span>
    </div>
    <p className="text-xs text-gray-500">
      Belum ada alamat.{" "}
      <a href="/alamat" className="text-blue-600 font-medium hover:underline">
        Klik di sini untuk menambahkan alamat
      </a>
    </p>
  </div>
);

  const handleSelect = (id) => {
    onSelect(id);
    setExpanded(false);
  };

  return (
    <div className="bg-white rounded-2xl border border-gray-100 mb-3 overflow-hidden">
      {/* Alamat aktif */}
      <div className="px-4 pt-4 pb-3">
        <div className="flex items-center gap-2 mb-3">
          <MapPin size={14} className="text-blue-500 shrink-0" />
          <span className="text-sm font-semibold text-gray-800">Alamat Pengiriman</span>
        </div>
        <AddressCard address={selected} isSelected />
      </div>

      {/* Tombol expand + daftar alamat lain */}
      {others.length > 0 && (
        <>
          <button
            type="button"
            onClick={() => setExpanded((v) => !v)}
            className="w-full flex items-center justify-center gap-1.5 text-xs text-blue-600 font-medium py-2.5 border-t border-gray-100 hover:bg-gray-50 transition-colors"
          >
            {expanded ? (
              <>Sembunyikan alamat lain <ChevronUp size={13} /></>
            ) : (
              <>Lihat {others.length} alamat lain <ChevronDown size={13} /></>
            )}
          </button>

          {expanded && (
            <div className="border-t border-gray-100 px-4 pb-3 pt-3 flex flex-col gap-2">
              {others.map((addr) => (
                <button
                  key={addr.address_id}
                  type="button"
                  onClick={() => handleSelect(addr.address_id)}
                  className="text-left w-full"
                >
                  <AddressCard address={addr} isSelected={false} />
                </button>
              ))}
            </div>
          )}
        </>
      )}
    </div>
  );
}