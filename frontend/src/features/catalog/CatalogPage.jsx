import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { useCatalog } from "./useCatalog";
import LongSearchBar from "../../shared/ui/LongSearchBar";
import ProductCard from "./components/ProductCard";

export default function CatalogPage() {
  const { products, loading, hasMore, setSearch, loadMore } = useCatalog();
  const navigate = useNavigate();
  const sentinelRef = useRef(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) loadMore();
      },
      { threshold: 0.1 },
    );
    if (sentinelRef.current) observer.observe(sentinelRef.current);
    return () => observer.disconnect();
  }, [loadMore]);

  return (
    <div className="p-4 max-w-5xl mx-auto">
      <div className="mb-6 flex justify-center">
        <div className="w-full max-w-xl">
          <LongSearchBar
            placeholder="Cari produk, kategori, atau brand..."
            onChange={(value) => setSearch(value)}
          />
        </div>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
        {products.map((p) => (
          <ProductCard
            key={p.ProductID}
            product={p}
            onClick={() => navigate(`/products/${p.ProductSlug}`)}
          />
        ))}
      </div>

      {hasMore && (
        <div ref={sentinelRef} className="flex justify-center py-6">
          {loading && (
            <div className="w-5 h-5 border-2 border-gray-200 border-t-gray-500 rounded-full animate-spin" />
          )}
        </div>
      )}

      {!hasMore && products.length > 0 && (
        <p className="text-center text-sm text-gray-400 py-6">
          Semua produk sudah ditampilkan
        </p>
      )}

      {!loading && products.length === 0 && (
        <p className="text-center text-sm text-gray-400 py-10">
          Produk tidak ditemukan
        </p>
      )}
    </div>
  );
}
