import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useCheckoutContext } from "./CheckoutContext";

export function useCheckout() {
  const { checkoutData, clearCheckout } = useCheckoutContext();
  const navigate = useNavigate();

  useEffect(() => {
    if (!checkoutData) {
      navigate("/cart");
    }
  }, [checkoutData]);

  return { checkoutData, clearCheckout };
}