import { useCallback, useEffect, useState } from "react";
import { useModal } from "../../shared/modal/ModalContext";
import ApiSalesOrderCustomer from "./apiSalesOrderCustomer";

export function useUserSalesOrder() {
  const { error: modalError, loading: modalLoading } = useModal();

  const [loading, setLoading] = useState(false);
  const [orders, setOrders] = useState([]);
  const [paginate, setPaginate] = useState(null);
  const [page, setPage] = useState(1);
  const [hasNext, setHasNext] = useState(false);
  const [hasPrev, setHasPrev] = useState(false);

  const [selectedStatuses, setSelectedStatuses] = useState([]);

  const applyPaginate = (data, pag) => {
    setOrders(data);
    if (pag) {
      setPaginate(pag);
      setHasNext(pag.has_next === "true");
      setHasPrev(pag.has_prev === "true");
    } else {
      setHasNext(false);
      setHasPrev(false);
    }
  };

  const fetchFirstPage = useCallback(async (statuses = []) => {
    setLoading(true);
    try {
      const res = await ApiSalesOrderCustomer.getFirstPage(statuses);
      applyPaginate(res.data.data ?? [], res.data.paginate);
      setPage(1);
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
    } finally {
      setLoading(false);
    }
  }, []);

  const next = async () => {
    if (!hasNext || !paginate) return;
    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    try {
      const res = await ApiSalesOrderCustomer.getNextPage(paginate, selectedStatuses);
      applyPaginate(res.data.data ?? [], res.data.paginate);
      setPage((p) => p + 1);
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const prev = async () => {
    if (!hasPrev || !paginate) return;
    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    try {
      const res = await ApiSalesOrderCustomer.getPrevPage(paginate, selectedStatuses);
      applyPaginate(res.data.data ?? [], res.data.paginate);
      setPage((p) => p - 1);
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const toggleStatus = useCallback((status) => {
    setSelectedStatuses((prev) => {
      const next = prev.includes(status)
        ? prev.filter((s) => s !== status)
        : [...prev, status];
      fetchFirstPage(next);
      return next;
    });
  }, [fetchFirstPage]);

  const clearStatuses = useCallback(() => {
    setSelectedStatuses([]);
    fetchFirstPage([]);
  }, [fetchFirstPage]);

  useEffect(() => {
    fetchFirstPage([]);
  }, []);

  return {
    loading,
    orders,
    page,
    hasNext,
    hasPrev,
    selectedStatuses,
    toggleStatus,
    clearStatuses,
    next,
    prev,
  };
}