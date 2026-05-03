import { useState } from "react";
import { useParams } from "react-router-dom";
import { useCatalogDetail } from "./useCatalogDetail";
import { Package, ShoppingCart, Check } from "lucide-react";
import { saveToCart } from "../../shared/util/cartStorage";

export default function CatalogDetailPage() {
  const { slug } = useParams();
  const { product, loading } = useCatalogDetail(slug);
  const [activeImg, setActiveImg] = useState(0);
  const [showMore, setShowMore] = useState(false);
  const [toast, setToast] = useState(false);

  const handleAddToCart = () => {
    const primaryImg = sortedImages[0]?.picture_path ?? null;

    saveToCart({
      id: product.product_id,
      productName: product.product_name,
      qty: 1,
      availableStock: product.available_stock,
      image: primaryImg,
    });

    setToast(true);
    setTimeout(() => setToast(false), 2000);
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="w-6 h-6 border-2 border-gray-200 border-t-gray-500 rounded-full animate-spin" />
      </div>
    );
  }

  if (!product) return null;

  const hasDiscount = product.best_discount > 0;
  const images = product.images ?? [];
  const sortedImages = [
    ...images.filter(i => i.is_primary === 1),
    ...images.filter(i => i.is_primary !== 1),
  ];

  return (
    <div className="min-h-screen bg-gray-50">

      <div className={`fixed top-4 right-4 z-50 flex items-center gap-2 bg-gray-900 text-white text-sm px-4 py-2.5 rounded-xl shadow-lg transition-all duration-300 ${
        toast ? "opacity-100 translate-y-0" : "opacity-0 -translate-y-2 pointer-events-none"
      }`}>
        <Check size={14} className="text-green-400" />
        Ditambahkan ke keranjang
      </div>

      <div className="max-w-4xl mx-auto px-4 py-6">
        <div className="bg-white rounded-2xl border border-gray-100 overflow-hidden">
          <div className="grid grid-cols-1 md:grid-cols-2">

            {/* Image section */}
            <div className="p-4 border-b md:border-b-0 md:border-r border-gray-100">
              <div className="aspect-square rounded-xl overflow-hidden bg-gray-50 mb-3">
                {sortedImages[activeImg] ? (
                  <img
                    src={sortedImages[activeImg].picture_path}
                    alt={product.product_name}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-gray-300">
                    <Package size={48} />
                  </div>
                )}
              </div>

              {sortedImages.length > 1 && (
                <div className="flex gap-2 overflow-x-auto pb-1">
                  {sortedImages.map((img, idx) => (
                    <button
                      key={img.image_id}
                      onClick={() => setActiveImg(idx)}
                      className={`shrink-0 w-16 h-16 rounded-lg overflow-hidden border-2 transition-colors ${
                        activeImg === idx
                          ? "border-gray-800"
                          : "border-transparent hover:border-gray-300"
                      }`}
                    >
                      <img
                        src={img.picture_path}
                        alt={`${product.product_name} ${idx + 1}`}
                        className="w-full h-full object-cover"
                      />
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Detail section */}
            <div className="p-4 md:p-6 flex flex-col gap-4">

              {/* Name + stock */}
              <div>
                <h1 className="text-xl font-semibold text-gray-900">
                  {product.product_name}
                </h1>
                <div className="mt-2">
                  <span className={`text-xs font-medium px-2.5 py-1 rounded-full ${
                    product.available_stock > 10
                      ? "bg-green-50 text-green-700"
                      : product.available_stock > 0
                      ? "bg-orange-50 text-orange-700"
                      : "bg-red-50 text-red-700"
                  }`}>
                    Stok: {product.available_stock}
                  </span>
                </div>
              </div>

              {/* Price */}
              <div>
                {hasDiscount && (
                  <p className="text-sm text-gray-400 line-through mb-0.5">
                    {product.product_price_format}
                  </p>
                )}
                <div className="flex items-center gap-2 flex-wrap">
                  <span className="text-2xl font-semibold text-gray-900">
                    {product.best_price_format}
                  </span>
                  {hasDiscount && (
                    <span className="text-xs bg-orange-50 text-orange-700 rounded-full px-2.5 py-1 font-medium">
                      Hemat {product.best_discount_format}
                    </span>
                  )}
                </div>
                {hasDiscount && (
                  <p className="text-xs text-gray-400 mt-1">{product.discount_name}</p>
                )}
              </div>

              {/* Tombol keranjang */}
              <button
                onClick={handleAddToCart}
                disabled={product.available_stock === 0}
                className="flex items-center justify-center gap-2 w-full py-2.5 rounded-xl bg-gray-900 text-white text-sm font-medium hover:bg-gray-700 transition-colors disabled:bg-gray-200 disabled:text-gray-400 disabled:cursor-not-allowed"
              >
                <ShoppingCart size={16} />
                {product.available_stock === 0 ? "Stok Habis" : "Masukkan Keranjang"}
              </button>

              <div className="border-t border-gray-100" />

              {/* Detail Produk */}
              <div>
                <p className="text-xs font-medium text-gray-400 uppercase tracking-wider mb-3">
                  Detail Produk
                </p>

                <div className="bg-gray-50 rounded-xl p-3 flex flex-wrap gap-x-6 gap-y-2 mb-3">
                  <div>
                    <p className="text-xs text-gray-400">Kode</p>
                    <p className="text-sm font-medium text-gray-800">{product.product_code}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-400">Tipe</p>
                    <p className="text-sm font-medium text-gray-800">{product.type_name}</p>
                  </div>
                  <div>
                    <p className="text-xs text-gray-400">Berat</p>
                    <p className="text-sm font-medium text-gray-800">{product.weight_gram}g</p>
                  </div>
                </div>

                <p className={`text-sm text-gray-600 leading-relaxed text-left ${!showMore ? "line-clamp-2" : ""}`}>
                  {product.description}
                </p>
                <button
                  onClick={() => setShowMore(prev => !prev)}
                  className="text-xs font-medium text-gray-800 underline underline-offset-2 mt-2 text-left"
                >
                  {showMore ? "Sembunyikan" : "Selengkapnya"}
                </button>
              </div>

            </div>
          </div>
        </div>
      </div>
    </div>
  );
}