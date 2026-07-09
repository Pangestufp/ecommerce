import { useParams, useNavigate } from "react-router-dom";
import { useCustomerOrderDetail } from "./useCustomerOrderDetail";
import CustomerOrderDetailView from "./component/CustomerOrderDetailView";

export default function UserSalesOrderDetailPage() {
  const { code } = useParams();
  const navigate = useNavigate();
  const { loading, detail } = useCustomerOrderDetail(code);

  return (
    <div className="p-4 sm:p-6 space-y-4 max-w-2xl mx-auto">
      <div className="flex items-center gap-3">
        <button
          onClick={() => navigate(-1)}
          className="text-sm text-gray-500 hover:text-gray-700 transition"
        >
          ← Kembali
        </button>
        <h1 className="text-lg font-semibold text-gray-800">Detail Pesanan</h1>
      </div>

      {loading ? (
        <div className="py-16 text-center text-sm text-gray-400">Memuat...</div>
      ) : detail ? (
        <CustomerOrderDetailView detail={detail} />
      ) : (
        <div className="py-16 text-center text-sm text-gray-400">Pesanan tidak ditemukan</div>
      )}
    </div>
  );
}