import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";

import ApiCheckout from "./apiCheckout";
import { getBestDiscount, getUnitPrice } from "./checkoutHelpers";

export function useCheckout() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [loading, setLoading] = useState(true);
  const [checkoutData, setCheckoutData] = useState(null);

  const [productStates, setProductStates] = useState({});

  const [selectedAddressId, setSelectedAddressId] = useState(null);

  const [courierData, setCourierData] = useState(null);
  const [selectedCourier, setSelectedCourier] = useState(null);
  const [loadingCourier, setLoadingCourier] = useState(false);

  const [note, setNote] = useState("");

  useEffect(() => {
    const fetchCheckout = async () => {
      try {
        const res = await ApiCheckout.getCheckout(id);
        const data = res.data.data;

        setCheckoutData(data);

        const state = {};
        data.product_price?.forEach((item) => {
          const best = getBestDiscount(item.discount);
          state[item.product_id] = {
            qty: item.qty,
            selectedDiscountId: best?.discount_id ?? null,
          };
        });
        setProductStates(state);

        const primary = data.user_address?.find((a) => a.is_primary === 1);
        setSelectedAddressId(
          primary?.address_id ?? data.user_address?.[0]?.address_id ?? null,
        );
      } catch (err) {
        navigate("/cart");
      } finally {
        setLoading(false);
      }
    };

    fetchCheckout();
  }, [id]);

  useEffect(() => {
    setCourierData(null);
    setSelectedCourier(null);
  }, [selectedAddressId]);

  const products = checkoutData?.product_price ?? [];
  const addresses = checkoutData?.user_address ?? [];

  const loadCourier = async () => {
    if (!selectedAddressId) return;

    try {
      setLoadingCourier(true);

      const res = await ApiCheckout.getCourier({
        checkout_id: id,
        address_id: selectedAddressId,
      });

      setCourierData(res.data.data ?? null);
    } catch (err) {
      console.error(err);
    } finally {
      setLoadingCourier(false);
    }
  };

  const handleDiscount = (productId, discountId) => {
    setProductStates((prev) => ({
      ...prev,
      [productId]: {
        ...prev[productId],
        selectedDiscountId: discountId,
      },
    }));
  };

  const subtotal = useMemo(() => {
    return products.reduce((sum, item) => {
      const state = productStates[item.product_id];
      if (!state) return sum;

      const disc =
        item.discount?.find((d) => d.discount_id === state.selectedDiscountId) ?? null;

      return sum + getUnitPrice(item, disc) * state.qty;
    }, 0);
  }, [products, productStates]);

  const grandTotal = useMemo(() => {
    return subtotal + (selectedCourier?.cost ?? 0);
  }, [subtotal, selectedCourier]);

  const handleOrder = () => {
    const payload = {
      checkout_id: id,
      address_id: selectedAddressId,
      courier_code: selectedCourier?.code,
      courier_service: selectedCourier?.service,
      courier_name: selectedCourier?.display_name,
      shipping_cost: selectedCourier?.cost,
      note: note.trim() || null,
      items: products.map((item) => {
        const state = productStates[item.product_id];
        const disc =
          item.discount?.find((d) => d.discount_id === state?.selectedDiscountId) ?? null;

        return {
          product_id: item.product_id,
          qty: state?.qty ?? item.qty,
          discount_id: state?.selectedDiscountId ?? null,
          unit_price: getUnitPrice(item, disc),
        };
      }),
      subtotal,
      total: grandTotal,
    };

    console.log(payload);
  };

  return {
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
  };
}