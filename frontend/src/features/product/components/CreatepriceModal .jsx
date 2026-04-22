import { useState } from "react";
import Button from "../../../shared/ui/Button";

export default function CreatePriceModal({
  productID,
  onSubmit,
  onClose,
  loading,
}) {
  const [form, setForm] = useState({
    product_id: productID,
    product_price: "",
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    await onSubmit({
      ...form,
      product_price: parseFloat(form.product_price),
    });
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/40 z-50">
      <div className="bg-white rounded-xl p-6 w-[360px] shadow-lg">
        <h2 className="text-base font-semibold text-gray-800 mb-1">
          Tambah Harga
        </h2>
        <p className="text-xs text-gray-400 mb-4">
          Harga terbaru akan menjadi harga aktif produk.
        </p>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1">
            <label className="text-sm text-gray-600">Harga Produk</label>
            <div className="flex items-center border border-gray-300 rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-blue-500">
              <span className="px-3 py-2 text-sm text-gray-500 border-r border-gray-300 bg-gray-50">
                Rp
              </span>
              <input
                type="number"
                min={0}
                step={1000}
                placeholder="0"
                value={form.product_price}
                onChange={(e) =>
                  setForm((prev) => ({ ...prev, product_price: e.target.value }))
                }
                required
                className="flex-1 px-3 py-2 text-sm outline-none bg-white"
              />
            </div>
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