export default function ProductCard({ product, onClick }) {
  const primaryImg =
    product.Images?.find((i) => i.is_primary === 1) || product.Images?.[0];
  const hasDiscount = product.BestDiscount > 0;

  return (
    <div
      onClick={onClick}
      className="bg-white border border-gray-100 rounded-xl overflow-hidden cursor-pointer hover:border-gray-300 transition-colors"
    >
      {primaryImg ? (
        <img
          src={primaryImg.picture_path}
          alt={product.ProductName}
          className="w-full aspect-square object-cover"
          loading="lazy"
        />
      ) : (
        <div className="w-full aspect-square bg-gray-100 flex items-center justify-center">
          <span className="text-gray-400 text-xs">No Image</span>
        </div>
      )}

      <div className="p-3">
        <p className="text-sm font-medium text-gray-900 truncate mb-1">
          {product.ProductName}
        </p>

        <div className="flex items-center gap-1.5 flex-wrap">
          <span className="text-sm font-medium text-gray-900">
            {product.BestPriceFormat}
          </span>
          {hasDiscount && (
            <span className="text-xs bg-orange-50 text-orange-700 rounded px-1.5 py-0.5">
              -{product.BestDiscountFormat}
            </span>
          )}
        </div>

        <p className="text-xs text-gray-400 mt-1.5">
          Stok: {product.AvailableStock}
        </p>
      </div>
    </div>
  );
}
