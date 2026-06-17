import Button from "../../shared/ui/Button";
import { formatRupiah } from "./checkoutHelpers";

import { AddressSection } from "./components/AddressSection";
import ProductItem from "./components/Productitem ";
import OrderSummary from "./components/Ordersummary";
import CourierSection from "./components/CourierSection";
import NoteSection from "./components/NoteSection";

import { useCheckout } from "./useCheckout";

export default function CheckoutPage() {
  const {
    loading,
    checkoutData,

    products,
    addresses,

    productStates,

    selectedAddressId,
    setSelectedAddressId,

    courierData,
    selectedCourier,
    setSelectedCourier,

    loadingCourier,
    loadCourier,

    note,
    setNote,

    subtotal,
    grandTotal,

    handleDiscount,
    handleOrder,
  } = useCheckout();

  if (loading) {
    return <div className="p-4">Loading...</div>;
  }

  if (!checkoutData) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-xl mx-auto px-4 py-6 pb-32">

        <h1 className="text-lg font-semibold text-gray-900 mb-5">
          Konfirmasi Pesanan
        </h1>

        {checkoutData.is_note === 1 && (
          <div className="bg-orange-50 border border-orange-100 rounded-xl px-4 py-3 mb-4">
            <p className="text-xs text-orange-700">{checkoutData.note}</p>
          </div>
        )}

        <AddressSection
          addresses={addresses}
          selectedId={selectedAddressId}
          onSelect={setSelectedAddressId}
        />

        <CourierSection
          courierData={courierData}
          selectedCourier={selectedCourier}
          onSelect={setSelectedCourier}
          onLoad={loadCourier}
          loading={loadingCourier}
          disabled={!selectedAddressId}
        />

        <div className="bg-white rounded-2xl border border-gray-100 px-4 mb-3">
          {products.map((item, idx) => {
            const state = productStates[item.product_id];

            return (
              <ProductItem
                key={item.product_id}
                item={item}
                qty={state?.qty ?? item.qty}
                selectedDiscountId={state?.selectedDiscountId}
                onDiscountChange={(id) => handleDiscount(item.product_id, id)}
                isLast={idx === products.length - 1}
              />
            );
          })}
        </div>

        <NoteSection value={note} onChange={setNote} />

        <OrderSummary
          products={products}
          productStates={productStates}
          subtotal={subtotal}
          selectedCourier={selectedCourier}
          grandTotal={grandTotal}
        />
      </div>

      <div className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-100 px-4 py-3">
        <div className="max-w-xl mx-auto flex items-center justify-between gap-4">
          <div>
            <p className="text-xs text-gray-400">Total Pembayaran</p>
            <p className="text-base font-bold text-gray-900">{formatRupiah(grandTotal)}</p>
          </div>

          <Button
            onClick={handleOrder}
            disabled={!selectedAddressId || !selectedCourier}
            className="shrink-0 px-6 py-2.5 text-sm"
          >
            Buat Pesanan
          </Button>
        </div>
      </div>
    </div>
  );
}