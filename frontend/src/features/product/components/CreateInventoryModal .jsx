import { useState } from "react";
import TextField from "../../../shared/ui/TextField";
import Button from "../../../shared/ui/Button";

export default function CreateInventoryModal({
  productID,
  onSubmit,
  onClose,
  loading,
}) {
  const [form, setForm] = useState({
    product_id: productID,
    cost_price: "",
    stock: "",
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    await onSubmit({
      ...form,
      cost_price: parseFloat(form.cost_price),
      stock: parseInt(form.stock),
    });
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/40 z-50">
      <div className="bg-white rounded-xl p-6 w-[360px] shadow-lg">
        <h2 className="text-base font-semibold text-gray-800 mb-4">
          Tambah Batch Inventory
        </h2>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1">
            <label className="text-sm text-gray-600">Harga Modal</label>
            <div className="flex items-center border border-gray-300 rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-blue-500">
              <span className="px-3 py-2 text-sm text-gray-500 border-r border-gray-300 bg-gray-50">
                Rp
              </span>
              <input
                type="number"
                name="cost_price"
                min={0}
                step={1000}
                placeholder="0"
                value={form.cost_price}
                onChange={handleChange}
                required
                className="flex-1 px-3 py-2 text-sm outline-none bg-white"
              />
            </div>
          </div>

          <TextField
            label="Stok"
            name="stock"
            type="number"
            placeholder="Jumlah stok"
            value={form.stock}
            onChange={(e) => {
              let val = e.target.value;
              // stok tidak boleh negatif
              if (val !== "" && parseInt(val) < 0) val = "0";
              handleChange({ target: { name: "stock", value: val } });
            }}
            required
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