import { useState } from "react";
import { ChevronDown, ChevronUp, Truck, BadgeCheck, X, MapPin } from "lucide-react";
import { formatRupiah } from "../checkoutHelpers";

function CourierModal({ data, selectedCourier, onSelect, onClose }) {
  const [expandedCode, setExpandedCode] = useState(
    selectedCourier?.code ?? data.shipping_service[0]?.code ?? null
  );

  const toggleExpand = (code) => {
    setExpandedCode((prev) => (prev === code ? null : code));
  };

  const handleSelect = (courierName, courierCode, option) => {
    onSelect({
      name: courierName,
      code: courierCode,
      service: option.service,
      description: option.description,
      cost: option.cost,
      etd: option.etd,
      display_name: option.display_name,
      group: option.group,
      is_recommended: option.is_recommended,
    });
    onClose();
  };

  const isOptionSelected = (code, service) =>
    selectedCourier?.code === code && selectedCourier?.service === service;

  return (
    <div
      className="fixed inset-0 z-50 bg-black/40 flex items-end sm:items-center justify-center"
      onClick={(e) => e.target === e.currentTarget && onClose()}
    >
      <div
        className="bg-white w-full sm:max-w-md flex flex-col rounded-t-2xl sm:rounded-2xl"
        style={{ height: "80vh" }}
      >
        {/* Header — tidak ikut scroll */}
        <div className="flex items-center justify-between px-4 pt-4 pb-3 border-b border-gray-100 shrink-0">
          <div className="min-w-0 flex-1 pr-3">
            <h2 className="text-sm font-semibold text-gray-900">Pilih Layanan Pengiriman</h2>
            <p className="text-[11px] text-gray-400 mt-0.5 truncate">{data.origin_name}</p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="w-7 h-7 rounded-lg flex items-center justify-center text-gray-400 hover:bg-gray-100 transition shrink-0"
          >
            <X size={15} />
          </button>
        </div>

        {/* List kurir — area scroll */}
        <div className="flex-1 overflow-y-auto px-4 py-3 flex flex-col gap-2">
          {data.shipping_service.map((courier) => {
            const isExpanded = expandedCode === courier.code;
            const cheapest = courier.option[0];

            return (
              <div
                key={courier.code}
                className="border border-gray-200 rounded-xl overflow-hidden"
              >
                <button
                  type="button"
                  onClick={() => toggleExpand(courier.code)}
                  className="w-full flex items-center justify-between px-3 py-2.5 hover:bg-gray-50 transition"
                >
                  <div className="flex items-center gap-2 min-w-0">
                    <div className="w-7 h-7 rounded-lg bg-gray-100 flex items-center justify-center shrink-0">
                      <Truck size={13} className="text-gray-500" />
                    </div>
                    <div className="text-left min-w-0">
                      <p className="text-xs font-semibold text-gray-800 truncate">{courier.name}</p>
                      <p className="text-[11px] text-gray-400">
                        mulai {formatRupiah(cheapest.cost)} · {courier.option.length} layanan
                      </p>
                    </div>
                  </div>
                  {isExpanded ? (
                    <ChevronUp size={14} className="text-gray-400 shrink-0" />
                  ) : (
                    <ChevronDown size={14} className="text-gray-400 shrink-0" />
                  )}
                </button>

                {isExpanded && (
                  <div className="border-t border-gray-100 divide-y divide-gray-100">
                    {courier.option.map((option) => {
                      const selected = isOptionSelected(courier.code, option.service);
                      return (
                        <button
                          key={option.service}
                          type="button"
                          onClick={() => handleSelect(courier.name, courier.code, option)}
                          className={`w-full text-left px-3 py-2.5 flex items-center justify-between gap-3 transition ${
                            selected ? "bg-blue-50" : "hover:bg-gray-50"
                          }`}
                        >
                          <div className="min-w-0 flex-1">
                            <div className="flex items-center gap-1.5 flex-wrap">
                              <span
                                className={`text-xs font-medium ${
                                  selected ? "text-blue-700" : "text-gray-800"
                                }`}
                              >
                                {option.display_name}
                              </span>
                              {option.is_recommended && (
                                <span className="inline-flex items-center gap-0.5 text-[10px] px-1.5 py-0.5 rounded-full bg-green-100 text-green-700">
                                  <BadgeCheck size={9} />
                                  Rekomendasi
                                </span>
                              )}
                            </div>
                            <p className="text-[11px] text-gray-400 mt-0.5">
                              {option.description}
                              {option.etd ? ` · Est. ${option.etd}` : ""}
                            </p>
                          </div>
                          <div className="shrink-0 text-right">
                            <p
                              className={`text-xs font-semibold ${
                                selected ? "text-blue-700" : "text-gray-800"
                              }`}
                            >
                              {formatRupiah(option.cost)}
                            </p>
                            {selected && (
                              <span className="text-[10px] text-blue-500">Terpilih</span>
                            )}
                          </div>
                        </button>
                      );
                    })}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

export default function CourierSection({
  courierData,
  selectedCourier,
  onSelect,
  onLoad,
  loading,
  disabled,
}) {
  const [modalOpen, setModalOpen] = useState(false);

  const openModal = async () => {
    if (!courierData) {
      await onLoad();
    }
    setModalOpen(true);
  };

  return (
    <>
      <div className="bg-white rounded-2xl border border-gray-100 p-4 mb-3">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <Truck size={14} className="text-blue-500 shrink-0" />
            <h3 className="text-sm font-semibold text-gray-800">Kurir Pengiriman</h3>
          </div>

          <button
            type="button"
            onClick={openModal}
            disabled={disabled || loading}
            className="text-xs text-blue-600 font-medium disabled:text-gray-300 hover:text-blue-700 transition"
          >
            {loading ? "Memuat..." : selectedCourier ? "Ganti" : "Pilih Kurir"}
          </button>
        </div>

        {selectedCourier ? (
          <div className="border border-green-200 bg-green-50 rounded-xl px-3 py-2.5">
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0 flex-1">
                <div className="flex items-center gap-1.5 flex-wrap">
                  <span className="text-xs font-semibold text-gray-800">
                    {selectedCourier.display_name}
                  </span>
                  {selectedCourier.is_recommended && (
                    <span className="inline-flex items-center gap-0.5 text-[10px] px-1.5 py-0.5 rounded-full bg-green-100 text-green-700">
                      <BadgeCheck size={9} />
                      Rekomendasi
                    </span>
                  )}
                </div>
                <p className="text-[11px] text-gray-500 mt-0.5">
                  {selectedCourier.description}
                  {selectedCourier.etd ? ` · Est. ${selectedCourier.etd}` : ""}
                </p>

                {courierData && (
                  <div className="flex items-start gap-1.5 mt-2 text-[11px] text-gray-400">
                    <MapPin size={11} className="mt-0.5 shrink-0 text-gray-400" />
                    <div className="min-w-0">
                      <p className="truncate">{courierData.origin_name}</p>
                      <p className="text-center leading-none my-0.5">↓</p>
                      <p className="truncate">{courierData.destination_name}</p>
                    </div>
                  </div>
                )}
              </div>

              <p className="text-xs font-bold text-gray-900 whitespace-nowrap shrink-0">
                {formatRupiah(selectedCourier.cost)}
              </p>
            </div>
          </div>
        ) : (
          <button
            type="button"
            onClick={openModal}
            disabled={disabled || loading}
            className="w-full border border-dashed border-gray-200 rounded-xl py-4 flex flex-col items-center gap-1.5 text-gray-400 hover:border-blue-300 hover:text-blue-400 disabled:opacity-40 disabled:cursor-not-allowed transition"
          >
            <Truck size={18} />
            <span className="text-xs">
              {disabled ? "Pilih alamat terlebih dahulu" : "Ketuk untuk memilih kurir"}
            </span>
          </button>
        )}
      </div>

      {modalOpen && courierData && (
        <CourierModal
          data={courierData}
          selectedCourier={selectedCourier}
          onSelect={onSelect}
          onClose={() => setModalOpen(false)}
        />
      )}
    </>
  );
}