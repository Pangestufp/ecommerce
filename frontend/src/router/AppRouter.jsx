import { Routes, Route, Navigate } from "react-router-dom";
import ProtectedRoute from "./ProtectedRoute";
import ProtectedAdminOwner from "./ProtectedAdminOwner";
import AppLayout from "../shared/layout/AppLayout";
import { getToken } from "../shared/util/token";

import LoginPage from "../features/login/LoginPage";
import TypePage from "../features/type/TypePage";
import ProductPage from "../features/product/ProductPage";
import ProductDetailPage from "../features/product/ProductDetailPage";
import CatalogPage from "../features/catalog/CatalogPage";
import CatalogDetailPage from "../features/catalog/CatalogDetailPage";
import CartPage from "../features/cart/CartPage";
import CheckoutPage from "../features/checkout/CheckoutPage";

function AppRouter() {
  const token = getToken();

  return (
    <Routes>
      {/* Login tidak pakai layout */}
      <Route
        path="/login"
        element={token ? <Navigate to="/products" replace /> : <LoginPage />}
      />

      {/* Admin routes */}
      <Route
        path="/admin/tipe"
        element={
          <ProtectedRoute>
            <ProtectedAdminOwner>
              <AppLayout>
                <TypePage />
              </AppLayout>
            </ProtectedAdminOwner>
          </ProtectedRoute>
        }
      />
      <Route
        path="/admin/produk"
        element={
          <ProtectedRoute>
            <ProtectedAdminOwner>
              <AppLayout>
                <ProductPage />
              </AppLayout>
            </ProtectedAdminOwner>
          </ProtectedRoute>
        }
      />
      <Route
        path="/admin/produk/:id"
        element={
          <ProtectedRoute>
            <ProtectedAdminOwner>
              <AppLayout>
                <ProductDetailPage />
              </AppLayout>
            </ProtectedAdminOwner>
          </ProtectedRoute>
        }
      />

      {/* User routes */}
      <Route
        path="/products"
        element={
          <ProtectedRoute>
            <AppLayout>
              <CatalogPage />
            </AppLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/products/:slug"
        element={
          <ProtectedRoute>
            <AppLayout>
              <CatalogDetailPage />
            </AppLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/keranjang"
        element={
          <ProtectedRoute>
            <AppLayout>
              <CartPage />
            </AppLayout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/checkout"
        element={
          <ProtectedRoute>
            <AppLayout>
              <CheckoutPage />
            </AppLayout>
          </ProtectedRoute>
        }
      />

      <Route path="*" element={<Navigate to="/products" replace />} />
    </Routes>
  );
}

export default AppRouter;