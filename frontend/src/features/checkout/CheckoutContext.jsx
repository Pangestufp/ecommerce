import { createContext, useContext, useState } from "react";

const CheckoutContext = createContext(null);

export function CheckoutProvider({ children }) {
  const [checkoutData, setCheckoutData] = useState(null);

  const setCheckout = (data) => {
    setCheckoutData(data);
  };

  const clearCheckout = () => {
    setCheckoutData(null);
  };

  return (
    <CheckoutContext.Provider value={{ checkoutData, setCheckout, clearCheckout }}>
      {children}
    </CheckoutContext.Provider>
  );
}

export function useCheckoutContext() {
  const ctx = useContext(CheckoutContext);
  if (!ctx) throw new Error("useCheckoutContext must be used within CheckoutProvider");
  return ctx;
}