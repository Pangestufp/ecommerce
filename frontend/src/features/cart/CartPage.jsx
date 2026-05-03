import { useNavigate } from "react-router-dom";
import { useCart } from "./useCart";
import { useCheckoutContext } from "../checkout/CheckoutContext";
import CartLineItem from "./components/CartLineItem";
import { ShoppingCart } from "lucide-react";
import ApiCart from "./apiCart";
import { useModal } from "../../shared/modal/ModalContext";

function formatRupiah(n) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency", currency: "IDR", maximumFractionDigits: 0,
  }).format(n);
}

export default function CartPage() {
  const {
    verifiedItems,
    checkedIDs,
    checkedItems,
    checkedTotal,
    note,
    loading,
    toggleCheck,
    updateQty,
    removeItem,
    emptyCart,
  } = useCart();

  const { setCheckout } = useCheckoutContext();
  const navigate = useNavigate();
  const { error: modalError, loading: modalLoading } = useModal();

  const handleCheckout = async () => {
    if (checkedItems.length === 0) return;

    const closeLoading = modalLoading("Memproses checkout...");
    try {
      const res = await ApiCart.checkout(
        checkedItems.map(p => ({ id: p.product_id, qty: p.qty }))
      );
      closeLoading();
      setCheckout(res.data.data);
      navigate("/checkout");
    } catch (err) {
      closeLoading();
      await modalError(err);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-xl mx-auto px-4 py-6">

        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-lg font-semibold text-gray-900">
            Keranjang
            {verifiedItems.length > 0 && (
              <span className="ml-2 text-sm font-normal text-gray-400">
                ({verifiedItems.length} produk)
              </span>
            )}
          </h1>
          {verifiedItems.length > 0 && (
            <button
              onClick={emptyCart}
              className="text-xs text-red-400 hover:text-red-600 transition-colors"
            >
              Hapus semua
            </button>
          )}
        </div>

        {/* Note dari backend */}
        {note && (
          <div className="bg-orange-50 border border-orange-100 rounded-xl px-4 py-3 mb-4">
            <p className="text-xs text-orange-700">{note}</p>
          </div>
        )}

        {/* Empty state */}
        {!loading && verifiedItems.length === 0 ? (
          <div className="bg-white rounded-2xl border border-gray-100 p-12 flex flex-col items-center gap-3">
            <ShoppingCart size={36} className="text-gray-200" />
            <p className="text-sm text-gray-400">Keranjang masih kosong</p>
            <button
              onClick={() => navigate("/products")}
              className="text-xs font-medium text-gray-800 underline underline-offset-2 mt-1"
            >
              Lihat produk
            </button>
          </div>
        ) : (
          <>
            {/* List item */}
            <div className="bg-white rounded-2xl border border-gray-100 px-4">
              {verifiedItems.map(item => (
                <CartLineItem
                  key={item.product_id}
                  item={item}
                  checked={checkedIDs.has(item.product_id)}
                  onCheck={toggleCheck}
                  onQtyChange={updateQty}
                  onRemove={removeItem}
                />
              ))}
            </div>

            {/* Total + checkout */}
            <div className="bg-white rounded-2xl border border-gray-100 px-4 py-4 mt-3">
              <div className="flex justify-between items-center mb-4">
                <span className="text-sm text-gray-500">
                  Total ({checkedIDs.size} dipilih)
                </span>
                <span className="text-base font-semibold text-gray-900">
                  {formatRupiah(checkedTotal)}
                </span>
              </div>
              <button
                onClick={handleCheckout}
                disabled={checkedItems.length === 0}
                className="w-full py-3 rounded-xl bg-gray-900 text-white text-sm font-medium hover:bg-gray-700 transition-colors disabled:bg-gray-200 disabled:text-gray-400 disabled:cursor-not-allowed"
              >
                {checkedItems.length === 0
                  ? "Pilih produk untuk checkout"
                  : `Checkout (${checkedItems.length} produk)`}
              </button>
            </div>
          </>
        )}

      </div>
    </div>
  );
}