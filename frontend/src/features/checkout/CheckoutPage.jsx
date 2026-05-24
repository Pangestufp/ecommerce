import { useState, useMemo } from "react";
import { useCheckout } from "./useCheckout";
import Button from "../../shared/ui/Button"
import { getBestDiscount, getUnitPrice, formatRupiah } from "./checkoutHelpers";
import { AddressSection } from "./components/AddressSection";
import ProductItem from "./components/Productitem ";
import OrderSummary from "./components/Ordersummary";

export default function CheckoutPage() {
  const { checkoutData } = useCheckout();

  // ── State produk: qty + diskon terpilih ───────────────────────────────────
  const [productStates, setProductStates] = useState(() => {
    if (!checkoutData) return {};
    const init = {};
    checkoutData.product_price?.forEach((item) => {
      const best = getBestDiscount(item.discount);
      init[item.product_id] = {
        qty: item.qty,
        selectedDiscountId: best ? best.discount_id : null,
      };
    });
    return init;
  });

  // ── State alamat: primary atau index pertama ──────────────────────────────
  const [selectedAddressId, setSelectedAddressId] = useState(() => {
    if (!checkoutData) return null;
    const primary = checkoutData.user_address?.find((a) => a.is_primary === 1);
    return primary?.address_id ?? checkoutData.user_address?.[0]?.address_id ?? null;
  });

  if (!checkoutData) return null;

  const products = checkoutData.product_price ?? [];
  const addresses = checkoutData.user_address ?? [];

  // ── Grand total ───────────────────────────────────────────────────────────
  const grandTotal = useMemo(() => {
    return products.reduce((sum, item) => {
      const state = productStates[item.product_id];
      if (!state) return sum;
      const disc = item.discount?.find((d) => d.discount_id === state.selectedDiscountId) ?? null;
      return sum + getUnitPrice(item, disc) * state.qty;
    }, 0);
  }, [productStates, products]);

  // ── Handlers ─────────────────────────────────────────────────────────────
  const handleQty = (productId, delta, maxStock) => {
    setProductStates((prev) => {
      const next = prev[productId].qty + delta;
      return {
        ...prev,
        [productId]: {
          ...prev[productId],
          qty: Math.min(Math.max(1, next), maxStock),
        },
      };
    });
  };

  const handleDiscount = (productId, discountId) => {
    setProductStates((prev) => ({
      ...prev,
      [productId]: { ...prev[productId], selectedDiscountId: discountId },
    }));
  };

  const handleOrder = () => {
    const payload = {
      address_id: selectedAddressId,
      items: products.map((item) => {
        const state = productStates[item.product_id];
        const disc = item.discount?.find((d) => d.discount_id === state?.selectedDiscountId) ?? null;
        return {
          product_id: item.product_id,
          qty: state?.qty ?? item.qty,
          discount_id: state?.selectedDiscountId ?? null,
          unit_price: getUnitPrice(item, disc),
        };
      }),
      total: grandTotal,
    };
    console.log("Order payload:", payload);
    // TODO: call API order
  };

  // ─────────────────────────────────────────────────────────────────────────

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-xl mx-auto px-4 py-6 pb-32">

        <h1 className="text-lg font-semibold text-gray-900 mb-5">
          Konfirmasi Pesanan
        </h1>

        {/* Notifikasi dari backend */}
        {checkoutData.is_note === 1 && (
          <div className="bg-orange-50 border border-orange-100 rounded-xl px-4 py-3 mb-4">
            <p className="text-xs text-orange-700">{checkoutData.note}</p>
          </div>
        )}

        {/* Alamat pengiriman */}
        <AddressSection
          addresses={addresses}
          selectedId={selectedAddressId}
          onSelect={setSelectedAddressId}
        />

        {/* List produk */}
        <div className="bg-white rounded-2xl border border-gray-100 px-4 mb-3">
          {products.map((item, idx) => {
            const state = productStates[item.product_id] ?? {
              qty: item.qty,
              selectedDiscountId: null,
            };
            return (
              <ProductItem
                key={item.product_id}
                item={item}
                qty={state.qty}
                selectedDiscountId={state.selectedDiscountId}
                onQtyChange={(delta) => handleQty(item.product_id, delta, item.available_stock)}
                onDiscountChange={(id) => handleDiscount(item.product_id, id)}
                isLast={idx === products.length - 1}
              />
            );
          })}
        </div>

        {/* Ringkasan harga */}
        <OrderSummary
          products={products}
          productStates={productStates}
          grandTotal={grandTotal}
        />

      </div>

      {/* Sticky footer */}
      <div className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-100 px-4 py-3">
        <div className="max-w-xl mx-auto flex items-center justify-between gap-4">
          <div>
            <p className="text-xs text-gray-400">Total Pembayaran</p>
            <p className="text-base font-bold text-gray-900">{formatRupiah(grandTotal)}</p>
          </div>
          <Button
            onClick={handleOrder}
            disabled={!selectedAddressId}
            className="shrink-0 px-6 py-2.5 text-sm"
          >
            Buat Pesanan
          </Button>
        </div>
      </div>
    </div>
  );
}