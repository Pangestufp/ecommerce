import { useEffect, useState } from "react";
import { useModal } from "../../shared/modal/ModalContext";
import ApiSalesOrderCustomer from "./apiSalesOrderCustomer";

export function useCustomerOrderDetail(code) {
  const { error: modalError } = useModal();

  const [loading, setLoading] = useState(false);
  const [detail, setDetail] = useState(null); // SalesOrderDetailResponse

  useEffect(() => {
    if (!code) return;
    const fetch = async () => {
      setLoading(true);
      try {
        const res = await ApiSalesOrderCustomer.getByCode(code);
        setDetail(res.data.data);
      } catch (err) {
        await modalError(err.message || "Terjadi kesalahan");
      } finally {
        setLoading(false);
      }
    };
    fetch();
  }, [code]);

  return { loading, detail };
}