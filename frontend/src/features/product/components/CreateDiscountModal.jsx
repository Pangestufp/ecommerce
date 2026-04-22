import { useState } from "react";
import TextField from "../../../shared/ui/TextField";
import Button from "../../../shared/ui/Button";
import Textarea from "../../../shared/ui/TextArea";
import SearchableSelect from "../../../shared/ui/SearchableSelect";
import DateRangeField from "../../../shared/ui/DateRangeField";

export default function CreateDiscountModal({
  productID,
  onSubmit,
  onClose,
  types,
  loading,
}) {
  const [form, setForm] = useState({
    product_id: productID,
    discount_name: "",
    discount_type: "",
    discount_value: "",
    start_at: "",
    expired_at: "",
  });

  const isPercentage = form.discount_type === "percentage";

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((prev) => ({
      ...prev,
      [name]: value,
      ...(name === "discount_type" ? { discount_value: "" } : {}),
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const rawValue = parseFloat(form.discount_value);
    await onSubmit({
      ...form,
      discount_value: isPercentage
        ? rawValue / 100
        : rawValue,
    });
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/40 z-50">
      <div className="bg-white rounded-xl p-6 w-[360px] shadow-lg">
        <h2 className="text-base font-semibold text-gray-800 mb-4">
          Buat Diskon
        </h2>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <TextField
            label="Nama Diskon"
            name="discount_name"
            placeholder="Masukkan Nama Diskon"
            value={form.discount_name}
            onChange={handleChange}
          />

          <SearchableSelect
            label="Tipe"
            name="discount_type"
            value={form.discount_type}
            onChange={handleChange}
            options={types}
            placeholder="Cari tipe..."
            required
            getValue={(t) => t.discount_type}
            getLabel={(t) => t.discount_type}
          />

          {form.discount_type && (
            <div className="flex flex-col gap-1">
              <label className="text-sm text-gray-600">Jumlah Diskon</label>

              {isPercentage ? (
                <>
                  <div className="flex items-center border border-gray-300 rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-blue-500">
                    <input
                      type="number"
                      name="discount_value"
                      min={1}
                      max={100}
                      step={1}
                      placeholder="1 – 100"
                      value={form.discount_value}
                      onChange={(e) => {
                        let val = e.target.value;
                        if (val !== "" && parseInt(val) > 100) val = "100";
                        if (val !== "" && parseInt(val) < 1) val = "1";
                        handleChange({ target: { name: "discount_value", value: val } });
                      }}
                      className="flex-1 px-3 py-2 text-sm outline-none bg-white"
                    />
                    <span className="px-3 py-2 text-sm text-gray-500 border-l border-gray-300 bg-gray-50">
                      %
                    </span>
                  </div>
                </>
              ) : (
                <div className="flex items-center border border-gray-300 rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-blue-500">
                  <span className="px-3 py-2 text-sm text-gray-500 border-r border-gray-300 bg-gray-50">
                    Rp
                  </span>
                  <input
                    type="number"
                    name="discount_value"
                    min={0}
                    step={1000}
                    placeholder="0"
                    value={form.discount_value}
                    onChange={handleChange}
                    className="flex-1 px-3 py-2 text-sm outline-none bg-white"
                  />
                </div>
              )}
            </div>
          )}

          <DateRangeField
            label="Periode Diskon"
            startName="start_at"
            endName="expired_at"
            startValue={form.start_at}
            endValue={form.expired_at}
            onChange={handleChange}
          />

          <div className="flex justify-end gap-2 mt-2">
            <Button variant="secondary" onClick={onClose} disabled={loading}>
              Batal
            </Button>
            <Button type="submit" loading={loading}>
              Simpan
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}