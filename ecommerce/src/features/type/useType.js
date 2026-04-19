import { useEffect, useState } from "react";
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
  const [cursorHistory, setCursorHistory] = useState([]);
  const [nextCursor, setNextCursor] = useState({ id: "", createdAt: "" });

  const create = async (payload) => {
    setLoading(true);
    const closeLoading = modalLoading("Create...");
    setError(null);

    try {
      const res = await ApiType.create(payload);
      setTypes((prev) => [res.data.data, ...prev]);
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
    console.log("ini next");
    console.log(cursorHistory);

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);

    try {
      const cursorBefore = nextCursor;
      const res = await ApiType.getAll(
        nextCursor.id,
        nextCursor.createdAt,
      );

      if (res.data.data.length > 0) {
        setCursorHistory((prev) => [...prev, cursorBefore]);
        setNextCursor({
          id: res.data.paginate.last_id,
          createdAt: res.data.paginate.last_created_at,
        });
        setTypes(res.data.data);
      }
    } catch (err) {
      await modalError(err.message);
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const prev = async () => {
    if (cursorHistory.length === 0) return;
    console.log("ini prev");
    console.log(cursorHistory);

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const history = [...cursorHistory];
      history.pop();
      const fetchCursor =
        history.length > 0
          ? history[history.length - 1]
          : { id: "", createdAt: "" };

      const res = await ApiType.getAll(fetchCursor.id, fetchCursor.createdAt);

      if (res.data.data.length > 0) {
        setCursorHistory(history);
        setNextCursor({
          id: res.data.paginate.last_id,
          createdAt: res.data.paginate.last_created_at,
        });
        setTypes(res.data.data);
      }
    } catch (err) {
      await modalError(err.message);
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const fetch = async () => {
    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const res = await ApiType.getAll("", "");
      if (res.data.data.length > 0) {
        setNextCursor({
          id: res.data.paginate.last_id,
          createdAt: res.data.paginate.last_created_at,
        });
        setTypes(res.data.data);
      }
    } catch (err) {
      await modalError(err.message);
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  useEffect(() => {
    fetch();
  }, []);

  return { create, next, prev, update, del, cursorHistory, types, loading, error };
}
