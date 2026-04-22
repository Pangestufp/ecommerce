import { useNavigate } from "react-router-dom";
import Button from "../../../shared/ui/Button";

export default function TableProduct({ data = [], onUpdate, onDelete }) {
  const navigate = useNavigate();

  const columns = [
    { key: "product_code", label: "Kode Produk" },
    { key: "product_name", label: "Nama Produk" },
    { key: "type_name", label: "Tipe" },
    { key: "weight_gram", label: "Berat (g)" },
    { key: "stock", label: "Stok" },
    { key: "reserved_stock", label: "Reserved" },
    { key: "product_price", label: "Harga" },
    { key: "available_discount", label: "Diskon Aktif" },
  ];

  const formatPrice = (row) => {
    if (row.is_price_set !== 1) return <span className="text-xs text-gray-400">Belum diset</span>;
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(row.product_price);
  };

  const renderCell = (row, key) => {
    if (key === "product_price") return formatPrice(row);

    if (key === "stock") return row.is_stock_set !== 1
      ? <span className="text-xs text-gray-400">Belum diset</span>
      : row.stock;

    if (key === "available_discount") return row.available_discount > 0
      ? <span className="px-2 py-0.5 bg-blue-100 text-blue-700 rounded-full text-xs">{row.available_discount} diskon</span>
      : <span className="text-xs text-gray-400">-</span>;

    return row[key] ?? "-";
  };

  return (
    <div className="overflow-x-auto rounded-lg border border-gray-200">
      <table className="min-w-full text-sm text-left text-gray-700">
        <thead className="bg-gray-50 text-xs text-gray-500 uppercase">
          <tr>
            {columns.map((col) => (
              <th key={col.key} className="px-4 py-3 font-medium whitespace-nowrap">
                {col.label}
              </th>
            ))}
            <th className="px-4 py-3 font-medium text-right">Aksi</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100">
          {data.length === 0 ? (
            <tr>
              <td colSpan={columns.length + 1} className="px-4 py-6 text-center text-gray-400">
                Tidak ada data
              </td>
            </tr>
          ) : (
            data.map((row) => (
              <tr key={row.product_id} className="hover:bg-gray-50 transition">
                {columns.map((col) => (
                  <td key={col.key} className="px-4 py-3 whitespace-nowrap">
                    {renderCell(row, col.key)}
                  </td>
                ))}
                <td className="px-4 py-3 text-right">
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="secondary"
                      className="text-xs px-3 py-1"
                      onClick={() => navigate(`/produk/${row.product_id}`)}
                    >
                      Detail
                    </Button>
                    <Button
                      variant="secondary"
                      className="text-xs px-3 py-1"
                      onClick={() => onUpdate(row)}
                    >
                      Edit
                    </Button>
                    <Button
                      variant="danger"
                      className="text-xs px-3 py-1"
                      onClick={() => onDelete(row.product_id)}
                    >
                      Hapus
                    </Button>
                  </div>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}