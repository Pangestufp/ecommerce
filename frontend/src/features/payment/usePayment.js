import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";

import api from "../../api/axios";
import Endpoints from "../../shared/util/endpoint";
import { generateIdempotencyKey } from "../checkout/checkoutHelpers";

function loadSnapScript() {
  if (document.getElementById("midtrans-snap")) return;
  const script = document.createElement("script");
  script.id = "midtrans-snap";
  script.src =
    import.meta.env.VITE_MIDTRANS_SNAP_URL ??
    "https://app.sandbox.midtrans.com/snap/snap.js";
  script.setAttribute(
    "data-client-key",
    import.meta.env.VITE_MIDTRANS_CLIENT_KEY ?? ""
  );
  script.async = true;
  document.head.appendChild(script);
}

export function usePayment(salesOrderCode) {
  const navigate = useNavigate();

  const [order, setOrder] = useState(null);
  const [details, setDetails] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [payLoading, setPayLoading] = useState(false);

  const idempotencyKeyRef = useRef(generateIdempotencyKey());

  // Muat Snap script sekali saat hook pertama kali dipakai
  useEffect(() => {
    loadSnapScript();
  }, []);

  useEffect(() => {
    if (!salesOrderCode) return;

    const fetchOrder = async () => {
      setLoading(true);
      setError(null);
      try {
        const res = await api.get(
          Endpoints.SALESORDER_USER.GET_BY_CODE(salesOrderCode)
        );
        setOrder(res.data.data.sales_order);
        setDetails(res.data.data.details);
      } catch (err) {
        setError(err.message || "Gagal memuat detail pesanan");
      } finally {
        setLoading(false);
      }
    };

    fetchOrder();
  }, [salesOrderCode]);

  const handlePay = async () => {
    if (payLoading) return;
    setPayLoading(true);
    try {
      const res = await api.post(
        Endpoints.PAYMENT.CREATE_SNAP,
        { sales_order_code: salesOrderCode },
        { idempotencyKey: idempotencyKeyRef.current }
      );
      const { snap_token } = res.data.data;

      window.snap.pay(snap_token, {
        onSuccess: () => navigate(`/pesanan/${salesOrderCode}`),
        onPending: () => navigate(`/pesanan/${salesOrderCode}`),
        onError: () => navigate(`/pesanan/${salesOrderCode}`),
        onClose: () => setPayLoading(false),
      });
    } catch (err) {
      setPayLoading(false);
      setError(err.message || "Gagal memulai pembayaran");
    }
  };

  const isPending = order?.status === "PENDING_PAYMENT";

  return {
    order,
    details,
    loading,
    error,
    payLoading,
    isPending,
    handlePay,
  };
}