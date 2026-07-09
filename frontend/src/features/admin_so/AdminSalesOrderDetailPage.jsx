import { useParams, useNavigate } from "react-router-dom";
import Button from "../../shared/ui/Button";
import AdminOrderDetailView from "./component/AdminOrderDetailView";
import { useAdminOrderDetail } from "./useAdminOrderDetail";

export default function AdminSalesOrderDetailPage() {
  const { code } = useParams();
  const navigate = useNavigate();
  const { loading, detail } = useAdminOrderDetail(code);

  return (
    <div className="p-6 space-y-4 max-w-3xl mx-auto">
      <div className="flex items-center gap-3">
        <button
          onClick={() => navigate(-1)}
          className="text-sm text-gray-500 hover:text-gray-700 transition"
        >
          ← Kembali
        </button>
        <h1 className="text-lg font-semibold text-gray-800">Detail Order</h1>
      </div>

      {loading ? (
        <div className="py-16 text-center text-sm text-gray-400">Memuat...</div>
      ) : detail ? (
        <AdminOrderDetailView detail={detail} />
      ) : (
        <div className="py-16 text-center text-sm text-gray-400">Order tidak ditemukan</div>
      )}
    </div>
  );
}