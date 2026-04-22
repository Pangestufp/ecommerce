import { useState } from "react";
import Button from "../../../shared/ui/Button";

export default function UpdateInventoryModal({
  inventory,
  onSubmit,
  onClose,
  loading,
}) {
  const [form, setForm] = useState({
    cost_price: inventory.cost_price ?? "",
    stock: inventory.stock ?? "",
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    await onSubmit(inventory.batch_id, {
      cost_price: parseFloat(form.cost_price),
      stock: parseInt(form.stock),
    });
  };

  // stok tidak boleh kurang dari reserved_stock
  const minStock = inventory.reserved_stock ?? 0;

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/40 z-50">
      <div className="bg-white rounded-xl p-6 w-[360px] shadow-lg">
        <h2 className="text-base font-semibold text-gray-800 mb-1">
          Update Batch Inventory
        </h2>
        <p className="text-xs text-gray-400 mb-4">
          Batch:{" "}
          <span className="font-medium text-gray-600">
            {inventory.batch_code}
          </span>
        </p>

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

          <div className="flex flex-col gap-1">
            <label className="text-sm text-gray-600">Stok</label>
            <input
              type="number"
              name="stock"
              min={minStock}
              step={1}
              placeholder={String(minStock)}
              value={form.stock}
              onChange={(e) => {
                let val = e.target.value;
                // clamp: tidak boleh di bawah reserved_stock
                if (val !== "" && parseInt(val) < minStock)
                  val = String(minStock);
                handleChange({ target: { name: "stock", value: val } });
              }}
              required
              className="border border-gray-300 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-blue-500"
            />
            {minStock > 0 && (
              <p className="text-xs text-gray-400">
                Minimum{" "}
                <span className="font-medium text-gray-600">{minStock}</span>{" "}
                (stok yang sedang direservasi)
              </p>
            )}
          </div>

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