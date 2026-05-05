import { useCallback, useEffect, useState } from "react";
import { useModal } from "../../shared/modal/ModalContext";
import ApiProduct from "./apiProduct ";

export function useProductDetail(id) {
  const {
    confirm: modalConfirm,
    success: modalSuccess,
    error: modalError,
    loading: modalLoading,
  } = useModal();

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const [prices, setPrices] = useState([]);
  const [pricePaginate, setPricePaginate] = useState(null);
  const [pricePage, setPricePage] = useState(1);
  const [hasNextPrice, setHasNextPrice] = useState(true);
  const [hasPrevPrice, setHasPrevPrice] = useState(false);

  const [discounts, setDiscounts] = useState([]);
  const [DiscountPaginate, setDiscountPaginate] = useState(null);
  const [pageDiscount, setPageDiscount] = useState(1);
  const [hasNextDiscount, setHasNextDiscount] = useState(true);
  const [hasPrevDiscount, setHasPrevDiscount] = useState(false);
  const [searchParaDiscount, setSearchParaDiscount] = useState("");

  const [inventories, setInventories] = useState([]);
  const [InventoryPaginate, setInventoryPaginate] = useState(null);
  const [pageInventory, setPageInventory] = useState(1);
  const [hasNextInventory, setHasNextInventory] = useState(true);
  const [hasPrevInventory, setHasPrevInventory] = useState(false);
  const [searchParaInventory, setSearchParaInventory] = useState("");

  const [types, setTypes] = useState([]);

  const [product, setProduct] = useState(null);

  const fetchProduct = async () => {
    try {
      const res = await ApiProduct.get(id);
      setProduct(res.data.data);
      console.log(res.data.data)
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
    }
  };

  const fetchTypes = async () => {
    try {
      const res = await ApiProduct.getAllDiscountType();
      setTypes(res.data.data || []);
    } catch (err) {}
  };

  const createPrice = async (payload) => {
    setLoading(true);
    const closeLoading = modalLoading("Create...");
    setError(null);

    try {
      const res = await ApiProduct.createPrice(payload);
      setPrices((prev) => [res.data.data, ...prev]);
      fetchFirstPricePage();
      await modalSuccess("Harga berhasil dibuat");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const nextPrice = async () => {
    if (pricePaginate.has_next === "false") return;

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllPricePaginate(
        id,
        pricePaginate.last_id,
        pricePaginate.last_created_at,
        "next",
      );
      setPrices(res.data.data);
      setPricePaginate(res.data.paginate);
      setPricePage((prev) => prev + 1);

      setHasNextPrice(res.data.paginate.has_next === "true");
      setHasPrevPrice(res.data.paginate.has_prev === "true");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const prevPrice = async () => {
    if (pricePaginate.has_prev === "false") return;

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllPricePaginate(
        id,
        pricePaginate.first_id,
        pricePaginate.first_created_at,
        "prev",
      );
      setPrices(res.data.data);
      setPricePaginate(res.data.paginate);
      setPricePage((prev) => prev - 1);

      setHasNextPrice(res.data.paginate.has_next === "true");
      setHasPrevPrice(res.data.paginate.has_prev === "true");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const fetchFirstPricePage = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllPricePaginate(id, "", "", "next");

      setPrices(res.data.data);

      if (res.data.data.length > 0) {
        setPricePaginate(res.data.paginate);
        setPricePage(1);

        setHasNextPrice(res.data.paginate.has_next === "true");
        setHasPrevPrice(res.data.paginate.has_prev === "true");
      } else {
        setHasNextPrice(false);
        setHasPrevPrice(false);
      }
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
    } finally {
      setLoading(false);
    }
  };

  const createDiscount = async (payload) => {
    setLoading(true);
    const closeLoading = modalLoading("Create...");
    setError(null);

    try {
      const res = await ApiProduct.createDiscount(payload);
      setDiscounts((prev) => [res.data.data, ...prev]);
      fetchFirstDiscountPage(searchParaDiscount);
      await modalSuccess("Diskon berhasil dibuat");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const delDiscount = async (id) => {
    const confirmed = await modalConfirm("Apakah anda yakin menghapus Diskon?");
    if (!confirmed) return;

    const closeLoading = modalLoading("Update...");
    setLoading(true);
    setError(null);
    try {
      await ApiProduct.deleteDiscount(id);
      setDiscounts((prev) => prev.filter((t) => t.discount_id !== id));
      await modalSuccess("Tipe berhasil dihapus");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const nextDiscount = async () => {
    if (DiscountPaginate.has_next === "false") return;

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllDiscountPaginate(
        id,
        DiscountPaginate.last_id,
        DiscountPaginate.last_created_at,
        "next",
        searchParaDiscount,
      );
      setDiscounts(res.data.data);
      setDiscountPaginate(res.data.paginate);
      setPageDiscount((prev) => prev + 1);

      setHasNextDiscount(res.data.paginate.has_next === "true");
      setHasPrevDiscount(res.data.paginate.has_prev === "true");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const prevDiscount = async () => {
    if (DiscountPaginate.has_prev === "false") return;

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllDiscountPaginate(
        id,
        DiscountPaginate.first_id,
        DiscountPaginate.first_created_at,
        "prev",
        searchParaDiscount,
      );
      setDiscounts(res.data.data);
      setDiscountPaginate(res.data.paginate);
      setPageDiscount((prev) => prev - 1);

      setHasNextDiscount(res.data.paginate.has_next === "true");
      setHasPrevDiscount(res.data.paginate.has_prev === "true");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const fetchFirstDiscountPage = async (value) => {
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllDiscountPaginate(
        id,
        "",
        "",
        "next",
        value,
      );

      setDiscounts(res.data.data);

      if (res.data.data.length > 0) {
        setDiscountPaginate(res.data.paginate);
        setPageDiscount(1);

        setHasNextDiscount(res.data.paginate.has_next === "true");
        setHasPrevDiscount(res.data.paginate.has_prev === "true");
      } else {
        setHasNextDiscount(false);
        setHasPrevDiscount(false);
      }
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
    } finally {
      setLoading(false);
    }
  };

  const setSearchDiscount = useCallback((value) => {
    setSearchParaDiscount(value);
    if (value.length >= 2 || value === "") {
      fetchFirstDiscountPage(value);
    }
  }, []);

  const createInventory = async (payload) => {
    setLoading(true);
    const closeLoading = modalLoading("Create...");
    setError(null);

    try {
      const res = await ApiProduct.createInventory(payload);
      setInventories((prev) => [res.data.data, ...prev]);
      fetchFirstInventoryPage(searchParaInventory);
      await modalSuccess("Inventory berhasil dibuat");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const updateInventory = async (id, payload) => {
    const confirmed = await modalConfirm(
      "Apakah anda yakin mengubah Batch ini?",
    );
    if (!confirmed) return;

    setLoading(true);
    const closeLoading = modalLoading("Update...");
    setError(null);
    try {
      const res = await ApiProduct.updateInventory(id, payload);
      setInventories((prev) =>
        prev.map((t) => (t.batch_id === id ? res.data.data : t)),
      );
      fetchFirstInventoryPage(searchParaInventory);
      await modalSuccess("Inventory berhasil diubah");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const nextInventory = async () => {
    if (InventoryPaginate.has_next === "false") return;

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllInventoryPaginate(
        id,
        InventoryPaginate.last_id,
        InventoryPaginate.last_created_at,
        "next",
        searchParaInventory,
      );
      setInventories(res.data.data);
      setInventoryPaginate(res.data.paginate);
      setPageInventory((prev) => prev + 1);

      setHasNextInventory(res.data.paginate.has_next === "true");
      setHasPrevInventory(res.data.paginate.has_prev === "true");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const prevInventory = async () => {
    if (InventoryPaginate.has_prev === "false") return;

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllInventoryPaginate(
        id,
        InventoryPaginate.first_id,
        InventoryPaginate.first_created_at,
        "prev",
        searchParaInventory,
      );
      setInventories(res.data.data);
      setInventoryPaginate(res.data.paginate);
      setPageInventory((prev) => prev - 1);

      setHasNextInventory(res.data.paginate.has_next === "true");
      setHasPrevInventory(res.data.paginate.has_prev === "true");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const fetchFirstInventoryPage = async (value) => {
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllInventoryPaginate(
        id,
        "",
        "",
        "next",
        value,
      );

      setInventories(res.data.data);

      if (res.data.data.length > 0) {
        setInventoryPaginate(res.data.paginate);
        setPageInventory(1);

        setHasNextInventory(res.data.paginate.has_next === "true");
        setHasPrevInventory(res.data.paginate.has_prev === "true");
      } else {
        setHasNextInventory(false);
        setHasPrevInventory(false);
      }
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
    } finally {
      setLoading(false);
    }
  };

  const setSearchInventory = useCallback((value) => {
    setSearchParaInventory(value);
    if (value.length >= 2 || value === "") {
      fetchFirstInventoryPage(value);
    }
  }, []);

  useEffect(() => {
    fetchProduct();
    fetchFirstPricePage("");
    fetchFirstInventoryPage("");
    fetchFirstDiscountPage("");
    fetchTypes();
  }, []);

  return {
    product,
    types,
    loading,
    error,
    prices,
    pricePage,
    hasNextPrice,
    hasPrevPrice,
    createPrice,
    nextPrice,
    prevPrice,
    discounts,
    pageDiscount,
    hasNextDiscount,
    hasPrevDiscount,
    createDiscount,
    delDiscount,
    nextDiscount,
    prevDiscount,
    setSearchDiscount,
    inventories,
    pageInventory,
    hasNextInventory,
    hasPrevInventory,
    createInventory,
    updateInventory,
    nextInventory,
    prevInventory,
    setSearchInventory,
  };
}