import { Routes, Route, Navigate } from "react-router-dom";
import ProtectedRoute from "./ProtectedRoute";
import {getToken } from "../shared/util/token";
import LoginPage from "../features/login/LoginPage";
import TypePage from "../features/type/TypePage";
import ProductPage from "../features/product/ProductPage";
import ProductDetailPage from "../features/product/ProductDetailPage";
import CatalogPage from "../features/catalog/CatalogPage";
import CatalogDetailPage from "../features/catalog/CatalogDetailPage";
function AppRouter() {
  const token = getToken();

  return (
    <Routes>
      <Route
        path="/login"
        element={token ? <Navigate to="/" replace /> : <LoginPage />}
      />

      <Route
        path="/tipe"
        element={
          <ProtectedRoute>
            <TypePage/>
          </ProtectedRoute>
        }
      />

      <Route
        path="/produk"
        element={
          <ProtectedRoute>
            <ProductPage/>
          </ProtectedRoute>
        }
      />

      <Route
        path="/products"
        element={
          <ProtectedRoute>
            <CatalogPage/>
          </ProtectedRoute>
        }
      />

      <Route
          path="/produk/:id"
          element={
            <ProtectedRoute>
              <ProductDetailPage/>
            </ProtectedRoute>
          }
        />

        <Route
          path="/products/:slug"
          element={
            <ProtectedRoute>
              <CatalogDetailPage/>
            </ProtectedRoute>
          }
        />


      <Route path="*" element={<Navigate to="/dashboard" />} />
    </Routes>
  );
}

export default AppRouter;
