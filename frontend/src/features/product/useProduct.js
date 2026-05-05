import { useCallback, useEffect, useState } from "react";
import { useModal } from "../../shared/modal/ModalContext";
import ApiProduct from "./apiProduct ";
import ApiType from "../type/apiType";

export function useProduct() {
  const {
    confirm: modalConfirm,
    success: modalSuccess,
    error: modalError,
    loading: modalLoading,
  } = useModal();

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [products, setProducts] = useState([]);
  const [types, setTypes] = useState([]);
  const [paginate, setPaginate] = useState(null);
  const [page, setPage] = useState(1);
  const [hasNext, setHasNext] = useState(false);
  const [hasPrev, setHasPrev] = useState(false);
  const [searchPara, setSearchPara] = useState("");

  const fetchTypes = async () => {
    try {
      const res = await ApiType.getAll();
      setTypes(res.data.data || []);
    } catch (err) {
    }
  };

  // files: array of File object
  // return: array of { upload_url, object_name }
  const generatePresignedURLs = async (files) => {
    const payload = {
      files: files.map((f) => ({
        file_name: f.name,
      })),
    };
    
    var res;

    try {
      res = await ApiProduct.generatePresignedURLs(payload);
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
    }
    return res.data.uploads;
  };

  const uploadToPresignedURL = async (uploadURL, file) => {
    await ApiProduct.uploadToPresignedURL(uploadURL, file);
  };

  const applyPaginate = (paginateData) => {
    setPaginate(paginateData);
    setHasNext(paginateData?.has_next === "true");
    setHasPrev(paginateData?.has_prev === "true");
  };

  const fetchFirstPage = async (search = searchPara) => {
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllPaginate("", "", "next", search);
      setProducts(res.data.data || []);
      console.log(res.data)

      if (res.data.data?.length > 0) {
        applyPaginate(res.data.paginate);
        setPage(1);
        setHasNext(res.data.paginate?.has_next === "true");
        setHasPrev(false);
      } else {
        if (search) await modalError(`Data "${search}" tidak ditemukan`);
        setHasNext(false);
        setHasPrev(false);
      }
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
    } finally {
      setLoading(false);
    }
  };

  const next = async () => {
    if (paginate?.has_next !== "true") return;

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllPaginate(
        paginate.last_id,
        paginate.last_created_at,
        "next",
        searchPara
      );
      setProducts(res.data.data || []);
      applyPaginate(res.data.paginate);
      setPage((p) => p + 1);
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const prev = async () => {
    if (paginate?.has_prev !== "true") return;

    const closeLoading = modalLoading("Loading...");
    setLoading(true);
    setError(null);
    try {
      const res = await ApiProduct.getAllPaginate(
        paginate.first_id,
        paginate.first_created_at,
        "prev",
        searchPara
      );
      setProducts(res.data.data || []);
      applyPaginate(res.data.paginate);
      setPage((p) => p - 1);
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  // CRUD

  // imageItems: [{ file: File, objectName: string }] — sudah di-upload ke MinIO
  const create = async (formData, imageItems) => {
    if (imageItems.length === 0) {
      await modalError("Minimal harus ada 1 gambar");
      throw err;
    }

    setLoading(true);
    const closeLoading = modalLoading("Menyimpan produk...");
    setError(null);
    try {
      const payload = {
        product_code: formData.product_code,
        product_name: formData.product_name,
        weight_gram: Number(formData.weight_gram),
        type_id: formData.type_id,
        description: formData.description,
        images: imageItems.map((i) => i.objectName),
      };
      await ApiProduct.create(payload);
      await fetchFirstPage(searchPara);
      await modalSuccess("Produk berhasil dibuat");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const update = async (id, formData, imageItems) => {
    if (imageItems.length === 0) {
      await modalError("Minimal harus ada 1 gambar");
      throw err;
    }

    const confirmed = await modalConfirm("Apakah anda yakin mengubah produk?");
    if (!confirmed) return;

    setLoading(true);
    const closeLoading = modalLoading("Update...");
    setError(null);
    try {
      const payload = {
        product_code: formData.product_code,
        product_name: formData.product_name,
        weight_gram: Number(formData.weight_gram),
        type_id: formData.type_id,
        description: formData.description,
        images: imageItems.map((i) => i.objectName),
      };
      await ApiProduct.update(id, payload);
      await fetchFirstPage(searchPara);
      await modalSuccess("Produk berhasil diubah");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const del = async (id) => {
    const confirmed = await modalConfirm("Apakah anda yakin menghapus produk?");
    if (!confirmed) return;

    const closeLoading = modalLoading("Menghapus...");
    setLoading(true);
    setError(null);
    try {
      await ApiProduct.delete(id);
      await fetchFirstPage(searchPara);
      await modalSuccess("Produk berhasil dihapus");
    } catch (err) {
      await modalError(err.message || "Terjadi kesalahan");
      throw err;
    } finally {
      setLoading(false);
      closeLoading();
    }
  };

  const getById = async (id) => {
    const res = await ApiProduct.get(id);
    return res.data.data;
  };

  const setSearch = useCallback((value) => {
    setSearchPara(value);
    if (value.length >= 2 || value === "") {
      fetchFirstPage(value);
    }
  }, []);

  useEffect(() => {
    fetchFirstPage("");
    fetchTypes();
  }, []);

  return {
    products,
    types,
    loading,
    error,
    page,
    hasNext,
    hasPrev,
    setSearch,
    create,
    update,
    del,
    getById,
    next,
    prev,
    generatePresignedURLs,
    uploadToPresignedURL,
  };
}