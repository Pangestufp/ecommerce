import { useState, useEffect } from "react";
import ApiProduct from "./apiProduct ";

export default function useBatchTransaction(batchID) {
  const [loading, setLoading] = useState(false);
  const [transactions, setTransactions] = useState([]);
  const [paginate, setPaginate] = useState(null);
  const [page, setPage] = useState(1);

  const fetchTransactions = async (trxId = "", direction = "next", createdAt = "") => {
    setLoading(true);
    try {
      const res = await ApiProduct.getTransactionsByBatchPaginate(batchID, trxId, direction, createdAt);
      setTransactions(res.data.data || []);
      setPaginate(res.data.paginate);
    } catch (err) {
      console.error("Gagal mengambil transaksi batch:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (batchID) {
      fetchTransactions();
      setPage(1); // Reset halaman jika batchID berubah
    }
  }, [batchID]);

  const handleNext = () => {
    if (paginate?.has_next === "true") {
      fetchTransactions(paginate.last_id, "next", paginate.last_created_at);
      setPage((prev) => prev + 1);
    }
  };

  const handlePrev = () => {
    if (paginate?.has_prev === "true") {
      fetchTransactions(paginate.first_id, "prev", paginate.first_created_at);
      setPage((prev) => prev - 1);
    }
  };

  return {
    transactions,
    loading,
    paginate,
    page,
    handleNext,
    handlePrev,
  };
}