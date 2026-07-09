import { useParams } from "react-router-dom";

export default function PaymentPage() {
  const { salesOrderCode } = useParams();

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center px-4">
      <div className="bg-white rounded-2xl border border-gray-100 p-6 max-w-md w-full text-center">
        <p className="text-xs text-gray-400 mb-1">Kode Pesanan</p>
        <h1 className="text-lg font-semibold text-gray-900 mb-4">{salesOrderCode}</h1>
        <p className="text-sm text-gray-500">
          Pesanan Anda berhasil dibuat. Detail pembayaran untuk pesanan ini akan
          ditampilkan di sini.
        </p>
      </div>
    </div>
  );
}