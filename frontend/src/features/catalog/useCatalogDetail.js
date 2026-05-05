import { useEffect, useState } from "react";
import { useModal } from "../../shared/modal/ModalContext";
import ApiCatalog from "./apiCatalog ";

export function useCatalogDetail(slug) {
  const { error: modalError, loading: modalLoading } = useModal();

  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!slug) return;

    const fetch = async () => {
      setLoading(true);
      const closeLoading = modalLoading("Memuat produk...");
      try {
        const res = await ApiCatalog.getBySlug(slug);
        setProduct(res.data.data);
        closeLoading();
      } catch (err) {
        closeLoading();
        await modalError(err.message || "Terjadi kesalahan");
      } finally {
        setLoading(false);
      }
    };

    fetch();
  }, [slug]);

  return { product, loading };
}