import { useCallback, useEffect, useState } from "react";
import ApiType from "./apiType";
import { useModal } from "../../shared/modal/ModalContext";

export function useType() {
  const {
    confirm: modalConfirm,
    success: modalSuccess,
    error: modalError,
    loading: modalLoading,
  } = useModal();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [types, setTypes] = useState([]);
  const [paginate, setPaginate] = useState(null);
  const [page, setPage] = useState(1);
  const [hasNext, setHasNext] = useState(true);
  const [hasPrev, setHasPrev] = useState(false);
  const [searchPara, setSearchPara] = useState("");

  const create = async (payload) => {
    setLoading(true);
    const closeLoading = modalLoading("Create...");
    setError(null);

    try {
      const res = await ApiType.create(payload);
      setTypes((prev) => [res.data.data, ...prev]);
      fetchFirstPage(searchPara);
      await modalSuccess("Tipe berhasil dibuat");
    } catch (err) {
      await modalError(err.message);
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const update = async (id, payload) => {
    const confirmed = await modalConfirm("Apakah anda yakin mengubah Tipe?");
    if (!confirmed) return;

    setLoading(true);
    const closeLoading = modalLoading("Update...");
    setError(null);
    try {
      const res = await ApiType.update(id, payload);
      setTypes((prev) =>
        prev.map((t) => (t.type_id === id ? res.data.data : t)),
      );
      fetchFirstPage(searchPara)
      await modalSuccess("Tipe berhasil diubah");
    } catch (err) {
      await modalError(err.message);
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const del = async (id) => {
    const confirmed = await modalConfirm("Apakah anda yakin menghapus Tipe?");
    if (!confirmed) return;

    const closeLoading = modalLoading("Update...");
    setLoading(true);
    setError(null);
    try {
      await ApiType.delete(id);
      setTypes((prev) => prev.filter((t) => t.type_id !== id));
      await modalSuccess("Tipe berhasil dihapus");
    } catch (err) {
      await modalError(err.message);
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const next = async () => {
    if (paginate.has_next === "false") return;

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const res = await ApiType.getAllPaginate(
        paginate.last_id,
        paginate.last_created_at,
        "next",
        searchPara
      );
      setTypes(res.data.data);
      setPaginate(res.data.paginate);
      setPage(prev => prev + 1)

      setHasNext(res.data.paginate.has_next === "true")
      setHasPrev(res.data.paginate.has_prev === "true")

    } catch (err) {
      await modalError(err.message);
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const prev = async () => {
    if (paginate.has_prev === "false") return;

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const res = await ApiType.getAllPaginate(
        paginate.first_id,
        paginate.first_created_at,
        "prev",
        searchPara
      );
      setTypes(res.data.data);
      setPaginate(res.data.paginate);
      setPage(prev => prev - 1)
      
      setHasNext(res.data.paginate.has_next === "true")
      setHasPrev(res.data.paginate.has_prev === "true")

    } catch (err) {
      await modalError(err.message);
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };
  
  const fetchFirstPage = async (value) => {
    setLoading(true);
    setError(null);
    try {
      const res = await ApiType.getAllPaginate(
        "","","next", value
      );

      setTypes(res.data.data);

      if(res.data.data.length>0){

        setPaginate(res.data.paginate);
        setPage(1);

        setHasNext(res.data.paginate.has_next === "true")
        setHasPrev(res.data.paginate.has_prev === "true")
      }else{
        if (value !="") {
        await modalError("Data ("+value+") tidak ditemukan");
        }
        setHasNext(false)
        setHasPrev(false)
      }

    } catch (err) {
      await modalError(err.message);
    } finally {
      setLoading(false);
    }
  };


  const setSearch = useCallback((value) => {
    setSearchPara(value);
    if (value.length >= 2 || value === "") {
        fetchFirstPage(value);
    }
  }, []);


  useEffect(() => {
    fetchFirstPage("");
  }, []);

  return { create, next, prev, update, del, setSearch, hasNext, hasPrev, page, types, loading, error };
}
