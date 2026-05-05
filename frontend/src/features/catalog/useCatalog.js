import { useEffect, useRef, useState } from "react";
import { useModal } from "../../shared/modal/ModalContext";
import ApiCatalog from "./apiCatalog ";

const LIMIT = 10;

export function useCatalog() {
  const { error: modalError, loading: modalLoading } = useModal();

  const [loading, setLoading] = useState(false);
  const [products, setProducts] = useState([]);
  const [searchPara, setSearchPara] = useState("");
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const searchRef = useRef("");

  const fetchAll = async (search, pageNum) => {
    if (loading) return;
    setLoading(true);
    const closeLoading = modalLoading("Memuat produk...");

    try {
      const res = await ApiCatalog.getBySearch(search, pageNum, LIMIT);
      const data = res.data.data || [];

      if (pageNum === 1) {
        setProducts(data);
      } else {
        setProducts(prev => [...prev, ...data]);
      }

      setHasMore(data.length === LIMIT);

      closeLoading();
    } catch (err) {
      closeLoading();
      await modalError(err.message || "Terjadi kesalahan");
    } finally {
      setLoading(false);
    }
  };

  const setSearch = (value) => {
    setSearchPara(value);
    searchRef.current = value;
    if (value.length >= 2 || value === "") {
      setPage(1);
      setHasMore(true);
      fetchAll(value, 1);
    }
  };

  const loadMore = () => {
    if (!hasMore || loading) return;
    const nextPage = page + 1;
    setPage(nextPage);
    fetchAll(searchRef.current, nextPage);
  };

  useEffect(() => {
    fetchAll("", 1);
  }, []);

  return { products, loading, hasMore, setSearch, loadMore };
}