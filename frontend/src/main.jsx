import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.jsx";
import { ModalProvider } from "./shared/modal/ModalContext.jsx";
import { BrowserRouter } from "react-router-dom";
import { CheckoutProvider } from "./features/checkout/CheckoutContext.jsx";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <BrowserRouter>
      <CheckoutProvider>
        <ModalProvider>
          <App />
        </ModalProvider>
      </CheckoutProvider>
    </BrowserRouter>
  </StrictMode>,
);
