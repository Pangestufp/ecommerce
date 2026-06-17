import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import { useModal } from "../../shared/modal/ModalContext";

import { getCart, setCart, clearCart } from "../../shared/util/cartStorage";

import ApiCart from "./apiCart";

export function useCart() {
  const navigate = useNavigate();

  const { error: modalError, loading: modalLoading } = useModal();
  const [items, setItems] = useState(() => getCart());
  const [verifiedItems, setVerifiedItems] = useState([]);
  const [checkedIDs, setCheckedIDs] = useState(new Set());
  const [note, setNote] = useState("");
  const [checkedTotal, setCheckedTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [isCheckingOut, setIsCheckingOut] = useState(false);

  useEffect(() => {
    const localCart = getCart();

    if (localCart.length === 0) return;

    verify(localCart);
  }, []);

  const recalcTotal = (verified, ids) => {
    const total = verified
      .filter((p) => ids.has(p.product_id))
      .reduce((sum, p) => sum + parseFloat(p.Price) * p.qty, 0);

    setCheckedTotal(total);
  };

  const verify = async (cartItems) => {
    setLoading(true);

    const closeLoading = modalLoading("Memverifikasi keranjang...");

    try {
      const res = await ApiCart.verifyCart(cartItems);

      const data = res.data.data;

      if (data.list_save?.length > 0) {
        setCart(data.list_save);

        setItems(getCart());
      } else {
        clearCart();

        setItems([]);
      }

      const verified = data.list_product ?? [];

      setVerifiedItems(verified);

      setNote(data.is_note === 1 ? data.note : "");

      const allIDs = new Set(verified.map((p) => p.product_id));

      setCheckedIDs(allIDs);

      recalcTotal(verified, allIDs);
    } catch (err) {
      await modalError(err.message ?? "Terjadi kesalahan");
    } finally {
      closeLoading();
      setLoading(false);
    }
  };

  const toggleCheck = (id) => {
    setCheckedIDs((prev) => {
      const next = new Set(prev);

      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }

      recalcTotal(verifiedItems, next);

      return next;
    });
  };

  const updateQty = (id, newQty) => {
    if (newQty <= 0) {
      removeItem(id);
      return;
    }

    const verifiedItem = verifiedItems.find((p) => p.product_id === id);

    if (verifiedItem && newQty > verifiedItem.available_stock) {
      return;
    }

    const updatedVerified = verifiedItems.map((p) =>
      p.product_id === id
        ? {
            ...p,
            qty: newQty,
          }
        : p,
    );

    setVerifiedItems(updatedVerified);

    recalcTotal(updatedVerified, checkedIDs);

    const updatedLocal = items.map((i) =>
      i.id === id
        ? {
            ...i,
            qty: newQty,
          }
        : i,
    );

    setCart(updatedLocal);

    setItems(updatedLocal);
  };

  const removeItem = (id) => {
    const updatedLocal = items.filter((i) => i.id !== id);

    setCart(updatedLocal);

    setItems(updatedLocal);

    const updatedVerified = verifiedItems.filter((p) => p.product_id !== id);

    setVerifiedItems(updatedVerified);

    setCheckedIDs((prev) => {
      const next = new Set(prev);

      next.delete(id);

      recalcTotal(updatedVerified, next);

      return next;
    });

    if (updatedLocal.length === 0) {
      setCheckedTotal(0);

      setNote("");
    }
  };

  const emptyCart = () => {
    clearCart();

    setItems([]);

    setVerifiedItems([]);

    setCheckedIDs(new Set());

    setCheckedTotal(0);

    setNote("");
  };

  const checkedItems = verifiedItems.filter((p) =>
    checkedIDs.has(p.product_id),
  );

  const checkout = async () => {
    if (checkedItems.length === 0) return;

    setIsCheckingOut(true);

    const closeLoading = modalLoading("Memproses checkout...");

    try {
      const res = await ApiCart.checkout(
        checkedItems.map((p) => ({
          id: p.product_id,
          qty: p.qty,
        })),
      );

      const checkoutId = res.data.data.checkout_id;

      navigate(`/checkout/${checkoutId}`);
    } catch (err) {
      await modalError(err);
    } finally {
      closeLoading();
      setIsCheckingOut(false);
    }
  };

  return {
    verifiedItems,
    checkedIDs,
    checkedItems,

    checkedTotal,
    note,

    loading,
    isCheckingOut,

    toggleCheck,
    updateQty,
    removeItem,
    emptyCart,

    checkout,
  };
}
