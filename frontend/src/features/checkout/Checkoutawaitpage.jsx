import { useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import ApiCheckout from "./apiCheckout";

const POLL_INTERVAL_MS = 2000;

export default function CheckoutAwaitPage() {
  const { idempotencyKey } = useParams();
  const navigate = useNavigate();

  const [status, setStatus] = useState("pending");
  const [errorMessage, setErrorMessage] = useState(null);

  const intervalRef = useRef(null);
  const redirectedRef = useRef(false);

  useEffect(() => {
    redirectedRef.current = false;

    const stopPolling = () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };

    const poll = async () => {
      try {
        const res = await ApiCheckout.getCheckoutStatus(idempotencyKey);
        const data = res.data.data; // CheckoutStatusResponse

        if (redirectedRef.current) return;

        setStatus(data.status);

        // begitu order sudah tidak null, langsung pindah ke halaman payment
        if (data.order) {
          redirectedRef.current = true;
          stopPolling();
          navigate(`/payment/${data.order.sales_order_code}`, { replace: true });
          return;
        }

        if (data.status === "failed" || data.error_message) {
          stopPolling();
          setErrorMessage(data.error_message || "Checkout gagal diproses, silakan coba lagi.");
        }
      } catch (err) {
        stopPolling();
        setErrorMessage(err.message || "Terjadi kesalahan saat memeriksa status pesanan");
      }
    };

    poll();
    intervalRef.current = setInterval(poll, POLL_INTERVAL_MS);

    return () => stopPolling();
  }, [idempotencyKey, navigate]);

  if (errorMessage) {
    return (
      <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center px-4">
        <div className="bg-white rounded-2xl border border-gray-100 p-6 max-w-sm w-full text-center">
          <p className="text-sm text-red-600 mb-4">{errorMessage}</p>
          <button
            onClick={() => navigate("/keranjang")}
            className="text-sm font-medium text-blue-600"
          >
            Kembali ke Keranjang
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-center px-4">
      <div className="h-10 w-10 border-4 border-gray-200 border-t-blue-600 rounded-full animate-spin mb-4" />
      <p className="text-sm text-gray-600 mb-1">Memproses pesanan Anda...</p>
      <p className="text-xs text-gray-400">Status: {status}</p>
    </div>
  );
}